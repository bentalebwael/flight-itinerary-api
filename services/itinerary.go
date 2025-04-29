package services

import (
	"context"
	"errors"
	"sync"

	"flight-itinerary-api/models"

	"github.com/panjf2000/ants/v2"
)

// ItineraryService handles the business logic for processing flight tickets
type ItineraryService struct {
	pool *ants.Pool
}

// NewItineraryService creates a new instance of ItineraryService
func NewItineraryService(workerCount int) *ItineraryService {
	if workerCount <= 0 {
		workerCount = 1000 // default worker count
	}

	pool, _ := ants.NewPool(workerCount) // Ignoring error as we validate workerCount > 0
	return &ItineraryService{
		pool: pool,
	}
}

// Stop gracefully shuts down all workers
func (s *ItineraryService) Stop() {
	s.pool.Release()
}

// ReconstructItinerary processes the flight tickets and returns an ordered itinerary
func (s *ItineraryService) ReconstructItinerary(ctx context.Context, request *models.ItineraryRequest) ([]string, error) {
	var (
		result []string
		err    error
		wg     sync.WaitGroup
	)

	wg.Add(1)
	submitErr := s.pool.Submit(func() {
		defer wg.Done()
		result, err = processItinerary(request)
	})

	if submitErr != nil {
		return nil, submitErr
	}

	// Wait for completion or context cancellation
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		return result, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// processItinerary handles the actual itinerary reconstruction logic
func processItinerary(request *models.ItineraryRequest) ([]string, error) {
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
