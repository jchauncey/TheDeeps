package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// InventoryHandler handles inventory-related API endpoints
type InventoryHandler struct {
	characterRepo *repositories.CharacterRepository
	inventoryRepo *repositories.InventoryRepository
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(characterRepo *repositories.CharacterRepository, inventoryRepo *repositories.InventoryRepository) *InventoryHandler {
	return &InventoryHandler{
		characterRepo: characterRepo,
		inventoryRepo: inventoryRepo,
	}
}

// RegisterRoutes registers the inventory routes
func (h *InventoryHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/characters/{characterID}/inventory", h.GetInventory).Methods("GET")
	router.HandleFunc("/api/characters/{characterID}/inventory/{itemID}", h.GetInventoryItem).Methods("GET")
	router.HandleFunc("/api/characters/{characterID}/inventory/{itemID}/equip", h.EquipItem).Methods("POST")
	router.HandleFunc("/api/characters/{characterID}/inventory/{itemID}/unequip", h.UnequipItem).Methods("POST")
	router.HandleFunc("/api/characters/{characterID}/inventory/{itemID}/use", h.UseItem).Methods("POST")
	router.HandleFunc("/api/characters/{characterID}/equipment", h.GetEquipment).Methods("GET")
	router.HandleFunc("/api/items", h.GetAllItems).Methods("GET")
	router.HandleFunc("/api/items/generate", h.GenerateItems).Methods("POST")
}

// GetInventory returns a character's inventory
func (h *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["characterID"]

	character, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character.Inventory)
}

// GetInventoryItem returns a specific item from a character's inventory
func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["characterID"]
	itemID := vars["itemID"]

	character, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	item, found := character.GetInventoryItem(itemID)
	if !found {
		http.Error(w, "Item not found in inventory", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// EquipItem equips an item from a character's inventory
func (h *InventoryHandler) EquipItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["characterID"]
	itemID := vars["itemID"]

	character, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	success := character.EquipItem(itemID)
	if !success {
		http.Error(w, "Failed to equip item", http.StatusBadRequest)
		return
	}

	// Save the updated character
	h.characterRepo.Save(character)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// UnequipItem unequips an item
func (h *InventoryHandler) UnequipItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["characterID"]
	itemID := vars["itemID"]

	character, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Find the item to determine its type
	item, found := character.GetInventoryItem(itemID)
	if !found {
		// Check if it's equipped
		var itemType models.ItemType
		if character.Equipment.Weapon != nil && character.Equipment.Weapon.ID == itemID {
			itemType = models.ItemWeapon
		} else if character.Equipment.Armor != nil && character.Equipment.Armor.ID == itemID {
			itemType = models.ItemArmor
		} else if character.Equipment.Accessory != nil && character.Equipment.Accessory.ID == itemID {
			itemType = models.ItemArtifact
		} else {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		success := character.UnequipItem(itemType)
		if !success {
			http.Error(w, "Failed to unequip item", http.StatusBadRequest)
			return
		}
	} else {
		success := character.UnequipItem(item.Type)
		if !success {
			http.Error(w, "Failed to unequip item", http.StatusBadRequest)
			return
		}
	}

	// Save the updated character
	h.characterRepo.Save(character)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// UseItem uses an item from a character's inventory
func (h *InventoryHandler) UseItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["characterID"]
	itemID := vars["itemID"]

	character, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	success := character.UseItem(itemID)
	if !success {
		http.Error(w, "Failed to use item", http.StatusBadRequest)
		return
	}

	// Save the updated character
	h.characterRepo.Save(character)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// GetEquipment returns a character's equipped items
func (h *InventoryHandler) GetEquipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID := vars["characterID"]

	character, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character.Equipment)
}

// GetAllItems returns all items in the repository
func (h *InventoryHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items := h.inventoryRepo.GetAllItems()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// GenerateItemsRequest represents a request to generate random items
type GenerateItemsRequest struct {
	Count      int `json:"count"`
	FloorLevel int `json:"floorLevel"`
}

// GenerateItems generates random items
func (h *InventoryHandler) GenerateItems(w http.ResponseWriter, r *http.Request) {
	var req GenerateItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Count <= 0 {
		req.Count = 1
	}

	if req.FloorLevel <= 0 {
		req.FloorLevel = 1
	}

	items := h.inventoryRepo.GenerateRandomItems(req.Count, req.FloorLevel)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
