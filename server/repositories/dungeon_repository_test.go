package repositories

import (
	"testing"

	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/stretchr/testify/assert"
)

func TestNewDungeonRepository(t *testing.T) {
	repo := NewDungeonRepository()
	assert.NotNil(t, repo, "Repository should not be nil")
	assert.NotNil(t, repo.dungeons, "Dungeons map should not be nil")
	assert.Equal(t, 0, len(repo.dungeons), "Dungeons map should be empty")
}

func TestDungeonRepositorySave(t *testing.T) {
	repo := NewDungeonRepository()
	dungeon := models.NewDungeon("TestDungeon", 5, 12345)

	// Save the dungeon
	err := repo.Save(dungeon)
	assert.NoError(t, err, "Save should not return an error")

	// Verify the dungeon was saved
	assert.Equal(t, 1, len(repo.dungeons), "Repository should have 1 dungeon")
	assert.Equal(t, dungeon, repo.dungeons[dungeon.ID], "Saved dungeon should match the original")

	// Update the dungeon
	dungeon.Name = "UpdatedName"
	err = repo.Save(dungeon)
	assert.NoError(t, err, "Save should not return an error when updating")

	// Verify the dungeon was updated
	assert.Equal(t, 1, len(repo.dungeons), "Repository should still have 1 dungeon")
	assert.Equal(t, "UpdatedName", repo.dungeons[dungeon.ID].Name, "Dungeon name should be updated")
}

func TestDungeonRepositoryGetByID(t *testing.T) {
	repo := NewDungeonRepository()
	dungeon := models.NewDungeon("TestDungeon", 5, 12345)

	// Save the dungeon
	_ = repo.Save(dungeon)

	// Get the dungeon by ID
	retrievedDungeon, err := repo.GetByID(dungeon.ID)
	assert.NoError(t, err, "GetByID should not return an error for existing dungeon")
	assert.Equal(t, dungeon, retrievedDungeon, "Retrieved dungeon should match the original")

	// Try to get a non-existent dungeon
	_, err = repo.GetByID("non-existent-id")
	assert.Error(t, err, "GetByID should return an error for non-existent dungeon")
	assert.Equal(t, "dungeon not found", err.Error(), "Error message should be 'dungeon not found'")
}

func TestDungeonRepositoryGetAll(t *testing.T) {
	repo := NewDungeonRepository()

	// Get all dungeons from empty repository
	dungeons := repo.GetAll()
	assert.Equal(t, 0, len(dungeons), "GetAll should return empty slice for empty repository")

	// Add some dungeons
	dungeon1 := models.NewDungeon("Dungeon1", 5, 12345)
	dungeon2 := models.NewDungeon("Dungeon2", 10, 67890)
	_ = repo.Save(dungeon1)
	_ = repo.Save(dungeon2)

	// Get all dungeons
	dungeons = repo.GetAll()
	assert.Equal(t, 2, len(dungeons), "GetAll should return 2 dungeons")

	// Verify the dungeons are in the slice
	found1, found2 := false, false
	for _, d := range dungeons {
		if d.ID == dungeon1.ID {
			found1 = true
		}
		if d.ID == dungeon2.ID {
			found2 = true
		}
	}
	assert.True(t, found1, "Dungeon1 should be in the result")
	assert.True(t, found2, "Dungeon2 should be in the result")
}

func TestDungeonRepositoryDelete(t *testing.T) {
	repo := NewDungeonRepository()
	dungeon := models.NewDungeon("TestDungeon", 5, 12345)

	// Save the dungeon
	_ = repo.Save(dungeon)

	// Delete the dungeon
	err := repo.Delete(dungeon.ID)
	assert.NoError(t, err, "Delete should not return an error for existing dungeon")
	assert.Equal(t, 0, len(repo.dungeons), "Repository should be empty after deletion")

	// Try to delete a non-existent dungeon
	err = repo.Delete("non-existent-id")
	assert.Error(t, err, "Delete should return an error for non-existent dungeon")
	assert.Equal(t, "dungeon not found", err.Error(), "Error message should be 'dungeon not found'")
}

