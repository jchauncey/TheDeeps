package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// Game constants
const (
	VisibilityRange = 10             // Adjust this value to change the fog of war size
	MaxInactivity   = 24 * time.Hour // Maximum inactivity before dungeon cleanup
)

// Debug flags
var (
	DebugMode = false // Set to true to enable debug features
)

// GameServer manages the game state and client connections
type GameServer struct {
	Game                *Game                      // Keep for backward compatibility
	Clients             map[*websocket.Conn]string // Map of connections to character IDs
	CharacterRepository *repositories.CharacterRepository
	DungeonRepository   *repositories.DungeonRepository
	MobSpawner          *MobSpawner
	Upgrader            websocket.Upgrader
}

// Game represents the game state
type Game struct {
	Dungeon *models.Dungeon
	Players map[string]*models.Character
}

// Message types
type DebugMessage struct {
	Message   string `json:"message"`
	Level     string `json:"level"`
	Timestamp int64  `json:"timestamp"`
}

type FloorMessage struct {
	Type         string            `json:"type"`
	Floor        *models.Floor     `json:"floor"`
	PlayerPos    models.Position   `json:"playerPosition"`
	CurrentFloor int               `json:"currentFloor"`
	PlayerData   *models.Character `json:"playerData"`
	DungeonID    string            `json:"dungeonId"`
}

type MoveMessage struct {
	Type      string `json:"type"`
	Direction string `json:"direction"`
}

type ActionMessage struct {
	Type   string `json:"type"`
	Action string `json:"action"`
}

type CreateCharacterRequest struct {
	Name           string       `json:"name"`
	CharacterClass string       `json:"characterClass"`
	Stats          models.Stats `json:"stats"`
}

// New message types for dungeon management
type CreateDungeonMessage struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	NumFloors int    `json:"numFloors"`
}

type JoinDungeonMessage struct {
	Type        string `json:"type"`
	DungeonID   string `json:"dungeonId"`
	CharacterID string `json:"characterId"`
}

type ListDungeonsMessage struct {
	Type string `json:"type"`
}

type DungeonListResponse struct {
	Type     string                    `json:"type"`
	Dungeons []DungeonListItemResponse `json:"dungeons"`
}

type DungeonListItemResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	PlayerCount  int       `json:"playerCount"`
	CreatedAt    time.Time `json:"createdAt"`
	LastActivity time.Time `json:"lastActivity"`
}

// LeaveDungeonMessage represents a request to leave a dungeon
type LeaveDungeonMessage struct {
	Type        string `json:"type"`
	DungeonID   string `json:"dungeonId"`
	CharacterID string `json:"characterId"`
}

// NewGameServer creates a new game server instance
func NewGameServer() *GameServer {
	// Create a default dungeon for backward compatibility
	defaultDungeon := models.NewDungeon(5)

	game := &Game{
		Dungeon: defaultDungeon,
		Players: make(map[string]*models.Character),
	}

	server := &GameServer{
		Game:                game,
		Clients:             make(map[*websocket.Conn]string),
		CharacterRepository: repositories.NewCharacterRepository(),
		DungeonRepository:   repositories.NewDungeonRepository(),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}

	// Initialize the mob spawner
	server.MobSpawner = NewMobSpawner(nil)

	// Start dungeon cleanup routine
	go server.cleanupInactiveDungeons()

	// Create a default dungeon instance with the default dungeon
	defaultInstance := server.DungeonRepository.Create("Default Dungeon", 5)
	defaultInstance.Dungeon = defaultDungeon

	return server
}

// Initialize game state
func initGame() *Game {
	log.Println("Initializing game...")
	dungeon := models.NewDungeon(10)

	log.Printf("Created dungeon with %d floors", len(dungeon.Floors))
	for i, floor := range dungeon.Floors {
		log.Printf("Floor %d: %d rooms", i+1, len(floor.Rooms))
	}

	game := &Game{
		Dungeon: dungeon,
		Players: make(map[string]*models.Character),
	}

	return game
}

// SetupRoutes configures the HTTP routes for the server
func (s *GameServer) SetupRoutes(handler any) *mux.Router {
	h := handler.(interface {
		HandleCreateCharacter(w http.ResponseWriter, r *http.Request)
		HandleGetCharacter(w http.ResponseWriter, r *http.Request)
		HandleGetCharacters(w http.ResponseWriter, r *http.Request)
		HandleGetFloor(w http.ResponseWriter, r *http.Request)
		HandleGetCharacterFloor(w http.ResponseWriter, r *http.Request)
		HandleCreateDungeon(w http.ResponseWriter, r *http.Request)
		HandleListDungeons(w http.ResponseWriter, r *http.Request)
		HandleJoinDungeon(w http.ResponseWriter, r *http.Request)
	})

	r := mux.NewRouter()

	// WebSocket endpoint
	r.HandleFunc("/ws", s.HandleWebSocket)

	// Character endpoints
	r.HandleFunc("/character", h.HandleCreateCharacter).Methods("POST")
	r.HandleFunc("/character/{id}", h.HandleGetCharacter).Methods("GET")
	r.HandleFunc("/characters", h.HandleGetCharacters).Methods("GET")
	r.HandleFunc("/character/{characterId}/floor", h.HandleGetCharacterFloor).Methods("GET")

	// Dungeon endpoints
	r.HandleFunc("/dungeon", h.HandleCreateDungeon).Methods("POST")
	r.HandleFunc("/dungeons", h.HandleListDungeons).Methods("GET")
	r.HandleFunc("/dungeon/{dungeonId}/join", h.HandleJoinDungeon).Methods("POST")
	r.HandleFunc("/dungeon/{dungeonId}/floor/{level}", h.HandleGetFloor).Methods("GET")

	return r
}

