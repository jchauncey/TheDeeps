package game

import (
	"testing"

	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/stretchr/testify/assert"
)

func TestHandlePickup(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	character.Attributes.Strength = 10 // Weight limit of 50
	characterRepo.Save(character)

	// Create a test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	floor := &models.Floor{
		Level:  1,
		Width:  10,
		Height: 10,
		Items:  make(map[string]models.Item),
	}
	dungeon.FloorData = make(map[int]*models.Floor)
	dungeon.FloorData[1] = floor
	dungeonRepo.Save(dungeon)

	// Set character's current dungeon and floor
	character.CurrentDungeon = dungeon.ID
	character.CurrentFloor = 1
	character.Position = models.Position{X: 5, Y: 5}
	characterRepo.Save(character)

	// Create a test client
	client := &Client{
		ID:        "test-client",
		Character: character,
		Manager:   manager,
		Send:      make(chan Message, 10),
	}

	// Test 1: Pickup a light item (should succeed)
	// Create a light item at the character's position
	lightItem := models.NewWeaponWithWeight("Light Sword", 5, 10, 2.0, 1, nil)
	lightItem.Position = character.Position
	floor.Items[lightItem.ID] = *lightItem
	dungeonRepo.Save(dungeon)

	// Send pickup message
	manager.handlePickup(client, Message{
		Type:   MsgPickup,
		ItemID: lightItem.ID,
	})

	// Check that the item was added to inventory
	updatedCharacter, _ := characterRepo.GetByID(character.ID)
	assert.Equal(t, 1, len(updatedCharacter.Inventory))
	assert.Equal(t, lightItem.ID, updatedCharacter.Inventory[0].ID)

	// Check that the item was removed from the floor
	updatedDungeon, _ := dungeonRepo.GetByID(dungeon.ID)
	_, exists := updatedDungeon.FloorData[1].Items[lightItem.ID]
	assert.False(t, exists)

	// Test 2: Pickup a heavy item that would exceed weight limit (should fail)
	// Create a very heavy item at the character's position
	heavyItem := models.NewArmorWithWeight("Heavy Armor", 10, 100, 49.0, 1, nil) // Just under the limit
	heavyItem.Position = character.Position
	updatedDungeon.FloorData[1].Items[heavyItem.ID] = *heavyItem
	dungeonRepo.Save(updatedDungeon)

	// Send pickup message
	manager.handlePickup(client, Message{
		Type:   MsgPickup,
		ItemID: heavyItem.ID,
	})

	// Check that the item was not added to inventory (still only 1 item)
	updatedCharacter, _ = characterRepo.GetByID(character.ID)
	assert.Equal(t, 1, len(updatedCharacter.Inventory))

	// Check that the item is still on the floor
	updatedDungeon, _ = dungeonRepo.GetByID(dungeon.ID)
	_, exists = updatedDungeon.FloorData[1].Items[heavyItem.ID]
	assert.True(t, exists)

	// Drain the message channel to avoid blocking
	for i := 0; i < len(client.Send); i++ {
		<-client.Send
	}
}
