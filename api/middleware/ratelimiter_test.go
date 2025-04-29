package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"flight-itinerary-api/config"
)

func TestRateLimiterActiveIPs(t *testing.T) {
	// Create config with rate limiting enabled
	cfg := &config.RateLimiterConfig{
		Enabled:       true,
		MaxReqsPerMin: 60,
	}

	// Create rate limiter
	ctx, cancel := context.WithCancel(context.Background())
	rl := NewRateLimiterMiddleware(ctx, cfg)

	// Create echo instance
	e := echo.New()
	handler := rl.Middleware()(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Test IP that's actively making requests
	ip := "1.1.1.1"
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Real-IP", ip)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = handler(c)

		// Sleep briefly between requests
		time.Sleep(100 * time.Millisecond)
	}

	// Verify the IP is still in the map after cleanup
	rl.cleanup()
	assert.Equal(t, 1, len(rl.ips))

	// Cancel context to clean up goroutines
	cancel()

	// Give a small amount of time for cleanup to complete
	time.Sleep(10 * time.Millisecond)

	// Verify the last access time was updated
	entry := rl.ips[ip]
	assert.True(t, time.Since(entry.lastAccess) < time.Second)
}

func TestRateLimiterStop(t *testing.T) {
	// Create config with rate limiting enabled
	cfg := &config.RateLimiterConfig{
		Enabled:       true,
		MaxReqsPerMin: 60,
	}

	// Create rate limiter
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rl := NewRateLimiterMiddleware(ctx, cfg)

	// Add some IPs
	ips := []string{"1.1.1.1", "2.2.2.2"}
	for _, ip := range ips {
		rl.GetLimiter(ip)
	}

	// Cancel context to trigger cleanup
	cancel()

	// Give a small amount of time for cleanup to complete
	time.Sleep(10 * time.Millisecond)

	// Verify cleanup was performed
	assert.Equal(t, 2, len(rl.ips))

	// Verify the done channel is closed
	select {
	case <-rl.done:
		// Channel is closed as expected
	default:
		t.Error("done channel should be closed")
	}
}
