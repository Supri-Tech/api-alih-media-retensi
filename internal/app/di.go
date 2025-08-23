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
	"github.com/go-chi/cors"
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

	dokumenRepo := repositories.NewRepoDokumen(db)
	dokumenService := services.NewServiceDokumen(dokumenRepo)

	kunjunganRepo := repositories.NewRepoKunjungan(db)
	kunjunganService := services.NewServiceKunjungan(kunjunganRepo)
	kunjunganHandler := handler.NewKunjunganHandler(kunjunganService, dokumenService)

	infoSistemRepo := repositories.NewRepoInfoSistem(db)
	infoSistemService := services.InfoSistemService(infoSistemRepo)
	InfoSistemHandler := handler.NewInfoSistemHandler(infoSistemService)

	aliMediaRepo := repositories.NewRepoAlihMedia(db)
	alihMediaService := services.NewServiceAlihMedia(aliMediaRepo)
	alihMediaHandler := handler.NewAlihMediaHandler(alihMediaService)

	retensiRepo := repositories.NewRepoRetensi(db)
	retensiService := services.NewServiceRetensi(retensiRepo)
	retensiHandler := handler.NewRetensiHandler(retensiService)

	pemusnahanRepo := repositories.NewRepoPemusnahan(db)
	pemusnahanService := services.NewServicePemusnahan(pemusnahanRepo)
	pemusnahanHandler := handler.NewPemusnahanHandler(pemusnahanService)

	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.RealIP,
		middleware.Recoverer,
		customMiddleware.SecurityHeaders,
	)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
	}))

	router.Route("/api/v2", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			pkg.Success(w, "Miaw", nil)
		})
		userHandler.UserRoutes(r)
		kasusHandler.KasusRoutes(r)
		PasienHandler.PasienRoutes(r)
		kunjunganHandler.KunjunganRoutes(r)
		InfoSistemHandler.InfoSistemRoutes(r)
		alihMediaHandler.AlihMediaRoutes(r)
		retensiHandler.RetensiRoutes(r)
		pemusnahanHandler.PemusnahanRoutes(r)
	})

	return &App{
		Router: router,
	}
}
