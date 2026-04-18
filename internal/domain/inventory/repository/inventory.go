package repository

import (
	"context"
	"finance/internal/domain/inventory/entity"
)

// InventoryAlertRepository defines the interface for inventory alert data access
type InventoryAlertRepository interface {
	Create(ctx context.Context, alert *entity.InventoryAlert) error
	GetByID(ctx context.Context, id int64) (*entity.InventoryAlert, error)
	List(ctx context.Context, productID *int64, isRead *bool, limit, offset int) ([]*entity.InventoryAlert, error)
	Update(ctx context.Context, alert *entity.InventoryAlert) error
	Delete(ctx context.Context, id int64) error
	MarkAsRead(ctx context.Context, id int64) error
	GetUnreadCount(ctx context.Context, productID *int64) (int, error)
}

// InventoryThresholdRepository defines the interface for inventory threshold data access
type InventoryThresholdRepository interface {
	Create(ctx context.Context, threshold *entity.InventoryThreshold) error
	GetByID(ctx context.Context, id int64) (*entity.InventoryThreshold, error)
	GetByProductID(ctx context.Context, productID int64) (*entity.InventoryThreshold, error)
	Update(ctx context.Context, threshold *entity.InventoryThreshold) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*entity.InventoryThreshold, error)
}
