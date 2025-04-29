package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"flight-itinerary-api/tests/attack"
	"flight-itinerary-api/tests/payloads"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	// Command line flags for test configuration
	rate := flag.Int("rate", 100, "Requests per second")
	duration := flag.Int("duration", 10, "Test duration in seconds")
	target := flag.String("target", "http://localhost:8080/api/itinerary", "Target URL")
	flag.Parse()

	fmt.Printf("\nStarting load test with rate: %d req/s, duration: %ds, target: %s\n",
		*rate, *duration, *target)

	// Configure the load test
	cfg := attack.Config{
		Rate:     *rate,
		Duration: time.Duration(*duration) * time.Second,
		Target:   *target,
	}

	// Get test payloads
	testPayloads := payloads.GetTestPayloads()

	// Execute the load test
	fmt.Println("Executing load test...")
	metrics, err := attack.ExecuteAttack(cfg, testPayloads)
	if err != nil {
		log.Fatalf("Failed to execute load test: %v", err)
	}

	// Print the formatted results
	printMetrics(metrics)
}

// printMetrics formats and prints the load test metrics in a readable way
func printMetrics(metrics *vegeta.Metrics) {
	fmt.Println("\n=== Load Test Results ===")
	fmt.Println("Success Metrics:")
	fmt.Printf("  Success Rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("  Total Requests: %d\n", metrics.Requests)

	fmt.Println("\nTiming Metrics:")
	fmt.Printf("  Latency (mean): %s\n", metrics.Latencies.Mean)
	fmt.Printf("  Latency (P50): %s\n", metrics.Latencies.P50)
	fmt.Printf("  Latency (P90): %s\n", metrics.Latencies.P90)
	fmt.Printf("  Latency (P99): %s\n", metrics.Latencies.P99)

	fmt.Println("\nThroughput Metrics:")
	fmt.Printf("  Mean Throughput: %.2f requests/sec\n", metrics.Throughput)
	fmt.Printf("  Max Throughput: %.2f requests/sec\n", metrics.Rate)
	fmt.Printf("  Total Bytes Read: %d bytes\n", metrics.BytesIn.Total)
	fmt.Printf("  Total Bytes Written: %d bytes\n\n", metrics.BytesOut.Total)

	if len(metrics.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range metrics.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}
}
