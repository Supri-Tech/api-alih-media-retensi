package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/middleware"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/go-chi/chi/v5"
)

type KunjunganHandler struct {
	service        services.KunjunganService
	dokumenService services.DokumenService
}

func NewKunjunganHandler(service services.KunjunganService, dokumenService services.DokumenService) *KunjunganHandler {
	return &KunjunganHandler{service: service, dokumenService: dokumenService}
}

func (hdl *KunjunganHandler) KunjunganRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)

		r.Get("/kunjungan", hdl.GetAll)
		r.Get("/kunjungan/{id}", hdl.GetByID)
		r.Post("/kunjungan", hdl.Create)
		r.Put("/kunjungan/{id}", hdl.Update)
		r.Delete("/kunjungan/{id}", hdl.Delete)
	})
}

func (hdl *KunjunganHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	kunjungan, err := hdl.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		fmt.Printf("Error: %s", err)
		pkg.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	pkg.Success(w, "Data fetched successfully", kunjungan)
}

func (hdl *KunjunganHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	kunjungan, err := hdl.service.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "Kunjungan not found" {
			pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	pkg.Success(w, "Data found", kunjungan)
}

func (hdl *KunjunganHandler) Create(w http.ResponseWriter, r *http.Request) {
	idPasien, _ := strconv.Atoi(r.FormValue("IdPasien"))
	idKasus, _ := strconv.Atoi(r.FormValue("IdKasus"))
	jenisKunjungan := r.FormValue("JenisKunjungan")

	tanggalMasuk, err := time.Parse("2006-01-02", r.FormValue("TglMasuk"))
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid TglMasuk format, must be YYYY-MM-DD")
		return
	}

	kunjungan := models.Kunjungan{
		IDPasien:       idPasien,
		IDKasus:        idKasus,
		TanggalMasuk:   tanggalMasuk,
		JenisKunjungan: jenisKunjungan,
	}

	newKunjungan, err := hdl.service.Create(r.Context(), kunjungan)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	file, header, err := r.FormFile("File")
	if err == nil {
		defer file.Close()

		_, err := hdl.dokumenService.UploadDokumen(r.Context(), newKunjungan.ID, file, header)
		if err != nil {
			pkg.Error(w, http.StatusInternalServerError, "Kunjungan created but failed to upload file: "+err.Error())
			return
		}
	}

	pkg.Success(w, "Kunjungan created", newKunjungan)
}

// func (hdl *KunjunganHandler) Create(w http.ResponseWriter, r *http.Request) {
// 	type CreateKunjungan struct {
// 		IDPasien       int       `json:"IdPasien"`
// 		IDKasus        int       `json:"IdKasus"`
// 		TanggalMasuk   time.Time `json:"TglMasuk"`
// 		JenisKunjungan string    `json:"JenisKunjungan"`
// 	}
//
// 	var req CreateKunjungan
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}
//
// 	kunjungan := models.Kunjungan{
// 		IDPasien:       req.IDPasien,
// 		IDKasus:        req.IDKasus,
// 		TanggalMasuk:   req.TanggalMasuk,
// 		JenisKunjungan: req.JenisKunjungan,
// 	}
//
// 	newKunjungan, err := hdl.service.Create(r.Context(), kunjungan)
// 	if err != nil {
// 		pkg.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
//
// 	pkg.Success(w, "Pasien created", newKunjungan)
// }

func (hdl *KunjunganHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	idPasien, _ := strconv.Atoi(r.FormValue("IdPasien"))
	idKasus, _ := strconv.Atoi(r.FormValue("IdKasus"))
	jenisKunjungan := r.FormValue("JenisKunjungan")

	tanggalMasuk, err := time.Parse("2006-01-02", r.FormValue("TglMasuk"))
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid TglMasuk format, must be YYYY-MM-DD")
		return
	}

	kunjungan := models.Kunjungan{
		ID:             id,
		IDPasien:       idPasien,
		IDKasus:        idKasus,
		TanggalMasuk:   tanggalMasuk,
		JenisKunjungan: jenisKunjungan,
	}

	updatedKunjungan, err := hdl.service.Update(r.Context(), kunjungan)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	file, header, err := r.FormFile("File")
	if err == nil {
		defer file.Close()

		_, err := hdl.dokumenService.UpdateDokumen(r.Context(), id, file, header)
		if err != nil {
			pkg.Error(w, http.StatusInternalServerError, "Kunjungan updated but failed to upload file: "+err.Error())
			return
		}
	}

	pkg.Success(w, "Kunjungan updated", updatedKunjungan)
}

func (hdl *KunjunganHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := hdl.dokumenService.DeleteDokumen(r.Context(), id); err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Failed to delete dokumen: "+err.Error())
		return
	}

	pkg.Success(w, "Data deleted", nil)
}
