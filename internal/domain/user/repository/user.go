package repository

import (
	"context"

	"erp-system/internal/domain/user/entity"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
}
