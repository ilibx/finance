package repository

import (
	"context"
	"database/sql"
	"time"

	"erp-system/internal/domain/product/entity"
	"erp-system/internal/domain/product/repository"
)

// productRepositoryImpl implements repository.ProductRepository
type productRepositoryImpl struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepositoryImpl{db: db}
}

// Create creates a new product
func (r *productRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {
	query := `INSERT INTO products (name, sku, price_amount, price_currency, cost_amount, cost_currency, 
		stock, description, status_code, status_description, status_updated_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`
	
	return r.db.QueryRowContext(ctx, query,
		product.Name,
		product.SKU,
		product.Price.Amount,
		product.Price.Currency,
		product.Cost.Amount,
		product.Cost.Currency,
		product.Stock,
		product.Description,
		product.Status.Code,
		product.Status.Description,
		product.Status.UpdatedAt,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)
}

// GetByID gets a product by ID
func (r *productRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	query := `SELECT id, name, sku, price_amount, price_currency, cost_amount, cost_currency,
		stock, description, status_code, status_description, status_updated_at, created_at, updated_at
		FROM products WHERE id = $1`
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	var product entity.Product
	var updatedAt time.Time
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Price.Amount,
		&product.Price.Currency,
		&product.Cost.Amount,
		&product.Cost.Currency,
		&product.Stock,
		&product.Description,
		&product.Status.Code,
		&product.Status.Description,
		&updatedAt,
		&product.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	product.Status.UpdatedAt = updatedAt
	product.UpdatedAt = updatedAt
	
	return &product, nil
}

// GetBySKU gets a product by SKU
func (r *productRepositoryImpl) GetBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	query := `SELECT id, name, sku, price_amount, price_currency, cost_amount, cost_currency,
		stock, description, status_code, status_description, status_updated_at, created_at, updated_at
		FROM products WHERE sku = $1`
	
	row := r.db.QueryRowContext(ctx, query, sku)
	
	var product entity.Product
	var updatedAt time.Time
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Price.Amount,
		&product.Price.Currency,
		&product.Cost.Amount,
		&product.Cost.Currency,
		&product.Stock,
		&product.Description,
		&product.Status.Code,
		&product.Status.Description,
		&updatedAt,
		&product.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	product.Status.UpdatedAt = updatedAt
	product.UpdatedAt = updatedAt
	
	return &product, nil
}

// Update updates a product
func (r *productRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	query := `UPDATE products SET name=$1, sku=$2, price_amount=$3, price_currency=$4,
		cost_amount=$5, cost_currency=$6, stock=$7, description=$8, status_code=$9,
		status_description=$10, status_updated_at=$11, updated_at=$12 WHERE id=$13`
	
	_, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.SKU,
		product.Price.Amount,
		product.Price.Currency,
		product.Cost.Amount,
		product.Cost.Currency,
		product.Stock,
		product.Description,
		product.Status.Code,
		product.Status.Description,
		product.Status.UpdatedAt,
		time.Now(),
		product.ID,
	)
	
	return err
}

// List lists products with pagination
func (r *productRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	query := `SELECT id, name, sku, price_amount, price_currency, cost_amount, cost_currency,
		stock, description, status_code, status_description, status_updated_at, created_at, updated_at
		FROM products ORDER BY id LIMIT $1 OFFSET $2`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		var updatedAt time.Time
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.SKU,
			&product.Price.Amount,
			&product.Price.Currency,
			&product.Cost.Amount,
			&product.Cost.Currency,
			&product.Stock,
			&product.Description,
			&product.Status.Code,
			&product.Status.Description,
			&updatedAt,
			&product.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		product.Status.UpdatedAt = updatedAt
		product.UpdatedAt = updatedAt
		products = append(products, &product)
	}
	
	return products, rows.Err()
}

// Ensure interface compliance
var _ repository.ProductRepository = (*productRepositoryImpl)(nil)
