package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/rs/cors"
)

// Game represents the game state
type Game struct {
	// Add game state here
}

// DebugMessage represents a debug message
type DebugMessage struct {
	Message   string `json:"message"`
	Level     string `json:"level"`
	Timestamp int64  `json:"timestamp"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
	clients             = make(map[*websocket.Conn]bool)
	characterRepository = repositories.NewCharacterRepository()
)

// CreateCharacterRequest represents a request to create a character
type CreateCharacterRequest struct {
	Name           string       `json:"name"`
	CharacterClass string       `json:"characterClass"`
	Stats          models.Stats `json:"stats"`
}

// CharacterResponse represents a character response
type CharacterResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CharacterClass string `json:"characterClass"`
}

func main() {
	r := mux.NewRouter()

	// WebSocket endpoint
	r.HandleFunc("/ws", handleWebSocket)

	// Character endpoints
	r.HandleFunc("/character", handleCreateCharacter).Methods("POST")
	r.HandleFunc("/character/{id}", handleGetCharacter).Methods("GET")
	r.HandleFunc("/characters", handleGetCharacters).Methods("GET")

	// Use CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	// Start server
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Register client
	clients[conn] = true
	defer delete(clients, conn)

	// Send welcome message
	debugMsg := DebugMessage{
		Message:   "Connected to game server",
		Level:     "info",
		Timestamp: time.Now().Unix(),
	}
	conn.WriteJSON(debugMsg)

	// Message handling loop
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Echo the message back for now
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func handleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var req CreateCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" || req.CharacterClass == "" {
		http.Error(w, "Name and class are required", http.StatusBadRequest)
		return
	}

	// Create character
	character := models.NewCharacter(req.Name, req.CharacterClass, req.Stats)
	if err := characterRepository.Create(character); err != nil {
		http.Error(w, "Failed to create character", http.StatusInternalServerError)
		return
	}

	// Log character creation to console
	logMessage := fmt.Sprintf("[%s] Character created: %s (ID: %s)",
		time.Now().Format("2006-01-02 15:04:05"),
		formatCharacterInfo(character),
		character.ID)
	log.Println(logMessage)

	// Return character ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id": character.ID,
	})

	// Broadcast debug message to all clients
	debugMsg := DebugMessage{
		Message:   fmt.Sprintf("New character created: %s (%s)", character.Name, character.CharacterClass),
		Level:     "info",
		Timestamp: time.Now().Unix(),
	}
	broadcastDebugMessage(debugMsg)
}

func handleGetCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	character, err := characterRepository.GetByID(id)
	if err != nil {
		if err == repositories.ErrCharacterNotFound {
			http.Error(w, "Character not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get character", http.StatusInternalServerError)
		}
		return
	}

	// Log character loading to console
	logMessage := fmt.Sprintf("[%s] Character loaded: %s (ID: %s)",
		time.Now().Format("2006-01-02 15:04:05"),
		formatCharacterInfo(character),
		character.ID)
	log.Println(logMessage)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

func handleGetCharacters(w http.ResponseWriter, r *http.Request) {
	characters := characterRepository.GetAll()

	// Log character list request to console
	logMessage := fmt.Sprintf("[%s] Character list requested: %d characters found",
		time.Now().Format("2006-01-02 15:04:05"),
		len(characters))
	log.Println(logMessage)

	// Convert to response format
	response := make([]CharacterResponse, 0, len(characters))
	for _, character := range characters {
		response = append(response, CharacterResponse{
			ID:             character.ID,
			Name:           character.Name,
			CharacterClass: character.CharacterClass,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func broadcastDebugMessage(message DebugMessage) {
	for client := range clients {
		if err := client.WriteJSON(message); err != nil {
			log.Printf("Error broadcasting message: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// Helper function to format character information for logging
func formatCharacterInfo(character *models.Character) string {
	return fmt.Sprintf("%s the %s (STR:%d, DEX:%d, CON:%d, INT:%d, WIS:%d, CHA:%d)",
		character.Name,
		character.CharacterClass,
		character.Stats.Strength,
		character.Stats.Dexterity,
		character.Stats.Constitution,
		character.Stats.Intelligence,
		character.Stats.Wisdom,
		character.Stats.Charisma)
}
