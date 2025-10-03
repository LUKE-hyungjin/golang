package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Base model with common fields
type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User represents a bank customer
type User struct {
	Base
	Email       string    `gorm:"uniqueIndex;not null" json:"email"`
	Username    string    `gorm:"uniqueIndex;not null" json:"username"`
	Password    string    `json:"-"` // Never expose password
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Phone       string    `json:"phone"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	PostalCode  string    `json:"postal_code"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	Role        string    `gorm:"default:'customer'" json:"role"` // customer, admin, staff
	LastLoginAt *time.Time `json:"last_login_at"`
	Accounts    []Account `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
}

// Account represents a bank account
type Account struct {
	Base
	AccountNumber string          `gorm:"uniqueIndex;not null" json:"account_number"`
	AccountType   string          `json:"account_type"` // savings, checking, credit
	Currency      string          `gorm:"default:'USD'" json:"currency"`
	Balance       float64         `gorm:"default:0" json:"balance"`
	AvailableBalance float64      `gorm:"default:0" json:"available_balance"`
	IsActive      bool            `gorm:"default:true" json:"is_active"`
	IsFrozen      bool            `gorm:"default:false" json:"is_frozen"`
	UserID        uint            `json:"user_id"`
	User          User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions  []Transaction   `gorm:"foreignKey:AccountID" json:"transactions,omitempty"`
	Version       int             `gorm:"default:0" json:"version"` // For optimistic locking
}

// Transaction represents a financial transaction
type Transaction struct {
	Base
	TransactionID   string    `gorm:"uniqueIndex;not null" json:"transaction_id"`
	Type            string    `json:"type"` // deposit, withdrawal, transfer, payment
	Status          string    `json:"status"` // pending, processing, completed, failed, cancelled
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Description     string    `json:"description"`
	Reference       string    `json:"reference"`
	AccountID       uint      `json:"account_id"`
	Account         Account   `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	ToAccountID     *uint     `json:"to_account_id,omitempty"` // For transfers
	ToAccount       *Account  `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
	ProcessingTime  int64     `json:"processing_time_ms"`
	CompletedAt     *time.Time `json:"completed_at"`
	FailureReason   string    `json:"failure_reason,omitempty"`
	Metadata        string    `gorm:"type:json" json:"metadata,omitempty"`
}

// AuditLog represents system audit logs
type AuditLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     *uint     `json:"user_id"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID uint      `json:"resource_id"`
	OldValue   string    `gorm:"type:json" json:"old_value,omitempty"`
	NewValue   string    `gorm:"type:json" json:"new_value,omitempty"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	RequestID  string    `json:"request_id"`
	StatusCode int       `json:"status_code"`
	Duration   int64     `json:"duration_ms"`
	CreatedAt  time.Time `json:"created_at"`
}

// Migration represents database migrations
type Migration struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Version   string    `gorm:"uniqueIndex;not null" json:"version"`
	Name      string    `json:"name"`
	AppliedAt time.Time `json:"applied_at"`
}

// Session represents user sessions (for session-based auth)
type Session struct {
	ID        string    `gorm:"primarykey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Notification represents user notifications
type Notification struct {
	Base
	UserID   uint   `json:"user_id"`
	User     User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type     string `json:"type"` // email, sms, push
	Title    string `json:"title"`
	Message  string `json:"message"`
	IsRead   bool   `gorm:"default:false" json:"is_read"`
	ReadAt   *time.Time `json:"read_at"`
	Metadata string `gorm:"type:json" json:"metadata,omitempty"`
}

// TableName specifications
func (User) TableName() string { return "users" }
func (Account) TableName() string { return "accounts" }
func (Transaction) TableName() string { return "transactions" }
func (AuditLog) TableName() string { return "audit_logs" }
func (Migration) TableName() string { return "migrations" }
func (Session) TableName() string { return "sessions" }
func (Notification) TableName() string { return "notifications" }

// Hooks for User model
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Hash password before creating user
	// This would be implemented in the service layer
	return nil
}

// Hooks for Account model
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	// Generate account number if not provided
	if a.AccountNumber == "" {
		a.AccountNumber = generateAccountNumber()
	}
	return nil
}

// Hooks for Transaction model
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	// Generate transaction ID if not provided
	if t.TransactionID == "" {
		t.TransactionID = generateTransactionID()
	}
	return nil
}

// Helper functions
func generateAccountNumber() string {
	// In production, use a proper algorithm
	return fmt.Sprintf("ACC%d", time.Now().UnixNano())
}

func generateTransactionID() string {
	// In production, use a proper algorithm
	return fmt.Sprintf("TXN%d", time.Now().UnixNano())
}

