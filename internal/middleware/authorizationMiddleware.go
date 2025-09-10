package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
)

// func VerifyToken(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			pkg.Error(w, http.StatusUnauthorized, "Missing Authorization header")
// 			return
// 		}

// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid Authorization header format")
// 			return
// 		}

// 		claims, err := pkg.VerifyToken(parts[1])
// 		if err != nil {
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		uidFloat, ok := claims["user_id"].(float64)
// 		if !ok {
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid token payload")
// 			return
// 		}
// 		userID := int(uidFloat)

// 		ctx := context.WithValue(r.Context(), "userID", userID)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			pkg.Error(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			pkg.Error(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		claims, err := pkg.VerifyToken(parts[1])
		if err != nil {
			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		uidFloat, ok := claims["user_id"].(float64)
		role, okRole := claims["role"].(string)
		if !ok || !okRole {
			pkg.Error(w, http.StatusUnauthorized, "Invalid token payload")
			return
		}
		userID := int(uidFloat)

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		ctx = context.WithValue(ctx, "userRole", role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func VerifyAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := pkg.GetUserRoleFromCtx(r.Context())
		if role != "admin" {
			pkg.Error(w, http.StatusForbidden, "Access denied")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// func VerifyToken(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			pkg.Error(w, http.StatusUnauthorized, "Authorization header required")
// 			return
// 		}

// 		token := strings.TrimPrefix(authHeader, "Bearer ")
// 		if token == authHeader {
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid token format")
// 			return
// 		}

// 		claims, err := pkg.VerifyToken(token)
// 		if err != nil {
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		var userID string
// 		switch v := claims["user_id"].(type) {
// 		case float64:
// 			userID = strconv.Itoa(int(v))
// 		case int:
// 			userID = strconv.Itoa(v)
// 		case string:
// 			userID = v
// 		default:
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		if userID == "" {
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		email, ok := claims["email"].(string)
// 		if !ok || email == "" {
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		status, ok := claims["status"].(string)
// 		if !ok || status == "" {
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		role, ok := claims["role"].(string)
// 		if !ok || role == "" {
// 			pkg.Error(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		ctx := r.Context()
// 		ctx = context.WithValue(ctx, "user_id", userID)
// 		ctx = context.WithValue(ctx, "email", email)
// 		ctx = context.WithValue(ctx, "role", role)
// 		ctx = context.WithValue(ctx, "status", status)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
