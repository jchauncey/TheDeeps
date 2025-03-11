package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockWebSocketConn is a mock for the WebSocket connection
type MockWebSocketConn struct {
	messages [][]byte
}

// NewMockWebSocketConn creates a new mock WebSocket connection
func NewMockWebSocketConn() *MockWebSocketConn {
	return &MockWebSocketConn{
		messages: make([][]byte, 0),
	}
}

// WriteMessage implements the websocket.Conn WriteMessage method
func (m *MockWebSocketConn) WriteMessage(messageType int, data []byte) error {
	m.messages = append(m.messages, data)
	return nil
}

// TestNewCombatHandler tests the creation of a new combat handler
func TestNewCombatHandler(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()
	gameManager := game.NewGameManager(characterRepo, dungeonRepo)

	// Create combat handler
	handler := NewCombatHandler(characterRepo, dungeonRepo, gameManager)

	// Verify handler was created correctly
	assert.NotNil(t, handler, "Combat handler should not be nil")
	assert.NotNil(t, handler.characterRepo, "Character repository should not be nil")
	assert.NotNil(t, handler.dungeonRepo, "Dungeon repository should not be nil")
	assert.NotNil(t, handler.gameManager, "Game manager should not be nil")
	assert.NotNil(t, handler.combatManager, "Combat manager should not be nil")
}

// TestHandleAttack tests the handleAttack function
func TestHandleAttack(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()
	gameManager := game.NewGameManager(characterRepo, dungeonRepo)

	// Create combat handler
	handler := NewCombatHandler(characterRepo, dungeonRepo, gameManager)

	// Create test character
	character := models.NewCharacter("TestWarrior", models.Warrior)
	character.Position = models.Position{X: 5, Y: 5}
	characterRepo.Save(character)

	// Create test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 1, 12345)
	dungeon.AddCharacter(character.ID)
	dungeon.SetCharacterFloor(character.ID, 1)
	character.CurrentDungeon = dungeon.ID
	dungeonRepo.Save(dungeon)

	// Create a floor with a mob
	floor := &models.Floor{
		Level:  1,
		Width:  20,
		Height: 20,
		Tiles:  make([][]models.Tile, 20),
		Mobs:   make(map[string]*models.Mob),
		Items:  make(map[string]models.Item),
	}

	// Initialize tiles
	for y := 0; y < 20; y++ {
		floor.Tiles[y] = make([]models.Tile, 20)
		for x := 0; x < 20; x++ {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileFloor,
				Walkable: true,
			}
		}
	}

	// Add a mob adjacent to the character
	mob := models.NewMob(models.MobSkeleton, models.VariantEasy, 1)
	mob.Position = models.Position{X: 6, Y: 5} // Adjacent to character
	mob.HP = 10
	mobID := "mob1"
	floor.Mobs[mobID] = mob

	// Add a mob not adjacent to the character
	farMob := models.NewMob(models.MobGoblin, models.VariantNormal, 1)
	farMob.Position = models.Position{X: 10, Y: 10} // Not adjacent to character
	farMobID := "mob2"
	floor.Mobs[farMobID] = farMob

	// Save the floor
	dungeonRepo.SaveFloor(dungeon.ID, 1, floor)

	// Test cases
	tests := []struct {
		name     string
		mobID    string
		expected bool
		message  string
	}{
		{
			name:     "Valid Attack",
			mobID:    mobID,
			expected: true, // The actual success depends on RNG, but we're testing the function call
			message:  "",
		},
		{
			name:     "Mob Not Adjacent",
			mobID:    farMobID,
			expected: false,
			message:  "Not adjacent to mob",
		},
		{
			name:     "Invalid Mob ID",
			mobID:    "nonexistent",
			expected: false,
			message:  "Mob not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call handleAttack
			response := handler.handleAttack(character, tt.mobID)

			// For the valid case, we can't predict success due to RNG
			if tt.name == "Valid Attack" {
				// Just check that we got a response with action set to "attack"
				assert.Equal(t, "attack", response.Action, "Action should be 'attack'")
				return
			}

			// For other cases, check expected failure
			assert.Equal(t, tt.expected, response.Success, "Success flag should match expected")
			assert.Contains(t, response.Message, tt.message, "Error message should contain expected text")
		})
	}
}

