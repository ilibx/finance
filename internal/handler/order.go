package handler

import (
	"encoding/json"
	"net/http"

	"erp-system/internal/common/valueobject"
	"erp-system/internal/domain/order/entity"
	orderService "erp-system/internal/domain/order/service"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderService *orderService.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService *orderService.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	UserID int64              `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}

// OrderItemRequest represents an order item in the request
type OrderItemRequest struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

// CreateOrder handles order creation
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate input
	if req.UserID <= 0 || len(req.Items) == 0 {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "User ID and at least one item are required",
		})
		return
	}

	// Convert request items to entity items
	items := make([]entity.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		unitPrice := valueobject.NewMoney(item.UnitPrice, "CNY")
		subtotal := valueobject.NewMoney(item.UnitPrice*float64(item.Quantity), "CNY")
		orderItem := entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: unitPrice,
			Subtotal:  subtotal,
		}
		items = append(items, orderItem)
	}

	order, err := h.orderService.CreateOrder(r.Context(), req.UserID, items)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to create order: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    order,
	})
}
