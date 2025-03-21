package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/log"
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
				return true // Allow all origins for now
			},
		},
	}
}

// HandleWebSocket handles WebSocket connections
func (h *GameHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get character ID from query params
	characterID := r.URL.Query().Get("characterId")
	if characterID == "" {
		log.Warn("Character ID not provided")
		http.Error(w, "Character ID is required", http.StatusBadRequest)
		return
	}

	// Check if character exists
	character, err := h.characterRepo.GetByID(characterID)
	if err != nil || character == nil {
		log.Warn("Character not found: %s", characterID)
		http.Error(w, "character not found", http.StatusNotFound)
		return
	}

	// Now that we've validated the character, let the game manager handle the connection
	// The game manager will upgrade the connection and manage the WebSocket
	h.manager.HandleConnection(w, r)
}

// StartGameManager starts the game manager
func (h *GameHandler) StartGameManager() {
	// Start any background processes for the game manager
	// This could include things like periodic updates, AI processing, etc.
}
