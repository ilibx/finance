package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"finance/internal/domain/purchase/entity"
	purchaseService "finance/internal/domain/purchase/service"
)

// PurchaseHandler handles purchase-related HTTP requests
type PurchaseHandler struct {
	purchaseService *purchaseService.PurchaseService
}

// NewPurchaseHandler creates a new purchase handler
func NewPurchaseHandler(purchaseService *purchaseService.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseService: purchaseService,
	}
}

// CreatePurchaseRequest represents the request body for creating a purchase order
type CreatePurchaseRequest struct {
	SupplierID   int64                    `json:"supplier_id"`
	Items        []purchaseService.PurchaseItemRequest `json:"items"`
	CreatedBy    int64                    `json:"created_by"`
	Notes        string                   `json:"notes"`
	DeliveryDate *string                  `json:"delivery_date"`
}

// CreatePurchase handles purchase order creation
func (h *PurchaseHandler) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate input
	if req.SupplierID <= 0 || len(req.Items) == 0 {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Supplier ID and at least one item are required",
		})
		return
	}

	// Parse delivery date if provided
	var deliveryDate *time.Time
	if req.DeliveryDate != nil && *req.DeliveryDate != "" {
		t, err := time.Parse("2006-01-02", *req.DeliveryDate)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Invalid delivery date format. Use YYYY-MM-DD",
			})
			return
		}
		deliveryDate = &t
	}

	purchase, err := h.purchaseService.CreatePurchase(r.Context(), req.SupplierID, req.Items, req.CreatedBy, req.Notes, deliveryDate)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to create purchase order: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    purchase,
	})
}

// GetPurchase handles getting a purchase order by ID
func (h *PurchaseHandler) GetPurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Purchase ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid purchase ID",
		})
		return
	}

	purchase, err := h.purchaseService.GetPurchaseByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Purchase order not found: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    purchase,
	})
}

// SubmitPurchase handles submitting a purchase order for approval
func (h *PurchaseHandler) SubmitPurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Purchase ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid purchase ID",
		})
		return
	}

	if err := h.purchaseService.SubmitPurchase(r.Context(), id); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Failed to submit purchase order: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Purchase order submitted for approval"},
	})
}

// ApprovePurchase handles approving a purchase order
func (h *PurchaseHandler) ApprovePurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	approverStr := r.URL.Query().Get("approver_id")
	if idStr == "" || approverStr == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Purchase ID and approver ID are required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid purchase ID",
		})
		return
	}

	approverID, err := strconv.ParseInt(approverStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid approver ID",
		})
		return
	}

	if err := h.purchaseService.ApprovePurchase(r.Context(), id, approverID); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Failed to approve purchase order: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Purchase order approved"},
	})
}

// RejectPurchase handles rejecting a purchase order
func (h *PurchaseHandler) RejectPurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Purchase ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid purchase ID",
		})
		return
	}

	if err := h.purchaseService.RejectPurchase(r.Context(), id); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Failed to reject purchase order: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Purchase order rejected"},
	})
}

// OrderPurchase handles placing a purchase order to supplier
func (h *PurchaseHandler) OrderPurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Purchase ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid purchase ID",
		})
		return
	}

	if err := h.purchaseService.OrderPurchase(r.Context(), id); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Failed to place purchase order: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Purchase order placed with supplier"},
	})
}

// ReceivePurchase handles receiving items from a purchase order
func (h *PurchaseHandler) ReceivePurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PurchaseID int64 `json:"purchase_id"`
		ItemID     int64 `json:"item_id"`
		Quantity   int   `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.PurchaseID <= 0 || req.ItemID <= 0 || req.Quantity <= 0 {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Purchase ID, Item ID, and Quantity are required",
		})
		return
	}

	if err := h.purchaseService.ReceivePurchase(r.Context(), req.PurchaseID, req.ItemID, req.Quantity); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Failed to receive items: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Items received successfully"},
	})
}

// CancelPurchase handles cancelling a purchase order
func (h *PurchaseHandler) CancelPurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Purchase ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid purchase ID",
		})
		return
	}

	if err := h.purchaseService.CancelPurchase(r.Context(), id); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Failed to cancel purchase order: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Purchase order cancelled"},
	})
}

// ListPurchases handles listing purchase orders
func (h *PurchaseHandler) ListPurchases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	statusStr := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var status entity.PurchaseStatus
	if statusStr != "" {
		status = entity.PurchaseStatus(statusStr)
	}

	purchases, err := h.purchaseService.ListPurchases(r.Context(), status, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to list purchase orders: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    purchases,
	})
}

// DeletePurchase handles deleting a purchase order
func (h *PurchaseHandler) DeletePurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Purchase ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid purchase ID",
		})
		return
	}

	if err := h.purchaseService.DeletePurchase(r.Context(), id); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Failed to delete purchase order: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Purchase order deleted"},
	})
}
