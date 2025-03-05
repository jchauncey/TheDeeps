package mapgen

import (
	"math/rand"
	"time"
)

// TileType represents the type of a tile
type TileType int

const (
	TileWall TileType = iota
	TileFloor
	TileEntrance
	TileExit
	TileHallway
	TilePillar
	TileWater
	TileRubble
)

// Tile represents a single tile in the dungeon
type Tile struct {
	Type     TileType
	Symbol   rune
	Walkable bool
}

// Position represents a position in the dungeon
type Position struct {
	X int
	Y int
}

// Room represents a room in the dungeon
type Room struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Level represents a single level of the dungeon
type Level struct {
	Width  int
	Height int
	Tiles  [][]Tile
	Rooms  []Room
}

// Floor represents a floor of the dungeon
type Floor struct {
	Level    Level
	Entrance Position
	Exit     Position
}

// Dungeon represents the entire dungeon
type Dungeon struct {
	Floors []Floor
}

// DungeonConfig holds configuration for dungeon generation
type DungeonConfig struct {
	Width       int
	Height      int
	TotalFloors int
	MinRooms    int
	MaxRooms    int
	MinRoomSize int
	MaxRoomSize int
}

// NewDefaultConfig creates a default dungeon configuration
func NewDefaultConfig() *DungeonConfig {
	return &DungeonConfig{
		Width:       80,
		Height:      40,
		TotalFloors: 5,
		MinRooms:    5,
		MaxRooms:    10,
		MinRoomSize: 5,
		MaxRoomSize: 12,
	}
}

// NewDungeon creates a new dungeon with the given number of floors
func NewDungeon(numFloors int) *Dungeon {
	rand.Seed(time.Now().UnixNano())

	dungeon := &Dungeon{
		Floors: make([]Floor, numFloors),
	}

	// Generate each floor
	for i := 0; i < numFloors; i++ {
		dungeon.Floors[i] = generateFloor(80, 40)
	}

	return dungeon
}

