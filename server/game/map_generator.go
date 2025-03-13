package game

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/jchauncey/TheDeeps/server/models"
)

// MapGenerator handles the procedural generation of dungeon maps
type MapGenerator struct {
	rng *rand.Rand
}

// NewMapGenerator creates a new map generator with the given seed
func NewMapGenerator(seed int64) *MapGenerator {
	return &MapGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// GenerateFloor generates a complete floor for a dungeon
func (g *MapGenerator) GenerateFloor(floor *models.Floor, level int, isFinalFloor bool) {
	// Initialize the floor with walls
	for y := 0; y < floor.Height; y++ {
		for x := 0; x < floor.Width; x++ {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileWall,
				Walkable: false,
				Explored: false,
			}
		}
	}

	// Determine number of rooms based on floor level
	minRooms := 5
	maxRooms := 10 + level
	if maxRooms > 20 {
		maxRooms = 20
	}

	numRooms := minRooms + g.rng.Intn(maxRooms-minRooms+1)

	// Generate rooms
	rooms := g.generateRooms(floor, numRooms, level, isFinalFloor)
	floor.Rooms = rooms

	// Connect rooms with corridors
	g.connectRooms(floor, rooms)

	// Place stairs
	g.placeStairs(floor, rooms, level, isFinalFloor)

	// Place mobs
	g.placeMobs(floor, rooms, level, isFinalFloor)

	// Place items
	g.placeItems(floor, rooms, level)
}

// GenerateFloorWithDifficulty generates a complete floor for a dungeon with the specified difficulty
func (g *MapGenerator) GenerateFloorWithDifficulty(floor *models.Floor, level int, isFinalFloor bool, difficulty string) {
	// Initialize the floor with walls
	for y := 0; y < floor.Height; y++ {
		for x := 0; x < floor.Width; x++ {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileWall,
				Walkable: false,
				Explored: false,
			}
		}
	}

	// Determine number of rooms based on floor level
	minRooms := 5
	maxRooms := 10 + level
	if maxRooms > 20 {
		maxRooms = 20
	}

	numRooms := minRooms + g.rng.Intn(maxRooms-minRooms+1)

	// Generate rooms
	rooms := g.generateRooms(floor, numRooms, level, isFinalFloor)
	floor.Rooms = rooms

	// Connect rooms with corridors
	g.connectRooms(floor, rooms)

	// Place stairs
	g.placeStairs(floor, rooms, level, isFinalFloor)

	// Place mobs with difficulty
	g.placeMobsWithDifficulty(floor, rooms, level, isFinalFloor, difficulty)

	// Place items
	g.placeItems(floor, rooms, level)
}

