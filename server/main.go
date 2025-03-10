package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jchauncey/TheDeeps/server/log"
	"github.com/rs/cors"
)

func main() {
	// Initialize logger
	logLevel := log.InfoLevel
	if os.Getenv("DEBUG") == "true" {
		logLevel = log.DebugLevel
		log.Info("Debug logging enabled")
	}
	log.SetLevel(logLevel)
	log.Info("Logger initialized with level: %s", log.LevelNames[logLevel])

	// Parse command line flags
	port := flag.String("port", "8080", "port to run the server on")
	flag.Parse()

	// Create and set up server
	server := NewServer()
	server.SetupRoutes()

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Update with your client URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + *port,
		Handler: c.Handler(server.router),
	}

	// Start the server in a goroutine
	go func() {
		log.Info("Server starting on port %s...", *port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Could not start server: %v", err)
		}
	}()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal
	<-stop

	log.Info("Shutting down server...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}
