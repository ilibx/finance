package entity

import (
	"time"

	"finance/internal/common/valueobject"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user aggregate root
type User struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	Username  string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string     `gorm:"size:255;not null" json:"-"`
	Email     string     `gorm:"size:100" json:"email"`
	Phone     string     `gorm:"size:20" json:"phone"`
	Nickname  string     `gorm:"size:50" json:"nickname"`
	Avatar    string     `gorm:"size:255" json:"avatar"`
	Balance   valueobject.Money `json:"balance"`
	Status    valueobject.Status `json:"status"`
	RoleID    uint       `gorm:"index" json:"role_id"`
	Role      *Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Role represents a user role
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string    `gorm:"size:255" json:"description"`
	Status      int       `gorm:"default:1" json:"status"` // 1:正常 0:禁用
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission represents a permission
type Permission struct {
	ID       uint     `gorm:"primaryKey" json:"id"`
	Name     string   `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Code     string   `gorm:"uniqueIndex;size:100;not null" json:"code"` // e.g., "user:create", "order:read"
	Resource string   `gorm:"size:50" json:"resource"`
	Action   string   `gorm:"size:20" json:"action"` // create, read, update, delete
	Status   int      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user
func NewUser(username, email, phone string) *User {
	now := time.Now()
	return &User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Balance:  valueobject.NewMoney(0, "CNY"),
		Status: valueobject.Status{
			Code:        "active",
			Description: "Active",
			UpdatedAt:   now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetPassword sets encrypted password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Recharge adds money to user balance
func (u *User) Recharge(amount float64) {
	u.Balance = u.Balance.Add(valueobject.NewMoney(amount, u.Balance.Currency))
	u.UpdatedAt = time.Now()
}

// Deduct subtracts money from user balance
func (u *User) Deduct(amount float64) {
	u.Balance = u.Balance.Subtract(valueobject.NewMoney(amount, u.Balance.Currency))
	u.UpdatedAt = time.Now()
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status.Code == "active"
}
