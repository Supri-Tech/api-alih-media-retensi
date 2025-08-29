package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

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
		r.Get("/kasus/search", hdl.Search)
		r.Get("/kasus/{id}", hdl.GetByID)
		r.Post("/kasus", hdl.Create)
		r.Put("/kasus/{id}", hdl.Update)
		r.Delete("/kasus/{id}", hdl.Delete)
		r.Post("/kasus/import", hdl.Import)
	})
	router.Get("/kasus/export", hdl.Export)
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

func (hdl *KasusHandler) Import(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	// Get the file from form data
	file, header, err := r.FormFile("File")
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Failed to get file from form data. Use key 'file'")
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		pkg.Error(w, http.StatusBadRequest, "Only Excel files (.xlsx) are allowed")
		return
	}

	tempFile, err := os.CreateTemp("", "kasus-upload-*.xlsx")
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Failed to create temporary file")
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	err = hdl.service.Import(r.Context(), tempFile.Name())
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Failed to import data: "+err.Error())
		return
	}

	pkg.Success(w, "Excel file imported successfully", nil)
}

func (hdl *KasusHandler) Export(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filter := services.KasusFilter{
		JenisKasus: query.Get("JenisKasus"),
	}

	excelData, err := hdl.service.Export(r.Context(), filter)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Failed to export data: "+err.Error())
		return
	}

	filename := "data_kasus_" + time.Now().Format("20060102_150405") + ".xlsx"
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Length", strconv.Itoa(len(excelData)))

	if _, err := w.Write(excelData); err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Failed to write Excel file")
		return
	}
}
