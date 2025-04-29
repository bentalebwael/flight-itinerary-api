package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"

	"flight-itinerary-api/config"
)

// IPRateLimiter stores rate limiters for different IPs
type IPRateLimiter struct {
	rateLimiterEnabled bool
	ips                map[string]*rateLimiterEntry
	mu                 sync.RWMutex
	rate               rate.Limit
	burst              int
	done               chan struct{} // Channel to signal cleanup goroutine to stop
}

// rateLimiterEntry stores the rate limiter and its last access time
type rateLimiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

// cleanup duration constants
const (
	cleanupInterval = 5 * time.Minute // How often cleanup runs
)

// NewRateLimiterMiddleware creates a new rate limiter middleware instance
func NewRateLimiterMiddleware(ctx context.Context, cfg *config.RateLimiterConfig) *IPRateLimiter {
	rl := &IPRateLimiter{
		rateLimiterEnabled: cfg.Enabled,
		rate:               rate.Limit(cfg.MaxReqsPerMin) / 60,
		burst:              cfg.MaxReqsPerMin,
		ips:                make(map[string]*rateLimiterEntry),
		done:               make(chan struct{}),
	}

	// Start cleanup goroutine if rate limiting is enabled
	if cfg.Enabled {
		go rl.cleanupLoop(ctx)
	}

	return rl
}

// GetLimiter returns rate limiter for provided IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	entry, exists := i.ips[ip]
	now := time.Now()

	if !exists {
		entry = &rateLimiterEntry{
			limiter:    rate.NewLimiter(i.rate, i.burst),
			lastAccess: now,
		}
		i.ips[ip] = entry
	} else {
		entry.lastAccess = now
	}

	return entry.limiter
}

// Middleware returns the Echo middleware handler
func (i *IPRateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !i.rateLimiterEnabled {
				return next(c)
			}
			ip := c.RealIP()
			if !i.GetLimiter(ip).Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests",
				})
			}
			return next(c)
		}
	}
}

// cleanup removes entries that haven't been accessed for longer than maxIdleTime
func (i *IPRateLimiter) cleanup() {
	i.mu.Lock()
	defer i.mu.Unlock()

	threshold := time.Now().Add(-cleanupInterval)
	for ip, entry := range i.ips {
		if entry.lastAccess.Before(threshold) {
			delete(i.ips, ip)
		}
	}
}

// cleanupLoop runs the cleanup process periodically
func (i *IPRateLimiter) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			i.cleanup()
		case <-i.done:
			return
		case <-ctx.Done():
			i.Stop()
			return
		}
	}
}

// Stop gracefully stops the rate limiter cleanup goroutine
func (i *IPRateLimiter) Stop() {
	if i.rateLimiterEnabled {
		close(i.done)
		// Perform one final cleanup
		i.cleanup()
	}
}
