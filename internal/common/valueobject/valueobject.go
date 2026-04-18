package valueobject

import (
	"time"
)

// Money represents a monetary value
type Money struct {
	Amount   float64
	Currency string
}

// NewMoney creates a new Money value object
func NewMoney(amount float64, currency string) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}

// Add adds two money values
func (m Money) Add(other Money) Money {
	if m.Currency != other.Currency {
		panic("cannot add money with different currencies")
	}
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}
}

// Subtract subtracts two money values
func (m Money) Subtract(other Money) Money {
	if m.Currency != other.Currency {
		panic("cannot subtract money with different currencies")
	}
	return Money{
		Amount:   m.Amount - other.Amount,
		Currency: m.Currency,
	}
}

// AuditLog represents an audit log entry
type AuditLog struct {
	Action    string
	Timestamp time.Time
	UserID    int64
	Details   map[string]interface{}
}

// Status represents a status value object
type Status struct {
	Code        string
	Description string
	UpdatedAt   time.Time
}

// Address represents an address value object
type Address struct {
	Street  string
	City    string
	State   string
	Country string
	ZipCode string
}

// String returns the string representation of the address
func (a Address) String() string {
	return a.Street + ", " + a.City + ", " + a.State + " " + a.ZipCode + ", " + a.Country
}

// ContactInfo represents contact information
type ContactInfo struct {
	Phone string
	Email string
	Name  string
}
