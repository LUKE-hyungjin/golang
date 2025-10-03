package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/example/banking-system/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TransactionServiceImpl implements TransactionService (Lesson 17)
type TransactionServiceImpl struct {
	container *Container
}

func NewTransactionService(container *Container) TransactionService {
	return &TransactionServiceImpl{container: container}
}

// Transfer performs money transfer between accounts with ACID guarantees
func (s *TransactionServiceImpl) Transfer(ctx context.Context, fromAccountID, toAccountID uint, amount float64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	if fromAccountID == toAccountID {
		return nil, errors.New("cannot transfer to the same account")
	}

	// Create transaction record
	txRecord := &models.Transaction{
		Type:          "transfer",
		Status:        "pending",
		Amount:        amount,
		Currency:      "USD",
		AccountID:     fromAccountID,
		ToAccountID:   &toAccountID,
		Description:   fmt.Sprintf("Transfer from account %d to account %d", fromAccountID, toAccountID),
	}

	startTime := time.Now()

	// Start database transaction with context timeout
	err := s.container.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create transaction record
		if err := tx.Create(txRecord).Error; err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}

		// 2. Lock and fetch source account (pessimistic locking)
		var fromAccount models.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&fromAccount, fromAccountID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("source account not found")
			}
			return fmt.Errorf("failed to fetch source account: %w", err)
		}

		// 3. Lock and fetch destination account
		var toAccount models.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&toAccount, toAccountID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("destination account not found")
			}
			return fmt.Errorf("failed to fetch destination account: %w", err)
		}

		// 4. Check account status
		if !fromAccount.IsActive || fromAccount.IsFrozen {
			return errors.New("source account is not active or frozen")
		}

		if !toAccount.IsActive || toAccount.IsFrozen {
			return errors.New("destination account is not active or frozen")
		}

		// 5. Check balance
		if fromAccount.AvailableBalance < amount {
			return errors.New("insufficient balance")
		}

		// 6. Update balances
		fromAccount.Balance -= amount
		fromAccount.AvailableBalance -= amount
		toAccount.Balance += amount
		toAccount.AvailableBalance += amount

		// 7. Save updated accounts
		if err := tx.Save(&fromAccount).Error; err != nil {
			return fmt.Errorf("failed to update source account: %w", err)
		}

		if err := tx.Save(&toAccount).Error; err != nil {
			return fmt.Errorf("failed to update destination account: %w", err)
		}

		// 8. Update transaction status
		now := time.Now()
		txRecord.Status = "completed"
		txRecord.CompletedAt = &now
		txRecord.ProcessingTime = time.Since(startTime).Milliseconds()

		if err := tx.Save(txRecord).Error; err != nil {
			return fmt.Errorf("failed to update transaction record: %w", err)
		}

		// 9. Create audit log
		s.createAuditLog(tx, "transfer", txRecord)

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable, // Highest isolation level
	})

	if err != nil {
		// Update transaction as failed
		txRecord.Status = "failed"
		txRecord.FailureReason = err.Error()
		txRecord.ProcessingTime = time.Since(startTime).Milliseconds()
		s.container.DB.Save(txRecord)

		s.container.Logger.Error("Transfer failed", map[string]interface{}{
			"from_account": fromAccountID,
			"to_account":   toAccountID,
			"amount":       amount,
			"error":        err.Error(),
		})

		return nil, err
	}

	s.container.Logger.Info("Transfer completed", map[string]interface{}{
		"transaction_id": txRecord.TransactionID,
		"from_account":   fromAccountID,
		"to_account":     toAccountID,
		"amount":         amount,
		"processing_ms":  txRecord.ProcessingTime,
	})

	return txRecord, nil
}

