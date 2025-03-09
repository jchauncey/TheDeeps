package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

// HandleCreateCharacter handles character creation requests
func (h *Handler) HandleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var req game.CreateCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.CharacterClass == "" {
		http.Error(w, "Name and class are required", http.StatusBadRequest)
		return
	}

	character := models.NewCharacter(req.Name, req.CharacterClass, req.Stats)
	if err := h.Server.CharacterRepository.Create(character); err != nil {
		http.Error(w, "Failed to create character", http.StatusInternalServerError)
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
	json.NewEncoder(w).Encode(map[string]string{"id": character.ID})
}

// HandleGetCharacter handles requests to get a character by ID
func (h *Handler) HandleGetCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	character, err := h.Server.CharacterRepository.GetByID(id)
	if err != nil {
		if err == repositories.ErrCharacterNotFound {
			http.Error(w, "Character not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get character", http.StatusInternalServerError)
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

// HandleGetFloor handles requests to get a specific floor
func (h *Handler) HandleGetFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dungeonID := vars["dungeonId"]
	levelStr := vars["level"]
	characterID := r.URL.Query().Get("characterId")

	// Validate parameters
	if dungeonID == "" {
		http.Error(w, "Dungeon ID is required", http.StatusBadRequest)
		return
	}

	// Get the dungeon
	dungeon, err := h.Server.DungeonRepository.GetByID(dungeonID)
	if err != nil {
		if err == repositories.ErrDungeonNotFound {
			http.Error(w, "Dungeon not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get dungeon", http.StatusInternalServerError)
		}
		return
	}

	var level int
	if _, err := fmt.Sscanf(levelStr, "%d", &level); err != nil {
		http.Error(w, "Invalid floor level", http.StatusBadRequest)
		return
	}

	level-- // Adjust for 0-indexing
	if level < 0 || level >= len(dungeon.Dungeon.Floors) {
		http.Error(w, "Floor level out of range", http.StatusBadRequest)
		return
	}

	// If character ID is provided, update the character's current floor
	if characterID != "" {
		// Check if character exists
		character, err := h.Server.CharacterRepository.GetByID(characterID)
		if err != nil {
			if err == repositories.ErrCharacterNotFound {
				http.Error(w, "Character not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to get character", http.StatusInternalServerError)
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
			http.Error(w, "Character not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get character", http.StatusInternalServerError)
		}
		return
	}

	// Find the dungeon the character is in
	dungeon := h.Server.DungeonRepository.GetPlayerDungeon(characterID)
	if dungeon == nil {
		http.Error(w, "Character is not in any dungeon", http.StatusNotFound)
		return
	}

	// Get the character's floor
	floorIndex := dungeon.GetPlayerFloor(characterID)
	if floorIndex < 0 || floorIndex >= len(dungeon.Dungeon.Floors) {
		http.Error(w, "Invalid floor index", http.StatusInternalServerError)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Dungeon name is required", http.StatusBadRequest)
		return
	}

	if req.NumFloors <= 0 {
		http.Error(w, "Number of floors must be positive", http.StatusBadRequest)
		return
	}

	// Create the dungeon
	dungeon := h.Server.DungeonRepository.Create(req.Name, req.NumFloors)

	// Initialize mobs on all floors
	for i := range dungeon.Dungeon.Floors {
		h.Server.MobSpawner.SpawnMobsOnFloor(dungeon.Dungeon, i)
	}

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
	characterID := r.URL.Query().Get("characterId")

	// Validate parameters
	if dungeonID == "" {
		http.Error(w, "Dungeon ID is required", http.StatusBadRequest)
		return
	}

	if characterID == "" {
		http.Error(w, "Character ID is required", http.StatusBadRequest)
		return
	}

	// Check if dungeon exists
	dungeon, err := h.Server.DungeonRepository.GetByID(dungeonID)
	if err != nil {
		if err == repositories.ErrDungeonNotFound {
			http.Error(w, "Dungeon not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get dungeon", http.StatusInternalServerError)
		}
		return
	}

	// Check if character exists
	character, err := h.Server.CharacterRepository.GetByID(characterID)
	if err != nil {
		if err == repositories.ErrCharacterNotFound {
			http.Error(w, "Character not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get character", http.StatusInternalServerError)
		}
		return
	}

	// Add character to dungeon (this will place them on the first floor)
	if err := h.Server.DungeonRepository.AddPlayerToDungeon(dungeonID, characterID); err != nil {
		http.Error(w, "Failed to add player to dungeon", http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Character %s (%s) joined dungeon %s (%s) on floor 1",
		time.Now().Format("2006-01-02 15:04:05"),
		character.Name, character.ID, dungeon.Name, dungeon.ID)

	// Get the character's position and floor (should be floor 0)
	position := dungeon.GetPlayerPosition(character.ID)
	floorIndex := dungeon.GetPlayerFloor(character.ID)
	floor := dungeon.Dungeon.Floors[floorIndex]

	response := map[string]interface{}{
		"dungeonId":   dungeon.ID,
		"dungeonName": dungeon.Name,
		"characterId": character.ID,
		"floor":       floor,
		"floorLevel":  floorIndex + 1, // Adjust for 1-indexing in response
		"position":    position,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
