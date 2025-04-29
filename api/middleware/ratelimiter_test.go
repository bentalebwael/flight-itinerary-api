package middleware

import (
	"encoding/json"
	"flight-itinerary-api/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiterMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		enabled       bool
		maxReqsPerMin int
		requests      int
		expectedCode  int
	}{
		{
			name:          "Rate limiter disabled",
			enabled:       false,
			maxReqsPerMin: 60,
			requests:      100,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Under rate limit",
			enabled:       true,
			maxReqsPerMin: 60,
			requests:      30,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Exceeds rate limit",
			enabled:       true,
			maxReqsPerMin: 2,
			requests:      5,
			expectedCode:  http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create rate limiter middleware
			cfg := &config.RateLimiterConfig{
				Enabled:       tt.enabled,
				MaxReqsPerMin: tt.maxReqsPerMin,
			}
			rateLimiter := NewRateLimiterMiddleware(cfg)

			// Create echo instance
			e := echo.New()
			handler := rateLimiter.Middleware()(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})

			lastStatusCode := 0
			// Send multiple requests
			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				_ = handler(c)
				lastStatusCode = rec.Code
			}

			assert.Equal(t, tt.expectedCode, lastStatusCode)
		})
	}
}

func TestRateLimiterPerIP(t *testing.T) {
	cfg := &config.RateLimiterConfig{
		Enabled:       true,
		MaxReqsPerMin: 2,
	}
	rateLimiter := NewRateLimiterMiddleware(cfg)

	// Create echo instance
	e := echo.New()
	handler := rateLimiter.Middleware()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Test different IPs
	ips := []string{"1.1.1.1", "2.2.2.2"}

	for _, ip := range ips {
		// Send requests for each IP
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("X-Real-IP", ip)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = handler(c)

			if i < 2 {
				assert.Equal(t, http.StatusOK, rec.Code)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, rec.Code)
				var response map[string]string
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "Too many requests", response["error"])
			}
		}

		// Wait a bit to allow rate limiter to reset
		time.Sleep(time.Second)
	}
}

func TestGetLimiter(t *testing.T) {
	cfg := &config.RateLimiterConfig{
		Enabled:       true,
		MaxReqsPerMin: 60,
	}
	rateLimiter := NewRateLimiterMiddleware(cfg)

	// Get limiter for the same IP twice
	ip := "1.1.1.1"
	limiter1 := rateLimiter.GetLimiter(ip)
	limiter2 := rateLimiter.GetLimiter(ip)

	// Should return the same limiter instance
	assert.Equal(t, limiter1, limiter2)

	// Check if limiter has correct rate and burst
	assert.Equal(t, rate.Limit(1), limiter1.Limit())
	assert.Equal(t, 60, limiter1.Burst())
}
