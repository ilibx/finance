package repository

import (
	"context"
	"database/sql"
	"time"

	"erp-system/internal/domain/invoice/entity"
	"erp-system/internal/domain/invoice/repository"
)

// invoiceRepositoryImpl implements repository.InvoiceRepository
type invoiceRepositoryImpl struct {
	db *sql.DB
}

// NewInvoiceRepository creates a new invoice repository
func NewInvoiceRepository(db *sql.DB) repository.InvoiceRepository {
	return &invoiceRepositoryImpl{db: db}
}

// Create creates a new invoice
func (r *invoiceRepositoryImpl) Create(ctx context.Context, invoice *entity.Invoice) error {
	query := `INSERT INTO invoices (invoice_no, order_id, user_id, amount_amount, amount_currency,
		tax_amount_amount, tax_amount_currency, total_amount_amount, total_amount_currency,
		type, status, issued_at, paid_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`
	
	var issuedAt, paidAt sql.NullTime
	if invoice.IssuedAt != nil {
		issuedAt = sql.NullTime{Time: *invoice.IssuedAt, Valid: true}
	}
	if invoice.PaidAt != nil {
		paidAt = sql.NullTime{Time: *invoice.PaidAt, Valid: true}
	}
	
	return r.db.QueryRowContext(ctx, query,
		invoice.InvoiceNo,
		invoice.OrderID,
		invoice.UserID,
		invoice.Amount.Amount,
		invoice.Amount.Currency,
		invoice.TaxAmount.Amount,
		invoice.TaxAmount.Currency,
		invoice.TotalAmount.Amount,
		invoice.TotalAmount.Currency,
		string(invoice.Type),
		string(invoice.Status),
		issuedAt,
		paidAt,
		invoice.CreatedAt,
		invoice.UpdatedAt,
	).Scan(&invoice.ID)
}

// GetByID gets an invoice by ID
func (r *invoiceRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Invoice, error) {
	query := `SELECT id, invoice_no, order_id, user_id, amount_amount, amount_currency,
		tax_amount_amount, tax_amount_currency, total_amount_amount, total_amount_currency,
		type, status, issued_at, paid_at, created_at, updated_at
		FROM invoices WHERE id = $1`
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	var invoice entity.Invoice
	var issuedAt, paidAt sql.NullTime
	
	err := row.Scan(
		&invoice.ID,
		&invoice.InvoiceNo,
		&invoice.OrderID,
		&invoice.UserID,
		&invoice.Amount.Amount,
		&invoice.Amount.Currency,
		&invoice.TaxAmount.Amount,
		&invoice.TaxAmount.Currency,
		&invoice.TotalAmount.Amount,
		&invoice.TotalAmount.Currency,
		&invoice.Type,
		&invoice.Status,
		&issuedAt,
		&paidAt,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	if issuedAt.Valid {
		invoice.IssuedAt = &issuedAt.Time
	}
	if paidAt.Valid {
		invoice.PaidAt = &paidAt.Time
	}
	
	return &invoice, nil
}

// GetByInvoiceNo gets an invoice by invoice number
func (r *invoiceRepositoryImpl) GetByInvoiceNo(ctx context.Context, invoiceNo string) (*entity.Invoice, error) {
	query := `SELECT id, invoice_no, order_id, user_id, amount_amount, amount_currency,
		tax_amount_amount, tax_amount_currency, total_amount_amount, total_amount_currency,
		type, status, issued_at, paid_at, created_at, updated_at
		FROM invoices WHERE invoice_no = $1`
	
	row := r.db.QueryRowContext(ctx, query, invoiceNo)
	
	var invoice entity.Invoice
	var issuedAt, paidAt sql.NullTime
	
	err := row.Scan(
		&invoice.ID,
		&invoice.InvoiceNo,
		&invoice.OrderID,
		&invoice.UserID,
		&invoice.Amount.Amount,
		&invoice.Amount.Currency,
		&invoice.TaxAmount.Amount,
		&invoice.TaxAmount.Currency,
		&invoice.TotalAmount.Amount,
		&invoice.TotalAmount.Currency,
		&invoice.Type,
		&invoice.Status,
		&issuedAt,
		&paidAt,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	if issuedAt.Valid {
		invoice.IssuedAt = &issuedAt.Time
	}
	if paidAt.Valid {
		invoice.PaidAt = &paidAt.Time
	}
	
	return &invoice, nil
}

// Update updates an invoice
func (r *invoiceRepositoryImpl) Update(ctx context.Context, invoice *entity.Invoice) error {
	var issuedAt, paidAt sql.NullTime
	if invoice.IssuedAt != nil {
		issuedAt = sql.NullTime{Time: *invoice.IssuedAt, Valid: true}
	}
	if invoice.PaidAt != nil {
		paidAt = sql.NullTime{Time: *invoice.PaidAt, Valid: true}
	}
	
	query := `UPDATE invoices SET invoice_no=$1, order_id=$2, user_id=$3,
		amount_amount=$4, amount_currency=$5, tax_amount_amount=$6, tax_amount_currency=$7,
		total_amount_amount=$8, total_amount_currency=$9, type=$10, status=$11,
		issued_at=$12, paid_at=$13, updated_at=$14 WHERE id=$15`
	
	_, err := r.db.ExecContext(ctx, query,
		invoice.InvoiceNo,
		invoice.OrderID,
		invoice.UserID,
		invoice.Amount.Amount,
		invoice.Amount.Currency,
		invoice.TaxAmount.Amount,
		invoice.TaxAmount.Currency,
		invoice.TotalAmount.Amount,
		invoice.TotalAmount.Currency,
		string(invoice.Type),
		string(invoice.Status),
		issuedAt,
		paidAt,
		time.Now(),
		invoice.ID,
	)
	
	return err
}

// List lists invoices with pagination
func (r *invoiceRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Invoice, error) {
	query := `SELECT id, invoice_no, order_id, user_id, amount_amount, amount_currency,
		tax_amount_amount, tax_amount_currency, total_amount_amount, total_amount_currency,
		type, status, issued_at, paid_at, created_at, updated_at
		FROM invoices ORDER BY id LIMIT $1 OFFSET $2`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var invoices []*entity.Invoice
	for rows.Next() {
		var invoice entity.Invoice
		var issuedAt, paidAt sql.NullTime
		
		err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceNo,
			&invoice.OrderID,
			&invoice.UserID,
			&invoice.Amount.Amount,
			&invoice.Amount.Currency,
			&invoice.TaxAmount.Amount,
			&invoice.TaxAmount.Currency,
			&invoice.TotalAmount.Amount,
			&invoice.TotalAmount.Currency,
			&invoice.Type,
			&invoice.Status,
			&issuedAt,
			&paidAt,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if issuedAt.Valid {
			invoice.IssuedAt = &issuedAt.Time
		}
		if paidAt.Valid {
			invoice.PaidAt = &paidAt.Time
		}
		
		invoices = append(invoices, &invoice)
	}
	
	return invoices, rows.Err()
}

// Ensure interface compliance
var _ repository.InvoiceRepository = (*invoiceRepositoryImpl)(nil)
