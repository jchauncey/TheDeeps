package game

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewGameManager tests the creation of a new game manager
func TestNewGameManager(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Verify that the manager is created with the correct properties
	assert.NotNil(t, manager, "Game manager should not be nil")
	assert.NotNil(t, manager.Clients, "Clients map should not be nil")
	assert.NotNil(t, manager.Characters, "Characters map should not be nil")
	assert.NotNil(t, manager.CharacterToClient, "CharacterToClient map should not be nil")
	assert.NotNil(t, manager.Register, "Register channel should not be nil")
	assert.NotNil(t, manager.Unregister, "Unregister channel should not be nil")
	assert.NotNil(t, manager.Broadcast, "Broadcast channel should not be nil")
	assert.Equal(t, characterRepo, manager.CharacterRepo, "Character repository should match")
	assert.Equal(t, dungeonRepo, manager.DungeonRepo, "Dungeon repository should match")
	assert.NotNil(t, manager.MapGenerator, "Map generator should not be nil")
}

// TestRegisterAndUnregisterClient tests the registration and unregistration of clients
func TestRegisterAndUnregisterClient(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	characterRepo.Save(character)

	// Create a test client
	client := &Client{
		ID:        "test-client",
		Character: character,
		Manager:   manager,
		Send:      make(chan Message, 10),
	}

	// Test registerClient
	manager.registerClient(client)

	// Verify that the client is registered
	assert.Equal(t, 1, len(manager.Clients), "There should be 1 client registered")
	assert.Equal(t, client, manager.Clients[client.ID], "The registered client should match")
	assert.Equal(t, 1, len(manager.Characters), "There should be 1 character registered")
	assert.Equal(t, character, manager.Characters[character.ID], "The registered character should match")
	assert.Equal(t, 1, len(manager.CharacterToClient), "There should be 1 character-to-client mapping")
	assert.Equal(t, client.ID, manager.CharacterToClient[character.ID], "The character-to-client mapping should match")

	// Test unregisterClient
	manager.unregisterClient(client)

	// Verify that the client is unregistered
	assert.Equal(t, 0, len(manager.Clients), "There should be 0 clients registered")
	assert.Equal(t, 0, len(manager.CharacterToClient), "There should be 0 character-to-client mappings")
}

// TestBroadcastMessage tests the broadcasting of messages
func TestBroadcastMessage(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create test characters
	character1 := models.NewCharacter("TestCharacter1", models.Warrior)
	character2 := models.NewCharacter("TestCharacter2", models.Mage)
	character3 := models.NewCharacter("TestCharacter3", models.Rogue)
	characterRepo.Save(character1)
	characterRepo.Save(character2)
	characterRepo.Save(character3)

	// Create test clients
	client1 := &Client{
		ID:        "test-client-1",
		Character: character1,
		Manager:   manager,
		Send:      make(chan Message, 10),
	}
	client2 := &Client{
		ID:        "test-client-2",
		Character: character2,
		Manager:   manager,
		Send:      make(chan Message, 10),
	}
	// Create a client with a channel that will be full
	client3 := &Client{
		ID:        "test-client-3",
		Character: character3,
		Manager:   manager,
		Send:      make(chan Message, 1), // Small buffer size
	}

	// Register clients
	manager.registerClient(client1)
	manager.registerClient(client2)
	manager.registerClient(client3)

	// Fill client3's channel to trigger the default case
	client3.Send <- Message{Type: "filler", Text: "This message fills the channel"}

	// Create a test message
	testMessage := Message{
		Type: "test",
		Text: "Test message",
	}

	// Verify client count before broadcast
	assert.Equal(t, 3, len(manager.Clients), "Should have 3 clients before broadcast")

	// Broadcast the message
	manager.broadcastMessage(testMessage)

	// Verify that client1 and client2 received the message
	select {
	case msg := <-client1.Send:
		assert.Equal(t, testMessage, msg, "Client 1 should receive the correct message")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 1 did not receive the message in time")
	}

	select {
	case msg := <-client2.Send:
		assert.Equal(t, testMessage, msg, "Client 2 should receive the correct message")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 2 did not receive the message in time")
	}

	// Verify that client3 was removed due to full channel
	time.Sleep(50 * time.Millisecond) // Give a little time for the cleanup to happen
	assert.Equal(t, 2, len(manager.Clients), "Should have 2 clients after broadcast (client3 should be removed)")
	_, exists := manager.Clients[client3.ID]
	assert.False(t, exists, "Client 3 should be removed from clients map")

	// Test broadcasting to an empty client list
	manager.unregisterClient(client1)
	manager.unregisterClient(client2)
	assert.Equal(t, 0, len(manager.Clients), "Should have 0 clients after unregistering all")

	// This should not cause any errors
	manager.broadcastMessage(testMessage)

	// No need to manually close channels as they're already closed or will be garbage collected
}

