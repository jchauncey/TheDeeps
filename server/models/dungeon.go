package models

import (
	"fmt"
	"math/rand"
	"time"
)

// TileType represents the type of a dungeon tile
type TileType string

const (
	// TileWall represents a wall tile
	TileWall TileType = "wall"
	// TileFloor represents a floor tile
	TileFloor TileType = "floor"
	// TileDoor represents a door tile
	TileDoor TileType = "door"
	// TileStairsUp represents stairs going up
	TileStairsUp TileType = "stairs_up"
	// TileStairsDown represents stairs going down
	TileStairsDown TileType = "stairs_down"
)

// Tile represents a single tile in the dungeon
type Tile struct {
	Type     TileType `json:"type"`
	Explored bool     `json:"explored"`
	Visible  bool     `json:"visible"`
	Entity   *Entity  `json:"entity,omitempty"`
	Item     *Item    `json:"item,omitempty"`
}

// Position represents a 2D position in the dungeon
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Room represents a room in the dungeon
type Room struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Contains checks if a position is inside the room
func (r *Room) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}

// Center returns the center position of the room
func (r *Room) Center() (int, int) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

// Entity represents an entity in the dungeon (monster, NPC, etc.)
type Entity struct {
	ID             string   `json:"id"`
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	Position       Position `json:"position"`
	CharacterClass string   `json:"characterClass,omitempty"`
	Health         int      `json:"health,omitempty"`
	MaxHealth      int      `json:"maxHealth,omitempty"`
	Status         []string `json:"status,omitempty"`
	// Add more entity properties as needed
}

// Item represents an item in the dungeon
type Item struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Position Position `json:"position"`
	// Add more item properties as needed
}

