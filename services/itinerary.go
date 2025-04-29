package services

import (
	"errors"

	"professional-me-takehome/models"
)

// ItineraryService handles the business logic for processing flight tickets
type ItineraryService struct{}

// NewItineraryService creates a new instance of ItineraryService
func NewItineraryService() *ItineraryService {
	return &ItineraryService{}
}

// ReconstructItinerary processes the flight tickets and returns an ordered itinerary
func (s *ItineraryService) ReconstructItinerary(request *models.ItineraryRequest) ([]string, error) {
	// Build graph representation of flights
	graph := make(map[string]string)
	inDegree := make(map[string]int)

	// Populate the graph and count incoming edges
	for _, ticket := range request.Tickets {
		src, dst := ticket[0], ticket[1]

		// Check for duplicate flights from same source
		if _, exists := graph[src]; exists {
			return nil, errors.New("invalid tickets: multiple flights from same source")
		}

		graph[src] = dst
		inDegree[dst]++
		// Ensure source is in inDegree map
		if _, exists := inDegree[src]; !exists {
			inDegree[src] = 0
		}
	}

	// Find starting airport (node with no incoming edges)
	var start string
	for airport, degree := range inDegree {
		if degree == 0 {
			if start != "" {
				return nil, errors.New("invalid tickets: multiple starting points found")
			}
			start = airport
		}
	}

	if start == "" {
		return nil, errors.New("invalid tickets: no starting point found")
	}

	// Construct itinerary by following the graph
	itinerary := make([]string, 0, len(request.Tickets)+1)
	current := start

	for len(itinerary) <= len(request.Tickets) {
		itinerary = append(itinerary, current)

		next, exists := graph[current]
		if !exists {
			break // We've reached the final destination
		}
		current = next
	}

	// Verify we used all tickets
	if len(itinerary)-1 != len(request.Tickets) {
		return nil, errors.New("invalid tickets: disconnected route")
	}

	return itinerary, nil
}