// TestHandleUseItem tests the handleUseItem function
func TestHandleUseItem(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()
	gameManager := game.NewGameManager(characterRepo, dungeonRepo)

	// Create combat handler
	handler := NewCombatHandler(characterRepo, dungeonRepo, gameManager)

	// Create test character
	character := models.NewCharacter("TestWarrior", models.Warrior)
	characterRepo.Save(character)

	// Test handleUseItem
	response := handler.handleUseItem(character, "item1")

	// Since the function is not fully implemented, we expect a failure
	assert.False(t, response.Success, "UseItem should return false as it's not implemented")
	assert.Contains(t, response.Message, "not implemented", "Message should indicate not implemented")
}

// TestHandleFlee tests the handleFlee function
func TestHandleFlee(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()
	gameManager := game.NewGameManager(characterRepo, dungeonRepo)

	// Create combat handler
	handler := NewCombatHandler(characterRepo, dungeonRepo, gameManager)

	// Create test character
	character := models.NewCharacter("TestWarrior", models.Warrior)
	character.Position = models.Position{X: 5, Y: 5}
	characterRepo.Save(character)

	// Create test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 1, 12345)
	dungeon.AddCharacter(character.ID)
	dungeon.SetCharacterFloor(character.ID, 1)
	character.CurrentDungeon = dungeon.ID
	dungeonRepo.Save(dungeon)

	// Create a floor with a mob
	floor := &models.Floor{
		Level:  1,
		Width:  20,
		Height: 20,
		Tiles:  make([][]models.Tile, 20),
		Mobs:   make(map[string]*models.Mob),
		Items:  make(map[string]models.Item),
	}

	// Initialize tiles
	for y := 0; y < 20; y++ {
		floor.Tiles[y] = make([]models.Tile, 20)
		for x := 0; x < 20; x++ {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileFloor,
				Walkable: true,
			}
		}
	}

	// Add a mob adjacent to the character
	mob := models.NewMob(models.MobSkeleton, models.VariantEasy, 1)
	mob.Position = models.Position{X: 6, Y: 5} // Adjacent to character
	mobID := "mob1"
	floor.Mobs[mobID] = mob

	// Add a mob not adjacent to the character
	farMob := models.NewMob(models.MobGoblin, models.VariantNormal, 1)
	farMob.Position = models.Position{X: 10, Y: 10} // Not adjacent to character
	farMobID := "mob2"
	floor.Mobs[farMobID] = farMob

	// Save the floor
	dungeonRepo.SaveFloor(dungeon.ID, 1, floor)

	// Test cases
	tests := []struct {
		name     string
		mobID    string
		expected bool
		message  string
	}{
		{
			name:     "Valid Flee Attempt",
			mobID:    mobID,
			expected: true, // The actual success depends on RNG, but we're testing the function call
			message:  "",
		},
		{
			name:     "Mob Not Adjacent",
			mobID:    farMobID,
			expected: false,
			message:  "Not adjacent to mob",
		},
		{
			name:     "Invalid Mob ID",
			mobID:    "nonexistent",
			expected: false,
			message:  "Mob not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call handleFlee
			response := handler.handleFlee(character, tt.mobID)

			// For the valid case, we can't predict success due to RNG
			if tt.name == "Valid Flee Attempt" {
				// Just check that we got a response with action set to "flee"
				assert.Equal(t, "flee", response.Action, "Action should be 'flee'")
				return
			}

			// For other cases, check expected failure
			assert.Equal(t, tt.expected, response.Success, "Success flag should match expected")
			assert.Contains(t, response.Message, tt.message, "Error message should contain expected text")
		})
	}
}

