package game

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/log"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// WebSocket connection parameters
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait
	maxMessageSize = 512 * 1024          // Maximum message size allowed from peer (512KB)

	// Client to server message types
	MsgMove        MessageType = "move"
	MsgAttack      MessageType = "attack"
	MsgPickup      MessageType = "pickup"
	MsgUseItem     MessageType = "useItem"
	MsgDropItem    MessageType = "dropItem"
	MsgEquipItem   MessageType = "equipItem"
	MsgUnequipItem MessageType = "unequipItem"
	MsgAscend      MessageType = "ascend"
	MsgDescend     MessageType = "descend"

	// Server to client message types
	MsgUpdateMap    MessageType = "updateMap"
	MsgUpdatePlayer MessageType = "updatePlayer"
	MsgUpdateMob    MessageType = "updateMob"
	MsgRemoveMob    MessageType = "removeMob"
	MsgAddItem      MessageType = "addItem"
	MsgRemoveItem   MessageType = "removeItem"
	MsgNotification MessageType = "notification"
	MsgFloorUpdate  MessageType = "floorUpdate"
	MsgFloorChange  MessageType = "floorChange"
	MsgError        MessageType = "error"
	MsgInitialState MessageType = "initialState"
)

// Direction represents a movement direction
type Direction string

const (
	DirUp    Direction = "up"
	DirDown  Direction = "down"
	DirLeft  Direction = "left"
	DirRight Direction = "right"
)