func TestDungeonRepositoryGetFloor(t *testing.T) {
	repo := NewDungeonRepository()
	dungeon := models.NewDungeon("TestDungeon", 5, 12345)

	// Save the dungeon
	_ = repo.Save(dungeon)

	// Get a floor that doesn't exist yet (should generate it)
	floor, err := repo.GetFloor(dungeon.ID, 1)
	assert.NoError(t, err, "GetFloor should not return an error")
	assert.NotNil(t, floor, "Floor should not be nil")
	assert.Equal(t, 1, floor.Level, "Floor level should be 1")

	// Get the same floor again (should return the existing one)
	floor2, err := repo.GetFloor(dungeon.ID, 1)
	assert.NoError(t, err, "GetFloor should not return an error")
	assert.Equal(t, floor, floor2, "Should return the same floor object")

	// Try to get a floor from a non-existent dungeon
	_, err = repo.GetFloor("non-existent-id", 1)
	assert.Error(t, err, "GetFloor should return an error for non-existent dungeon")
	assert.Equal(t, "dungeon not found", err.Error(), "Error message should be 'dungeon not found'")
}

func TestDungeonRepositorySaveFloor(t *testing.T) {
	repo := NewDungeonRepository()
	dungeon := models.NewDungeon("TestDungeon", 5, 12345)

	// Save the dungeon
	_ = repo.Save(dungeon)

	// Create a floor
	floor := &models.Floor{
		Level:  1,
		Width:  20,
		Height: 20,
		Tiles:  make([][]models.Tile, 20),
	}

	// Initialize tiles
	for i := range floor.Tiles {
		floor.Tiles[i] = make([]models.Tile, 20)
		for j := range floor.Tiles[i] {
			floor.Tiles[i][j] = models.Tile{
				Type: models.TileWall,
			}
		}
	}

	// Save the floor
	err := repo.SaveFloor(dungeon.ID, 1, floor)
	assert.NoError(t, err, "SaveFloor should not return an error")

	// Verify the floor was saved
	retrievedFloor, err := repo.GetFloor(dungeon.ID, 1)
	assert.NoError(t, err, "GetFloor should not return an error")
	assert.Equal(t, floor, retrievedFloor, "Retrieved floor should match the saved one")

	// Try to save a floor to a non-existent dungeon
	err = repo.SaveFloor("non-existent-id", 1, floor)
	assert.Error(t, err, "SaveFloor should return an error for non-existent dungeon")
	assert.Equal(t, "dungeon not found", err.Error(), "Error message should be 'dungeon not found'")
}

func TestDungeonRepositoryCharacterManagement(t *testing.T) {
	repo := NewDungeonRepository()
	dungeon := models.NewDungeon("TestDungeon", 5, 12345)
	characterID := "test-character-id"

	// Save the dungeon
	_ = repo.Save(dungeon)

	// Add character to dungeon
	err := repo.AddCharacterToDungeon(dungeon.ID, characterID)
	assert.NoError(t, err, "AddCharacterToDungeon should not return an error")
	assert.Contains(t, dungeon.Characters, characterID, "Character should be added to dungeon")

	// Set character floor
	err = repo.SetCharacterFloor(dungeon.ID, characterID, 2)
	assert.NoError(t, err, "SetCharacterFloor should not return an error")

	// Get character floor
	floor, err := repo.GetCharacterFloor(dungeon.ID, characterID)
	assert.NoError(t, err, "GetCharacterFloor should not return an error")
	assert.Equal(t, 2, floor, "Character floor should be 2")

	// Remove character from dungeon
	err = repo.RemoveCharacterFromDungeon(dungeon.ID, characterID)
	assert.NoError(t, err, "RemoveCharacterFromDungeon should not return an error")
	assert.NotContains(t, dungeon.Characters, characterID, "Character should be removed from dungeon")

	// Try to get floor of removed character
	_, err = repo.GetCharacterFloor(dungeon.ID, characterID)
	assert.Error(t, err, "GetCharacterFloor should return an error for removed character")
	assert.Equal(t, "character not found in dungeon", err.Error(), "Error message should be 'character not found in dungeon'")
}

func TestDungeonRepositoryErrorCases(t *testing.T) {
	repo := NewDungeonRepository()
	characterID := "test-character-id"

	// Try operations on non-existent dungeon
	_, err := repo.GetFloor("non-existent-id", 1)
	assert.Error(t, err, "GetFloor should return an error for non-existent dungeon")

	err = repo.AddCharacterToDungeon("non-existent-id", characterID)
	assert.Error(t, err, "AddCharacterToDungeon should return an error for non-existent dungeon")

	err = repo.RemoveCharacterFromDungeon("non-existent-id", characterID)
	assert.Error(t, err, "RemoveCharacterFromDungeon should return an error for non-existent dungeon")

	_, err = repo.GetCharacterFloor("non-existent-id", characterID)
	assert.Error(t, err, "GetCharacterFloor should return an error for non-existent dungeon")

	err = repo.SetCharacterFloor("non-existent-id", characterID, 1)
	assert.Error(t, err, "SetCharacterFloor should return an error for non-existent dungeon")
}
