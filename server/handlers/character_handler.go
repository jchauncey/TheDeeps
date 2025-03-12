package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// CharacterHandler handles character-related HTTP requests
type CharacterHandler struct {
	characterRepo *repositories.CharacterRepository
}

// NewCharacterHandler creates a new character handler
func NewCharacterHandler(characterRepo *repositories.CharacterRepository) *CharacterHandler {
	return &CharacterHandler{
		characterRepo: characterRepo,
	}
}

// GetCharacters handles GET /characters
func (h *CharacterHandler) GetCharacters(w http.ResponseWriter, r *http.Request) {
	characters := h.characterRepo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}

// GetCharacter handles GET /characters/{id}
func (h *CharacterHandler) GetCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	character, err := h.characterRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

// CreateCharacter handles POST /characters
func (h *CharacterHandler) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	// Check if we've reached the character limit
	if h.characterRepo.Count() >= 10 {
		http.Error(w, "Maximum number of characters reached (10)", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request struct {
		Name       string                `json:"name"`
		Class      models.CharacterClass `json:"class"`
		Attributes models.Attributes     `json:"attributes,omitempty"`
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

	// Create character
	character := models.NewCharacter(request.Name, request.Class)

	// Apply custom attributes if provided
	if request.Attributes.Strength > 0 {
		character.Attributes = request.Attributes

		// Recalculate HP and Mana based on new attributes
		character.MaxHP = 10 + character.Attributes.Constitution
		character.CurrentHP = character.MaxHP

		switch character.Class {
		case models.Warrior, models.Barbarian:
			character.MaxHP += 5
			character.MaxMana = 0
		case models.Mage, models.Sorcerer, models.Warlock:
			character.MaxHP -= 2
			character.MaxMana = 10 + character.Attributes.Intelligence
		case models.Cleric, models.Druid:
			character.MaxMana = 10 + character.Attributes.Wisdom
		case models.Bard:
			character.MaxMana = 8 + character.Attributes.Charisma
		case models.Paladin:
			character.MaxMana = 5 + character.Attributes.Charisma
		}

		character.CurrentMana = character.MaxMana
	}

	// Save character
	if err := h.characterRepo.Save(character); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return character
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(character)
}

// DeleteCharacter handles DELETE /characters/{id}
func (h *CharacterHandler) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.characterRepo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SaveCharacter handles POST /characters/{id}/save
func (h *CharacterHandler) SaveCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get character
	character, err := h.characterRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Parse request body
	var request struct {
		Position       models.Position `json:"position"`
		CurrentHP      int             `json:"currentHp"`
		CurrentMana    int             `json:"currentMana"`
		Gold           int             `json:"gold"`
		Experience     int             `json:"experience"`
		CurrentFloor   int             `json:"currentFloor"`
		CurrentDungeon string          `json:"currentDungeon"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update character
	character.Position = request.Position
	character.CurrentHP = request.CurrentHP
	character.CurrentMana = request.CurrentMana
	character.Gold = request.Gold
	character.Experience = request.Experience
	character.CurrentFloor = request.CurrentFloor
	character.CurrentDungeon = request.CurrentDungeon

	// Save character
	if err := h.characterRepo.Save(character); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return character
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

// GetCharacterFloor handles GET /characters/{id}/floor
func (h *CharacterHandler) GetCharacterFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get character
	character, err := h.characterRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return floor
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"floor": character.CurrentFloor,
	})
}
