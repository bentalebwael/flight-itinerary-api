package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogger(t *testing.T) {
	// Create an observer for the zap logger
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	// Create echo instance
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	mw := Logger(logger)
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Execute request
	err := handler(c)
	assert.NoError(t, err)

	// Check logs
	assert.Equal(t, 1, recorded.Len())
	log := recorded.All()[0]
	assert.Equal(t, "HTTP Request", log.Message)
	assert.Equal(t, "GET", log.ContextMap()["method"])
	assert.Equal(t, "/test", log.ContextMap()["uri"])
	assert.EqualValues(t, 200, log.ContextMap()["status"])
}
