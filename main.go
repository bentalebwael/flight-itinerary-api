package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"flight-itinerary-api/api"
	"flight-itinerary-api/config"
	"flight-itinerary-api/services"
)

func main() {
	// Initialize logger
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize the logger: %v", err)
	}
	defer l.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Echo
	e := echo.New()

	// Create app context that will be canceled on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize services with configured worker count
	itineraryService := services.NewItineraryService(ctx, cfg)

	// Setup router
	router := api.NewRouter(ctx, cfg, l, itineraryService)
	router.SetupRoutes(e)

	// Start server in a goroutine
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
			l.Info("Shutting down the server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	l.Info("Received shutdown signal")

	// Create a deadline for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop accepting new requests
	if err := e.Shutdown(shutdownCtx); err != nil {
		l.Error("Failed to shutdown server gracefully", zap.Error(err))
	} else {
		l.Info("Server shut down gracefully")
	}

	// Final cleanup
	if err := l.Sync(); err != nil {
		log.Printf("Failed to sync logger: %v", err)
	}
}
