package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"flight-itinerary-api/models"
	"flight-itinerary-api/services"
)

// ItineraryHandler handles HTTP requests for flight itinerary operations
type ItineraryHandler struct {
	service *services.ItineraryService
}

// NewItineraryHandler creates a new instance of ItineraryHandler
func NewItineraryHandler(service *services.ItineraryService) *ItineraryHandler {
	return &ItineraryHandler{
		service: service,
	}
}

// ProcessItinerary handles the POST request to process flight tickets and return an ordered itinerary
func (h *ItineraryHandler) ProcessItinerary(c echo.Context) error {
	var request models.ItineraryRequest

	// Parse request body
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Process the itinerary with context
	itinerary, err := h.service.ReconstructItinerary(c.Request().Context(), &request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Return the response
	response := models.ItineraryResponse{
		Itinerary: itinerary,
	}

	return c.JSON(http.StatusOK, response)
}
