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

	// Create handler using the constructor
	handler := NewDungeonHandler(dungeonRepo, characterRepo)

	// Create a test dungeon
	dungeon := models.NewDungeon("Test Dungeon", 5, 0)
	err := handler.dungeonRepo.Save(dungeon)
	assert.NoError(t, err)

	// Create a test character
	character := models.NewCharacter("Test Character", models.Warrior)
	err = handler.characterRepo.Save(character)
	assert.NoError(t, err)

	// Test cases
	tests := []struct {
		name           string
		dungeonID      string
		characterID    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid join",
			dungeonID:      dungeon.ID,
			characterID:    character.ID,
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Invalid dungeon ID",
			dungeonID:      "non-existent-dungeon",
			characterID:    character.ID,
			expectedStatus: http.StatusNotFound,
			expectedError:  "dungeon not found",
		},
		{
			name:           "Invalid character ID",
			dungeonID:      dungeon.ID,
			characterID:    "non-existent-character",
			expectedStatus: http.StatusNotFound,
			expectedError:  "character not found",
		},
		{
			name:           "Missing character ID",
			dungeonID:      dungeon.ID,
			characterID:    "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Character ID is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			requestBody := map[string]string{
				"characterId": tc.characterID,
			}
			bodyBytes, err := json.Marshal(requestBody)
			assert.NoError(t, err)

			// Create request
			req, err := http.NewRequest("POST", "/dungeons/"+tc.dungeonID+"/join", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Set up router with mux vars
			router := mux.NewRouter()
			router.HandleFunc("/dungeons/{id}/join", handler.JoinDungeon).Methods("POST")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// Check error message if expected
			if tc.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tc.expectedError)
			}

			// If successful join, verify character was added to dungeon
			if tc.expectedStatus == http.StatusOK {
				// Verify character is in dungeon
				updatedDungeon, err := handler.dungeonRepo.GetByID(tc.dungeonID)
				assert.NoError(t, err)
				assert.Contains(t, updatedDungeon.Characters, tc.characterID)

				// Verify character has dungeon reference
				updatedCharacter, err := handler.characterRepo.GetByID(tc.characterID)
				assert.NoError(t, err)
				assert.Equal(t, tc.dungeonID, updatedCharacter.CurrentDungeon)
				assert.Equal(t, 1, updatedCharacter.CurrentFloor)
			}
		})
	}
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
