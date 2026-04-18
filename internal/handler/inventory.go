package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	inventoryService "finance/internal/domain/inventory/service"
)

// InventoryHandler handles inventory alert-related HTTP requests
type InventoryHandler struct {
	inventoryService *inventoryService.InventoryAlertService
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(inventoryService *inventoryService.InventoryAlertService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// SetThresholdRequest represents the request body for setting inventory threshold
type SetThresholdRequest struct {
	ProductID    int64 `json:"product_id"`
	MinStock     int   `json:"min_stock"`
	MaxStock     int   `json:"max_stock"`
	SafetyStock  int   `json:"safety_stock"`
	ReorderPoint int   `json:"reorder_point"`
}

// SetThreshold handles setting inventory threshold
func (h *InventoryHandler) SetThreshold(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SetThresholdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate input
	if req.ProductID <= 0 {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Product ID must be positive",
		})
		return
	}

	threshold, err := h.inventoryService.SetThreshold(r.Context(), req.ProductID, req.MinStock, req.MaxStock, req.SafetyStock, req.ReorderPoint)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to set threshold: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    threshold,
	})
}

// GetThreshold handles getting inventory threshold for a product
func (h *InventoryHandler) GetThreshold(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productIDStr := r.URL.Query().Get("product_id")
	if productIDStr == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Product ID is required",
		})
		return
	}

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid product ID",
		})
		return
	}

	threshold, err := h.inventoryService.GetThreshold(r.Context(), productID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to get threshold: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    threshold,
	})
}

// ListAlerts handles listing inventory alerts
func (h *InventoryHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	var productID *int64
	if pid := r.URL.Query().Get("product_id"); pid != "" {
		id, err := strconv.ParseInt(pid, 10, 64)
		if err == nil {
			productID = &id
		}
	}

	var isRead *bool
	if ir := r.URL.Query().Get("is_read"); ir != "" {
		b := ir == "true"
		isRead = &b
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	alerts, err := h.inventoryService.ListAlerts(r.Context(), productID, isRead, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to list alerts: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    alerts,
	})
}

// MarkAlertAsRead handles marking an alert as read
func (h *InventoryHandler) MarkAlertAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	alertIDStr := r.URL.Query().Get("alert_id")
	if alertIDStr == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Alert ID is required",
		})
		return
	}

	alertID, err := strconv.ParseInt(alertIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid alert ID",
		})
		return
	}

	if err := h.inventoryService.MarkAlertAsRead(r.Context(), alertID); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to mark alert as read: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"status": "marked_as_read"},
	})
}

// GetUnreadCount handles getting unread alert count
func (h *InventoryHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var productID *int64
	if pid := r.URL.Query().Get("product_id"); pid != "" {
		id, err := strconv.ParseInt(pid, 10, 64)
		if err == nil {
			productID = &id
		}
	}

	count, err := h.inventoryService.GetUnreadAlertCount(r.Context(), productID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to get unread count: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]int{"unread_count": count},
	})
}

// CheckAllProducts handles manual trigger of stock check for all products
func (h *InventoryHandler) CheckAllProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.inventoryService.CheckAllProducts(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to check all products: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"status": "check_completed"},
	})
}

// ListThresholds handles listing all inventory thresholds
func (h *InventoryHandler) ListThresholds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	thresholds, err := h.inventoryService.ListThresholds(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to list thresholds: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    thresholds,
	})
}