// Floor represents a floor in the dungeon
type Floor struct {
	Level    int      `json:"level"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	Tiles    [][]Tile `json:"tiles"`
	Rooms    []Room   `json:"rooms"`
	Entities []Entity `json:"entities"`
	Items    []Item   `json:"items"`
}

// Dungeon represents the entire dungeon
type Dungeon struct {
	Floors         []*Floor `json:"floors"`
	CurrentFloor   int      `json:"currentFloor"`
	PlayerPosition Position `json:"playerPosition"`
}

// NewDungeon creates a new dungeon with the specified number of floors
func NewDungeon(numFloors int) *Dungeon {
	rand.Seed(time.Now().UnixNano())

	dungeon := &Dungeon{
		Floors:       make([]*Floor, numFloors),
		CurrentFloor: 0,
	}

	// Generate each floor
	for i := 0; i < numFloors; i++ {
		dungeon.Floors[i] = GenerateFloor(i+1, 80, 50)
	}

	// Set player position to the center of the first room on the first floor
	firstRoom := dungeon.Floors[0].Rooms[0]
	centerX, centerY := firstRoom.Center()
	dungeon.PlayerPosition = Position{X: centerX, Y: centerY}

	return dungeon
}

// GenerateFloor generates a new dungeon floor
func GenerateFloor(level, width, height int) *Floor {
	// Increase floor size based on level
	width = 100 + level*5 // Larger width, scaling with level
	height = 70 + level*3 // Larger height, scaling with level

	floor := &Floor{
		Level:  level,
		Width:  width,
		Height: height,
		Tiles:  make([][]Tile, height),
	}

	// Initialize all tiles as walls
	for y := 0; y < height; y++ {
		floor.Tiles[y] = make([]Tile, width)
		for x := 0; x < width; x++ {
			floor.Tiles[y][x] = Tile{Type: TileWall}
		}
	}

	// Generate rooms - more rooms for deeper levels
	minRooms := 8 + level  // Minimum rooms increases with level
	maxRooms := 15 + level // Maximum rooms increases with level
	floor.Rooms = generateRooms(floor, maxRooms, 6, 15, 2, 30)

	// If we didn't generate enough rooms, try again with more relaxed parameters
	if len(floor.Rooms) < minRooms {
		floor.Rooms = generateRooms(floor, maxRooms, 5, 12, 1, 40)
	}

	// Connect rooms with corridors
	connectRooms(floor)

	// Add stairs
	addStairs(floor, level)

	// Add entities and items - more on deeper levels
	floor.Entities = generateEntities(floor, level)
	floor.Items = generateItems(floor, level)

	return floor
}

// generateRooms generates random rooms on the floor
func generateRooms(floor *Floor, maxRooms, minSize, maxSize, minDistance, maxTries int) []Room {
	rooms := make([]Room, 0, maxRooms)

	for i := 0; i < maxRooms; i++ {
		// Try to place a room
		for try := 0; try < maxTries; try++ {
			// Random room size
			w := rand.Intn(maxSize-minSize+1) + minSize
			h := rand.Intn(maxSize-minSize+1) + minSize

			// Random room position
			x := rand.Intn(floor.Width-w-2) + 1
			y := rand.Intn(floor.Height-h-2) + 1

			newRoom := Room{X: x, Y: y, Width: w, Height: h}

			// Check if the room overlaps with existing rooms
			overlaps := false
			for _, room := range rooms {
				if roomsOverlap(&newRoom, &room, minDistance) {
					overlaps = true
					break
				}
			}

			if !overlaps {
				// Room doesn't overlap, add it
				rooms = append(rooms, newRoom)

				// Carve out the room
				for y := newRoom.Y; y < newRoom.Y+newRoom.Height; y++ {
					for x := newRoom.X; x < newRoom.X+newRoom.Width; x++ {
						floor.Tiles[y][x].Type = TileFloor
					}
				}

				break
			}
		}
	}

	return rooms
}

// roomsOverlap checks if two rooms overlap or are too close
func roomsOverlap(r1, r2 *Room, minDistance int) bool {
	return r1.X-minDistance < r2.X+r2.Width+minDistance &&
		r1.X+r1.Width+minDistance > r2.X-minDistance &&
		r1.Y-minDistance < r2.Y+r2.Height+minDistance &&
		r1.Y+r1.Height+minDistance > r2.Y-minDistance
}

// connectRooms connects rooms with corridors
func connectRooms(floor *Floor) {
	// First, connect rooms in sequence
	for i := 0; i < len(floor.Rooms)-1; i++ {
		// Connect each room to the next one
		room1 := floor.Rooms[i]
		room2 := floor.Rooms[i+1]

		// Get center points of each room
		x1, y1 := room1.Center()
		x2, y2 := room2.Center()

		// Determine which direction has the greater distance
		dx := abs(x2 - x1)
		dy := abs(y2 - y1)

		// Choose corridor creation strategy based on the room positions
		if dx > dy {
			// Rooms are more horizontally separated, so go horizontal first
			createHorizontalCorridor(floor, x1, x2, y1)
			createVerticalCorridor(floor, y1, y2, x2)
		} else {
			// Rooms are more vertically separated, so go vertical first
			createVerticalCorridor(floor, y1, y2, x1)
			createHorizontalCorridor(floor, x1, x2, y2)
		}
	}

	// Add some additional random connections for better connectivity
	// But limit to avoid too many redundant corridors
	if len(floor.Rooms) > 3 {
		numExtraConnections := len(floor.Rooms) / 4 // Reduce the number of extra connections

		// Keep track of connections to avoid duplicates
		connections := make(map[string]bool)

		for i := 0; i < numExtraConnections; i++ {
			// Try several times to find a valid connection
			for attempt := 0; attempt < 5; attempt++ {
				// Pick two random rooms that are not adjacent in the sequence
				r1 := rand.Intn(len(floor.Rooms))
				r2 := rand.Intn(len(floor.Rooms))

				// Make sure they're different rooms and not adjacent in the sequence
				if r1 == r2 || abs(r1-r2) == 1 {
					continue
				}

				// Create a unique key for this room pair
				connectionKey := ""
				if r1 < r2 {
					connectionKey = fmt.Sprintf("%d-%d", r1, r2)
				} else {
					connectionKey = fmt.Sprintf("%d-%d", r2, r1)
				}

				// Skip if we already have this connection
				if connections[connectionKey] {
					continue
				}

				room1 := floor.Rooms[r1]
				room2 := floor.Rooms[r2]

				// Get center points of each room
				x1, y1 := room1.Center()
				x2, y2 := room2.Center()

				// Determine which direction has the greater distance
				dx := abs(x2 - x1)
				dy := abs(y2 - y1)

				// Create a corridor between them based on their relative positions
				if dx > dy {
					createHorizontalCorridor(floor, x1, x2, y1)
					createVerticalCorridor(floor, y1, y2, x2)
				} else {
					createVerticalCorridor(floor, y1, y2, x1)
					createHorizontalCorridor(floor, x1, x2, y2)
				}

				// Mark this connection as created
				connections[connectionKey] = true
				break
			}
		}
	}
}

// createHorizontalCorridor creates a horizontal corridor
func createHorizontalCorridor(floor *Floor, x1, x2, y int) {
	for x := min(x1, x2); x <= max(x1, x2); x++ {
		if y >= 0 && y < floor.Height && x >= 0 && x < floor.Width {
			floor.Tiles[y][x].Type = TileFloor
		}
	}
}

// createVerticalCorridor creates a vertical corridor
func createVerticalCorridor(floor *Floor, y1, y2, x int) {
	for y := min(y1, y2); y <= max(y1, y2); y++ {
		if y >= 0 && y < floor.Height && x >= 0 && x < floor.Width {
			floor.Tiles[y][x].Type = TileFloor
		}
	}
}

// addStairs adds stairs to the floor
func addStairs(floor *Floor, level int) {
	// Add stairs up (except on first floor)
	if level > 1 {
		room := floor.Rooms[0]
		x, y := room.Center()
		floor.Tiles[y][x].Type = TileStairsUp
	}

	// Add stairs down (except on last floor)
	if level < 10 { // Assuming 10 floors total
		room := floor.Rooms[len(floor.Rooms)-1]
		x, y := room.Center()
		floor.Tiles[y][x].Type = TileStairsDown
	}
}

// generateEntities generates random entities for the floor
func generateEntities(floor *Floor, level int) []Entity {
	// Scale number of entities with floor level and size
	numEntities := 5 + level*2 + len(floor.Rooms)/2
	entities := make([]Entity, 0, numEntities)

	// Get available mob types for this floor level
	availableMobTypes := GetMobsForFloorLevel(level)
	if len(availableMobTypes) == 0 {
		// Fallback to basic mobs if none are available
		availableMobTypes = []MobType{MobRatman, MobGoblin, MobSkeleton}
	}

	// Generate entities
	for i := 0; i < numEntities; i++ {
		// Pick a random room (excluding the first room which is the player's starting point)
		roomIndex := 0
		if len(floor.Rooms) > 1 {
			roomIndex = 1 + rand.Intn(len(floor.Rooms)-1)
		}
		room := floor.Rooms[roomIndex]

		// Pick a random position within the room
		x := room.X + rand.Intn(room.Width)
		y := room.Y + rand.Intn(room.Height)

		// Pick a random mob type from available types
		mobTypeIndex := rand.Intn(len(availableMobTypes))
		mobType := availableMobTypes[mobTypeIndex]

		// Determine difficulty based on floor level
		difficulty := GetRandomDifficulty(level)

		// Create mob instance
		mobPosition := Position{X: x, Y: y}
		mobInstance := CreateMobInstance(mobType, difficulty, level, mobPosition)

		// Convert to Entity for the floor
		entity := Entity{
			ID:        mobInstance.ID,
			Type:      string(mobInstance.Type),
			Name:      mobInstance.Name,
			Position:  mobInstance.Position,
			Health:    mobInstance.Health,
			MaxHealth: mobInstance.MaxHealth,
			Status:    mobInstance.Status,
		}

		entities = append(entities, entity)
	}

	return entities
}

// generateItems generates random items for the floor
func generateItems(floor *Floor, level int) []Item {
	// Scale number of items with floor level and size
	numItems := 3 + level + len(floor.Rooms)/3
	items := make([]Item, 0, numItems)

	// Item types with weights (higher level = better items)
	itemTypes := []string{"potion", "scroll", "coin", "gem"}

	// Add better items on deeper levels
	if level > 2 {
		itemTypes = append(itemTypes, "weapon", "armor")
	}
	if level > 5 {
		itemTypes = append(itemTypes, "wand", "amulet")
	}

	// Generate items
	for i := 0; i < numItems; i++ {
		// Pick a random room (excluding the first room on the first level)
		roomIndex := 0
		if len(floor.Rooms) > 1 && (level > 1 || i > 0) {
			roomIndex = rand.Intn(len(floor.Rooms))
		}
		room := floor.Rooms[roomIndex]

		// Pick a random position within the room
		x := room.X + rand.Intn(room.Width)
		y := room.Y + rand.Intn(room.Height)

		// Pick a random item type (weighted toward better items on deeper levels)
		typeIndex := rand.Intn(len(itemTypes))
		if level > 3 && rand.Intn(10) < 4 {
			// 40% chance to pick from the second half of the list on deeper levels
			typeIndex = len(itemTypes)/2 + rand.Intn(len(itemTypes)/2+1)
			if typeIndex >= len(itemTypes) {
				typeIndex = len(itemTypes) - 1
			}
		}
		itemType := itemTypes[typeIndex]

		// Create the item
		item := Item{
			ID:   generateID(),
			Type: itemType,
			Name: itemType, // Simple name for now
			Position: Position{
				X: x,
				Y: y,
			},
		}

		items = append(items, item)
	}

	return items
}

// generateID generates a random ID
func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
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

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
