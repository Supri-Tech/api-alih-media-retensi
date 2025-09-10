package handler

import (
	"net/http"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/go-chi/chi/v5"
)

type GeneralHandler struct {
	service services.GeneralService
}

func NewGeneralHandler(service services.GeneralService) *GeneralHandler {
	return &GeneralHandler{service: service}
}

func (h *GeneralHandler) GeneralRoutes(r chi.Router) {
	r.Get("/general/statistik", h.GetStatistik)
}

func (h *GeneralHandler) GetStatistik(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := h.service.GetStatistik(ctx)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(w, "Statistik berhasil diambil", data)
}
