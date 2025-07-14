package handler

import (
	"net/http"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/middleware"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/go-chi/chi/v5"
)

type KunjunganHandler struct {
	service services.KunjunganService
}

func NewKunjunganHandler(service services.KunjunganService) *KunjunganHandler {
	return &KunjunganHandler{service: service}
}

func (hdl *KunjunganHandler) KunjunganRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)

		r.Get("/kunjungan", func(w http.ResponseWriter, r *http.Request) {})
	})
}
