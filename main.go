package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/rs/cors"
)

// Game represents the game state
type Game struct {
	Dungeon *models.Dungeon
	Players map[string]*models.Character
}

// DebugMessage represents a debug message
type DebugMessage struct {
	Message   string `json:"message"`
	Level     string `json:"level"`
	Timestamp int64  `json:"timestamp"`
}

// FloorMessage represents a floor data message
type FloorMessage struct {
	Type         string          `json:"type"`
	Floor        *models.Floor   `json:"floor"`
	PlayerPos    models.Position `json:"playerPosition"`
	CurrentFloor int             `json:"currentFloor"`
}

// MoveMessage represents a movement message from the client
type MoveMessage struct {
	Type      string `json:"type"`
	Direction string `json:"direction"`
}

// ActionMessage represents an action message from the client
type ActionMessage struct {
	Type   string `json:"type"`
	Action string `json:"action"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
	clients             = make(map[*websocket.Conn]bool)
	characterRepository = repositories.NewCharacterRepository()
	game                = initGame()
)

// CreateCharacterRequest represents a request to create a character
type CreateCharacterRequest struct {
	Name           string       `json:"name"`
	CharacterClass string       `json:"characterClass"`
	Stats          models.Stats `json:"stats"`
}

// CharacterResponse represents a character response
type CharacterResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CharacterClass string `json:"characterClass"`
}

// initGame initializes the game state
func initGame() *Game {
	log.Println("Initializing game...")

	// Create a new dungeon with 10 floors
	dungeon := models.NewDungeon(10)

	log.Printf("Created dungeon with %d floors\n", len(dungeon.Floors))
	for i, floor := range dungeon.Floors {
		log.Printf("Floor %d: %d rooms, %d entities, %d items\n",
			i+1, len(floor.Rooms), len(floor.Entities), len(floor.Items))
	}

	// Mark tiles around player as visible
	updateVisibility(dungeon)

	return &Game{
		Dungeon: dungeon,
		Players: make(map[string]*models.Character),
	}
}

func main() {
	r := mux.NewRouter()

	// WebSocket endpoint
	r.HandleFunc("/ws", handleWebSocket)

	// Character endpoints
	r.HandleFunc("/character", handleCreateCharacter).Methods("POST")
	r.HandleFunc("/character/{id}", handleGetCharacter).Methods("GET")
	r.HandleFunc("/characters", handleGetCharacters).Methods("GET")

	// Dungeon endpoints
	r.HandleFunc("/dungeon/floor/{level}", handleGetFloor).Methods("GET")

	// Use CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	// Start server
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Register client
	clients[conn] = true
	defer delete(clients, conn)

	// Send welcome message
	debugMsg := DebugMessage{
		Message:   "Connected to game server",
		Level:     "info",
		Timestamp: time.Now().Unix(),
	}
	conn.WriteJSON(debugMsg)

	// Send current floor data
	sendFloorData(conn)

	// Message handling loop
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Parse the message
		var message map[string]interface{}
		if err := json.Unmarshal(p, &message); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		// Handle different message types
		msgType, ok := message["type"].(string)
		if !ok {
			log.Println("Message missing 'type' field")
			continue
		}

		switch msgType {
		case "get_floor":
			sendFloorData(conn)
		case "move":
			handleMove(conn, p)
		case "action":
			handleAction(conn, p)
		default:
			log.Printf("Unknown message type: %s", msgType)
		}
	}
}

// handleMove processes a move message
func handleMove(conn *websocket.Conn, payload []byte) {
	var moveMsg MoveMessage
	if err := json.Unmarshal(payload, &moveMsg); err != nil {
		log.Printf("Error parsing move message: %v", err)
		return
	}

	// Get current player position
	currentPos := game.Dungeon.PlayerPosition
	newPos := currentPos

	// Calculate new position based on direction
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

	// Check if the new position is valid
	if isValidMove(newPos) {
		// Update player position
		game.Dungeon.PlayerPosition = newPos

		// Update visibility
		updateVisibility(game.Dungeon)

		// Log movement
		log.Printf("Player moved %s to (%d, %d)", moveMsg.Direction, newPos.X, newPos.Y)

		// Send updated floor data to all clients
		broadcastFloorData()
	} else {
		log.Printf("Invalid move to (%d, %d)", newPos.X, newPos.Y)
	}
}

// handleAction processes an action message
func handleAction(conn *websocket.Conn, payload []byte) {
	var actionMsg ActionMessage
	if err := json.Unmarshal(payload, &actionMsg); err != nil {
		log.Printf("Error parsing action message: %v", err)
		return
	}

	// Process different actions
	switch actionMsg.Action {
	case "wait":
		log.Println("Player waited a turn")
		// TODO: Implement turn-based logic
	case "pickup":
		pickupItem()
	case "inventory":
		log.Println("Player opened inventory")
		// TODO: Implement inventory
	case "attack":
		attackEntity()
	case "use":
		log.Println("Player used an item")
		// TODO: Implement item usage
	case "descend":
		descendStairs()
	case "ascend":
		ascendStairs()
	case "menu":
		log.Println("Player opened menu")
		// TODO: Implement menu
	default:
		log.Printf("Unknown action: %s", actionMsg.Action)
	}
}

// isValidMove checks if a move is valid
func isValidMove(pos models.Position) bool {
	currentFloor := game.Dungeon.CurrentFloor
	floor := game.Dungeon.Floors[currentFloor]

	// Check if position is within bounds
	if pos.X < 0 || pos.X >= floor.Width || pos.Y < 0 || pos.Y >= floor.Height {
		return false
	}

	// Check if the tile is walkable
	tile := floor.Tiles[pos.Y][pos.X]
	return tile.Type == models.TileFloor ||
		tile.Type == models.TileStairsUp ||
		tile.Type == models.TileStairsDown ||
		tile.Type == models.TileDoor
}

// updateVisibility updates which tiles are visible to the player
func updateVisibility(dungeon *models.Dungeon) {
	currentFloor := dungeon.CurrentFloor
	floor := dungeon.Floors[currentFloor]
	playerPos := dungeon.PlayerPosition

	// Reset visibility
	for y := 0; y < floor.Height; y++ {
		for x := 0; x < floor.Width; x++ {
			floor.Tiles[y][x].Visible = false
		}
	}

	// Set tiles within visibility range to visible
	visibilityRange := 8
	for y := max(0, playerPos.Y-visibilityRange); y <= min(floor.Height-1, playerPos.Y+visibilityRange); y++ {
		for x := max(0, playerPos.X-visibilityRange); x <= min(floor.Width-1, playerPos.X+visibilityRange); x++ {
			// Simple distance check for now
			dx := playerPos.X - x
			dy := playerPos.Y - y
			distance := dx*dx + dy*dy

			if distance <= visibilityRange*visibilityRange {
				// Mark as visible and explored
				floor.Tiles[y][x].Visible = true
				floor.Tiles[y][x].Explored = true
			}
		}
	}
}

// pickupItem handles item pickup
func pickupItem() {
	currentFloor := game.Dungeon.CurrentFloor
	floor := game.Dungeon.Floors[currentFloor]
	playerPos := game.Dungeon.PlayerPosition

	// Check for items at player position
	for i, item := range floor.Items {
		if item.Position.X == playerPos.X && item.Position.Y == playerPos.Y {
			log.Printf("Player picked up %s", item.Name)

			// Remove item from floor
			floor.Items = append(floor.Items[:i], floor.Items[i+1:]...)

			// TODO: Add item to player inventory

			// Send updated floor data
			broadcastFloorData()
			return
		}
	}

	log.Println("No item to pick up")
}

// attackEntity handles attacking entities
func attackEntity() {
	currentFloor := game.Dungeon.CurrentFloor
	floor := game.Dungeon.Floors[currentFloor]
	playerPos := game.Dungeon.PlayerPosition

	// Check for entities adjacent to player
	for i, entity := range floor.Entities {
		dx := abs(entity.Position.X - playerPos.X)
		dy := abs(entity.Position.Y - playerPos.Y)

		// If entity is adjacent (including diagonals)
		if dx <= 1 && dy <= 1 {
			log.Printf("Player attacked %s", entity.Name)

			// Remove entity (for now, just kill it)
			floor.Entities = append(floor.Entities[:i], floor.Entities[i+1:]...)

			// Send updated floor data
			broadcastFloorData()
			return
		}
	}

	log.Println("No entity to attack")
}

// descendStairs handles descending stairs
func descendStairs() {
	currentFloor := game.Dungeon.CurrentFloor
	floor := game.Dungeon.Floors[currentFloor]
	playerPos := game.Dungeon.PlayerPosition

	// Check if player is on stairs down
	if floor.Tiles[playerPos.Y][playerPos.X].Type == models.TileStairsDown {
		// Check if there's a floor below
		if currentFloor < len(game.Dungeon.Floors)-1 {
			game.Dungeon.CurrentFloor++

			// Find stairs up on the new floor
			newFloor := game.Dungeon.Floors[game.Dungeon.CurrentFloor]
			for y := 0; y < newFloor.Height; y++ {
				for x := 0; x < newFloor.Width; x++ {
					if newFloor.Tiles[y][x].Type == models.TileStairsUp {
						game.Dungeon.PlayerPosition = models.Position{X: x, Y: y}
						break
					}
				}
			}

			// Update visibility
			updateVisibility(game.Dungeon)

			log.Printf("Player descended to floor %d", game.Dungeon.CurrentFloor+1)

			// Send updated floor data
			broadcastFloorData()
		}
	} else {
		log.Println("No stairs to descend")
	}
}

// ascendStairs handles ascending stairs
func ascendStairs() {
	currentFloor := game.Dungeon.CurrentFloor
	floor := game.Dungeon.Floors[currentFloor]
	playerPos := game.Dungeon.PlayerPosition

	// Check if player is on stairs up
	if floor.Tiles[playerPos.Y][playerPos.X].Type == models.TileStairsUp {
		// Check if there's a floor above
		if currentFloor > 0 {
			game.Dungeon.CurrentFloor--

			// Find stairs down on the new floor
			newFloor := game.Dungeon.Floors[game.Dungeon.CurrentFloor]
			for y := 0; y < newFloor.Height; y++ {
				for x := 0; x < newFloor.Width; x++ {
					if newFloor.Tiles[y][x].Type == models.TileStairsDown {
						game.Dungeon.PlayerPosition = models.Position{X: x, Y: y}
						break
					}
				}
			}

			// Update visibility
			updateVisibility(game.Dungeon)

			log.Printf("Player ascended to floor %d", game.Dungeon.CurrentFloor+1)

			// Send updated floor data
			broadcastFloorData()
		}
	} else {
		log.Println("No stairs to ascend")
	}
}

// sendFloorData sends the current floor data to the client
func sendFloorData(conn *websocket.Conn) {
	currentFloor := game.Dungeon.CurrentFloor
	floor := game.Dungeon.Floors[currentFloor]

	// Create floor message
	floorMsg := FloorMessage{
		Type:         "floor_data",
		Floor:        floor,
		PlayerPos:    game.Dungeon.PlayerPosition,
		CurrentFloor: currentFloor + 1, // 1-indexed for display
	}

	// Send floor data
	if err := conn.WriteJSON(floorMsg); err != nil {
		log.Printf("Error sending floor data: %v", err)
	} else {
		log.Printf("Sent floor %d data to client", currentFloor+1)
	}
}

// broadcastFloorData sends the current floor data to all clients
func broadcastFloorData() {
	for client := range clients {
		sendFloorData(client)
	}
}

func handleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var req CreateCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" || req.CharacterClass == "" {
		http.Error(w, "Name and class are required", http.StatusBadRequest)
		return
	}

	// Create character
	character := models.NewCharacter(req.Name, req.CharacterClass, req.Stats)
	if err := characterRepository.Create(character); err != nil {
		http.Error(w, "Failed to create character", http.StatusInternalServerError)
		return
	}

	// Log character creation to console
	logMessage := fmt.Sprintf("[%s] Character created: %s (ID: %s)",
		time.Now().Format("2006-01-02 15:04:05"),
		formatCharacterInfo(character),
		character.ID)
	log.Println(logMessage)

	// Return character ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id": character.ID,
	})

	// Broadcast debug message to all clients
	debugMsg := DebugMessage{
		Message:   fmt.Sprintf("New character created: %s (%s)", character.Name, character.CharacterClass),
		Level:     "info",
		Timestamp: time.Now().Unix(),
	}
	broadcastDebugMessage(debugMsg)
}

func handleGetCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	character, err := characterRepository.GetByID(id)
	if err != nil {
		if err == repositories.ErrCharacterNotFound {
			http.Error(w, "Character not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get character", http.StatusInternalServerError)
		}
		return
	}

	// Log character loading to console
	logMessage := fmt.Sprintf("[%s] Character loaded: %s (ID: %s)",
		time.Now().Format("2006-01-02 15:04:05"),
		formatCharacterInfo(character),
		character.ID)
	log.Println(logMessage)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

func handleGetCharacters(w http.ResponseWriter, r *http.Request) {
	characters := characterRepository.GetAll()

	// Log character list request to console
	logMessage := fmt.Sprintf("[%s] Character list requested: %d characters found",
		time.Now().Format("2006-01-02 15:04:05"),
		len(characters))
	log.Println(logMessage)

	// Convert to response format
	response := make([]CharacterResponse, 0, len(characters))
	for _, character := range characters {
		response = append(response, CharacterResponse{
			ID:             character.ID,
			Name:           character.Name,
			CharacterClass: character.CharacterClass,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	levelStr := vars["level"]

	var level int
	if _, err := fmt.Sscanf(levelStr, "%d", &level); err != nil {
		http.Error(w, "Invalid floor level", http.StatusBadRequest)
		return
	}

	// Adjust for 0-indexing
	level--

	if level < 0 || level >= len(game.Dungeon.Floors) {
		http.Error(w, "Floor level out of range", http.StatusBadRequest)
		return
	}

	floor := game.Dungeon.Floors[level]

	// Log floor request
	logMessage := fmt.Sprintf("[%s] Floor %d requested",
		time.Now().Format("2006-01-02 15:04:05"),
		level+1)
	log.Println(logMessage)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(floor)
}

func broadcastDebugMessage(message DebugMessage) {
	for client := range clients {
		if err := client.WriteJSON(message); err != nil {
			log.Printf("Error broadcasting message: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// Helper function to format character information for logging
func formatCharacterInfo(character *models.Character) string {
	return fmt.Sprintf("%s the %s (STR:%d, DEX:%d, CON:%d, INT:%d, WIS:%d, CHA:%d)",
		character.Name,
		character.CharacterClass,
		character.Stats.Strength,
		character.Stats.Dexterity,
		character.Stats.Constitution,
		character.Stats.Intelligence,
		character.Stats.Wisdom,
		character.Stats.Charisma)
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

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
