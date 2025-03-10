package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// GameHandler handles WebSocket connections for real-time game mechanics
type GameHandler struct {
	manager       *game.GameManager
	characterRepo *repositories.CharacterRepository
	upgrader      websocket.Upgrader
}

// NewGameHandler creates a new game handler
func NewGameHandler() *GameHandler {
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	return &GameHandler{
		manager:       game.NewGameManager(characterRepo, dungeonRepo),
		characterRepo: characterRepo,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
	}
}

// HandleWebSocket handles WebSocket connections
func (h *GameHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get character ID from query parameters
	characterID := r.URL.Query().Get("characterId")
	if characterID == "" {
		http.Error(w, "Character ID is required", http.StatusBadRequest)
		return
	}

	// Get character
	character, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Create client
	client := &game.Client{
		ID:         uuid.New().String(),
		Connection: conn,
		Character:  character,
		Send:       make(chan game.Message, 256),
		Manager:    h.manager,
	}

	// Register client with game manager
	h.manager.Register <- client

	// Start client
	go client.Run()
}

// StartGameManager starts the game manager
func (h *GameHandler) StartGameManager() {
	go h.manager.Start()
}
