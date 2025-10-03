# Lesson 17: 트랜잭션과 컨텍스트 타임아웃 ⏱️

> ACID 트랜잭션, 동시성 제어, 타임아웃 처리 완벽 가이드

## 📌 이번 레슨에서 배우는 내용

데이터 일관성은 애플리케이션의 신뢰성을 좌우합니다. 트랜잭션은 여러 작업을 하나의 원자적 단위로 처리하고, 컨텍스트 타임아웃은 장시간 실행되는 작업을 제어합니다. 이번 레슨에서는 실전 시나리오를 통해 이러한 개념을 학습합니다.

### 핵심 학습 목표
- ✅ ACID 트랜잭션 처리
- ✅ 비관적/낙관적 잠금
- ✅ Context를 활용한 타임아웃
- ✅ 데드락 방지 전략
- ✅ Saga 패턴 구현
- ✅ 동시성 테스트

## 🏗 트랜잭션 아키텍처

### ACID 속성
- **Atomicity (원자성)**: 모두 성공하거나 모두 실패
- **Consistency (일관성)**: 데이터 무결성 유지
- **Isolation (격리성)**: 트랜잭션 간 간섭 방지
- **Durability (지속성)**: 커밋된 데이터 영구 저장

### 격리 수준 (Isolation Level)
```sql
-- SQLite는 SERIALIZABLE만 지원
-- PostgreSQL/MySQL은 다음 수준 지원:
READ UNCOMMITTED
READ COMMITTED
REPEATABLE READ
SERIALIZABLE
```

## 🛠 구현된 기능

### 1. **계좌 이체 시스템**
- 비관적 잠금으로 동시성 제어
- 잔액 확인 및 업데이트
- 트랜잭션 이력 기록
- 타임아웃 처리

### 2. **주문 처리 시스템**
- 복잡한 다단계 트랜잭션
- 재고 확인 및 예약
- 결제 처리
- Saga 패턴 구현

### 3. **재고 관리**
- 낙관적 잠금 (Version 필드)
- 재시도 로직
- 동시 업데이트 감지

## 🎯 주요 API 엔드포인트

### 트랜잭션 처리
```bash
POST /transactions/transfer  # 계좌 이체
POST /transactions/order     # 주문 처리
POST /transactions/stock     # 재고 업데이트
GET  /transactions/history   # 트랜잭션 이력
```

### 테스트 엔드포인트
```bash
GET  /tests/concurrency      # 동시성 테스트
GET  /tests/deadlock         # 데드락 테스트
```

### 데이터 조회
```bash
GET  /accounts               # 계좌 목록
GET  /products               # 제품 목록
```

## 💻 실습 가이드

### 1. 실행
```bash
cd gin/17
go run main.go

# 초기 데이터 자동 생성
# - 5개 계좌 (잔액 포함)
# - 5개 제품 (재고 포함)
```

### 2. 계좌 이체

#### 기본 이체
```bash
curl -X POST http://localhost:8080/transactions/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 100
  }'

# 응답
{
  "id": 1,
  "transaction_id": "TXN1234567890",
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 100,
  "status": "completed",
  "processing_time_ms": 105
}
```

#### 타임아웃 설정
```bash
# 500ms 타임아웃
curl -X POST http://localhost:8080/transactions/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 50,
    "timeout_ms": 500
  }'

# 타임아웃 발생 시
{
  "error": "Transaction timeout"
}
```

### 3. 주문 처리

#### 복잡한 트랜잭션
```bash
curl -X POST http://localhost:8080/transactions/order \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "total_amount": 1329.97,
    "items": [
      {"product_id": 1, "quantity": 1},
      {"product_id": 2, "quantity": 2},
      {"product_id": 4, "quantity": 1}
    ]
  }'

# 응답
{
  "id": 1,
  "order_number": "ORD1234567890",
  "status": "completed",
  "payment": {
    "payment_id": "PAY1234567890",
    "status": "completed"
  }
}
```

