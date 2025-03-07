package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/handlers"
	"github.com/rs/cors"
)

func main() {
	// Create game server
	server := game.NewGameServer()

	// Create HTTP handlers
	handler := handlers.NewHandler(server)

	// Setup routes
	router := server.SetupRoutes(handler)

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start server
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), corsHandler.Handler(router)))
}
