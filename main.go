package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"flight-itinerary-api/api"
	"flight-itinerary-api/config"
	"flight-itinerary-api/services"
)

func main() {
	// Initialize logger
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to intialize the logger: %v", err)
	}
	defer l.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Echo
	e := echo.New()

	// Initialize services with configured worker count
	itineraryService := services.NewItineraryService(cfg.WorkerPool.WorkerCount)

	// Setup router
	router := api.NewRouter(cfg, l, itineraryService)
	router.SetupRoutes(e)

	// Start server
	err = e.Start(fmt.Sprintf(":%s", cfg.Server.Port))
	l.Fatal("Error while running the server", zap.Error(err))

	// Graceful shutdown
	defer func() {
		// Stop worker pool
		itineraryService.Stop()

		if err := e.Shutdown(nil); err != nil {
			l.Fatal("Failed to shutdown server", zap.Error(err))
		}
	}()
}