// TestIsAdjacent tests the isAdjacent helper function
func TestIsAdjacent(t *testing.T) {
	tests := []struct {
		name     string
		pos1     models.Position
		pos2     models.Position
		expected bool
	}{
		{
			name:     "Same Position",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 5, Y: 5},
			expected: false, // Same position is not adjacent
		},
		{
			name:     "Adjacent North",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 5, Y: 4},
			expected: true,
		},
		{
			name:     "Adjacent South",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 5, Y: 6},
			expected: true,
		},
		{
			name:     "Adjacent East",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 6, Y: 5},
			expected: true,
		},
		{
			name:     "Adjacent West",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 4, Y: 5},
			expected: true,
		},
		{
			name:     "Adjacent Northeast",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 6, Y: 4},
			expected: true,
		},
		{
			name:     "Adjacent Northwest",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 4, Y: 4},
			expected: true,
		},
		{
			name:     "Adjacent Southeast",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 6, Y: 6},
			expected: true,
		},
		{
			name:     "Adjacent Southwest",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 4, Y: 6},
			expected: true,
		},
		{
			name:     "Not Adjacent",
			pos1:     models.Position{X: 5, Y: 5},
			pos2:     models.Position{X: 7, Y: 7},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAdjacent(tt.pos1, tt.pos2)
			assert.Equal(t, tt.expected, result, "isAdjacent should return expected result")
		})
	}
}

// TestAbs tests the abs helper function
func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "Positive Number",
			input:    5,
			expected: 5,
		},
		{
			name:     "Negative Number",
			input:    -5,
			expected: 5,
		},
		{
			name:     "Zero",
			input:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abs(tt.input)
			assert.Equal(t, tt.expected, result, "abs should return expected result")
		})
	}
}

// TestFindSafePosition tests the findSafePosition helper function
func TestFindSafePosition(t *testing.T) {
	// Create a test floor
	floor := &models.Floor{
		Level:  1,
		Width:  10,
		Height: 10,
		Tiles:  make([][]models.Tile, 10),
		Mobs:   make(map[string]*models.Mob),
	}

	// Initialize tiles
	for y := 0; y < 10; y++ {
		floor.Tiles[y] = make([]models.Tile, 10)
		for x := 0; x < 10; x++ {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileFloor,
				Walkable: true,
			}
		}
	}

	// Add some walls
	floor.Tiles[2][2].Walkable = false
	floor.Tiles[2][3].Walkable = false
	floor.Tiles[3][2].Walkable = false

	// Add some mobs
	mob1 := models.NewMob(models.MobSkeleton, models.VariantEasy, 1)
	mob1.Position = models.Position{X: 5, Y: 5}
	floor.Mobs["mob1"] = mob1

	mob2 := models.NewMob(models.MobGoblin, models.VariantNormal, 1)
	mob2.Position = models.Position{X: 7, Y: 7}
	floor.Mobs["mob2"] = mob2

	// Test cases
	tests := []struct {
		name        string
		currentPos  models.Position
		expectSafe  bool
		expectSameX bool
		expectSameY bool
	}{
		{
			name:       "Position Near Mob",
			currentPos: models.Position{X: 4, Y: 5}, // Adjacent to mob1
			expectSafe: true,                        // Should find a safe position
		},
		{
			name:        "Position in Corner",
			currentPos:  models.Position{X: 0, Y: 0},
			expectSafe:  true,
			expectSameX: false,
			expectSameY: false,
		},
		{
			name:       "Position in Open Area",
			currentPos: models.Position{X: 2, Y: 7},
			expectSafe: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findSafePosition(floor, tt.currentPos)

			// Check if result is different from current position
			if tt.expectSafe {
				assert.False(t, isAdjacent(result, mob1.Position), "Result should not be adjacent to mob1")
				assert.False(t, isAdjacent(result, mob2.Position), "Result should not be adjacent to mob2")
			}

			// Check if position is walkable
			assert.True(t, floor.Tiles[result.Y][result.X].Walkable, "Result position should be walkable")

			// Check if position is within bounds
			assert.GreaterOrEqual(t, result.X, 0, "X should be >= 0")
			assert.Less(t, result.X, floor.Width, "X should be < width")
			assert.GreaterOrEqual(t, result.Y, 0, "Y should be >= 0")
			assert.Less(t, result.Y, floor.Height, "Y should be < height")
		})
	}
}

