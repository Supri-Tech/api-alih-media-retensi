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

type AlihMediaHandler struct {
	service services.AlihMediaService
}

func NewAlihMediaHandler(service services.AlihMediaService) *AlihMediaHandler {
	return &AlihMediaHandler{service: service}
}

func (hdl *AlihMediaHandler) AlihMediaRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)

		r.Get("/alih-media", hdl.GetAll)
		r.Get("/alih-media/search", hdl.Search)
		r.Get("/alih-media/{id}", hdl.GetByID)
		r.Post("/alih-media", hdl.Create)
		r.Put("/alih-media/{id}", hdl.Update)
		r.Delete("/alih-media/{id}", hdl.Delete)
	})
	router.Get("/alih-media/export", hdl.Export)
}

func (hdl *AlihMediaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	alihMedia, err := hdl.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	pkg.Success(w, "Data fetched successfully", alihMedia)
}

func (hdl *AlihMediaHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))

	filter := services.AlihMediaFilter{
		NoRM:       query.Get("NoRM"),
		NamaPasien: query.Get("NamaPasien"),
		Limit:      limit,
	}

	if filter.NoRM == "" && filter.NamaPasien == "" {
		pkg.Error(w, http.StatusBadRequest, "At least one search parameter is required")
		return
	}

	alihMedia, err := hdl.service.Search(r.Context(), filter)
	if err != nil {
		if err.Error() == "No alih media found" {
			pkg.Error(w, http.StatusNotFound, err.Error())
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	pkg.Success(w, "Data found", alihMedia)
}

func (hdl *AlihMediaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	alihMedia, err := hdl.service.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "Alih media not found" {
			pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	pkg.Success(w, "Data found", alihMedia)
}

func (hdl *AlihMediaHandler) Create(w http.ResponseWriter, r *http.Request) {
	type CreateAlihMedia struct {
		ID             int        `json:"IdKunjungan"`
		TanggalLaporan *time.Time `json:"TglLaporan"`
		Status         string     `json:"Status"`
	}

	var req CreateAlihMedia
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	alihMedia := models.AlihMedia{
		ID:         req.ID,
		TglLaporan: req.TanggalLaporan,
		Status:     req.Status,
	}

	newAlihMedia, err := hdl.service.Create(r.Context(), alihMedia)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "Pasien created", newAlihMedia)
}

func (hdl *AlihMediaHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	type UpdateAlihMedia struct {
		ID             int        `json:"IdKunjungan"`
		TanggalLaporan *time.Time `json:"TglLaporan"`
		Status         string     `json:"Status"`
	}

	var req UpdateAlihMedia
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	alihMedia := models.AlihMedia{
		ID:         req.ID,
		TglLaporan: req.TanggalLaporan,
		Status:     req.Status,
	}

	updatedAlihMedia, err := hdl.service.Update(r.Context(), alihMedia)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(w, "Alih media updated", updatedAlihMedia)
}

func (hdl *AlihMediaHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *AlihMediaHandler) Export(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.service.Export(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=alih_media.xlsx")
	w.Write(data)
}
