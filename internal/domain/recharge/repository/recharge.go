package repository

import (
	"context"
	"database/sql"
	"time"

	"finance/internal/domain/recharge/entity"
)

// RechargeRepository defines the interface for recharge data access
type RechargeRepository interface {
	Create(ctx context.Context, record *entity.RechargeRecord) error
	GetByID(ctx context.Context, id int64) (*entity.RechargeRecord, error)
	Update(ctx context.Context, record *entity.RechargeRecord) error
	ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.RechargeRecord, error)
}

// postgresRechargeRepository implements RechargeRepository
type postgresRechargeRepository struct {
	db *sql.DB
}

// NewRechargeRepository creates a new PostgreSQL recharge repository
func NewRechargeRepository(db *sql.DB) RechargeRepository {
	return &postgresRechargeRepository{db: db}
}

func (r *postgresRechargeRepository) Create(ctx context.Context, record *entity.RechargeRecord) error {
	query := `
		INSERT INTO recharge_records (user_id, amount, type, method, status, remark, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		record.UserID,
		record.Amount.Amount,
		string(record.Type),
		string(record.Method),
		string(record.Status),
		record.Remark,
		record.CreatedAt,
		record.UpdatedAt,
	).Scan(&record.ID)
}

func (r *postgresRechargeRepository) GetByID(ctx context.Context, id int64) (*entity.RechargeRecord, error) {
	query := `
		SELECT id, user_id, amount, type, method, status, remark, created_at, updated_at
		FROM recharge_records WHERE id = $1
	`
	record := &entity.RechargeRecord{}
	var typeStr, methodStr, statusStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.UserID,
		&record.Amount.Amount,
		&typeStr,
		&methodStr,
		&statusStr,
		&record.Remark,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	record.Type = entity.RechargeType(typeStr)
	record.Method = entity.RechargeMethod(methodStr)
	record.Status = entity.RechargeStatus(statusStr)
	record.Amount.Currency = "CNY"
	return record, nil
}

func (r *postgresRechargeRepository) Update(ctx context.Context, record *entity.RechargeRecord) error {
	query := `
		UPDATE recharge_records SET
			user_id = $1, amount = $2, type = $3, method = $4,
			status = $5, remark = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.ExecContext(ctx, query,
		record.UserID,
		record.Amount.Amount,
		string(record.Type),
		string(record.Method),
		string(record.Status),
		record.Remark,
		time.Now(),
		record.ID,
	)
	return err
}

func (r *postgresRechargeRepository) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.RechargeRecord, error) {
	query := `
		SELECT id, user_id, amount, type, method, status, remark, created_at, updated_at
		FROM recharge_records WHERE user_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*entity.RechargeRecord
	for rows.Next() {
		record := &entity.RechargeRecord{}
		var typeStr, methodStr, statusStr string
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.Amount.Amount,
			&typeStr,
			&methodStr,
			&statusStr,
			&record.Remark,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		record.Type = entity.RechargeType(typeStr)
		record.Method = entity.RechargeMethod(methodStr)
		record.Status = entity.RechargeStatus(statusStr)
		record.Amount.Currency = "CNY"
		records = append(records, record)
	}
	return records, rows.Err()
}
