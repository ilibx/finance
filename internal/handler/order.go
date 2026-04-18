package handler

import (
	"encoding/json"
	"net/http"

	"erp-system/internal/service"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	UserID int64                  `json:"user_id"`
	Items  []service.OrderItemInput `json:"items"`
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

	order, err := h.orderService.CreateOrder(r.Context(), req.UserID, req.Items)
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
