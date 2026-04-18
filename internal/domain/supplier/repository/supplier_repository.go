package repository

import (
"context"

"erp-system/internal/domain/supplier/entity"
)

// SupplierRepository defines the interface for supplier data access
type SupplierRepository interface {
Create(ctx context.Context, supplier *entity.Supplier) error
GetByID(ctx context.Context, id int64) (*entity.Supplier, error)
Update(ctx context.Context, supplier *entity.Supplier) error
List(ctx context.Context, limit, offset int) ([]*entity.Supplier, error)
}