// generateRooms creates a set of rooms for the floor
func (g *MapGenerator) generateRooms(floor *models.Floor, numRooms int, level int, isFinalFloor bool) []models.Room {
	rooms := make([]models.Room, 0, numRooms)

	// Room size ranges
	minRoomSize := 5
	maxRoomSize := 10

	// First, create the entrance room on the first floor with a consistent size
	if level == 1 {
		// Entrance room has a consistent size of 8x8
		entranceWidth := 8
		entranceHeight := 8

		// Place the entrance room in a more central location
		// Calculate a position near the center of the map
		centerX := floor.Width / 2
		centerY := floor.Height / 2

		// Add some randomness but keep it near the center
		offsetX := -4 + g.rng.Intn(9) // -4 to +4
		offsetY := -4 + g.rng.Intn(9) // -4 to +4

		entranceX := centerX + offsetX - entranceWidth/2
		entranceY := centerY + offsetY - entranceHeight/2

		// Ensure the room is within bounds
		if entranceX < 1 {
			entranceX = 1
		}
		if entranceY < 1 {
			entranceY = 1
		}
		if entranceX+entranceWidth >= floor.Width-1 {
			entranceX = floor.Width - entranceWidth - 2
		}
		if entranceY+entranceHeight >= floor.Height-1 {
			entranceY = floor.Height - entranceHeight - 2
		}

		// Create the entrance room
		entranceRoom := models.Room{
			ID:       uuid.New().String(),
			Type:     models.RoomEntrance,
			X:        entranceX,
			Y:        entranceY,
			Width:    entranceWidth,
			Height:   entranceHeight,
			Explored: true, // Entrance room starts explored
		}

		// Carve out the entrance room
		for ry := 0; ry < entranceHeight; ry++ {
			for rx := 0; rx < entranceWidth; rx++ {
				floor.Tiles[entranceY+ry][entranceX+rx] = models.Tile{
					Type:     models.TileFloor,
					Walkable: true,
					Explored: true, // Entrance room tiles start explored
					RoomID:   entranceRoom.ID,
				}
			}
		}

		rooms = append(rooms, entranceRoom)
	}

	// Try to place the rest of the rooms
	for i := 0; i < numRooms*3 && len(rooms) < numRooms; i++ {
		// Random room size
		width := minRoomSize + g.rng.Intn(maxRoomSize-minRoomSize+1)
		height := minRoomSize + g.rng.Intn(maxRoomSize-minRoomSize+1)

		// Random position (leaving a 1-tile border)
		x := 1 + g.rng.Intn(floor.Width-width-2)
		y := 1 + g.rng.Intn(floor.Height-height-2)

		// Check if the room overlaps with existing rooms
		overlaps := false
		for _, room := range rooms {
			if x+width > room.X-2 && x < room.X+room.Width+2 &&
				y+height > room.Y-2 && y < room.Y+room.Height+2 {
				overlaps = true
				break
			}
		}

		if !overlaps {
			// Determine room type
			roomType := models.RoomStandard

			// Special rooms
			if level == 1 && len(rooms) == 0 {
				// Skip this iteration if we're on the first floor and already created the entrance room
				continue
			} else if isFinalFloor && len(rooms) == 0 {
				// First room on final floor is the boss room
				roomType = models.RoomBoss
			} else if g.rng.Float64() < 0.1 {
				// 10% chance for treasure room
				roomType = models.RoomTreasure
			} else if g.rng.Float64() < 0.05 {
				// 5% chance for safe room
				roomType = models.RoomSafe
			} else if g.rng.Float64() < 0.05 && level > 2 {
				// 5% chance for shop room on deeper floors
				roomType = models.RoomShop
			}

			// Create the room
			room := models.Room{
				ID:       uuid.New().String(),
				Type:     roomType,
				X:        x,
				Y:        y,
				Width:    width,
				Height:   height,
				Explored: false,
			}

			// Carve out the room
			for ry := 0; ry < height; ry++ {
				for rx := 0; rx < width; rx++ {
					floor.Tiles[y+ry][x+rx] = models.Tile{
						Type:     models.TileFloor,
						Walkable: true,
						Explored: false,
						RoomID:   room.ID,
					}
				}
			}

			rooms = append(rooms, room)
		}
	}

	return rooms
}

// connectRooms connects rooms with corridors
func (g *MapGenerator) connectRooms(floor *models.Floor, rooms []models.Room) {
	// Connect each room to the next one
	for i := 0; i < len(rooms)-1; i++ {
		// Get center points of rooms
		x1 := rooms[i].X + rooms[i].Width/2
		y1 := rooms[i].Y + rooms[i].Height/2
		x2 := rooms[i+1].X + rooms[i+1].Width/2
		y2 := rooms[i+1].Y + rooms[i+1].Height/2

		// Randomly decide whether to go horizontal-then-vertical or vertical-then-horizontal
		if g.rng.Intn(2) == 0 {
			// Horizontal then vertical
			g.createHorizontalCorridor(floor, x1, x2, y1)
			g.createVerticalCorridor(floor, y1, y2, x2)
		} else {
			// Vertical then horizontal
			g.createVerticalCorridor(floor, y1, y2, x1)
			g.createHorizontalCorridor(floor, x1, x2, y2)
		}
	}
}

// createHorizontalCorridor creates a horizontal corridor
func (g *MapGenerator) createHorizontalCorridor(floor *models.Floor, x1, x2, y int) {
	for x := min(x1, x2); x <= max(x1, x2); x++ {
		if x >= 0 && x < floor.Width && y >= 0 && y < floor.Height {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileFloor,
				Walkable: true,
				Explored: false,
			}
		}
	}
}

