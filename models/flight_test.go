package models

import "testing"

func TestItineraryRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request ItineraryRequest
		wantErr bool
	}{
		{
			name:    "empty tickets",
			request: ItineraryRequest{Tickets: []TicketPair{}},
			wantErr: true,
		},
		{
			name: "invalid ticket format",
			request: ItineraryRequest{
				Tickets: []TicketPair{{"SFO"}},
			},
			wantErr: true,
		},
		{
			name: "empty airport code",
			request: ItineraryRequest{
				Tickets: []TicketPair{{"SFO", ""}},
			},
			wantErr: true,
		},
		{
			name: "invalid airport code length",
			request: ItineraryRequest{
				Tickets: []TicketPair{{"SFOO", "LAX"}},
			},
			wantErr: true,
		},
		{
			name: "valid tickets",
			request: ItineraryRequest{
				Tickets: []TicketPair{{"SFO", "LAX"}, {"LAX", "JFK"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ItineraryRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}