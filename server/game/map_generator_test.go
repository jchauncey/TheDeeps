package game

import (
	"fmt"
	"testing"

	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/stretchr/testify/assert"
)

func TestNewMapGenerator(t *testing.T) {
	// Test with a specific seed
	seed := int64(12345)
	generator := NewMapGenerator(seed)

	// Ensure the generator is created with the correct seed
	assert.NotNil(t, generator)
	assert.NotNil(t, generator.rng)
}

func TestGenerateFloor(t *testing.T) {
	// Create a map generator with a fixed seed for reproducibility
	generator := NewMapGenerator(42)

	// Test cases for different floor levels and final floor status
	testCases := []struct {
		name         string
		level        int
		isFinalFloor bool
		width        int
		height       int
	}{
		{"First Floor", 1, false, 80, 40},
		{"Middle Floor", 5, false, 80, 40},
		{"Final Floor", 10, true, 80, 40},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new floor
			floor := &models.Floor{
				Level:  tc.level,
				Width:  tc.width,
				Height: tc.height,
				Tiles:  make([][]models.Tile, tc.height),
			}

			// Initialize tiles
			for y := 0; y < tc.height; y++ {
				floor.Tiles[y] = make([]models.Tile, tc.width)
			}

			// Generate the floor
			generator.GenerateFloor(floor, tc.level, tc.isFinalFloor)

			// Verify the floor was generated correctly
			assert.Equal(t, tc.level, floor.Level)
			assert.Equal(t, tc.width, floor.Width)
			assert.Equal(t, tc.height, floor.Height)

			// Check that rooms were generated
			assert.Greater(t, len(floor.Rooms), 0)

			// Check that tiles were modified (not all walls)
			hasFloorTile := false
			for y := 0; y < tc.height; y++ {
				for x := 0; x < tc.width; x++ {
					if floor.Tiles[y][x].Type == models.TileFloor {
						hasFloorTile = true
						break
					}
				}
				if hasFloorTile {
					break
				}
			}
			assert.True(t, hasFloorTile, "No floor tiles were generated")

			// Check for stairs
			if tc.level > 1 {
				assert.Greater(t, len(floor.UpStairs), 0, "No up stairs on level > 1")
			}

			if !tc.isFinalFloor {
				assert.Greater(t, len(floor.DownStairs), 0, "No down stairs on non-final floor")
			} else {
				assert.Equal(t, 0, len(floor.DownStairs), "Down stairs should not exist on final floor")
			}

			// Check for mobs and items
			assert.NotNil(t, floor.Mobs)
			assert.NotNil(t, floor.Items)
		})
	}
}

func TestGenerateRooms(t *testing.T) {
	// Create a map generator with a fixed seed
	generator := NewMapGenerator(42)

	// Create a test floor
	width, height := 80, 40
	floor := &models.Floor{
		Width:  width,
		Height: height,
		Tiles:  make([][]models.Tile, height),
	}

	// Initialize tiles
	for y := 0; y < height; y++ {
		floor.Tiles[y] = make([]models.Tile, width)
		for x := 0; x < width; x++ {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileWall,
				Walkable: false,
			}
		}
	}

	// Test cases
	testCases := []struct {
		name         string
		numRooms     int
		level        int
		isFinalFloor bool
		minExpected  int // Minimum number of rooms expected (may be less than numRooms due to overlap)
	}{
		{"Few Rooms", 5, 1, false, 3},
		{"Many Rooms", 15, 5, false, 8},
		{"Final Floor", 10, 10, true, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rooms := generator.generateRooms(floor, tc.numRooms, tc.level, tc.isFinalFloor)

			// Check that some rooms were generated
			assert.GreaterOrEqual(t, len(rooms), tc.minExpected)

			// Check room properties
			for _, room := range rooms {
				assert.NotEmpty(t, room.ID)
				assert.GreaterOrEqual(t, room.Width, 5)
				assert.GreaterOrEqual(t, room.Height, 5)
				assert.GreaterOrEqual(t, room.X, 1)
				assert.GreaterOrEqual(t, room.Y, 1)
				assert.Less(t, room.X+room.Width, width-1)
				assert.Less(t, room.Y+room.Height, height-1)

				// Check that room tiles are floor tiles
				for y := room.Y; y < room.Y+room.Height; y++ {
					for x := room.X; x < room.X+room.Width; x++ {
						assert.Equal(t, models.TileFloor, floor.Tiles[y][x].Type)
						assert.True(t, floor.Tiles[y][x].Walkable)
						assert.Equal(t, room.ID, floor.Tiles[y][x].RoomID)
					}
				}
			}

			// Check for boss room on final floor
			if tc.isFinalFloor {
				hasBossRoom := false
				for _, room := range rooms {
					if room.Type == models.RoomBoss {
						hasBossRoom = true
						break
					}
				}
				assert.True(t, hasBossRoom, "Final floor should have a boss room")
			}
		})
	}
}