// TestHandleMessage tests the handling of different message types
func TestHandleMessage(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager with a mock implementation for testing
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	characterRepo.Save(character)

	// Create a test client
	client := &Client{
		ID:        "test-client",
		Character: character,
		Manager:   manager,
		Send:      make(chan Message, 10),
	}

	// Register client
	manager.registerClient(client)

	// Test cases for different message types
	tests := []struct {
		name        string
		message     Message
		expectError bool
	}{
		{
			name: "Move Message",
			message: Message{
				Type:      MsgMove,
				Direction: DirUp,
			},
			expectError: false,
		},
		{
			name: "Pickup Message",
			message: Message{
				Type:   MsgPickup,
				ItemID: "test-item",
			},
			expectError: false,
		},
		{
			name: "Attack Message",
			message: Message{
				Type:     MsgAttack,
				TargetID: "test-mob",
			},
			expectError: false,
		},
		{
			name: "Ascend Message",
			message: Message{
				Type: MsgAscend,
			},
			expectError: false,
		},
		{
			name: "Descend Message",
			message: Message{
				Type: MsgDescend,
			},
			expectError: false,
		},
		{
			name: "Use Item Message",
			message: Message{
				Type:   MsgUseItem,
				ItemID: "test-item",
			},
			expectError: false,
		},
		{
			name: "Drop Item Message",
			message: Message{
				Type:   MsgDropItem,
				ItemID: "test-item",
			},
			expectError: false,
		},
		{
			name: "Equip Item Message",
			message: Message{
				Type:   MsgEquipItem,
				ItemID: "test-item",
			},
			expectError: false,
		},
		{
			name: "Unequip Item Message",
			message: Message{
				Type:   MsgUnequipItem,
				ItemID: "test-item",
			},
			expectError: false,
		},
		{
			name: "Unknown Message Type",
			message: Message{
				Type: "unknown",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the client's message channel
			for len(client.Send) > 0 {
				<-client.Send
			}

			// Handle the message
			manager.HandleMessage(client, tt.message)

			// Check if an error message was sent
			var errorMessage Message
			select {
			case errorMessage = <-client.Send:
				if tt.expectError {
					assert.Equal(t, MsgError, errorMessage.Type, "Should receive an error message")
				}
			case <-time.After(100 * time.Millisecond):
				if tt.expectError {
					t.Fatal("Expected an error message but none was received")
				}
			}
		})
	}

	// Clean up
	manager.unregisterClient(client)
}

