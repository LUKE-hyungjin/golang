package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// ============================================================================
// 모델 정의
// ============================================================================

type Account struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Number    string         `gorm:"uniqueIndex;not null" json:"number"`
	Name      string         `json:"name"`
	Balance   float64        `json:"balance"`
	Currency  string         `json:"currency"`
	IsLocked  bool           `gorm:"default:false" json:"is_locked"`
	Version   int            `gorm:"default:0" json:"version"` // 낙관적 잠금용
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transaction struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	TransactionID   string    `gorm:"uniqueIndex;not null" json:"transaction_id"`
	FromAccountID   uint      `json:"from_account_id"`
	FromAccount     Account   `gorm:"foreignKey:FromAccountID" json:"from_account,omitempty"`
	ToAccountID     uint      `json:"to_account_id"`
	ToAccount       Account   `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Type            string    `json:"type"` // transfer, deposit, withdrawal
	Status          string    `json:"status"` // pending, completed, failed, timeout
	Description     string    `json:"description"`
	ErrorMessage    string    `json:"error_message,omitempty"`
	ProcessingTime  int64     `json:"processing_time_ms"` // 밀리초
	CreatedAt       time.Time `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at"`
}

type Order struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	OrderNumber   string         `gorm:"uniqueIndex;not null" json:"order_number"`
	CustomerID    uint           `json:"customer_id"`
	TotalAmount   float64        `json:"total_amount"`
	Status        string         `json:"status"` // pending, processing, completed, cancelled
	Items         []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
	PaymentID     *uint          `json:"payment_id"`
	Payment       *Payment       `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID        uint    `gorm:"primarykey" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Product struct {
	ID        uint    `gorm:"primarykey" json:"id"`
	Name      string  `json:"name"`
	SKU       string  `gorm:"uniqueIndex;not null" json:"sku"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	Reserved  int     `json:"reserved"` // 예약된 재고
	Version   int     `gorm:"default:0" json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Payment struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	PaymentID     string    `gorm:"uniqueIndex;not null" json:"payment_id"`
	OrderID       uint      `json:"order_id"`
	Amount        float64   `json:"amount"`
	Method        string    `json:"method"` // card, bank, wallet
	Status        string    `json:"status"` // pending, processing, completed, failed
	TransactionID *string   `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	ProcessedAt   *time.Time `json:"processed_at"`
}

// ============================================================================
// 트랜잭션 서비스
// ============================================================================

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

// 계좌 이체 (트랜잭션 처리)
func (s *TransactionService) Transfer(ctx context.Context, fromAccountID, toAccountID uint, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	txRecord := &Transaction{
		TransactionID: fmt.Sprintf("TXN%d", time.Now().UnixNano()),
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		Type:          "transfer",
		Status:        "pending",
		Currency:      "USD",
	}

	startTime := time.Now()

	// 트랜잭션 시작 (타임아웃 설정)
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 트랜잭션 레코드 생성
		if err := tx.Create(txRecord).Error; err != nil {
			return err
		}

		// 2. 송금 계좌 조회 및 잠금 (비관적 잠금)
		var fromAccount Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&fromAccount, fromAccountID).Error; err != nil {
			return fmt.Errorf("from account not found: %w", err)
		}

		// 3. 수신 계좌 조회 및 잠금
		var toAccount Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&toAccount, toAccountID).Error; err != nil {
			return fmt.Errorf("to account not found: %w", err)
		}

		// 4. 잔액 확인
		if fromAccount.Balance < amount {
			return errors.New("insufficient balance")
		}

		// 5. 계좌 잠금 상태 확인
		if fromAccount.IsLocked || toAccount.IsLocked {
			return errors.New("account is locked")
		}

		// 6. 잔액 업데이트
		fromAccount.Balance -= amount
		toAccount.Balance += amount

		if err := tx.Save(&fromAccount).Error; err != nil {
			return fmt.Errorf("failed to update from account: %w", err)
		}

		if err := tx.Save(&toAccount).Error; err != nil {
			return fmt.Errorf("failed to update to account: %w", err)
		}

		// 7. 트랜잭션 상태 업데이트
		now := time.Now()
		txRecord.Status = "completed"
		txRecord.CompletedAt = &now
		txRecord.ProcessingTime = time.Since(startTime).Milliseconds()

		if err := tx.Save(txRecord).Error; err != nil {
			return fmt.Errorf("failed to update transaction record: %w", err)
		}

		// 인위적 지연 (테스트용)
		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable, // 최고 격리 수준
	})

	if err != nil {
		// 트랜잭션 실패 기록
		txRecord.Status = "failed"
		txRecord.ErrorMessage = err.Error()
		txRecord.ProcessingTime = time.Since(startTime).Milliseconds()
		s.db.Save(txRecord)
		return nil, err
	}

	return txRecord, nil
}

// 낙관적 잠금을 사용한 재고 업데이트
func (s *TransactionService) UpdateStock(ctx context.Context, productID uint, quantity int) error {
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var product Product
			if err := tx.First(&product, productID).Error; err != nil {
				return err
			}

			// 재고 확인
			if product.Stock < quantity {
				return errors.New("insufficient stock")
			}

			// 낙관적 잠금: Version 체크와 업데이트
			result := tx.Model(&Product{}).
				Where("id = ? AND version = ?", productID, product.Version).
				Updates(map[string]interface{}{
					"stock":   product.Stock - quantity,
					"version": product.Version + 1,
				})

			if result.RowsAffected == 0 {
				return errors.New("concurrent update detected")
			}

			return result.Error
		})

		if err == nil {
			return nil
		}

		if err.Error() == "concurrent update detected" {
			// 재시도 대기
			time.Sleep(time.Duration(i*50) * time.Millisecond)
			continue
		}

		return err
	}

	return errors.New("max retries exceeded")
}

// 주문 처리 (복잡한 트랜잭션)
func (s *TransactionService) ProcessOrder(ctx context.Context, order *Order) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 주문 생성
		order.Status = "processing"
		order.OrderNumber = fmt.Sprintf("ORD%d", time.Now().UnixNano())

		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// 2. 재고 확인 및 예약
		for _, item := range order.Items {
			var product Product

			// 비관적 잠금으로 제품 조회
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&product, item.ProductID).Error; err != nil {
				return fmt.Errorf("product not found: %w", err)
			}

			// 재고 확인
			availableStock := product.Stock - product.Reserved
			if availableStock < item.Quantity {
				return fmt.Errorf("insufficient stock for product %s", product.Name)
			}

			// 재고 예약
			product.Reserved += item.Quantity
			if err := tx.Save(&product).Error; err != nil {
				return fmt.Errorf("failed to reserve stock: %w", err)
			}

			// 주문 아이템 저장
			item.OrderID = order.ID
			item.Price = product.Price
			if err := tx.Create(&item).Error; err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}

		// 3. 결제 처리
		payment := &Payment{
			PaymentID: fmt.Sprintf("PAY%d", time.Now().UnixNano()),
			OrderID:   order.ID,
			Amount:    order.TotalAmount,
			Method:    "card",
			Status:    "processing",
		}

		if err := tx.Create(payment).Error; err != nil {
			return fmt.Errorf("failed to create payment: %w", err)
		}

		// 결제 처리 시뮬레이션
		select {
		case <-time.After(200 * time.Millisecond):
			payment.Status = "completed"
			now := time.Now()
			payment.ProcessedAt = &now
			tx.Save(payment)
		case <-ctx.Done():
			return ctx.Err()
		}

		// 4. 주문 상태 업데이트
		order.Status = "completed"
		order.PaymentID = &payment.ID
		if err := tx.Save(order).Error; err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}

		// 5. 재고 확정 (예약 → 실제 차감)
		for _, item := range order.Items {
			result := tx.Model(&Product{}).
				Where("id = ?", item.ProductID).
				Updates(map[string]interface{}{
					"stock":    gorm.Expr("stock - ?", item.Quantity),
					"reserved": gorm.Expr("reserved - ?", item.Quantity),
				})

			if result.Error != nil {
				return fmt.Errorf("failed to update stock: %w", result.Error)
			}
		}

		return nil
	})
}

// Saga 패턴 예시
func (s *TransactionService) ProcessOrderSaga(ctx context.Context, order *Order) error {
	// 각 단계를 독립적으로 처리하고 실패 시 보상 트랜잭션 실행

	// Step 1: 재고 예약
	if err := s.reserveStock(ctx, order); err != nil {
		return err
	}

	// Step 2: 결제 처리
	payment, err := s.processPayment(ctx, order)
	if err != nil {
		// 보상: 재고 예약 취소
		s.cancelStockReservation(ctx, order)
		return err
	}

	// Step 3: 주문 확정
	if err := s.confirmOrder(ctx, order, payment); err != nil {
		// 보상: 결제 취소, 재고 예약 취소
		s.cancelPayment(ctx, payment)
		s.cancelStockReservation(ctx, order)
		return err
	}

	return nil
}

func (s *TransactionService) reserveStock(ctx context.Context, order *Order) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range order.Items {
			var product Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&product, item.ProductID).Error; err != nil {
				return err
			}

			if product.Stock < item.Quantity {
				return errors.New("insufficient stock")
			}

			product.Reserved += item.Quantity
			if err := tx.Save(&product).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionService) cancelStockReservation(ctx context.Context, order *Order) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range order.Items {
			tx.Model(&Product{}).
				Where("id = ?", item.ProductID).
				Update("reserved", gorm.Expr("reserved - ?", item.Quantity))
		}
		return nil
	})
}

func (s *TransactionService) processPayment(ctx context.Context, order *Order) (*Payment, error) {
	payment := &Payment{
		PaymentID: fmt.Sprintf("PAY%d", time.Now().UnixNano()),
		OrderID:   order.ID,
		Amount:    order.TotalAmount,
		Method:    "card",
		Status:    "completed",
	}

	if err := s.db.WithContext(ctx).Create(payment).Error; err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *TransactionService) cancelPayment(ctx context.Context, payment *Payment) error {
	return s.db.WithContext(ctx).
		Model(payment).
		Update("status", "cancelled").Error
}

func (s *TransactionService) confirmOrder(ctx context.Context, order *Order, payment *Payment) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order.Status = "completed"
		order.PaymentID = &payment.ID
		return tx.Save(order).Error
	})
}

// ============================================================================
// 동시성 테스트 서비스
// ============================================================================

type ConcurrencyTestService struct {
	db      *gorm.DB
	service *TransactionService
}

func NewConcurrencyTestService(db *gorm.DB, service *TransactionService) *ConcurrencyTestService {
	return &ConcurrencyTestService{db: db, service: service}
}

// 동시 이체 테스트
func (s *ConcurrencyTestService) TestConcurrentTransfers(numWorkers int) map[string]interface{} {
	results := make(map[string]interface{})
	successCount := 0
	failureCount := 0
	timeoutCount := 0

	var wg sync.WaitGroup
	var mu sync.Mutex

	// 테스트 계좌 생성
	account1 := Account{Number: "TEST001", Name: "Test Account 1", Balance: 10000}
	account2 := Account{Number: "TEST002", Name: "Test Account 2", Balance: 10000}
	s.db.Create(&account1)
	s.db.Create(&account2)

	startTime := time.Now()

	// 동시 이체 실행
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// 랜덤 타임아웃 설정
			timeout := time.Duration(rand.Intn(500)+500) * time.Millisecond
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			// 양방향 이체
			var fromID, toID uint
			if workerID%2 == 0 {
				fromID, toID = account1.ID, account2.ID
			} else {
				fromID, toID = account2.ID, account1.ID
			}

			amount := float64(rand.Intn(100) + 1)

			_, err := s.service.Transfer(ctx, fromID, toID, amount)

			mu.Lock()
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					timeoutCount++
				} else {
					failureCount++
				}
			} else {
				successCount++
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// 최종 잔액 확인
	s.db.First(&account1, account1.ID)
	s.db.First(&account2, account2.ID)

	results["duration_ms"] = duration.Milliseconds()
	results["total_workers"] = numWorkers
	results["success_count"] = successCount
	results["failure_count"] = failureCount
	results["timeout_count"] = timeoutCount
	results["final_balance_1"] = account1.Balance
	results["final_balance_2"] = account2.Balance
	results["total_balance"] = account1.Balance + account2.Balance

	// 테스트 데이터 정리
	s.db.Delete(&account1)
	s.db.Delete(&account2)

	return results
}

// 데드락 테스트
func (s *ConcurrencyTestService) TestDeadlock() map[string]interface{} {
	results := make(map[string]interface{})

	// 테스트 계좌 생성
	account1 := Account{Number: "DEAD001", Name: "Deadlock Test 1", Balance: 1000}
	account2 := Account{Number: "DEAD002", Name: "Deadlock Test 2", Balance: 1000}
	s.db.Create(&account1)
	s.db.Create(&account2)

	var wg sync.WaitGroup
	deadlockDetected := false

	// 데드락 유발 시나리오
	wg.Add(2)

	// Worker 1: Account1 → Account2
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Account1 잠금
			var acc1 Account
			tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&acc1, account1.ID)

			// 대기
			time.Sleep(100 * time.Millisecond)

			// Account2 잠금 시도
			var acc2 Account
			return tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&acc2, account2.ID).Error
		})

		if err != nil {
			deadlockDetected = true
		}
	}()

	// Worker 2: Account2 → Account1 (역순)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Account2 잠금
			var acc2 Account
			tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&acc2, account2.ID)

			// 대기
			time.Sleep(100 * time.Millisecond)

			// Account1 잠금 시도
			var acc1 Account
			return tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&acc1, account1.ID).Error
		})

		if err != nil {
			deadlockDetected = true
		}
	}()

	wg.Wait()

	results["deadlock_detected"] = deadlockDetected
	results["test_scenario"] = "Circular wait condition"

	// 테스트 데이터 정리
	s.db.Delete(&account1)
	s.db.Delete(&account2)

	return results
}

// ============================================================================
// HTTP Handlers
// ============================================================================

type Handler struct {
	service     *TransactionService
	testService *ConcurrencyTestService
}

func NewHandler(db *gorm.DB) *Handler {
	service := NewTransactionService(db)
	testService := NewConcurrencyTestService(db, service)

	return &Handler{
		service:     service,
		testService: testService,
	}
}

// 계좌 이체
func (h *Handler) Transfer(c *gin.Context) {
	var req struct {
		FromAccountID uint    `json:"from_account_id" binding:"required"`
		ToAccountID   uint    `json:"to_account_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		Timeout       int     `json:"timeout_ms"` // 밀리초
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 타임아웃 설정
	timeout := 5 * time.Second
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Millisecond
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
	defer cancel()

	transaction, err := h.service.Transfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(408, gin.H{"error": "Transaction timeout"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, transaction)
}

// 주문 처리
func (h *Handler) ProcessOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 타임아웃 설정
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.ProcessOrder(ctx, &order); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(408, gin.H{"error": "Order processing timeout"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, order)
}

// 재고 업데이트
func (h *Handler) UpdateStock(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.UpdateStock(ctx, req.ProductID, req.Quantity); err != nil {
		if err.Error() == "max retries exceeded" {
			c.JSON(409, gin.H{"error": "Conflict: Too many concurrent updates"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Stock updated successfully"})
}

// 동시성 테스트
func (h *Handler) TestConcurrency(c *gin.Context) {
	workers := c.DefaultQuery("workers", "10")
	var numWorkers int
	fmt.Sscanf(workers, "%d", &numWorkers)

	if numWorkers < 1 || numWorkers > 100 {
		numWorkers = 10
	}

	results := h.testService.TestConcurrentTransfers(numWorkers)
	c.JSON(200, results)
}

// 데드락 테스트
func (h *Handler) TestDeadlock(c *gin.Context) {
	results := h.testService.TestDeadlock()
	c.JSON(200, results)
}

// 트랜잭션 이력 조회
func (h *Handler) GetTransactionHistory(c *gin.Context) {
	var transactions []Transaction

	query := h.service.db.Order("created_at DESC").Limit(100)

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Find(&transactions)

	c.JSON(200, gin.H{
		"transactions": transactions,
		"count":       len(transactions),
	})
}

// ============================================================================
// 초기 데이터 생성
// ============================================================================

func InitializeData(db *gorm.DB) {
	// 계좌 생성
	accounts := []Account{
		{Number: "ACC001", Name: "Alice Johnson", Balance: 5000, Currency: "USD"},
		{Number: "ACC002", Name: "Bob Smith", Balance: 3000, Currency: "USD"},
		{Number: "ACC003", Name: "Charlie Brown", Balance: 10000, Currency: "USD"},
		{Number: "ACC004", Name: "Diana Prince", Balance: 7500, Currency: "USD"},
		{Number: "ACC005", Name: "Eve Adams", Balance: 2000, Currency: "USD"},
	}

	for _, acc := range accounts {
		db.Create(&acc)
	}

	// 제품 생성
	products := []Product{
		{Name: "Laptop", SKU: "SKU001", Price: 999.99, Stock: 50},
		{Name: "Mouse", SKU: "SKU002", Price: 29.99, Stock: 200},
		{Name: "Keyboard", SKU: "SKU003", Price: 79.99, Stock: 150},
		{Name: "Monitor", SKU: "SKU004", Price: 299.99, Stock: 75},
		{Name: "Headphones", SKU: "SKU005", Price: 149.99, Stock: 100},
	}

	for _, product := range products {
		db.Create(&product)
	}

	log.Println("✅ Initial data created")
}

// ============================================================================
// Router Setup
// ============================================================================

func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	// Middleware for request ID
	router.Use(func(c *gin.Context) {
		c.Set("request_id", fmt.Sprintf("REQ%d", time.Now().UnixNano()))
		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	// Info endpoint
	router.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"database": "SQLite with Transactions",
			"features": []string{
				"ACID Transactions",
				"Context Timeout",
				"Optimistic Locking",
				"Pessimistic Locking",
				"Deadlock Detection",
				"Saga Pattern",
			},
		})
	})

	// Transaction routes
	transactions := router.Group("/transactions")
	{
		transactions.POST("/transfer", handler.Transfer)
		transactions.POST("/order", handler.ProcessOrder)
		transactions.POST("/stock", handler.UpdateStock)
		transactions.GET("/history", handler.GetTransactionHistory)
	}

	// Test routes
	tests := router.Group("/tests")
	{
		tests.GET("/concurrency", handler.TestConcurrency)
		tests.GET("/deadlock", handler.TestDeadlock)
	}

	// Account management
	router.GET("/accounts", func(c *gin.Context) {
		var accounts []Account
		handler.service.db.Find(&accounts)
		c.JSON(200, accounts)
	})

	// Product management
	router.GET("/products", func(c *gin.Context) {
		var products []Product
		handler.service.db.Find(&products)
		c.JSON(200, products)
	})

	return router
}

// ============================================================================
// Main
// ============================================================================

func main() {
	// Database connection
	db, err := gorm.Open(sqlite.Open("transaction.db?_journal_mode=WAL"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Enable WAL mode for better concurrency
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Auto migrate
	db.AutoMigrate(&Account{}, &Transaction{}, &Order{}, &OrderItem{}, &Product{}, &Payment{})

	// Initialize data
	var count int64
	db.Model(&Account{}).Count(&count)
	if count == 0 {
		InitializeData(db)
	}

	// Initialize handler
	handler := NewHandler(db)

	// Setup router
	router := SetupRouter(handler)

	// Start server
	log.Println("🚀 Server starting on :8080")
	log.Println("💰 Transaction service ready")
	log.Println("🔒 ACID compliance enabled")
	log.Println("⏱️ Context timeout support")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}