package main

import (
	// "fmt"
	"log"
	"net/http"
	"os"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/app"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/database"
	"github.com/joho/godotenv"
	// "github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.InitDB()

	// router := chi.NewRouter()
	app := app.NewApplication(db)

	// router.Use(middleware.Logger)
	// router.Use(middleware.Recoverer)
	// router.Use(middleware.RealIP)

	// router.Route("/api/v1", func(r chi.Router) {
	// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 		pkg.Success(w, "Miaw", nil)
	// 	})
	// })

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
	// if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
	// 	log.Fatal("Failed to start server: ", err)
	// }
}
