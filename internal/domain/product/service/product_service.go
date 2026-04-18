package service

import (
"context"

"erp-system/internal/common/errors"
"erp-system/internal/domain/product/entity"
"erp-system/internal/domain/product/repository"
)

// ProductService handles product business logic
type ProductService struct {
productRepo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo repository.ProductRepository) *ProductService {
return &ProductService{
productRepo: productRepo,
}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, name, sku string, price, cost float64, stock int, description string) (*entity.Product, error) {
// Check if product already exists
existingProduct, err := s.productRepo.GetBySKU(ctx, sku)
if err == nil && existingProduct != nil {
return nil, errors.ErrProductAlreadyExists
}

product := entity.NewProduct(name, sku, price, cost, stock, description)
if err := s.productRepo.Create(ctx, product); err != nil {
return nil, err
}

return product, nil
}

// GetProductByID gets a product by ID
func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*entity.Product, error) {
product, err := s.productRepo.GetByID(ctx, id)
if err != nil {
return nil, errors.ErrProductNotFound
}
return product, nil
}

// ListProducts lists products with pagination
func (s *ProductService) ListProducts(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
return s.productRepo.List(ctx, limit, offset)
}

// UpdateStock updates product stock
func (s *ProductService) UpdateStock(ctx context.Context, productID, quantity int) error {
product, err := s.productRepo.GetByID(ctx, int64(productID))
if err != nil {
return errors.ErrProductNotFound
}

if err := product.UpdateStock(quantity); err != nil {
return err
}

return s.productRepo.Update(ctx, product)
}
