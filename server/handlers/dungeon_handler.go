package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// DungeonHandler handles dungeon-related HTTP requests
type DungeonHandler struct {
	dungeonRepo   *repositories.DungeonRepository
	characterRepo *repositories.CharacterRepository
	mapGenerator  *game.MapGenerator
}

// NewDungeonHandler creates a new dungeon handler
func NewDungeonHandler() *DungeonHandler {
	return &DungeonHandler{
		dungeonRepo:   repositories.NewDungeonRepository(),
		characterRepo: repositories.NewCharacterRepository(),
		mapGenerator:  game.NewMapGenerator(time.Now().UnixNano()),
	}
}

// GetDungeons handles GET /dungeons
func (h *DungeonHandler) GetDungeons(w http.ResponseWriter, r *http.Request) {
	dungeons := h.dungeonRepo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dungeons)
}

// CreateDungeon handles POST /dungeons
func (h *DungeonHandler) CreateDungeon(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request struct {
		Name       string `json:"name"`
		Floors     int    `json:"floors"`
		Difficulty string `json:"difficulty"`
		Seed       int64  `json:"seed,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if request.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if request.Floors < 1 {
		http.Error(w, "Floors must be at least 1", http.StatusBadRequest)
		return
	}

	if request.Difficulty == "" {
		request.Difficulty = "normal" // Default difficulty
	}

	// Create dungeon
	dungeon := models.NewDungeon(request.Name, request.Floors, request.Seed)
	dungeon.Difficulty = request.Difficulty

	// Generate first floor
	floor := dungeon.GenerateFloor(1)
	h.mapGenerator.GenerateFloorWithDifficulty(floor, 1, request.Floors == 1, request.Difficulty)

	// Save dungeon
	if err := h.dungeonRepo.Save(dungeon); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return dungeon
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dungeon)
}

// JoinDungeon handles POST /dungeons/{id}/join
func (h *DungeonHandler) JoinDungeon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dungeonID := vars["id"]

	// Parse request body
	var request struct {
		CharacterID string `json:"characterId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if request.CharacterID == "" {
		http.Error(w, "Character ID is required", http.StatusBadRequest)
		return
	}

	// Get dungeon
	_, err := h.dungeonRepo.GetByID(dungeonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get character
	character, err := h.characterRepo.GetByID(request.CharacterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Add character to dungeon
	if err := h.dungeonRepo.AddCharacterToDungeon(dungeonID, request.CharacterID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get first floor
	floor, err := h.dungeonRepo.GetFloor(dungeonID, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Place character on first floor at the first up stairs
	if len(floor.UpStairs) > 0 {
		character.Position = floor.UpStairs[0]
	} else {
		// If no up stairs, place at a random walkable tile
		for y := 0; y < floor.Height; y++ {
			for x := 0; x < floor.Width; x++ {
				if floor.Tiles[y][x].Walkable && floor.Tiles[y][x].MobID == "" {
					character.Position = models.Position{X: x, Y: y}
					break
				}
			}
			if character.Position.X != 0 || character.Position.Y != 0 {
				break
			}
		}
	}

	// Update character
	character.CurrentFloor = 1
	character.CurrentDungeon = dungeonID

	// Save character
	if err := h.characterRepo.Save(character); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return floor
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(floor)
}

// GetFloor handles GET /dungeons/{id}/floor/{level}
func (h *DungeonHandler) GetFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dungeonID := vars["id"]
	levelStr := vars["level"]

	// Parse level
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		http.Error(w, "Invalid floor level", http.StatusBadRequest)
		return
	}

	// Get dungeon
	dungeon, err := h.dungeonRepo.GetByID(dungeonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate level
	if level < 1 || level > dungeon.Floors {
		http.Error(w, "Floor level out of range", http.StatusBadRequest)
		return
	}

	// Get floor
	floor, err := h.dungeonRepo.GetFloor(dungeonID, level)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If floor hasn't been generated yet, generate it
	if len(floor.Rooms) == 0 {
		h.mapGenerator.GenerateFloorWithDifficulty(floor, level, level == dungeon.Floors, dungeon.Difficulty)
	}

	// Return floor
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(floor)
}
