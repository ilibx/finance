package handler

import (
"encoding/json"
"fmt"
"net/http"

userService "finance/internal/domain/user/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
userService *userService.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *userService.UserService) *UserHandler {
return &UserHandler{
userService: userService,
}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
Username string `json:"username"`
Email    string `json:"email"`
Phone    string `json:"phone"`
}

// APIResponse represents a standard API response
type APIResponse struct {
Success bool        `json:"success"`
Data    interface{} `json:"data,omitempty"`
Error   string      `json:"error,omitempty"`
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req CreateUserRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Invalid request body: " + err.Error(),
})
return
}

// Validate input
if req.Username == "" || req.Email == "" {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Username and email are required",
})
return
}

user, err := h.userService.CreateUser(r.Context(), req.Username, req.Email, req.Phone)
if err != nil {
writeJSON(w, http.StatusInternalServerError, APIResponse{
Success: false,
Error:   "Failed to create user: " + err.Error(),
})
return
}

writeJSON(w, http.StatusOK, APIResponse{
Success: true,
Data:    user,
})
}

// GetUser handles getting a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodGet {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

idStr := r.URL.Query().Get("id")
if idStr == "" {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "User ID is required",
})
return
}

var id int64
if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
writeJSON(w, http.StatusBadRequest, APIResponse{
Success: false,
Error:   "Invalid user ID",
})
return
}

user, err := h.userService.GetUserByID(r.Context(), id)
if err != nil {
writeJSON(w, http.StatusNotFound, APIResponse{
Success: false,
Error:   "User not found",
})
return
}

writeJSON(w, http.StatusOK, APIResponse{
Success: true,
Data:    user,
})
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(status)
json.NewEncoder(w).Encode(data)
}
