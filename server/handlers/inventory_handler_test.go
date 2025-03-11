package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/stretchr/testify/assert"
)

func TestInventoryHandler(t *testing.T) {
	// Setup
	characterRepo := repositories.NewCharacterRepository()
	inventoryRepo := repositories.NewInventoryRepository()
	handler := NewInventoryHandler(characterRepo, inventoryRepo)

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	characterRepo.Save(character)

	// Create test items
	sword := models.NewWeapon("Test Sword", 10, 100, 1, nil)
	armor := models.NewArmor("Test Armor", 5, 80, 1, nil)
	potion := models.NewPotion("Health Potion", 20, 30)

	// Add items to inventory
	character.AddToInventory(sword)
	character.AddToInventory(armor)
	character.AddToInventory(potion)

	// Save items to repository
	inventoryRepo.SaveItem(sword)
	inventoryRepo.SaveItem(armor)
	inventoryRepo.SaveItem(potion)

	// Save the updated character
	characterRepo.Save(character)

	// Test GetInventory
	t.Run("GetInventory", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/characters/"+character.ID+"/inventory", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var inventory []*models.Item
		json.Unmarshal(rr.Body.Bytes(), &inventory)
		assert.Equal(t, 3, len(inventory))
	})

	// Test GetInventoryItem
	t.Run("GetInventoryItem", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/characters/"+character.ID+"/inventory/"+sword.ID, nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var item models.Item
		json.Unmarshal(rr.Body.Bytes(), &item)
		assert.Equal(t, sword.ID, item.ID)
		assert.Equal(t, "Test Sword", item.Name)
	})

	// Test EquipItem
	t.Run("EquipItem", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/characters/"+character.ID+"/inventory/"+sword.ID+"/equip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify the item is equipped
		updatedCharacter, _ := characterRepo.GetByID(character.ID)
		assert.NotNil(t, updatedCharacter.Equipment.Weapon)
		assert.Equal(t, sword.ID, updatedCharacter.Equipment.Weapon.ID)
		assert.True(t, updatedCharacter.Equipment.Weapon.Equipped)
	})

	// Test that equipping a second weapon replaces the first one
	t.Run("EquipSecondWeapon", func(t *testing.T) {
		// Create a second weapon
		sword2 := models.NewWeapon("Second Sword", 15, 150, 1, nil)
		character.AddToInventory(sword2)
		inventoryRepo.SaveItem(sword2)
		characterRepo.Save(character)

		// First, make sure the first sword is equipped
		character.EquipItem(sword.ID)
		characterRepo.Save(character)

		// Now equip the second sword
		req, _ := http.NewRequest("POST", "/api/characters/"+character.ID+"/inventory/"+sword2.ID+"/equip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify the second sword replaced the first one
		updatedCharacter, _ := characterRepo.GetByID(character.ID)
		assert.NotNil(t, updatedCharacter.Equipment.Weapon)
		assert.Equal(t, sword2.ID, updatedCharacter.Equipment.Weapon.ID)
		assert.True(t, updatedCharacter.Equipment.Weapon.Equipped)

		// Verify the first sword is no longer equipped
		firstSword, found := updatedCharacter.GetInventoryItem(sword.ID)
		assert.True(t, found)
		assert.False(t, firstSword.Equipped)
	})

	// Test that equipping a second armor replaces the first one
	t.Run("EquipSecondArmor", func(t *testing.T) {
		// Create a second armor
		armor2 := models.NewArmor("Second Armor", 8, 120, 1, nil)
		character.AddToInventory(armor2)
		inventoryRepo.SaveItem(armor2)
		characterRepo.Save(character)

		// First, make sure the first armor is equipped
		character.EquipItem(armor.ID)
		characterRepo.Save(character)

		// Now equip the second armor
		req, _ := http.NewRequest("POST", "/api/characters/"+character.ID+"/inventory/"+armor2.ID+"/equip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify the second armor replaced the first one
		updatedCharacter, _ := characterRepo.GetByID(character.ID)
		assert.NotNil(t, updatedCharacter.Equipment.Armor)
		assert.Equal(t, armor2.ID, updatedCharacter.Equipment.Armor.ID)
		assert.True(t, updatedCharacter.Equipment.Armor.Equipped)

		// Verify the first armor is no longer equipped
		firstArmor, found := updatedCharacter.GetInventoryItem(armor.ID)
		assert.True(t, found)
		assert.False(t, firstArmor.Equipped)
	})

	// Test UnequipItem
	t.Run("UnequipItem", func(t *testing.T) {
		// First, make sure the item is equipped
		character.EquipItem(sword.ID)
		characterRepo.Save(character)

		req, _ := http.NewRequest("POST", "/api/characters/"+character.ID+"/inventory/"+sword.ID+"/unequip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify the item is unequipped
		updatedCharacter, _ := characterRepo.GetByID(character.ID)
		assert.Nil(t, updatedCharacter.Equipment.Weapon)
	})

	// Test UnequipItem for non-existent character
	t.Run("UnequipItem_CharacterNotFound", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/characters/nonexistent-id/inventory/"+sword.ID+"/unequip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Character not found")
	})

	// Test UnequipItem for non-existent item
	t.Run("UnequipItem_ItemNotFound", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/characters/"+character.ID+"/inventory/nonexistent-id/unequip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Item not found")
	})

	// Test UnequipItem for armor
	t.Run("UnequipItem_Armor", func(t *testing.T) {
		// First, make sure the armor is equipped
		character.EquipItem(armor.ID)
		characterRepo.Save(character)

		req, _ := http.NewRequest("POST", "/api/characters/"+character.ID+"/inventory/"+armor.ID+"/unequip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify the armor is unequipped
		updatedCharacter, _ := characterRepo.GetByID(character.ID)
		assert.Nil(t, updatedCharacter.Equipment.Armor)
	})

	// Test UnequipItem for accessory
	t.Run("UnequipItem_Accessory", func(t *testing.T) {
		// Create and equip an accessory (artifact)
		accessory := &models.Item{
			ID:          uuid.New().String(),
			Type:        models.ItemArtifact,
			Name:        "Test Accessory",
			Description: "A test accessory",
			Value:       50,
			Power:       5,
			Weight:      0.5,
		}

		character.AddToInventory(accessory)
		inventoryRepo.SaveItem(accessory)

		// Directly set the accessory in the equipment slot
		character.Equipment.Accessory = accessory
		accessory.Equipped = true
		characterRepo.Save(character)

		req, _ := http.NewRequest("POST", "/api/characters/"+character.ID+"/inventory/"+accessory.ID+"/unequip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify the accessory is unequipped
		updatedCharacter, _ := characterRepo.GetByID(character.ID)
		assert.Nil(t, updatedCharacter.Equipment.Accessory)
	})

	// Test UnequipItem for item that's equipped but not in inventory
	t.Run("UnequipItem_EquippedButNotInInventory", func(t *testing.T) {
		// Create a new weapon and equip it directly to the equipment slot
		directWeapon := models.NewWeapon("Direct Weapon", 12, 120, 1, nil)
		inventoryRepo.SaveItem(directWeapon)

		// Manually set the equipment without adding to inventory
		character.Equipment.Weapon = directWeapon
		characterRepo.Save(character)

		req, _ := http.NewRequest("POST", "/api/characters/"+character.ID+"/inventory/"+directWeapon.ID+"/unequip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify the weapon is unequipped
		updatedCharacter, _ := characterRepo.GetByID(character.ID)
		assert.Nil(t, updatedCharacter.Equipment.Weapon)
	})

	// Test UnequipItem failure case for item not found
	t.Run("UnequipItem_FailureCase", func(t *testing.T) {
		// Create a mock character with a special setup that will cause unequip to fail
		mockCharacter := models.NewCharacter("MockCharacter", models.Warrior)
		characterRepo.Save(mockCharacter)

		// Try to unequip an item type that isn't equipped
		// This will cause the unequip operation to fail
		req, _ := http.NewRequest("POST", "/api/characters/"+mockCharacter.ID+"/inventory/nonexistent-id/unequip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		// Since there's no item to unequip, this should fail
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Item not found")
	})

	// Test UnequipItem failure case when character model's UnequipItem returns false
	t.Run("UnequipItem_ModelFailure", func(t *testing.T) {
		// Create a test character
		mockCharacter := models.NewCharacter("MockForFailure", models.Warrior)
		characterRepo.Save(mockCharacter)

		// Create a test item that's in inventory but not equipped
		testItem := models.NewWeapon("Test Weapon", 10, 100, 1, nil)
		mockCharacter.AddToInventory(testItem)
		inventoryRepo.SaveItem(testItem)
		characterRepo.Save(mockCharacter)

		// When we try to unequip an item that's in inventory but not equipped,
		// the character model's UnequipItem should return false
		req, _ := http.NewRequest("POST", "/api/characters/"+mockCharacter.ID+"/inventory/"+testItem.ID+"/unequip", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		// This should fail with a bad request
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Failed to unequip item")
	})

	// Test UseItem
	t.Run("UseItem", func(t *testing.T) {
		// Reduce HP to test healing
		character.CurrentHP = character.MaxHP - 10
		characterRepo.Save(character)

		req, _ := http.NewRequest("POST", "/api/characters/"+character.ID+"/inventory/"+potion.ID+"/use", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify the potion was used and HP increased
		updatedCharacter, _ := characterRepo.GetByID(character.ID)
		assert.Equal(t, character.MaxHP, updatedCharacter.CurrentHP) // HP should be restored

		// Verify the potion was removed from inventory
		_, found := updatedCharacter.GetInventoryItem(potion.ID)
		assert.False(t, found)
	})

	// Test GetEquipment
	t.Run("GetEquipment", func(t *testing.T) {
		// First, equip an item
		character.EquipItem(armor.ID)
		characterRepo.Save(character)

		req, _ := http.NewRequest("GET", "/api/characters/"+character.ID+"/equipment", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var equipment models.Equipment
		json.Unmarshal(rr.Body.Bytes(), &equipment)
		assert.NotNil(t, equipment.Armor)
		assert.Equal(t, armor.ID, equipment.Armor.ID)
	})

	// Test GenerateItems
	t.Run("GenerateItems", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"count": 5, "floorLevel": 3}`)
		req, _ := http.NewRequest("POST", "/api/items/generate", reqBody)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var items []*models.Item
		json.Unmarshal(rr.Body.Bytes(), &items)
		assert.Equal(t, 5, len(items))
	})
}
