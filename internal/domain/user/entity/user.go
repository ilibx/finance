package entity

import (
	"time"

	"finance/internal/common/valueobject"
)

// User represents a user aggregate root
type User struct {
	ID        int64
	Username  string
	Email     string
	Phone     string
	Balance   valueobject.Money
	Status    valueobject.Status
	CreatedAt time.Time
	UpdatedAt time.Time
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
