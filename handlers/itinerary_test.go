package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"flight-itinerary-api/config"
	"flight-itinerary-api/models"
	"flight-itinerary-api/services"
)

func TestProcessItinerary(t *testing.T) {
	// Setup
	e := echo.New()
	cfg := &config.AppConfig{
		WorkerPool: config.WorkerPoolConfig{
			WorkerCount: 5,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service := services.NewItineraryService(ctx, cfg)
	handler := NewItineraryHandler(service)

	tests := []struct {
		name       string
		input      models.ItineraryRequest
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid request",
			input: models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"LAX", "JFK"},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid request - empty tickets",
			input: models.ItineraryRequest{
				Tickets: []models.TicketPair{},
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid request - disconnected route",
			input: models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"JFK", "MCO"},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/itinerary", strings.NewReader(string(body))).WithContext(context.Background())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// Create response recorder
			rec := httptest.NewRecorder()

			// Create echo context
			c := e.NewContext(req, rec)

			// Serve request
			err = handler.ProcessItinerary(c)
			if err != nil {
				// Echo's error handler would normally handle this
				he, ok := err.(*echo.HTTPError)
				if ok {
					rec.Code = he.Code
				} else {
					rec.Code = http.StatusInternalServerError
				}
			}

			// Check status code
			if rec.Code != tt.wantStatus {
				t.Errorf("ProcessItinerary() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			// Parse response
			var response struct {
				Itinerary []string `json:"itinerary,omitempty"`
				Error     string   `json:"error,omitempty"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			// Verify response structure
			if tt.wantErr {
				if response.Error == "" {
					t.Error("Expected error in response, got none")
				}
				if response.Itinerary != nil {
					t.Error("Expected no itinerary in error response")
				}
			} else {
				if response.Error != "" {
					t.Errorf("Unexpected error in response: %v", response.Error)
				}
				if response.Itinerary == nil {
					t.Error("Expected itinerary in successful response")
				}
			}
		})
	}
}
