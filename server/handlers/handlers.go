package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// Handler contains dependencies for HTTP handlers
type Handler struct {
	Server *game.GameServer
}

// NewHandler creates a new handler instance
func NewHandler(server *game.GameServer) *Handler {
	return &Handler{
		Server: server,
	}
}

// sendJSONError sends a standardized JSON error response
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "false",
		"message": message,
	})
}

// HandleCreateCharacter handles character creation requests
func (h *Handler) HandleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var req game.CreateCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.CharacterClass == "" {
		sendJSONError(w, "Name and class are required", http.StatusBadRequest)
		return
	}

	character := models.NewCharacter(req.Name, req.CharacterClass, req.Stats)
	if err := h.Server.CharacterRepository.Create(character); err != nil {
		sendJSONError(w, "Failed to create character", http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Character created: %s the %s (STR:%d, DEX:%d, CON:%d, INT:%d, WIS:%d, CHA:%d) (ID: %s)",
		time.Now().Format("2006-01-02 15:04:05"),
		character.Name, character.CharacterClass,
		character.Stats.Strength, character.Stats.Dexterity, character.Stats.Constitution,
		character.Stats.Intelligence, character.Stats.Wisdom, character.Stats.Charisma,
		character.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":   character.ID,
		"name": character.Name,
	})
}

