package service

import (
	"context"
	"fmt"
	"time"

	"finance/internal/common/errors"
	"finance/internal/common/valueobject"
	productRepo "finance/internal/domain/product/repository"
	purchaseEntity "finance/internal/domain/purchase/entity"
	purchaseRepo "finance/internal/domain/purchase/repository"
	supplierRepo "finance/internal/domain/supplier/repository"
)

// PurchaseService handles purchase order business logic
type PurchaseService struct {
	purchaseRepo   purchaseRepo.PurchaseRepository
	productRepo    productRepo.ProductRepository
	supplierRepo   supplierRepo.SupplierRepository
}

// NewPurchaseService creates a new purchase service
func NewPurchaseService(
	purchaseRepo purchaseRepo.PurchaseRepository,
	productRepo productRepo.ProductRepository,
	supplierRepo supplierRepo.SupplierRepository,
) *PurchaseService {
	return &PurchaseService{
		purchaseRepo:   purchaseRepo,
		productRepo:    productRepo,
		supplierRepo:   supplierRepo,
	}
}

// CreatePurchase creates a new purchase order
func (s *PurchaseService) CreatePurchase(ctx context.Context, supplierID int64, items []PurchaseItemRequest, createdBy int64, notes string, deliveryDate *time.Time) (*purchaseEntity.Purchase, error) {
	// Get supplier
	supplier, err := s.supplierRepo.GetByID(ctx, supplierID)
	if err != nil {
		return nil, errors.ErrSupplierNotFound
	}

	// Build purchase items
	purchaseItems := make([]purchaseEntity.PurchaseItem, 0, len(items))
	for _, itemReq := range items {
		// Get product to validate and get name
		product, err := s.productRepo.GetByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found: %w", itemReq.ProductID, err)
		}

		unitPrice := valueobject.NewMoney(itemReq.UnitPrice, "CNY")
		subtotal := valueobject.NewMoney(itemReq.UnitPrice*float64(itemReq.Quantity), "CNY")

		purchaseItem := purchaseEntity.PurchaseItem{
			ProductID:      itemReq.ProductID,
			ProductName:    product.Name,
			Quantity:       itemReq.Quantity,
			ReceivedQty:    0,
			UnitPrice:      unitPrice,
			Subtotal:       subtotal,
			SupplierPartNo: itemReq.SupplierPartNo,
		}
		purchaseItems = append(purchaseItems, purchaseItem)
	}

	purchaseNo := generatePurchaseNo()
	purchase := purchaseEntity.NewPurchase(purchaseNo, supplier.ID, supplier.Name, purchaseItems, createdBy, notes)
	purchase.DeliveryDate = deliveryDate

	if err := s.purchaseRepo.Create(ctx, purchase); err != nil {
		return nil, err
	}

	return purchase, nil
}

// GetPurchaseByID gets a purchase order by ID
func (s *PurchaseService) GetPurchaseByID(ctx context.Context, id int64) (*purchaseEntity.Purchase, error) {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrOrderNotFound
	}
	return purchase, nil
}

// SubmitPurchase submits a purchase order for approval
func (s *PurchaseService) SubmitPurchase(ctx context.Context, id int64) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrOrderNotFound
	}

	if err := purchase.Submit(); err != nil {
		return err
	}

	return s.purchaseRepo.Update(ctx, purchase)
}

// ApprovePurchase approves a purchase order
func (s *PurchaseService) ApprovePurchase(ctx context.Context, id int64, approverID int64) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrOrderNotFound
	}

	if err := purchase.Approve(approverID); err != nil {
		return err
	}

	return s.purchaseRepo.Update(ctx, purchase)
}

// RejectPurchase rejects a purchase order
func (s *PurchaseService) RejectPurchase(ctx context.Context, id int64) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrOrderNotFound
	}

	if err := purchase.Reject(); err != nil {
		return err
	}

	return s.purchaseRepo.Update(ctx, purchase)
}

// OrderPurchase places the purchase order to supplier
func (s *PurchaseService) OrderPurchase(ctx context.Context, id int64) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrOrderNotFound
	}

	if err := purchase.Order(); err != nil {
		return err
	}

	return s.purchaseRepo.Update(ctx, purchase)
}

// ReceivePurchase receives items from a purchase order and updates stock
func (s *PurchaseService) ReceivePurchase(ctx context.Context, id int64, itemID int64, quantity int) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrOrderNotFound
	}

	// Find the item to receive
	var targetItem *purchaseEntity.PurchaseItem
	for i := range purchase.Items {
		if purchase.Items[i].ID == itemID {
			targetItem = &purchase.Items[i]
			break
		}
	}
	if targetItem == nil {
		return errors.NewDomainError("ITEM_NOT_FOUND", "Purchase item not found", nil)
	}

	// Validate quantity
	if targetItem.ReceivedQty+quantity > targetItem.Quantity {
		return errors.NewDomainError("EXCESS_RECEIPT", "Received quantity exceeds ordered quantity", nil)
	}

	// Update purchase status and received quantity
	if err := purchase.Receive(itemID, quantity); err != nil {
		return err
	}

	// Update product stock
	product, err := s.productRepo.GetByID(ctx, targetItem.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	product.AddStock(quantity)
	if err := s.productRepo.Update(ctx, product); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	// Check if this is the last receipt
	if purchase.IsFullyReceived() {
		now := time.Now()
		purchase.ReceivedDate = &now
	}

	return s.purchaseRepo.Update(ctx, purchase)
}

// CancelPurchase cancels a purchase order
func (s *PurchaseService) CancelPurchase(ctx context.Context, id int64) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrOrderNotFound
	}

	if err := purchase.Cancel(); err != nil {
		return err
	}

	return s.purchaseRepo.Update(ctx, purchase)
}

// ListPurchases lists purchase orders with optional status filter
func (s *PurchaseService) ListPurchases(ctx context.Context, status purchaseEntity.PurchaseStatus, limit, offset int) ([]*purchaseEntity.Purchase, error) {
	return s.purchaseRepo.List(ctx, status, limit, offset)
}

// ListPurchasesBySupplier lists purchase orders by supplier
func (s *PurchaseService) ListPurchasesBySupplier(ctx context.Context, supplierID int64, limit, offset int) ([]*purchaseEntity.Purchase, error) {
	return s.purchaseRepo.ListBySupplier(ctx, supplierID, limit, offset)
}

// DeletePurchase deletes a purchase order (only draft status)
func (s *PurchaseService) DeletePurchase(ctx context.Context, id int64) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrOrderNotFound
	}

	if purchase.Status != purchaseEntity.PurchaseStatusDraft {
		return errors.NewDomainError("INVALID_STATUS", "Only draft purchases can be deleted", nil)
	}

	return s.purchaseRepo.Delete(ctx, id)
}

// PurchaseItemRequest represents a purchase item in the request
type PurchaseItemRequest struct {
	ProductID      int64   `json:"product_id"`
	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	SupplierPartNo string  `json:"supplier_part_no"`
}

func generatePurchaseNo() string {
	return fmt.Sprintf("PO%s", time.Now().Format("20060102150405"))
}