func TestConnectRooms(t *testing.T) {
	// Create a map generator with a fixed seed
	generator := NewMapGenerator(42)

	// Create a test floor
	width, height := 80, 40
	floor := &models.Floor{
		Width:  width,
		Height: height,
		Tiles:  make([][]models.Tile, height),
	}

	// Initialize tiles with walls
	for y := 0; y < height; y++ {
		floor.Tiles[y] = make([]models.Tile, width)
		for x := 0; x < width; x++ {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileWall,
				Walkable: false,
			}
		}
	}

	// Create some test rooms
	rooms := []models.Room{
		{ID: "room1", X: 10, Y: 10, Width: 8, Height: 8},
		{ID: "room2", X: 30, Y: 15, Width: 10, Height: 6},
		{ID: "room3", X: 50, Y: 25, Width: 7, Height: 9},
	}

	// Carve out the rooms
	for _, room := range rooms {
		for y := room.Y; y < room.Y+room.Height; y++ {
			for x := room.X; x < room.X+room.Width; x++ {
				floor.Tiles[y][x] = models.Tile{
					Type:     models.TileFloor,
					Walkable: true,
					RoomID:   room.ID,
				}
			}
		}
	}

	// Connect the rooms
	generator.connectRooms(floor, rooms)

	// Helper function to check if two points are connected by floor tiles
	isConnected := func(x1, y1, x2, y2 int) bool {
		// Simple BFS to check connectivity
		visited := make(map[string]bool)
		queue := []struct{ x, y int }{{x1, y1}}

		for len(queue) > 0 {
			pos := queue[0]
			queue = queue[1:]

			if pos.x == x2 && pos.y == y2 {
				return true
			}

			key := fmt.Sprintf("%d,%d", pos.x, pos.y)
			if visited[key] {
				continue
			}
			visited[key] = true

			// Check adjacent tiles
			directions := []struct{ dx, dy int }{
				{-1, 0}, {1, 0}, {0, -1}, {0, 1},
			}

			for _, dir := range directions {
				nx, ny := pos.x+dir.dx, pos.y+dir.dy
				if nx >= 0 && nx < width && ny >= 0 && ny < height &&
					floor.Tiles[ny][nx].Type == models.TileFloor {
					queue = append(queue, struct{ x, y int }{nx, ny})
				}
			}
		}

		return false
	}

	// Check that each room is connected to the next
	for i := 0; i < len(rooms)-1; i++ {
		room1 := rooms[i]
		room2 := rooms[i+1]

		// Get center points of rooms
		x1 := room1.X + room1.Width/2
		y1 := room1.Y + room1.Height/2
		x2 := room2.X + room2.Width/2
		y2 := room2.Y + room2.Height/2

		assert.True(t, isConnected(x1, y1, x2, y2),
			"Room %s should be connected to room %s", room1.ID, room2.ID)
	}
}

