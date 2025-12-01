package middleware

import (
	"context"
	"net/http"
	"strings"

	"jira-clone-be/internal/service"
	"jira-clone-be/pkg/logger"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate token
			userID, err := authService.ValidateToken(token)
			if err != nil {
				logger := logger.New("info")
				logger.Error("Token validation failed", map[string]interface{}{
					"error": err.Error(),
					"token": token[:10] + "...", // Log only first 10 chars for security
				})
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
