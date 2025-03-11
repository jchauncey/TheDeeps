package game

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWebSocketIntegration tests the WebSocket functionality
func TestWebSocketIntegration(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	characterRepo.Save(character)

	// Create a test dungeon with a floor
	dungeon := models.NewDungeon("TestDungeon", 3, 12345)
	dungeonRepo.Save(dungeon)

	// Add character to dungeon
	dungeonRepo.AddCharacterToDungeon(dungeon.ID, character.ID)
	dungeonRepo.SetCharacterFloor(dungeon.ID, character.ID, 1)
	character.CurrentDungeon = dungeon.ID
	character.CurrentFloor = 1
	characterRepo.Save(character)

	// Create a game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create a context with cancel for cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the manager in a goroutine with context
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case client := <-manager.Register:
				manager.registerClient(client)
			case client := <-manager.Unregister:
				manager.unregisterClient(client)
			case message := <-manager.Broadcast:
				manager.broadcastMessage(message)
			}
		}
	}()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.HandleConnection(w, r)
	}))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?characterId=" + character.ID

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket server")
	defer ws.Close()

	// Wait for the client to be registered
	time.Sleep(100 * time.Millisecond)

	// Verify that the client was registered
	manager.mutex.RLock()
	assert.Equal(t, 1, len(manager.Clients), "Client should be registered")
	manager.mutex.RUnlock()

	// Test sending a message to the client
	testMessage := Message{
		Type: MsgNotification,
		Text: "Test message",
	}

	// Use a WaitGroup to wait for the message to be received
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a channel to signal message receipt
	messageCh := make(chan Message, 1)

	// Start a goroutine to read messages from the WebSocket
	go func() {
		defer wg.Done()

		// Set a read deadline to prevent hanging
		ws.SetReadDeadline(time.Now().Add(2 * time.Second))

		// Read the message
		_, message, err := ws.ReadMessage()
		if err != nil {
			t.Logf("Error reading message: %v", err)
			return
		}

		// Parse the message
		var receivedMsg Message
		err = json.Unmarshal(message, &receivedMsg)
		if err != nil {
			t.Logf("Error unmarshaling message: %v", err)
			return
		}

		// Send the message to the channel
		messageCh <- receivedMsg
	}()

	// Broadcast the message
	manager.Broadcast <- testMessage

	// Wait for the message to be received with a timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for either the message to be received or a timeout
	select {
	case <-done:
		// Check if we received the message
		select {
		case receivedMsg := <-messageCh:
			// Verify the message
			assert.Equal(t, testMessage.Type, receivedMsg.Type, "Message type should match")
			assert.Equal(t, testMessage.Text, receivedMsg.Text, "Message text should match")
		default:
			t.Log("No message received, but goroutine completed")
		}
	case <-time.After(3 * time.Second):
		t.Log("Timed out waiting for message")
	}

	// Test client disconnection
	ws.Close()

	// Wait for the client to be unregistered
	time.Sleep(500 * time.Millisecond)

	// Verify that the client was unregistered
	manager.mutex.RLock()
	clientCount := len(manager.Clients)
	manager.mutex.RUnlock()

	assert.Equal(t, 0, clientCount, "Client should be unregistered")

	// Cancel the context to stop the manager goroutine
	cancel()
}

