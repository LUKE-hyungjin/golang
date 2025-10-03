package services

import (
	"context"

	"github.com/example/banking-system/internal/config"
	"github.com/example/banking-system/internal/models"
	"github.com/example/banking-system/pkg/logger"
	"gorm.io/gorm"
)

// Container is the dependency injection container (Lesson 13)
type Container struct {
	DB     *gorm.DB
	Config *config.Config
	Logger *logger.Logger
}

// Service interfaces
type UserService interface {
	Register(email, username, password string) (*models.User, error)
	Login(email, password string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.User, int64, error)
}

type AccountService interface {
	Create(userID uint, accountType string) (*models.Account, error)
	GetByID(id uint) (*models.Account, error)
	GetByAccountNumber(number string) (*models.Account, error)
	Deposit(accountID uint, amount float64) error
	Withdraw(accountID uint, amount float64) error
	GetBalance(accountID uint) (float64, error)
	List(userID uint, offset, limit int) ([]*models.Account, int64, error)
}

type TransactionService interface {
	Transfer(ctx context.Context, fromAccountID, toAccountID uint, amount float64) (*models.Transaction, error)
	GetByID(id uint) (*models.Transaction, error)
	List(accountID uint, offset, limit int) ([]*models.Transaction, int64, error)
	ProcessTransaction(tx *models.Transaction) error
}

