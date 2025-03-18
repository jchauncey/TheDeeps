package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	// Create a repository
	repo := repositories.NewCharacterRepository()

	// Create handler with the repository
	handler := NewCharacterHandler(repo)

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
	// Create a repository
	repo := repositories.NewCharacterRepository()

	// Create a test character
	character := models.NewCharacter("Test Character", models.Warrior)
	repo.Save(character)

	// Create handler with the repository
	handler := NewCharacterHandler(repo)

	tests := []struct {
		name           string
		characterID    string
		expectedStatus int
		validateFunc   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid Character ID",
			characterID:    character.ID,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var responseChar models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &responseChar)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, character.ID, responseChar.ID, "Character ID should match")
				assert.Equal(t, character.Name, responseChar.Name, "Character name should match")
				assert.Equal(t, character.Class, responseChar.Class, "Character class should match")
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
	// Create a repository
	repo := repositories.NewCharacterRepository()

	// Create test characters
	character1 := models.NewCharacter("Character 1", models.Warrior)
	character2 := models.NewCharacter("Character 2", models.Mage)
	repo.Save(character1)
	repo.Save(character2)

	// Create handler with the repository
	handler := NewCharacterHandler(repo)

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
		if c.ID == character1.ID {
			foundChar1 = true
			assert.Equal(t, character1.Name, c.Name, "Character 1 name should match")
			assert.Equal(t, character1.Class, c.Class, "Character 1 class should match")
		} else if c.ID == character2.ID {
			foundChar2 = true
			assert.Equal(t, character2.Name, c.Name, "Character 2 name should match")
			assert.Equal(t, character2.Class, c.Class, "Character 2 class should match")
		}
	}

	assert.True(t, foundChar1, "Character 1 should be in response")
	assert.True(t, foundChar2, "Character 2 should be in response")
}

// TestDeleteCharacter tests the DeleteCharacter handler
func TestDeleteCharacter(t *testing.T) {
	// Create a repository
	repo := repositories.NewCharacterRepository()

	// Create a test character
	character := models.NewCharacter("Test Character", models.Warrior)
	repo.Save(character)

	// Create handler with the repository
	handler := NewCharacterHandler(repo)

	tests := []struct {
		name           string
		characterID    string
		expectedStatus int
		validateFunc   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid Character ID",
			characterID:    character.ID,
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
			handler = NewCharacterHandler(repo)

			// Add test character
			if tt.name == "Valid Character ID" {
				repo.Save(character)
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

// TestNewCharacterHandler tests the NewCharacterHandler function
func TestNewCharacterHandler(t *testing.T) {
	// Create a repository
	repo := repositories.NewCharacterRepository()

	// Create handler using the constructor
	handler := NewCharacterHandler(repo)

	// Verify handler is initialized correctly
	assert.NotNil(t, handler, "Handler should not be nil")
	assert.NotNil(t, handler.characterRepo, "Character repository should not be nil")

	// Verify the repository is the one we passed in
	assert.Same(t, repo, handler.characterRepo, "Character repository should be the same instance")
}

// TestSaveCharacter tests the SaveCharacter handler
func TestSaveCharacter(t *testing.T) {
	// Create a repository
	repo := repositories.NewCharacterRepository()

	// Create a test character
	character := models.NewCharacter("Test Character", models.Warrior)
	repo.Save(character)

	// Create handler with the repository
	handler := NewCharacterHandler(repo)

	tests := []struct {
		name           string
		characterID    string
		requestBody    interface{}
		expectedStatus int
		validateFunc   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:        "Valid Save",
			characterID: character.ID,
			requestBody: map[string]interface{}{
				"position": map[string]interface{}{
					"x": 10,
					"y": 15,
				},
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var updatedChar models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &updatedChar)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate updated character properties
				assert.Equal(t, character.ID, updatedChar.ID, "Character ID should match")
				assert.Equal(t, 10, updatedChar.Position.X, "X position should be updated")
				assert.Equal(t, 15, updatedChar.Position.Y, "Y position should be updated")
			},
		},
		{
			name:        "Invalid Character ID",
			characterID: "invalid-id",
			requestBody: map[string]interface{}{
				"position": map[string]interface{}{
					"x": 10,
					"y": 15,
				},
			},
			expectedStatus: http.StatusNotFound,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Check error message in response body
				assert.NotEmpty(t, resp.Body.String(), "Expected error message in response")
			},
		},
		{
			name:           "Invalid Request Body",
			characterID:    character.ID,
			requestBody:    nil, // This will cause JSON decoding to fail
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
			var reqBody []byte
			var err error

			if tt.requestBody != nil {
				reqBody, err = json.Marshal(tt.requestBody)
				require.NoError(t, err, "Failed to marshal request body")
			} else {
				// Create an invalid JSON for the "Invalid Request Body" test
				reqBody = []byte("{invalid json")
			}

			req, err := http.NewRequest("POST", "/characters/"+tt.characterID+"/save", bytes.NewBuffer(reqBody))
			require.NoError(t, err, "Failed to create request")

			req.Header.Set("Content-Type", "application/json")

			// Set up router context with character ID
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.characterID,
			})

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.SaveCharacter(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code should match expected")

			// Run validation function
			tt.validateFunc(t, rr)
		})
	}
}

