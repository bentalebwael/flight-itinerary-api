package middleware

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"

	"flight-itinerary-api/config"
)

// IPRateLimiter stores rate limiters for different IPs
type IPRateLimiter struct {
	rateLimiterEnabled bool
	ips                map[string]*rate.Limiter
	mu                 sync.RWMutex
	rate               rate.Limit
	burst              int
}

// NewRateLimiterMiddleware creates a new rate limiter middleware instance
func NewRateLimiterMiddleware(cfg *config.RateLimiterConfig) *IPRateLimiter {
	return &IPRateLimiter{
		rateLimiterEnabled: cfg.Enabled,
		rate:               rate.Limit(cfg.MaxReqsPerMin) / 60,
		burst:              cfg.MaxReqsPerMin,
		ips:                make(map[string]*rate.Limiter),
	}
}

// GetLimiter returns rate limiter for provided IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.ips[ip] = limiter
	}

	return limiter
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
