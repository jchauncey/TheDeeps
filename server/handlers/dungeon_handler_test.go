package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateDungeon tests the CreateDungeon handler
func TestCreateDungeon(t *testing.T) {
	// Create repositories
	dungeonRepo := repositories.NewDungeonRepository()
	characterRepo := repositories.NewCharacterRepository()

	// Create handler using the constructor
	handler := NewDungeonHandler(dungeonRepo, characterRepo)

	// Override the map generator for deterministic tests
	handler.mapGenerator = game.NewMapGenerator(12345)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		validateFunc   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Dungeon",
			requestBody: map[string]interface{}{
				"name":   "Test Dungeon",
				"floors": 5,
				"seed":   12345,
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var dungeon models.Dungeon
				err := json.Unmarshal(resp.Body.Bytes(), &dungeon)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate dungeon properties
				assert.Equal(t, "Test Dungeon", dungeon.Name, "Dungeon name should match")
				assert.Equal(t, 5, dungeon.Floors, "Dungeon floors should match")
				assert.Equal(t, int64(12345), dungeon.Seed, "Dungeon seed should match")
				assert.NotEmpty(t, dungeon.ID, "Dungeon ID should be generated")
				// Default difficulty should be "normal"
				assert.Equal(t, "normal", dungeon.Difficulty, "Default difficulty should be 'normal'")
			},
		},
		{
			name: "Dungeon With Custom Difficulty",
			requestBody: map[string]interface{}{
				"name":       "Hard Dungeon",
				"floors":     3,
				"difficulty": "hard",
				"seed":       12345,
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var dungeon models.Dungeon
				err := json.Unmarshal(resp.Body.Bytes(), &dungeon)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate dungeon properties
				assert.Equal(t, "Hard Dungeon", dungeon.Name, "Dungeon name should match")
				assert.Equal(t, 3, dungeon.Floors, "Dungeon floors should match")
				assert.Equal(t, int64(12345), dungeon.Seed, "Dungeon seed should match")
				assert.NotEmpty(t, dungeon.ID, "Dungeon ID should be generated")
				// Custom difficulty should be preserved
				assert.Equal(t, "hard", dungeon.Difficulty, "Difficulty should be 'hard'")
			},
		},
		{
			name: "Missing Name",
			requestBody: map[string]interface{}{
				"floors": 5,
				"seed":   12345,
			},
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check error message in response body
				assert.NotEmpty(t, resp.Body.String(), "Expected error message in response")
			},
		},
		{
			name: "Invalid Floor Count",
			requestBody: map[string]interface{}{
				"name":   "Test Dungeon",
				"floors": 0,
				"seed":   12345,
			},
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check error message in response body
				assert.NotEmpty(t, resp.Body.String(), "Expected error message in response")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			reqBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err, "Failed to marshal request body")

			req, err := http.NewRequest("POST", "/dungeons", bytes.NewBuffer(reqBody))
			require.NoError(t, err, "Failed to create request")

			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.CreateDungeon(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code should match expected")

			// Run validation function
			tt.validateFunc(t, rr)
		})
	}
}

