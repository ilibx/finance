package repository

import (
"context"

"erp-system/internal/domain/product/entity"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
Create(ctx context.Context, product *entity.Product) error
GetByID(ctx context.Context, id int64) (*entity.Product, error)
GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
Update(ctx context.Context, product *entity.Product) error
List(ctx context.Context, limit, offset int) ([]*entity.Product, error)
}