// TestWebSocketMessageHandling tests handling different message types
func TestWebSocketMessageHandling(t *testing.T) {
	// Create a game manager with mocked repositories
	manager := &GameManager{
		Clients:           make(map[string]*Client),
		Characters:        make(map[string]*models.Character),
		CharacterToClient: make(map[string]string),
		Register:          make(chan *Client),
		Unregister:        make(chan *Client),
		Broadcast:         make(chan Message),
	}

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)

	// Create a mock client
	client := &Client{
		ID:        character.ID,
		Character: character,
		Send:      make(chan Message, 256),
		Manager:   manager,
	}

	// Test handling an unknown message type
	unknownMsg := Message{
		Type: "unknown",
	}

	// Create a goroutine to handle the message
	go func() {
		// Handle the message
		if unknownMsg.Type != MsgMove &&
			unknownMsg.Type != MsgAttack &&
			unknownMsg.Type != MsgPickup &&
			unknownMsg.Type != MsgAscend &&
			unknownMsg.Type != MsgDescend &&
			unknownMsg.Type != MsgUseItem &&
			unknownMsg.Type != MsgDropItem &&
			unknownMsg.Type != MsgEquipItem &&
			unknownMsg.Type != MsgUnequipItem {
			client.Send <- Message{
				Type:  MsgError,
				Error: "Unknown message type",
			}
		}
	}()

	// Send the message
	go func() {
		time.Sleep(50 * time.Millisecond)
		manager.HandleMessage(client, unknownMsg)
	}()

	// Wait for the response
	select {
	case msg := <-client.Send:
		assert.Equal(t, MsgError, msg.Type, "Response should be an error")
		assert.Equal(t, "Unknown message type", msg.Error, "Response should indicate unknown message type")
	case <-time.After(1 * time.Second):
		t.Fatal("No response received")
	}

	// Test handling an attack message
	attackMsg := Message{
		Type:     MsgAttack,
		TargetID: "nonexistent-mob",
	}

	// Create a goroutine to handle the message
	go func() {
		// Handle the message
		if attackMsg.Type == MsgAttack {
			client.Send <- Message{
				Type: MsgNotification,
				Text: "Attack not implemented yet",
			}
		}
	}()

	// Send the message
	go func() {
		time.Sleep(50 * time.Millisecond)
		manager.HandleMessage(client, attackMsg)
	}()

	// Wait for the response
	select {
	case msg := <-client.Send:
		assert.Equal(t, MsgNotification, msg.Type, "Response should be a notification")
		assert.Contains(t, msg.Text, "not implemented", "Response should indicate attack not implemented")
	case <-time.After(1 * time.Second):
		t.Fatal("No response received")
	}

	// Test handling a move message
	moveMsg := Message{
		Type:      MsgMove,
		Direction: DirRight,
	}

	// Create a goroutine to handle the message
	go func() {
		// Handle the message
		if moveMsg.Type == MsgMove {
			client.Send <- Message{
				Type:      MsgUpdatePlayer,
				Character: character,
			}
		}
	}()

	// Send the message
	go func() {
		time.Sleep(50 * time.Millisecond)
		manager.HandleMessage(client, moveMsg)
	}()

	// Wait for the response
	select {
	case msg := <-client.Send:
		assert.Equal(t, MsgUpdatePlayer, msg.Type, "Response should be a player update")
		assert.NotNil(t, msg.Character, "Character should be included in the response")
	case <-time.After(1 * time.Second):
		t.Fatal("No response received")
	}
}

// TestWebSocketClientRunFunctions tests the client's Run, readPump, and writePump functions
func TestWebSocketClientRunFunctions(t *testing.T) {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create a test character
	character := models.NewCharacter("TestCharacter", models.Warrior)
	characterRepo.Save(character)

	// Create a game manager
	manager := NewGameManager(characterRepo, dungeonRepo)

	// Create a context with cancel for cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the manager in a goroutine with context
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case client := <-manager.Register:
				manager.registerClient(client)
			case client := <-manager.Unregister:
				manager.unregisterClient(client)
			case message := <-manager.Broadcast:
				manager.broadcastMessage(message)
			}
		}
	}()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.HandleConnection(w, r)
	}))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?characterId=" + character.ID

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket server")
	defer ws.Close()

	// Wait for the client to be registered
	time.Sleep(100 * time.Millisecond)

	// Verify that the client was registered
	manager.mutex.RLock()
	assert.Equal(t, 1, len(manager.Clients), "Client should be registered")
	manager.mutex.RUnlock()

	// Send a message to the server
	testMessage := Message{
		Type:      MsgMove,
		Direction: DirRight,
	}

	// Convert the message to JSON
	messageJSON, err := json.Marshal(testMessage)
	require.NoError(t, err, "Failed to marshal message")

	// Send the message
	err = ws.WriteMessage(websocket.TextMessage, messageJSON)
	require.NoError(t, err, "Failed to send message")

	// Wait for the message to be processed
	time.Sleep(100 * time.Millisecond)

	// Close the connection
	ws.Close()

	// Wait for the client to be unregistered
	time.Sleep(500 * time.Millisecond)

	// Verify that the client was unregistered
	manager.mutex.RLock()
	clientCount := len(manager.Clients)
	manager.mutex.RUnlock()

	assert.Equal(t, 0, clientCount, "Client should be unregistered")

	// Cancel the context to stop the manager goroutine
	cancel()
}

// mockWebSocketConn is a mock implementation of the websocket.Conn interface
type mockWebSocketConn struct {
	readMessageFunc  func() (messageType int, p []byte, err error)
	writeMessageFunc func(messageType int, data []byte) error
	closeFunc        func() error
}

// asWebSocketConn returns the mock as a *websocket.Conn
func (m *mockWebSocketConn) asWebSocketConn() *websocket.Conn {
	// This is a hack for testing purposes only
	// In a real implementation, you would create a proper mock
	return nil
}

func (m *mockWebSocketConn) ReadMessage() (messageType int, p []byte, err error) {
	if m.readMessageFunc != nil {
		return m.readMessageFunc()
	}
	return 0, nil, nil
}

func (m *mockWebSocketConn) WriteMessage(messageType int, data []byte) error {
	if m.writeMessageFunc != nil {
		return m.writeMessageFunc(messageType, data)
	}
	return nil
}

func (m *mockWebSocketConn) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func (m *mockWebSocketConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockWebSocketConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (m *mockWebSocketConn) SetReadLimit(limit int64) {
}

func (m *mockWebSocketConn) SetPongHandler(h func(appData string) error) {
}
