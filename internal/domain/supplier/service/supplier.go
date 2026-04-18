package service

import (
"context"

"finance/internal/common/errors"
"finance/internal/common/valueobject"
"finance/internal/domain/supplier/entity"
"finance/internal/domain/supplier/repository"
)

// SupplierService handles supplier business logic
type SupplierService struct {
supplierRepo repository.SupplierRepository
}

// NewSupplierService creates a new supplier service
func NewSupplierService(supplierRepo repository.SupplierRepository) *SupplierService {
return &SupplierService{
supplierRepo: supplierRepo,
}
}

// CreateSupplier creates a new supplier
func (s *SupplierService) CreateSupplier(ctx context.Context, name, phone, email, address string) (*entity.Supplier, error) {
contact := valueobject.ContactInfo{
Name:  name,
Phone: phone,
Email: email,
}
addr := valueobject.Address{
Street: address,
Country: "China",
}
supplier := entity.NewSupplier(name, contact, addr)

if err := s.supplierRepo.Create(ctx, supplier); err != nil {
return nil, err
}

return supplier, nil
}

// GetSupplierByID gets a supplier by ID
func (s *SupplierService) GetSupplierByID(ctx context.Context, id int64) (*entity.Supplier, error) {
supplier, err := s.supplierRepo.GetByID(ctx, id)
if err != nil {
return nil, errors.ErrSupplierNotFound
}
return supplier, nil
}

// RechargeSupplier recharges a supplier's balance
func (s *SupplierService) RechargeSupplier(ctx context.Context, supplierID int64, amount float64) error {
supplier, err := s.supplierRepo.GetByID(ctx, supplierID)
if err != nil {
return errors.ErrSupplierNotFound
}

if amount <= 0 {
return errors.ErrInvalidRechargeAmount
}

supplier.AddBalance(amount)
return s.supplierRepo.Update(ctx, supplier)
}

// ListSuppliers lists suppliers with pagination
func (s *SupplierService) ListSuppliers(ctx context.Context, limit, offset int) ([]*entity.Supplier, error) {
return s.supplierRepo.List(ctx, limit, offset)
}
