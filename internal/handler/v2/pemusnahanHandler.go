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

type PemusnahanHandler struct {
	service services.PemusnahanService
}

func NewPemusnahanHandler(service services.PemusnahanService) *PemusnahanHandler {
	return &PemusnahanHandler{service: service}
}

func (hdl *PemusnahanHandler) PemusnahanRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)

		r.Get("/pemusnahan", hdl.GetAll)
		r.Get("/pemusnahan/search", hdl.Search)
		r.Get("/pemusnahan/{id}", hdl.GetByID)
		r.Post("/pemusnahan", hdl.Create)
		r.Put("/pemusnahan/{id}", hdl.Update)
		r.Delete("/pemusnahan/{id}", hdl.Delete)
	})
	router.Get("/pemusnahan/export", hdl.Export)
}

func (hdl *PemusnahanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	pemusnahan, err := hdl.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	pkg.Success(w, "Data fetched successfully", pemusnahan)
}

func (hdl *PemusnahanHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))

	filter := services.PemusnahanFilter{
		NoRM:       query.Get("NoRM"),
		NamaPasien: query.Get("NamaPasien"),
		Limit:      limit,
	}

	if filter.NoRM == "" && filter.NamaPasien == "" {
		pkg.Error(w, http.StatusBadRequest, "At least one search parameter is required")
		return
	}

	pemusnahan, err := hdl.service.Search(r.Context(), filter)
	if err != nil {
		if err.Error() == "No pemusnahan found" {
			pkg.Error(w, http.StatusNotFound, err.Error())
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	pkg.Success(w, "Data found", pemusnahan)
}

func (hdl *PemusnahanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	pemusnahan, err := hdl.service.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "Kunjungan not found" {
			pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	pkg.Success(w, "Data found", pemusnahan)
}

func (hdl *PemusnahanHandler) Create(w http.ResponseWriter, r *http.Request) {
	type CreatePemusnahan struct {
		ID             int        `json:"IdKunjungan"`
		TanggalLaporan *time.Time `json:"TglLaporan"`
		Status         string     `json:"Status"`
	}

	var req CreatePemusnahan
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	pemusnahan := models.Pemusnahan{
		ID:         req.ID,
		TglLaporan: req.TanggalLaporan,
		Status:     req.Status,
	}

	newPemusnahan, err := hdl.service.Create(r.Context(), pemusnahan)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "Pemusnahan created", newPemusnahan)
}

func (hdl *PemusnahanHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	type UpdatePemusnahan struct {
		ID             int        `json:"IdKunjungan"`
		TanggalLaporan *time.Time `json:"TglLaporan"`
		Status         string     `json:"Status"`
	}

	var req UpdatePemusnahan
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// var tglLaporan *time.Time
	// if req.TanggalLaporan != "" {
	// 	parsedTime, err := time.Parse("2006-01-02", req.TanggalLaporan)
	// 	if err != nil {
	// 		pkg.Error(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
	// 		return
	// 	}
	// 	tglLaporan = &parsedTime
	// }

	pemusnahan := models.Pemusnahan{
		ID:         id,
		TglLaporan: req.TanggalLaporan,
		Status:     req.Status,
	}

	updatedPemusnahan, err := hdl.service.Update(r.Context(), pemusnahan)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(w, "Pemusnahan updated", updatedPemusnahan)
}

func (hdl *PemusnahanHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *PemusnahanHandler) Export(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.service.Export(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=pemusnahan.xlsx")
	w.Write(data)
}
