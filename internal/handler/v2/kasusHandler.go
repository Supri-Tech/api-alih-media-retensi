package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/middleware"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/go-chi/chi/v5"
)

type KasusHandler struct {
	service services.KasusService
}

func NewKasusHandler(service services.KasusService) *KasusHandler {
	return &KasusHandler{service: service}
}

func (hdl *KasusHandler) KasusRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)

		r.Get("/kasus", hdl.GetAll)
		r.Get("/kasus/{id}", hdl.GetByID)
		r.Get("/kasus/search", hdl.Search)
		r.Post("/kasus", hdl.Create)
		r.Put("/kasus/{id}", hdl.Update)
		r.Delete("/kasus/{id}", hdl.Delete)
	})
}

func (hdl *KasusHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	kasus, err := hdl.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	pkg.Success(w, "Data fetched successfully", kasus)
}

func (hdl *KasusHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	kasus, err := hdl.service.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "Kasus not found" {
			pkg.Error(w, http.StatusNotFound, err.Error())
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	pkg.Success(w, "Data found", kasus)
}

func (hdl *KasusHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))

	filter := services.KasusFilter{
		JenisKasus: query.Get("JenisKasus"),
		Limit:      limit,
	}

	if filter.JenisKasus == "" {
		pkg.Error(w, http.StatusBadRequest, "At least one search parameter is required")
		return
	}

	kasus, err := hdl.service.Search(r.Context(), filter)
	if err != nil {
		if err.Error() == "No kasus found" {
			pkg.Error(w, http.StatusNotFound, err.Error())
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	pkg.Success(w, "Data found", kasus)
}

func (hdl *KasusHandler) Create(w http.ResponseWriter, r *http.Request) {
	type CreateKasus struct {
		JenisKasus    string `json:"JenisKasus"`
		MasaAktifRi   int    `json:"MasaAktifRi"`
		MasaInaktifRi int    `json:"MasaInaktifRi"`
		MasaAktifRj   int    `json:"MasaAktifRj"`
		MasaInaktifRj int    `json:"MasaInaktifRj"`
		InfoLain      string `json:"InfoLain"`
	}

	var req CreateKasus
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	kasus := models.Kasus{
		JenisKasus:    req.JenisKasus,
		MasaAktifRI:   req.MasaAktifRi,
		MasaInaktifRI: req.MasaInaktifRi,
		MasaAktifRJ:   req.MasaAktifRj,
		MasaInaktifRJ: req.MasaInaktifRj,
		InfoLain:      req.InfoLain,
	}

	newKasus, err := hdl.service.Create(r.Context(), kasus)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "Kasus created", newKasus)
}

func (hdl *KasusHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	type UpdateKasus struct {
		JenisKasus    string `json:"JenisKasus"`
		MasaAktifRi   int    `json:"MasaAktifRi"`
		MasaInaktifRi int    `json:"MasaInaktifRi"`
		MasaAktifRj   int    `json:"MasaAktifRj"`
		MasaInaktifRj int    `json:"MasaInaktifRj"`
		InfoLain      string `json:"InfoLain"`
	}

	var req UpdateKasus
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	kasus := models.Kasus{
		ID:            id,
		JenisKasus:    req.JenisKasus,
		MasaAktifRI:   req.MasaAktifRi,
		MasaInaktifRI: req.MasaInaktifRi,
		MasaAktifRJ:   req.MasaAktifRj,
		MasaInaktifRJ: req.MasaInaktifRj,
		InfoLain:      req.InfoLain,
	}

	updatedKasus, err := hdl.service.Update(r.Context(), kasus)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(w, "Data updated", updatedKasus)
}

func (hdl *KasusHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