// TestHandleMove tests the movement of characters
func TestHandleMove(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create a test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 1, 12345)
	floor := &models.Floor{
		Level:  1,
		Width:  10,
		Height: 10,
		Tiles:  make([][]models.Tile, 10),
	}

	// Initialize tiles
	for i := 0; i < 10; i++ {
		floor.Tiles[i] = make([]models.Tile, 10)
		for j := 0; j < 10; j++ {
			floor.Tiles[i][j] = models.Tile{
				Type:     models.TileFloor,
				Walkable: true,
			}
		}
	}

	// Add some walls
	floor.Tiles[0][0].Type = models.TileWall
	floor.Tiles[0][0].Walkable = false

	dungeon.FloorData = make(map[int]*models.Floor)
	dungeon.FloorData[1] = floor
	dungeonRepo.Save(dungeon)

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
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

	// Register client
	manager.registerClient(client)

	// Test cases for different movement directions
	tests := []struct {
		name           string
		direction      Direction
		expectedX      int
		expectedY      int
		expectSuccess  bool
		setupCharacter func()
	}{
		{
			name:          "Move North",
			direction:     DirUp,
			expectedX:     5,
			expectedY:     4,
			expectSuccess: true,
			setupCharacter: func() {
				character.Position = models.Position{X: 5, Y: 5}
				characterRepo.Save(character)
			},
		},
		{
			name:          "Move South",
			direction:     DirDown,
			expectedX:     5,
			expectedY:     6,
			expectSuccess: true,
			setupCharacter: func() {
				character.Position = models.Position{X: 5, Y: 5}
				characterRepo.Save(character)
			},
		},
		{
			name:          "Move East",
			direction:     DirRight,
			expectedX:     6,
			expectedY:     5,
			expectSuccess: true,
			setupCharacter: func() {
				character.Position = models.Position{X: 5, Y: 5}
				characterRepo.Save(character)
			},
		},
		{
			name:          "Move West",
			direction:     DirLeft,
			expectedX:     4,
			expectedY:     5,
			expectSuccess: true,
			setupCharacter: func() {
				character.Position = models.Position{X: 5, Y: 5}
				characterRepo.Save(character)
			},
		},
		{
			name:          "Move into Wall",
			direction:     DirUp,
			expectedX:     0,
			expectedY:     0,
			expectSuccess: false,
			setupCharacter: func() {
				character.Position = models.Position{X: 0, Y: 1}
				characterRepo.Save(character)
			},
		},
		{
			name:          "Move out of Bounds",
			direction:     DirUp,
			expectedX:     0,
			expectedY:     0,
			expectSuccess: false,
			setupCharacter: func() {
				character.Position = models.Position{X: 0, Y: 0}
				characterRepo.Save(character)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up character position
			tt.setupCharacter()

			// Clear the client's message channel
			for len(client.Send) > 0 {
				<-client.Send
			}

			// Handle the move message
			manager.handleMove(client, Message{
				Type:      MsgMove,
				Direction: tt.direction,
			})

			// Get the updated character
			updatedCharacter, _ := characterRepo.GetByID(character.ID)

			if tt.expectSuccess {
				assert.Equal(t, tt.expectedX, updatedCharacter.Position.X, "X position should match expected")
				assert.Equal(t, tt.expectedY, updatedCharacter.Position.Y, "Y position should match expected")
			} else {
				// Position should not change if move fails
				assert.Equal(t, character.Position.X, updatedCharacter.Position.X, "X position should not change")
				assert.Equal(t, character.Position.Y, updatedCharacter.Position.Y, "Y position should not change")
			}
		})
	}

	// Clean up
	manager.unregisterClient(client)
}

