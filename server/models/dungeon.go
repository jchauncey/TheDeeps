package models

import (
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
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Position Position `json:"position"`
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

	// Generate rooms - ensure at least 3 rooms per floor
	minRooms := 3
	maxRooms := 10 + level // More rooms on deeper levels
	floor.Rooms = generateRooms(floor, maxRooms, 5, 15, 3, 20)

	// If we didn't generate enough rooms, try again with more relaxed parameters
	if len(floor.Rooms) < minRooms {
		floor.Rooms = generateRooms(floor, maxRooms, 4, 12, 2, 30)
	}

	// Connect rooms with corridors
	connectRooms(floor)

	// Add stairs
	addStairs(floor, level)

	// Add entities and items
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

		// Randomly decide whether to go horizontal first or vertical first
		if rand.Intn(2) == 0 {
			createHorizontalCorridor(floor, x1, x2, y1)
			createVerticalCorridor(floor, y1, y2, x2)
		} else {
			createVerticalCorridor(floor, y1, y2, x1)
			createHorizontalCorridor(floor, x1, x2, y2)
		}
	}

	// Add some additional random connections for better connectivity
	// This helps ensure rooms have multiple exits
	if len(floor.Rooms) > 2 {
		numExtraConnections := len(floor.Rooms)/3 + 1
		for i := 0; i < numExtraConnections; i++ {
			// Pick two random rooms
			r1 := rand.Intn(len(floor.Rooms))
			r2 := rand.Intn(len(floor.Rooms))

			// Make sure they're different rooms
			if r1 == r2 {
				r2 = (r2 + 1) % len(floor.Rooms)
			}

			room1 := floor.Rooms[r1]
			room2 := floor.Rooms[r2]

			// Get center points of each room
			x1, y1 := room1.Center()
			x2, y2 := room2.Center()

			// Create a corridor between them
			if rand.Intn(2) == 0 {
				createHorizontalCorridor(floor, x1, x2, y1)
				createVerticalCorridor(floor, y1, y2, x2)
			} else {
				createVerticalCorridor(floor, y1, y2, x1)
				createHorizontalCorridor(floor, x1, x2, y2)
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

// generateEntities generates entities for the floor
func generateEntities(floor *Floor, level int) []Entity {
	entities := make([]Entity, 0)

	// Number of entities scales with floor level
	numEntities := 5 + level

	for i := 0; i < numEntities; i++ {
		// Pick a random room (skip the first room for player safety)
		roomIndex := rand.Intn(len(floor.Rooms)-1) + 1
		room := floor.Rooms[roomIndex]

		// Random position within the room
		x := rand.Intn(room.Width-2) + room.X + 1
		y := rand.Intn(room.Height-2) + room.Y + 1

		// Create entity
		entityTypes := []string{"goblin", "orc", "skeleton", "rat", "bat"}
		entityType := entityTypes[rand.Intn(len(entityTypes))]

		entity := Entity{
			ID:       generateID(),
			Type:     entityType,
			Name:     entityType, // Simple name for now
			Position: Position{X: x, Y: y},
		}

		entities = append(entities, entity)
	}

	return entities
}

// generateItems generates items for the floor
func generateItems(floor *Floor, level int) []Item {
	items := make([]Item, 0)

	// Number of items scales with floor level
	numItems := 3 + level/2

	for i := 0; i < numItems; i++ {
		// Pick a random room
		roomIndex := rand.Intn(len(floor.Rooms))
		room := floor.Rooms[roomIndex]

		// Random position within the room
		x := rand.Intn(room.Width-2) + room.X + 1
		y := rand.Intn(room.Height-2) + room.Y + 1

		// Create item
		itemTypes := []string{"potion", "scroll", "weapon", "armor", "gold"}
		itemType := itemTypes[rand.Intn(len(itemTypes))]

		item := Item{
			ID:       generateID(),
			Type:     itemType,
			Name:     itemType, // Simple name for now
			Position: Position{X: x, Y: y},
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
