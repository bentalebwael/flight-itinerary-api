package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Middleware returns the Echo middleware function
func Logger(logger *zap.Logger) echo.MiddlewareFunc {
	// Create a new logger without caller information
	l := logger.WithOptions(zap.WithCaller(false))
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			latency := time.Since(start)

			// Log the request details
			l.Info("HTTP Request",
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("client_ip", c.RealIP()),
				zap.Int("status", res.Status),
				zap.Duration("latency", latency),
				zap.String("user_agent", req.UserAgent()),
			)

			return err
		}
	}
}
