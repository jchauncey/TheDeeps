package models

import (
	"github.com/google/uuid"
)

// ItemType represents the type of item
type ItemType string

const (
	ItemWeapon   ItemType = "weapon"
	ItemArmor    ItemType = "armor"
	ItemPotion   ItemType = "potion"
	ItemScroll   ItemType = "scroll"
	ItemKey      ItemType = "key"
	ItemGold     ItemType = "gold"
	ItemArtifact ItemType = "artifact"
)

// Item represents an item in the game
type Item struct {
	ID          string           `json:"id"`
	Type        ItemType         `json:"type"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Value       int              `json:"value"` // Gold value
	Power       int              `json:"power"` // Damage for weapons, defense for armor, effect power for consumables
	Symbol      string           `json:"symbol"`
	Color       string           `json:"color"`
	Position    Position         `json:"position"`
	Equipped    bool             `json:"equipped"`
	ClassReq    []CharacterClass `json:"classReq,omitempty"` // Classes that can use this item
	LevelReq    int              `json:"levelReq,omitempty"` // Minimum level required to use
}

// NewWeapon creates a new weapon item
func NewWeapon(name string, damage int, value int, levelReq int, classReq []CharacterClass) *Item {
	return &Item{
		ID:          uuid.New().String(),
		Type:        ItemWeapon,
		Name:        name,
		Description: "A weapon that deals damage.",
		Value:       value,
		Power:       damage,
		Symbol:      "/",
		Color:       "#C0C0C0", // Silver
		Position:    Position{X: 0, Y: 0},
		Equipped:    false,
		ClassReq:    classReq,
		LevelReq:    levelReq,
	}
}

// NewArmor creates a new armor item
func NewArmor(name string, defense int, value int, levelReq int, classReq []CharacterClass) *Item {
	return &Item{
		ID:          uuid.New().String(),
		Type:        ItemArmor,
		Name:        name,
		Description: "Armor that provides protection.",
		Value:       value,
		Power:       defense,
		Symbol:      "[",
		Color:       "#808080", // Gray
		Position:    Position{X: 0, Y: 0},
		Equipped:    false,
		ClassReq:    classReq,
		LevelReq:    levelReq,
	}
}

// NewPotion creates a new potion item
func NewPotion(name string, power int, value int) *Item {
	return &Item{
		ID:          uuid.New().String(),
		Type:        ItemPotion,
		Name:        name,
		Description: "A potion with magical effects.",
		Value:       value,
		Power:       power,
		Symbol:      "!",
		Color:       "#FF00FF", // Magenta
		Position:    Position{X: 0, Y: 0},
	}
}

// NewGold creates a new gold item
func NewGold(amount int) *Item {
	return &Item{
		ID:          uuid.New().String(),
		Type:        ItemGold,
		Name:        "Gold",
		Description: "Shiny gold coins.",
		Value:       amount,
		Symbol:      "$",
		Color:       "#FFD700", // Gold
		Position:    Position{X: 0, Y: 0},
	}
}

// GenerateRandomItem creates a random item based on floor level
func GenerateRandomItem(floorLevel int) *Item {
	// This is a placeholder for a more sophisticated item generation system
	// In a real implementation, you would use the floor level to determine item quality

	// For now, just return a basic weapon
	return NewWeapon("Sword", 5+floorLevel, 10*floorLevel, floorLevel, nil)
}