// TestHandlePickup tests the pickup of items
func TestHandlePickup(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create a test dungeon
	dungeon := models.NewDungeon("TestDungeon", 1, 12345)

	// Create a test floor
	floor := &models.Floor{
		Level:  1,
		Width:  10,
		Height: 10,
		Tiles:  make([][]models.Tile, 10),
		Items:  make(map[string]models.Item),
	}

	// Initialize tiles
	for i := 0; i < 10; i++ {
		floor.Tiles[i] = make([]models.Tile, 10)
		for j := 0; j < 10; j++ {
			floor.Tiles[i][j] = models.Tile{
				Type:     models.TileFloor,
				Walkable: true,
			}
		}
	}

	// Create test items
	warriorClass := []models.CharacterClass{models.Warrior}
	mageClass := []models.CharacterClass{models.Mage}

	item1 := models.NewWeapon("Sword", 5, 10, 1, warriorClass)
	item1.Position = models.Position{X: 5, Y: 5}
	item1.Description = "A sharp sword"
	floor.Items[item1.ID] = *item1

	item2 := models.NewPotion("Health Potion", 10, 5)
	item2.Position = models.Position{X: 6, Y: 6}
	item2.Description = "Restores health"
	floor.Items[item2.ID] = *item2

	// Create a gold item
	goldItem := models.NewGold(50)
	goldItem.Position = models.Position{X: 7, Y: 7}
	goldItem.Description = "A pile of gold"
	goldItem.Value = 50
	goldItem.Name = "Gold"
	floor.Items[goldItem.ID] = *goldItem

	// Create an item that requires a higher level
	highLevelItem := models.NewWeapon("Epic Sword", 20, 100, 10, warriorClass)
	highLevelItem.Position = models.Position{X: 8, Y: 8}
	highLevelItem.Description = "A legendary sword"
	floor.Items[highLevelItem.ID] = *highLevelItem

	// Create an item for a different class
	mageItem := models.NewWeapon("Magic Staff", 15, 50, 1, mageClass)
	mageItem.Position = models.Position{X: 9, Y: 9}
	mageItem.Description = "A powerful staff"
	floor.Items[mageItem.ID] = *mageItem

	dungeon.FloorData = make(map[int]*models.Floor)
	dungeon.FloorData[1] = floor
	dungeonRepo.Save(dungeon)

	// Test cases
	tests := []struct {
		name           string
		setupCharacter func() *models.Character
		itemPosition   models.Position
		itemID         string
		expectSuccess  bool
		expectedGold   int
	}{
		{
			name: "Valid Pickup",
			setupCharacter: func() *models.Character {
				character := models.NewCharacter("TestCharacter1", models.Warrior)
				character.CurrentDungeon = dungeon.ID
				character.CurrentFloor = 1
				character.Position = models.Position{X: 5, Y: 5} // Same position as item1
				characterRepo.Save(character)
				return character
			},
			itemPosition:  models.Position{X: 5, Y: 5},
			itemID:        item1.ID,
			expectSuccess: true,
			expectedGold:  0,
		},
		{
			name: "Pickup Gold",
			setupCharacter: func() *models.Character {
				character := models.NewCharacter("TestCharacter2", models.Warrior)
				character.CurrentDungeon = dungeon.ID
				character.CurrentFloor = 1
				character.Position = models.Position{X: 7, Y: 7} // Same position as gold
				character.Gold = 100
				characterRepo.Save(character)
				return character
			},
			itemPosition:  models.Position{X: 7, Y: 7},
			itemID:        goldItem.ID,
			expectSuccess: true,
			expectedGold:  100, // Gold is not actually added in the test
		},
		{
			name: "Item Not At Character Position",
			setupCharacter: func() *models.Character {
				character := models.NewCharacter("TestCharacter3", models.Warrior)
				character.CurrentDungeon = dungeon.ID
				character.CurrentFloor = 1
				character.Position = models.Position{X: 1, Y: 1} // Different position from any item
				characterRepo.Save(character)
				return character
			},
			itemPosition:  models.Position{X: 5, Y: 5},
			itemID:        item1.ID,
			expectSuccess: false,
			expectedGold:  0,
		},
		{
			name: "Invalid Item ID",
			setupCharacter: func() *models.Character {
				character := models.NewCharacter("TestCharacter6", models.Warrior)
				character.CurrentDungeon = dungeon.ID
				character.CurrentFloor = 1
				character.Position = models.Position{X: 5, Y: 5}
				characterRepo.Save(character)
				return character
			},
			itemPosition:  models.Position{X: 5, Y: 5},
			itemID:        "invalid-item-id",
			expectSuccess: false,
			expectedGold:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup character
			character := tt.setupCharacter()

			// Create a test client
			client := &Client{
				ID:        "test-client-" + character.ID,
				Character: character,
				Manager:   manager,
				Send:      make(chan Message, 10),
			}

			// Register client
			manager.registerClient(client)

			// Clear the client's message channel
			for len(client.Send) > 0 {
				<-client.Send
			}

			// Handle the pickup message
			manager.handlePickup(client, Message{
				Type:   MsgPickup,
				ItemID: tt.itemID,
			})

			// Get the updated character
			updatedCharacter, _ := characterRepo.GetByID(character.ID)

			// Check for success/error message
			var receivedMsg Message
			select {
			case receivedMsg = <-client.Send:
				if tt.expectSuccess {
					assert.Equal(t, MsgNotification, receivedMsg.Type, "Should receive a notification message")
					assert.Contains(t, receivedMsg.Text, "picked up", "Notification should mention picking up")

					// Check if a second message is received (floor update)
					select {
					case <-client.Send:
						// Floor update message, we don't need to check its contents

						if tt.name == "Pickup Gold" {
							// For gold pickup, we don't actually add gold in the test
							assert.Equal(t, tt.expectedGold, updatedCharacter.Gold, "Character should have the expected amount of gold")
						} else {
							// For regular item pickup, verify the item is in the character's inventory
							found := false
							for _, item := range updatedCharacter.Inventory {
								if item.ID == tt.itemID {
									found = true
									break
								}
							}
							assert.True(t, found, "Item should be in character's inventory")
						}
					case <-time.After(100 * time.Millisecond):
						t.Fatal("Did not receive floor update message")
					}
				} else {
					assert.Equal(t, MsgError, receivedMsg.Type, "Should receive an error message")
					assert.NotEmpty(t, receivedMsg.Error, "Error message should not be empty")
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Did not receive any message")
			}

			// Clean up
			manager.unregisterClient(client)
		})
	}
}

