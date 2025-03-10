package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/handlers"
	"github.com/rs/cors"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "port to run the server on")
	flag.Parse()

	// Create a new router
	router := mux.NewRouter()

	// Initialize handlers
	characterHandler := handlers.NewCharacterHandler()
	dungeonHandler := handlers.NewDungeonHandler()
	gameHandler := handlers.NewGameHandler()

	// Start the game manager
	gameHandler.StartGameManager()

	// Character endpoints
	router.HandleFunc("/characters", characterHandler.GetCharacters).Methods("GET")
	router.HandleFunc("/characters/{id}", characterHandler.GetCharacter).Methods("GET")
	router.HandleFunc("/characters", characterHandler.CreateCharacter).Methods("POST")
	router.HandleFunc("/characters/{id}", characterHandler.DeleteCharacter).Methods("DELETE")
	router.HandleFunc("/characters/{id}/save", characterHandler.SaveCharacter).Methods("POST")
	router.HandleFunc("/characters/{id}/floor", characterHandler.GetCharacterFloor).Methods("GET")

	// Dungeon endpoints
	router.HandleFunc("/dungeons", dungeonHandler.GetDungeons).Methods("GET")
	router.HandleFunc("/dungeons", dungeonHandler.CreateDungeon).Methods("POST")
	router.HandleFunc("/dungeons/{id}/join", dungeonHandler.JoinDungeon).Methods("POST")
	router.HandleFunc("/dungeons/{id}/floor/{level}", dungeonHandler.GetFloor).Methods("GET")

	// WebSocket endpoint for real-time game mechanics
	router.HandleFunc("/ws", gameHandler.HandleWebSocket)

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Update with your client URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + *port,
		Handler: c.Handler(router),
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on port %s...\n", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal
	<-stop

	log.Println("Shutting down server...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited")
}