// TestGetDungeons tests the GetDungeons handler
func TestGetDungeons(t *testing.T) {
	// Create a new dungeon repository and handler
	dungeonRepo := repositories.NewDungeonRepository()
	characterRepo := repositories.NewCharacterRepository()
	mapGenerator := game.NewMapGenerator(12345)
	handler := &DungeonHandler{
		dungeonRepo:   dungeonRepo,
		characterRepo: characterRepo,
		mapGenerator:  mapGenerator,
	}

	// Create test dungeons
	dungeon1 := models.NewDungeon("Dungeon 1", 3, 12345)
	dungeon2 := models.NewDungeon("Dungeon 2", 5, 67890)
	dungeonRepo.Save(dungeon1)
	dungeonRepo.Save(dungeon2)

	// Create request
	req, err := http.NewRequest("GET", "/dungeons", nil)
	require.NoError(t, err, "Failed to create request")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetDungeons(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")

	// Unmarshal response
	var dungeons []*models.Dungeon
	err = json.Unmarshal(rr.Body.Bytes(), &dungeons)
	require.NoError(t, err, "Failed to unmarshal response")

	// Check number of dungeons
	assert.Len(t, dungeons, 2, "Should return 2 dungeons")

	// Check dungeon properties
	foundDungeon1 := false
	foundDungeon2 := false
	for _, d := range dungeons {
		if d.ID == dungeon1.ID {
			foundDungeon1 = true
			assert.Equal(t, dungeon1.Name, d.Name, "Dungeon 1 name should match")
			assert.Equal(t, dungeon1.Floors, d.Floors, "Dungeon 1 floors should match")
		} else if d.ID == dungeon2.ID {
			foundDungeon2 = true
			assert.Equal(t, dungeon2.Name, d.Name, "Dungeon 2 name should match")
			assert.Equal(t, dungeon2.Floors, d.Floors, "Dungeon 2 floors should match")
		}
	}

	assert.True(t, foundDungeon1, "Dungeon 1 should be in response")
	assert.True(t, foundDungeon2, "Dungeon 2 should be in response")
}

// TestGetFloor tests the GetFloor handler
func TestGetFloor(t *testing.T) {
	// Create a new dungeon repository and handler
	dungeonRepo := repositories.NewDungeonRepository()
	characterRepo := repositories.NewCharacterRepository()
	mapGenerator := game.NewMapGenerator(12345)
	handler := &DungeonHandler{
		dungeonRepo:   dungeonRepo,
		characterRepo: characterRepo,
		mapGenerator:  mapGenerator,
	}

	// Create a test dungeon
	testDungeon := models.NewDungeon("Test Dungeon", 5, 12345)

	// Generate a floor
	floor := testDungeon.GenerateFloor(1)
	mapGenerator.GenerateFloor(floor, 1, false)

	dungeonRepo.Save(testDungeon)

	tests := []struct {
		name           string
		dungeonID      string
		floorLevel     string
		expectedStatus int
		validateFunc   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid Floor",
			dungeonID:      testDungeon.ID,
			floorLevel:     "1",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var floor models.Floor
				err := json.Unmarshal(resp.Body.Bytes(), &floor)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate floor properties
				assert.Equal(t, 1, floor.Level, "Floor level should be 1")
				assert.Greater(t, floor.Width, 0, "Width should be positive")
				assert.Greater(t, floor.Height, 0, "Height should be positive")
				assert.NotEmpty(t, floor.Tiles, "Tiles should be initialized")
			},
		},
		{
			name:           "Invalid Dungeon ID",
			dungeonID:      "invalid-id",
			floorLevel:     "1",
			expectedStatus: http.StatusNotFound,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check error message in response body
				assert.NotEmpty(t, resp.Body.String(), "Expected error message in response")
			},
		},
		{
			name:           "Invalid Floor Level",
			dungeonID:      testDungeon.ID,
			floorLevel:     "invalid",
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check error message in response body
				assert.NotEmpty(t, resp.Body.String(), "Expected error message in response")
			},
		},
		{
			name:           "Floor Level Out of Range",
			dungeonID:      testDungeon.ID,
			floorLevel:     "10",
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check error message in response body
				assert.NotEmpty(t, resp.Body.String(), "Expected error message in response")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", "/dungeons/"+tt.dungeonID+"/floor/"+tt.floorLevel, nil)
			require.NoError(t, err, "Failed to create request")

			// Set up router context with dungeon ID and floor level
			req = mux.SetURLVars(req, map[string]string{
				"id":    tt.dungeonID,
				"level": tt.floorLevel,
			})

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.GetFloor(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code should match expected")

			// Run validation function
			tt.validateFunc(t, rr)
		})
	}
}

