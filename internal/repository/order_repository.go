package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"erp-system/internal/common/valueobject"
	"erp-system/internal/domain/order/entity"
	"erp-system/internal/domain/order/repository"
)

// orderRepositoryImpl implements repository.OrderRepository
type orderRepositoryImpl struct {
	db *sql.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &orderRepositoryImpl{db: db}
}

// Create creates a new order
func (r *orderRepositoryImpl) Create(ctx context.Context, order *entity.Order) error {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}
	
	query := `INSERT INTO orders (order_no, user_id, items_json, total_amount, total_currency, 
		status, paid_at, shipped_at, completed_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	
	var paidAt, shippedAt, completedAt sql.NullTime
	if order.PaidAt != nil {
		paidAt = sql.NullTime{Time: *order.PaidAt, Valid: true}
	}
	if order.ShippedAt != nil {
		shippedAt = sql.NullTime{Time: *order.ShippedAt, Valid: true}
	}
	if order.CompletedAt != nil {
		completedAt = sql.NullTime{Time: *order.CompletedAt, Valid: true}
	}
	
	return r.db.QueryRowContext(ctx, query,
		order.OrderNo,
		order.UserID,
		string(itemsJSON),
		order.TotalAmount.Amount,
		order.TotalAmount.Currency,
		string(order.Status),
		paidAt,
		shippedAt,
		completedAt,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID)
}

// GetByID gets an order by ID
func (r *orderRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Order, error) {
	query := `SELECT id, order_no, user_id, items_json, total_amount, total_currency,
		status, paid_at, shipped_at, completed_at, created_at, updated_at
		FROM orders WHERE id = $1`
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	var order entity.Order
	var itemsJSON string
	var paidAt, shippedAt, completedAt sql.NullTime
	
	err := row.Scan(
		&order.ID,
		&order.OrderNo,
		&order.UserID,
		&itemsJSON,
		&order.TotalAmount.Amount,
		&order.TotalAmount.Currency,
		&order.Status,
		&paidAt,
		&shippedAt,
		&completedAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	if err := json.Unmarshal([]byte(itemsJSON), &order.Items); err != nil {
		return nil, err
	}
	
	if paidAt.Valid {
		order.PaidAt = &paidAt.Time
	}
	if shippedAt.Valid {
		order.ShippedAt = &shippedAt.Time
	}
	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}
	
	return &order, nil
}

// GetByOrderNo gets an order by order number
func (r *orderRepositoryImpl) GetByOrderNo(ctx context.Context, orderNo string) (*entity.Order, error) {
	query := `SELECT id, order_no, user_id, items_json, total_amount, total_currency,
		status, paid_at, shipped_at, completed_at, created_at, updated_at
		FROM orders WHERE order_no = $1`
	
	row := r.db.QueryRowContext(ctx, query, orderNo)
	
	var order entity.Order
	var itemsJSON string
	var paidAt, shippedAt, completedAt sql.NullTime
	
	err := row.Scan(
		&order.ID,
		&order.OrderNo,
		&order.UserID,
		&itemsJSON,
		&order.TotalAmount.Amount,
		&order.TotalAmount.Currency,
		&order.Status,
		&paidAt,
		&shippedAt,
		&completedAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	if err := json.Unmarshal([]byte(itemsJSON), &order.Items); err != nil {
		return nil, err
	}
	
	if paidAt.Valid {
		order.PaidAt = &paidAt.Time
	}
	if shippedAt.Valid {
		order.ShippedAt = &shippedAt.Time
	}
	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}
	
	return &order, nil
}

// Update updates an order
func (r *orderRepositoryImpl) Update(ctx context.Context, order *entity.Order) error {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}
	
	var paidAt, shippedAt, completedAt sql.NullTime
	if order.PaidAt != nil {
		paidAt = sql.NullTime{Time: *order.PaidAt, Valid: true}
	}
	if order.ShippedAt != nil {
		shippedAt = sql.NullTime{Time: *order.ShippedAt, Valid: true}
	}
	if order.CompletedAt != nil {
		completedAt = sql.NullTime{Time: *order.CompletedAt, Valid: true}
	}
	
	query := `UPDATE orders SET order_no=$1, user_id=$2, items_json=$3, total_amount=$4,
		total_currency=$5, status=$6, paid_at=$7, shipped_at=$8, completed_at=$9,
		updated_at=$10 WHERE id=$11`
	
	_, err = r.db.ExecContext(ctx, query,
		order.OrderNo,
		order.UserID,
		string(itemsJSON),
		order.TotalAmount.Amount,
		order.TotalAmount.Currency,
		string(order.Status),
		paidAt,
		shippedAt,
		completedAt,
		time.Now(),
		order.ID,
	)
	
	return err
}

// List lists orders with pagination
func (r *orderRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Order, error) {
	query := `SELECT id, order_no, user_id, items_json, total_amount, total_currency,
		status, paid_at, shipped_at, completed_at, created_at, updated_at
		FROM orders ORDER BY id LIMIT $1 OFFSET $2`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var orders []*entity.Order
	for rows.Next() {
		var order entity.Order
		var itemsJSON string
		var paidAt, shippedAt, completedAt sql.NullTime
		
		err := rows.Scan(
			&order.ID,
			&order.OrderNo,
			&order.UserID,
			&itemsJSON,
			&order.TotalAmount.Amount,
			&order.TotalAmount.Currency,
			&order.Status,
			&paidAt,
			&shippedAt,
			&completedAt,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if err := json.Unmarshal([]byte(itemsJSON), &order.Items); err != nil {
			return nil, err
		}
		
		if paidAt.Valid {
			order.PaidAt = &paidAt.Time
		}
		if shippedAt.Valid {
			order.ShippedAt = &shippedAt.Time
		}
		if completedAt.Valid {
			order.CompletedAt = &completedAt.Time
		}
		
		orders = append(orders, &order)
	}
	
	return orders, rows.Err()
}

// Ensure interface compliance
var _ repository.OrderRepository = (*orderRepositoryImpl)(nil)
