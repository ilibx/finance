package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// ContextKey defines context keys
type ContextKey string

const (
	UserIDKey   ContextKey = "user_id"
	UsernameKey ContextKey = "username"
	RoleIDKey   ContextKey = "role_id"
	RoleCodeKey ContextKey = "role_code"
)

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
	ExpireAt int64  `json:"expire_at"`
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: secretKey,
	}
}

// GenerateToken generates a JWT token
func (m *AuthMiddleware) GenerateToken(userID int64, username string, roleID uint, roleCode string) (string, error) {
	claims := TokenClaims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		RoleCode: roleCode,
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	
	// Simple token for demo (replace with real JWT in production)
	token := generateToken(claims, m.secretKey)
	return token, nil
}

// ValidateToken validates a JWT token
func (m *AuthMiddleware) ValidateToken(tokenString string) (*TokenClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	
	claims, err := validateToken(tokenString, m.secretKey)
	if err != nil {
		return nil, err
	}
	
	if time.Now().Unix() > claims.ExpireAt {
		return nil, ErrTokenExpired
	}
	
	return claims, nil
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
			if rp == perm {
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

// Helper functions
func generateToken(claims TokenClaims, secret string) string {
	// Simplified token generation - replace with proper JWT in production
	return "token_" + claims.Username + "_" + string(rune(claims.UserID))
}

func validateToken(token string, secret string) (*TokenClaims, error) {
	// Simplified validation - replace with proper JWT validation in production
	if !strings.HasPrefix(token, "token_") {
		return nil, ErrInvalidToken
	}
	
	// Extract username from token (simplified)
	parts := strings.Split(strings.TrimPrefix(token, "token_"), "_")
	if len(parts) < 1 {
		return nil, ErrInvalidToken
	}
	
	return &TokenClaims{
		UserID:   1,
		Username: parts[0],
		RoleID:   1,
		RoleCode: "admin",
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	// Simple JSON encoding
	if data == nil {
		w.Write([]byte("{}"))
		return
	}
	
	// Use simple formatting for demo
	switch v := data.(type) {
	case map[string]interface{}:
		w.Write([]byte(`{"success":`))
		if success, ok := v["success"].(bool); ok && success {
			w.Write([]byte(`true`))
		} else {
			w.Write([]byte(`false`))
		}
		if errMsg, ok := v["error"].(string); ok {
			w.Write([]byte(`,"error":"` + errMsg + `"`))
		}
		w.Write([]byte(`}`))
	}
}
