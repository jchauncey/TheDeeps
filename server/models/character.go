package models

import (
	"github.com/google/uuid"
)

// CharacterClass represents the available character classes
type CharacterClass string

const (
	Warrior   CharacterClass = "warrior"
	Mage      CharacterClass = "mage"
	Rogue     CharacterClass = "rogue"
	Cleric    CharacterClass = "cleric"
	Druid     CharacterClass = "druid"
	Warlock   CharacterClass = "warlock"
	Bard      CharacterClass = "bard"
	Paladin   CharacterClass = "paladin"
	Ranger    CharacterClass = "ranger"
	Monk      CharacterClass = "monk"
	Barbarian CharacterClass = "barbarian"
	Sorcerer  CharacterClass = "sorcerer"
)

// Attributes represents the D&D standard attributes
type Attributes struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

// Character represents a player character
type Character struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Class          CharacterClass `json:"class"`
	Level          int            `json:"level"`
	Experience     int            `json:"experience"`
	Attributes     Attributes     `json:"attributes"`
	MaxHP          int            `json:"maxHp"`
	CurrentHP      int            `json:"currentHp"`
	MaxMana        int            `json:"maxMana"`
	CurrentMana    int            `json:"currentMana"`
	Gold           int            `json:"gold"`
	CurrentFloor   int            `json:"currentFloor"`
	CurrentDungeon string         `json:"currentDungeon,omitempty"`
	Position       Position       `json:"position"`
	// Equipment and inventory will be added later
}

// Position represents a character's position on the map
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// NewCharacter creates a new character with default values based on class
func NewCharacter(name string, class CharacterClass) *Character {
	// Default attributes
	attributes := Attributes{
		Strength:     10,
		Dexterity:    10,
		Constitution: 10,
		Intelligence: 10,
		Wisdom:       10,
		Charisma:     10,
	}

	// Adjust attributes based on class
	switch class {
	case Warrior:
		attributes.Strength += 2
		attributes.Constitution += 2
		attributes.Intelligence -= 1
	case Mage:
		attributes.Intelligence += 3
		attributes.Wisdom += 1
		attributes.Strength -= 1
	case Rogue:
		attributes.Dexterity += 3
		attributes.Charisma += 1
		attributes.Constitution -= 1
	case Cleric:
		attributes.Wisdom += 3
		attributes.Charisma += 1
		attributes.Dexterity -= 1
	case Druid:
		attributes.Wisdom += 2
		attributes.Constitution += 1
		attributes.Charisma -= 1
	case Warlock:
		attributes.Charisma += 2
		attributes.Constitution += 1
		attributes.Wisdom -= 1
	case Bard:
		attributes.Charisma += 3
		attributes.Dexterity += 1
		attributes.Strength -= 1
	case Paladin:
		attributes.Strength += 1
		attributes.Charisma += 2
		attributes.Intelligence -= 1
	case Ranger:
		attributes.Dexterity += 2
		attributes.Wisdom += 1
		attributes.Charisma -= 1
	case Monk:
		attributes.Dexterity += 2
		attributes.Wisdom += 1
		attributes.Intelligence -= 1
	case Barbarian:
		attributes.Strength += 3
		attributes.Constitution += 1
		attributes.Intelligence -= 2
	case Sorcerer:
		attributes.Charisma += 2
		attributes.Constitution += 1
		attributes.Wisdom -= 1
	}

	// Calculate HP and Mana based on attributes and class
	maxHP := 10 + attributes.Constitution
	maxMana := 10

	switch class {
	case Warrior, Barbarian:
		maxHP += 5
		maxMana = 0
	case Mage, Sorcerer, Warlock:
		maxHP -= 2
		maxMana = 10 + attributes.Intelligence
	case Cleric, Druid:
		maxMana = 10 + attributes.Wisdom
	case Bard:
		maxMana = 8 + attributes.Charisma
	case Paladin:
		maxMana = 5 + attributes.Charisma
	}

	return &Character{
		ID:           uuid.New().String(),
		Name:         name,
		Class:        class,
		Level:        1,
		Experience:   0,
		Attributes:   attributes,
		MaxHP:        maxHP,
		CurrentHP:    maxHP,
		MaxMana:      maxMana,
		CurrentMana:  maxMana,
		Gold:         0,
		CurrentFloor: 1,
		Position:     Position{X: 0, Y: 0},
	}
}

// GetModifier calculates the attribute modifier based on D&D rules
func GetModifier(attributeValue int) int {
	return (attributeValue - 10) / 2
}

// CalculateExperienceForNextLevel calculates the experience needed for the next level
func CalculateExperienceForNextLevel(level int) int {
	return level * 1000
}

// AddExperience adds experience to the character and levels up if necessary
func (c *Character) AddExperience(exp int) bool {
	c.Experience += exp
	leveledUp := false

	nextLevelExp := CalculateExperienceForNextLevel(c.Level)
	for c.Experience >= nextLevelExp && c.Level < 20 {
		c.Level++
		leveledUp = true
		nextLevelExp = CalculateExperienceForNextLevel(c.Level)

		// Increase stats on level up
		c.MaxHP += GetModifier(c.Attributes.Constitution) + 1
		c.CurrentHP = c.MaxHP

		if c.MaxMana > 0 {
			manaIncrease := 0
			switch c.Class {
			case Mage, Sorcerer, Warlock:
				manaIncrease = GetModifier(c.Attributes.Intelligence) + 1
			case Cleric, Druid:
				manaIncrease = GetModifier(c.Attributes.Wisdom) + 1
			case Bard, Paladin:
				manaIncrease = GetModifier(c.Attributes.Charisma) + 1
			default:
				manaIncrease = 1
			}
			c.MaxMana += manaIncrease
			c.CurrentMana = c.MaxMana
		}
	}

	return leveledUp
}
