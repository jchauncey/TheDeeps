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
	Level          int       `json:"level"`
	Health         int       `json:"health"`
	MaxHealth      int       `json:"maxHealth"`
	Mana           int       `json:"mana"`
	MaxMana        int       `json:"maxMana"`
	Experience     int       `json:"experience"`
	Gold           int       `json:"gold"`
	Abilities      []string  `json:"abilities"`
	Status         []string  `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// CharacterCreate represents the data needed to create a character
type CharacterCreate struct {
	Name           string   `json:"name"`
	CharacterClass string   `json:"characterClass"`
	Stats          Stats    `json:"stats"`
	Abilities      []string `json:"abilities"`
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
		Level:          1,
		Health:         100,
		MaxHealth:      100,
		Mana:           50,
		MaxMana:        50,
		Experience:     0,
		Gold:           10,
		Abilities:      []string{},
		Status:         []string{},
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