// TestJoinDungeon tests the JoinDungeon handler
func TestJoinDungeon(t *testing.T) {
	// Create repositories
	dungeonRepo := repositories.NewDungeonRepository()
	characterRepo := repositories.NewCharacterRepository()

	// Create a test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	dungeonRepo.Save(dungeon)

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	characterRepo.Save(character)

	// Generate the first floor with an entrance room
	floor := dungeon.GenerateFloor(1)
	mapGenerator := game.NewMapGenerator(12345)
	mapGenerator.GenerateFloorWithDifficulty(floor, 1, false, "normal")

	// Create the handler
	handler := NewDungeonHandler(dungeonRepo, characterRepo)

	// Create a request
	requestBody := map[string]string{
		"characterId": character.ID,
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/dungeons/"+dungeon.ID+"/join", bytes.NewBuffer(body))
	require.NoError(t, err)

	// Set up the router with the route parameter
	router := mux.NewRouter()
	router.HandleFunc("/dungeons/{id}/join", handler.JoinDungeon).Methods("POST")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response
	var responseFloor models.Floor
	err = json.Unmarshal(rr.Body.Bytes(), &responseFloor)
	require.NoError(t, err)

	// Get the updated character
	updatedCharacter, err := characterRepo.GetByID(character.ID)
	require.NoError(t, err)

	// Verify the character is in the dungeon
	assert.Equal(t, dungeon.ID, updatedCharacter.CurrentDungeon)
	assert.Equal(t, 1, updatedCharacter.CurrentFloor)

	// Find the entrance room
	var entranceRoom *models.Room
	for i := range responseFloor.Rooms {
		if responseFloor.Rooms[i].Type == models.RoomEntrance {
			entranceRoom = &responseFloor.Rooms[i]
			break
		}
	}
	require.NotNil(t, entranceRoom, "Entrance room should exist")

	// Verify the character is positioned in the entrance room
	assert.GreaterOrEqual(t, updatedCharacter.Position.X, entranceRoom.X)
	assert.Less(t, updatedCharacter.Position.X, entranceRoom.X+entranceRoom.Width)
	assert.GreaterOrEqual(t, updatedCharacter.Position.Y, entranceRoom.Y)
	assert.Less(t, updatedCharacter.Position.Y, entranceRoom.Y+entranceRoom.Height)

	// Verify the character is on a walkable tile
	assert.True(t, responseFloor.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Walkable)

	// Verify the character is not on stairs
	assert.NotEqual(t, models.TileDownStairs, responseFloor.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Type)
	assert.NotEqual(t, models.TileUpStairs, responseFloor.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Type)

	// Verify the tile has the character ID
	assert.Equal(t, character.ID, responseFloor.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Character)
}

// TestJoinDungeonWithObstaclesInEntranceRoom tests the JoinDungeon handler with obstacles in the entrance room
func TestJoinDungeonWithObstaclesInEntranceRoom(t *testing.T) {
	// Create repositories
	dungeonRepo := repositories.NewDungeonRepository()
	characterRepo := repositories.NewCharacterRepository()

	// Create a test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	dungeonRepo.Save(dungeon)

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	characterRepo.Save(character)

	// Generate the first floor with an entrance room
	floor := dungeon.GenerateFloor(1)
	mapGenerator := game.NewMapGenerator(12345)
	mapGenerator.GenerateFloorWithDifficulty(floor, 1, false, "normal")

	// Find the entrance room
	var entranceRoom *models.Room
	for i := range floor.Rooms {
		if floor.Rooms[i].Type == models.RoomEntrance {
			entranceRoom = &floor.Rooms[i]
			break
		}
	}
	require.NotNil(t, entranceRoom, "Entrance room should exist")

	// Place obstacles in the center of the entrance room
	centerX := entranceRoom.X + entranceRoom.Width/2
	centerY := entranceRoom.Y + entranceRoom.Height/2

	// Place a down stairs at the center
	floor.Tiles[centerY][centerX].Type = models.TileDownStairs
	floor.Tiles[centerY][centerX].Walkable = true
	floor.DownStairs = append(floor.DownStairs, models.Position{X: centerX, Y: centerY})

	// Place a mob near the center
	mobID := "test-mob"
	mob := models.NewMob(models.MobGoblin, models.VariantNormal, 1)
	mob.ID = mobID
	mob.Position = models.Position{X: centerX + 1, Y: centerY}
	floor.Mobs[mobID] = mob
	floor.Tiles[centerY][centerX+1].MobID = mobID

	// Create the handler
	handler := NewDungeonHandler(dungeonRepo, characterRepo)

	// Create a request
	requestBody := map[string]string{
		"characterId": character.ID,
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/dungeons/"+dungeon.ID+"/join", bytes.NewBuffer(body))
	require.NoError(t, err)

	// Set up the router with the route parameter
	router := mux.NewRouter()
	router.HandleFunc("/dungeons/{id}/join", handler.JoinDungeon).Methods("POST")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response
	var responseFloor models.Floor
	err = json.Unmarshal(rr.Body.Bytes(), &responseFloor)
	require.NoError(t, err)

	// Get the updated character
	updatedCharacter, err := characterRepo.GetByID(character.ID)
	require.NoError(t, err)

	// Verify the character is in the dungeon
	assert.Equal(t, dungeon.ID, updatedCharacter.CurrentDungeon)
	assert.Equal(t, 1, updatedCharacter.CurrentFloor)

	// Verify the character is positioned in the entrance room
	assert.GreaterOrEqual(t, updatedCharacter.Position.X, entranceRoom.X)
	assert.Less(t, updatedCharacter.Position.X, entranceRoom.X+entranceRoom.Width)
	assert.GreaterOrEqual(t, updatedCharacter.Position.Y, entranceRoom.Y)
	assert.Less(t, updatedCharacter.Position.Y, entranceRoom.Y+entranceRoom.Height)

	// Verify the character is on a walkable tile
	assert.True(t, responseFloor.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Walkable)

	// Verify the character is not on stairs
	assert.NotEqual(t, models.TileDownStairs, responseFloor.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Type)
	assert.NotEqual(t, models.TileUpStairs, responseFloor.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Type)

	// Verify the character is not on the same tile as a mob
	assert.Equal(t, "", responseFloor.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].MobID)

	// Verify the tile has the character ID
	assert.Equal(t, character.ID, responseFloor.Tiles[updatedCharacter.Position.Y][updatedCharacter.Position.X].Character)

	// Verify the character is not at the center (where we placed the stairs)
	assert.False(t, updatedCharacter.Position.X == centerX && updatedCharacter.Position.Y == centerY)

	// Verify the character is not at the mob position
	assert.False(t, updatedCharacter.Position.X == centerX+1 && updatedCharacter.Position.Y == centerY)
}

// TestNewDungeonHandler tests the NewDungeonHandler function
func TestNewDungeonHandler(t *testing.T) {
	// Create repositories
	dungeonRepo := repositories.NewDungeonRepository()
	characterRepo := repositories.NewCharacterRepository()

	// Create handler using the constructor
	handler := NewDungeonHandler(dungeonRepo, characterRepo)

	// Verify handler is initialized correctly
	assert.NotNil(t, handler, "Handler should not be nil")
	assert.NotNil(t, handler.dungeonRepo, "Dungeon repository should not be nil")
	assert.NotNil(t, handler.characterRepo, "Character repository should not be nil")
	assert.NotNil(t, handler.mapGenerator, "Map generator should not be nil")

	// Verify the repositories are the ones we passed in
	assert.Same(t, dungeonRepo, handler.dungeonRepo, "Dungeon repository should be the same instance")
	assert.Same(t, characterRepo, handler.characterRepo, "Character repository should be the same instance")
}

// TestSharedRepositories tests that the repositories are shared correctly
func TestSharedRepositories(t *testing.T) {
	// Create shared repositories
	dungeonRepo := repositories.NewDungeonRepository()
	characterRepo := repositories.NewCharacterRepository()

	// Create a test character in the character repository
	character := models.NewCharacter("Test Character", models.Warrior)
	err := characterRepo.Save(character)
	assert.NoError(t, err, "Failed to save character")

	// Create a test dungeon in the dungeon repository
	dungeon := models.NewDungeon("Test Dungeon", 5, 0)
	err = dungeonRepo.Save(dungeon)
	assert.NoError(t, err, "Failed to save dungeon")

	// Create handler using the shared repositories
	handler := NewDungeonHandler(dungeonRepo, characterRepo)

	// Create request to join the dungeon
	requestBody := map[string]string{
		"characterId": character.ID,
	}
	bodyBytes, err := json.Marshal(requestBody)
	assert.NoError(t, err, "Failed to marshal request body")

	req, err := http.NewRequest("POST", "/dungeons/"+dungeon.ID+"/join", bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	// Set up router with mux vars
	req = mux.SetURLVars(req, map[string]string{
		"id": dungeon.ID,
	})

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler directly
	handler.JoinDungeon(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")

	// Verify character was added to dungeon
	updatedDungeon, err := dungeonRepo.GetByID(dungeon.ID)
	assert.NoError(t, err, "Failed to get updated dungeon")
	assert.Contains(t, updatedDungeon.Characters, character.ID, "Character should be added to dungeon")

	// Verify character has dungeon reference
	updatedCharacter, err := characterRepo.GetByID(character.ID)
	assert.NoError(t, err, "Failed to get updated character")
	assert.Equal(t, dungeon.ID, updatedCharacter.CurrentDungeon, "Character should reference dungeon")
	assert.Equal(t, 1, updatedCharacter.CurrentFloor, "Character should be on floor 1")
}

func TestGenerateTestRoom(t *testing.T) {
	// Create repositories
	dungeonRepo := repositories.NewDungeonRepository()
	characterRepo := repositories.NewCharacterRepository()

	// Create handler
	handler := NewDungeonHandler(dungeonRepo, characterRepo)

	// Test cases
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkFunc      func(*testing.T, *models.Floor)
	}{
		{
			name:           "Default Entrance Room",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, floor *models.Floor) {
				// Check floor dimensions
				assert.Equal(t, 20, floor.Width)
				assert.Equal(t, 20, floor.Height)

				// Check that there's exactly one room
				assert.Equal(t, 1, len(floor.Rooms))

				// Check room properties
				room := floor.Rooms[0]
				assert.Equal(t, models.RoomEntrance, room.Type)
				assert.Equal(t, 8, room.Width)
				assert.Equal(t, 8, room.Height)
				assert.True(t, room.Explored)

				// Check that down stairs exist
				assert.Equal(t, 1, len(floor.DownStairs))

				// Check that a character is placed
				characterFound := false
				for y := 0; y < floor.Height; y++ {
					for x := 0; x < floor.Width; x++ {
						if floor.Tiles[y][x].Character != "" {
							characterFound = true
							break
						}
					}
					if characterFound {
						break
					}
				}
				assert.True(t, characterFound, "Character should be placed in the room")
			},
		},
		{
			name:           "Treasure Room",
			queryParams:    "?type=treasure",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, floor *models.Floor) {
				// Check room type
				assert.Equal(t, models.RoomTreasure, floor.Rooms[0].Type)

				// Check that items exist
				assert.Greater(t, len(floor.Items), 0, "Treasure room should have items")
			},
		},
		{
			name:           "Boss Room",
			queryParams:    "?type=boss",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, floor *models.Floor) {
				// Check room type
				assert.Equal(t, models.RoomBoss, floor.Rooms[0].Type)

				// Check that a boss mob exists
				assert.Equal(t, 1, len(floor.Mobs), "Boss room should have a boss mob")

				// Find the boss mob
				var boss *models.Mob
				for _, mob := range floor.Mobs {
					boss = mob
					break
				}
				assert.NotNil(t, boss)
				assert.Equal(t, models.VariantBoss, boss.Variant)
			},
		},
		{
			name:           "Custom Size Room",
			queryParams:    "?width=30&height=25&roomWidth=10&roomHeight=12",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, floor *models.Floor) {
				// Check floor dimensions
				assert.Equal(t, 30, floor.Width)
				assert.Equal(t, 25, floor.Height)

				// Check room dimensions
				assert.Equal(t, 10, floor.Rooms[0].Width)
				assert.Equal(t, 12, floor.Rooms[0].Height)
			},
		},
		{
			name:           "Invalid Room Type",
			queryParams:    "?type=invalid",
			expectedStatus: http.StatusBadRequest,
			checkFunc:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", "/test/room"+tt.queryParams, nil)
			assert.NoError(t, err)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create router and register handler
			router := mux.NewRouter()
			router.HandleFunc("/test/room", handler.GenerateTestRoom).Methods("GET")

			// Serve request
			router.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// If we expect success, check the response
			if tt.expectedStatus == http.StatusOK && tt.checkFunc != nil {
				var floor models.Floor
				err := json.Unmarshal(rr.Body.Bytes(), &floor)
				assert.NoError(t, err)

				// Run the check function
				tt.checkFunc(t, &floor)
			}
		})
	}
}