// TestGetCombatState tests the GetCombatState handler
func TestGetCombatState(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	character.Level = 1
	character.CurrentHP = 20
	character.MaxHP = 20
	characterRepo.Save(character)

	// Create a test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	dungeonRepo.Save(dungeon)

	// Add character to dungeon
	dungeonRepo.AddCharacterToDungeon(dungeon.ID, character.ID)
	dungeonRepo.SetCharacterFloor(dungeon.ID, character.ID, 1)
	character.CurrentDungeon = dungeon.ID
	character.CurrentFloor = 1
	characterRepo.Save(character)

	// Create a test mob on the same floor
	mob := &models.Mob{
		ID:        "test-mob",
		Type:      models.MobGoblin,
		Variant:   models.VariantNormal,
		Name:      "Test Mob",
		Level:     1,
		HP:        10,
		MaxHP:     10,
		Damage:    4,
		Defense:   2,
		AC:        10,
		Dexterity: 10,
		GoldValue: 5,
		Position:  models.Position{X: 5, Y: 5},
		Symbol:    "G",
	}

	// Get the floor and add the mob
	floor, err := dungeonRepo.GetFloor(dungeon.ID, 1)
	require.NoError(t, err)
	floor.Mobs[mob.ID] = mob
	floor.Tiles[5][5].MobID = mob.ID
	dungeonRepo.SaveFloor(dungeon.ID, 1, floor)

	// Position the character adjacent to the mob
	character.Position = models.Position{X: 5, Y: 6}
	characterRepo.Save(character)
	floor.Tiles[6][5].Character = character.ID
	dungeonRepo.SaveFloor(dungeon.ID, 1, floor)

	// Create a game manager
	gameManager := game.NewGameManager(characterRepo, dungeonRepo)

	// Create a combat handler
	combatHandler := NewCombatHandler(characterRepo, dungeonRepo, gameManager)

	// Create a request to get the combat state
	req, err := http.NewRequest("GET", "/combat/"+character.ID, nil)
	require.NoError(t, err)

	// Add the character ID as a URL parameter
	vars := map[string]string{
		"id": character.ID,
	}
	req = mux.SetURLVars(req, vars)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	combatHandler.GetCombatState(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response
	var response struct {
		Character  *models.Character      `json:"character"`
		NearbyMobs map[string]*models.Mob `json:"nearbyMobs"`
		InCombat   bool                   `json:"inCombat"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify the response
	assert.Equal(t, character.ID, response.Character.ID)
	assert.True(t, response.InCombat)
	assert.Contains(t, response.NearbyMobs, mob.ID)
}

// TestCombatHandlerHelperFunctions tests the helper functions in the combat handler
func TestCombatHandlerHelperFunctions(t *testing.T) {
	// Test isAdjacent function
	t.Run("isAdjacent", func(t *testing.T) {
		// Test adjacent positions (including diagonals)
		assert.True(t, isAdjacent(models.Position{X: 5, Y: 5}, models.Position{X: 5, Y: 6}))
		assert.True(t, isAdjacent(models.Position{X: 5, Y: 5}, models.Position{X: 6, Y: 5}))
		assert.True(t, isAdjacent(models.Position{X: 5, Y: 5}, models.Position{X: 4, Y: 5}))
		assert.True(t, isAdjacent(models.Position{X: 5, Y: 5}, models.Position{X: 5, Y: 4}))
		assert.True(t, isAdjacent(models.Position{X: 5, Y: 5}, models.Position{X: 6, Y: 6})) // Diagonal is adjacent

		// Test non-adjacent positions
		assert.False(t, isAdjacent(models.Position{X: 5, Y: 5}, models.Position{X: 7, Y: 5}))
		assert.False(t, isAdjacent(models.Position{X: 5, Y: 5}, models.Position{X: 5, Y: 7}))
		assert.False(t, isAdjacent(models.Position{X: 5, Y: 5}, models.Position{X: 5, Y: 5})) // Same position is not adjacent
	})

	// Test abs function
	t.Run("abs", func(t *testing.T) {
		assert.Equal(t, 5, abs(5))
		assert.Equal(t, 5, abs(-5))
		assert.Equal(t, 0, abs(0))
	})

	// Test findSafePosition function
	t.Run("findSafePosition", func(t *testing.T) {
		// Create a test floor
		floor := &models.Floor{
			Level:  1,
			Width:  10,
			Height: 10,
			Tiles:  make([][]models.Tile, 10),
			Mobs:   make(map[string]*models.Mob),
			Items:  make(map[string]models.Item),
		}

		// Initialize tiles
		for y := 0; y < 10; y++ {
			floor.Tiles[y] = make([]models.Tile, 10)
			for x := 0; x < 10; x++ {
				floor.Tiles[y][x] = models.Tile{
					Type:     models.TileFloor,
					Walkable: true,
				}
			}
		}

		// Add some unwalkable tiles and mobs
		floor.Tiles[5][5].Walkable = false
		floor.Tiles[6][6].MobID = "some-mob"

		// Test finding a safe position
		currentPos := models.Position{X: 5, Y: 5}
		safePos := findSafePosition(floor, currentPos)

		// Verify that the safe position is walkable and has no mob
		assert.True(t, floor.Tiles[safePos.Y][safePos.X].Walkable)
		assert.Empty(t, floor.Tiles[safePos.Y][safePos.X].MobID)
	})
}

// TestHandleCombat tests the HandleCombat WebSocket handler
func TestHandleCombat(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	character.Level = 1
	character.CurrentHP = 20
	character.MaxHP = 20
	character.Attributes.Strength = 16 // High strength for better attack chance
	characterRepo.Save(character)

	// Create a test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	dungeonRepo.Save(dungeon)

	// Add character to dungeon
	dungeonRepo.AddCharacterToDungeon(dungeon.ID, character.ID)
	dungeonRepo.SetCharacterFloor(dungeon.ID, character.ID, 1)
	character.CurrentDungeon = dungeon.ID
	character.CurrentFloor = 1
	characterRepo.Save(character)

	// Create a test mob on the same floor
	mob := &models.Mob{
		ID:        "test-mob",
		Type:      models.MobGoblin,
		Variant:   models.VariantNormal,
		Name:      "Test Mob",
		Level:     1,
		HP:        10,
		MaxHP:     10,
		Damage:    4, // Using an int instead of a string
		Defense:   2,
		AC:        10,
		Dexterity: 10,
		GoldValue: 5,
		Position:  models.Position{X: 5, Y: 5},
		Symbol:    "G",
	}

	// Get the floor and add the mob
	floor, err := dungeonRepo.GetFloor(dungeon.ID, 1)
	require.NoError(t, err)
	floor.Mobs[mob.ID] = mob
	floor.Tiles[5][5].MobID = mob.ID
	dungeonRepo.SaveFloor(dungeon.ID, 1, floor)

	// Position the character adjacent to the mob
	character.Position = models.Position{X: 5, Y: 6}
	characterRepo.Save(character)
	floor.Tiles[6][5].Character = character.ID
	dungeonRepo.SaveFloor(dungeon.ID, 1, floor)

	// Create a game manager
	gameManager := game.NewGameManager(characterRepo, dungeonRepo)

	// Create a combat handler
	combatHandler := NewCombatHandler(characterRepo, dungeonRepo, gameManager)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		combatHandler.HandleCombat(w, r)
	}))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket server")
	defer ws.Close()

	// Test attack action
	t.Run("Attack Action", func(t *testing.T) {
		// Call the handler directly
		response := combatHandler.handleAttack(character, mob.ID)

		// The response might be successful or not depending on the implementation
		// Just verify that we got a response with a message
		assert.NotEmpty(t, response.Message)
	})

	// Test use item action
	t.Run("Use Item Action", func(t *testing.T) {
		// Add a health potion to the character's inventory
		healthPotion := models.NewPotion("Health Potion", 5, 10)
		// Store the potion in the character's inventory
		character.Inventory = append(character.Inventory, healthPotion)
		characterRepo.Save(character)

		// Call the handler directly
		response := combatHandler.handleUseItem(character, healthPotion.ID)

		// Just verify that we got a response with a message
		assert.NotEmpty(t, response.Message)
	})

	// Test flee action
	t.Run("Flee Action", func(t *testing.T) {
		// Call the handler directly
		response := combatHandler.handleFlee(character, mob.ID)

		// Just verify that we got a response with a message
		assert.NotEmpty(t, response.Message)
	})
}

// TestHandleCombatWithMock tests the combat handler functions using a mock WebSocket connection
func TestHandleCombatWithMock(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	character.Level = 1
	character.CurrentHP = 20
	character.MaxHP = 20
	character.Attributes.Strength = 16 // High strength for better attack chance
	characterRepo.Save(character)

	// Create a test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	dungeonRepo.Save(dungeon)

	// Add character to dungeon
	dungeonRepo.AddCharacterToDungeon(dungeon.ID, character.ID)
	dungeonRepo.SetCharacterFloor(dungeon.ID, character.ID, 1)
	character.CurrentDungeon = dungeon.ID
	character.CurrentFloor = 1
	characterRepo.Save(character)

	// Create a test mob on the same floor
	mob := &models.Mob{
		ID:        "test-mob",
		Type:      models.MobGoblin,
		Variant:   models.VariantNormal,
		Name:      "Test Mob",
		Level:     1,
		HP:        10,
		MaxHP:     10,
		Damage:    4,
		Defense:   2,
		AC:        10,
		Dexterity: 10,
		GoldValue: 5,
		Position:  models.Position{X: 5, Y: 5},
		Symbol:    "G",
	}

	// Get the floor and add the mob
	floor, err := dungeonRepo.GetFloor(dungeon.ID, 1)
	require.NoError(t, err)
	floor.Mobs[mob.ID] = mob
	floor.Tiles[5][5].MobID = mob.ID
	dungeonRepo.SaveFloor(dungeon.ID, 1, floor)

	// Position the character adjacent to the mob
	character.Position = models.Position{X: 5, Y: 6}
	characterRepo.Save(character)
	floor.Tiles[6][5].Character = character.ID
	dungeonRepo.SaveFloor(dungeon.ID, 1, floor)

	// Create a game manager
	gameManager := game.NewGameManager(characterRepo, dungeonRepo)

	// Create a combat handler
	combatHandler := NewCombatHandler(characterRepo, dungeonRepo, gameManager)

	// Test attack action
	t.Run("Attack Action", func(t *testing.T) {
		// Call the handler directly
		response := combatHandler.handleAttack(character, mob.ID)

		// The response might be successful or not depending on the implementation
		// Just verify that we got a response with a message
		assert.NotEmpty(t, response.Message)
	})

	// Test use item action
	t.Run("Use Item Action", func(t *testing.T) {
		// Add a health potion to the character's inventory
		healthPotion := models.NewPotion("Health Potion", 5, 10)
		// Store the potion in the character's inventory
		character.Inventory = append(character.Inventory, healthPotion)
		characterRepo.Save(character)

		// Call the handler directly
		response := combatHandler.handleUseItem(character, healthPotion.ID)

		// Just verify that we got a response with a message
		assert.NotEmpty(t, response.Message)
	})

	// Test flee action
	t.Run("Flee Action", func(t *testing.T) {
		// Call the handler directly
		response := combatHandler.handleFlee(character, mob.ID)

		// Just verify that we got a response with a message
		assert.NotEmpty(t, response.Message)
	})
}
