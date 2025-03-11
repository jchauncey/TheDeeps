package repositories

import (
	"testing"

	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/stretchr/testify/assert"
)

func TestNewCharacterRepository(t *testing.T) {
	repo := NewCharacterRepository()
	assert.NotNil(t, repo, "Repository should not be nil")
	assert.NotNil(t, repo.characters, "Characters map should not be nil")
	assert.Equal(t, 0, len(repo.characters), "Characters map should be empty")
}

func TestCharacterRepositorySave(t *testing.T) {
	repo := NewCharacterRepository()
	character := models.NewCharacter("TestCharacter", models.Warrior)

	// Save the character
	err := repo.Save(character)
	assert.NoError(t, err, "Save should not return an error")

	// Verify the character was saved
	assert.Equal(t, 1, len(repo.characters), "Repository should have 1 character")
	assert.Equal(t, character, repo.characters[character.ID], "Saved character should match the original")

	// Update the character
	character.Name = "UpdatedName"
	err = repo.Save(character)
	assert.NoError(t, err, "Save should not return an error when updating")

	// Verify the character was updated
	assert.Equal(t, 1, len(repo.characters), "Repository should still have 1 character")
	assert.Equal(t, "UpdatedName", repo.characters[character.ID].Name, "Character name should be updated")
}

func TestCharacterRepositoryGetByID(t *testing.T) {
	repo := NewCharacterRepository()
	character := models.NewCharacter("TestCharacter", models.Warrior)

	// Save the character
	_ = repo.Save(character)

	// Get the character by ID
	retrievedCharacter, err := repo.GetByID(character.ID)
	assert.NoError(t, err, "GetByID should not return an error for existing character")
	assert.Equal(t, character, retrievedCharacter, "Retrieved character should match the original")

	// Try to get a non-existent character
	_, err = repo.GetByID("non-existent-id")
	assert.Error(t, err, "GetByID should return an error for non-existent character")
	assert.Equal(t, "character not found", err.Error(), "Error message should be 'character not found'")
}

func TestCharacterRepositoryGetAll(t *testing.T) {
	repo := NewCharacterRepository()

	// Get all characters from empty repository
	characters := repo.GetAll()
	assert.Equal(t, 0, len(characters), "GetAll should return empty slice for empty repository")

	// Add some characters
	character1 := models.NewCharacter("Character1", models.Warrior)
	character2 := models.NewCharacter("Character2", models.Mage)
	_ = repo.Save(character1)
	_ = repo.Save(character2)

	// Get all characters
	characters = repo.GetAll()
	assert.Equal(t, 2, len(characters), "GetAll should return 2 characters")

	// Verify the characters are in the slice
	found1, found2 := false, false
	for _, c := range characters {
		if c.ID == character1.ID {
			found1 = true
		}
		if c.ID == character2.ID {
			found2 = true
		}
	}
	assert.True(t, found1, "Character1 should be in the result")
	assert.True(t, found2, "Character2 should be in the result")
}

func TestCharacterRepositoryDelete(t *testing.T) {
	repo := NewCharacterRepository()
	character := models.NewCharacter("TestCharacter", models.Warrior)

	// Save the character
	_ = repo.Save(character)

	// Delete the character
	err := repo.Delete(character.ID)
	assert.NoError(t, err, "Delete should not return an error for existing character")
	assert.Equal(t, 0, len(repo.characters), "Repository should be empty after deletion")

	// Try to delete a non-existent character
	err = repo.Delete("non-existent-id")
	assert.Error(t, err, "Delete should return an error for non-existent character")
	assert.Equal(t, "character not found", err.Error(), "Error message should be 'character not found'")
}

func TestCharacterRepositoryCount(t *testing.T) {
	repo := NewCharacterRepository()

	// Count characters in empty repository
	count := repo.Count()
	assert.Equal(t, 0, count, "Count should return 0 for empty repository")

	// Add some characters
	character1 := models.NewCharacter("Character1", models.Warrior)
	character2 := models.NewCharacter("Character2", models.Mage)
	_ = repo.Save(character1)
	_ = repo.Save(character2)

	// Count characters
	count = repo.Count()
	assert.Equal(t, 2, count, "Count should return 2 after adding 2 characters")

	// Delete a character
	_ = repo.Delete(character1.ID)

	// Count characters again
	count = repo.Count()
	assert.Equal(t, 1, count, "Count should return 1 after deleting 1 character")
}

func TestCharacterRepositoryConcurrency(t *testing.T) {
	repo := NewCharacterRepository()
	character := models.NewCharacter("TestCharacter", models.Warrior)

	// Test concurrent reads and writes
	done := make(chan bool)

	// Save the character
	_ = repo.Save(character)

	// Start multiple goroutines to read the character
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = repo.GetByID(character.ID)
			done <- true
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify the character is still there
	retrievedCharacter, err := repo.GetByID(character.ID)
	assert.NoError(t, err, "Character should still exist after concurrent reads")
	assert.Equal(t, character, retrievedCharacter, "Character should be unchanged")
}
