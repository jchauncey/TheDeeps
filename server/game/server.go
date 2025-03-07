package game

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// Game constants
const (
	VisibilityRange = 10 // Adjust this value to change the fog of war size
)

// Debug flags
var (
	DebugMode = false // Set to true to enable debug features
)

// GameServer manages the game state and client connections
type GameServer struct {
	Game                *Game
	Clients             map[*websocket.Conn]bool
	CharacterRepository *repositories.CharacterRepository
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
	Type         string          `json:"type"`
	Floor        *models.Floor   `json:"floor"`
	PlayerPos    models.Position `json:"playerPosition"`
	CurrentFloor int             `json:"currentFloor"`
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

// NewGameServer creates a new game server instance
func NewGameServer() *GameServer {
	return &GameServer{
		Game:                initGame(),
		Clients:             make(map[*websocket.Conn]bool),
		CharacterRepository: repositories.NewCharacterRepository(),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

// Initialize game state
func initGame() *Game {
	log.Println("Initializing game...")
	dungeon := models.NewDungeon(10)

	log.Printf("Created dungeon with %d floors", len(dungeon.Floors))
	for i, floor := range dungeon.Floors {
		log.Printf("Floor %d: %d rooms, %d entities, %d items",
			i+1, len(floor.Rooms), len(floor.Entities), len(floor.Items))
	}

	game := &Game{
		Dungeon: dungeon,
		Players: make(map[string]*models.Character),
	}

	UpdateVisibility(game.Dungeon)
	return game
}

// SetupRoutes configures the HTTP routes for the server
func (s *GameServer) SetupRoutes(handler any) *mux.Router {
	h := handler.(interface {
		HandleCreateCharacter(w http.ResponseWriter, r *http.Request)
		HandleGetCharacter(w http.ResponseWriter, r *http.Request)
		HandleGetCharacters(w http.ResponseWriter, r *http.Request)
		HandleGetFloor(w http.ResponseWriter, r *http.Request)
	})

	r := mux.NewRouter()

	// WebSocket endpoint
	r.HandleFunc("/ws", s.HandleWebSocket)

	// Character endpoints
	r.HandleFunc("/character", h.HandleCreateCharacter).Methods("POST")
	r.HandleFunc("/character/{id}", h.HandleGetCharacter).Methods("GET")
	r.HandleFunc("/characters", h.HandleGetCharacters).Methods("GET")

	// Dungeon endpoints
	r.HandleFunc("/dungeon/floor/{level}", h.HandleGetFloor).Methods("GET")

	return r
}

// HandleWebSocket processes WebSocket connections
func (s *GameServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Register client
	s.Clients[conn] = true
	defer delete(s.Clients, conn)

	// Send welcome message
	conn.WriteJSON(DebugMessage{
		Message:   "Connected to game server",
		Level:     "info",
		Timestamp: time.Now().Unix(),
	})

	// Send current floor data
	s.SendFloorData(conn)

	// Message handling loop
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var message map[string]interface{}
		if err := json.Unmarshal(p, &message); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		msgType, ok := message["type"].(string)
		if !ok {
			log.Println("Message missing 'type' field")
			continue
		}

		switch msgType {
		case "get_floor":
			s.SendFloorData(conn)
		case "move":
			s.HandleMove(conn, p)
		case "action":
			s.HandleAction(conn, p)
		default:
			log.Printf("Unknown message type: %s", msgType)
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

	currentPos := s.Game.Dungeon.PlayerPosition
	newPos := currentPos

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

	// Move if valid
	if s.IsValidMove(newPos) {
		s.Game.Dungeon.PlayerPosition = newPos
		UpdateVisibility(s.Game.Dungeon)
		log.Printf("Player moved %s to (%d, %d)", moveMsg.Direction, newPos.X, newPos.Y)
		s.BroadcastFloorData()
	} else {
		log.Printf("Invalid move to (%d, %d)", newPos.X, newPos.Y)
	}
}

// HandleAction handles player actions
func (s *GameServer) HandleAction(conn *websocket.Conn, payload []byte) {
	if s.Game == nil || s.Game.Dungeon == nil {
		log.Printf("Error: Game or dungeon is nil in HandleAction")
		return
	}

	var message ActionMessage
	err := json.Unmarshal(payload, &message)
	if err != nil {
		log.Printf("Error unmarshaling action message: %v", err)
		return
	}

	log.Printf("Received action: %s", message.Action)

	switch message.Action {
	case "pickup":
		s.PickupItem()
	case "attack":
		s.AttackEntity()
	case "descend":
		s.DescendStairs()
	case "ascend":
		s.AscendStairs()
	case "toggle_debug":
		DebugMode = !DebugMode
		log.Printf("Debug mode: %v", DebugMode)
		if DebugMode {
			RevealEntireMap(s.Game.Dungeon)
		} else {
			// When turning debug mode off, update visibility normally
			UpdateVisibility(s.Game.Dungeon)
		}
		// Send updated floor data to the client
		if conn != nil {
			s.SendFloorData(conn)
		}
	default:
		log.Printf("Unknown action: %s", message.Action)
	}
}

// IsValidMove checks if a move is valid
func (s *GameServer) IsValidMove(pos models.Position) bool {
	currentFloor := s.Game.Dungeon.CurrentFloor
	floor := s.Game.Dungeon.Floors[currentFloor]

	// Check bounds
	if pos.X < 0 || pos.X >= floor.Width || pos.Y < 0 || pos.Y >= floor.Height {
		return false
	}

	// Check walkable
	tile := floor.Tiles[pos.Y][pos.X]
	return tile.Type == models.TileFloor ||
		tile.Type == models.TileStairsUp ||
		tile.Type == models.TileStairsDown ||
		tile.Type == models.TileDoor
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
			return
		}
	}
	log.Println("No entity to attack")
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

			UpdateVisibility(s.Game.Dungeon)
			log.Printf("Player descended to floor %d", s.Game.Dungeon.CurrentFloor+1)
			s.BroadcastFloorData()
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

			UpdateVisibility(s.Game.Dungeon)
			log.Printf("Player ascended to floor %d", s.Game.Dungeon.CurrentFloor+1)
			s.BroadcastFloorData()
		}
	} else {
		log.Println("No stairs to ascend")
	}
}

