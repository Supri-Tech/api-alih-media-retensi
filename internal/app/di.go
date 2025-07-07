package app

import (
	"database/sql"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/handler/v1"
	customMiddleware "github.com/cukiprit/api-sistem-alih-media-retensi/internal/middleware"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v1"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v1"
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

	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.RealIP,
		middleware.Recoverer,
		customMiddleware.SecurityHeaders,
	)

	router.Route("/api/v1", func(r chi.Router) {
		userHandler.UserRoutes(r)
		kasusHandler.KasusRoutes(r)
	})

	return &App{
		Router: router,
	}
}