### 4. 재고 업데이트 (낙관적 잠금)

```bash
curl -X POST http://localhost:8080/transactions/stock \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 5
  }'

# 동시 업데이트 시 자동 재시도
{
  "message": "Stock updated successfully"
}

# 재시도 실패 시
{
  "error": "Conflict: Too many concurrent updates"
}
```

### 5. 동시성 테스트

#### 동시 이체 테스트
```bash
# 20개 워커로 동시 이체
curl "http://localhost:8080/tests/concurrency?workers=20" | jq

# 응답
{
  "duration_ms": 1250,
  "total_workers": 20,
  "success_count": 15,
  "failure_count": 3,
  "timeout_count": 2,
  "final_balance_1": 9850,
  "final_balance_2": 10150,
  "total_balance": 20000  // 잔액 합계 유지
}
```

#### 데드락 테스트
```bash
curl http://localhost:8080/tests/deadlock | jq

# 응답
{
  "deadlock_detected": true,
  "test_scenario": "Circular wait condition"
}
```

### 6. 트랜잭션 이력 조회

```bash
# 모든 트랜잭션
curl http://localhost:8080/transactions/history | jq

# 실패한 트랜잭션만
curl "http://localhost:8080/transactions/history?status=failed" | jq

# 응답
{
  "transactions": [
    {
      "transaction_id": "TXN1234567890",
      "type": "transfer",
      "status": "failed",
      "error_message": "insufficient balance",
      "processing_time_ms": 25
    }
  ],
  "count": 1
}
```

## 🔍 코드 하이라이트

### 비관적 잠금 (Pessimistic Locking)
```go
func (s *TransactionService) Transfer(ctx context.Context, fromID, toID uint, amount float64) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // SELECT ... FOR UPDATE
        var fromAccount Account
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&fromAccount, fromID).Error; err != nil {
            return err
        }

        // 잔액 확인 후 업데이트
        if fromAccount.Balance < amount {
            return errors.New("insufficient balance")
        }

        fromAccount.Balance -= amount
        return tx.Save(&fromAccount).Error
    })
}
```

### 낙관적 잠금 (Optimistic Locking)
```go
func (s *TransactionService) UpdateStock(ctx context.Context, productID uint, quantity int) error {
    for i := 0; i < maxRetries; i++ {
        err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
            var product Product
            tx.First(&product, productID)

            // Version 체크와 업데이트
            result := tx.Model(&Product{}).
                Where("id = ? AND version = ?", productID, product.Version).
                Updates(map[string]interface{}{
                    "stock":   product.Stock - quantity,
                    "version": product.Version + 1,
                })

            if result.RowsAffected == 0 {
                return errors.New("concurrent update detected")
            }

            return nil
        })

        if err == nil {
            return nil
        }

        // 재시도
        time.Sleep(time.Duration(i*50) * time.Millisecond)
    }

    return errors.New("max retries exceeded")
}
```

### Context Timeout 처리
```go
func (h *Handler) Transfer(c *gin.Context) {
    // 타임아웃 설정
    timeout := 5 * time.Second
    if req.Timeout > 0 {
        timeout = time.Duration(req.Timeout) * time.Millisecond
    }

    ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
    defer cancel()

    transaction, err := h.service.Transfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount)

    if errors.Is(err, context.DeadlineExceeded) {
        c.JSON(408, gin.H{"error": "Transaction timeout"})
        return
    }
}
```

### Saga 패턴
```go
func (s *TransactionService) ProcessOrderSaga(ctx context.Context, order *Order) error {
    // Step 1: 재고 예약
    if err := s.reserveStock(ctx, order); err != nil {
        return err
    }

    // Step 2: 결제 처리
    payment, err := s.processPayment(ctx, order)
    if err != nil {
        // 보상 트랜잭션: 재고 예약 취소
        s.cancelStockReservation(ctx, order)
        return err
    }

    // Step 3: 주문 확정
    if err := s.confirmOrder(ctx, order, payment); err != nil {
        // 보상 트랜잭션: 결제 취소, 재고 예약 취소
        s.cancelPayment(ctx, payment)
        s.cancelStockReservation(ctx, order)
        return err
    }

    return nil
}
```

