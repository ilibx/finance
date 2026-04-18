package repository

import (
	"context"
	"database/sql"
	"time"

	"erp-system/internal/domain/supplier/entity"
	"erp-system/internal/domain/supplier/repository"
)

// supplierRepositoryImpl implements repository.SupplierRepository
type supplierRepositoryImpl struct {
	db *sql.DB
}

// NewSupplierRepository creates a new supplier repository
func NewSupplierRepository(db *sql.DB) repository.SupplierRepository {
	return &supplierRepositoryImpl{db: db}
}

// Create creates a new supplier
func (r *supplierRepositoryImpl) Create(ctx context.Context, supplier *entity.Supplier) error {
	query := `INSERT INTO suppliers (name, contact_name, contact_email, contact_phone, 
		address_street, address_city, address_state, address_country, address_zipcode,
		balance_amount, balance_currency, status_code, status_description, status_updated_at, 
		created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`
	
	return r.db.QueryRowContext(ctx, query,
		supplier.Name,
		supplier.Contact.Name,
		supplier.Contact.Email,
		supplier.Contact.Phone,
		supplier.Address.Street,
		supplier.Address.City,
		supplier.Address.State,
		supplier.Address.Country,
		supplier.Address.ZipCode,
		supplier.Balance.Amount,
		supplier.Balance.Currency,
		supplier.Status.Code,
		supplier.Status.Description,
		supplier.Status.UpdatedAt,
		supplier.CreatedAt,
		supplier.UpdatedAt,
	).Scan(&supplier.ID)
}

// GetByID gets a supplier by ID
func (r *supplierRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Supplier, error) {
	query := `SELECT id, name, contact_name, contact_email, contact_phone, 
		address_street, address_city, address_state, address_country, address_zipcode,
		balance_amount, balance_currency, status_code, status_description, status_updated_at, 
		created_at, updated_at
		FROM suppliers WHERE id = $1`
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	var supplier entity.Supplier
	var updatedAt time.Time
	
	err := row.Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.Contact.Name,
		&supplier.Contact.Email,
		&supplier.Contact.Phone,
		&supplier.Address.Street,
		&supplier.Address.City,
		&supplier.Address.State,
		&supplier.Address.Country,
		&supplier.Address.ZipCode,
		&supplier.Balance.Amount,
		&supplier.Balance.Currency,
		&supplier.Status.Code,
		&supplier.Status.Description,
		&updatedAt,
		&supplier.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	supplier.Status.UpdatedAt = updatedAt
	supplier.UpdatedAt = updatedAt
	
	return &supplier, nil
}

// Update updates a supplier
func (r *supplierRepositoryImpl) Update(ctx context.Context, supplier *entity.Supplier) error {
	query := `UPDATE suppliers SET name=$1, contact_name=$2, contact_email=$3, contact_phone=$4,
		address_street=$5, address_city=$6, address_state=$7, address_country=$8, address_zipcode=$9,
		balance_amount=$10, balance_currency=$11, status_code=$12, status_description=$13,
		status_updated_at=$14, updated_at=$15 WHERE id=$16`
	
	_, err := r.db.ExecContext(ctx, query,
		supplier.Name,
		supplier.Contact.Name,
		supplier.Contact.Email,
		supplier.Contact.Phone,
		supplier.Address.Street,
		supplier.Address.City,
		supplier.Address.State,
		supplier.Address.Country,
		supplier.Address.ZipCode,
		supplier.Balance.Amount,
		supplier.Balance.Currency,
		supplier.Status.Code,
		supplier.Status.Description,
		supplier.Status.UpdatedAt,
		time.Now(),
		supplier.ID,
	)
	
	return err
}

// List lists suppliers with pagination
func (r *supplierRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Supplier, error) {
	query := `SELECT id, name, contact_name, contact_email, contact_phone, 
		address_street, address_city, address_state, address_country, address_zipcode,
		balance_amount, balance_currency, status_code, status_description, status_updated_at, 
		created_at, updated_at
		FROM suppliers ORDER BY id LIMIT $1 OFFSET $2`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var suppliers []*entity.Supplier
	for rows.Next() {
		var supplier entity.Supplier
		var updatedAt time.Time
		
		err := rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.Contact.Name,
			&supplier.Contact.Email,
			&supplier.Contact.Phone,
			&supplier.Address.Street,
			&supplier.Address.City,
			&supplier.Address.State,
			&supplier.Address.Country,
			&supplier.Address.ZipCode,
			&supplier.Balance.Amount,
			&supplier.Balance.Currency,
			&supplier.Status.Code,
			&supplier.Status.Description,
			&updatedAt,
			&supplier.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		supplier.Status.UpdatedAt = updatedAt
		supplier.UpdatedAt = updatedAt
		suppliers = append(suppliers, &supplier)
	}
	
	return suppliers, rows.Err()
}

// Ensure interface compliance
var _ repository.SupplierRepository = (*supplierRepositoryImpl)(nil)
