package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/middleware"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/go-chi/chi/v5"
)

type RetensiHandler struct {
	service services.RetensiService
}

func NewRetensiHandler(service services.RetensiService) *RetensiHandler {
	return &RetensiHandler{service: service}
}

func (hdl *RetensiHandler) RetensiRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)

		r.Get("/retensi", hdl.GetAll)
		r.Get("/retensi/{id}", hdl.GetByID)
		r.Post("/retensi", hdl.Create)
		r.Put("/retensi/{id}", hdl.Update)
		r.Delete("/retensi/{id}", hdl.Delete)
	})
}

func (hdl *RetensiHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	retensi, err := hdl.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	pkg.Success(w, "Data fetched successfully", retensi)
}

func (hdl *RetensiHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	retensi, err := hdl.service.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "Kunjungan not found" {
			pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	pkg.Success(w, "Data found", retensi)
}

func (hdl *RetensiHandler) Create(w http.ResponseWriter, r *http.Request) {
	type CreateRetensi struct {
		ID             int        `json:"IdKunjugan"`
		TanggalLaporan *time.Time `json:"TglLaporan"`
		Status         string     `json:"Status"`
	}

	var req CreateRetensi
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	retensi := models.Retensi{
		ID:         req.ID,
		TglLaporan: req.TanggalLaporan,
		Status:     req.Status,
	}

	newRetensi, err := hdl.service.Create(r.Context(), retensi)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "Retensi created", newRetensi)
}

func (hdl *RetensiHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	type UpdateRetensi struct {
		TanggalLaporan string `json:"TglLaporan"`
		Status         string `json:"Status"`
	}

	var req UpdateRetensi
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var tglLaporan *time.Time
	if req.TanggalLaporan != "" {
		parsedTime, err := time.Parse("2006-01-02", req.TanggalLaporan)
		if err != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
			return
		}
		tglLaporan = &parsedTime
	}

	retensi := models.Retensi{
		ID:         id,
		TglLaporan: tglLaporan,
		Status:     req.Status,
	}

	updatedRetensi, err := hdl.service.Update(r.Context(), retensi)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(w, "Retensi updated", updatedRetensi)
}

func (hdl *RetensiHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := hdl.service.Delete(r.Context(), id); err != nil {
		pkg.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(w, "Data deleted", nil)
}
