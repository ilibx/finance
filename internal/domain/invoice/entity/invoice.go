package entity

import (
"time"

"finance/internal/common/errors"
"finance/internal/common/valueobject"
)

// InvoiceType represents invoice type
type InvoiceType string

const (
InvoiceTypeSales    InvoiceType = "sales"
InvoiceTypePurchase InvoiceType = "purchase"
)

// InvoiceStatus represents invoice status
type InvoiceStatus string

const (
InvoiceStatusDraft   InvoiceStatus = "draft"
InvoiceStatusIssued  InvoiceStatus = "issued"
InvoiceStatusPaid    InvoiceStatus = "paid"
InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

// Invoice represents an invoice aggregate root
type Invoice struct {
ID          int64
InvoiceNo   string
OrderID     int64
UserID      int64
Amount      valueobject.Money
TaxAmount   valueobject.Money
TotalAmount valueobject.Money
Type        InvoiceType
Status      InvoiceStatus
IssuedAt    *time.Time
PaidAt      *time.Time
CreatedAt   time.Time
UpdatedAt   time.Time
}

// NewInvoice creates a new invoice
func NewInvoice(invoiceNo string, orderID, userID int64, amount, taxAmount float64, invoiceType InvoiceType) *Invoice {
now := time.Now()
totalAmount := amount + taxAmount

return &Invoice{
InvoiceNo:   invoiceNo,
OrderID:     orderID,
UserID:      userID,
Amount:      valueobject.NewMoney(amount, "CNY"),
TaxAmount:   valueobject.NewMoney(taxAmount, "CNY"),
TotalAmount: valueobject.NewMoney(totalAmount, "CNY"),
Type:        invoiceType,
Status:      InvoiceStatusDraft,
CreatedAt:   now,
UpdatedAt:   now,
}
}

// Issue issues the invoice
func (i *Invoice) Issue() error {
if i.Status != InvoiceStatusDraft {
return errors.ErrInvoiceAlreadyIssued
}
now := time.Now()
i.Status = InvoiceStatusIssued
i.IssuedAt = &now
i.UpdatedAt = now
return nil
}

// Pay marks the invoice as paid
func (i *Invoice) Pay() error {
if i.Status != InvoiceStatusIssued {
return errors.ErrInvalidOrderStatus
}
now := time.Now()
i.Status = InvoiceStatusPaid
i.PaidAt = &now
i.UpdatedAt = now
return nil
}

// Cancel cancels the invoice
func (i *Invoice) Cancel() error {
if i.Status == InvoiceStatusPaid || i.Status == InvoiceStatusCancelled {
return errors.ErrInvalidOrderStatus
}
now := time.Now()
i.Status = InvoiceStatusCancelled
i.UpdatedAt = now
return nil
}
