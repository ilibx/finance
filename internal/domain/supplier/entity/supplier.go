package entity

import (
"time"

"finance/internal/common/valueobject"
)

// Supplier represents a supplier aggregate root
type Supplier struct {
ID        int64
Name      string
Contact   valueobject.ContactInfo
Address   valueobject.Address
Balance   valueobject.Money
Status    valueobject.Status
CreatedAt time.Time
UpdatedAt time.Time
}

// NewSupplier creates a new supplier
func NewSupplier(name string, contact valueobject.ContactInfo, address valueobject.Address) *Supplier {
now := time.Now()
return &Supplier{
Name:    name,
Contact: contact,
Address: address,
Balance: valueobject.NewMoney(0, "CNY"),
Status: valueobject.Status{
Code:        "active",
Description: "Active",
UpdatedAt:   now,
},
CreatedAt: now,
UpdatedAt: now,
}
}

// AddBalance adds balance to supplier
func (s *Supplier) AddBalance(amount float64) {
s.Balance = s.Balance.Add(valueobject.NewMoney(amount, s.Balance.Currency))
s.UpdatedAt = time.Now()
}

// DeductBalance deducts balance from supplier
func (s *Supplier) DeductBalance(amount float64) {
s.Balance = s.Balance.Subtract(valueobject.NewMoney(amount, s.Balance.Currency))
s.UpdatedAt = time.Now()
}

// IsActive checks if supplier is active
func (s *Supplier) IsActive() bool {
return s.Status.Code == "active"
}