// generateFloor generates a single floor of the dungeon
func generateFloor(width, height int) Floor {
	// Create empty level
	level := Level{
		Width:  width,
		Height: height,
		Tiles:  make([][]Tile, height),
		Rooms:  []Room{},
	}

	// Initialize all tiles as walls
	for y := 0; y < height; y++ {
		level.Tiles[y] = make([]Tile, width)
		for x := 0; x < width; x++ {
			level.Tiles[y][x] = Tile{
				Type:     TileWall,
				Symbol:   '#',
				Walkable: false,
			}
		}
	}

	// Generate rooms
	numRooms := rand.Intn(5) + 5 // 5-10 rooms
	for i := 0; i < numRooms; i++ {
		roomWidth := rand.Intn(8) + 5  // 5-12 width
		roomHeight := rand.Intn(6) + 4 // 4-9 height
		roomX := rand.Intn(width-roomWidth-2) + 1
		roomY := rand.Intn(height-roomHeight-2) + 1

		// Check if room overlaps with existing rooms
		overlaps := false
		for _, room := range level.Rooms {
			if roomX <= room.X+room.Width && roomX+roomWidth >= room.X &&
				roomY <= room.Y+room.Height && roomY+roomHeight >= room.Y {
				overlaps = true
				break
			}
		}

		if !overlaps {
			// Add room
			room := Room{
				X:      roomX,
				Y:      roomY,
				Width:  roomWidth,
				Height: roomHeight,
			}
			level.Rooms = append(level.Rooms, room)

			// Carve out room
			for y := roomY; y < roomY+roomHeight; y++ {
				for x := roomX; x < roomX+roomWidth; x++ {
					level.Tiles[y][x] = Tile{
						Type:     TileFloor,
						Symbol:   '.',
						Walkable: true,
					}
				}
			}

			// Add some pillars
			if roomWidth >= 8 && roomHeight >= 6 && rand.Intn(3) == 0 {
				pillarX := roomX + roomWidth/3
				pillarY := roomY + roomHeight/3
				level.Tiles[pillarY][pillarX] = Tile{
					Type:     TilePillar,
					Symbol:   'O',
					Walkable: false,
				}

				pillarX = roomX + (roomWidth*2)/3
				pillarY = roomY + (roomHeight*2)/3
				level.Tiles[pillarY][pillarX] = Tile{
					Type:     TilePillar,
					Symbol:   'O',
					Walkable: false,
				}
			}

			// Add water or rubble
			if rand.Intn(5) == 0 {
				featureType := TileWater
				featureSymbol := '~'
				if rand.Intn(2) == 0 {
					featureType = TileRubble
					featureSymbol = ','
				}

				featureX := roomX + rand.Intn(roomWidth-2) + 1
				featureY := roomY + rand.Intn(roomHeight-2) + 1
				featureWidth := rand.Intn(3) + 2
				featureHeight := rand.Intn(2) + 2

				for y := featureY; y < featureY+featureHeight && y < roomY+roomHeight-1; y++ {
					for x := featureX; x < featureX+featureWidth && x < roomX+roomWidth-1; x++ {
						level.Tiles[y][x] = Tile{
							Type:     featureType,
							Symbol:   featureSymbol,
							Walkable: featureType == TileRubble, // Water is not walkable
						}
					}
				}
			}
		}
	}

	// Connect rooms with hallways
	for i := 0; i < len(level.Rooms)-1; i++ {
		room1 := level.Rooms[i]
		room2 := level.Rooms[i+1]

		// Get center points of rooms
		x1 := room1.X + room1.Width/2
		y1 := room1.Y + room1.Height/2
		x2 := room2.X + room2.Width/2
		y2 := room2.Y + room2.Height/2

		// Randomly decide whether to go horizontal first or vertical first
		if rand.Intn(2) == 0 {
			// Horizontal first
			for x := min(x1, x2); x <= max(x1, x2); x++ {
				if level.Tiles[y1][x].Type == TileWall {
					level.Tiles[y1][x] = Tile{
						Type:     TileHallway,
						Symbol:   '.',
						Walkable: true,
					}
				}
			}
			for y := min(y1, y2); y <= max(y1, y2); y++ {
				if level.Tiles[y][x2].Type == TileWall {
					level.Tiles[y][x2] = Tile{
						Type:     TileHallway,
						Symbol:   '.',
						Walkable: true,
					}
				}
			}
		} else {
			// Vertical first
			for y := min(y1, y2); y <= max(y1, y2); y++ {
				if level.Tiles[y][x1].Type == TileWall {
					level.Tiles[y][x1] = Tile{
						Type:     TileHallway,
						Symbol:   '.',
						Walkable: true,
					}
				}
			}
			for x := min(x1, x2); x <= max(x1, x2); x++ {
				if level.Tiles[y2][x].Type == TileWall {
					level.Tiles[y2][x] = Tile{
						Type:     TileHallway,
						Symbol:   '.',
						Walkable: true,
					}
				}
			}
		}
	}

	// Place entrance and exit
	entranceRoom := level.Rooms[0]
	exitRoom := level.Rooms[len(level.Rooms)-1]

	entranceX := entranceRoom.X + entranceRoom.Width/2
	entranceY := entranceRoom.Y + entranceRoom.Height/2
	exitX := exitRoom.X + exitRoom.Width/2
	exitY := exitRoom.Y + exitRoom.Height/2

	level.Tiles[entranceY][entranceX] = Tile{
		Type:     TileEntrance,
		Symbol:   '<',
		Walkable: true,
	}

	level.Tiles[exitY][exitX] = Tile{
		Type:     TileExit,
		Symbol:   '>',
		Walkable: true,
	}

	return Floor{
		Level:    level,
		Entrance: Position{X: entranceX, Y: entranceY},
		Exit:     Position{X: exitX, Y: exitY},
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