// GetByID retrieves a transaction by ID
func (s *TransactionServiceImpl) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := s.container.DB.Preload("Account").Preload("ToAccount").First(&transaction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

// List retrieves transactions for an account
func (s *TransactionServiceImpl) List(accountID uint, offset, limit int) ([]*models.Transaction, int64, error) {
	var transactions []*models.Transaction
	var total int64

	query := s.container.DB.Model(&models.Transaction{})

	if accountID > 0 {
		query = query.Where("account_id = ? OR to_account_id = ?", accountID, accountID)
	}

	// Get total count
	query.Count(&total)

	// Get paginated results
	err := query.
		Preload("Account").
		Preload("ToAccount").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

// ProcessTransaction processes a single transaction (used for batch processing)
func (s *TransactionServiceImpl) ProcessTransaction(tx *models.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch tx.Type {
	case "deposit":
		return s.processDeposit(ctx, tx)
	case "withdrawal":
		return s.processWithdrawal(ctx, tx)
	case "transfer":
		if tx.ToAccountID == nil {
			return errors.New("transfer requires destination account")
		}
		_, err := s.Transfer(ctx, tx.AccountID, *tx.ToAccountID, tx.Amount)
		return err
	default:
		return fmt.Errorf("unsupported transaction type: %s", tx.Type)
	}
}

// processDeposit handles deposit transactions
func (s *TransactionServiceImpl) processDeposit(ctx context.Context, tx *models.Transaction) error {
	return s.container.DB.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		var account models.Account

		// Lock account for update
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&account, tx.AccountID).Error; err != nil {
			return err
		}

		// Update balance
		account.Balance += tx.Amount
		account.AvailableBalance += tx.Amount

		if err := db.Save(&account).Error; err != nil {
			return err
		}

		// Update transaction status
		tx.Status = "completed"
		now := time.Now()
		tx.CompletedAt = &now

		return db.Save(tx).Error
	})
}

// processWithdrawal handles withdrawal transactions
func (s *TransactionServiceImpl) processWithdrawal(ctx context.Context, tx *models.Transaction) error {
	return s.container.DB.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		var account models.Account

		// Lock account for update
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&account, tx.AccountID).Error; err != nil {
			return err
		}

		// Check balance
		if account.AvailableBalance < tx.Amount {
			return errors.New("insufficient balance")
		}

		// Update balance
		account.Balance -= tx.Amount
		account.AvailableBalance -= tx.Amount

		if err := db.Save(&account).Error; err != nil {
			return err
		}

		// Update transaction status
		tx.Status = "completed"
		now := time.Now()
		tx.CompletedAt = &now

		return db.Save(tx).Error
	})
}

// createAuditLog creates an audit log entry
func (s *TransactionServiceImpl) createAuditLog(tx *gorm.DB, action string, transaction *models.Transaction) {
	audit := &models.AuditLog{
		Action:     action,
		Resource:   "transaction",
		ResourceID: transaction.ID,
		NewValue:   fmt.Sprintf(`{"amount": %f, "status": "%s"}`, transaction.Amount, transaction.Status),
		CreatedAt:  time.Now(),
	}
	tx.Create(audit)
}

// Optimistic locking example for high-concurrency updates
func (s *TransactionServiceImpl) UpdateAccountWithOptimisticLocking(accountID uint, amount float64) error {
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		err := s.container.DB.Transaction(func(tx *gorm.DB) error {
			var account models.Account
			if err := tx.First(&account, accountID).Error; err != nil {
				return err
			}

			// Check version and update
			result := tx.Model(&models.Account{}).
				Where("id = ? AND version = ?", accountID, account.Version).
				Updates(map[string]interface{}{
					"balance":           account.Balance + amount,
					"available_balance": account.AvailableBalance + amount,
					"version":          account.Version + 1,
				})

			if result.RowsAffected == 0 {
				return errors.New("concurrent update detected")
			}

			return result.Error
		})

		if err == nil {
			return nil
		}

		if err.Error() == "concurrent update detected" && i < maxRetries-1 {
			// Retry with exponential backoff
			time.Sleep(time.Duration(i*50) * time.Millisecond)
			continue
		}

		return err
	}

	return errors.New("max retries exceeded")
}