// SendFloorData sends floor data to a client
func (s *GameServer) SendFloorData(conn *websocket.Conn) {
	currentFloor := s.Game.Dungeon.CurrentFloor
	floor := s.Game.Dungeon.Floors[currentFloor]

	floorMsg := FloorMessage{
		Type:         "floor_data",
		Floor:        floor,
		PlayerPos:    s.Game.Dungeon.PlayerPosition,
		CurrentFloor: currentFloor + 1, // 1-indexed for display
	}

	if err := conn.WriteJSON(floorMsg); err != nil {
		log.Printf("Error sending floor data: %v", err)
	} else {
		log.Printf("Sent floor %d data to client", currentFloor+1)
	}
}

// BroadcastFloorData sends floor data to all clients
func (s *GameServer) BroadcastFloorData() {
	for client := range s.Clients {
		s.SendFloorData(client)
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
	playerPos := dungeon.PlayerPosition

	if floor == nil || len(floor.Tiles) == 0 {
		log.Printf("Warning: Cannot update visibility, invalid floor state")
		return
	}

	// Reset visibility
	for y := 0; y < floor.Height; y++ {
		if y >= len(floor.Tiles) {
			continue
		}
		for x := 0; x < floor.Width; x++ {
			if x >= len(floor.Tiles[y]) {
				continue
			}
			floor.Tiles[y][x].Visible = false
		}
	}

	// If in debug mode, make the entire floor visible
	if DebugMode {
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
		return
	}

	// Find which room the player is in
	var currentRoom *models.Room
	for i := range floor.Rooms {
		room := &floor.Rooms[i]
		if room.Contains(playerPos.X, playerPos.Y) {
			currentRoom = room
			break
		}
	}

	// If player is in a room, make the entire room visible
	if currentRoom != nil {
		// Make all tiles in the room visible and explored
		for y := currentRoom.Y; y < currentRoom.Y+currentRoom.Height; y++ {
			if y < 0 || y >= floor.Height || y >= len(floor.Tiles) {
				continue
			}
			for x := currentRoom.X; x < currentRoom.X+currentRoom.Width; x++ {
				if x < 0 || x >= floor.Width || x >= len(floor.Tiles[y]) {
					continue
				}
				floor.Tiles[y][x].Visible = true
				floor.Tiles[y][x].Explored = true
			}
		}
	} else {
		// Player is in a corridor, make a small area around them visible
		visibilityRange := 3 // Smaller range for corridors

		// Mark tiles around the player as visible
		for y := playerPos.Y - visibilityRange; y <= playerPos.Y+visibilityRange; y++ {
			if y < 0 || y >= floor.Height || y >= len(floor.Tiles) {
				continue
			}
			for x := playerPos.X - visibilityRange; x <= playerPos.X+visibilityRange; x++ {
				// Skip if out of bounds
				if x < 0 || x >= floor.Width || x >= len(floor.Tiles[y]) {
					continue
				}

				// Calculate distance from player
				dx := playerPos.X - x
				dy := playerPos.Y - y
				distance := Abs(dx) + Abs(dy) // Manhattan distance

				// Only make tiles within range visible
				if distance <= visibilityRange {
					floor.Tiles[y][x].Visible = true
					floor.Tiles[y][x].Explored = true
				}
			}
		}
	}
}

// RevealEntireMap reveals the entire current floor
func RevealEntireMap(dungeon *models.Dungeon) {
	if dungeon == nil || len(dungeon.Floors) == 0 || dungeon.CurrentFloor < 0 || dungeon.CurrentFloor >= len(dungeon.Floors) {
		log.Printf("Warning: Cannot reveal map, invalid dungeon state")
		return
	}

	currentFloor := dungeon.CurrentFloor
	floor := dungeon.Floors[currentFloor]

	if floor == nil || len(floor.Tiles) == 0 {
		log.Printf("Warning: Cannot reveal map, invalid floor state")
		return
	}

	// Mark all tiles as explored
	for y := 0; y < floor.Height; y++ {
		for x := 0; x < floor.Width; x++ {
			if y < len(floor.Tiles) && x < len(floor.Tiles[y]) {
				floor.Tiles[y][x].Explored = true
				floor.Tiles[y][x].Visible = true
			}
		}
	}
}

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
