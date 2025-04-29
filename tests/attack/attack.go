package attack

import (
	"encoding/json"
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Config holds the configuration for the load test
type Config struct {
	Rate     int
	Duration time.Duration
	Target   string
}

// CreateTarget creates a Vegeta target from a URL and payload
func CreateTarget(url string, payload map[string]interface{}) (vegeta.Target, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return vegeta.Target{}, fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	return vegeta.Target{
		Method: "POST",
		URL:    url,
		Body:   jsonPayload,
		Header: map[string][]string{
			"Content-Type":   {"application/json"},
			"Accept":         {"application/json"},
			"Content-Length": {fmt.Sprintf("%d", len(jsonPayload))},
		},
	}, nil
}

// ExecuteAttack performs the load test with the given configuration and payloads
func ExecuteAttack(cfg Config, payloads []map[string]interface{}) (*vegeta.Metrics, error) {
	targets := make([]vegeta.Target, 0, len(payloads))
	for _, payload := range payloads {
		target, err := CreateTarget(cfg.Target, payload)
		if err != nil {
			return nil, fmt.Errorf("failed to create target: %v", err)
		}
		targets = append(targets, target)
	}

	// Create attack rate
	rate := vegeta.Rate{Freq: cfg.Rate, Per: time.Second}

	// Create target iterator - randomizes the selection from our target array
	targeter := vegeta.NewStaticTargeter(targets...)

	// Create attacker with default options
	attacker := vegeta.NewAttacker()

	// Run the attack and capture metrics
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, cfg.Duration, "Flight Itinerary Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	return &metrics, nil
}