func TestCreateCorridors(t *testing.T) {
	// Create a map generator
	generator := NewMapGenerator(42)

	// Create a test floor
	width, height := 50, 50
	floor := &models.Floor{
		Width:  width,
		Height: height,
		Tiles:  make([][]models.Tile, height),
	}

	// Initialize tiles with walls
	for y := 0; y < height; y++ {
		floor.Tiles[y] = make([]models.Tile, width)
		for x := 0; x < width; x++ {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileWall,
				Walkable: false,
			}
		}
	}

	// Test horizontal corridor
	x1, x2, y := 10, 30, 25
	generator.createHorizontalCorridor(floor, x1, x2, y)

	// Check that the horizontal corridor was created
	for x := min(x1, x2); x <= max(x1, x2); x++ {
		assert.Equal(t, models.TileFloor, floor.Tiles[y][x].Type)
		assert.True(t, floor.Tiles[y][x].Walkable)
	}

	// Test vertical corridor
	y1, y2, x := 10, 40, 15
	generator.createVerticalCorridor(floor, y1, y2, x)

	// Check that the vertical corridor was created
	for y := min(y1, y2); y <= max(y1, y2); y++ {
		assert.Equal(t, models.TileFloor, floor.Tiles[y][x].Type)
		assert.True(t, floor.Tiles[y][x].Walkable)
	}
}

func TestPlaceStairs(t *testing.T) {
	// Create a map generator
	generator := NewMapGenerator(42)

	// Test cases
	testCases := []struct {
		name         string
		level        int
		isFinalFloor bool
		expectUp     bool
		expectDown   bool
	}{
		{"First Floor", 1, false, false, true},
		{"Middle Floor", 5, false, true, true},
		{"Final Floor", 10, true, true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test floor
			width, height := 80, 40
			floor := &models.Floor{
				Level:      tc.level,
				Width:      width,
				Height:     height,
				Tiles:      make([][]models.Tile, height),
				UpStairs:   []models.Position{},
				DownStairs: []models.Position{},
			}

			// Initialize tiles
			for y := 0; y < height; y++ {
				floor.Tiles[y] = make([]models.Tile, width)
				for x := 0; x < width; x++ {
					floor.Tiles[y][x] = models.Tile{
						Type:     models.TileWall,
						Walkable: false,
					}
				}
			}

			// Create some test rooms
			rooms := []models.Room{
				{ID: "room1", X: 10, Y: 10, Width: 8, Height: 8},
				{ID: "room2", X: 30, Y: 15, Width: 10, Height: 6},
			}

			// Carve out the rooms
			for _, room := range rooms {
				for y := room.Y; y < room.Y+room.Height; y++ {
					for x := room.X; x < room.X+room.Width; x++ {
						floor.Tiles[y][x] = models.Tile{
							Type:     models.TileFloor,
							Walkable: true,
							RoomID:   room.ID,
						}
					}
				}
			}

			// Place stairs
			generator.placeStairs(floor, rooms, tc.level, tc.isFinalFloor)

			// Check for up stairs
			hasUpStairs := len(floor.UpStairs) > 0
			assert.Equal(t, tc.expectUp, hasUpStairs)

			if hasUpStairs {
				pos := floor.UpStairs[0]
				assert.Equal(t, models.TileUpStairs, floor.Tiles[pos.Y][pos.X].Type)
				assert.True(t, floor.Tiles[pos.Y][pos.X].Walkable)
			}

			// Check for down stairs
			hasDownStairs := len(floor.DownStairs) > 0
			assert.Equal(t, tc.expectDown, hasDownStairs)

			if hasDownStairs {
				pos := floor.DownStairs[0]
				assert.Equal(t, models.TileDownStairs, floor.Tiles[pos.Y][pos.X].Type)
				assert.True(t, floor.Tiles[pos.Y][pos.X].Walkable)
			}
		})
	}
}

