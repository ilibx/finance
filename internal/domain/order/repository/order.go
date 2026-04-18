package repository

import (
	"context"
	"database/sql"
	"time"

	"finance/internal/domain/order/entity"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id int64) (*entity.Order, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	List(ctx context.Context, limit, offset int) ([]*entity.Order, error)
}

// postgresOrderRepository implements OrderRepository
type postgresOrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new PostgreSQL order repository
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	query := `
		INSERT INTO orders (order_no, user_id, total_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		order.OrderNo,
		order.UserID,
		order.TotalAmount.Amount,
		string(order.Status),
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID)
}

func (r *postgresOrderRepository) GetByID(ctx context.Context, id int64) (*entity.Order, error) {
	query := `
		SELECT id, order_no, user_id, total_amount, status, paid_at, shipped_at, completed_at, created_at, updated_at
		FROM orders WHERE id = $1
	`
	order := &entity.Order{}
	var statusStr string
	var paidAt, shippedAt, completedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.OrderNo,
		&order.UserID,
		&order.TotalAmount.Amount,
		&statusStr,
		&paidAt,
		&shippedAt,
		&completedAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	order.Status = entity.OrderStatus(statusStr)
	order.TotalAmount.Currency = "CNY"
	if paidAt.Valid {
		order.PaidAt = &paidAt.Time
	}
	if shippedAt.Valid {
		order.ShippedAt = &shippedAt.Time
	}
	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}
	return order, nil
}

func (r *postgresOrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*entity.Order, error) {
	query := `
		SELECT id, order_no, user_id, total_amount, status, paid_at, shipped_at, completed_at, created_at, updated_at
		FROM orders WHERE order_no = $1
	`
	order := &entity.Order{}
	var statusStr string
	var paidAt, shippedAt, completedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, orderNo).Scan(
		&order.ID,
		&order.OrderNo,
		&order.UserID,
		&order.TotalAmount.Amount,
		&statusStr,
		&paidAt,
		&shippedAt,
		&completedAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	order.Status = entity.OrderStatus(statusStr)
	order.TotalAmount.Currency = "CNY"
	if paidAt.Valid {
		order.PaidAt = &paidAt.Time
	}
	if shippedAt.Valid {
		order.ShippedAt = &shippedAt.Time
	}
	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}
	return order, nil
}

func (r *postgresOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	query := `
		UPDATE orders SET
			user_id = $1, total_amount = $2, status = $3,
			paid_at = $4, shipped_at = $5, completed_at = $6, updated_at = $7
		WHERE id = $8
	`
	var paidAt, shippedAt, completedAt interface{}
	if order.PaidAt != nil {
		paidAt = *order.PaidAt
	} else {
		paidAt = nil
	}
	if order.ShippedAt != nil {
		shippedAt = *order.ShippedAt
	} else {
		shippedAt = nil
	}
	if order.CompletedAt != nil {
		completedAt = *order.CompletedAt
	} else {
		completedAt = nil
	}
	_, err := r.db.ExecContext(ctx, query,
		order.UserID,
		order.TotalAmount.Amount,
		string(order.Status),
		paidAt,
		shippedAt,
		completedAt,
		time.Now(),
		order.ID,
	)
	return err
}

func (r *postgresOrderRepository) List(ctx context.Context, limit, offset int) ([]*entity.Order, error) {
	query := `
		SELECT id, order_no, user_id, total_amount, status, paid_at, shipped_at, completed_at, created_at, updated_at
		FROM orders ORDER BY id DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		order := &entity.Order{}
		var statusStr string
		var paidAt, shippedAt, completedAt sql.NullTime
		err := rows.Scan(
			&order.ID,
			&order.OrderNo,
			&order.UserID,
			&order.TotalAmount.Amount,
			&statusStr,
			&paidAt,
			&shippedAt,
			&completedAt,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		order.Status = entity.OrderStatus(statusStr)
		order.TotalAmount.Currency = "CNY"
		if paidAt.Valid {
			order.PaidAt = &paidAt.Time
		}
		if shippedAt.Valid {
			order.ShippedAt = &shippedAt.Time
		}
		if completedAt.Valid {
			order.CompletedAt = &completedAt.Time
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}
