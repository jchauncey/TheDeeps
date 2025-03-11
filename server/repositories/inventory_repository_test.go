package repositories

import (
	"testing"

	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/stretchr/testify/assert"
)

func TestInventoryRepository(t *testing.T) {
	repo := NewInventoryRepository()

	// Test saving and retrieving an item
	t.Run("SaveAndGetItem", func(t *testing.T) {
		item := models.NewWeapon("Test Sword", 10, 100, 1, nil)
		err := repo.SaveItem(item)
		assert.NoError(t, err)

		retrievedItem, exists := repo.GetItem(item.ID)
		assert.True(t, exists)
		assert.Equal(t, item.ID, retrievedItem.ID)
		assert.Equal(t, "Test Sword", retrievedItem.Name)
	})

	// Test deleting an item
	t.Run("DeleteItem", func(t *testing.T) {
		item := models.NewWeapon("Delete Sword", 5, 50, 1, nil)
		err := repo.SaveItem(item)
		assert.NoError(t, err)

		success := repo.DeleteItem(item.ID)
		assert.True(t, success)

		_, exists := repo.GetItem(item.ID)
		assert.False(t, exists)

		// Try deleting non-existent item
		success = repo.DeleteItem("non-existent-id")
		assert.False(t, success)
	})

	// Test getting all items
	t.Run("GetAllItems", func(t *testing.T) {
		// Clear existing items
		repo = NewInventoryRepository()

		item1 := models.NewWeapon("Sword 1", 10, 100, 1, nil)
		item2 := models.NewArmor("Armor 1", 5, 80, 1, nil)
		item3 := models.NewPotion("Health Potion", 20, 30)

		repo.SaveItem(item1)
		repo.SaveItem(item2)
		repo.SaveItem(item3)

		allItems := repo.GetAllItems()
		assert.Equal(t, 3, len(allItems))

		// Check if all items are in the result
		itemIDs := make(map[string]bool)
		for _, item := range allItems {
			itemIDs[item.ID] = true
		}

		assert.True(t, itemIDs[item1.ID])
		assert.True(t, itemIDs[item2.ID])
		assert.True(t, itemIDs[item3.ID])
	})

	// Test generating random items
	t.Run("GenerateRandomItems", func(t *testing.T) {
		items := repo.GenerateRandomItems(5, 3)
		assert.Equal(t, 5, len(items))

		// Check if all items are saved in the repository
		for _, item := range items {
			retrievedItem, exists := repo.GetItem(item.ID)
			assert.True(t, exists)
			assert.Equal(t, item.ID, retrievedItem.ID)
		}
	})
}
