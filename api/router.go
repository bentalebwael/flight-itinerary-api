package api

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"flight-itinerary-api/api/middleware"
	"flight-itinerary-api/config"
	"flight-itinerary-api/handlers"
	"flight-itinerary-api/services"
)

// Router handles the setup of all API routes
type Router struct {
	logger           *zap.Logger
	rateLimiter      *middleware.IPRateLimiter
	itineraryHandler *handlers.ItineraryHandler
}

// NewRouter creates a new instance of Router
func NewRouter(cfg *config.AppConfig, logger *zap.Logger, itineraryService *services.ItineraryService) *Router {
	// Create rate limiter
	rateLimiter := middleware.NewRateLimiterMiddleware(&cfg.RateLimiter)

	return &Router{
		logger:           logger,
		rateLimiter:      rateLimiter,
		itineraryHandler: handlers.NewItineraryHandler(itineraryService),
	}
}

// SetupRoutes configures all the routes for the application
func (r *Router) SetupRoutes(e *echo.Echo) {
	// Middleware
	e.Use(middleware.Logger(r.logger))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// API group with rate limiting
	api := e.Group("/api", r.rateLimiter.Middleware())

	// Itinerary routes
	api.POST("/itinerary", r.itineraryHandler.ProcessItinerary)
}
