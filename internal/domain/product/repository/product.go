package repository

import (
	"context"
	"database/sql"
	"time"

	"erp-system/internal/domain/product/entity"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id int64) (*entity.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	List(ctx context.Context, limit, offset int) ([]*entity.Product, error)
}

// postgresProductRepository implements ProductRepository
type postgresProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new PostgreSQL product repository
func NewProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(ctx context.Context, product *entity.Product) error {
	query := `
		INSERT INTO products (name, sku, price, cost, stock, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		product.Name,
		product.SKU,
		product.Price,
		product.Cost,
		product.Stock,
		product.Description,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)
}

func (r *postgresProductRepository) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	query := `
		SELECT id, name, sku, price, cost, stock, description, created_at, updated_at
		FROM products WHERE id = $1
	`
	product := &entity.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Price,
		&product.Cost,
		&product.Stock,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *postgresProductRepository) GetBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	query := `
		SELECT id, name, sku, price, cost, stock, description, created_at, updated_at
		FROM products WHERE sku = $1
	`
	product := &entity.Product{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Price,
		&product.Cost,
		&product.Stock,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *postgresProductRepository) Update(ctx context.Context, product *entity.Product) error {
	query := `
		UPDATE products SET
			name = $1, sku = $2, price = $3, cost = $4, stock = $5,
			description = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.SKU,
		product.Price,
		product.Cost,
		product.Stock,
		product.Description,
		time.Now(),
		product.ID,
	)
	return err
}

func (r *postgresProductRepository) List(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	query := `
		SELECT id, name, sku, price, cost, stock, description, created_at, updated_at
		FROM products ORDER BY id DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		product := &entity.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.SKU,
			&product.Price,
			&product.Cost,
			&product.Stock,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}
