package middleware

import (
	"net/http"

	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
)

func AuthorizeRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value("role").(string)
			if !ok || role == "" {
				pkg.Error(w, http.StatusForbidden, "Role not found in token")
				return
			}

			status, ok := r.Context().Value("status").(string)
			if !ok || status != "aktif" {
				pkg.Error(w, http.StatusForbidden, "User is not active")
				return
			}

			authorized := false
			for _, allowed := range allowedRoles {
				if role == allowed {
					authorized = true
					break
				}
			}

			if !authorized {
				pkg.Error(w, http.StatusForbidden, "You are not allowed to access this resource")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
