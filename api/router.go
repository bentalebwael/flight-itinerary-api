package api

import (
	"github.com/labstack/echo/v4"

	"flight-itinerary-api/api/middleware"
	"flight-itinerary-api/config"
	"flight-itinerary-api/handlers"
	"flight-itinerary-api/services"
)

// Router handles the setup of all API routes
type Router struct {
	itineraryHandler *handlers.ItineraryHandler
	rateLimiter      *middleware.IPRateLimiter
}

// NewRouter creates a new instance of Router
func NewRouter(cfg *config.AppConfig, itineraryService *services.ItineraryService) *Router {
	// Create rate limiter - 100 requests per minute
	rateLimiter := middleware.NewRateLimiterMiddleware(&cfg.RateLimiter)

	return &Router{
		itineraryHandler: handlers.NewItineraryHandler(itineraryService),
		rateLimiter:      rateLimiter,
	}
}

// SetupRoutes configures all the routes for the application
func (r *Router) SetupRoutes(e *echo.Echo) {
	// API group with rate limiting
	api := e.Group("/api", r.rateLimiter.Middleware())

	// Itinerary routes
	api.POST("/itinerary", r.itineraryHandler.ProcessItinerary)
}
