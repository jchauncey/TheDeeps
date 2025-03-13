package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// DungeonHandler handles dungeon-related HTTP requests
type DungeonHandler struct {
	dungeonRepo   *repositories.DungeonRepository
	characterRepo *repositories.CharacterRepository
	mapGenerator  *game.MapGenerator
}

// NewDungeonHandler creates a new dungeon handler
func NewDungeonHandler(dungeonRepo *repositories.DungeonRepository, characterRepo *repositories.CharacterRepository) *DungeonHandler {
	return &DungeonHandler{
		dungeonRepo:   dungeonRepo,
		characterRepo: characterRepo,
		mapGenerator:  game.NewMapGenerator(time.Now().UnixNano()),
	}
}

// GetDungeons handles GET /dungeons
func (h *DungeonHandler) GetDungeons(w http.ResponseWriter, r *http.Request) {
	dungeons := h.dungeonRepo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dungeons)
}

// CreateDungeon handles POST /dungeons
func (h *DungeonHandler) CreateDungeon(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request struct {
		Name       string `json:"name"`
		Floors     int    `json:"floors"`
		Difficulty string `json:"difficulty"`
		Seed       int64  `json:"seed,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if request.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if request.Floors < 1 {
		http.Error(w, "Floors must be at least 1", http.StatusBadRequest)
		return
	}

	if request.Difficulty == "" {
		request.Difficulty = "normal" // Default difficulty
	}

	// Create dungeon
	dungeon := models.NewDungeon(request.Name, request.Floors, request.Seed)
	dungeon.Difficulty = request.Difficulty

	// Generate first floor
	floor := dungeon.GenerateFloor(1)
	h.mapGenerator.GenerateFloorWithDifficulty(floor, 1, request.Floors == 1, request.Difficulty)

	// Save dungeon
	if err := h.dungeonRepo.Save(dungeon); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return dungeon
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dungeon)
}

// JoinDungeon handles POST /dungeons/{id}/join
func (h *DungeonHandler) JoinDungeon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dungeonID := vars["id"]

	// Parse request body
	var request struct {
		CharacterID string `json:"characterId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if request.CharacterID == "" {
		http.Error(w, "Character ID is required", http.StatusBadRequest)
		return
	}

	// Get dungeon
	_, err := h.dungeonRepo.GetByID(dungeonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get character
	character, err := h.characterRepo.GetByID(request.CharacterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Add character to dungeon
	if err := h.dungeonRepo.AddCharacterToDungeon(dungeonID, request.CharacterID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get first floor
	floor, err := h.dungeonRepo.GetFloor(dungeonID, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find the entrance room and place the character there
	var entranceRoom *models.Room
	for i := range floor.Rooms {
		if floor.Rooms[i].Type == models.RoomEntrance {
			entranceRoom = &floor.Rooms[i]
			break
		}
	}

	if entranceRoom != nil {
		// Place character in the center of the entrance room
		centerX := entranceRoom.X + entranceRoom.Width/2
		centerY := entranceRoom.Y + entranceRoom.Height/2
		character.Position = models.Position{X: centerX, Y: centerY}
	} else if len(floor.UpStairs) > 0 {
		// Fallback: Place character on first floor at the first up stairs
		character.Position = floor.UpStairs[0]
	} else {
		// Fallback: If no entrance room or up stairs, place at a random walkable tile
		for y := 0; y < floor.Height; y++ {
			for x := 0; x < floor.Width; x++ {
				if floor.Tiles[y][x].Walkable && floor.Tiles[y][x].MobID == "" {
					character.Position = models.Position{X: x, Y: y}
					break
				}
			}
			if character.Position.X != 0 || character.Position.Y != 0 {
				break
			}
		}
	}

	// Update the tile to mark the character's position
	floor.Tiles[character.Position.Y][character.Position.X].Character = character.ID

	// Update character
	character.CurrentFloor = 1
	character.CurrentDungeon = dungeonID

	// Save character
	if err := h.characterRepo.Save(character); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return floor
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(floor)
}

// GetFloor handles GET /dungeons/{id}/floor/{level}
func (h *DungeonHandler) GetFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dungeonID := vars["id"]
	levelStr := vars["level"]

	// Parse level
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		http.Error(w, "Invalid floor level", http.StatusBadRequest)
		return
	}

	// Get dungeon
	dungeon, err := h.dungeonRepo.GetByID(dungeonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate level
	if level < 1 || level > dungeon.Floors {
		http.Error(w, "Floor level out of range", http.StatusBadRequest)
		return
	}

	// Get floor
	floor, err := h.dungeonRepo.GetFloor(dungeonID, level)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If floor hasn't been generated yet, generate it
	if len(floor.Rooms) == 0 {
		h.mapGenerator.GenerateFloorWithDifficulty(floor, level, level == dungeon.Floors, dungeon.Difficulty)
	}

	// Return floor
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(floor)
}

// GenerateTestRoom handles GET /test/room
// This endpoint generates a single room for testing client rendering
func (h *DungeonHandler) GenerateTestRoom(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Get room type (default to entrance)
	roomType := models.RoomType(query.Get("type"))
	if roomType == "" {
		roomType = models.RoomEntrance
	}

	// Validate room type
	validTypes := map[models.RoomType]bool{
		models.RoomEntrance: true,
		models.RoomStandard: true,
		models.RoomTreasure: true,
		models.RoomBoss:     true,
		models.RoomSafe:     true,
		models.RoomShop:     true,
	}

	if !validTypes[roomType] {
		http.Error(w, "Invalid room type. Valid types: entrance, standard, treasure, boss, safe, shop", http.StatusBadRequest)
		return
	}

	// Get width and height (default to 20x20)
	width := 20
	height := 20

	if widthStr := query.Get("width"); widthStr != "" {
		if w, err := strconv.Atoi(widthStr); err == nil && w > 0 && w <= 100 {
			width = w
		}
	}

	if heightStr := query.Get("height"); heightStr != "" {
		if h, err := strconv.Atoi(heightStr); err == nil && h > 0 && h <= 100 {
			height = h
		}
	}

	// Get room size (default to 8x8 for entrance, 7x7 for others)
	roomWidth := 8
	roomHeight := 8

	if roomType != models.RoomEntrance {
		roomWidth = 7
		roomHeight = 7
	}

	if roomWidthStr := query.Get("roomWidth"); roomWidthStr != "" {
		if rw, err := strconv.Atoi(roomWidthStr); err == nil && rw > 0 && rw < width {
			roomWidth = rw
		}
	}

	if roomHeightStr := query.Get("roomHeight"); roomHeightStr != "" {
		if rh, err := strconv.Atoi(roomHeightStr); err == nil && rh > 0 && rh < height {
			roomHeight = rh
		}
	}

	// Create a test floor
	floor := &models.Floor{
		Level:      1,
		Width:      width,
		Height:     height,
		Tiles:      make([][]models.Tile, height),
		UpStairs:   []models.Position{},
		DownStairs: []models.Position{},
		Mobs:       make(map[string]*models.Mob),
		Items:      make(map[string]models.Item),
	}

	// Initialize tiles with walls
	for y := 0; y < height; y++ {
		floor.Tiles[y] = make([]models.Tile, width)
		for x := 0; x < width; x++ {
			floor.Tiles[y][x] = models.Tile{
				Type:     models.TileWall,
				Walkable: false,
				Explored: false,
			}
		}
	}

	// Calculate room position (center of the floor)
	roomX := (width - roomWidth) / 2
	roomY := (height - roomHeight) / 2

	// Create the room
	room := models.Room{
		ID:       "test-room",
		Type:     roomType,
		X:        roomX,
		Y:        roomY,
		Width:    roomWidth,
		Height:   roomHeight,
		Explored: true,
	}

	// Carve out the room
	for y := 0; y < roomHeight; y++ {
		for x := 0; x < roomWidth; x++ {
			floor.Tiles[roomY+y][roomX+x] = models.Tile{
				Type:     models.TileFloor,
				Walkable: true,
				Explored: true,
				RoomID:   room.ID,
			}
		}
	}

	floor.Rooms = []models.Room{room}

	// Add stairs if it's an entrance room
	if roomType == models.RoomEntrance {
		// Add down stairs in the bottom right corner
		stairsX := roomX + roomWidth - 2
		stairsY := roomY + roomHeight - 2

		floor.Tiles[stairsY][stairsX] = models.Tile{
			Type:     models.TileDownStairs,
			Walkable: true,
			Explored: true,
			RoomID:   room.ID,
		}

		floor.DownStairs = append(floor.DownStairs, models.Position{X: stairsX, Y: stairsY})
	}

	// Add a test character in the center of the room
	characterX := roomX + roomWidth/2
	characterY := roomY + roomHeight/2

	// Set the character in the tile
	floor.Tiles[characterY][characterX].Character = "test-character"

	// Add items for treasure rooms
	if roomType == models.RoomTreasure {
		// Add some test items
		for i := 0; i < 3; i++ {
			itemX := roomX + 2 + i*2
			itemY := roomY + 2

			itemID := "test-item-" + strconv.Itoa(i)
			floor.Items[itemID] = models.Item{
				ID:       itemID,
				Type:     models.ItemWeapon,
				Name:     "Test Item " + strconv.Itoa(i),
				Position: models.Position{X: itemX, Y: itemY},
			}

			floor.Tiles[itemY][itemX].ItemID = itemID
		}
	}

	// Add a boss for boss rooms
	if roomType == models.RoomBoss {
		mobID := "test-boss"
		mob := models.NewMob(models.MobDragon, models.VariantBoss, 10)
		mob.ID = mobID
		mob.Position = models.Position{X: characterX + 2, Y: characterY}

		floor.Mobs[mobID] = mob
		floor.Tiles[characterY][characterX+2].MobID = mobID
	}

	// Return the floor with the single room
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(floor)
}