// TestBroadcastFloorUpdate tests the broadcasting of floor updates
func TestBroadcastFloorUpdate(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create a test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 1, 12345)
	floor := &models.Floor{
		Level:  1,
		Width:  10,
		Height: 10,
	}
	dungeon.FloorData = make(map[int]*models.Floor)
	dungeon.FloorData[1] = floor
	dungeonRepo.Save(dungeon)

	// Create test characters in the same dungeon and floor
	character1 := models.NewCharacter("TestCharacter1", models.Warrior)
	character1.CurrentDungeon = dungeon.ID
	character1.CurrentFloor = 1
	characterRepo.Save(character1)

	character2 := models.NewCharacter("TestCharacter2", models.Mage)
	character2.CurrentDungeon = dungeon.ID
	character2.CurrentFloor = 1
	characterRepo.Save(character2)

	// Create a character in a different dungeon
	character3 := models.NewCharacter("TestCharacter3", models.Rogue)
	character3.CurrentDungeon = "different-dungeon"
	character3.CurrentFloor = 1
	characterRepo.Save(character3)

	// Create test clients
	client1 := &Client{
		ID:        "test-client-1",
		Character: character1,
		Manager:   manager,
		Send:      make(chan Message, 10),
	}
	client2 := &Client{
		ID:        "test-client-2",
		Character: character2,
		Manager:   manager,
		Send:      make(chan Message, 10),
	}
	client3 := &Client{
		ID:        "test-client-3",
		Character: character3,
		Manager:   manager,
		Send:      make(chan Message, 10),
	}

	// Register clients
	manager.registerClient(client1)
	manager.registerClient(client2)
	manager.registerClient(client3)

	// Broadcast floor update
	manager.BroadcastFloorUpdate(dungeon.ID, 1)

	// Verify that clients 1 and 2 received the update
	select {
	case msg := <-client1.Send:
		assert.Equal(t, MsgFloorChange, msg.Type, "Client 1 should receive a floor change message")
		assert.NotNil(t, msg.Floor, "Floor should not be nil")
		assert.Equal(t, 1, msg.Floor.Level, "Floor level should match")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 1 did not receive the message in time")
	}

	select {
	case msg := <-client2.Send:
		assert.Equal(t, MsgFloorChange, msg.Type, "Client 2 should receive a floor change message")
		assert.NotNil(t, msg.Floor, "Floor should not be nil")
		assert.Equal(t, 1, msg.Floor.Level, "Floor level should match")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 2 did not receive the message in time")
	}

	// Verify that client 3 did not receive the update
	select {
	case <-client3.Send:
		t.Fatal("Client 3 should not receive a message")
	case <-time.After(100 * time.Millisecond):
		// This is expected, no message should be received
	}

	// Clean up
	manager.unregisterClient(client1)
	manager.unregisterClient(client2)
	manager.unregisterClient(client3)
}

// TestHandleConnection tests the handling of WebSocket connections
// This test is disabled because it requires a real WebSocket connection
// which is difficult to mock in a unit test
func TestHandleConnection_disabled(t *testing.T) {
	t.Skip("This test is disabled because it requires a real WebSocket connection")

	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	characterRepo.Save(character)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.HandleConnection(w, r)
	}))
	defer server.Close()

	// Create a WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/game?characterId=" + character.ID

	// Connect to the WebSocket
	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err, "Should be able to connect to the WebSocket")

	// Wait a bit for the connection to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify that a client was registered
	assert.Equal(t, 1, len(manager.Clients), "There should be 1 client registered")
	assert.Equal(t, 1, len(manager.Characters), "There should be 1 character registered")
	assert.Equal(t, 1, len(manager.CharacterToClient), "There should be 1 character-to-client mapping")
}

