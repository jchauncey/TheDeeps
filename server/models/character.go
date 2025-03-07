package models

import (
	"time"

	"github.com/google/uuid"
)

// Character represents a player character in the game
type Character struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	CharacterClass string    `json:"characterClass"`
	Stats          Stats     `json:"stats"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Stats represents character statistics
type Stats struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

// NewCharacter creates a new character with the given parameters
func NewCharacter(name, class string, stats Stats) *Character {
	now := time.Now()
	return &Character{
		ID:             uuid.New().String(),
		Name:           name,
		CharacterClass: class,
		Stats:          stats,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Update updates the character's information
func (c *Character) Update(name, class string, stats Stats) {
	c.Name = name
	c.CharacterClass = class
	c.Stats = stats
	c.UpdatedAt = time.Now()
}
