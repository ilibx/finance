package handler

import (
	"encoding/json"
	"net/http"

	"erp-system/internal/service"
)

// InvoiceHandler handles invoice-related HTTP requests
type InvoiceHandler struct {
	invoiceService *service.InvoiceService
	orderService   *service.OrderService
}

// NewInvoiceHandler creates a new invoice handler
func NewInvoiceHandler(invoiceService *service.InvoiceService, orderService *service.OrderService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
		orderService:   orderService,
	}
}

// GenerateInvoiceRequest represents the request body for generating an invoice
type GenerateInvoiceRequest struct {
	OrderID int64   `json:"order_id"`
	TaxRate float64 `json:"tax_rate"`
}

// GenerateInvoice handles invoice generation
func (h *InvoiceHandler) GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate input
	if req.OrderID <= 0 {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Order ID is required",
		})
		return
	}

	// Get order first to get user ID and amount
	order, err := h.orderService.GetOrderByID(r.Context(), req.OrderID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to get order: " + err.Error(),
		})
		return
	}

	invoice, err := h.invoiceService.CreateInvoice(r.Context(), req.OrderID, order.UserID, order.TotalAmount, req.TaxRate)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to generate invoice: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    invoice,
	})
}
