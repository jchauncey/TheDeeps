package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/repositories"
	"github.com/stretchr/testify/assert"
)

// MockWebSocketConnection is a mock implementation of the websocket connection
type MockWebSocketConnection struct {
	ReadJSONFunc  func(v interface{}) error
	WriteJSONFunc func(v interface{}) error
	CloseFunc     func() error
}

func (m *MockWebSocketConnection) ReadJSON(v interface{}) error {
	if m.ReadJSONFunc != nil {
		return m.ReadJSONFunc(v)
	}
	return nil
}

func (m *MockWebSocketConnection) WriteJSON(v interface{}) error {
	if m.WriteJSONFunc != nil {
		return m.WriteJSONFunc(v)
	}
	return nil
}

func (m *MockWebSocketConnection) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// TestNewGameHandler tests the creation of a new game handler
func TestNewGameHandler(t *testing.T) {
	handler := NewGameHandler()

	assert.NotNil(t, handler, "Handler should not be nil")
	assert.NotNil(t, handler.manager, "Game manager should not be nil")
	assert.NotNil(t, handler.characterRepo, "Character repository should not be nil")
	assert.NotNil(t, handler.upgrader, "WebSocket upgrader should not be nil")
}

// TestHandleWebSocketWithoutCharacterID tests handling a WebSocket connection without a character ID
func TestHandleWebSocketWithoutCharacterID(t *testing.T) {
	handler := NewGameHandler()

	// Create a request without a character ID
	req, err := http.NewRequest("GET", "/ws/game", nil)
	assert.NoError(t, err, "Creating request should not return an error")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.HandleWebSocket(rr, req)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Response code should be 400 Bad Request")
	assert.Equal(t, "Character ID is required\n", rr.Body.String(), "Response body should indicate character ID is required")
}

// TestHandleWebSocketWithInvalidCharacterID tests handling a WebSocket connection with an invalid character ID
func TestHandleWebSocketWithInvalidCharacterID(t *testing.T) {
	// Create a real character repository (empty)
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create a handler with the repository
	handler := &GameHandler{
		manager:       game.NewGameManager(characterRepo, dungeonRepo),
		characterRepo: characterRepo,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	// Create a request with an invalid character ID
	req, err := http.NewRequest("GET", "/ws/game?characterId=invalid-id", nil)
	assert.NoError(t, err, "Creating request should not return an error")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.HandleWebSocket(rr, req)

	// Check the response
	assert.Equal(t, http.StatusNotFound, rr.Code, "Response code should be 404 Not Found")
	assert.Equal(t, "character not found\n", rr.Body.String(), "Response body should indicate character not found")
}

// TestStartGameManager tests starting the game manager
func TestStartGameManager(t *testing.T) {
	handler := NewGameHandler()

	// This is a simple test to ensure the method doesn't panic
	assert.NotPanics(t, func() {
		handler.StartGameManager()
	}, "StartGameManager should not panic")
}
