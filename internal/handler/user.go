package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	userService "finance/internal/domain/user/service"
	"finance/internal/middleware"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *userService.UserService
	auth        *middleware.AuthMiddleware
	db          *sql.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *userService.UserService, auth *middleware.AuthMiddleware, db *sql.DB) *UserHandler {
	return &UserHandler{
		userService: userService,
		auth:        auth,
		db:          db,
	}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	RoleID   uint   `json:"role_id"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Login handles user login and returns JWT token
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Username and password are required",
		})
		return
	}

	// Query user from database
	query := `SELECT id, username, password, role_id FROM users WHERE username = $1 AND deleted_at IS NULL`
	var userID int64
	var username string
	var hashedPassword string
	var roleID uint

	err := h.db.QueryRow(query, req.Username).Scan(&userID, &username, &hashedPassword, &roleID)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid username or password",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database error: " + err.Error(),
		})
		return
	}

	// Verify password using bcrypt
	if err := checkPassword(hashedPassword, req.Password); err != nil {
		writeJSON(w, http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid username or password",
		})
		return
	}

	// Get role code
	roleCode := "user"
	roleQuery := `SELECT code FROM roles WHERE id = $1`
	h.db.QueryRow(roleQuery, roleID).Scan(&roleCode)

	// Generate JWT token
	token, err := h.auth.GenerateToken(userID, username, roleID, roleCode)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":       userID,
				"username": username,
				"role_id":  roleID,
				"role":     roleCode,
			},
		},
	})
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeJSON(w, http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Authorization header required",
		})
		return
	}

	newToken, err := h.auth.RefreshToken(authHeader)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "Invalid or expired token",
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"token": newToken,
		},
	})
}

// Logout handles user logout (client should remove token)
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT system, logout is handled client-side by removing the token
	// Optionally, you can maintain a blacklist of revoked tokens in Redis
	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Logged out successfully"},
	})
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

	ctx := context.WithValue(r.Context(), middleware.UserIDKey, int64(1)) // Default admin ID for now
	user, err := h.userService.CreateUser(ctx, req.Username, req.Email, req.Phone)
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

// ListUsers handles listing all users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `SELECT id, username, email, phone, nickname, balance_amount, status_code, role_id, created_at, updated_at 
			  FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC`
	
	rows, err := h.db.Query(query)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database error: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	type User struct {
		ID        int64      `json:"id"`
		Username  string     `json:"username"`
		Email     string     `json:"email"`
		Phone     string     `json:"phone"`
		Nickname  string     `json:"nickname"`
		Balance   float64    `json:"balance"`
		StatusCode string    `json:"status"`
		RoleID    uint       `json:"role_id"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
	}

	users := []User{}
	for rows.Next() {
		var u User
		var nullableNickname, nullablePhone sql.NullString
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &nullablePhone, &nullableNickname, 
			&u.Balance, &u.StatusCode, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "Error scanning user: " + err.Error(),
			})
			return
		}
		u.Phone = nullablePhone.String
		u.Nickname = nullableNickname.String
		users = append(users, u)
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    users,
	})
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// checkPassword compares a plain text password with a hashed password
func checkPassword(hashedPassword, password string) error {
	// Import bcrypt here or use the one from user entity
	return userService.CheckPasswordHash(hashedPassword, password)
}
