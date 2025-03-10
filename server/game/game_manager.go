package game

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Client to server message types
	MsgMove        MessageType = "move"
	MsgAttack      MessageType = "attack"
	MsgPickup      MessageType = "pickup"
	MsgAscend      MessageType = "ascend"
	MsgDescend     MessageType = "descend"
	MsgUseItem     MessageType = "useItem"
	MsgDropItem    MessageType = "dropItem"
	MsgEquipItem   MessageType = "equipItem"
	MsgUnequipItem MessageType = "unequipItem"

	// Server to client message types
	MsgUpdateMap    MessageType = "updateMap"
	MsgUpdatePlayer MessageType = "updatePlayer"
	MsgUpdateMob    MessageType = "updateMob"
	MsgRemoveMob    MessageType = "removeMob"
	MsgAddItem      MessageType = "addItem"
	MsgRemoveItem   MessageType = "removeItem"
	MsgNotification MessageType = "notification"
	MsgFloorChange  MessageType = "floorChange"
	MsgError        MessageType = "error"
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
	// This is a placeholder for the pickup logic
	// In a real implementation, you would check if the item is valid,
	// add it to the character's inventory, remove it from the floor, etc.
	client.Send <- Message{
		Type: MsgNotification,
		Text: "Pickup not implemented yet",
	}
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

	// Find the down stairs on the new floor that corresponds to our up stairs
	// In a real implementation, you would need to track which stairs connect to which
	// For now, just place the character at the first down stairs
	if len(newFloor.DownStairs) == 0 {
		client.Send <- Message{
			Type:  MsgError,
			Error: "No down stairs found on the floor above",
		}
		return
	}

	// Update character position
	client.Character.Position.X = newFloor.DownStairs[0].X
	client.Character.Position.Y = newFloor.DownStairs[0].Y

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
		Text: "You ascend to floor " + string(client.Character.CurrentFloor),
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

	// Find the up stairs on the new floor that corresponds to our down stairs
	// In a real implementation, you would need to track which stairs connect to which
	// For now, just place the character at the first up stairs
	if len(newFloor.UpStairs) == 0 {
		client.Send <- Message{
			Type:  MsgError,
			Error: "No up stairs found on the floor below",
		}
		return
	}

	// Update character position
	client.Character.Position.X = newFloor.UpStairs[0].X
	client.Character.Position.Y = newFloor.UpStairs[0].Y

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
		Text: "You descend to floor " + string(client.Character.CurrentFloor),
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

// readPump pumps messages from the WebSocket connection to the manager
func (c *Client) readPump() {
	defer func() {
		c.Manager.Unregister <- c
		c.Connection.Close()
	}()

	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error: %v", err)
			continue
		}

		c.Manager.HandleMessage(c, msg)
	}
}

// writePump pumps messages from the manager to the WebSocket connection
func (c *Client) writePump() {
	defer func() {
		c.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// The manager closed the channel
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}

			if err := c.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("error: %v", err)
				return
			}
		}
	}
}
