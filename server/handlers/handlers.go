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
	levelStr := vars["level"]

	var level int
	if _, err := fmt.Sscanf(levelStr, "%d", &level); err != nil {
		http.Error(w, "Invalid floor level", http.StatusBadRequest)
		return
	}

	level-- // Adjust for 0-indexing
	if level < 0 || level >= len(h.Server.Game.Dungeon.Floors) {
		http.Error(w, "Floor level out of range", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Server.Game.Dungeon.Floors[level])
}