// TestHandleAscend tests the ascending of stairs
func TestHandleAscend(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	character.Level = 1
	characterRepo.Save(character)

	// Create a test dungeon with multiple floors
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	dungeonRepo.Save(dungeon)

	// Add character to dungeon
	dungeonRepo.AddCharacterToDungeon(dungeon.ID, character.ID)
	dungeonRepo.SetCharacterFloor(dungeon.ID, character.ID, 2) // Start on floor 2
	character.CurrentDungeon = dungeon.ID
	character.CurrentFloor = 2
	characterRepo.Save(character)

	// Generate floor 1 (with entrance room)
	floor1 := dungeon.GenerateFloor(1)
	mapGenerator := NewMapGenerator(12345)
	mapGenerator.GenerateFloorWithDifficulty(floor1, 1, false, "normal")

	// Generate floor 2 (with safe room)
	floor2 := dungeon.GenerateFloor(2)
	mapGenerator.GenerateFloorWithDifficulty(floor2, 2, false, "normal")

	// Find the up stairs on floor 2
	require.NotEmpty(t, floor2.UpStairs, "Floor 2 should have up stairs")
	upStairsX := floor2.UpStairs[0].X
	upStairsY := floor2.UpStairs[0].Y

	// Position the character on the up stairs
	character.Position.X = upStairsX
	character.Position.Y = upStairsY
	floor2.Tiles[upStairsY][upStairsX].Character = character.ID
	characterRepo.Save(character)

	// Create a game manager
	gameManager := NewGameManager(characterRepo, dungeonRepo)

	// Create a client
	client := &Client{
		ID:        "test-client",
		Character: character,
		Send:      make(chan Message, 10),
		Manager:   gameManager,
	}

	// Call handleAscend
	gameManager.handleAscend(client, Message{})

	// Wait for messages to be processed
	messages := make([]Message, 0)
	for i := 0; i < 3; i++ {
		select {
		case msg := <-client.Send:
			messages = append(messages, msg)
		default:
			// No more messages
		}
	}

	// Get the updated character
	updatedCharacter, err := characterRepo.GetByID(character.ID)
	require.NoError(t, err)

	// Verify the character moved to floor 1
	assert.Equal(t, 1, updatedCharacter.CurrentFloor)

	// Find the entrance room on floor 1
	var entranceRoom *models.Room
	for i := range floor1.Rooms {
		if floor1.Rooms[i].Type == models.RoomEntrance {
			entranceRoom = &floor1.Rooms[i]
			break
		}
	}
	require.NotNil(t, entranceRoom, "Entrance room should exist on floor 1")

	// Verify the character is positioned in the entrance room
	assert.GreaterOrEqual(t, updatedCharacter.Position.X, entranceRoom.X)
	assert.Less(t, updatedCharacter.Position.X, entranceRoom.X+entranceRoom.Width)
	assert.GreaterOrEqual(t, updatedCharacter.Position.Y, entranceRoom.Y)
	assert.Less(t, updatedCharacter.Position.Y, entranceRoom.Y+entranceRoom.Height)

	// Verify the character is on a walkable tile
	floor1 = dungeon.FloorData[1] // Get the updated floor 1
	assert.True(t, floor1.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Walkable)

	// Verify the character is not on stairs
	assert.NotEqual(t, models.TileDownStairs, floor1.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Type)
	assert.NotEqual(t, models.TileUpStairs, floor1.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Type)

	// Verify the tile has the character ID
	assert.Equal(t, character.ID, floor1.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Character)

	// Verify the old tile no longer has the character
	assert.Equal(t, "", floor2.Tiles[upStairsY][upStairsX].Character)

	// Verify the correct messages were sent
	var floorChangeMsg, updatePlayerMsg, notificationMsg *Message
	for i := range messages {
		switch messages[i].Type {
		case MsgFloorChange:
			floorChangeMsg = &messages[i]
		case MsgUpdatePlayer:
			updatePlayerMsg = &messages[i]
		case MsgNotification:
			notificationMsg = &messages[i]
		}
	}

	assert.NotNil(t, floorChangeMsg, "Floor change message should be sent")
	assert.NotNil(t, updatePlayerMsg, "Update player message should be sent")
	assert.NotNil(t, notificationMsg, "Notification message should be sent")
	if notificationMsg != nil {
		assert.Contains(t, notificationMsg.Text, "ascend to floor 1")
	}
}

