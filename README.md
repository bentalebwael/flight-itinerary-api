# Flight Itinerary API

A Go web application that reconstructs a complete flight itinerary from a sequence of flight tickets.

## Features

- Accepts flight tickets as pairs of airports (source and destination)
- Reconstructs the complete travel itinerary
- Input validation and error handling
- RESTful API endpoint

## Technologies Used

- Go 1.21
- Echo framework
- Standard library packages for core functionality

## Getting Started

### Prerequisites

- Go 1.21 or later

### Installation

1. Clone the repository
```bash
git clone <repository-url>
cd flight-itinerary-api
```

2. Install dependencies
```bash
go mod download
```

3. Run the application
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Usage

### Reconstruct Flight Itinerary

**Endpoint:** `POST /api/itinerary`

**Request Body:**
```json
{
    "tickets": [
        ["LAX", "DXB"],
        ["JFK", "LAX"],
        ["SFO", "SJC"],
        ["DXB", "SFO"]
    ]
}
```

**Success Response:**
```json
{
    "itinerary": ["JFK", "LAX", "DXB", "SFO", "SJC"]
}
```

**Error Response:**
```json
{
    "error": "error message here"
}
```

### Example using cURL

```bash
curl -X POST http://localhost:8080/api/itinerary \
-H "Content-Type: application/json" \
-d '{
    "tickets": [
        ["LAX", "DXB"],
        ["JFK", "LAX"],
        ["SFO", "SJC"],
        ["DXB", "SFO"]
    ]
}'
```

## Algorithm

The application uses a graph-based approach to reconstruct the itinerary:

1. Builds a directed graph from the flight tickets
2. Identifies the starting point (airport with no incoming flights)
3. Traverses the graph to reconstruct the complete itinerary

Time Complexity: O(n) where n is the number of tickets
Space Complexity: O(n) for storing the graph

## Error Handling

The API handles various error cases:
- Invalid JSON format
- Missing or malformed ticket data
- Invalid airport codes (must be 3 characters)
- Disconnected routes
- Multiple starting points
- Multiple flights from the same source

## Development Choices

- **Echo Framework**: Chosen for its simplicity, performance, and built-in middleware support
- **Graph-based Solution**: Provides efficient O(n) time complexity for itinerary reconstruction
- **Modular Structure**: Separates concerns into models, services, and handlers for better maintainability