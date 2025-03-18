package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TileType represents the type of tile on the map
type TileType string

const (
	TileWall       TileType = "#"
	TileFloor      TileType = "."
	TileUpStairs   TileType = "<"
	TileDownStairs TileType = ">"
	TileDoor       TileType = "+"
	TileChest      TileType = "C"
	TileTrap       TileType = "^"
)

// RoomType represents the type of room
type RoomType string

const (
	RoomStandard RoomType = "standard"
	RoomTreasure RoomType = "treasure"
	RoomBoss     RoomType = "boss"
	RoomPuzzle   RoomType = "puzzle"
	RoomSafe     RoomType = "safe"
	RoomShop     RoomType = "shop"
	RoomEntrance RoomType = "entrance" // New room type for dungeon entrance
)

// Room represents a room in the dungeon
type Room struct {
	ID       string   `json:"id"`
	Type     RoomType `json:"type"`
	X        int      `json:"x"`
	Y        int      `json:"y"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	Explored bool     `json:"explored"`
}

// Tile represents a single tile on the map
type Tile struct {
	Type      TileType `json:"type"`
	Walkable  bool     `json:"walkable"`
	Explored  bool     `json:"explored"`
	RoomID    string   `json:"roomId,omitempty"`
	MobID     string   `json:"mobId,omitempty"`
	ItemID    string   `json:"itemId,omitempty"`
	Character string   `json:"character,omitempty"` // Character ID if a player is on this tile
}

// Floor represents a single floor of the dungeon
type Floor struct {
	Level      int             `json:"level"`
	Width      int             `json:"width"`
	Height     int             `json:"height"`
	Tiles      [][]Tile        `json:"tiles"`
	Rooms      []Room          `json:"rooms"`
	UpStairs   []Position      `json:"upStairs"`
	DownStairs []Position      `json:"downStairs"`
	Mobs       map[string]*Mob `json:"mobs"`
	Items      map[string]Item `json:"items"`
}

// Dungeon represents a complete dungeon
type Dungeon struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Floors      int               `json:"floors"`
	Difficulty  string            `json:"difficulty"`
	CreatedAt   time.Time         `json:"createdAt"`
	FloorData   map[int]*Floor    `json:"floorData"`
	Characters  map[string]string `json:"characters"` // Map of character ID to floor level
	Seed        int64             `json:"seed"`
	PlayerCount int               `json:"playerCount"`
}

// NewDungeon creates a new dungeon with the specified number of floors
func NewDungeon(name string, floors int, seed int64) *Dungeon {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	return &Dungeon{
		ID:         uuid.New().String(),
		Name:       name,
		Floors:     floors,
		Difficulty: "normal", // Default difficulty
		CreatedAt:  time.Now(),
		FloorData:  make(map[int]*Floor),
		Characters: make(map[string]string),
		Seed:       seed,
	}
}

// GenerateFloor generates a new floor for the dungeon
func (d *Dungeon) GenerateFloor(level int) *Floor {
	// Floor dimensions based on level (deeper floors can be larger)
	width := 50 + (level * 2)
	height := 50 + (level * 2)
	if width > 100 {
		width = 100
	}
	if height > 100 {
		height = 100
	}

	// Initialize tiles with walls
	tiles := make([][]Tile, height)
	for y := range tiles {
		tiles[y] = make([]Tile, width)
		for x := range tiles[y] {
			tiles[y][x] = Tile{
				Type:     TileWall,
				Walkable: false,
				Explored: false,
			}
		}
	}

	// Create a new floor
	floor := &Floor{
		Level:      level,
		Width:      width,
		Height:     height,
		Tiles:      tiles,
		Rooms:      []Room{},
		UpStairs:   []Position{},
		DownStairs: []Position{},
		Mobs:       make(map[string]*Mob),
		Items:      make(map[string]Item),
	}

	// Store the floor in the dungeon
	d.FloorData[level] = floor

	return floor
}

// AddCharacter adds a character to the dungeon
func (d *Dungeon) AddCharacter(characterID string) {
	d.Characters[characterID] = "1" // Start at floor 1
	d.PlayerCount++
}

// RemoveCharacter removes a character from the dungeon
func (d *Dungeon) RemoveCharacter(characterID string) {
	delete(d.Characters, characterID)
	d.PlayerCount--
}

// GetCharacterFloor gets the floor level for a character
func (d *Dungeon) GetCharacterFloor(characterID string) int {
	floorStr, exists := d.Characters[characterID]
	if !exists {
		return 0
	}

	// Convert string to int (in a real implementation, you'd handle errors)
	var floor int
	_, err := fmt.Sscanf(floorStr, "%d", &floor)
	if err != nil {
		return 0
	}

	return floor
}

// SetCharacterFloor sets the floor level for a character
func (d *Dungeon) SetCharacterFloor(characterID string, floor int) {
	d.Characters[characterID] = fmt.Sprintf("%d", floor)
}