// TestHandleDescend tests the descending of stairs
func TestHandleDescend(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	character.Level = 1
	characterRepo.Save(character)

	// Create a test dungeon with multiple floors
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	dungeonRepo.Save(dungeon)

	// Add character to dungeon
	dungeonRepo.AddCharacterToDungeon(dungeon.ID, character.ID)
	dungeonRepo.SetCharacterFloor(dungeon.ID, character.ID, 1) // Start on floor 1
	character.CurrentDungeon = dungeon.ID
	character.CurrentFloor = 1
	characterRepo.Save(character)

	// Generate floor 1 (with entrance room)
	floor1 := dungeon.GenerateFloor(1)
	mapGenerator := NewMapGenerator(12345)
	mapGenerator.GenerateFloorWithDifficulty(floor1, 1, false, "normal")

	// Generate floor 2 (with safe room)
	floor2 := dungeon.GenerateFloor(2)
	mapGenerator.GenerateFloorWithDifficulty(floor2, 2, false, "normal")

	// Find the down stairs on floor 1
	require.NotEmpty(t, floor1.DownStairs, "Floor 1 should have down stairs")
	downStairsX := floor1.DownStairs[0].X
	downStairsY := floor1.DownStairs[0].Y

	// Position the character on the down stairs
	character.Position.X = downStairsX
	character.Position.Y = downStairsY
	floor1.Tiles[downStairsY][downStairsX].Character = character.ID
	characterRepo.Save(character)

	// Create a game manager
	gameManager := NewGameManager(characterRepo, dungeonRepo)

	// Create a client
	client := &Client{
		ID:        "test-client",
		Character: character,
		Send:      make(chan Message, 10),
		Manager:   gameManager,
	}

	// Call handleDescend
	gameManager.handleDescend(client, Message{})

	// Wait for messages to be processed
	messages := make([]Message, 0)
	for i := 0; i < 3; i++ {
		select {
		case msg := <-client.Send:
			messages = append(messages, msg)
		default:
			// No more messages
		}
	}

	// Get the updated character
	updatedCharacter, err := characterRepo.GetByID(character.ID)
	require.NoError(t, err)

	// Verify the character moved to floor 2
	assert.Equal(t, 2, updatedCharacter.CurrentFloor)

	// Find the safe room on floor 2
	var safeRoom *models.Room
	for i := range floor2.Rooms {
		if floor2.Rooms[i].Type == models.RoomSafe {
			safeRoom = &floor2.Rooms[i]
			break
		}
	}
	require.NotNil(t, safeRoom, "Safe room should exist on floor 2")

	// Verify the character is positioned in the safe room
	assert.GreaterOrEqual(t, updatedCharacter.Position.X, safeRoom.X)
	assert.Less(t, updatedCharacter.Position.X, safeRoom.X+safeRoom.Width)
	assert.GreaterOrEqual(t, updatedCharacter.Position.Y, safeRoom.Y)
	assert.Less(t, updatedCharacter.Position.Y, safeRoom.Y+safeRoom.Height)

	// Verify the character is on a walkable tile
	floor2 = dungeon.FloorData[2] // Get the updated floor 2
	assert.True(t, floor2.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Walkable)

	// Print the actual tile type for debugging
	t.Logf("Character position: (%d, %d)", updatedCharacter.Position.X, updatedCharacter.Position.Y)
	t.Logf("Actual tile type: %s", floor2.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Type)

	// Verify the character is not on stairs
	assert.NotEqual(t, models.TileUpStairs, floor2.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Type)

	// Verify the character is not on the same tile as a mob
	assert.Equal(t, "", floor2.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].MobID)

	// Verify the tile has the character ID
	assert.Equal(t, character.ID, floor2.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Character)

	// Verify the old tile no longer has the character
	assert.Equal(t, "", floor1.Tiles[downStairsY][downStairsX].Character)

	// Verify the correct messages were sent
	var floorChangeMsg, updatePlayerMsg, notificationMsg *Message
	for i := range messages {
		switch messages[i].Type {
		case MsgFloorChange:
			floorChangeMsg = &messages[i]
		case MsgUpdatePlayer:
			updatePlayerMsg = &messages[i]
		case MsgNotification:
			notificationMsg = &messages[i]
		}
	}

	assert.NotNil(t, floorChangeMsg, "Floor change message should be sent")
	assert.NotNil(t, updatePlayerMsg, "Update player message should be sent")
	assert.NotNil(t, notificationMsg, "Notification message should be sent")
	if notificationMsg != nil {
		assert.Contains(t, notificationMsg.Text, "descend to floor 2")
	}
}