// HandleGetCharacter handles requests to get a character by ID
func (h *Handler) HandleGetCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	character, err := h.Server.CharacterRepository.GetByID(id)
	if err != nil {
		if err == repositories.ErrCharacterNotFound {
			sendJSONError(w, "Character not found", http.StatusNotFound)
		} else {
			sendJSONError(w, "Failed to get character", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

// HandleGetCharacters handles requests to get all characters
func (h *Handler) HandleGetCharacters(w http.ResponseWriter, r *http.Request) {
	characters := h.Server.CharacterRepository.GetAll()

	response := make([]map[string]string, 0, len(characters))
	for _, character := range characters {
		response = append(response, map[string]string{
			"id":             character.ID,
			"name":           character.Name,
			"characterClass": character.CharacterClass,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleGetFloor handles requests to get a specific floor of a dungeon
func (h *Handler) HandleGetFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dungeonID := vars["dungeonId"]
	if dungeonID == "" {
		sendJSONError(w, "Dungeon ID is required", http.StatusBadRequest)
		return
	}

	// Get the dungeon
	dungeon, err := h.Server.DungeonRepository.GetByID(dungeonID)
	if err != nil {
		if err == repositories.ErrDungeonNotFound {
			sendJSONError(w, "Dungeon not found", http.StatusNotFound)
		} else {
			sendJSONError(w, "Failed to get dungeon", http.StatusInternalServerError)
		}
		return
	}

	// Get the floor level
	levelStr := vars["level"]
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		sendJSONError(w, "Invalid floor level", http.StatusBadRequest)
		return
	}

	// Adjust for 1-indexed floors in the API
	level--

	// Check if the floor level is valid
	if level < 0 || level >= len(dungeon.Dungeon.Floors) {
		sendJSONError(w, "Floor level out of range", http.StatusBadRequest)
		return
	}

	// If a character ID is provided, update their floor
	characterID := r.URL.Query().Get("characterId")
	if characterID != "" {
		// Get the character
		character, err := h.Server.CharacterRepository.GetByID(characterID)
		if err != nil {
			if err == repositories.ErrCharacterNotFound {
				sendJSONError(w, "Character not found", http.StatusNotFound)
			} else {
				sendJSONError(w, "Failed to get character", http.StatusInternalServerError)
			}
			return
		}

		// Update character's floor in the dungeon
		dungeon.UpdatePlayerFloor(character.ID, level)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dungeon.Dungeon.Floors[level])
}

// HandleGetCharacterFloor handles requests to get the floor of a specific character
func (h *Handler) HandleGetCharacterFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["characterId"]

	// Get the character
	character, err := h.Server.CharacterRepository.GetByID(characterID)
	if err != nil {
		if err == repositories.ErrCharacterNotFound {
			sendJSONError(w, "Character not found", http.StatusNotFound)
		} else {
			sendJSONError(w, "Failed to get character", http.StatusInternalServerError)
		}
		return
	}

	// Find the dungeon the character is in
	dungeon := h.Server.DungeonRepository.GetPlayerDungeon(characterID)
	if dungeon == nil {
		sendJSONError(w, "Character is not in any dungeon", http.StatusNotFound)
		return
	}

	// Get the character's floor
	floorIndex := dungeon.GetPlayerFloor(characterID)
	if floorIndex < 0 || floorIndex >= len(dungeon.Dungeon.Floors) {
		sendJSONError(w, "Invalid floor index", http.StatusInternalServerError)
		return
	}

	floor := dungeon.Dungeon.Floors[floorIndex]

	// Get the character's position
	position := dungeon.GetPlayerPosition(characterID)

	response := map[string]interface{}{
		"floor":         floor,
		"floorLevel":    floorIndex + 1, // Adjust for 1-indexing in response
		"position":      position,
		"dungeonId":     dungeon.ID,
		"dungeonName":   dungeon.Name,
		"characterId":   character.ID,
		"characterName": character.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleCreateDungeon handles requests to create a new dungeon
func (h *Handler) HandleCreateDungeon(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string `json:"name"`
		NumFloors int    `json:"numFloors"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		sendJSONError(w, "Dungeon name is required", http.StatusBadRequest)
		return
	}

	if req.NumFloors <= 0 {
		sendJSONError(w, "Number of floors must be positive", http.StatusBadRequest)
		return
	}

	// Create the dungeon
	dungeon := h.Server.DungeonRepository.Create(req.Name, req.NumFloors)

	// Initialize mobs on each floor
	h.Server.InitializeGameWorld()

	log.Printf("[%s] Dungeon created: %s with %d floors (ID: %s)",
		time.Now().Format("2006-01-02 15:04:05"),
		dungeon.Name, req.NumFloors, dungeon.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":   dungeon.ID,
		"name": dungeon.Name,
	})
}

// HandleListDungeons handles requests to list all dungeons
func (h *Handler) HandleListDungeons(w http.ResponseWriter, r *http.Request) {
	dungeons := h.Server.DungeonRepository.GetAll()

	response := make([]map[string]interface{}, 0, len(dungeons))
	for _, dungeon := range dungeons {
		response = append(response, map[string]interface{}{
			"id":          dungeon.ID,
			"name":        dungeon.Name,
			"playerCount": len(dungeon.Players),
			"createdAt":   dungeon.CreatedAt,
			"numFloors":   len(dungeon.Dungeon.Floors),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleJoinDungeon handles requests to join a dungeon
func (h *Handler) HandleJoinDungeon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dungeonID := vars["dungeonId"]

	var req struct {
		CharacterID string `json:"characterId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if dungeonID == "" {
		sendJSONError(w, "Dungeon ID is required", http.StatusBadRequest)
		return
	}

	if req.CharacterID == "" {
		sendJSONError(w, "Character ID is required", http.StatusBadRequest)
		return
	}

	// Get the dungeon
	dungeon, err := h.Server.DungeonRepository.GetByID(dungeonID)
	if err != nil {
		if err == repositories.ErrDungeonNotFound {
			sendJSONError(w, "Dungeon not found", http.StatusNotFound)
		} else {
			sendJSONError(w, "Failed to get dungeon", http.StatusInternalServerError)
		}
		return
	}

	// Get the character
	character, err := h.Server.CharacterRepository.GetByID(req.CharacterID)
	if err != nil {
		if err == repositories.ErrCharacterNotFound {
			sendJSONError(w, "Character not found", http.StatusNotFound)
		} else {
			sendJSONError(w, "Failed to get character", http.StatusInternalServerError)
		}
		return
	}

	// Add the player to the dungeon
	if err := h.Server.DungeonRepository.AddPlayerToDungeon(dungeonID, req.CharacterID); err != nil {
		sendJSONError(w, "Failed to add player to dungeon", http.StatusInternalServerError)
		return
	}

	// Get the first floor of the dungeon
	floor := dungeon.Dungeon.Floors[0]

	// Set the player's position to the center of the first room
	firstRoom := floor.Rooms[0]
	centerX, centerY := firstRoom.Center()
	position := models.Position{X: centerX, Y: centerY}

	// Update the player's position in the dungeon
	dungeon.UpdatePlayerPosition(req.CharacterID, position)

	// Update the player's floor in the dungeon (floor 0)
	dungeon.UpdatePlayerFloor(req.CharacterID, 0)

	log.Printf("[%s] Character %s joined dungeon %s",
		time.Now().Format("2006-01-02 15:04:05"),
		character.Name, dungeon.Name)

	// Return the floor data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(floor)
}
