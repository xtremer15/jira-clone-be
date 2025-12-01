// Package main provides the entry point for the Jira Clone Backend application.
// This application serves as a RESTful API backend for a Jira-like project management system,
// built with Go and following clean architecture principles.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jira-clone-be/internal/config"
	"jira-clone-be/internal/routes"
	"jira-clone-be/pkg/logger"

	"github.com/joho/godotenv"
)

// main is the entry point of the application.
// It initializes the server configuration, sets up logging, configures routes,
// and starts the HTTP server with graceful shutdown support.
//
// The application follows these steps:
//  1. Load environment variables from .env file (if present)
//  2. Initialize application configuration
//  3. Set up structured logging
//  4. Create and configure the routes service
//  5. Start the HTTP server with proper timeouts
//  6. Handle graceful shutdown on system signals
//
// Environment variables can be configured in a .env file or set as system environment variables.
// See .env.example for available configuration options.
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Initialize routes
	routesService := routes.NewRoutes(cfg, logger)
	router := routesService.SetupRoutes()

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info("Starting server on port " + cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: " + err.Error())
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: " + err.Error())
	}

	logger.Info("Server exited")
}
