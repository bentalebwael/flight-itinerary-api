package services

import (
	"context"
	"reflect"
	"testing"

	"flight-itinerary-api/config"
	"flight-itinerary-api/models"
)

func TestReconstructItinerary(t *testing.T) {
	cfg := &config.AppConfig{
		WorkerPool: config.WorkerPoolConfig{
			WorkerCount: 5,
		},
	}

	service := NewItineraryService(cfg)
	defer service.Stop()

	tests := []struct {
		name    string
		request *models.ItineraryRequest
		want    []string
		wantErr bool
	}{
		{
			name: "valid linear itinerary",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"LAX", "JFK"},
				},
			},
			want:    []string{"SFO", "LAX", "JFK"},
			wantErr: false,
		},
		{
			name: "multiple flights from same source",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"SFO", "JFK"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "multiple starting points",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"JFK", "MCO"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no starting point (cyclic route)",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"LAX", "JFK"},
					{"JFK", "SFO"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "disconnected route",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"MCO", "JFK"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "complex valid itinerary",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "ATL"},
					{"ATL", "MCO"},
					{"MCO", "JFK"},
				},
			},
			want:    []string{"SFO", "ATL", "MCO", "JFK"},
			wantErr: false,
		},
		{
			name: "empty tickets array",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "single ticket",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "JFK"},
				},
			},
			want:    []string{"SFO", "JFK"},
			wantErr: false,
		},
		{
			name: "large complex itinerary",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"LAX", "DFW"},
					{"DFW", "ORD"},
					{"ORD", "JFK"},
					{"JFK", "LHR"},
					{"LHR", "CDG"},
				},
			},
			want:    []string{"SFO", "LAX", "DFW", "ORD", "JFK", "LHR", "CDG"},
			wantErr: false,
		},
		{
			name: "repeated destinations with different sources",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"LAX", "DFW"},
					{"DFW", "LAX"}, // Return to LAX
					{"LAX", "JFK"},
				},
			},
			want:    nil,
			wantErr: true, // Should fail because LAX is used twice as destination
		},
		{
			name: "shuffle order but valid itinerary",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"ATL", "MCO"},
					{"SFO", "ATL"},
					{"MCO", "JFK"},
				},
			},
			want:    []string{"SFO", "ATL", "MCO", "JFK"},
			wantErr: false,
		},
		{
			name: "international multi-city trip",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"JFK", "LHR"}, // New York to London
					{"CDG", "FCO"}, // Paris to Rome
					{"FCO", "ATH"}, // Rome to Athens
					{"LHR", "CDG"}, // London to Paris
				},
			},
			want:    []string{"JFK", "LHR", "CDG", "FCO", "ATH"},
			wantErr: false,
		},
		{
			name: "complex with missing segment",
			request: &models.ItineraryRequest{
				Tickets: []models.TicketPair{
					{"SFO", "LAX"},
					{"LAX", "DFW"},
					// Missing DFW to somewhere
					{"ORD", "JFK"},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.ReconstructItinerary(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconstructItinerary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReconstructItinerary() = %v, want %v", got, tt.want)
			}
		})
	}
}
