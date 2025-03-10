package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockCharacterRepository is a mock implementation of the character repository for testing
type MockCharacterRepository struct {
	characters map[string]*models.Character
}

// NewMockCharacterRepository creates a new mock character repository
func NewMockCharacterRepository() *MockCharacterRepository {
	return &MockCharacterRepository{
		characters: make(map[string]*models.Character),
	}
}

// GetAll returns all characters
func (r *MockCharacterRepository) GetAll() []*models.Character {
	characters := make([]*models.Character, 0, len(r.characters))
	for _, character := range r.characters {
		characters = append(characters, character)
	}
	return characters
}

// GetByID returns a character by ID
func (r *MockCharacterRepository) GetByID(id string) (*models.Character, error) {
	character, exists := r.characters[id]
	if !exists {
		return nil, errors.New("character not found")
	}
	return character, nil
}

// Save saves a character
func (r *MockCharacterRepository) Save(character *models.Character) error {
	r.characters[character.ID] = character
	return nil
}

// Delete deletes a character
func (r *MockCharacterRepository) Delete(id string) error {
	if _, exists := r.characters[id]; !exists {
		return errors.New("character not found")
	}
	delete(r.characters, id)
	return nil
}

// Count returns the number of characters
func (r *MockCharacterRepository) Count() int {
	return len(r.characters)
}

// TestCreateCharacter tests the CreateCharacter handler
func TestCreateCharacter(t *testing.T) {
	// Create a new character handler with a mock repository
	handler := &CharacterHandler{
		characterRepo: repositories.NewCharacterRepository(),
	}

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		validateFunc   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Warrior Character",
			requestBody: map[string]interface{}{
				"name":  "TestWarrior",
				"class": "warrior",
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var character models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &character)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, "TestWarrior", character.Name, "Character name should match")
				assert.Equal(t, models.Warrior, character.Class, "Character class should be warrior")
				assert.Equal(t, 1, character.Level, "Character level should be 1")
				assert.NotEmpty(t, character.ID, "Character ID should be generated")
			},
		},
		{
			name: "Valid Mage Character",
			requestBody: map[string]interface{}{
				"name":  "TestMage",
				"class": "mage",
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var character models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &character)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, "TestMage", character.Name, "Character name should match")
				assert.Equal(t, models.Mage, character.Class, "Character class should be mage")
			},
		},
		{
			name: "Missing Name",
			requestBody: map[string]interface{}{
				"class": "warrior",
			},
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check error message in response body
				assert.NotEmpty(t, resp.Body.String(), "Expected error message in response")
			},
		},
		{
			name: "Invalid Class",
			requestBody: map[string]interface{}{
				"name":  "TestInvalid",
				"class": "InvalidClass",
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var character models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &character)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, "TestInvalid", character.Name, "Character name should match")
				assert.Equal(t, models.CharacterClass("InvalidClass"), character.Class, "Character class should match input")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			reqBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err, "Failed to marshal request body")

			req, err := http.NewRequest("POST", "/characters", bytes.NewBuffer(reqBody))
			require.NoError(t, err, "Failed to create request")

			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.CreateCharacter(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code should match expected")

			// Run validation function
			tt.validateFunc(t, rr)
		})
	}
}

// TestGetCharacter tests the GetCharacter handler
func TestGetCharacter(t *testing.T) {
	// Create a new character repository and handler
	repo := repositories.NewCharacterRepository()
	handler := &CharacterHandler{
		characterRepo: repo,
	}

	// Create a test character
	testChar := models.NewCharacter("TestWarrior", models.Warrior)
	repo.Save(testChar)

	tests := []struct {
		name           string
		characterID    string
		expectedStatus int
		validateFunc   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid Character ID",
			characterID:    testChar.ID,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var character models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &character)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, testChar.ID, character.ID, "Character ID should match")
				assert.Equal(t, testChar.Name, character.Name, "Character name should match")
				assert.Equal(t, testChar.Class, character.Class, "Character class should match")
			},
		},
		{
			name:           "Invalid Character ID",
			characterID:    "invalid-id",
			expectedStatus: http.StatusNotFound,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check error message in response body
				assert.NotEmpty(t, resp.Body.String(), "Expected error message in response")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", "/characters/"+tt.characterID, nil)
			require.NoError(t, err, "Failed to create request")

			// Set up router context with character ID
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.characterID,
			})

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.GetCharacter(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code should match expected")

			// Run validation function
			tt.validateFunc(t, rr)
		})
	}
}

// TestGetCharacters tests the GetCharacters handler
func TestGetCharacters(t *testing.T) {
	// Create a new character repository and handler
	repo := repositories.NewCharacterRepository()
	handler := &CharacterHandler{
		characterRepo: repo,
	}

	// Create test characters
	char1 := models.NewCharacter("Warrior1", models.Warrior)
	char2 := models.NewCharacter("Mage1", models.Mage)
	repo.Save(char1)
	repo.Save(char2)

	// Create request
	req, err := http.NewRequest("GET", "/characters", nil)
	require.NoError(t, err, "Failed to create request")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetCharacters(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")

	// Unmarshal response
	var characters []*models.Character
	err = json.Unmarshal(rr.Body.Bytes(), &characters)
	require.NoError(t, err, "Failed to unmarshal response")

	// Check number of characters
	assert.Len(t, characters, 2, "Should return 2 characters")

	// Check character properties
	foundChar1 := false
	foundChar2 := false
	for _, c := range characters {
		if c.ID == char1.ID {
			foundChar1 = true
			assert.Equal(t, char1.Name, c.Name, "Character 1 name should match")
			assert.Equal(t, char1.Class, c.Class, "Character 1 class should match")
		} else if c.ID == char2.ID {
			foundChar2 = true
			assert.Equal(t, char2.Name, c.Name, "Character 2 name should match")
			assert.Equal(t, char2.Class, c.Class, "Character 2 class should match")
		}
	}

	assert.True(t, foundChar1, "Character 1 should be in response")
	assert.True(t, foundChar2, "Character 2 should be in response")
}

// TestDeleteCharacter tests the DeleteCharacter handler
func TestDeleteCharacter(t *testing.T) {
	// Create a new character repository and handler
	repo := repositories.NewCharacterRepository()
	handler := &CharacterHandler{
		characterRepo: repo,
	}

	// Create a test character
	testChar := models.NewCharacter("TestWarrior", models.Warrior)
	repo.Save(testChar)

	tests := []struct {
		name           string
		characterID    string
		expectedStatus int
		validateFunc   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid Character ID",
			characterID:    testChar.ID,
			expectedStatus: http.StatusNoContent,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check that response body is empty (NoContent)
				assert.Empty(t, resp.Body.String(), "Response body should be empty")
			},
		},
		{
			name:           "Invalid Character ID",
			characterID:    "invalid-id",
			expectedStatus: http.StatusNotFound,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check error message in response body
				assert.NotEmpty(t, resp.Body.String(), "Expected error message in response")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset repository for each test
			repo = repositories.NewCharacterRepository()
			handler = &CharacterHandler{
				characterRepo: repo,
			}

			// Add test character
			if tt.name == "Valid Character ID" {
				repo.Save(testChar)
			}

			// Create request
			req, err := http.NewRequest("DELETE", "/characters/"+tt.characterID, nil)
			require.NoError(t, err, "Failed to create request")

			// Set up router context with character ID
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.characterID,
			})

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.DeleteCharacter(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code should match expected")

			// Run validation function
			tt.validateFunc(t, rr)
		})
	}
}