// createVerticalCorridor creates a vertical corridor
func (g *MapGenerator) createVerticalCorridor(floor *models.Floor, y1, y2, x int) {
	for y := min(y1, y2); y <= max(y1, y2); y++ {
		if x >= 0 && x < floor.Width && y >= 0 && y < floor.Height {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileFloor,
				Walkable: true,
				Explored: false,
			}
		}
	}
}

// placeStairs places up and down stairs on the floor
func (g *MapGenerator) placeStairs(floor *models.Floor, rooms []models.Room, level int, isFinalFloor bool) {
	// Find entrance room for first floor
	var entranceRoom *models.Room
	if level == 1 {
		for i := range rooms {
			if rooms[i].Type == models.RoomEntrance {
				entranceRoom = &rooms[i]
				break
			}
		}
	}

	// Place up stairs in the first room (except for the first floor)
	if level > 1 {
		room := rooms[0]
		x := room.X + g.rng.Intn(room.Width)
		y := room.Y + g.rng.Intn(room.Height)

		floor.Tiles[y][x] = models.Tile{
			Type:     models.TileUpStairs,
			Walkable: true,
			Explored: false,
			RoomID:   room.ID,
		}

		floor.UpStairs = append(floor.UpStairs, models.Position{X: x, Y: y})
	}

	// Place down stairs in the entrance room if it's the first floor, otherwise in the last room
	// (except for the final floor which has no down stairs)
	if !isFinalFloor {
		var room models.Room
		var x, y int

		if level == 1 && entranceRoom != nil {
			room = *entranceRoom
			// Place down stairs in a consistent location in the entrance room
			// Bottom right corner, but not at the very edge
			x = room.X + room.Width - 2
			y = room.Y + room.Height - 2
		} else {
			room = rooms[len(rooms)-1]
			// Random position for other rooms
			x = room.X + g.rng.Intn(room.Width)
			y = room.Y + g.rng.Intn(room.Height)
		}

		floor.Tiles[y][x] = models.Tile{
			Type:     models.TileDownStairs,
			Walkable: true,
			Explored: level == 1 && entranceRoom != nil, // Explored if in entrance room
			RoomID:   room.ID,
		}

		floor.DownStairs = append(floor.DownStairs, models.Position{X: x, Y: y})
	}
}

// placeMobs places mobs on the floor
func (g *MapGenerator) placeMobs(floor *models.Floor, rooms []models.Room, level int, isFinalFloor bool) {
	// Default to normal difficulty
	g.placeMobsWithDifficulty(floor, rooms, level, isFinalFloor, "normal")
}

