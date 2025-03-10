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
	// Create a new dungeon handler with a repository
	dungeonRepo := repositories.NewDungeonRepository()
	characterRepo := repositories.NewCharacterRepository()
	mapGenerator := game.NewMapGenerator(12345)
	handler := &DungeonHandler{
		dungeonRepo:   dungeonRepo,
		characterRepo: characterRepo,
		mapGenerator:  mapGenerator,
	}

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