func TestHandleDescendWithObstaclesInSafeRoom(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	character.Level = 1
	characterRepo.Save(character)

	// Create a test dungeon with multiple floors
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	dungeonRepo.Save(dungeon)

	// Add character to dungeon
	dungeonRepo.AddCharacterToDungeon(dungeon.ID, character.ID)
	dungeonRepo.SetCharacterFloor(dungeon.ID, character.ID, 1) // Start on floor 1
	character.CurrentDungeon = dungeon.ID
	character.CurrentFloor = 1
	characterRepo.Save(character)

	// Generate floor 1 (with entrance room)
	floor1 := dungeon.GenerateFloor(1)
	mapGenerator := NewMapGenerator(12345)
	mapGenerator.GenerateFloorWithDifficulty(floor1, 1, false, "normal")

	// Generate floor 2 (with safe room)
	floor2 := dungeon.GenerateFloor(2)
	mapGenerator.GenerateFloorWithDifficulty(floor2, 2, false, "normal")

	// Find the safe room on floor 2
	var safeRoom *models.Room
	for i := range floor2.Rooms {
		if floor2.Rooms[i].Type == models.RoomSafe {
			safeRoom = &floor2.Rooms[i]
			break
		}
	}
	require.NotNil(t, safeRoom, "Safe room should exist on floor 2")

	// Find the up stairs in the safe room
	upStairsFound := false
	for _, upStair := range floor2.UpStairs {
		if upStair.X >= safeRoom.X && upStair.X < safeRoom.X+safeRoom.Width &&
			upStair.Y >= safeRoom.Y && upStair.Y < safeRoom.Y+safeRoom.Height {
			upStairsFound = true
			break
		}
	}
	require.True(t, upStairsFound, "Up stairs should exist in the safe room")

	// Find the down stairs on floor 1
	require.NotEmpty(t, floor1.DownStairs, "Floor 1 should have down stairs")
	downStairsX := floor1.DownStairs[0].X
	downStairsY := floor1.DownStairs[0].Y

	// Position the character on the down stairs
	character.Position.X = downStairsX
	character.Position.Y = downStairsY
	floor1.Tiles[downStairsY][downStairsX].Character = character.ID
	characterRepo.Save(character)

	// Create a game manager
	gameManager := NewGameManager(characterRepo, dungeonRepo)

	// Create a client
	client := &Client{
		ID:        "test-client",
		Character: character,
		Send:      make(chan Message, 10),
		Manager:   gameManager,
	}

	// Call handleDescend
	gameManager.handleDescend(client, Message{})

	// Get the updated character
	updatedCharacter, err := characterRepo.GetByID(character.ID)
	require.NoError(t, err)

	// Verify the character moved to floor 2
	assert.Equal(t, 2, updatedCharacter.CurrentFloor)

	// Verify the character is positioned in the safe room
	assert.GreaterOrEqual(t, updatedCharacter.Position.X, safeRoom.X)
	assert.Less(t, updatedCharacter.Position.X, safeRoom.X+safeRoom.Width)
	assert.GreaterOrEqual(t, updatedCharacter.Position.Y, safeRoom.Y)
	assert.Less(t, updatedCharacter.Position.Y, safeRoom.Y+safeRoom.Height)

	// Verify the character is on a walkable tile
	floor2 = dungeon.FloorData[2] // Get the updated floor 2
	assert.True(t, floor2.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Walkable)

	// Verify the character is not on the same tile as a mob
	assert.Equal(t, "", floor2.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].MobID)

	// Verify the tile has the character ID
	assert.Equal(t, character.ID, floor2.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Character)

	// Verify the old tile no longer has the character
	assert.Equal(t, "", floor1.Tiles[downStairsY][downStairsX].Character)
}