// TestGetCharacterFloor tests the GetCharacterFloor handler
func TestGetCharacterFloor(t *testing.T) {
	// Create a repository
	repo := repositories.NewCharacterRepository()

	// Create a test character
	character := models.NewCharacter("Test Character", models.Warrior)
	character.CurrentFloor = 3
	repo.Save(character)

	// Create handler with the repository
	handler := NewCharacterHandler(repo)

	tests := []struct {
		name           string
		characterID    string
		expectedStatus int
		expectedFloor  int
	}{
		{
			name:           "Valid Character ID",
			characterID:    character.ID,
			expectedStatus: http.StatusOK,
			expectedFloor:  character.CurrentFloor,
		},
		{
			name:           "Invalid Character ID",
			characterID:    "invalid-id",
			expectedStatus: http.StatusNotFound,
			expectedFloor:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", "/characters/"+tt.characterID+"/floor", nil)
			require.NoError(t, err, "Failed to create request")

			// Set up router context with character ID
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.characterID,
			})

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.GetCharacterFloor(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code should match expected")

			// If successful, check the floor value
			if tt.expectedStatus == http.StatusOK {
				var response map[string]int
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err, "Failed to unmarshal response")
				assert.Equal(t, tt.expectedFloor, response["floor"], "Floor should match expected")
			}
		})
	}
}

