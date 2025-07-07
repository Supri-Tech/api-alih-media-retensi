package middleware

import (
	"net/http"

	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
)

func CheckActiveUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		email, ok := ctx.Value("email").(string)
		if !ok || email == "" {
			pkg.Error(w, http.StatusUnauthorized, "User information is not found")
			return
		}

	})
}
