package app

import (
	"database/sql"
	"net/http"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/handler/v2"
	customMiddleware "github.com/cukiprit/api-sistem-alih-media-retensi/internal/middleware"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	Router *chi.Mux
}

func NewApplication(db *sql.DB) *App {
	kasusRepo := repositories.NewRepoKasus(db)
	kasusService := services.NewServiceKasus(kasusRepo)
	kasusHandler := handler.NewKasusHandler(kasusService)

	userRepo := repositories.NewRepoUser(db)
	userService := services.NewServiceUser(userRepo)
	userHandler := handler.NewUserHandler(userService)

	pasienRepo := repositories.NewRepoPasien(db)
	pasienService := services.NewServicePasien(pasienRepo)
	PasienHandler := handler.NewPasienHandler(pasienService)

	kunjunganRepo := repositories.NewRepoKunjungan(db)
	kunjunganService := services.NewServiceKunjungan(kunjunganRepo)
	kunjunganHandler := handler.NewKunjunganHandler(kunjunganService)

	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.RealIP,
		middleware.Recoverer,
		customMiddleware.SecurityHeaders,
	)

	router.Route("/api/v2", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			pkg.Success(w, "Miaw", nil)
		})
		userHandler.UserRoutes(r)
		kasusHandler.KasusRoutes(r)
		PasienHandler.PasienRoutes(r)
		kunjunganHandler.KunjunganRoutes(r)
	})

	return &App{
		Router: router,
	}
}