### 데드락 방지
```go
// 항상 동일한 순서로 잠금 획득
func (s *TransactionService) TransferSafe(ctx context.Context, acc1ID, acc2ID uint) error {
    // ID 순서로 정렬
    if acc1ID > acc2ID {
        acc1ID, acc2ID = acc2ID, acc1ID
    }

    return s.db.Transaction(func(tx *gorm.DB) error {
        // 항상 작은 ID부터 잠금
        tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&account1, acc1ID)
        tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&account2, acc2ID)
        // ...
    })
}
```

## 🎨 트랜잭션 패턴

### 1. **2단계 커밋 (2PC)**
```go
// Prepare Phase
prepared := prepareAllResources()
if !prepared {
    rollbackAll()
    return error
}

// Commit Phase
commitAll()
```

### 2. **Saga Pattern**
- 긴 트랜잭션을 작은 단계로 분할
- 각 단계마다 보상 트랜잭션 정의
- 실패 시 역순으로 보상 실행

### 3. **Outbox Pattern**
```go
// 트랜잭션과 이벤트 발행 원자성 보장
tx.Begin()
tx.Save(order)
tx.Save(outboxEvent)
tx.Commit()

// 별도 프로세스에서 이벤트 발행
publishEventsFromOutbox()
```

## 📝 베스트 프랙티스

### 1. **트랜잭션 범위 최소화**
```go
// ❌ Bad: 긴 트랜잭션
tx.Begin()
processOrder()      // 10초
sendEmail()         // 5초
updateAnalytics()   // 3초
tx.Commit()

// ✅ Good: 짧은 트랜잭션
tx.Begin()
processOrder()
tx.Commit()

// 트랜잭션 외부에서 처리
go sendEmail()
go updateAnalytics()
```

### 2. **적절한 잠금 선택**
- **비관적 잠금**: 충돌이 자주 발생하는 경우
- **낙관적 잠금**: 충돌이 드문 경우
- **No Lock**: 읽기 전용 작업

### 3. **타임아웃 설정**
```go
// 항상 타임아웃 설정
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

db.WithContext(ctx).Transaction(...)
```

### 4. **멱등성 보장**
```go
// 트랜잭션 ID로 중복 방지
if exists := checkTransactionExists(txID); exists {
    return previousResult
}
```

## 🚀 성능 최적화

### 연결 풀 설정
```go
sqlDB.SetMaxOpenConns(25)      // 최대 연결 수
sqlDB.SetMaxIdleConns(5)       // 유휴 연결 수
sqlDB.SetConnMaxLifetime(5*time.Minute)  // 연결 수명
```

### WAL 모드 (SQLite)
```go
// Write-Ahead Logging으로 동시성 향상
db, err := gorm.Open(sqlite.Open("file.db?_journal_mode=WAL"))
```

### 배치 처리
```go
// 여러 작업을 하나의 트랜잭션으로
tx.Begin()
for _, item := range items {
    tx.Create(&item)
}
tx.Commit()
```

## 📚 추가 학습 자료

- [Database Transactions](https://www.postgresql.org/docs/current/tutorial-transactions.html)
- [Concurrency Control](https://en.wikipedia.org/wiki/Concurrency_control)
- [Saga Pattern](https://microservices.io/patterns/data/saga.html)
- [Context in Go](https://go.dev/blog/context)

## 🎯 정리

트랜잭션과 타임아웃 관리는 신뢰할 수 있는 애플리케이션의 핵심입니다:

- **ACID 보장**: 데이터 일관성 유지
- **동시성 제어**: 적절한 잠금 전략 선택
- **타임아웃**: 리소스 고갈 방지
- **에러 처리**: 보상 트랜잭션으로 복구

이제 프로덕션급 트랜잭션 처리를 구현할 수 있습니다! ⏱️