package models

import "errors"

// TicketPair represents a single flight ticket with source and destination airports
type TicketPair []string

// ItineraryRequest represents the incoming request containing flight tickets
type ItineraryRequest struct {
    Tickets []TicketPair `json:"tickets"`
}

// ItineraryResponse represents the API response with the ordered itinerary
type ItineraryResponse struct {
    Itinerary []string `json:"itinerary"`
}

// Validate checks if the request contains valid ticket data
func (r *ItineraryRequest) Validate() error {
    if len(r.Tickets) == 0 {
        return errors.New("no tickets provided")
    }

    for _, ticket := range r.Tickets {
        if len(ticket) != 2 {
            return errors.New("invalid ticket format: each ticket must have exactly source and destination")
        }
        if ticket[0] == "" || ticket[1] == "" {
            return errors.New("invalid ticket: airport codes cannot be empty")
        }
        // Basic IATA airport code validation (3 uppercase letters)
        for _, code := range ticket {
            if len(code) != 3 {
                return errors.New("invalid airport code: must be 3 characters")
            }
        }
    }

    return nil
}