package repository

import (
	"context"
	"database/sql"
	"time"

	"finance/internal/domain/purchase/entity"
)

// PurchaseRepository defines the interface for purchase order data access
type PurchaseRepository interface {
	Create(ctx context.Context, purchase *entity.Purchase) error
	GetByID(ctx context.Context, id int64) (*entity.Purchase, error)
	GetByPurchaseNo(ctx context.Context, purchaseNo string) (*entity.Purchase, error)
	Update(ctx context.Context, purchase *entity.Purchase) error
	List(ctx context.Context, status entity.PurchaseStatus, limit, offset int) ([]*entity.Purchase, error)
	ListBySupplier(ctx context.Context, supplierID int64, limit, offset int) ([]*entity.Purchase, error)
	Delete(ctx context.Context, id int64) error
}

// postgresPurchaseRepository implements PurchaseRepository
type postgresPurchaseRepository struct {
	db *sql.DB
}

// NewPurchaseRepository creates a new PostgreSQL purchase repository
func NewPurchaseRepository(db *sql.DB) PurchaseRepository {
	return &postgresPurchaseRepository{db: db}
}

func (r *postgresPurchaseRepository) Create(ctx context.Context, purchase *entity.Purchase) error {
	query := `
		INSERT INTO purchase_orders 
		(purchase_no, supplier_id, supplier_name, total_amount, status, created_by, notes, delivery_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	var deliveryDate interface{}
	if purchase.DeliveryDate != nil {
		deliveryDate = *purchase.DeliveryDate
	} else {
		deliveryDate = nil
	}

	err := r.db.QueryRowContext(ctx, query,
		purchase.PurchaseNo,
		purchase.SupplierID,
		purchase.SupplierName,
		purchase.TotalAmount.Amount,
		string(purchase.Status),
		purchase.CreatedBy,
		purchase.Notes,
		deliveryDate,
		purchase.CreatedAt,
		purchase.UpdatedAt,
	).Scan(&purchase.ID)
	if err != nil {
		return err
	}

	// Insert purchase items
	for _, item := range purchase.Items {
		itemQuery := `
			INSERT INTO purchase_order_items 
			(purchase_id, product_id, product_name, quantity, received_qty, unit_price, subtotal, supplier_part_no)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		err := r.db.QueryRowContext(ctx, itemQuery,
			purchase.ID,
			item.ProductID,
			item.ProductName,
			item.Quantity,
			item.ReceivedQty,
			item.UnitPrice.Amount,
			item.Subtotal.Amount,
			item.SupplierPartNo,
		).Scan(&item.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *postgresPurchaseRepository) GetByID(ctx context.Context, id int64) (*entity.Purchase, error) {
	query := `
		SELECT id, purchase_no, supplier_id, supplier_name, total_amount, status, created_by, 
		       approved_by, approved_at, ordered_at, completed_at, notes, delivery_date, received_date, created_at, updated_at
		FROM purchase_orders WHERE id = $1
	`
	purchase := &entity.Purchase{}
	var statusStr string
	var approvedBy sql.NullInt64
	var approvedAt, orderedAt, completedAt, deliveryDate, receivedDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&purchase.ID,
		&purchase.PurchaseNo,
		&purchase.SupplierID,
		&purchase.SupplierName,
		&purchase.TotalAmount.Amount,
		&statusStr,
		&purchase.CreatedBy,
		&approvedBy,
		&approvedAt,
		&orderedAt,
		&completedAt,
		&purchase.Notes,
		&deliveryDate,
		&receivedDate,
		&purchase.CreatedAt,
		&purchase.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	purchase.Status = entity.PurchaseStatus(statusStr)
	purchase.TotalAmount.Currency = "CNY"
	if approvedBy.Valid {
		purchase.ApprovedBy = &approvedBy.Int64
	}
	if approvedAt.Valid {
		purchase.ApprovedAt = &approvedAt.Time
	}
	if orderedAt.Valid {
		purchase.OrderedAt = &orderedAt.Time
	}
	if completedAt.Valid {
		purchase.CompletedAt = &completedAt.Time
	}
	if deliveryDate.Valid {
		purchase.DeliveryDate = &deliveryDate.Time
	}
	if receivedDate.Valid {
		purchase.ReceivedDate = &receivedDate.Time
	}

	// Load purchase items
	items, err := r.getPurchaseItems(ctx, id)
	if err != nil {
		return nil, err
	}
	purchase.Items = items

	return purchase, nil
}

func (r *postgresPurchaseRepository) GetByPurchaseNo(ctx context.Context, purchaseNo string) (*entity.Purchase, error) {
	query := `
		SELECT id, purchase_no, supplier_id, supplier_name, total_amount, status, created_by, 
		       approved_by, approved_at, ordered_at, completed_at, notes, delivery_date, received_date, created_at, updated_at
		FROM purchase_orders WHERE purchase_no = $1
	`
	purchase := &entity.Purchase{}
	var statusStr string
	var approvedBy sql.NullInt64
	var approvedAt, orderedAt, completedAt, deliveryDate, receivedDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, purchaseNo).Scan(
		&purchase.ID,
		&purchase.PurchaseNo,
		&purchase.SupplierID,
		&purchase.SupplierName,
		&purchase.TotalAmount.Amount,
		&statusStr,
		&purchase.CreatedBy,
		&approvedBy,
		&approvedAt,
		&orderedAt,
		&completedAt,
		&purchase.Notes,
		&deliveryDate,
		&receivedDate,
		&purchase.CreatedAt,
		&purchase.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	purchase.Status = entity.PurchaseStatus(statusStr)
	purchase.TotalAmount.Currency = "CNY"
	if approvedBy.Valid {
		purchase.ApprovedBy = &approvedBy.Int64
	}
	if approvedAt.Valid {
		purchase.ApprovedAt = &approvedAt.Time
	}
	if orderedAt.Valid {
		purchase.OrderedAt = &orderedAt.Time
	}
	if completedAt.Valid {
		purchase.CompletedAt = &completedAt.Time
	}
	if deliveryDate.Valid {
		purchase.DeliveryDate = &deliveryDate.Time
	}
	if receivedDate.Valid {
		purchase.ReceivedDate = &receivedDate.Time
	}

	// Load purchase items
	items, err := r.getPurchaseItems(ctx, purchase.ID)
	if err != nil {
		return nil, err
	}
	purchase.Items = items

	return purchase, nil
}

func (r *postgresPurchaseRepository) Update(ctx context.Context, purchase *entity.Purchase) error {
	query := `
		UPDATE purchase_orders SET
			supplier_id = $1, supplier_name = $2, total_amount = $3, status = $4,
			approved_by = $5, approved_at = $6, ordered_at = $7, completed_at = $8,
			notes = $9, delivery_date = $10, received_date = $11, updated_at = $12
		WHERE id = $13
	`
	var approvedBy, orderedAt, completedAt, deliveryDate, receivedDate interface{}
	if purchase.ApprovedBy != nil {
		approvedBy = *purchase.ApprovedBy
	} else {
		approvedBy = nil
	}
	if purchase.ApprovedAt != nil {
		orderedAt = *purchase.ApprovedAt
	} else {
		orderedAt = nil
	}
	if purchase.OrderedAt != nil {
		completedAt = *purchase.OrderedAt
	} else {
		completedAt = nil
	}
	if purchase.CompletedAt != nil {
		deliveryDate = *purchase.CompletedAt
	} else {
		deliveryDate = nil
	}
	if purchase.DeliveryDate != nil {
		receivedDate = *purchase.DeliveryDate
	} else {
		receivedDate = nil
	}

	_, err := r.db.ExecContext(ctx, query,
		purchase.SupplierID,
		purchase.SupplierName,
		purchase.TotalAmount.Amount,
		string(purchase.Status),
		approvedBy,
		orderedAt,
		completedAt,
		deliveryDate,
		purchase.Notes,
		receivedDate,
		time.Now(),
		purchase.ID,
	)
	if err != nil {
		return err
	}

	// Update purchase items
	for _, item := range purchase.Items {
		itemQuery := `
			UPDATE purchase_order_items SET
				product_id = $1, product_name = $2, quantity = $3, received_qty = $4,
				unit_price = $5, subtotal = $6, supplier_part_no = $7
			WHERE id = $8 AND purchase_id = $9
		`
		_, err := r.db.ExecContext(ctx, itemQuery,
			item.ProductID,
			item.ProductName,
			item.Quantity,
			item.ReceivedQty,
			item.UnitPrice.Amount,
			item.Subtotal.Amount,
			item.SupplierPartNo,
			item.ID,
			purchase.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *postgresPurchaseRepository) List(ctx context.Context, status entity.PurchaseStatus, limit, offset int) ([]*entity.Purchase, error) {
	query := `
		SELECT id, purchase_no, supplier_id, supplier_name, total_amount, status, created_by, 
		       approved_by, approved_at, ordered_at, completed_at, notes, delivery_date, received_date, created_at, updated_at
		FROM purchase_orders
	`
	args := []interface{}{}
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, string(status))
		query += " ORDER BY id DESC LIMIT $" + string(rune(len(args)+1)) + " OFFSET $" + string(rune(len(args)+2))
	} else {
		query += " ORDER BY id DESC LIMIT $" + string(rune(len(args)+1)) + " OFFSET $" + string(rune(len(args)+2))
	}
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*entity.Purchase
	for rows.Next() {
		purchase := &entity.Purchase{}
		var statusStr string
		var approvedBy sql.NullInt64
		var approvedAt, orderedAt, completedAt, deliveryDate, receivedDate sql.NullTime

		scanArgs := []interface{}{
			&purchase.ID,
			&purchase.PurchaseNo,
			&purchase.SupplierID,
			&purchase.SupplierName,
			&purchase.TotalAmount.Amount,
			&statusStr,
			&purchase.CreatedBy,
			&approvedBy,
			&approvedAt,
			&orderedAt,
			&completedAt,
			&purchase.Notes,
			&deliveryDate,
			&receivedDate,
			&purchase.CreatedAt,
			&purchase.UpdatedAt,
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		purchase.Status = entity.PurchaseStatus(statusStr)
		purchase.TotalAmount.Currency = "CNY"
		if approvedBy.Valid {
			purchase.ApprovedBy = &approvedBy.Int64
		}
		if approvedAt.Valid {
			purchase.ApprovedAt = &approvedAt.Time
		}
		if orderedAt.Valid {
			purchase.OrderedAt = &orderedAt.Time
		}
		if completedAt.Valid {
			purchase.CompletedAt = &completedAt.Time
		}
		if deliveryDate.Valid {
			purchase.DeliveryDate = &deliveryDate.Time
		}
		if receivedDate.Valid {
			purchase.ReceivedDate = &receivedDate.Time
		}

		// Load purchase items
		items, err := r.getPurchaseItems(ctx, purchase.ID)
		if err != nil {
			return nil, err
		}
		purchase.Items = items

		purchases = append(purchases, purchase)
	}
	return purchases, rows.Err()
}

func (r *postgresPurchaseRepository) ListBySupplier(ctx context.Context, supplierID int64, limit, offset int) ([]*entity.Purchase, error) {
	query := `
		SELECT id, purchase_no, supplier_id, supplier_name, total_amount, status, created_by, 
		       approved_by, approved_at, ordered_at, completed_at, notes, delivery_date, received_date, created_at, updated_at
		FROM purchase_orders WHERE supplier_id = $1
		ORDER BY id DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, supplierID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*entity.Purchase
	for rows.Next() {
		purchase := &entity.Purchase{}
		var statusStr string
		var approvedBy sql.NullInt64
		var approvedAt, orderedAt, completedAt, deliveryDate, receivedDate sql.NullTime

		err := rows.Scan(
			&purchase.ID,
			&purchase.PurchaseNo,
			&purchase.SupplierID,
			&purchase.SupplierName,
			&purchase.TotalAmount.Amount,
			&statusStr,
			&purchase.CreatedBy,
			&approvedBy,
			&approvedAt,
			&orderedAt,
			&completedAt,
			&purchase.Notes,
			&deliveryDate,
			&receivedDate,
			&purchase.CreatedAt,
			&purchase.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		purchase.Status = entity.PurchaseStatus(statusStr)
		purchase.TotalAmount.Currency = "CNY"
		if approvedBy.Valid {
			purchase.ApprovedBy = &approvedBy.Int64
		}
		if approvedAt.Valid {
			purchase.ApprovedAt = &approvedAt.Time
		}
		if orderedAt.Valid {
			purchase.OrderedAt = &orderedAt.Time
		}
		if completedAt.Valid {
			purchase.CompletedAt = &completedAt.Time
		}
		if deliveryDate.Valid {
			purchase.DeliveryDate = &deliveryDate.Time
		}
		if receivedDate.Valid {
			purchase.ReceivedDate = &receivedDate.Time
		}

		// Load purchase items
		items, err := r.getPurchaseItems(ctx, purchase.ID)
		if err != nil {
			return nil, err
		}
		purchase.Items = items

		purchases = append(purchases, purchase)
	}
	return purchases, rows.Err()
}

func (r *postgresPurchaseRepository) Delete(ctx context.Context, id int64) error {
	// First delete items
	_, err := r.db.ExecContext(ctx, "DELETE FROM purchase_order_items WHERE purchase_id = $1", id)
	if err != nil {
		return err
	}
	// Then delete the purchase order
	_, err = r.db.ExecContext(ctx, "DELETE FROM purchase_orders WHERE id = $1", id)
	return err
}

func (r *postgresPurchaseRepository) getPurchaseItems(ctx context.Context, purchaseID int64) ([]entity.PurchaseItem, error) {
	query := `
		SELECT id, product_id, product_name, quantity, received_qty, unit_price, subtotal, supplier_part_no
		FROM purchase_order_items WHERE purchase_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.PurchaseItem
	for rows.Next() {
		item := entity.PurchaseItem{}
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.ReceivedQty,
			&item.UnitPrice.Amount,
			&item.Subtotal.Amount,
			&item.SupplierPartNo,
		)
		if err != nil {
			return nil, err
		}
		item.UnitPrice.Currency = "CNY"
		item.Subtotal.Currency = "CNY"
		items = append(items, item)
	}
	return items, rows.Err()
}
