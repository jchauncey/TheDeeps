package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDungeon(t *testing.T) {
	tests := []struct {
		name        string
		dungeonName string
		floors      int
		seed        int64
	}{
		{
			name:        "Standard Dungeon",
			dungeonName: "Test Dungeon",
			floors:      5,
			seed:        12345,
		},
		{
			name:        "Single Floor Dungeon",
			dungeonName: "Small Dungeon",
			floors:      1,
			seed:        0, // Should generate a seed based on time
		},
		{
			name:        "Large Dungeon",
			dungeonName: "Mega Dungeon",
			floors:      20,
			seed:        99999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeTime := time.Now()
			dungeon := NewDungeon(tt.dungeonName, tt.floors, tt.seed)
			afterTime := time.Now()

			// Check basic properties
			assert.Equal(t, tt.dungeonName, dungeon.Name, "Dungeon name should match")
			assert.Equal(t, tt.floors, dungeon.Floors, "Floor count should match")
			assert.NotEmpty(t, dungeon.ID, "Dungeon ID should be generated")

			// Check seed
			if tt.seed != 0 {
				assert.Equal(t, tt.seed, dungeon.Seed, "Seed should match provided value")
			} else {
				assert.NotZero(t, dungeon.Seed, "Seed should be generated when not provided")
			}

			// Check creation time
			assert.True(t, dungeon.CreatedAt.After(beforeTime) || dungeon.CreatedAt.Equal(beforeTime),
				"Creation time should be after or equal to before time")
			assert.True(t, dungeon.CreatedAt.Before(afterTime) || dungeon.CreatedAt.Equal(afterTime),
				"Creation time should be before or equal to after time")

			// Check maps are initialized
			assert.NotNil(t, dungeon.FloorData, "Floor data map should be initialized")
			assert.Empty(t, dungeon.FloorData, "Floor data should start empty")
			assert.NotNil(t, dungeon.Characters, "Characters map should be initialized")
			assert.Empty(t, dungeon.Characters, "Characters should start empty")
			assert.Zero(t, dungeon.PlayerCount, "Player count should start at 0")
		})
	}
}

func TestGenerateFloor(t *testing.T) {
	tests := []struct {
		name        string
		floorLevel  int
		expectTiles bool
	}{
		{
			name:        "First Floor",
			floorLevel:  1,
			expectTiles: true,
		},
		{
			name:        "Middle Floor",
			floorLevel:  5,
			expectTiles: true,
		},
		{
			name:        "Deep Floor",
			floorLevel:  10,
			expectTiles: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dungeon := NewDungeon("Test Dungeon", 10, 12345)
			floor := dungeon.GenerateFloor(tt.floorLevel)

			// Check floor properties
			assert.Equal(t, tt.floorLevel, floor.Level, "Floor level should match")
			assert.Greater(t, floor.Width, 0, "Width should be positive")
			assert.Greater(t, floor.Height, 0, "Height should be positive")

			// Check tiles are initialized
			assert.NotEmpty(t, floor.Tiles, "Tiles should be initialized")
			assert.Len(t, floor.Tiles, floor.Height, "Tiles should have correct height")
			for y := 0; y < floor.Height; y++ {
				assert.Len(t, floor.Tiles[y], floor.Width, "Each row should have correct width")
				for x := 0; x < floor.Width; x++ {
					assert.Equal(t, TileWall, floor.Tiles[y][x].Type, "All tiles should start as walls")
					assert.False(t, floor.Tiles[y][x].Walkable, "Wall tiles should not be walkable")
					assert.False(t, floor.Tiles[y][x].Explored, "Tiles should start unexplored")
				}
			}

			// Check floor is stored in dungeon
			assert.Contains(t, dungeon.FloorData, tt.floorLevel, "Floor should be stored in dungeon")
			assert.Equal(t, floor, dungeon.FloorData[tt.floorLevel], "Stored floor should match generated floor")

			// Check maps are initialized
			assert.NotNil(t, floor.Mobs, "Mobs map should be initialized")
			assert.Empty(t, floor.Mobs, "Mobs should start empty")
			assert.NotNil(t, floor.Items, "Items map should be initialized")
			assert.Empty(t, floor.Items, "Items should start empty")
			assert.Empty(t, floor.Rooms, "Rooms should start empty")
			assert.Empty(t, floor.UpStairs, "Up stairs should start empty")
			assert.Empty(t, floor.DownStairs, "Down stairs should start empty")
		})
	}
}

func TestCharacterManagement(t *testing.T) {
	dungeon := NewDungeon("Test Dungeon", 5, 12345)

	// Test adding a character
	dungeon.AddCharacter("char1")
	assert.Equal(t, 1, dungeon.PlayerCount, "Player count should be 1 after adding a character")
	assert.Contains(t, dungeon.Characters, "char1", "Character should be in the map")
	assert.Equal(t, "1", dungeon.Characters["char1"], "Character should start on floor 1")

	// Test adding another character
	dungeon.AddCharacter("char2")
	assert.Equal(t, 2, dungeon.PlayerCount, "Player count should be 2 after adding another character")
	assert.Contains(t, dungeon.Characters, "char2", "Second character should be in the map")

	// Test getting character floor
	floor := dungeon.GetCharacterFloor("char1")
	assert.Equal(t, 1, floor, "Character should be on floor 1")

	// Test setting character floor
	dungeon.SetCharacterFloor("char1", 3)
	floor = dungeon.GetCharacterFloor("char1")
	assert.Equal(t, 3, floor, "Character should now be on floor 3")

	// Test getting non-existent character
	floor = dungeon.GetCharacterFloor("nonexistent")
	assert.Equal(t, 0, floor, "Non-existent character should return floor 0")

	// Test removing a character
	dungeon.RemoveCharacter("char1")
	assert.Equal(t, 1, dungeon.PlayerCount, "Player count should be 1 after removing a character")
	assert.NotContains(t, dungeon.Characters, "char1", "Removed character should not be in the map")

	// Test removing another character
	dungeon.RemoveCharacter("char2")
	assert.Equal(t, 0, dungeon.PlayerCount, "Player count should be 0 after removing all characters")
	assert.Empty(t, dungeon.Characters, "Characters map should be empty")
}
