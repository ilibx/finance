package entity

import (
	"time"

	"finance/internal/common/errors"
	"finance/internal/common/valueobject"
)

// PurchaseStatus represents purchase order status
type PurchaseStatus string

const (
	PurchaseStatusDraft     PurchaseStatus = "draft"
	PurchaseStatusPending   PurchaseStatus = "pending"
	PurchaseStatusApproved  PurchaseStatus = "approved"
	PurchaseStatusOrdered   PurchaseStatus = "ordered"
	PurchaseStatusPartial   PurchaseStatus = "partial_received"
	PurchaseStatusCompleted PurchaseStatus = "completed"
	PurchaseStatusCancelled PurchaseStatus = "cancelled"
)

// PurchaseItem represents a purchase order item
type PurchaseItem struct {
	ID             int64
	ProductID      int64
	ProductName    string
	Quantity       int
	ReceivedQty    int
	UnitPrice      valueobject.Money
	Subtotal       valueobject.Money
	SupplierPartNo string
}

// Purchase represents a purchase order aggregate root
type Purchase struct {
	ID            int64
	PurchaseNo    string
	SupplierID    int64
	SupplierName  string
	Items         []PurchaseItem
	TotalAmount   valueobject.Money
	Status        PurchaseStatus
	CreatedBy     int64
	ApprovedBy    *int64
	ApprovedAt    *time.Time
	OrderedAt     *time.Time
	CompletedAt   *time.Time
	Notes         string
	DeliveryDate  *time.Time
	ReceivedDate  *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewPurchase creates a new purchase order
func NewPurchase(purchaseNo string, supplierID int64, supplierName string, items []PurchaseItem, createdBy int64, notes string) *Purchase {
	now := time.Now()
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Subtotal.Amount
	}

	return &Purchase{
		PurchaseNo:   purchaseNo,
		SupplierID:   supplierID,
		SupplierName: supplierName,
		Items:        items,
		TotalAmount:  valueobject.NewMoney(totalAmount, "CNY"),
		Status:       PurchaseStatusDraft,
		CreatedBy:    createdBy,
		Notes:        notes,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Submit submits the purchase order for approval
func (p *Purchase) Submit() error {
	if p.Status != PurchaseStatusDraft {
		return errors.NewDomainError("INVALID_STATUS", "Only draft purchases can be submitted", nil)
	}
	if len(p.Items) == 0 {
		return errors.NewDomainError("EMPTY_ITEMS", "Purchase order must have at least one item", nil)
	}
	p.Status = PurchaseStatusPending
	p.UpdatedAt = time.Now()
	return nil
}

// Approve approves the purchase order
func (p *Purchase) Approve(approverID int64) error {
	if p.Status != PurchaseStatusPending {
		return errors.NewDomainError("INVALID_STATUS", "Only pending purchases can be approved", nil)
	}
	p.Status = PurchaseStatusApproved
	p.ApprovedBy = &approverID
	now := time.Now()
	p.ApprovedAt = &now
	p.UpdatedAt = now
	return nil
}

// Reject rejects the purchase order
func (p *Purchase) Reject() error {
	if p.Status != PurchaseStatusPending {
		return errors.NewDomainError("INVALID_STATUS", "Only pending purchases can be rejected", nil)
	}
	p.Status = PurchaseStatusDraft
	p.UpdatedAt = time.Now()
	return nil
}

// Order places the purchase order to supplier
func (p *Purchase) Order() error {
	if p.Status != PurchaseStatusApproved {
		return errors.NewDomainError("INVALID_STATUS", "Only approved purchases can be ordered", nil)
	}
	p.Status = PurchaseStatusOrdered
	now := time.Now()
	p.OrderedAt = &now
	p.UpdatedAt = now
	return nil
}

// Receive receives items from the purchase order
func (p *Purchase) Receive(itemID int64, quantity int) error {
	if p.Status != PurchaseStatusOrdered && p.Status != PurchaseStatusPartial {
		return errors.NewDomainError("INVALID_STATUS", "Can only receive items from ordered purchases", nil)
	}

	for i, item := range p.Items {
		if item.ID == itemID {
			if item.ReceivedQty+quantity > item.Quantity {
				return errors.NewDomainError("EXCESS_RECEIPT", "Received quantity exceeds ordered quantity", nil)
			}
			p.Items[i].ReceivedQty += quantity
			
			// Check if all items are fully received
			allReceived := true
			for _, it := range p.Items {
				if it.ReceivedQty < it.Quantity {
					allReceived = false
					break
				}
			}
			
			if allReceived {
				p.Status = PurchaseStatusCompleted
				now := time.Now()
				p.CompletedAt = &now
			} else {
				p.Status = PurchaseStatusPartial
			}
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.NewDomainError("ITEM_NOT_FOUND", "Purchase item not found", nil)
}

// Cancel cancels the purchase order
func (p *Purchase) Cancel() error {
	if p.Status == PurchaseStatusCompleted || p.Status == PurchaseStatusCancelled {
		return errors.NewDomainError("INVALID_STATUS", "Cannot cancel completed or cancelled purchases", nil)
	}
	p.Status = PurchaseStatusCancelled
	p.UpdatedAt = time.Now()
	return nil
}

// GetReceivedQuantity returns total received quantity
func (p *Purchase) GetReceivedQuantity() int {
	total := 0
	for _, item := range p.Items {
		total += item.ReceivedQty
	}
	return total
}

// GetTotalQuantity returns total ordered quantity
func (p *Purchase) GetTotalQuantity() int {
	total := 0
	for _, item := range p.Items {
		total += item.Quantity
	}
	return total
}

// IsFullyReceived checks if all items are fully received
func (p *Purchase) IsFullyReceived() bool {
	for _, item := range p.Items {
		if item.ReceivedQty < item.Quantity {
			return false
		}
	}
	return true
}
