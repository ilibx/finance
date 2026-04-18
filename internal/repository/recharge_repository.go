package repository

import (
	"context"
	"database/sql"
	"time"

	"erp-system/internal/domain/recharge/entity"
	"erp-system/internal/domain/recharge/repository"
)

// rechargeRepositoryImpl implements repository.RechargeRepository
type rechargeRepositoryImpl struct {
	db *sql.DB
}

// NewRechargeRecordRepository creates a new recharge repository
func NewRechargeRecordRepository(db *sql.DB) repository.RechargeRepository {
	return &rechargeRepositoryImpl{db: db}
}

// Create creates a new recharge record
func (r *rechargeRepositoryImpl) Create(ctx context.Context, record *entity.RechargeRecord) error {
	query := `INSERT INTO recharge_records (user_id, amount_amount, amount_currency, type, method, 
		status, remark, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	
	return r.db.QueryRowContext(ctx, query,
		record.UserID,
		record.Amount.Amount,
		record.Amount.Currency,
		string(record.Type),
		string(record.Method),
		string(record.Status),
		record.Remark,
		record.CreatedAt,
		record.UpdatedAt,
	).Scan(&record.ID)
}

// GetByID gets a recharge record by ID
func (r *rechargeRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.RechargeRecord, error) {
	query := `SELECT id, user_id, amount_amount, amount_currency, type, method, status, remark, created_at, updated_at
		FROM recharge_records WHERE id = $1`
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	var record entity.RechargeRecord
	
	err := row.Scan(
		&record.ID,
		&record.UserID,
		&record.Amount.Amount,
		&record.Amount.Currency,
		&record.Type,
		&record.Method,
		&record.Status,
		&record.Remark,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return &record, nil
}

// Update updates a recharge record
func (r *rechargeRepositoryImpl) Update(ctx context.Context, record *entity.RechargeRecord) error {
	query := `UPDATE recharge_records SET user_id=$1, amount_amount=$2, amount_currency=$3,
		type=$4, method=$5, status=$6, remark=$7, updated_at=$8 WHERE id=$9`
	
	_, err := r.db.ExecContext(ctx, query,
		record.UserID,
		record.Amount.Amount,
		record.Amount.Currency,
		string(record.Type),
		string(record.Method),
		string(record.Status),
		record.Remark,
		time.Now(),
		record.ID,
	)
	
	return err
}

// ListByUserID lists recharge records by user ID with pagination
func (r *rechargeRepositoryImpl) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.RechargeRecord, error) {
	query := `SELECT id, user_id, amount_amount, amount_currency, type, method, status, remark, created_at, updated_at
		FROM recharge_records WHERE user_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*entity.RechargeRecord
	for rows.Next() {
		var record entity.RechargeRecord
		
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.Amount.Amount,
			&record.Amount.Currency,
			&record.Type,
			&record.Method,
			&record.Status,
			&record.Remark,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		records = append(records, &record)
	}
	
	return records, rows.Err()
}

// Ensure interface compliance
var _ repository.RechargeRepository = (*rechargeRepositoryImpl)(nil)
