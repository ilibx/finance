package entity

import (
	"time"

	"erp-system/internal/common/errors"
	"erp-system/internal/common/valueobject"
)

// Product represents a product aggregate root
type Product struct {
	ID          int64
	Name        string
	SKU         string
	Price       valueobject.Money
	Cost        valueobject.Money
	Stock       int
	Description string
	Status      valueobject.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProduct creates a new product
func NewProduct(name, sku string, price, cost float64, stock int, description string) *Product {
	now := time.Now()
	return &Product{
		Name:        name,
		SKU:         sku,
		Price:       valueobject.NewMoney(price, "CNY"),
		Cost:        valueobject.NewMoney(cost, "CNY"),
		Stock:       stock,
		Description: description,
		Status: valueobject.Status{
			Code:        "active",
			Description: "Active",
			UpdatedAt:   now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateStock updates product stock
func (p *Product) UpdateStock(quantity int) error {
	if p.Stock < quantity {
		return errors.ErrInsufficientStock
	}
	p.Stock -= quantity
	p.UpdatedAt = time.Now()
	return nil
}

// AddStock adds stock to product
func (p *Product) AddStock(quantity int) {
	p.Stock += quantity
	p.UpdatedAt = time.Now()
}

// IsActive checks if product is active
func (p *Product) IsActive() bool {
	return p.Status.Code == "active"
}
