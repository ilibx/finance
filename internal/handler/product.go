package handler

import (
"encoding/json"
"net/http"

productService "finance/internal/domain/product/service"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
productService *productService.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *productService.ProductService) *ProductHandler {
return &ProductHandler{
productService: productService,
}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
Name        string  `json:"name"`
SKU         string  `json:"sku"`
Price       float64 `json:"price"`
Cost        float64 `json:"cost"`
Stock       int     `json:"stock"`
Description string  `json:"description"`
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req CreateProductRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Invalid request body: " + err.Error(),
})
return
}

// Validate input
if req.Name == "" || req.SKU == "" {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Name and SKU are required",
})
return
}

product, err := h.productService.CreateProduct(r.Context(), req.Name, req.SKU, req.Price, req.Cost, req.Stock, req.Description)
if err != nil {
writeJSON(w, http.StatusInternalServerError, APIResponse{
Success: false,
Error:   "Failed to create product: " + err.Error(),
})
return
}

writeJSON(w, http.StatusOK, APIResponse{
Success: true,
Data:    product,
})
}
