package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/log"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// CombatHandler handles combat-related WebSocket messages
type CombatHandler struct {
	characterRepo *repositories.CharacterRepository
	dungeonRepo   *repositories.DungeonRepository
	gameManager   *game.GameManager
	combatManager *game.CombatManager
	upgrader      websocket.Upgrader
}

// NewCombatHandler creates a new combat handler
func NewCombatHandler(characterRepo *repositories.CharacterRepository, dungeonRepo *repositories.DungeonRepository, gameManager *game.GameManager) *CombatHandler {
	return &CombatHandler{
		characterRepo: characterRepo,
		dungeonRepo:   dungeonRepo,
		gameManager:   gameManager,
		combatManager: game.NewCombatManager(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
	}
}

// CombatMessage represents a combat action from the client
type CombatMessage struct {
	Action      string `json:"action"`
	CharacterID string `json:"characterId"`
	MobID       string `json:"mobId,omitempty"`
	ItemID      string `json:"itemId,omitempty"`
}

// CombatResponse represents the server's response to a combat action
type CombatResponse struct {
	Action  string            `json:"action"`
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Result  game.CombatResult `json:"result,omitempty"`
}

// HandleCombat handles WebSocket connections for combat
func (h *CombatHandler) HandleCombat(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Main message loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("WebSocket error: %v", err)
			break
		}

		// Parse the combat message
		var combatMsg CombatMessage
		if err := json.Unmarshal(message, &combatMsg); err != nil {
			log.Error("Failed to parse combat message: %v", err)
			continue
		}

		// Get the character
		character, err := h.characterRepo.GetByID(combatMsg.CharacterID)
		if err != nil {
			log.Warn("Character not found: %s", combatMsg.CharacterID)
			response := CombatResponse{
				Action:  combatMsg.Action,
				Success: false,
				Message: "Character not found",
			}
			sendResponse(conn, response)
			continue
		}

		// Handle the combat action
		var response CombatResponse
		switch combatMsg.Action {
		case "attack":
			response = h.handleAttack(character, combatMsg.MobID)
		case "useItem":
			response = h.handleUseItem(character, combatMsg.ItemID)
		case "flee":
			response = h.handleFlee(character, combatMsg.MobID)
		default:
			response = CombatResponse{
				Action:  combatMsg.Action,
				Success: false,
				Message: "Unknown action",
			}
		}

		// Send the response
		sendResponse(conn, response)
	}
}

// sendResponse sends a combat response to the client
func sendResponse(conn *websocket.Conn, response CombatResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Error("Failed to marshal response: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Error("Failed to send response: %v", err)
	}
}

// handleAttack processes an attack action
func (h *CombatHandler) handleAttack(character *models.Character, mobID string) CombatResponse {
	// Get dungeon and floor
	dungeon, err := h.dungeonRepo.GetByID(character.CurrentDungeon)
	if err != nil {
		return CombatResponse{
			Success: false,
			Message: "Dungeon not found",
		}
	}

	// Get floor
	floorLevel := dungeon.GetCharacterFloor(character.ID)
	floor, err := h.dungeonRepo.GetFloor(dungeon.ID, floorLevel)
	if err != nil {
		return CombatResponse{
			Success: false,
			Message: "Floor not found",
		}
	}

	// Get mob
	mob, exists := floor.Mobs[mobID]
	if !exists {
		return CombatResponse{
			Success: false,
			Message: "Mob not found",
		}
	}

	// Check if character is adjacent to mob
	if !isAdjacent(character.Position, mob.Position) {
		return CombatResponse{
			Success: false,
			Message: "Not adjacent to mob",
		}
	}

	// Process attack
	result := h.combatManager.AttackMob(character, mob)

	// Update mob in floor data
	if result.Killed {
		delete(floor.Mobs, mobID)
	} else {
		floor.Mobs[mobID] = mob
	}

	// Save character and floor
	h.characterRepo.Save(character)
	h.dungeonRepo.SaveFloor(dungeon.ID, floorLevel, floor)

	// Send response
	return CombatResponse{
		Action:  "attack",
		Success: result.Success,
		Message: result.Message,
		Result:  result,
	}
}

// handleUseItem processes a use item action
func (h *CombatHandler) handleUseItem(character *models.Character, itemID string) CombatResponse {
	// TODO: Implement inventory system
	// For now, just send an error
	return CombatResponse{
		Action:  "useItem",
		Success: false,
		Message: "Inventory system not implemented yet",
	}
}

