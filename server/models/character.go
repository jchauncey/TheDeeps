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
	Inventory      []*Item        `json:"inventory"`
	Equipment      Equipment      `json:"equipment"`
}

// Position represents a character's position on the map
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Equipment represents the items a character has equipped
type Equipment struct {
	Weapon    *Item `json:"weapon,omitempty"`
	Armor     *Item `json:"armor,omitempty"`
	Accessory *Item `json:"accessory,omitempty"`
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
		Inventory:    make([]*Item, 0),
		Equipment:    Equipment{},
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

// AddToInventory adds an item to the character's inventory
func (c *Character) AddToInventory(item *Item) {
	c.Inventory = append(c.Inventory, item)
}

// RemoveFromInventory removes an item from the character's inventory by ID
// Returns the removed item and a boolean indicating success
func (c *Character) RemoveFromInventory(itemID string) (*Item, bool) {
	for i, item := range c.Inventory {
		if item.ID == itemID {
			// Remove the item from the inventory
			removedItem := item
			c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			return removedItem, true
		}
	}
	return nil, false
}

// GetInventoryItem retrieves an item from the inventory by ID
func (c *Character) GetInventoryItem(itemID string) (*Item, bool) {
	for _, item := range c.Inventory {
		if item.ID == itemID {
			return item, true
		}
	}
	return nil, false
}

// EquipItem equips an item from the inventory
// Returns true if successful, false otherwise
func (c *Character) EquipItem(itemID string) bool {
	item, found := c.GetInventoryItem(itemID)
	if !found {
		return false
	}

	// Check if the character meets the requirements
	if item.LevelReq > c.Level {
		return false
	}

	if len(item.ClassReq) > 0 {
		classAllowed := false
		for _, allowedClass := range item.ClassReq {
			if c.Class == allowedClass {
				classAllowed = true
				break
			}
		}
		if !classAllowed {
			return false
		}
	}

	// Equip the item based on its type
	switch item.Type {
	case ItemWeapon:
		// Unequip current weapon if any
		if c.Equipment.Weapon != nil {
			c.Equipment.Weapon.Equipped = false
		}
		c.Equipment.Weapon = item
	case ItemArmor:
		// Unequip current armor if any
		if c.Equipment.Armor != nil {
			c.Equipment.Armor.Equipped = false
		}
		c.Equipment.Armor = item
	default:
		// Item type cannot be equipped
		return false
	}

	item.Equipped = true
	return true
}

// UnequipItem unequips an item and returns it to the inventory
// Returns true if successful, false otherwise
func (c *Character) UnequipItem(itemType ItemType) bool {
	switch itemType {
	case ItemWeapon:
		if c.Equipment.Weapon != nil {
			c.Equipment.Weapon.Equipped = false
			c.Equipment.Weapon = nil
			return true
		}
	case ItemArmor:
		if c.Equipment.Armor != nil {
			c.Equipment.Armor.Equipped = false
			c.Equipment.Armor = nil
			return true
		}
	case ItemArtifact:
		if c.Equipment.Accessory != nil {
			c.Equipment.Accessory.Equipped = false
			c.Equipment.Accessory = nil
			return true
		}
	}
	return false
}

// UseItem uses a consumable item from the inventory
// Returns true if successful, false otherwise
func (c *Character) UseItem(itemID string) bool {
	item, found := c.GetInventoryItem(itemID)
	if !found {
		return false
	}

	// Handle different item types
	switch item.Type {
	case ItemPotion:
		// Heal the character
		c.CurrentHP += item.Power
		if c.CurrentHP > c.MaxHP {
			c.CurrentHP = c.MaxHP
		}
		// Remove the potion from inventory after use
		c.RemoveFromInventory(itemID)
		return true
	case ItemScroll:
		// Restore mana
		c.CurrentMana += item.Power
		if c.CurrentMana > c.MaxMana {
			c.CurrentMana = c.MaxMana
		}
		// Remove the scroll from inventory after use
		c.RemoveFromInventory(itemID)
		return true
	default:
		// Item type cannot be used
		return false
	}
}

// CalculateAttackPower calculates the character's attack power based on attributes and equipment
func (c *Character) CalculateAttackPower() int {
	basePower := GetModifier(c.Attributes.Strength) + c.Level

	// Add weapon power if equipped
	if c.Equipment.Weapon != nil {
		basePower += c.Equipment.Weapon.Power
	}

	return basePower
}

// CalculateDefensePower calculates the character's defense power based on attributes and equipment
func (c *Character) CalculateDefensePower() int {
	basePower := GetModifier(c.Attributes.Constitution) + (c.Level / 2)

	// Add armor power if equipped
	if c.Equipment.Armor != nil {
		basePower += c.Equipment.Armor.Power
	}

	return basePower
}
