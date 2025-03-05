package ui

import (
	"fmt"
	"time"
)

// MessageLog manages game messages
type MessageLog struct {
	Messages    []string
	MaxMessages int
}

// NewMessageLog creates a new message log
func NewMessageLog(maxMessages int) *MessageLog {
	return &MessageLog{
		Messages:    []string{},
		MaxMessages: maxMessages,
	}
}

// AddMessage adds a message to the log
func (ml *MessageLog) AddMessage(format string, args ...interface{}) {
	// Format the message
	message := fmt.Sprintf(format, args...)

	// Add timestamp
	timestamp := time.Now().Format("15:04:05")
	formattedMessage := fmt.Sprintf("[%s] %s", timestamp, message)

	// Add to messages
	ml.Messages = append(ml.Messages, formattedMessage)

	// Trim if exceeding max messages
	if len(ml.Messages) > ml.MaxMessages {
		ml.Messages = ml.Messages[len(ml.Messages)-ml.MaxMessages:]
	}
}

// GetRecentMessages returns the n most recent messages
func (ml *MessageLog) GetRecentMessages(count int) []string {
	if count >= len(ml.Messages) {
		return ml.Messages
	}
	return ml.Messages[len(ml.Messages)-count:]
}

// GetLastMessage returns the most recent message
func (ml *MessageLog) GetLastMessage() string {
	if len(ml.Messages) == 0 {
		return ""
	}
	return ml.Messages[len(ml.Messages)-1]
}

// ClearMessages clears all messages
func (ml *MessageLog) ClearMessages() {
	ml.Messages = []string{}
}
