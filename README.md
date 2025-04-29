# Flight Itinerary API

A Go web application that reconstructs a complete flight itinerary from a sequence of flight tickets.

## Features

Core Features:
- RESTful API endpoint for flight itinerary reconstruction
- Graph-based algorithm for efficient route calculation
- Comprehensive input validation and error handling

Performance Features:
- High-performance **worker pool** for concurrent processing
- Configurable rate limiting for API protection
- Optimized memory management
- Graceful shutdown with 30s timeout for:
  * Completion of in-flight requests
  * Worker pool termination
  * Resource cleanup

Quality Assurance:
- Comprehensive unit test coverage
- Built-in load testing capabilities
- Detailed performance metrics

## Technologies Used

- Go 1.21
- Echo framework
- Standard library packages for core functionality
- High-performance worker pool for concurrent request handling
- Configurable rate limiting
- Comprehensive test coverage of all major components

## Getting Started

### Prerequisites

- Go 1.21 or later

### Installation

1. Clone the repository
```bash
git clone https://github.com/bentalebwael/flight-itinerary-api
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

## Error Handling

The API handles various error cases:
- Invalid JSON format
- Missing or malformed ticket data
- Invalid airport codes (must be 3 characters)
- Disconnected routes
- Multiple starting points
- Multiple flights from the same source

## Configuration

The application can be configured using environment variables:

| Variable | Description | Default Value |
|----------|-------------|---------------|
| SERVER_PORT | Port number for the HTTP server | 8080 |
| WORKER_COUNT | Number of workers in the pool | 500 |
| RATE_LIMITER | Enable/disable rate limiting | disabled |
| MAX_REQUESTS_PER_MIN | Maximum requests per minute per IP | 10 |

Example configuration for high-performance setup:
```bash
export WORKER_COUNT=100
export RATE_LIMITER=enabled
export MAX_REQUESTS_PER_MIN=100
```

Note: The worker pool size can be adjusted based on your server's resources and expected load. The default configuration of 500 workers provides a good balance between performance and resource usage for most use cases.


## Algorithm

The application uses a graph-based approach to reconstruct the itinerary:

1. Builds a directed graph from the flight tickets
2. Identifies the starting point (airport with no incoming flights)
3. Traverses the graph to reconstruct the complete itinerary

Time Complexity: O(n) where n is the number of tickets
Space Complexity: O(n) for storing the graph

## Development Choices

- **Echo Framework**: Chosen for its simplicity, performance, and built-in middleware support
- **Graph-based Solution**: Provides efficient O(n) time complexity for itinerary reconstruction
- **Modular Structure**: Separates concerns into models, services, and handlers for better maintainability
- **Worker Pool**: Implements a highly efficient concurrent request handling system
- **Load Testing**: Built-in load testing using Vegeta for performance validation
- **Test Coverage**: Extensive unit testing covering all critical paths

## Performance

The application utilizes a configurable worker pool architecture that can easily **handle 1 Million requests** per minute. This is achieved through:
- Concurrent request processing using a fixed worker pool
- Efficient memory management
- Optimized request handling pipeline
- Graceful shutdown handling:
  * Zero request loss during shutdown
  * Clean worker pool termination
  * Proper resource cleanup
  * Configurable shutdown timeout (default: 30s)

## Unit Testing

The application maintains comprehensive test coverage across all major components:

- **Models**: Tests for flight data validation and processing
- **Services**: Coverage of itinerary reconstruction logic
- **Handlers**: API endpoint testing with various input scenarios
- **Middleware**: Tests for logging and rate limiting functionality

### Running Tests

Run all tests:
```bash
# Run all tests
go test ./...

# Run all tests with coverage
go test ./... -cover

# Run all tests with verbose output
go test -v ./...
```

The test suite covers:
- Input validation
- Error handling
- Edge cases
- Concurrent operations
- Middleware functionality
- API response formats
- Business logic correctness

## Load Testing

The application includes built-in load testing capabilities using **Vegeta, a versatile HTTP load testing tool**. The load tests can be executed using the provided test suite in the `tests` directory.

### Running Load Tests

To run the load tests:

```bash
# Build and run the API server first
go run main.go

# In a separate terminal, run the load tests
go run tests/main.go [flags]
```

Available flags:
- `-rate`: Number of requests per second (default: 100)
- `-duration`: Test duration in seconds (default: 10)
- `-target`: Target URL (default: http://localhost:8080/api/itinerary)

Example:
```bash
# Run load test with 1000 requests/second for 30 seconds
go run tests/main.go -rate 1000 -duration 30
```

### Test Results

The load test results provide detailed metrics in three categories:

1. **Success Metrics**
   ```
   Success Rate: 100.00%
   Total Requests: 30000
   ```

2. **Timing Metrics**
   ```
   Latency (mean): 1.5ms
   Latency (P50): 1.2ms
   Latency (P90): 2.1ms
   Latency (P99): 3.5ms
   ```

3. **Throughput Metrics**
   ```
   Mean Throughput: 999.82 requests/sec
   Max Throughput: 1000.00 requests/sec
   ```

These metrics help evaluate the application's performance under various load conditions.