// handleFlee processes a flee action
func (h *CombatHandler) handleFlee(character *models.Character, mobID string) CombatResponse {
	// Get dungeon and floor
	dungeon, err := h.dungeonRepo.GetByID(character.CurrentDungeon)
	if err != nil {
		return CombatResponse{
			Success: false,
			Message: "Dungeon not found",
		}
	}

	// Get floor
	floorLevel := dungeon.GetCharacterFloor(character.ID)
	floor, err := h.dungeonRepo.GetFloor(dungeon.ID, floorLevel)
	if err != nil {
		return CombatResponse{
			Success: false,
			Message: "Floor not found",
		}
	}

	// Get mob
	mob, exists := floor.Mobs[mobID]
	if !exists {
		return CombatResponse{
			Success: false,
			Message: "Mob not found",
		}
	}

	// Check if character is adjacent to mob
	if !isAdjacent(character.Position, mob.Position) {
		return CombatResponse{
			Success: false,
			Message: "Not adjacent to mob",
		}
	}

	// Process flee attempt
	result := h.combatManager.Flee(character, mob)

	// Save character
	h.characterRepo.Save(character)

	// Send response
	return CombatResponse{
		Action:  "flee",
		Success: result.Success,
		Message: result.Message,
		Result:  result,
	}
}

// Helper functions

// isAdjacent checks if two positions are adjacent (including diagonals)
func isAdjacent(pos1, pos2 models.Position) bool {
	dx := abs(pos1.X - pos2.X)
	dy := abs(pos1.Y - pos2.Y)
	return dx <= 1 && dy <= 1 && !(dx == 0 && dy == 0)
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// findSafePosition finds a safe position not adjacent to any mob
func findSafePosition(floor *models.Floor, currentPos models.Position) models.Position {
	// Start with current position
	pos := currentPos

	// Try to find a safe position within 5 tiles
	for distance := 1; distance <= 5; distance++ {
		// Try positions in a spiral pattern
		for dx := -distance; dx <= distance; dx++ {
			for dy := -distance; dy <= distance; dy++ {
				// Skip positions that are not on the edge of the spiral
				if abs(dx) != distance && abs(dy) != distance {
					continue
				}

				// Calculate new position
				newPos := models.Position{
					X: pos.X + dx,
					Y: pos.Y + dy,
				}

				// Check if position is valid
				if newPos.X < 0 || newPos.X >= floor.Width || newPos.Y < 0 || newPos.Y >= floor.Height {
					continue
				}

				// Check if position is walkable
				if !floor.Tiles[newPos.Y][newPos.X].Walkable {
					continue
				}

				// Check if position is safe (not adjacent to any mob)
				safe := true
				for _, mob := range floor.Mobs {
					if isAdjacent(newPos, mob.Position) {
						safe = false
						break
					}
				}

				if safe {
					return newPos
				}
			}
		}
	}

	// If no safe position found, return current position
	return pos
}

// GetCombatState handles GET /characters/{id}/combat
func (h *CombatHandler) GetCombatState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["id"]

	// Get character
	character, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get dungeon
	dungeon, err := h.dungeonRepo.GetByID(character.CurrentDungeon)
	if err != nil {
		http.Error(w, "Character is not in a dungeon", http.StatusBadRequest)
		return
	}

	// Get floor
	floorLevel := dungeon.GetCharacterFloor(character.ID)
	floor, err := h.dungeonRepo.GetFloor(dungeon.ID, floorLevel)
	if err != nil {
		http.Error(w, "Floor not found", http.StatusInternalServerError)
		return
	}

	// Find nearby mobs (adjacent to character)
	nearbyMobs := make(map[string]*models.Mob)
	for mobID, mob := range floor.Mobs {
		if isAdjacent(character.Position, mob.Position) {
			nearbyMobs[mobID] = mob
		}
	}

	// Create response
	response := struct {
		Character  *models.Character      `json:"character"`
		NearbyMobs map[string]*models.Mob `json:"nearbyMobs"`
		InCombat   bool                   `json:"inCombat"`
	}{
		Character:  character,
		NearbyMobs: nearbyMobs,
		InCombat:   len(nearbyMobs) > 0,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