func TestPlaceMobs(t *testing.T) {
	// Create a map generator
	generator := NewMapGenerator(42)

	// Test cases
	testCases := []struct {
		name         string
		level        int
		isFinalFloor bool
		roomTypes    []models.RoomType
	}{
		{
			"Standard Floor",
			3,
			false,
			[]models.RoomType{models.RoomStandard, models.RoomStandard, models.RoomTreasure},
		},
		{
			"Boss Floor",
			10,
			true,
			[]models.RoomType{models.RoomBoss, models.RoomStandard, models.RoomSafe},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test floor
			width, height := 80, 40
			floor := &models.Floor{
				Level:  tc.level,
				Width:  width,
				Height: height,
				Tiles:  make([][]models.Tile, height),
			}

			// Initialize tiles
			for y := 0; y < height; y++ {
				floor.Tiles[y] = make([]models.Tile, width)
				for x := 0; x < width; x++ {
					floor.Tiles[y][x] = models.Tile{
						Type:     models.TileWall,
						Walkable: false,
					}
				}
			}

			// Create rooms with specified types
			rooms := make([]models.Room, len(tc.roomTypes))
			for i, roomType := range tc.roomTypes {
				rooms[i] = models.Room{
					ID:     fmt.Sprintf("room%d", i+1),
					Type:   roomType,
					X:      10 + i*20,
					Y:      10,
					Width:  10,
					Height: 10,
				}

				// Carve out the room
				for y := rooms[i].Y; y < rooms[i].Y+rooms[i].Height; y++ {
					for x := rooms[i].X; x < rooms[i].X+rooms[i].Width; x++ {
						floor.Tiles[y][x] = models.Tile{
							Type:     models.TileFloor,
							Walkable: true,
							RoomID:   rooms[i].ID,
						}
					}
				}
			}

			// Place mobs
			generator.placeMobs(floor, rooms, tc.level, tc.isFinalFloor)

			// Check that mobs were created
			assert.NotNil(t, floor.Mobs)

			// Check that mobs are placed in appropriate rooms
			for _, mob := range floor.Mobs {
				// Get the tile where the mob is placed
				tile := floor.Tiles[mob.Position.Y][mob.Position.X]

				// Find the room this mob is in
				var room models.Room
				for _, r := range rooms {
					if r.ID == tile.RoomID {
						room = r
						break
					}
				}

				// Skip if we couldn't find the room (shouldn't happen)
				if room.ID == "" {
					continue
				}

				// Check mob placement rules
				switch room.Type {
				case models.RoomSafe:
					// Safe rooms should not have mobs
					assert.Fail(t, "Safe room should not have mobs")
				case models.RoomBoss:
					if tc.isFinalFloor {
						// Boss room on final floor should have a boss mob
						assert.Equal(t, models.MobVariant("boss"), mob.Variant)
					}
				}

				// Check that the mob is properly referenced in the tile
				assert.Equal(t, mob.ID, tile.MobID)
			}
		})
	}
}

