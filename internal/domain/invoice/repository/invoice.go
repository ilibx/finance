package repository

import (
	"context"
	"database/sql"
	"time"

	"finance/internal/domain/invoice/entity"
)

// InvoiceRepository defines the interface for invoice data access
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *entity.Invoice) error
	GetByID(ctx context.Context, id int64) (*entity.Invoice, error)
	GetByInvoiceNo(ctx context.Context, invoiceNo string) (*entity.Invoice, error)
	Update(ctx context.Context, invoice *entity.Invoice) error
	List(ctx context.Context, limit, offset int) ([]*entity.Invoice, error)
}

// postgresInvoiceRepository implements InvoiceRepository
type postgresInvoiceRepository struct {
	db *sql.DB
}

// NewInvoiceRepository creates a new PostgreSQL invoice repository
func NewInvoiceRepository(db *sql.DB) InvoiceRepository {
	return &postgresInvoiceRepository{db: db}
}

func (r *postgresInvoiceRepository) Create(ctx context.Context, invoice *entity.Invoice) error {
	query := `
		INSERT INTO invoices (invoice_no, order_id, user_id, amount, tax_amount, total_amount, type, status, issued_at, paid_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	var issuedAt, paidAt interface{}
	if invoice.IssuedAt != nil {
		issuedAt = *invoice.IssuedAt
	} else {
		issuedAt = nil
	}
	if invoice.PaidAt != nil {
		paidAt = *invoice.PaidAt
	} else {
		paidAt = nil
	}
	return r.db.QueryRowContext(ctx, query,
		invoice.InvoiceNo,
		invoice.OrderID,
		invoice.UserID,
		invoice.Amount.Amount,
		invoice.TaxAmount.Amount,
		invoice.TotalAmount.Amount,
		string(invoice.Type),
		string(invoice.Status),
		issuedAt,
		paidAt,
		invoice.CreatedAt,
		invoice.UpdatedAt,
	).Scan(&invoice.ID)
}

func (r *postgresInvoiceRepository) GetByID(ctx context.Context, id int64) (*entity.Invoice, error) {
	query := `
		SELECT id, invoice_no, order_id, user_id, amount, tax_amount, total_amount, type, status, issued_at, paid_at, created_at, updated_at
		FROM invoices WHERE id = $1
	`
	invoice := &entity.Invoice{}
	var typeStr, statusStr string
	var issuedAt, paidAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&invoice.ID,
		&invoice.InvoiceNo,
		&invoice.OrderID,
		&invoice.UserID,
		&invoice.Amount.Amount,
		&invoice.TaxAmount.Amount,
		&invoice.TotalAmount.Amount,
		&typeStr,
		&statusStr,
		&issuedAt,
		&paidAt,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	invoice.Type = entity.InvoiceType(typeStr)
	invoice.Status = entity.InvoiceStatus(statusStr)
	invoice.Amount.Currency = "CNY"
	invoice.TaxAmount.Currency = "CNY"
	invoice.TotalAmount.Currency = "CNY"
	if issuedAt.Valid {
		invoice.IssuedAt = &issuedAt.Time
	}
	if paidAt.Valid {
		invoice.PaidAt = &paidAt.Time
	}
	return invoice, nil
}

func (r *postgresInvoiceRepository) GetByInvoiceNo(ctx context.Context, invoiceNo string) (*entity.Invoice, error) {
	query := `
		SELECT id, invoice_no, order_id, user_id, amount, tax_amount, total_amount, type, status, issued_at, paid_at, created_at, updated_at
		FROM invoices WHERE invoice_no = $1
	`
	invoice := &entity.Invoice{}
	var typeStr, statusStr string
	var issuedAt, paidAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, invoiceNo).Scan(
		&invoice.ID,
		&invoice.InvoiceNo,
		&invoice.OrderID,
		&invoice.UserID,
		&invoice.Amount.Amount,
		&invoice.TaxAmount.Amount,
		&invoice.TotalAmount.Amount,
		&typeStr,
		&statusStr,
		&issuedAt,
		&paidAt,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	invoice.Type = entity.InvoiceType(typeStr)
	invoice.Status = entity.InvoiceStatus(statusStr)
	invoice.Amount.Currency = "CNY"
	invoice.TaxAmount.Currency = "CNY"
	invoice.TotalAmount.Currency = "CNY"
	if issuedAt.Valid {
		invoice.IssuedAt = &issuedAt.Time
	}
	if paidAt.Valid {
		invoice.PaidAt = &paidAt.Time
	}
	return invoice, nil
}

func (r *postgresInvoiceRepository) Update(ctx context.Context, invoice *entity.Invoice) error {
	query := `
		UPDATE invoices SET
			order_id = $1, user_id = $2, amount = $3, tax_amount = $4, total_amount = $5,
			type = $6, status = $7, issued_at = $8, paid_at = $9, updated_at = $10
		WHERE id = $11
	`
	var issuedAt, paidAt interface{}
	if invoice.IssuedAt != nil {
		issuedAt = *invoice.IssuedAt
	} else {
		issuedAt = nil
	}
	if invoice.PaidAt != nil {
		paidAt = *invoice.PaidAt
	} else {
		paidAt = nil
	}
	_, err := r.db.ExecContext(ctx, query,
		invoice.OrderID,
		invoice.UserID,
		invoice.Amount.Amount,
		invoice.TaxAmount.Amount,
		invoice.TotalAmount.Amount,
		string(invoice.Type),
		string(invoice.Status),
		issuedAt,
		paidAt,
		time.Now(),
		invoice.ID,
	)
	return err
}

func (r *postgresInvoiceRepository) List(ctx context.Context, limit, offset int) ([]*entity.Invoice, error) {
	query := `
		SELECT id, invoice_no, order_id, user_id, amount, tax_amount, total_amount, type, status, issued_at, paid_at, created_at, updated_at
		FROM invoices ORDER BY id DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []*entity.Invoice
	for rows.Next() {
		invoice := &entity.Invoice{}
		var typeStr, statusStr string
		var issuedAt, paidAt sql.NullTime
		err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceNo,
			&invoice.OrderID,
			&invoice.UserID,
			&invoice.Amount.Amount,
			&invoice.TaxAmount.Amount,
			&invoice.TotalAmount.Amount,
			&typeStr,
			&statusStr,
			&issuedAt,
			&paidAt,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invoice.Type = entity.InvoiceType(typeStr)
		invoice.Status = entity.InvoiceStatus(statusStr)
		invoice.Amount.Currency = "CNY"
		invoice.TaxAmount.Currency = "CNY"
		invoice.TotalAmount.Currency = "CNY"
		if issuedAt.Valid {
			invoice.IssuedAt = &issuedAt.Time
		}
		if paidAt.Valid {
			invoice.PaidAt = &paidAt.Time
		}
		invoices = append(invoices, invoice)
	}
	return invoices, rows.Err()
}
