package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ContextKey defines context keys
type ContextKey string

const (
	UserIDKey   ContextKey = "user_id"
	UsernameKey ContextKey = "username"
	RoleIDKey   ContextKey = "role_id"
	RoleCodeKey ContextKey = "role_code"
)

var jwtKey = []byte("erp_system_secret_key_2024_change_in_production")

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	secretKey string
}

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	RoleID   uint   `json:"role_id"`
	RoleCode string `json:"role_code"`
	jwt.RegisteredClaims
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: secretKey,
	}
}

// GenerateToken generates a JWT token
func (m *AuthMiddleware) GenerateToken(userID int64, username string, roleID uint, roleCode string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &TokenClaims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		RoleCode: roleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (m *AuthMiddleware) ValidateToken(tokenString string) (*TokenClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken generates a new token with extended expiration
func (m *AuthMiddleware) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same claims
	return m.GenerateToken(claims.UserID, claims.Username, claims.RoleID, claims.RoleCode)
}

// Middleware returns HTTP middleware for authentication
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Authorization header required",
			})
			return
		}

		claims, err := m.ValidateToken(authHeader)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Invalid or expired token",
			})
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, RoleIDKey, claims.RoleID)
		ctx = context.WithValue(ctx, RoleCodeKey, claims.RoleCode)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequirePermission returns middleware that checks for specific permissions
func (m *AuthMiddleware) RequirePermission(requiredPermissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleCode, ok := r.Context().Value(RoleCodeKey).(string)
			if !ok || roleCode == "" {
				writeJSON(w, http.StatusForbidden, map[string]interface{}{
					"success": false,
					"error":   "Role information missing",
				})
				return
			}

			// Admin has all permissions
			if roleCode != "admin" {
				// Check permissions for non-admin users
				if !checkPermissions(roleCode, requiredPermissions...) {
					writeJSON(w, http.StatusForbidden, map[string]interface{}{
						"success": false,
						"error":   "Insufficient permissions",
					})
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func checkPermissions(roleCode string, permissions ...string) bool {
	// Simplified permission check
	// In production, query database for role permissions
	rolePermissions := getRolePermissions(roleCode)

	for _, perm := range permissions {
		found := false
		for _, rp := range rolePermissions {
			if rp == perm || rp == "*" {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func getRolePermissions(roleCode string) []string {
	// Mock permissions by role
	permissions := map[string][]string{
		"admin":   {"*"},
		"manager": {"user:read", "user:create", "order:*", "product:*"},
		"user":    {"user:read", "order:read", "product:read"},
	}

	if perms, ok := permissions[roleCode]; ok {
		return perms
	}
	return []string{}
}

// Errors
var (
	ErrTokenExpired = &AuthError{Message: "Token has expired"}
	ErrInvalidToken = &AuthError{Message: "Invalid token"}
)

// AuthError represents an authentication error
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
