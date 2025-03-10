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

// sendJSONResponse sends a standardized JSON success response
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
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

	sendJSONResponse(w, map[string]string{
		"id":   character.ID,
		"name": character.Name,
	}, http.StatusCreated)
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

	sendJSONResponse(w, character, http.StatusOK)
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

	sendJSONResponse(w, response, http.StatusOK)
}

// HandleGetFloor handles requests to get a specific floor of a dungeon
func (h *Handler) HandleGetFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dungeonID := vars["id"]
	floorLevel, err := strconv.Atoi(vars["level"])
	if err != nil {
		sendJSONError(w, "Invalid floor level", http.StatusBadRequest)
		return
	}

	// Get character ID from query parameter if provided
	characterID := r.URL.Query().Get("characterId")

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

	// Check if floor level is valid
	if floorLevel < 0 || floorLevel >= len(dungeon.Floors) {
		sendJSONError(w, "Invalid floor level", http.StatusBadRequest)
		return
	}

	// Get the floor
	floor := dungeon.Floors[floorLevel]

	// If character ID is provided, update visibility for that character
	var playerPos models.Position
	var playerData *models.Character
	if characterID != "" {
		character, err := h.Server.CharacterRepository.GetByID(characterID)
		if err != nil {
			if err == repositories.ErrCharacterNotFound {
				sendJSONError(w, "Character not found", http.StatusNotFound)
			} else {
				sendJSONError(w, "Failed to get character", http.StatusInternalServerError)
			}
			return
		}

		// Get player position
		playerPos = character.Position

		// Set player data
		playerData = character
	}

	// Create response
	response := game.FloorMessage{
		Type:         "floor_data",
		Floor:        floor,
		PlayerPos:    playerPos,
		CurrentFloor: floorLevel,
		PlayerData:   playerData,
		DungeonID:    dungeonID,
	}

	sendJSONResponse(w, response, http.StatusOK)
}

// HandleGetCharacterFloor handles requests to get the current floor of a character
func (h *Handler) HandleGetCharacterFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["id"]

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

	// Check if character is in a dungeon
	if character.DungeonID == "" {
		sendJSONError(w, "Character is not in a dungeon", http.StatusBadRequest)
		return
	}

	// Get the dungeon
	dungeon, err := h.Server.DungeonRepository.GetByID(character.DungeonID)
	if err != nil {
		if err == repositories.ErrDungeonNotFound {
			sendJSONError(w, "Dungeon not found", http.StatusNotFound)
		} else {
			sendJSONError(w, "Failed to get dungeon", http.StatusInternalServerError)
		}
		return
	}

	// Check if floor level is valid
	if character.CurrentFloor < 0 || character.CurrentFloor >= len(dungeon.Floors) {
		sendJSONError(w, "Invalid floor level", http.StatusBadRequest)
		return
	}

	// Get the floor
	floor := dungeon.Floors[character.CurrentFloor]

	// Create response
	response := game.FloorMessage{
		Type:         "floor_data",
		Floor:        floor,
		PlayerPos:    character.Position,
		CurrentFloor: character.CurrentFloor,
		PlayerData:   character,
		DungeonID:    character.DungeonID,
	}

	sendJSONResponse(w, response, http.StatusOK)
}

// HandleCreateDungeon handles dungeon creation requests
func (h *Handler) HandleCreateDungeon(w http.ResponseWriter, r *http.Request) {
	var req game.CreateDungeonMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		sendJSONError(w, "Dungeon name is required", http.StatusBadRequest)
		return
	}

	if req.NumFloors <= 0 {
		sendJSONError(w, "Number of floors must be greater than 0", http.StatusBadRequest)
		return
	}

	// Create the dungeon
	dungeon := models.NewDungeon(req.Name, req.NumFloors)
	if err := h.Server.DungeonRepository.Create(dungeon); err != nil {
		sendJSONError(w, "Failed to create dungeon", http.StatusInternalServerError)
		return
	}

	// Initialize dungeon mobs
	h.Server.InitializeDungeonMobs(dungeon)

	log.Printf("[%s] Dungeon created: %s with %d floors (ID: %s)",
		time.Now().Format("2006-01-02 15:04:05"),
		dungeon.Name, len(dungeon.Floors), dungeon.ID)

	sendJSONResponse(w, map[string]string{
		"id":   dungeon.ID,
		"name": dungeon.Name,
	}, http.StatusCreated)
}

// HandleListDungeons handles requests to list all available dungeons
func (h *Handler) HandleListDungeons(w http.ResponseWriter, r *http.Request) {
	dungeons := h.Server.DungeonRepository.GetAll()

	response := make([]game.DungeonListItemResponse, 0, len(dungeons))
	for _, dungeon := range dungeons {
		response = append(response, game.DungeonListItemResponse{
			ID:           dungeon.ID,
			Name:         dungeon.Name,
			PlayerCount:  len(dungeon.Players),
			CreatedAt:    dungeon.CreatedAt,
			LastActivity: dungeon.LastActivity,
		})
	}

	sendJSONResponse(w, response, http.StatusOK)
}

// HandleJoinDungeon handles requests to join a dungeon
func (h *Handler) HandleJoinDungeon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dungeonID := vars["id"]

	// Parse request body
	var req struct {
		CharacterID string `json:"characterId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CharacterID == "" {
		sendJSONError(w, "Character ID is required", http.StatusBadRequest)
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

	// Add character to dungeon
	character.DungeonID = dungeonID
	character.CurrentFloor = 0

	// Get the starting position (center of first room on first floor)
	firstRoom := dungeon.Floors[0].Rooms[0]
	centerX, centerY := firstRoom.Center()
	character.Position = models.Position{X: centerX, Y: centerY}

	// Update character in repository
	if err := h.Server.CharacterRepository.Update(character); err != nil {
		sendJSONError(w, "Failed to update character", http.StatusInternalServerError)
		return
	}

	// Add character to dungeon players
	dungeon.Players[character.ID] = character

	// Update dungeon in repository
	dungeon.LastActivity = time.Now()
	if err := h.Server.DungeonRepository.Update(dungeon); err != nil {
		sendJSONError(w, "Failed to update dungeon", http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Character %s (%s) joined dungeon %s (%s)",
		time.Now().Format("2006-01-02 15:04:05"),
		character.Name, character.ID, dungeon.Name, dungeon.ID)

	// Create response
	response := game.FloorMessage{
		Type:         "floor_data",
		Floor:        dungeon.Floors[0],
		PlayerPos:    character.Position,
		CurrentFloor: 0,
		PlayerData:   character,
		DungeonID:    dungeonID,
	}

	sendJSONResponse(w, response, http.StatusOK)
}

// HandleSaveGame handles requests to save the game state for a character
func (h *Handler) HandleSaveGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["id"]

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

	// Save character state
	if err := h.Server.CharacterRepository.Update(character); err != nil {
		sendJSONError(w, "Failed to save character state", http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Saved game state for character %s (%s)",
		time.Now().Format("2006-01-02 15:04:05"),
		character.Name, character.ID)

	sendJSONResponse(w, map[string]string{
		"message": "Game state saved successfully",
	}, http.StatusOK)
}
