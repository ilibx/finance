package repository

import (
	"context"
	"database/sql"
	"time"

	"finance/internal/domain/inventory/entity"
)

// postgresInventoryAlertRepository implements InventoryAlertRepository
type postgresInventoryAlertRepository struct {
	db *sql.DB
}

// NewInventoryAlertRepository creates a new PostgreSQL inventory alert repository
func NewInventoryAlertRepository(db *sql.DB) InventoryAlertRepository {
	return &postgresInventoryAlertRepository{db: db}
}

func (r *postgresInventoryAlertRepository) Create(ctx context.Context, alert *entity.InventoryAlert) error {
	query := `
		INSERT INTO inventory_alerts 
		(product_id, product_name, product_sku, alert_type, alert_level, current_stock, threshold_value, message, is_read, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		alert.ProductID,
		alert.ProductName,
		alert.ProductSKU,
		alert.AlertType,
		alert.AlertLevel,
		alert.CurrentStock,
		alert.ThresholdValue,
		alert.Message,
		alert.IsRead,
		alert.CreatedAt,
		alert.UpdatedAt,
	).Scan(&alert.ID)
}

func (r *postgresInventoryAlertRepository) GetByID(ctx context.Context, id int64) (*entity.InventoryAlert, error) {
	query := `
		SELECT id, product_id, product_name, product_sku, alert_type, alert_level, current_stock, threshold_value, message, is_read, created_at, updated_at
		FROM inventory_alerts WHERE id = $1
	`
	alert := &entity.InventoryAlert{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&alert.ID,
		&alert.ProductID,
		&alert.ProductName,
		&alert.ProductSKU,
		&alert.AlertType,
		&alert.AlertLevel,
		&alert.CurrentStock,
		&alert.ThresholdValue,
		&alert.Message,
		&alert.IsRead,
		&alert.CreatedAt,
		&alert.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return alert, nil
}

func (r *postgresInventoryAlertRepository) List(ctx context.Context, productID *int64, isRead *bool, limit, offset int) ([]*entity.InventoryAlert, error) {
	query := `
		SELECT id, product_id, product_name, product_sku, alert_type, alert_level, current_stock, threshold_value, message, is_read, created_at, updated_at
		FROM inventory_alerts
		WHERE ($1::bigint IS NULL OR product_id = $1)
		  AND ($2::boolean IS NULL OR is_read = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(ctx, query, productID, isRead, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*entity.InventoryAlert
	for rows.Next() {
		alert := &entity.InventoryAlert{}
		err := rows.Scan(
			&alert.ID,
			&alert.ProductID,
			&alert.ProductName,
			&alert.ProductSKU,
			&alert.AlertType,
			&alert.AlertLevel,
			&alert.CurrentStock,
			&alert.ThresholdValue,
			&alert.Message,
			&alert.IsRead,
			&alert.CreatedAt,
			&alert.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, rows.Err()
}

func (r *postgresInventoryAlertRepository) Update(ctx context.Context, alert *entity.InventoryAlert) error {
	query := `
		UPDATE inventory_alerts SET
			product_id = $1, product_name = $2, product_sku = $3, alert_type = $4, alert_level = $5,
			current_stock = $6, threshold_value = $7, message = $8, is_read = $9, updated_at = $10
		WHERE id = $11
	`
	_, err := r.db.ExecContext(ctx, query,
		alert.ProductID,
		alert.ProductName,
		alert.ProductSKU,
		alert.AlertType,
		alert.AlertLevel,
		alert.CurrentStock,
		alert.ThresholdValue,
		alert.Message,
		alert.IsRead,
		time.Now(),
		alert.ID,
	)
	return err
}

func (r *postgresInventoryAlertRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM inventory_alerts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresInventoryAlertRepository) MarkAsRead(ctx context.Context, id int64) error {
	query := `UPDATE inventory_alerts SET is_read = true, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *postgresInventoryAlertRepository) GetUnreadCount(ctx context.Context, productID *int64) (int, error) {
	query := `SELECT COUNT(*) FROM inventory_alerts WHERE is_read = false AND ($1::bigint IS NULL OR product_id = $1)`
	var count int
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&count)
	return count, err
}

// postgresInventoryThresholdRepository implements InventoryThresholdRepository
type postgresInventoryThresholdRepository struct {
	db *sql.DB
}

// NewInventoryThresholdRepository creates a new PostgreSQL inventory threshold repository
func NewInventoryThresholdRepository(db *sql.DB) InventoryThresholdRepository {
	return &postgresInventoryThresholdRepository{db: db}
}

func (r *postgresInventoryThresholdRepository) Create(ctx context.Context, threshold *entity.InventoryThreshold) error {
	query := `
		INSERT INTO inventory_thresholds 
		(product_id, min_stock, max_stock, safety_stock, reorder_point, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		threshold.ProductID,
		threshold.MinStock,
		threshold.MaxStock,
		threshold.SafetyStock,
		threshold.ReorderPoint,
		threshold.Enabled,
		threshold.CreatedAt,
		threshold.UpdatedAt,
	).Scan(&threshold.ID)
}

func (r *postgresInventoryThresholdRepository) GetByID(ctx context.Context, id int64) (*entity.InventoryThreshold, error) {
	query := `
		SELECT id, product_id, min_stock, max_stock, safety_stock, reorder_point, enabled, created_at, updated_at
		FROM inventory_thresholds WHERE id = $1
	`
	threshold := &entity.InventoryThreshold{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&threshold.ID,
		&threshold.ProductID,
		&threshold.MinStock,
		&threshold.MaxStock,
		&threshold.SafetyStock,
		&threshold.ReorderPoint,
		&threshold.Enabled,
		&threshold.CreatedAt,
		&threshold.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return threshold, nil
}

func (r *postgresInventoryThresholdRepository) GetByProductID(ctx context.Context, productID int64) (*entity.InventoryThreshold, error) {
	query := `
		SELECT id, product_id, min_stock, max_stock, safety_stock, reorder_point, enabled, created_at, updated_at
		FROM inventory_thresholds WHERE product_id = $1
	`
	threshold := &entity.InventoryThreshold{}
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&threshold.ID,
		&threshold.ProductID,
		&threshold.MinStock,
		&threshold.MaxStock,
		&threshold.SafetyStock,
		&threshold.ReorderPoint,
		&threshold.Enabled,
		&threshold.CreatedAt,
		&threshold.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return threshold, nil
}

func (r *postgresInventoryThresholdRepository) Update(ctx context.Context, threshold *entity.InventoryThreshold) error {
	query := `
		UPDATE inventory_thresholds SET
			min_stock = $1, max_stock = $2, safety_stock = $3, reorder_point = $4, enabled = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		threshold.MinStock,
		threshold.MaxStock,
		threshold.SafetyStock,
		threshold.ReorderPoint,
		threshold.Enabled,
		time.Now(),
		threshold.ID,
	)
	return err
}

func (r *postgresInventoryThresholdRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM inventory_thresholds WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresInventoryThresholdRepository) List(ctx context.Context, limit, offset int) ([]*entity.InventoryThreshold, error) {
	query := `
		SELECT id, product_id, min_stock, max_stock, safety_stock, reorder_point, enabled, created_at, updated_at
		FROM inventory_thresholds ORDER BY id DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var thresholds []*entity.InventoryThreshold
	for rows.Next() {
		threshold := &entity.InventoryThreshold{}
		err := rows.Scan(
			&threshold.ID,
			&threshold.ProductID,
			&threshold.MinStock,
			&threshold.MaxStock,
			&threshold.SafetyStock,
			&threshold.ReorderPoint,
			&threshold.Enabled,
			&threshold.CreatedAt,
			&threshold.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		thresholds = append(thresholds, threshold)
	}
	return thresholds, rows.Err()
}
