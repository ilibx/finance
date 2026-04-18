package repository

import (
	"context"
	"database/sql"
	"time"

	"finance/internal/domain/supplier/entity"
)

// SupplierRepository defines the interface for supplier data access
type SupplierRepository interface {
	Create(ctx context.Context, supplier *entity.Supplier) error
	GetByID(ctx context.Context, id int64) (*entity.Supplier, error)
	Update(ctx context.Context, supplier *entity.Supplier) error
	List(ctx context.Context, limit, offset int) ([]*entity.Supplier, error)
	Delete(ctx context.Context, id int64) error
}

// postgresSupplierRepository implements SupplierRepository
type postgresSupplierRepository struct {
	db *sql.DB
}

// NewSupplierRepository creates a new PostgreSQL supplier repository
func NewSupplierRepository(db *sql.DB) SupplierRepository {
	return &postgresSupplierRepository{db: db}
}

func (r *postgresSupplierRepository) Create(ctx context.Context, supplier *entity.Supplier) error {
	query := `
		INSERT INTO suppliers (name, contact_name, contact_phone, contact_email, address, balance_amount, status_code, status_description, status_updated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		supplier.Name,
		supplier.Contact.Name,
		supplier.Contact.Phone,
		supplier.Contact.Email,
		supplier.Address.String(),
		supplier.Balance.Amount,
		supplier.Status.Code,
		supplier.Status.Description,
		supplier.Status.UpdatedAt,
		supplier.CreatedAt,
		supplier.UpdatedAt,
	).Scan(&supplier.ID)
}

func (r *postgresSupplierRepository) GetByID(ctx context.Context, id int64) (*entity.Supplier, error) {
	query := `
		SELECT id, name, contact_name, contact_phone, contact_email, address, balance_amount, status_code, status_description, status_updated_at, created_at, updated_at
		FROM suppliers WHERE id = $1
	`
	supplier := &entity.Supplier{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.Contact.Name,
		&supplier.Contact.Phone,
		&supplier.Contact.Email,
		&supplier.Address.Street,
		&supplier.Balance.Amount,
		&supplier.Status.Code,
		&supplier.Status.Description,
		&supplier.Status.UpdatedAt,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	supplier.Balance.Currency = "CNY"
	return supplier, nil
}

func (r *postgresSupplierRepository) Update(ctx context.Context, supplier *entity.Supplier) error {
	query := `
		UPDATE suppliers SET
			name = $1, contact_name = $2, contact_phone = $3, contact_email = $4, address = $5,
			balance_amount = $6, status_code = $7, status_description = $8, status_updated_at = $9, updated_at = $10
		WHERE id = $11
	`
	_, err := r.db.ExecContext(ctx, query,
		supplier.Name,
		supplier.Contact.Name,
		supplier.Contact.Phone,
		supplier.Contact.Email,
		supplier.Address.String(),
		supplier.Balance.Amount,
		supplier.Status.Code,
		supplier.Status.Description,
		supplier.Status.UpdatedAt,
		time.Now(),
		supplier.ID,
	)
	return err
}

func (r *postgresSupplierRepository) List(ctx context.Context, limit, offset int) ([]*entity.Supplier, error) {
	query := `
		SELECT id, name, contact_name, contact_phone, contact_email, address, balance_amount, status_code, status_description, status_updated_at, created_at, updated_at
		FROM suppliers ORDER BY id DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []*entity.Supplier
	for rows.Next() {
		supplier := &entity.Supplier{}
		err := rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.Contact.Name,
			&supplier.Contact.Phone,
			&supplier.Contact.Email,
			&supplier.Address.Street,
			&supplier.Balance.Amount,
			&supplier.Status.Code,
			&supplier.Status.Description,
			&supplier.Status.UpdatedAt,
			&supplier.CreatedAt,
			&supplier.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		supplier.Balance.Currency = "CNY"
		suppliers = append(suppliers, supplier)
	}
	return suppliers, rows.Err()
}

func (r *postgresSupplierRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM suppliers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
