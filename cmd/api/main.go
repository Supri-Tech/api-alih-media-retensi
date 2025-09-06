package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/app"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/database"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := database.InitDB()
	defer db.Close()

	app := app.NewApplication(db)

	startCronScheduler(app.CronService)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server started on :%s", port)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      app.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	waitForShutdown(server)
}

func startCronScheduler(cronService services.CronService) {
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
	)
	if err != nil {
		log.Printf("Failed to create scheduler: %v", err)
		return
	}

	cronSchedule := os.Getenv("CRON_SCHEDULE")
	if cronSchedule == "" {
		cronSchedule = "08:00"
	}

	log.Printf("Setting up cron job to run daily at %s", cronSchedule)

	_, err = scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(8, 0, 0),
			),
		),
		gocron.NewTask(
			func() {
				log.Println("Starting scheduled cron job...")
				ctx := context.Background()
				startTime := time.Now()

				if err := cronService.CheckAndProcessKunjungan(ctx); err != nil {
					log.Printf("Cron job failed: %v", err)
				} else {
					duration := time.Since(startTime)
					log.Printf("Cron job completed successfully in %v", duration)
				}
			},
		),
	)

	if err != nil {
		log.Printf("Failed to schedule cron job: %v", err)
		return
	}

	runInitialCheck := os.Getenv("RUN_INITIAL_CRON")
	if runInitialCheck == "true" || runInitialCheck == "1" {
		log.Println("🔍 Running initial cron job check...")
		ctx := context.Background()
		startTime := time.Now()

		if err := cronService.CheckAndProcessKunjungan(ctx); err != nil {
			log.Printf("Initial cron job failed: %v", err)
		} else {
			duration := time.Since(startTime)
			log.Printf("Initial cron job completed successfully in %v", duration)
		}
	}

	scheduler.Start()
	log.Printf("Cron job scheduler started")

	jobs := scheduler.Jobs()
	if len(jobs) > 0 {
		nextRun, err := jobs[0].NextRun()
		if err == nil {
			log.Printf("Next run at: %v", nextRun.Format("2006-01-02 15:04:05"))
		} else {
			log.Printf("Failed to get next run time: %v", err)
		}
	}

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			// Get jobs - only one return value in v2
			jobs := scheduler.Jobs()
			if len(jobs) > 0 {
				nextRun, err := jobs[0].NextRun()
				if err == nil {
					log.Printf("Cron service heartbeat - Next run: %v", nextRun.Format("2006-01-02 15:04:05"))
				} else {
					log.Printf("Cron service heartbeat - Failed to get next run: %v", err)
				}
			} else {
				log.Printf("Cron service heartbeat - No jobs scheduled")
			}
		}
	}()
}

func waitForShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received. Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
