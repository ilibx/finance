package repository

import (
"context"

"erp-system/internal/domain/recharge/entity"
)

// RechargeRepository defines the interface for recharge data access
type RechargeRepository interface {
Create(ctx context.Context, record *entity.RechargeRecord) error
GetByID(ctx context.Context, id int64) (*entity.RechargeRecord, error)
Update(ctx context.Context, record *entity.RechargeRecord) error
ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.RechargeRecord, error)
}
