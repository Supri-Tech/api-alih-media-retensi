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
	"github.com/xuri/excelize/v2"
)

type PasienHandler struct {
	service services.PasienService
}

func NewPasienHandler(service services.PasienService) *PasienHandler {
	return &PasienHandler{service: service}
}

func (hdl *PasienHandler) PasienRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)

		r.Get("/pasien", hdl.GetAll)
		r.Get("/pasien/search", hdl.Search)
		r.Get("/pasien/{id}", hdl.GetByID)
		r.Post("/pasien", hdl.Create)
		r.Put("/pasien/{id}", hdl.Update)
		r.Delete("/pasien/{id}", hdl.Delete)
	})
}

func (hdl *PasienHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	pasien, err := hdl.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	pkg.Success(w, "Data fetched successfully", pasien)
}

func (hdl *PasienHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	pasien, err := hdl.service.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "Pasien not found" {
			pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	pkg.Success(w, "Data found", pasien)
}

func (hdl *PasienHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))

	filter := services.PasienFilter{
		NoRM:       query.Get("NoRM"),
		NamaPasien: query.Get("NamaPasien"),
		NIK:        query.Get("NIK"),
		Limit:      limit,
	}

	if filter.NoRM == "" && filter.NamaPasien == "" && filter.NIK == "" {
		pkg.Error(w, http.StatusBadRequest, "At least one search parameter is required")
		return
	}

	pasien, err := hdl.service.Search(r.Context(), filter)
	if err != nil {
		if err.Error() == "No pasien found" {
			pkg.Error(w, http.StatusNotFound, err.Error())
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	pkg.Success(w, "Data found", pasien)
}

func (hdl *PasienHandler) Create(w http.ResponseWriter, r *http.Request) {
	type CreatePasien struct {
		NoRM         string    `json:"NoRM"`
		NamaPasien   string    `json:"NamaPasien"`
		JenisKelamin string    `json:"JenisKelamin"`
		TglLahir     time.Time `json:"TglLahir"`
		NIK          string    `json:"NIK"`
		Alamat       string    `json:"Alamat"`
		Status       string    `json:"Status"`
	}

	var req CreatePasien
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	pasien := models.Pasien{
		NoRM:         req.NoRM,
		NamaPasien:   req.NamaPasien,
		JenisKelamin: req.JenisKelamin,
		TanggalLahir: req.TglLahir,
		NIK:          req.NIK,
		Alamat:       req.Alamat,
		Status:       req.Status,
		CreatedAt:    time.Now(),
	}

	newPasien, err := hdl.service.Create(r.Context(), pasien)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "Pasien created", newPasien)
}

func (hdl *PasienHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	type UpdatePasien struct {
		NoRM         string    `json:"NoRM"`
		NamaPasien   string    `json:"NamaPasien"`
		JenisKelamin string    `json:"JenisKelamin"`
		TglLahir     time.Time `json:"TglLahir"`
		NIK          string    `json:"NIK"`
		Alamat       string    `json:"Alamat"`
		Status       string    `json:"Status"`
	}

	var req UpdatePasien
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	pasien := models.Pasien{
		ID:           id,
		NoRM:         req.NoRM,
		NamaPasien:   req.NamaPasien,
		JenisKelamin: req.JenisKelamin,
		TanggalLahir: req.TglLahir,
		NIK:          req.NIK,
		Alamat:       req.Alamat,
		Status:       req.Status,
	}

	updatedPasien, err := hdl.service.Update(r.Context(), pasien)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(w, "Pasien updated", updatedPasien)
}

func (hdl *PasienHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