// TestCreateCharacterWithCustomAttributes tests creating characters with custom attributes
func TestCreateCharacterWithCustomAttributes(t *testing.T) {
	// Create a repository
	repo := repositories.NewCharacterRepository()

	// Create handler with the repository
	handler := NewCharacterHandler(repo)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		validateFunc   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "Warrior with Custom Attributes",
			requestBody: map[string]interface{}{
				"name":  "CustomWarrior",
				"class": "warrior",
				"attributes": map[string]interface{}{
					"strength":     16,
					"dexterity":    14,
					"constitution": 15,
					"intelligence": 10,
					"wisdom":       12,
					"charisma":     8,
				},
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var character models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &character)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, "CustomWarrior", character.Name, "Character name should match")
				assert.Equal(t, models.Warrior, character.Class, "Character class should be warrior")
				assert.Equal(t, 16, character.Attributes.Strength, "Strength should match")
				assert.Equal(t, 14, character.Attributes.Dexterity, "Dexterity should match")
				assert.Equal(t, 15, character.Attributes.Constitution, "Constitution should match")
				assert.Equal(t, 10, character.Attributes.Intelligence, "Intelligence should match")
				assert.Equal(t, 12, character.Attributes.Wisdom, "Wisdom should match")
				assert.Equal(t, 8, character.Attributes.Charisma, "Charisma should match")
				assert.Equal(t, 30, character.MaxHP, "MaxHP should be calculated correctly for warrior")
				assert.Equal(t, 0, character.MaxMana, "MaxMana should be 0 for warrior")
			},
		},
		{
			name: "Mage with Custom Attributes",
			requestBody: map[string]interface{}{
				"name":  "CustomMage",
				"class": "mage",
				"attributes": map[string]interface{}{
					"strength":     8,
					"dexterity":    14,
					"constitution": 12,
					"intelligence": 18,
					"wisdom":       15,
					"charisma":     10,
				},
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var character models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &character)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, "CustomMage", character.Name, "Character name should match")
				assert.Equal(t, models.Mage, character.Class, "Character class should be mage")
				assert.Equal(t, 8, character.Attributes.Strength, "Strength should match")
				assert.Equal(t, 14, character.Attributes.Dexterity, "Dexterity should match")
				assert.Equal(t, 12, character.Attributes.Constitution, "Constitution should match")
				assert.Equal(t, 18, character.Attributes.Intelligence, "Intelligence should match")
				assert.Equal(t, 15, character.Attributes.Wisdom, "Wisdom should match")
				assert.Equal(t, 10, character.Attributes.Charisma, "Charisma should match")
				assert.Equal(t, 20, character.MaxHP, "MaxHP should be calculated correctly for mage")
				assert.Equal(t, 28, character.MaxMana, "MaxMana should be calculated correctly for mage")
			},
		},
		{
			name: "Cleric with Custom Attributes",
			requestBody: map[string]interface{}{
				"name":  "CustomCleric",
				"class": "cleric",
				"attributes": map[string]interface{}{
					"strength":     12,
					"dexterity":    10,
					"constitution": 14,
					"intelligence": 12,
					"wisdom":       18,
					"charisma":     14,
				},
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var character models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &character)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, "CustomCleric", character.Name, "Character name should match")
				assert.Equal(t, models.Cleric, character.Class, "Character class should be cleric")
				assert.Equal(t, 12, character.Attributes.Strength, "Strength should match")
				assert.Equal(t, 10, character.Attributes.Dexterity, "Dexterity should match")
				assert.Equal(t, 14, character.Attributes.Constitution, "Constitution should match")
				assert.Equal(t, 12, character.Attributes.Intelligence, "Intelligence should match")
				assert.Equal(t, 18, character.Attributes.Wisdom, "Wisdom should match")
				assert.Equal(t, 14, character.Attributes.Charisma, "Charisma should match")
				assert.Equal(t, 24, character.MaxHP, "MaxHP should be calculated correctly for cleric")
				assert.Equal(t, 28, character.MaxMana, "MaxMana should be calculated correctly for cleric")
			},
		},
		{
			name: "Bard with Custom Attributes",
			requestBody: map[string]interface{}{
				"name":  "CustomBard",
				"class": "bard",
				"attributes": map[string]interface{}{
					"strength":     10,
					"dexterity":    14,
					"constitution": 12,
					"intelligence": 14,
					"wisdom":       10,
					"charisma":     18,
				},
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var character models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &character)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, "CustomBard", character.Name, "Character name should match")
				assert.Equal(t, models.Bard, character.Class, "Character class should be bard")
				assert.Equal(t, 10, character.Attributes.Strength, "Strength should match")
				assert.Equal(t, 14, character.Attributes.Dexterity, "Dexterity should match")
				assert.Equal(t, 12, character.Attributes.Constitution, "Constitution should match")
				assert.Equal(t, 14, character.Attributes.Intelligence, "Intelligence should match")
				assert.Equal(t, 10, character.Attributes.Wisdom, "Wisdom should match")
				assert.Equal(t, 18, character.Attributes.Charisma, "Charisma should match")
				assert.Equal(t, 22, character.MaxHP, "MaxHP should be calculated correctly for bard")
				assert.Equal(t, 26, character.MaxMana, "MaxMana should be calculated correctly for bard")
			},
		},
		{
			name: "Paladin with Custom Attributes",
			requestBody: map[string]interface{}{
				"name":  "CustomPaladin",
				"class": "paladin",
				"attributes": map[string]interface{}{
					"strength":     16,
					"dexterity":    12,
					"constitution": 14,
					"intelligence": 10,
					"wisdom":       14,
					"charisma":     16,
				},
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var character models.Character
				err := json.Unmarshal(resp.Body.Bytes(), &character)
				require.NoError(t, err, "Failed to unmarshal response")

				// Validate character properties
				assert.Equal(t, "CustomPaladin", character.Name, "Character name should match")
				assert.Equal(t, models.Paladin, character.Class, "Character class should be paladin")
				assert.Equal(t, 16, character.Attributes.Strength, "Strength should match")
				assert.Equal(t, 12, character.Attributes.Dexterity, "Dexterity should match")
				assert.Equal(t, 14, character.Attributes.Constitution, "Constitution should match")
				assert.Equal(t, 10, character.Attributes.Intelligence, "Intelligence should match")
				assert.Equal(t, 14, character.Attributes.Wisdom, "Wisdom should match")
				assert.Equal(t, 16, character.Attributes.Charisma, "Charisma should match")
				assert.Equal(t, 24, character.MaxHP, "MaxHP should be calculated correctly for paladin")
				assert.Equal(t, 21, character.MaxMana, "MaxMana should be calculated correctly for paladin")
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

// TestCreateCharacterLimit tests the character creation limit
func TestCreateCharacterLimit(t *testing.T) {
	// Create a repository that already has MAX_CHARACTERS characters
	repo := repositories.NewCharacterRepository()

	// Add MAX_CHARACTERS characters to the repository
	for i := 0; i < 10; i++ {
		character := models.NewCharacter(fmt.Sprintf("Character %d", i), models.Warrior)
		repo.Save(character)
	}

	// Create handler with the repository
	handler := NewCharacterHandler(repo)

	// Now try to create one more character, which should fail
	reqBody, err := json.Marshal(map[string]interface{}{
		"name":  "LimitTest",
		"class": "warrior",
	})
	require.NoError(t, err, "Failed to marshal request body")

	req, err := http.NewRequest("POST", "/characters", bytes.NewBuffer(reqBody))
	require.NoError(t, err, "Failed to create request")

	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.CreateCharacter(rr, req)

	// Check status code
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should match expected")

	// Check error message
	assert.Contains(t, rr.Body.String(), "Maximum number of characters", "Expected error message about character limit")
}
