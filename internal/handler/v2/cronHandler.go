package handler

import (
	"net/http"
	"strconv"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/go-chi/chi/v5"
)

type CronHandler struct {
	cronService services.CronService
}

func NewCronHandler(cronService services.CronService) *CronHandler {
	return &CronHandler{cronService: cronService}
}

func (hdl *CronHandler) CronRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Post("/cron/check-inactive", hdl.CheckInactiveKunjungen)
		r.Post("/cron/process-kunjungan/{id}", hdl.ProcessSingleKunjungan)
		r.Post("/cron/run-now", hdl.RunCronNow)
	})
}

func (hdl *CronHandler) CheckInactiveKunjungen(w http.ResponseWriter, r *http.Request) {
	err := hdl.cronService.CheckAndProcessKunjungan(r.Context())
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Failed to process inactive kunjungen: "+err.Error())
		return
	}

	pkg.Success(w, "Inactive kunjungen check completed successfully", nil)
}

func (hdl *CronHandler) ProcessSingleKunjungan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	err = hdl.cronService.ProcessKunjungan(r.Context(), id)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Failed to process kunjungan: "+err.Error())
		return
	}

	pkg.Success(w, "Kunjungan processed successfully", nil)
}

func (hdl *CronHandler) RunCronNow(w http.ResponseWriter, r *http.Request) {
	err := hdl.cronService.CheckAndProcessKunjungan(r.Context())
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, "Failed to run cron job: "+err.Error())
		return
	}

	pkg.Success(w, "Cron job executed successfully", nil)
}
