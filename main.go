package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"flight-itinerary-api/handlers"
	"flight-itinerary-api/services"
)

func main() {
	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize services and handlers
	itineraryService := services.NewItineraryService()
	itineraryHandler := handlers.NewItineraryHandler(itineraryService)

	// Routes
	api := e.Group("/api")
	api.POST("/itinerary", itineraryHandler.ProcessItinerary)

	// Start server
	log.Fatal(e.Start(":8080"))
}
