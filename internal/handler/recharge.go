package handler

import (
"encoding/json"
"net/http"

rechargeService "erp-system/internal/domain/recharge/service"
)

// RechargeHandler handles recharge-related HTTP requests
type RechargeHandler struct {
rechargeService *rechargeService.RechargeService
}

// NewRechargeHandler creates a new recharge handler
func NewRechargeHandler(rechargeService *rechargeService.RechargeService) *RechargeHandler {
return &RechargeHandler{
rechargeService: rechargeService,
}
}

// ProcessRechargeRequest represents the request body for processing a recharge
type ProcessRechargeRequest struct {
UserID     int64   `json:"user_id"`
Amount     float64 `json:"amount"`
Method     string  `json:"method"`
Remark     string  `json:"remark"`
IsSupplier bool    `json:"is_supplier"`
}

// ProcessRecharge handles recharge processing
func (h *RechargeHandler) ProcessRecharge(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req ProcessRechargeRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Invalid request body: " + err.Error(),
})
return
}

// Validate input
if req.UserID <= 0 || req.Amount <= 0 {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "User ID and positive amount are required",
})
return
}

if req.Method == "" {
req.Method = "bank_transfer"
}

var err error
if req.IsSupplier {
err = h.rechargeService.ProcessSupplierRecharge(r.Context(), req.UserID, req.Amount, req.Method, req.Remark)
} else {
err = h.rechargeService.ProcessUserRecharge(r.Context(), req.UserID, req.Amount, req.Method, req.Remark)
}

if err != nil {
writeJSON(w, http.StatusInternalServerError, APIResponse{
Success: false,
Error:   "Failed to process recharge: " + err.Error(),
})
return
}

writeJSON(w, http.StatusOK, APIResponse{
Success: true,
Data:    map[string]string{"message": "Recharge processed successfully"},
})
}
