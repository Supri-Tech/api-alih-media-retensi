package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
)

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			pkg.Error(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			pkg.Error(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		claims, err := pkg.VerifyToken(token)
		if err != nil {
			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "role", role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