func TestPlaceItems(t *testing.T) {
	// Create a map generator
	generator := NewMapGenerator(42)

	// Test cases
	testCases := []struct {
		name      string
		level     int
		roomTypes []models.RoomType
	}{
		{
			"Standard Floor",
			3,
			[]models.RoomType{models.RoomStandard, models.RoomTreasure, models.RoomBoss},
		},
		{
			"Deep Floor",
			8,
			[]models.RoomType{models.RoomStandard, models.RoomShop, models.RoomSafe},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test floor
			width, height := 80, 40
			floor := &models.Floor{
				Level:  tc.level,
				Width:  width,
				Height: height,
				Tiles:  make([][]models.Tile, height),
			}

			// Initialize tiles
			for y := 0; y < height; y++ {
				floor.Tiles[y] = make([]models.Tile, width)
				for x := 0; x < width; x++ {
					floor.Tiles[y][x] = models.Tile{
						Type:     models.TileWall,
						Walkable: false,
					}
				}
			}

			// Create rooms with specified types
			rooms := make([]models.Room, len(tc.roomTypes))
			for i, roomType := range tc.roomTypes {
				rooms[i] = models.Room{
					ID:     fmt.Sprintf("room%d", i+1),
					Type:   roomType,
					X:      10 + i*20,
					Y:      10,
					Width:  10,
					Height: 10,
				}

				// Carve out the room
				for y := rooms[i].Y; y < rooms[i].Y+rooms[i].Height; y++ {
					for x := rooms[i].X; x < rooms[i].X+rooms[i].Width; x++ {
						floor.Tiles[y][x] = models.Tile{
							Type:     models.TileFloor,
							Walkable: true,
							RoomID:   rooms[i].ID,
						}
					}
				}
			}

			// Place items
			generator.placeItems(floor, rooms, tc.level)

			// Check that items were created
			assert.NotNil(t, floor.Items)

			// Check item placement rules
			treasureRoomItems := 0
			standardRoomItems := 0
			bossRoomItems := 0

			for itemID := range floor.Items {
				// Find the tile where the item is placed
				var itemTile models.Tile
				var itemPos models.Position

				// Find the position of the item
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						if floor.Tiles[y][x].ItemID == itemID {
							itemTile = floor.Tiles[y][x]
							itemPos = models.Position{X: x, Y: y}
							break
						}
					}
				}

				// Find the room this item is in
				var roomType models.RoomType
				for _, room := range rooms {
					if room.ID == itemTile.RoomID {
						roomType = room.Type
						break
					}
				}

				// Count items by room type
				switch roomType {
				case models.RoomTreasure:
					treasureRoomItems++
				case models.RoomStandard:
					standardRoomItems++
				case models.RoomBoss:
					bossRoomItems++
				}

				// Check that the item is on a walkable tile
				assert.True(t, floor.Tiles[itemPos.Y][itemPos.X].Walkable)

				// Check that the item is properly referenced in the tile
				assert.Equal(t, itemID, floor.Tiles[itemPos.Y][itemPos.X].ItemID)
			}

			// Treasure rooms should have more items
			if treasureRoomItems > 0 {
				assert.GreaterOrEqual(t, treasureRoomItems, 3, "Treasure rooms should have at least 3 items")
			}

			// Boss rooms should have some items
			if bossRoomItems > 0 {
				assert.GreaterOrEqual(t, bossRoomItems, 2, "Boss rooms should have at least 2 items")
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test min function
	assert.Equal(t, 5, min(5, 10))
	assert.Equal(t, 5, min(10, 5))
	assert.Equal(t, -5, min(-5, 5))
	assert.Equal(t, -10, min(-5, -10))

	// Test max function
	assert.Equal(t, 10, max(5, 10))
	assert.Equal(t, 10, max(10, 5))
	assert.Equal(t, 5, max(-5, 5))
	assert.Equal(t, -5, max(-5, -10))
}

// TestGenerateFloorWithDifficulty tests the GenerateFloorWithDifficulty function
func TestGenerateFloorWithDifficulty(t *testing.T) {
	// Create a map generator with a fixed seed for deterministic tests
	generator := NewMapGenerator(12345)

	tests := []struct {
		name         string
		level        int
		isFinalFloor bool
		difficulty   string
	}{
		{
			name:         "Easy Difficulty",
			level:        3,
			isFinalFloor: false,
			difficulty:   "easy",
		},
		{
			name:         "Normal Difficulty",
			level:        3,
			isFinalFloor: false,
			difficulty:   "normal",
		},
		{
			name:         "Hard Difficulty",
			level:        3,
			isFinalFloor: false,
			difficulty:   "hard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a floor
			floor := &models.Floor{
				Level:  tt.level,
				Width:  50,
				Height: 50,
				Tiles:  make([][]models.Tile, 50),
			}

			// Initialize tiles
			for y := range floor.Tiles {
				floor.Tiles[y] = make([]models.Tile, 50)
			}

			// Generate the floor
			generator.GenerateFloorWithDifficulty(floor, tt.level, tt.isFinalFloor, tt.difficulty)

			// Verify that the floor was generated
			assert.NotEmpty(t, floor.Rooms, "Rooms should be generated")
			assert.NotEmpty(t, floor.Mobs, "Mobs should be generated")

			// Count mobs by variant
			easyCount := 0
			normalCount := 0
			hardCount := 0

			for _, mob := range floor.Mobs {
				switch mob.Variant {
				case models.VariantEasy:
					easyCount++
				case models.VariantNormal:
					normalCount++
				case models.VariantHard:
					hardCount++
				}
			}

			// Verify that the difficulty affects mob variants
			switch tt.difficulty {
			case "easy":
				// Easy difficulty should have more easy mobs
				assert.True(t, easyCount > normalCount, "Easy difficulty should have more easy mobs than normal mobs")
				assert.True(t, normalCount > hardCount, "Easy difficulty should have more normal mobs than hard mobs")
			case "hard":
				// Hard difficulty should have more hard mobs
				assert.True(t, hardCount > 0, "Hard difficulty should have at least some hard mobs")
			}
		})
	}
}