// Message represents a WebSocket message
type Message struct {
	Type        MessageType       `json:"type"`
	CharacterID string            `json:"characterId,omitempty"`
	Direction   Direction         `json:"direction,omitempty"`
	TargetID    string            `json:"targetId,omitempty"`
	ItemID      string            `json:"itemId,omitempty"`
	Floor       *models.Floor     `json:"floor,omitempty"`
	Character   *models.Character `json:"character,omitempty"`
	Mob         *models.Mob       `json:"mob,omitempty"`
	Item        *models.Item      `json:"item,omitempty"`
	Text        string            `json:"text,omitempty"`
	Error       string            `json:"error,omitempty"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID         string
	Connection *websocket.Conn
	Character  *models.Character
	Send       chan Message
	Manager    *GameManager
}

// GameManager handles the game state and WebSocket connections
type GameManager struct {
	Clients           map[string]*Client
	Characters        map[string]*models.Character
	CharacterToClient map[string]string
	Register          chan *Client
	Unregister        chan *Client
	Broadcast         chan Message
	CharacterRepo     *repositories.CharacterRepository
	DungeonRepo       *repositories.DungeonRepository
	MapGenerator      *MapGenerator
	mutex             sync.RWMutex
}

// NewGameManager creates a new game manager
func NewGameManager(characterRepo *repositories.CharacterRepository, dungeonRepo *repositories.DungeonRepository) *GameManager {
	return &GameManager{
		Clients:           make(map[string]*Client),
		Characters:        make(map[string]*models.Character),
		CharacterToClient: make(map[string]string),
		Register:          make(chan *Client),
		Unregister:        make(chan *Client),
		Broadcast:         make(chan Message),
		CharacterRepo:     characterRepo,
		DungeonRepo:       dungeonRepo,
		MapGenerator:      NewMapGenerator(time.Now().UnixNano()),
	}
}

// Start starts the game manager
func (manager *GameManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.registerClient(client)
		case client := <-manager.Unregister:
			manager.unregisterClient(client)
		case message := <-manager.Broadcast:
			manager.broadcastMessage(message)
		}
	}
}

// registerClient registers a new client
func (manager *GameManager) registerClient(client *Client) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.Clients[client.ID] = client

	if client.Character != nil {
		manager.Characters[client.Character.ID] = client.Character
		manager.CharacterToClient[client.Character.ID] = client.ID

		// Send initial game state to the client
		if client.Character.CurrentDungeon != "" {
			// Get the current floor
			floorLevel := client.Character.CurrentFloor
			floor, err := manager.DungeonRepo.GetFloor(client.Character.CurrentDungeon, floorLevel)
			if err == nil {
				// Send the floor data
				client.Send <- Message{
					Type:  MsgFloorChange,
					Floor: floor,
				}

				// Send the character data
				client.Send <- Message{
					Type:      MsgUpdatePlayer,
					Character: client.Character,
				}
			}
		}
	}
}

// unregisterClient unregisters a client
func (manager *GameManager) unregisterClient(client *Client) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if _, ok := manager.Clients[client.ID]; ok {
		close(client.Send)
		delete(manager.Clients, client.ID)

		if client.Character != nil {
			delete(manager.CharacterToClient, client.Character.ID)

			// Save character state
			manager.CharacterRepo.Save(client.Character)
		}
	}
}

// broadcastMessage broadcasts a message to all clients
func (manager *GameManager) broadcastMessage(message Message) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, client := range manager.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(manager.Clients, client.ID)
		}
	}
}

// HandleMessage handles a message from a client
func (manager *GameManager) HandleMessage(client *Client, message Message) {
	// Validate that the character ID in the message matches the client's character
	if message.CharacterID != "" && client.Character != nil && message.CharacterID != client.Character.ID {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Invalid character ID",
		}
		return
	}

	// Handle different message types
	switch message.Type {
	case MsgMove:
		manager.handleMove(client, message)
	case MsgAttack:
		manager.handleAttack(client, message)
	case MsgPickup:
		manager.handlePickup(client, message)
	case MsgAscend:
		manager.handleAscend(client, message)
	case MsgDescend:
		manager.handleDescend(client, message)
	case MsgUseItem:
		manager.handleUseItem(client, message)
	case MsgDropItem:
		manager.handleDropItem(client, message)
	case MsgEquipItem:
		manager.handleEquipItem(client, message)
	case MsgUnequipItem:
		manager.handleUnequipItem(client, message)
	default:
		client.Send <- Message{
			Type:  MsgError,
			Error: "Unknown message type",
		}
	}
}

// handleMove handles a move message
func (manager *GameManager) handleMove(client *Client, message Message) {
	if client.Character == nil || client.Character.CurrentDungeon == "" {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Character not in a dungeon",
		}
		return
	}

	// Get the current floor
	_, err := manager.DungeonRepo.GetByID(client.Character.CurrentDungeon)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Dungeon not found",
		}
		return
	}

	floor, err := manager.DungeonRepo.GetFloor(client.Character.CurrentDungeon, client.Character.CurrentFloor)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Floor not found",
		}
		return
	}

	// Calculate new position
	newX, newY := client.Character.Position.X, client.Character.Position.Y

	switch message.Direction {
	case DirUp:
		newY--
	case DirDown:
		newY++
	case DirLeft:
		newX--
	case DirRight:
		newX++
	}

	// Check if the new position is valid
	if newX < 0 || newX >= floor.Width || newY < 0 || newY >= floor.Height {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Invalid move: out of bounds",
		}
		return
	}

	// Check if the tile is walkable
	if !floor.Tiles[newY][newX].Walkable {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Invalid move: tile not walkable",
		}
		return
	}

	// Check if there's a mob on the tile
	if floor.Tiles[newY][newX].MobID != "" {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Invalid move: tile occupied by mob",
		}
		return
	}

	// Update the old tile
	oldX, oldY := client.Character.Position.X, client.Character.Position.Y
	floor.Tiles[oldY][oldX].Character = ""

	// Update the new tile
	floor.Tiles[newY][newX].Character = client.Character.ID

	// Update character position
	client.Character.Position.X = newX
	client.Character.Position.Y = newY

	// Save the character
	manager.CharacterRepo.Save(client.Character)

	// Notify the client
	client.Send <- Message{
		Type:      MsgUpdatePlayer,
		Character: client.Character,
	}

	// Check if the character is on stairs
	if floor.Tiles[newY][newX].Type == models.TileUpStairs {
		client.Send <- Message{
			Type: MsgNotification,
			Text: "You are standing on stairs leading up. Press 'u' to ascend.",
		}
	} else if floor.Tiles[newY][newX].Type == models.TileDownStairs {
		client.Send <- Message{
			Type: MsgNotification,
			Text: "You are standing on stairs leading down. Press 'd' to descend.",
		}
	}

	// Check if there's an item on the tile
	if floor.Tiles[newY][newX].ItemID != "" {
		item := floor.Items[floor.Tiles[newY][newX].ItemID]
		client.Send <- Message{
			Type: MsgNotification,
			Text: "You see a " + item.Name + " here. Press 'g' to pick it up.",
		}
	}
}

// handleAttack handles an attack message
func (manager *GameManager) handleAttack(client *Client, message Message) {
	// This is a placeholder for the attack logic
	// In a real implementation, you would check if the target is valid,
	// calculate damage, update mob health, etc.
	client.Send <- Message{
		Type: MsgNotification,
		Text: "Attack not implemented yet",
	}
}

// handlePickup handles a pickup message
func (manager *GameManager) handlePickup(client *Client, message Message) {
	// Get the character
	character := client.Character
	if character == nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Character not found",
		}
		return
	}

	// Get the item ID from the message
	itemID := message.ItemID
	if itemID == "" {
		client.Send <- Message{
			Type:  MsgError,
			Error: "No item specified",
		}
		return
	}

	// Get the current floor
	dungeon, err := manager.DungeonRepo.GetByID(character.CurrentDungeon)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Dungeon not found",
		}
		return
	}

	floor := dungeon.FloorData[character.CurrentFloor]
	if floor == nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Floor not found",
		}
		return
	}

	// Find the item on the floor
	item, exists := floor.Items[itemID]
	if !exists {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Item not found on this floor",
		}
		return
	}

	// Check if the character is at the same position as the item
	if character.Position.X != item.Position.X || character.Position.Y != item.Position.Y {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Item is not at your position",
		}
		return
	}

	// Create a pointer to the item for adding to inventory
	itemPtr := &item

	// Check if adding this item would exceed the character's weight limit
	if !character.CanAddItem(itemPtr) {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Cannot pick up item: weight limit exceeded",
		}
		return
	}

	// Add the item to the character's inventory
	success := character.AddToInventory(itemPtr)
	if !success {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Failed to add item to inventory",
		}
		return
	}

	// Remove the item from the floor
	delete(floor.Items, itemID)

	// Save the updated character
	err = manager.CharacterRepo.Save(character)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Failed to save character",
		}
		return
	}

	// Save the updated dungeon
	err = manager.DungeonRepo.Save(dungeon)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Failed to save dungeon",
		}
		return
	}

	// Send success message to the client
	client.Send <- Message{
		Type:      MsgNotification,
		Text:      "You picked up " + item.Name,
		Character: character,
		Item:      itemPtr,
	}

	// Broadcast the floor update to all clients on this floor
	manager.BroadcastFloorUpdate(character.CurrentDungeon, character.CurrentFloor)
}

// handleAscend handles an ascend message
func (manager *GameManager) handleAscend(client *Client, message Message) {
	if client.Character == nil || client.Character.CurrentDungeon == "" {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Character not in a dungeon",
		}
		return
	}

	// Get the current floor
	_, err := manager.DungeonRepo.GetByID(client.Character.CurrentDungeon)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Dungeon not found",
		}
		return
	}

	floor, err := manager.DungeonRepo.GetFloor(client.Character.CurrentDungeon, client.Character.CurrentFloor)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Floor not found",
		}
		return
	}

	// Check if the character is on up stairs
	x, y := client.Character.Position.X, client.Character.Position.Y
	if floor.Tiles[y][x].Type != models.TileUpStairs {
		client.Send <- Message{
			Type:  MsgError,
			Error: "You are not on stairs leading up",
		}
		return
	}

	// Check if we're already at the top floor
	if client.Character.CurrentFloor == 1 {
		client.Send <- Message{
			Type:  MsgError,
			Error: "You are already at the top floor",
		}
		return
	}

	// Update the old tile
	floor.Tiles[y][x].Character = ""

	// Update character floor
	client.Character.CurrentFloor--

	// Get the new floor
	newFloor, err := manager.DungeonRepo.GetFloor(client.Character.CurrentDungeon, client.Character.CurrentFloor)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Floor not found",
		}
		return
	}

	// Check if we're going to floor 1 (entrance floor)
	if client.Character.CurrentFloor == 1 {
		// Find the entrance room
		var entranceRoom *models.Room
		for i := range newFloor.Rooms {
			if newFloor.Rooms[i].Type == models.RoomEntrance {
				entranceRoom = &newFloor.Rooms[i]
				break
			}
		}

		if entranceRoom != nil {
			// Find the down stairs in the entrance room
			var stairsX, stairsY int
			var stairsFound bool

			for _, downStair := range newFloor.DownStairs {
				// Check if this stair is in the entrance room
				if downStair.X >= entranceRoom.X && downStair.X < entranceRoom.X+entranceRoom.Width &&
					downStair.Y >= entranceRoom.Y && downStair.Y < entranceRoom.Y+entranceRoom.Height {
					stairsX = downStair.X
					stairsY = downStair.Y
					stairsFound = true
					break
				}
			}

			if stairsFound {
				// Place character one tile away from the stairs
				// Try different directions until we find a walkable tile
				directions := []struct{ dx, dy int }{
					{0, 1}, {1, 0}, {0, -1}, {-1, 0}, // Cardinal directions
					{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // Diagonals
				}

				placed := false
				for _, dir := range directions {
					newX, newY := stairsX+dir.dx, stairsY+dir.dy

					// Check if position is valid and walkable
					if newX >= 0 && newX < newFloor.Width &&
						newY >= 0 && newY < newFloor.Height &&
						newFloor.Tiles[newY][newX].Walkable &&
						newFloor.Tiles[newY][newX].Character == "" &&
						newFloor.Tiles[newY][newX].MobID == "" {
						client.Character.Position.X = newX
						client.Character.Position.Y = newY
						placed = true
						break
					}
				}

				// If we couldn't place adjacent to stairs, use center of room
				if !placed {
					client.Character.Position.X = entranceRoom.X + entranceRoom.Width/2
					client.Character.Position.Y = entranceRoom.Y + entranceRoom.Height/2
				}
			} else {
				// Fallback to center of entrance room
				client.Character.Position.X = entranceRoom.X + entranceRoom.Width/2
				client.Character.Position.Y = entranceRoom.Y + entranceRoom.Height/2
			}
		} else if len(newFloor.DownStairs) > 0 {
			// Fallback: Place character at the first down stairs
			client.Character.Position.X = newFloor.DownStairs[0].X
			client.Character.Position.Y = newFloor.DownStairs[0].Y
		} else {
			// Emergency fallback: find any walkable tile
			for y := 0; y < newFloor.Height; y++ {
				for x := 0; x < newFloor.Width; x++ {
					if newFloor.Tiles[y][x].Walkable &&
						newFloor.Tiles[y][x].Character == "" &&
						newFloor.Tiles[y][x].MobID == "" {
						client.Character.Position.X = x
						client.Character.Position.Y = y
						break
					}
				}
				if client.Character.Position.X != 0 || client.Character.Position.Y != 0 {
					break
				}
			}
		}
	} else {
		// For other floors, find a room with down stairs
		// First try to find the down stairs that correspond to our up stairs
		if len(newFloor.DownStairs) == 0 {
			client.Send <- Message{
				Type:  MsgError,
				Error: "No down stairs found on the floor above",
			}
			return
		}

		// Find which room contains the down stairs
		var stairsRoom *models.Room
		var stairsPosition models.Position

		// Use the first down stairs as default
		stairsPosition = newFloor.DownStairs[0]

		// Try to find which room contains these stairs
		for i := range newFloor.Rooms {
			room := &newFloor.Rooms[i]
			if stairsPosition.X >= room.X && stairsPosition.X < room.X+room.Width &&
				stairsPosition.Y >= room.Y && stairsPosition.Y < room.Y+room.Height {
				stairsRoom = room
				break
			}
		}

		if stairsRoom != nil {
			// Place character one tile away from the stairs
			directions := []struct{ dx, dy int }{
				{0, 1}, {1, 0}, {0, -1}, {-1, 0}, // Cardinal directions
				{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // Diagonals
			}

			placed := false
			for _, dir := range directions {
				newX, newY := stairsPosition.X+dir.dx, stairsPosition.Y+dir.dy

				// Check if position is valid and walkable
				if newX >= 0 && newX < newFloor.Width &&
					newY >= 0 && newY < newFloor.Height &&
					newFloor.Tiles[newY][newX].Walkable &&
					newFloor.Tiles[newY][newX].Character == "" &&
					newFloor.Tiles[newY][newX].MobID == "" {
					client.Character.Position.X = newX
					client.Character.Position.Y = newY
					placed = true
					break
				}
			}

			// If we couldn't place adjacent to stairs, use center of room
			if !placed {
				client.Character.Position.X = stairsRoom.X + stairsRoom.Width/2
				client.Character.Position.Y = stairsRoom.Y + stairsRoom.Height/2
			}
		} else {
			// Fallback: Place character at the stairs position
			client.Character.Position.X = stairsPosition.X
			client.Character.Position.Y = stairsPosition.Y
		}
	}

	// Update the new tile
	newFloor.Tiles[client.Character.Position.Y][client.Character.Position.X].Character = client.Character.ID

	// Update the character's floor in the dungeon
	manager.DungeonRepo.SetCharacterFloor(client.Character.CurrentDungeon, client.Character.ID, client.Character.CurrentFloor)

	// Save the character
	manager.CharacterRepo.Save(client.Character)

	// Notify the client
	client.Send <- Message{
		Type:  MsgFloorChange,
		Floor: newFloor,
	}

	client.Send <- Message{
		Type:      MsgUpdatePlayer,
		Character: client.Character,
	}

	client.Send <- Message{
		Type: MsgNotification,
		Text: "You ascend to floor " + strconv.Itoa(client.Character.CurrentFloor),
	}
}

// handleDescend handles a descend message
func (manager *GameManager) handleDescend(client *Client, message Message) {
	if client.Character == nil || client.Character.CurrentDungeon == "" {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Character not in a dungeon",
		}
		return
	}

	// Get the current floor
	dungeon, err := manager.DungeonRepo.GetByID(client.Character.CurrentDungeon)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Dungeon not found",
		}
		return
	}

	floor, err := manager.DungeonRepo.GetFloor(client.Character.CurrentDungeon, client.Character.CurrentFloor)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Floor not found",
		}
		return
	}

	// Check if the character is on down stairs
	x, y := client.Character.Position.X, client.Character.Position.Y
	if floor.Tiles[y][x].Type != models.TileDownStairs {
		client.Send <- Message{
			Type:  MsgError,
			Error: "You are not on stairs leading down",
		}
		return
	}

	// Check if we're already at the bottom floor
	if client.Character.CurrentFloor == dungeon.Floors {
		client.Send <- Message{
			Type:  MsgError,
			Error: "You are already at the bottom floor",
		}
		return
	}

	// Update the old tile
	floor.Tiles[y][x].Character = ""

	// Update character floor
	client.Character.CurrentFloor++

	// Get the new floor
	newFloor, err := manager.DungeonRepo.GetFloor(client.Character.CurrentDungeon, client.Character.CurrentFloor)
	if err != nil {
		client.Send <- Message{
			Type:  MsgError,
			Error: "Floor not found",
		}
		return
	}

	// Find a safe room with up stairs if possible
	var safeRoom *models.Room
	for i := range newFloor.Rooms {
		if newFloor.Rooms[i].Type == models.RoomSafe {
			safeRoom = &newFloor.Rooms[i]
			break
		}
	}

	// Place character at appropriate position
	if safeRoom != nil {
		// Place character in the center of the safe room, slightly offset from the stairs
		// First find the up stairs in this room
		var stairsX, stairsY int
		var stairsFound bool

		for _, upStair := range newFloor.UpStairs {
			// Check if this stair is in the safe room
			if upStair.X >= safeRoom.X && upStair.X < safeRoom.X+safeRoom.Width &&
				upStair.Y >= safeRoom.Y && upStair.Y < safeRoom.Y+safeRoom.Height {
				stairsX = upStair.X
				stairsY = upStair.Y
				stairsFound = true
				break
			}
		}

		if stairsFound {
			// Place character one tile away from the stairs
			// Try different directions until we find a walkable tile
			directions := []struct{ dx, dy int }{
				{0, 1}, {1, 0}, {0, -1}, {-1, 0}, // Cardinal directions
				{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // Diagonals
			}

			placed := false
			for _, dir := range directions {
				newX, newY := stairsX+dir.dx, stairsY+dir.dy

				// Check if position is valid and walkable
				if newX >= 0 && newX < newFloor.Width &&
					newY >= 0 && newY < newFloor.Height &&
					newFloor.Tiles[newY][newX].Walkable &&
					newFloor.Tiles[newY][newX].Character == "" &&
					newFloor.Tiles[newY][newX].MobID == "" {
					client.Character.Position.X = newX
					client.Character.Position.Y = newY
					placed = true
					break
				}
			}

			// If we couldn't place adjacent to stairs, use center of room
			if !placed {
				client.Character.Position.X = safeRoom.X + safeRoom.Width/2
				client.Character.Position.Y = safeRoom.Y + safeRoom.Height/2
			}
		} else {
			// Fallback to center of safe room
			client.Character.Position.X = safeRoom.X + safeRoom.Width/2
			client.Character.Position.Y = safeRoom.Y + safeRoom.Height/2
		}
	} else if len(newFloor.UpStairs) > 0 {
		// Fallback: Place character at the first up stairs
		client.Character.Position.X = newFloor.UpStairs[0].X
		client.Character.Position.Y = newFloor.UpStairs[0].Y
	} else {
		// Emergency fallback: find any walkable tile
		for y := 0; y < newFloor.Height; y++ {
			for x := 0; x < newFloor.Width; x++ {
				if newFloor.Tiles[y][x].Walkable &&
					newFloor.Tiles[y][x].Character == "" &&
					newFloor.Tiles[y][x].MobID == "" {
					client.Character.Position.X = x
					client.Character.Position.Y = y
					break
				}
			}
			if client.Character.Position.X != 0 || client.Character.Position.Y != 0 {
				break
			}
		}
	}

	// Update the new tile
	newFloor.Tiles[client.Character.Position.Y][client.Character.Position.X].Character = client.Character.ID

	// Update the character's floor in the dungeon
	manager.DungeonRepo.SetCharacterFloor(client.Character.CurrentDungeon, client.Character.ID, client.Character.CurrentFloor)

	// Save the character
	manager.CharacterRepo.Save(client.Character)

	// Notify the client
	client.Send <- Message{
		Type:  MsgFloorChange,
		Floor: newFloor,
	}

	client.Send <- Message{
		Type:      MsgUpdatePlayer,
		Character: client.Character,
	}

	client.Send <- Message{
		Type: MsgNotification,
		Text: "You descend to floor " + strconv.Itoa(client.Character.CurrentFloor),
	}
}

// handleUseItem handles a use item message
func (manager *GameManager) handleUseItem(client *Client, message Message) {
	// This is a placeholder for the use item logic
	client.Send <- Message{
		Type: MsgNotification,
		Text: "Use item not implemented yet",
	}
}

// handleDropItem handles a drop item message
func (manager *GameManager) handleDropItem(client *Client, message Message) {
	// This is a placeholder for the drop item logic
	client.Send <- Message{
		Type: MsgNotification,
		Text: "Drop item not implemented yet",
	}
}

// handleEquipItem handles an equip item message
func (manager *GameManager) handleEquipItem(client *Client, message Message) {
	// This is a placeholder for the equip item logic
	client.Send <- Message{
		Type: MsgNotification,
		Text: "Equip item not implemented yet",
	}
}

// handleUnequipItem handles an unequip item message
func (manager *GameManager) handleUnequipItem(client *Client, message Message) {
	// This is a placeholder for the unequip item logic
	client.Send <- Message{
		Type: MsgNotification,
		Text: "Unequip item not implemented yet",
	}
}

// Run starts the client's read and write pumps
func (c *Client) Run() {
	go c.readPump()
	go c.writePump()
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.Manager.Unregister <- c
		c.Connection.Close()
	}()

	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error {
		c.Connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("error: %v", err)
			}
			break
		}

		// Process the message
		// TODO: Parse and handle the message
		// ...
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()

	for {
		select {
		case _, ok := <-c.Send:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The manager closed the channel.
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write the message to the websocket
			// TODO: Implement message writing
			// ...

		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// BroadcastFloorUpdate broadcasts a floor update to all clients on the specified floor
func (gm *GameManager) BroadcastFloorUpdate(dungeonID string, floorLevel int) {
	// Get the floor using the repository
	floor, err := gm.DungeonRepo.GetFloor(dungeonID, floorLevel)
	if err != nil {
		log.Error("Failed to get floor: %v", err)
		return
	}

	// Find all clients on this floor
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()

	for _, client := range gm.Clients {
		if client.Character != nil &&
			client.Character.CurrentDungeon == dungeonID &&
			client.Character.CurrentFloor == floorLevel {

			// Send floor update to this client
			client.Send <- Message{
				Type:  MsgFloorChange,
				Floor: floor,
			}
		}
	}
}

// HandleConnection handles a new WebSocket connection
func (gm *GameManager) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Failed to upgrade connection: %v", err)
		return
	}

	// Get character ID from query parameters
	characterID := r.URL.Query().Get("characterId")
	if characterID == "" {
		log.Warn("Character ID not provided")
		conn.Close()
		return
	}

	// Get the character
	character, err := gm.CharacterRepo.GetByID(characterID)
	if err != nil {
		log.Warn("Character not found: %s", characterID)
		conn.Close()
		return
	}

	// Create a new client
	client := &Client{
		ID:         characterID,
		Connection: conn,
		Character:  character,
		Send:       make(chan Message, 256),
		Manager:    gm,
	}

	// Register the client
	gm.Register <- client

	// Start the client's read and write pumps
	go client.writePump()
	go client.readPump()

	// Send the initial game state
	client.Send <- Message{
		Type:      MsgInitialState,
		Character: character,
	}
}
