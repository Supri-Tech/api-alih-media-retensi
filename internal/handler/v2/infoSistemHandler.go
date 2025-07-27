package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/middleware"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/go-chi/chi/v5"
)

type InfoSistemHandler struct {
	service services.InfoSistemService
}

func NewInfoSistemHandler(service services.InfoSistemService) *InfoSistemHandler {
	return &InfoSistemHandler{service: service}
}

func (hdl *InfoSistemHandler) InfoSistemRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)

		r.Get("/info-sistem", hdl.GetAll)
		r.Get("/info-sistem/{id}", hdl.GetByID)
		r.Post("/info-sistem", hdl.Create)
		r.Put("/info-sistem/{id}", hdl.Update)
	})
}

func (hdl *InfoSistemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var infoSistem models.InfoSistem
	info, err := hdl.service.GetAllInfoSistem(r.Context(), infoSistem)
	log.Print(err)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	pkg.Success(w, "Data fetched successfully", info)
}

func (hdl *InfoSistemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	info, err := hdl.service.GetInfoSistem(r.Context(), id)
	if err != nil {
		if err.Error() == "Kasus not found" {
			pkg.Error(w, http.StatusNotFound, err.Error())
		} else {
			pkg.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	pkg.Success(w, "Data found", info)
}

func (hdl *InfoSistemHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(int64(pkg.MaxFileSize)); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	infoSistem := models.InfoSistem{
		NamaAplikasi: r.FormValue("NamaAplikasi"),
	}

	logoBase64, err := pkg.ParseImage(r, "Logo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	infoSistem.Logo = logoBase64

	info, err := hdl.service.CreateInfoSistem(r.Context(), infoSistem)
	log.Print(err)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	pkg.Success(w, "Data created", info)
}

func (hdl *InfoSistemHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := r.ParseMultipartForm(int64(pkg.MaxFileSize)); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	infoSistem := models.InfoSistem{
		ID:           id,
		NamaAplikasi: r.FormValue("NamaAplikasi"),
	}

	logoBase64, err := pkg.ParseImage(r, "Logo")
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if logoBase64 != "" {
		infoSistem.Logo = logoBase64
	}

	info, err := hdl.service.UpdateInfoSistem(r.Context(), infoSistem)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	pkg.Success(w, "Data updated", info)
}
