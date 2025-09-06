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
	"golang.org/x/time/rate"
)

type App struct {
	Router      *chi.Mux
	CronService services.CronService
}

func NewApplication(db *sql.DB) *App {
	kasusRepo := repositories.NewRepoKasus(db)
	dokumenRepo := repositories.NewRepoDokumen(db)
	userRepo := repositories.NewRepoUser(db)
	pasienRepo := repositories.NewRepoPasien(db)
	kunjunganRepo := repositories.NewRepoKunjungan(db)
	infoSistemRepo := repositories.NewRepoInfoSistem(db)
	aliMediaRepo := repositories.NewRepoAlihMedia(db)
	retensiRepo := repositories.NewRepoRetensi(db)
	pemusnahanRepo := repositories.NewRepoPemusnahan(db)

	kasusService := services.NewServiceKasus(kasusRepo)
	userService := services.NewServiceUser(userRepo)
	pasienService := services.NewServicePasien(pasienRepo)
	kunjunganService := services.NewServiceKunjungan(kunjunganRepo, pasienRepo, kasusRepo)
	dokumenService := services.NewServiceDokumen(dokumenRepo)
	infoSistemService := services.InfoSistemService(infoSistemRepo)
	alihMediaService := services.NewServiceAlihMedia(aliMediaRepo, kunjunganRepo, kasusRepo)
	retensiService := services.NewServiceRetensi(retensiRepo)
	pemusnahanService := services.NewServicePemusnahan(pemusnahanRepo)

	kasusHandler := handler.NewKasusHandler(kasusService)
	userHandler := handler.NewUserHandler(userService)
	PasienHandler := handler.NewPasienHandler(pasienService)
	kunjunganHandler := handler.NewKunjunganHandler(kunjunganService, dokumenService, alihMediaService)
	infoSistemHandler := handler.NewInfoSistemHandler(infoSistemService)
	alihMediaHandler := handler.NewAlihMediaHandler(alihMediaService)
	retensiHandler := handler.NewRetensiHandler(retensiService)
	pemusnahanHandler := handler.NewPemusnahanHandler(pemusnahanService)

	cronService := services.NewCronService(kunjunganRepo, kasusRepo, aliMediaRepo)
	cronHandler := handler.NewCronHandler(cronService)

	router := chi.NewRouter()

	limiter := rate.NewLimiter(rate.Limit(1), 5)

	router.Use(
		middleware.Logger,
		middleware.RealIP,
		middleware.Recoverer,
		customMiddleware.SecurityHeaders,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !limiter.Allow() {
					pkg.Error(w, http.StatusTooManyRequests, "Too many request")
					return
				}
				next.ServeHTTP(w, r)
			})
		},
	)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
	}))

	router.Handle("/uploads/*", http.StripPrefix("/uploads", http.FileServer(http.Dir("uploads"))))

	router.Route("/api/v2", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			pkg.Success(w, "Miaw", nil)
		})
		userHandler.UserRoutes(r)
		kasusHandler.KasusRoutes(r)
		PasienHandler.PasienRoutes(r)
		kunjunganHandler.KunjunganRoutes(r)
		infoSistemHandler.InfoSistemRoutes(r)
		alihMediaHandler.AlihMediaRoutes(r)
		retensiHandler.RetensiRoutes(r)
		pemusnahanHandler.PemusnahanRoutes(r)
		cronHandler.CronRoutes(r)
	})

	return &App{
		Router:      router,
		CronService: cronService,
	}
}
