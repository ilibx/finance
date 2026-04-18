package service

import (
"context"
"fmt"
"time"

"erp-system/internal/common/errors"
"erp-system/internal/domain/order/entity"
"erp-system/internal/domain/order/repository"
"erp-system/internal/domain/product/repository"
)

// OrderService handles order business logic
type OrderService struct {
orderRepo   repository.OrderRepository
productRepo repository.ProductRepository
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *OrderService {
return &OrderService{
orderRepo:   orderRepo,
productRepo: productRepo,
}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, userID int64, items []entity.OrderItem) (*entity.Order, error) {
orderNo := generateOrderNo()
order := entity.NewOrder(orderNo, userID, items)

if err := s.orderRepo.Create(ctx, order); err != nil {
return nil, err
}

return order, nil
}

// GetOrderByID gets an order by ID
func (s *OrderService) GetOrderByID(ctx context.Context, id int64) (*entity.Order, error) {
order, err := s.orderRepo.GetByID(ctx, id)
if err != nil {
return nil, errors.ErrOrderNotFound
}
return order, nil
}

// PayOrder pays an order
func (s *OrderService) PayOrder(ctx context.Context, orderID int64) error {
order, err := s.orderRepo.GetByID(ctx, orderID)
if err != nil {
return errors.ErrOrderNotFound
}

if err := order.Pay(); err != nil {
return err
}

return s.orderRepo.Update(ctx, order)
}

// ListOrders lists orders with pagination
func (s *OrderService) ListOrders(ctx context.Context, limit, offset int) ([]*entity.Order, error) {
return s.orderRepo.List(ctx, limit, offset)
}

func generateOrderNo() string {
return fmt.Sprintf("ORD%s", time.Now().Format("20060102150405"))
}
