package middleware

import (
	"context"
	"net/http"
	"poshta/internal/service"
	"poshta/pkg/logger"
	"strings"
)

// Key for user context
type contextKey string

const UserContextKey contextKey = "user"

// JWTMiddleware is middleware for JWT authentication
type JWTMiddleware struct {
	authService service.AuthService
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware(authService service.AuthService) *JWTMiddleware {
	return &JWTMiddleware{
		authService: authService,
	}
}

// Authenticate verifies JWT token and adds user to request context
func (m *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if the header has the Bearer prefix
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]

		// Validate token
		token, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			logger.Error("Invalid token", err, nil)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Get user from token
		user, err := m.authService.GetUserFromToken(token)
		if err != nil {
			logger.Error("Failed to get user from token", err, nil)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateAuthenticatedHandler is a convenience function that wraps a handler with authentication
func (m *JWTMiddleware) CreateAuthenticatedHandler(handler http.HandlerFunc) http.Handler {
	return m.Authenticate(handler)
}