// HandleWebSocket handles WebSocket connections
func (s *GameServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, Upgrade, Connection")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Register client
	s.Clients[conn] = ""

	// Ensure proper cleanup when function returns
	defer func() {
		// Get character ID associated with this connection
		characterID := s.Clients[conn]

		// If there's a character associated with this connection, remove them from any dungeons
		if characterID != "" {
			// Find which dungeon the character is in
			dungeonInstance := s.DungeonRepository.GetPlayerDungeon(characterID)
			if dungeonInstance != nil {
				// Remove the player from the dungeon
				err := s.DungeonRepository.RemovePlayerFromDungeon(dungeonInstance.ID, characterID)
				if err != nil {
					log.Printf("Error removing character %s from dungeon %s during disconnect: %v",
						characterID, dungeonInstance.ID, err)
				} else {
					log.Printf("Character %s removed from dungeon %s during disconnect",
						characterID, dungeonInstance.ID)
				}
			}
		}

		// Clean up connection
		conn.Close()
		delete(s.Clients, conn)

		// Log the disconnection
		if characterID != "" {
			log.Printf("Client disconnected: character ID %s", characterID)
		} else {
			log.Printf("Client disconnected: no character associated")
		}

		log.Println("WebSocket connection closed and cleaned up")
	}()

	// Send a welcome message
	welcomeMsg := map[string]interface{}{
		"type":      "welcome",
		"message":   "Welcome to The Deeps!",
		"timestamp": time.Now().Unix(),
	}

	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Println("Error sending welcome message:", err)
		return
	}

	// Message handling loop
	for {
		// Reset read deadline for each message
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		_, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			} else {
				log.Printf("Client disconnected: %v", err)
			}
			break
		}

		var message map[string]interface{}
		if err := json.Unmarshal(p, &message); err != nil {
			log.Println("Error parsing message:", err)
			conn.WriteJSON(DebugMessage{
				Message:   fmt.Sprintf("Error parsing message: %v", err),
				Level:     "error",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		msgType, ok := message["type"].(string)
		if !ok {
			log.Println("Message missing 'type' field")
			conn.WriteJSON(DebugMessage{
				Message:   "Message missing 'type' field",
				Level:     "error",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		// Reset write deadline for response
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		switch msgType {
		case "get_floor":
			s.SendFloorData(conn)
		case "move":
			s.HandleMove(conn, p)
		case "action":
			s.HandleAction(conn, p)
		case "create_character":
			s.HandleCreateCharacter(conn, p)
		case "create_dungeon":
			s.HandleCreateDungeon(conn, p)
		case "join_dungeon":
			s.HandleJoinDungeon(conn, p)
		case "leave_dungeon":
			s.HandleLeaveDungeon(conn, p)
		case "list_dungeons":
			s.HandleListDungeons(conn)
		default:
			log.Printf("Unknown message type: %s", msgType)
			conn.WriteJSON(DebugMessage{
				Message:   fmt.Sprintf("Unknown message type: %s", msgType),
				Level:     "error",
				Timestamp: time.Now().Unix(),
			})
		}
	}
}

// HandleMove processes movement commands
func (s *GameServer) HandleMove(conn *websocket.Conn, payload []byte) {
	var moveMsg MoveMessage
	if err := json.Unmarshal(payload, &moveMsg); err != nil {
		log.Printf("Error parsing move message: %v", err)
		return
	}

	log.Printf("Received move command: %s", moveMsg.Direction)

	// Get the character ID associated with this connection
	characterID, exists := s.Clients[conn]
	if !exists {
		log.Printf("No character associated with connection")
		return
	}

	// Find the dungeon the player is in
	dungeon := s.DungeonRepository.GetPlayerDungeon(characterID)
	if dungeon == nil {
		log.Printf("No dungeon found for character %s", characterID)
		return
	}

	// Get the player's position and floor
	currentPos := dungeon.GetPlayerPosition(characterID)
	if currentPos == nil {
		log.Printf("No position found for character %s", characterID)
		return
	}

	floorIndex := dungeon.GetPlayerFloor(characterID)
	floor := dungeon.Dungeon.Floors[floorIndex]

	newPos := *currentPos

	// Calculate new position
	switch moveMsg.Direction {
	case "up":
		newPos.Y--
	case "down":
		newPos.Y++
	case "left":
		newPos.X--
	case "right":
		newPos.X++
	default:
		log.Printf("Invalid direction: %s", moveMsg.Direction)
		return
	}

	log.Printf("Character %s - Current position: (%d, %d), New position: (%d, %d)",
		characterID, currentPos.X, currentPos.Y, newPos.X, newPos.Y)

	// Check if the move is valid
	if newPos.X >= 0 && newPos.X < floor.Width &&
		newPos.Y >= 0 && newPos.Y < floor.Height &&
		floor.Tiles[newPos.Y][newPos.X].Type != models.TileWall {

		// Update player position in the dungeon
		*currentPos = newPos

		// Find and update player entity
		for i, entity := range floor.Entities {
			if entity.Type == "player" && entity.ID == characterID {
				floor.Entities[i].Position = newPos
				break
			}
		}

		// Get the character
		character, err := s.CharacterRepository.GetByID(characterID)
		if err != nil {
			log.Printf("Error getting character: %v", err)
			return
		}

		// If no player entity exists for this character, create one
		playerEntityExists := false
		for _, entity := range floor.Entities {
			if entity.Type == "player" && entity.ID == characterID {
				playerEntityExists = true
				break
			}
		}

		if !playerEntityExists {
			playerEntity := models.Entity{
				ID:             characterID,
				Type:           "player",
				Name:           character.Name,
				Position:       newPos,
				CharacterClass: character.CharacterClass,
				Health:         character.Health,
				MaxHealth:      character.MaxHealth,
			}
			floor.Entities = append(floor.Entities, playerEntity)
		}

		// Update visibility
		// TODO: Update this to work with the new dungeon model
		// UpdateVisibility(dungeon.Dungeon)

		log.Printf("Player %s moved %s to (%d, %d)", characterID, moveMsg.Direction, newPos.X, newPos.Y)

		// Send updated floor data to the player
		s.SendFloorData(conn)
	} else {
		log.Printf("Invalid move to (%d, %d) for character %s", newPos.X, newPos.Y, characterID)
	}
}

// HandleAction handles player actions
func (s *GameServer) HandleAction(conn *websocket.Conn, payload []byte) {
	var message ActionMessage
	err := json.Unmarshal(payload, &message)
	if err != nil {
		log.Printf("Error unmarshaling action message: %v", err)
		return
	}

	log.Printf("Received action: %s", message.Action)

	// Get the character ID associated with this connection
	characterID, exists := s.Clients[conn]
	if !exists {
		log.Printf("No character associated with connection")
		return
	}

	// Find the dungeon the player is in
	dungeon := s.DungeonRepository.GetPlayerDungeon(characterID)
	if dungeon == nil {
		log.Printf("No dungeon found for character %s", characterID)
		return
	}

	// Get the player's position
	currentPos := dungeon.GetPlayerPosition(characterID)
	if currentPos == nil {
		log.Printf("No position found for character %s", characterID)
		return
	}

	// Get the player's current floor
	currentFloorIndex := dungeon.GetPlayerFloor(characterID)
	currentFloor := dungeon.Dungeon.Floors[currentFloorIndex]

	switch message.Action {
	case "pickup":
		// TODO: Update this to work with the new dungeon model
		log.Printf("Pickup action not yet implemented for new dungeon model")
	case "attack":
		// TODO: Update this to work with the new dungeon model
		log.Printf("Attack action not yet implemented for new dungeon model")
	case "descend":
		// Check if player is standing on stairs down
		if currentFloor.Tiles[currentPos.Y][currentPos.X].Type == models.TileStairsDown {
			if currentFloorIndex < len(dungeon.Dungeon.Floors)-1 {
				// Move player to next floor
				nextFloorIndex := currentFloorIndex + 1
				nextFloor := dungeon.Dungeon.Floors[nextFloorIndex]

				// Find stairs up on the next floor
				var stairsUpPos *models.Position
				for y := 0; y < nextFloor.Height; y++ {
					for x := 0; x < nextFloor.Width; x++ {
						if nextFloor.Tiles[y][x].Type == models.TileStairsUp {
							stairsUpPos = &models.Position{X: x, Y: y}
							break
						}
					}
					if stairsUpPos != nil {
						break
					}
				}

				if stairsUpPos != nil {
					// Update player position and floor
					dungeon.UpdatePlayerPosition(characterID, *stairsUpPos)
					dungeon.UpdatePlayerFloor(characterID, nextFloorIndex)

					// Update last activity time
					dungeon.LastActivity = time.Now()

					log.Printf("Player %s descended to floor %d", characterID, nextFloorIndex+1)
				} else {
					log.Printf("Could not find stairs up on floor %d", nextFloorIndex+1)
				}
			} else {
				log.Println("Already at the bottom floor")
			}
		} else {
			log.Println("Not standing on stairs down")
		}
	case "ascend":
		// Check if player is standing on stairs up
		if currentFloor.Tiles[currentPos.Y][currentPos.X].Type == models.TileStairsUp {
			if currentFloorIndex > 0 {
				// Move player to previous floor
				prevFloorIndex := currentFloorIndex - 1
				prevFloor := dungeon.Dungeon.Floors[prevFloorIndex]

				// Find stairs down on the previous floor
				var stairsDownPos *models.Position
				for y := 0; y < prevFloor.Height; y++ {
					for x := 0; x < prevFloor.Width; x++ {
						if prevFloor.Tiles[y][x].Type == models.TileStairsDown {
							stairsDownPos = &models.Position{X: x, Y: y}
							break
						}
					}
					if stairsDownPos != nil {
						break
					}
				}

				if stairsDownPos != nil {
					// Update player position and floor
					dungeon.UpdatePlayerPosition(characterID, *stairsDownPos)
					dungeon.UpdatePlayerFloor(characterID, prevFloorIndex)

					// Update last activity time
					dungeon.LastActivity = time.Now()

					log.Printf("Player %s ascended to floor %d", characterID, prevFloorIndex+1)
				} else {
					log.Printf("Could not find stairs down on floor %d", prevFloorIndex+1)
				}
			} else {
				log.Println("Already at the top floor")
			}
		} else {
			log.Println("Not standing on stairs up")
		}
	default:
		log.Printf("Unknown action: %s", message.Action)
	}

	// Send updated floor data to the player
	s.SendFloorData(conn)
}

// IsValidMove checks if a move is valid
func (s *GameServer) IsValidMove(pos models.Position) bool {
	currentFloor := s.Game.Dungeon.CurrentFloor
	floor := s.Game.Dungeon.Floors[currentFloor]

	// Check bounds
	if pos.X < 0 || pos.X >= floor.Width || pos.Y < 0 || pos.Y >= floor.Height {
		log.Printf("Move out of bounds: (%d, %d), floor dimensions: %dx%d", pos.X, pos.Y, floor.Width, floor.Height)
		return false
	}

	// Check walkable
	tile := floor.Tiles[pos.Y][pos.X]
	isWalkable := tile.Type == models.TileFloor ||
		tile.Type == models.TileStairsUp ||
		tile.Type == models.TileStairsDown ||
		tile.Type == models.TileDoor

	if !isWalkable {
		log.Printf("Tile at (%d, %d) is not walkable: %s", pos.X, pos.Y, tile.Type)
	}

	return isWalkable
}

// PickupItem handles item pickup
func (s *GameServer) PickupItem() {
	currentFloor := s.Game.Dungeon.CurrentFloor
	floor := s.Game.Dungeon.Floors[currentFloor]
	playerPos := s.Game.Dungeon.PlayerPosition

	for i, item := range floor.Items {
		if item.Position.X == playerPos.X && item.Position.Y == playerPos.Y {
			log.Printf("Player picked up %s", item.Name)
			floor.Items = append(floor.Items[:i], floor.Items[i+1:]...)
			s.BroadcastFloorData()
			return
		}
	}
	log.Println("No item to pick up")
}

// AttackEntity handles attacking entities
func (s *GameServer) AttackEntity() {
	currentFloor := s.Game.Dungeon.CurrentFloor
	floor := s.Game.Dungeon.Floors[currentFloor]
	playerPos := s.Game.Dungeon.PlayerPosition

	for i, entity := range floor.Entities {
		dx := Abs(entity.Position.X - playerPos.X)
		dy := Abs(entity.Position.Y - playerPos.Y)

		if dx <= 1 && dy <= 1 {
			log.Printf("Player attacked %s", entity.Name)
			floor.Entities = append(floor.Entities[:i], floor.Entities[i+1:]...)
			s.BroadcastFloorData()
			s.handleEntityDeath(entity)
			return
		}
	}
	log.Println("No entity to attack")
}

// handleEntityDeath handles the death of an entity
func (s *GameServer) handleEntityDeath(entity models.Entity) {
	log.Printf("Entity %s (%s) has been defeated", entity.Name, entity.ID)

	// Get the current floor
	currentFloor := s.Game.Dungeon.CurrentFloor
	if currentFloor < 0 || currentFloor >= len(s.Game.Dungeon.Floors) {
		return
	}
	floor := s.Game.Dungeon.Floors[currentFloor]

	// Create a mob instance from the entity for loot generation
	mobType := models.MobType(entity.Type)

	// Determine difficulty from the entity name
	// This is a simple approach - in a real game, you might store this in the entity data
	difficulty := models.DifficultyNormal
	if strings.HasPrefix(entity.Name, "easy") {
		difficulty = models.DifficultyEasy
	} else if strings.HasPrefix(entity.Name, "hard") {
		difficulty = models.DifficultyHard
	} else if strings.HasPrefix(entity.Name, "elite") {
		difficulty = models.DifficultyElite
	} else if strings.HasPrefix(entity.Name, "boss") {
		difficulty = models.DifficultyBoss
	}

	// Create a temporary mob instance for loot generation
	mobInstance := &models.MobInstance{
		ID:         entity.ID,
		Type:       mobType,
		Name:       entity.Name,
		Difficulty: difficulty,
		Position:   entity.Position,
		Health:     entity.Health,
		MaxHealth:  entity.MaxHealth,
		Damage:     entity.Damage,
		Defense:    entity.Defense,
		Speed:      entity.Speed,
		GoldDrop:   calculateGoldDrop(difficulty, currentFloor+1),
		Status:     entity.Status,
	}

	// Generate loot
	loot := models.GenerateLootFromMob(mobInstance, currentFloor+1)

	// Add items to the floor
	for _, item := range loot {
		// Convert to the floor's Item type
		floorItem := models.Item{
			ID:       item.ID,
			Type:     string(item.Type),
			Name:     item.Name,
			Position: item.Position,
		}

		floor.Items = append(floor.Items, floorItem)
		log.Printf("Added item %s at position (%d, %d)", item.Name, item.Position.X, item.Position.Y)
	}

	// Remove the entity from the floor
	for i, e := range floor.Entities {
		if e.ID == entity.ID {
			// Remove the entity by replacing it with the last one and truncating the slice
			floor.Entities[i] = floor.Entities[len(floor.Entities)-1]
			floor.Entities = floor.Entities[:len(floor.Entities)-1]
			break
		}
	}
}

// calculateGoldDrop calculates the gold drop based on difficulty and floor level
func calculateGoldDrop(difficulty models.MobDifficulty, floorLevel int) int {
	baseGold := 0
	switch difficulty {
	case models.DifficultyEasy:
		baseGold = 5
	case models.DifficultyNormal:
		baseGold = 10
	case models.DifficultyHard:
		baseGold = 20
	case models.DifficultyElite:
		baseGold = 50
	case models.DifficultyBoss:
		baseGold = 100
	}

	// Scale with floor level
	return baseGold * (1 + floorLevel/3)
}

// DescendStairs handles descending stairs
func (s *GameServer) DescendStairs() {
	currentFloor := s.Game.Dungeon.CurrentFloor
	floor := s.Game.Dungeon.Floors[currentFloor]
	playerPos := s.Game.Dungeon.PlayerPosition

	if floor.Tiles[playerPos.Y][playerPos.X].Type == models.TileStairsDown {
		if currentFloor < len(s.Game.Dungeon.Floors)-1 {
			s.Game.Dungeon.CurrentFloor++

			// Find stairs up on new floor
			newFloor := s.Game.Dungeon.Floors[s.Game.Dungeon.CurrentFloor]
			for y := 0; y < newFloor.Height; y++ {
				for x := 0; x < newFloor.Width; x++ {
					if newFloor.Tiles[y][x].Type == models.TileStairsUp {
						s.Game.Dungeon.PlayerPosition = models.Position{X: x, Y: y}
						break
					}
				}
			}

			// Get the player character for the current connection
			playerAddr := ""
			for client := range s.Clients {
				if s.Clients[client] != "" {
					playerAddr = client.RemoteAddr().String()
					break
				}
			}

			// If we have a player character, create or update the player entity on the new floor
			if playerAddr != "" && s.Game.Players[playerAddr] != nil {
				character := s.Game.Players[playerAddr]

				// Check if a player entity already exists on the new floor
				playerEntityExists := false
				for i, entity := range newFloor.Entities {
					if entity.Type == "player" {
						// Update existing player entity with character info
						newFloor.Entities[i].CharacterClass = character.CharacterClass
						newFloor.Entities[i].Health = character.Health
						newFloor.Entities[i].MaxHealth = character.MaxHealth
						newFloor.Entities[i].Position = s.Game.Dungeon.PlayerPosition
						playerEntityExists = true
						break
					}
				}

				// If no player entity exists, create one
				if !playerEntityExists {
					playerEntity := models.Entity{
						ID:             uuid.New().String(),
						Type:           "player",
						Name:           character.Name,
						Position:       s.Game.Dungeon.PlayerPosition,
						CharacterClass: character.CharacterClass,
						Health:         character.Health,
						MaxHealth:      character.MaxHealth,
					}
					newFloor.Entities = append(newFloor.Entities, playerEntity)
				}
			}

			UpdateVisibility(s.Game.Dungeon)
			log.Printf("Player descended to floor %d", s.Game.Dungeon.CurrentFloor+1)
			s.BroadcastFloorData()

			// After changing floors, respawn mobs on the new floor
			s.MobSpawner.SpawnMobsOnFloor(s.Game.Dungeon, s.Game.Dungeon.CurrentFloor)
		}
	} else {
		log.Println("No stairs to descend")
	}
}

// AscendStairs handles ascending stairs
func (s *GameServer) AscendStairs() {
	currentFloor := s.Game.Dungeon.CurrentFloor
	floor := s.Game.Dungeon.Floors[currentFloor]
	playerPos := s.Game.Dungeon.PlayerPosition

	if floor.Tiles[playerPos.Y][playerPos.X].Type == models.TileStairsUp {
		if currentFloor > 0 {
			s.Game.Dungeon.CurrentFloor--

			// Find stairs down on new floor
			newFloor := s.Game.Dungeon.Floors[s.Game.Dungeon.CurrentFloor]
			for y := 0; y < newFloor.Height; y++ {
				for x := 0; x < newFloor.Width; x++ {
					if newFloor.Tiles[y][x].Type == models.TileStairsDown {
						s.Game.Dungeon.PlayerPosition = models.Position{X: x, Y: y}
						break
					}
				}
			}

			// Get the player character for the current connection
			playerAddr := ""
			for client := range s.Clients {
				if s.Clients[client] != "" {
					playerAddr = client.RemoteAddr().String()
					break
				}
			}

			// If we have a player character, create or update the player entity on the new floor
			if playerAddr != "" && s.Game.Players[playerAddr] != nil {
				character := s.Game.Players[playerAddr]

				// Check if a player entity already exists on the new floor
				playerEntityExists := false
				for i, entity := range newFloor.Entities {
					if entity.Type == "player" {
						// Update existing player entity with character info
						newFloor.Entities[i].CharacterClass = character.CharacterClass
						newFloor.Entities[i].Health = character.Health
						newFloor.Entities[i].MaxHealth = character.MaxHealth
						newFloor.Entities[i].Position = s.Game.Dungeon.PlayerPosition
						playerEntityExists = true
						break
					}
				}

				// If no player entity exists, create one
				if !playerEntityExists {
					playerEntity := models.Entity{
						ID:             uuid.New().String(),
						Type:           "player",
						Name:           character.Name,
						Position:       s.Game.Dungeon.PlayerPosition,
						CharacterClass: character.CharacterClass,
						Health:         character.Health,
						MaxHealth:      character.MaxHealth,
					}
					newFloor.Entities = append(newFloor.Entities, playerEntity)
				}
			}

			UpdateVisibility(s.Game.Dungeon)
			log.Printf("Player ascended to floor %d", s.Game.Dungeon.CurrentFloor+1)
			s.BroadcastFloorData()

			// After changing floors, respawn mobs on the new floor
			s.MobSpawner.SpawnMobsOnFloor(s.Game.Dungeon, s.Game.Dungeon.CurrentFloor)
		}
	} else {
		log.Println("No stairs to ascend")
	}
}

// SendFloorData sends floor data to a client
func (s *GameServer) SendFloorData(conn *websocket.Conn) {
	// Get the character ID associated with this connection
	characterID, exists := s.Clients[conn]
	if !exists {
		log.Printf("No character associated with connection")
		return
	}

	// Find the dungeon the player is in
	dungeon := s.DungeonRepository.GetPlayerDungeon(characterID)
	if dungeon == nil {
		log.Printf("No dungeon found for character %s", characterID)
		return
	}

	// Get the player's position and floor
	position := dungeon.GetPlayerPosition(characterID)
	floorIndex := dungeon.GetPlayerFloor(characterID)
	floor := dungeon.Dungeon.Floors[floorIndex]

	// Get the character
	character, err := s.CharacterRepository.GetByID(characterID)
	if err != nil {
		log.Printf("Error getting character: %v", err)
		return
	}

	// Create floor message
	floorMsg := FloorMessage{
		Type:         "floor_data",
		Floor:        floor,
		PlayerPos:    *position,
		CurrentFloor: floorIndex + 1, // Convert to 1-indexed for the client
		PlayerData:   character,
		DungeonID:    dungeon.ID,
	}

	if err := conn.WriteJSON(floorMsg); err != nil {
		log.Printf("Error sending floor data: %v", err)
	}
}

// BroadcastFloorData sends floor data to all clients
func (s *GameServer) BroadcastFloorData() {
	for client, characterID := range s.Clients {
		// Only send floor data to clients that have a character ID associated with them
		if characterID != "" {
			s.SendFloorData(client)
		}
	}
}

// UpdateVisibility updates which tiles are visible to the player
func UpdateVisibility(dungeon *models.Dungeon) {
	if dungeon == nil || len(dungeon.Floors) == 0 || dungeon.CurrentFloor < 0 || dungeon.CurrentFloor >= len(dungeon.Floors) {
		log.Printf("Warning: Cannot update visibility, invalid dungeon state")
		return
	}

	currentFloor := dungeon.CurrentFloor
	floor := dungeon.Floors[currentFloor]

	if floor == nil || len(floor.Tiles) == 0 {
		log.Printf("Warning: Cannot update visibility, invalid floor state")
		return
	}

	// Make all tiles visible and explored
	for y := 0; y < floor.Height; y++ {
		if y >= len(floor.Tiles) {
			continue
		}
		for x := 0; x < floor.Width; x++ {
			if x >= len(floor.Tiles[y]) {
				continue
			}
			floor.Tiles[y][x].Visible = true
			floor.Tiles[y][x].Explored = true
		}
	}
}

// RevealEntireMap is no longer needed as UpdateVisibility now reveals everything

// Helper functions
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// HandleCreateCharacter processes character creation requests
func (s *GameServer) HandleCreateCharacter(conn *websocket.Conn, payload []byte) {
	var createCharMsg struct {
		Type      string                 `json:"type"`
		Character models.CharacterCreate `json:"character"`
	}

	if err := json.Unmarshal(payload, &createCharMsg); err != nil {
		log.Printf("Error parsing create character message: %v", err)
		sendError(conn, fmt.Sprintf("Error creating character: %v", err))
		return
	}

	// Validate character data
	if createCharMsg.Character.Name == "" {
		log.Println("Character name is required")
		sendError(conn, "Character name is required")
		return
	}

	if createCharMsg.Character.CharacterClass == "" {
		log.Println("Character class is required")
		sendError(conn, "Character class is required")
		return
	}

	// Log the character creation
	log.Printf("Creating character: %s (%s)", createCharMsg.Character.Name, createCharMsg.Character.CharacterClass)

	// Create a new character
	character := &models.Character{
		ID:             uuid.New().String(),
		Name:           createCharMsg.Character.Name,
		CharacterClass: createCharMsg.Character.CharacterClass,
		Stats:          createCharMsg.Character.Stats,
		Level:          1,
		Health:         100,
		MaxHealth:      100,
		Mana:           50,
		MaxMana:        50,
		Experience:     0,
		Gold:           10,
		Abilities:      []string{},
		Status:         []string{},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Set abilities based on class
	switch character.CharacterClass {
	case "warrior":
		character.Abilities = []string{"Slash", "Block", "Charge"}
	case "mage":
		character.Abilities = []string{"Fireball", "Frost Nova", "Blink"}
	case "rogue":
		character.Abilities = []string{"Backstab", "Stealth", "Evade"}
	case "cleric":
		character.Abilities = []string{"Heal", "Smite", "Bless"}
	case "bard":
		character.Abilities = []string{"Inspire", "Charm", "Mockery"}
	case "druid":
		character.Abilities = []string{"Shapeshift", "Entangle", "Rejuvenate"}
	case "paladin":
		character.Abilities = []string{"Smite", "Lay on Hands", "Divine Shield"}
	case "ranger":
		character.Abilities = []string{"Aimed Shot", "Track", "Animal Companion"}
	case "warlock":
		character.Abilities = []string{"Eldritch Blast", "Hex", "Dark Pact"}
	case "monk":
		character.Abilities = []string{"Flurry of Blows", "Stunning Strike", "Patient Defense"}
	case "barbarian":
		character.Abilities = []string{"Rage", "Reckless Attack", "Danger Sense"}
	case "sorcerer":
		character.Abilities = []string{"Chaos Bolt", "Metamagic", "Wild Magic"}
	default:
		character.Abilities = []string{"Attack"}
	}

	// Store the character
	if s.CharacterRepository != nil {
		err := s.CharacterRepository.Create(character)
		if err != nil {
			log.Printf("Error saving character: %v", err)
			sendError(conn, fmt.Sprintf("Error saving character: %v", err))
			return
		}
	}

	// Associate the character with the connection
	s.Game.Players[conn.RemoteAddr().String()] = character

	// Create or update player entity on the current floor with character class info
	currentFloor := s.Game.Dungeon.CurrentFloor
	floor := s.Game.Dungeon.Floors[currentFloor]
	playerPos := s.Game.Dungeon.PlayerPosition

	// Check if a player entity already exists
	playerEntityExists := false
	for i, entity := range floor.Entities {
		if entity.Type == "player" {
			// Update existing player entity with new character class
			floor.Entities[i].CharacterClass = character.CharacterClass
			floor.Entities[i].Health = character.Health
			floor.Entities[i].MaxHealth = character.MaxHealth
			playerEntityExists = true
			break
		}
	}

	// If no player entity exists, create one
	if !playerEntityExists {
		playerEntity := models.Entity{
			ID:             uuid.New().String(),
			Type:           "player",
			Name:           character.Name,
			Position:       playerPos,
			CharacterClass: character.CharacterClass,
			Health:         character.Health,
			MaxHealth:      character.MaxHealth,
		}
		floor.Entities = append(floor.Entities, playerEntity)
	}

	// Send success message
	response := map[string]interface{}{
		"type":      "character_created",
		"character": character,
		"message":   fmt.Sprintf("Character %s created successfully", character.Name),
		"timestamp": time.Now().Unix(),
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending character creation response: %v", err)
	} else {
		log.Printf("Character creation response sent successfully for %s", character.Name)
	}
}

// sendError sends an error message to the client
func sendError(conn *websocket.Conn, message string) {
	errorMsg := map[string]interface{}{
		"type":      "error",
		"message":   message,
		"timestamp": time.Now().Unix(),
	}

	if err := conn.WriteJSON(errorMsg); err != nil {
		log.Printf("Error sending error message: %v", err)
	}
}

// SaveGame saves the current game state for a player
func (s *GameServer) SaveGame(conn *websocket.Conn) {
	// Get the player character associated with this connection
	character := s.Game.Players[conn.RemoteAddr().String()]
	if character == nil {
		log.Printf("Error: No character found for connection")
		sendError(conn, "No character found to save")
		return
	}

	// Update the character's position and other state
	character.UpdatedAt = time.Now()

	// Save the character to the repository
	if s.CharacterRepository != nil {
		err := s.CharacterRepository.Update(character)
		if err != nil {
			log.Printf("Error saving character: %v", err)
			sendError(conn, fmt.Sprintf("Error saving game: %v", err))
			return
		}
	}

	// Send success message
	response := map[string]interface{}{
		"type":      "game_saved",
		"message":   fmt.Sprintf("Game saved successfully for %s", character.Name),
		"timestamp": time.Now().Unix(),
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending save game response: %v", err)
	} else {
		log.Printf("Game saved successfully for %s", character.Name)
	}
}

// InitializeGameWorld initializes the game world
func (s *GameServer) InitializeGameWorld() {
	log.Println("Initializing game world...")

	// Spawn mobs on all floors
	s.MobSpawner.SpawnMobsOnAllFloors(s.Game.Dungeon)

	// Spawn a boss on the last floor
	s.SpawnBossOnLastFloor()

	// Update visibility
	UpdateVisibility(s.Game.Dungeon)

	// Log initialization complete
	for i, floor := range s.Game.Dungeon.Floors {
		log.Printf("Floor %d: %d rooms, %d entities, %d items",
			i+1, len(floor.Rooms), len(floor.Entities), len(floor.Items))
	}

	log.Println("Game world initialization complete")
}

// SpawnBossOnLastFloor spawns a boss mob on the last floor
func (s *GameServer) SpawnBossOnLastFloor() {
	// Get the last floor
	lastFloorIndex := len(s.Game.Dungeon.Floors) - 1
	if lastFloorIndex < 0 {
		return
	}

	floor := s.Game.Dungeon.Floors[lastFloorIndex]
	floorLevel := lastFloorIndex + 1

	// Find a suitable room for the boss (preferably the last room)
	if len(floor.Rooms) == 0 {
		return
	}

	bossRoom := floor.Rooms[len(floor.Rooms)-1]

	// Get the center of the room for the boss position
	x, y := bossRoom.Center()
	bossPosition := models.Position{X: x, Y: y}

	// Choose a boss type based on the floor level
	var bossType models.MobType

	// Higher level bosses for deeper floors
	if floorLevel >= 8 {
		bossType = models.MobDragon
	} else if floorLevel >= 6 {
		bossType = models.MobLich
	} else if floorLevel >= 4 {
		bossType = models.MobOgre
	} else {
		bossType = models.MobTroll
	}

	// Create the boss with boss difficulty
	bossInstance := models.CreateMobInstance(bossType, models.DifficultyBoss, floorLevel, bossPosition)

	// Enhance the boss stats further (make it more challenging)
	bossInstance.Health *= 2
	bossInstance.MaxHealth *= 2
	bossInstance.Damage *= 2

	// Add a special prefix to the boss name
	bossInstance.Name = "Ancient " + bossInstance.Name

	// Convert to Entity for the floor
	bossEntity := models.Entity{
		ID:        bossInstance.ID,
		Type:      string(bossInstance.Type),
		Name:      bossInstance.Name,
		Position:  bossInstance.Position,
		Health:    bossInstance.Health,
		MaxHealth: bossInstance.MaxHealth,
		Damage:    bossInstance.Damage,
		Defense:   bossInstance.Defense,
		Speed:     bossInstance.Speed,
		Status:    bossInstance.Status,
	}

	// Add the boss to the floor
	floor.Entities = append(floor.Entities, bossEntity)

	log.Printf("Spawned boss %s on floor %d", bossEntity.Name, floorLevel)
}

// cleanupInactiveDungeons periodically removes inactive dungeons
func (s *GameServer) cleanupInactiveDungeons() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		count := s.DungeonRepository.CleanupInactive(MaxInactivity)
		if count > 0 {
			log.Printf("Cleaned up %d inactive dungeons", count)
		}
	}
}

// HandleCreateDungeon handles creating a new dungeon
func (s *GameServer) HandleCreateDungeon(conn *websocket.Conn, payload []byte) {
	var createDungeonMsg CreateDungeonMessage
	if err := json.Unmarshal(payload, &createDungeonMsg); err != nil {
		log.Printf("Error parsing create dungeon message: %v", err)
		sendError(conn, fmt.Sprintf("Error creating dungeon: %v", err))
		return
	}

	// Validate dungeon data
	if createDungeonMsg.Name == "" {
		log.Println("Dungeon name is required")
		sendError(conn, "Dungeon name is required")
		return
	}

	if createDungeonMsg.NumFloors <= 0 {
		log.Println("Invalid number of floors")
		sendError(conn, "Invalid number of floors")
		return
	}

	// Log the dungeon creation
	log.Printf("Creating dungeon: %s with %d floors", createDungeonMsg.Name, createDungeonMsg.NumFloors)

	// Create a new dungeon instance
	dungeon := s.DungeonRepository.Create(createDungeonMsg.Name, createDungeonMsg.NumFloors)

	// Initialize mobs on all floors
	s.initializeDungeonMobs(dungeon)

	// Send success message
	response := map[string]interface{}{
		"type":      "dungeon_created",
		"dungeonId": dungeon.ID,
		"message":   fmt.Sprintf("Dungeon %s created successfully with %d floors", createDungeonMsg.Name, createDungeonMsg.NumFloors),
		"timestamp": time.Now().Unix(),
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending dungeon creation response: %v", err)
	} else {
		log.Printf("Dungeon creation response sent successfully for %s", createDungeonMsg.Name)
	}
}

// initializeDungeonMobs initializes mobs on all floors of a dungeon
func (s *GameServer) initializeDungeonMobs(dungeon *models.DungeonInstance) {
	for i := range dungeon.Dungeon.Floors {
		s.MobSpawner.SpawnMobsOnFloor(dungeon.Dungeon, i)
	}

	// Spawn a boss on the last floor
	s.spawnBossOnLastFloor(dungeon.Dungeon)
}

// spawnBossOnLastFloor spawns a boss mob on the last floor of a dungeon
func (s *GameServer) spawnBossOnLastFloor(dungeon *models.Dungeon) {
	// Get the last floor
	lastFloorIndex := len(dungeon.Floors) - 1
	if lastFloorIndex < 0 {
		return
	}

	floor := dungeon.Floors[lastFloorIndex]
	floorLevel := lastFloorIndex + 1

	// Choose a random room (not the first one)
	roomIndex := 0
	if len(floor.Rooms) > 1 {
		roomIndex = 1 + rand.Intn(len(floor.Rooms)-1)
	}
	room := floor.Rooms[roomIndex]

	// Get boss position
	centerX, centerY := room.Center()
	bossPosition := models.Position{X: centerX, Y: centerY}

	// Create a boss entity - use a specific mob type based on floor level
	var bossType models.MobType
	switch {
	case floorLevel >= 10:
		bossType = models.MobDragon
	case floorLevel >= 7:
		bossType = models.MobLich
	case floorLevel >= 5:
		bossType = models.MobOgre
	case floorLevel >= 3:
		bossType = models.MobTroll
	default:
		bossType = models.MobGoblin
	}

	bossInstance := models.CreateMobInstance(bossType, models.DifficultyBoss, floorLevel, bossPosition)

	// Convert to Entity
	bossEntity := models.Entity{
		ID:        bossInstance.ID,
		Type:      string(bossInstance.Type),
		Name:      bossInstance.Name,
		Position:  bossInstance.Position,
		Health:    bossInstance.Health,
		MaxHealth: bossInstance.MaxHealth,
		Damage:    bossInstance.Damage,
		Defense:   bossInstance.Defense,
		Speed:     bossInstance.Speed,
		Status:    bossInstance.Status,
	}

	// Add boss to floor
	floor.Entities = append(floor.Entities, bossEntity)

	log.Printf("Spawned boss %s on floor %d", bossEntity.Name, floorLevel)
}

// HandleJoinDungeon handles joining a dungeon
func (s *GameServer) HandleJoinDungeon(conn *websocket.Conn, payload []byte) {
	log.Printf("HandleJoinDungeon called with payload: %s", string(payload))

	var joinDungeonMsg JoinDungeonMessage
	if err := json.Unmarshal(payload, &joinDungeonMsg); err != nil {
		log.Printf("Error parsing join dungeon message: %v", err)
		sendError(conn, fmt.Sprintf("Error joining dungeon: %v", err))
		return
	}

	log.Printf("Parsed join dungeon message: %+v", joinDungeonMsg)

	// Validate dungeon ID
	if joinDungeonMsg.DungeonID == "" {
		log.Println("Dungeon ID is required")
		sendError(conn, "Dungeon ID is required")
		return
	}

	// Validate character ID
	if joinDungeonMsg.CharacterID == "" {
		log.Println("Character ID is required")
		sendError(conn, "Character ID is required")
		return
	}

	// Get the dungeon
	dungeon, err := s.DungeonRepository.GetByID(joinDungeonMsg.DungeonID)
	if err != nil {
		log.Printf("Error getting dungeon: %v", err)
		sendError(conn, fmt.Sprintf("Error getting dungeon: %v", err))
		return
	}
	log.Printf("Found dungeon: %s (%s)", dungeon.Name, dungeon.ID)

	// Get the character
	character, err := s.CharacterRepository.GetByID(joinDungeonMsg.CharacterID)
	if err != nil {
		log.Printf("Error getting character: %v", err)
		sendError(conn, fmt.Sprintf("Error getting character: %v", err))
		return
	}
	log.Printf("Found character: %s (%s)", character.Name, character.ID)

	// Add the player to the dungeon (this will place them on the first floor)
	if err := s.DungeonRepository.AddPlayerToDungeon(joinDungeonMsg.DungeonID, joinDungeonMsg.CharacterID); err != nil {
		log.Printf("Error adding player to dungeon: %v", err)
		sendError(conn, fmt.Sprintf("Error adding player to dungeon: %v", err))
		return
	}
	log.Printf("Added player %s to dungeon %s", joinDungeonMsg.CharacterID, joinDungeonMsg.DungeonID)

	// Log the dungeon join
	log.Printf("Player %s joining dungeon: %s", joinDungeonMsg.CharacterID, joinDungeonMsg.DungeonID)

	// Associate the character with the connection
	s.Clients[conn] = joinDungeonMsg.CharacterID
	log.Printf("Associated character %s with WebSocket connection", joinDungeonMsg.CharacterID)

	// Get the player's position and floor (should be floor 0)
	position := dungeon.GetPlayerPosition(joinDungeonMsg.CharacterID)
	floorIndex := dungeon.GetPlayerFloor(joinDungeonMsg.CharacterID)
	floor := dungeon.Dungeon.Floors[floorIndex]
	log.Printf("Player position: (%d, %d), floor: %d", position.X, position.Y, floorIndex)

	// Send floor data
	floorData := FloorMessage{
		Type:         "floor_data",
		Floor:        floor,
		PlayerPos:    *position,
		CurrentFloor: floorIndex + 1, // Convert to 1-indexed for the client
		PlayerData:   character,
		DungeonID:    joinDungeonMsg.DungeonID,
	}

	log.Printf("Sending floor data for dungeon %s, floor %d", joinDungeonMsg.DungeonID, floorIndex)
	if err := conn.WriteJSON(floorData); err != nil {
		log.Printf("Error sending floor data: %v", err)
	} else {
		log.Printf("Successfully sent floor data")
	}

	// Send success message
	response := map[string]interface{}{
		"type":        "dungeon_joined",
		"dungeonId":   joinDungeonMsg.DungeonID,
		"characterId": joinDungeonMsg.CharacterID,
		"message":     fmt.Sprintf("Player %s joined dungeon %s on floor 1", joinDungeonMsg.CharacterID, joinDungeonMsg.DungeonID),
		"timestamp":   time.Now().Unix(),
	}

	log.Printf("Sending dungeon joined response")
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending dungeon join response: %v", err)
	} else {
		log.Printf("Player %s joined dungeon %s on floor 1", joinDungeonMsg.CharacterID, joinDungeonMsg.DungeonID)
	}
}

// HandleListDungeons handles listing all dungeons
func (s *GameServer) HandleListDungeons(conn *websocket.Conn) {
	// Get all dungeons from the repository
	dungeons := s.DungeonRepository.GetAll()

	// Create a response with dungeon list
	response := DungeonListResponse{
		Type:     "dungeon_list",
		Dungeons: make([]DungeonListItemResponse, len(dungeons)),
	}

	for i, dungeon := range dungeons {
		response.Dungeons[i] = DungeonListItemResponse{
			ID:           dungeon.ID,
			Name:         dungeon.Name,
			PlayerCount:  len(dungeon.Players),
			CreatedAt:    dungeon.CreatedAt,
			LastActivity: dungeon.LastActivity,
		}
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending dungeon list response: %v", err)
	} else {
		log.Printf("Dungeon list sent successfully")
	}
}

// HandleLeaveDungeon handles a request to leave a dungeon
func (s *GameServer) HandleLeaveDungeon(conn *websocket.Conn, payload []byte) {
	// Parse the message
	var leaveDungeonMsg LeaveDungeonMessage
	if err := json.Unmarshal(payload, &leaveDungeonMsg); err != nil {
		log.Printf("Error parsing leave_dungeon message: %v", err)
		sendError(conn, "Invalid leave_dungeon message format")
		return
	}

	// Validate the message
	if leaveDungeonMsg.DungeonID == "" {
		sendError(conn, "Missing dungeonId in leave_dungeon message")
		return
	}

	if leaveDungeonMsg.CharacterID == "" {
		sendError(conn, "Missing characterId in leave_dungeon message")
		return
	}

	// Remove the player from the dungeon
	err := s.DungeonRepository.RemovePlayerFromDungeon(leaveDungeonMsg.DungeonID, leaveDungeonMsg.CharacterID)
	if err != nil {
		log.Printf("Error removing character %s from dungeon %s: %v", leaveDungeonMsg.CharacterID, leaveDungeonMsg.DungeonID, err)
		sendError(conn, fmt.Sprintf("Error leaving dungeon: %v", err))
		return
	}

	// Remove the character association from the connection
	s.Clients[conn] = ""

	log.Printf("Character %s left dungeon %s", leaveDungeonMsg.CharacterID, leaveDungeonMsg.DungeonID)

	// Send a success response
	response := map[string]interface{}{
		"type":    "leave_dungeon_response",
		"success": true,
		"message": "Successfully left the dungeon",
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending leave_dungeon response: %v", err)
	}
}
