package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"flight-itinerary-api/api"
	"flight-itinerary-api/config"
	"flight-itinerary-api/services"
)

func main() {
	// Initialize logger

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize services
	itineraryService := services.NewItineraryService()

	// Setup router
	router := api.NewRouter(cfg, itineraryService)
	router.SetupRoutes(e)

	// Start server
	log.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Server.Port)))

	// Graceful shutdown
	defer func() {
		if err := e.Shutdown(nil); err != nil {
			log.Fatal("Failed to shutdown server", zap.Error(err))
		}
	}()
}