// placeMobsWithDifficulty places mobs on the floor with the specified difficulty
func (g *MapGenerator) placeMobsWithDifficulty(floor *models.Floor, rooms []models.Room, level int, isFinalFloor bool, difficulty string) {
	floor.Mobs = make(map[string]*models.Mob)

	// Determine mob types based on floor level
	mobTypes := []models.MobType{models.MobSkeleton, models.MobGoblin, models.MobRatman}

	if level >= 3 {
		mobTypes = append(mobTypes, models.MobOrc, models.MobOoze)
	}

	if level >= 5 {
		mobTypes = append(mobTypes, models.MobTroll, models.MobWraith)
	}

	if level >= 8 {
		mobTypes = append(mobTypes, models.MobOgre, models.MobDrake)
	}

	if level >= 10 {
		mobTypes = append(mobTypes, models.MobLich, models.MobElemental)
	}

	// Place mobs in each room
	for i, room := range rooms {
		// Skip the first room (where the player starts), safe rooms, and entrance rooms
		if i == 0 || room.Type == models.RoomSafe || room.Type == models.RoomEntrance {
			continue
		}

		// Boss room gets a boss
		if room.Type == models.RoomBoss {
			// Create a boss mob
			mobType := models.MobDragon
			if level < 10 {
				mobType = models.MobOgre
			}

			mob := models.NewMob(mobType, models.VariantBoss, level)

			// Place in center of room
			x := room.X + room.Width/2
			y := room.Y + room.Height/2
			mob.Position = models.Position{X: x, Y: y}

			// Add to floor
			floor.Mobs[mob.ID] = mob
			floor.Tiles[y][x].MobID = mob.ID

			continue
		}

		// Regular rooms get random mobs
		numMobs := 1 + g.rng.Intn(3) // 1-3 mobs per room

		// Adjust based on room size
		roomArea := room.Width * room.Height
		if roomArea > 80 {
			numMobs += 2
		} else if roomArea > 50 {
			numMobs += 1
		}

		// Adjust based on floor level
		numMobs += level / 3

		// Adjust based on difficulty
		if difficulty == "hard" {
			numMobs += 1
		} else if difficulty == "easy" {
			numMobs = max(1, numMobs-1)
		}

		// Cap at reasonable number
		if numMobs > 8 {
			numMobs = 8
		}

		for j := 0; j < numMobs; j++ {
			// Pick a random mob type
			mobType := mobTypes[g.rng.Intn(len(mobTypes))]

			// Determine variant based on difficulty
			variant := models.VariantNormal
			variantRoll := g.rng.Float64()

			switch difficulty {
			case "easy":
				if variantRoll < 0.8 {
					variant = models.VariantEasy
				} else {
					variant = models.VariantNormal
				}
			case "normal":
				if variantRoll < 0.6 {
					variant = models.VariantEasy
				} else if variantRoll < 0.9 {
					variant = models.VariantNormal
				} else {
					variant = models.VariantHard
				}
			case "hard":
				if variantRoll < 0.3 {
					variant = models.VariantEasy
				} else if variantRoll < 0.7 {
					variant = models.VariantNormal
				} else {
					variant = models.VariantHard
				}
			default:
				// Default to normal difficulty distribution
				if variantRoll < 0.6 {
					variant = models.VariantEasy
				} else if variantRoll < 0.9 {
					variant = models.VariantNormal
				} else {
					variant = models.VariantHard
				}
			}

			// Create the mob
			mob := models.NewMob(mobType, variant, level)

			// Find a valid position
			var x, y int
			for {
				x = room.X + g.rng.Intn(room.Width)
				y = room.Y + g.rng.Intn(room.Height)

				// Check if the tile is walkable and doesn't have a mob
				if floor.Tiles[y][x].Walkable && floor.Tiles[y][x].MobID == "" &&
					floor.Tiles[y][x].Type != models.TileUpStairs &&
					floor.Tiles[y][x].Type != models.TileDownStairs {
					break
				}
			}

			mob.Position = models.Position{X: x, Y: y}

			// Add to floor
			floor.Mobs[mob.ID] = mob
			floor.Tiles[y][x].MobID = mob.ID
		}
	}
}

// placeItems places items on the floor
func (g *MapGenerator) placeItems(floor *models.Floor, rooms []models.Room, level int) {
	floor.Items = make(map[string]models.Item)

	// Place items in each room
	for _, room := range rooms {
		// Determine number of items based on room type
		numItems := 0

		switch room.Type {
		case models.RoomTreasure:
			numItems = 3 + g.rng.Intn(3) // 3-5 items
		case models.RoomStandard:
			if g.rng.Float64() < 0.3 { // 30% chance for an item
				numItems = 1
			}
		case models.RoomBoss:
			numItems = 2 + g.rng.Intn(3) // 2-4 items
		case models.RoomEntrance:
			if level == 1 {
				// First floor entrance might have a starter item
				if g.rng.Float64() < 0.5 { // 50% chance for a starter item
					numItems = 1
				}
			}
		}

		// Skip if no items
		if numItems == 0 {
			continue
		}

		for j := 0; j < numItems; j++ {
			// Generate a random item
			item := models.GenerateRandomItem(level)

			// Find a valid position
			var x, y int
			for {
				x = room.X + g.rng.Intn(room.Width)
				y = room.Y + g.rng.Intn(room.Height)

				// Check if the tile is walkable and doesn't have an item or mob
				if floor.Tiles[y][x].Walkable && floor.Tiles[y][x].ItemID == "" &&
					floor.Tiles[y][x].MobID == "" &&
					floor.Tiles[y][x].Type != models.TileUpStairs &&
					floor.Tiles[y][x].Type != models.TileDownStairs {
					break
				}
			}

			item.Position = models.Position{X: x, Y: y}

			// Add to floor
			floor.Items[item.ID] = *item
			floor.Tiles[y][x].ItemID = item.ID
		}
	}
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
