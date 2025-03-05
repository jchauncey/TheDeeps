package ui

import (
	"fmt"
	"time"
)

// MessageType represents the type of message
type MessageType int

const (
	MessageTypeGame MessageType = iota
	MessageTypeSystem
)

// Message represents a single message in the log
type Message struct {
	Text      string
	Type      MessageType
	Timestamp string
}

// MessageLog manages game messages
type MessageLog struct {
	Messages    []Message
	MaxMessages int
}

// NewMessageLog creates a new message log
func NewMessageLog(maxMessages int) *MessageLog {
	return &MessageLog{
		Messages:    []Message{},
		MaxMessages: maxMessages,
	}
}

// AddMessage adds a generic message to the log (backward compatibility)
func (ml *MessageLog) AddMessage(format string, args ...interface{}) {
	// Use game message type by default
	ml.AddGameMessage(format, args...)
}

// AddGameMessage adds a game message to the log
func (ml *MessageLog) AddGameMessage(format string, args ...interface{}) {
	// Format the message
	messageText := fmt.Sprintf(format, args...)

	// Add timestamp
	timestamp := time.Now().Format("15:04:05")

	// Create message
	message := Message{
		Text:      messageText,
		Type:      MessageTypeGame,
		Timestamp: timestamp,
	}

	// Add to messages
	ml.Messages = append(ml.Messages, message)

	// Trim if exceeding max messages
	if len(ml.Messages) > ml.MaxMessages {
		ml.Messages = ml.Messages[len(ml.Messages)-ml.MaxMessages:]
	}
}

// AddSystemMessage adds a system message to the log
func (ml *MessageLog) AddSystemMessage(format string, args ...interface{}) {
	// Format the message
	messageText := fmt.Sprintf(format, args...)

	// Add timestamp
	timestamp := time.Now().Format("15:04:05")

	// Create message
	message := Message{
		Text:      messageText,
		Type:      MessageTypeSystem,
		Timestamp: timestamp,
	}

	// Add to messages
	ml.Messages = append(ml.Messages, message)

	// Trim if exceeding max messages
	if len(ml.Messages) > ml.MaxMessages {
		ml.Messages = ml.Messages[len(ml.Messages)-ml.MaxMessages:]
	}
}

// GetRecentMessages returns the n most recent messages
func (ml *MessageLog) GetRecentMessages(count int) []Message {
	if count >= len(ml.Messages) {
		return ml.Messages
	}
	return ml.Messages[len(ml.Messages)-count:]
}

// GetLastMessage returns the most recent message as a formatted string
func (ml *MessageLog) GetLastMessage() string {
	if len(ml.Messages) == 0 {
		return ""
	}

	lastMsg := ml.Messages[len(ml.Messages)-1]
	return ml.FormatMessage(lastMsg)
}

// FormatMessage formats a message for display
func (ml *MessageLog) FormatMessage(msg Message) string {
	prefix := "[Game]"
	if msg.Type == MessageTypeSystem {
		prefix = "[System]"
	}

	return fmt.Sprintf("[%s] %s %s", msg.Timestamp, prefix, msg.Text)
}

// ClearMessages clears all messages
func (ml *MessageLog) ClearMessages() {
	ml.Messages = []Message{}
}
