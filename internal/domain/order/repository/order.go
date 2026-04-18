package repository

import (
"context"

"erp-system/internal/domain/order/entity"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
Create(ctx context.Context, order *entity.Order) error
GetByID(ctx context.Context, id int64) (*entity.Order, error)
GetByOrderNo(ctx context.Context, orderNo string) (*entity.Order, error)
Update(ctx context.Context, order *entity.Order) error
List(ctx context.Context, limit, offset int) ([]*entity.Order, error)
}
