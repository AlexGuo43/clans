package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/AlexGuo43/clans/api-gateway/internal/services"
)

func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicEndpoint(r.URL.Path, r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := authService.ValidateJWT(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			r.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublicEndpoint(path, method string) bool {
	publicEndpoints := map[string][]string{
		"/api/auth/signup":     {"POST"},
		"/api/auth/login":      {"POST"},
		"/api/posts":           {"GET"},
		"/api/clans":           {"GET"},
		"/health":              {"GET"},
	}

	if strings.HasPrefix(path, "/api/posts/") && method == "GET" {
		return true
	}

	if strings.HasPrefix(path, "/api/comments/") && method == "GET" {
		return true
	}

	if strings.HasPrefix(path, "/api/clans/") && method == "GET" {
		return true
	}

	allowedMethods, exists := publicEndpoints[path]
	if !exists {
		return false
	}

	for _, allowedMethod := range allowedMethods {
		if method == allowedMethod {
			return true
		}
	}

	return false
}