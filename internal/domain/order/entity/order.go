package entity

import (
"time"

"finance/internal/common/errors"
"finance/internal/common/valueobject"
)

// OrderStatus represents order status
type OrderStatus string

const (
OrderStatusPending   OrderStatus = "pending"
OrderStatusPaid      OrderStatus = "paid"
OrderStatusShipped   OrderStatus = "shipped"
OrderStatusCompleted OrderStatus = "completed"
OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem represents an order item
type OrderItem struct {
ID        int64
ProductID int64
Quantity  int
UnitPrice valueobject.Money
Subtotal  valueobject.Money
}

// Order represents an order aggregate root
type Order struct {
ID          int64
OrderNo     string
UserID      int64
Items       []OrderItem
TotalAmount valueobject.Money
Status      OrderStatus
PaidAt      *time.Time
ShippedAt   *time.Time
CompletedAt *time.Time
CreatedAt   time.Time
UpdatedAt   time.Time
}

// NewOrder creates a new order
func NewOrder(orderNo string, userID int64, items []OrderItem) *Order {
now := time.Now()
var totalAmount float64
for _, item := range items {
totalAmount += item.Subtotal.Amount
}

return &Order{
OrderNo:     orderNo,
UserID:      userID,
Items:       items,
TotalAmount: valueobject.NewMoney(totalAmount, "CNY"),
Status:      OrderStatusPending,
CreatedAt:   now,
UpdatedAt:   now,
}
}

// Pay marks the order as paid
func (o *Order) Pay() error {
if o.Status != OrderStatusPending {
return errors.ErrInvalidOrderStatus
}
now := time.Now()
o.Status = OrderStatusPaid
o.PaidAt = &now
o.UpdatedAt = now
return nil
}

// Ship marks the order as shipped
func (o *Order) Ship() error {
if o.Status != OrderStatusPaid {
return errors.ErrInvalidOrderStatus
}
now := time.Now()
o.Status = OrderStatusShipped
o.ShippedAt = &now
o.UpdatedAt = now
return nil
}

// Complete marks the order as completed
func (o *Order) Complete() error {
if o.Status != OrderStatusShipped {
return errors.ErrInvalidOrderStatus
}
now := time.Now()
o.Status = OrderStatusCompleted
o.CompletedAt = &now
o.UpdatedAt = now
return nil
}

// Cancel cancels the order
func (o *Order) Cancel() error {
if o.Status == OrderStatusCompleted || o.Status == OrderStatusCancelled {
return errors.ErrInvalidOrderStatus
}
now := time.Now()
o.Status = OrderStatusCancelled
o.UpdatedAt = now
return nil
}
