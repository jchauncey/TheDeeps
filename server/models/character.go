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
	Skills         *Skills        `json:"skills"`
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

	// Initialize skills based on class
	skills := NewSkills(class)

	return &Character{
		ID:           uuid.New().String(),
		Name:         name,
		Class:        class,
		Level:        1,
		Experience:   0,
		Attributes:   attributes,
		Skills:       skills,
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

// NewCharacterWithSkills creates a new character with default values based on class
// This version includes the Skills field initialization
func NewCharacterWithSkills(name string, class CharacterClass) *Character {
	// Create a character using the existing NewCharacter function
	character := NewCharacter(name, class)

	// Initialize skills based on class
	character.Skills = NewSkills(class)

	return character
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

// CalculateInventoryWeight calculates the total weight of all items in the character's inventory
func (c *Character) CalculateInventoryWeight() float64 {
	totalWeight := 0.0
	for _, item := range c.Inventory {
		if !item.Equipped { // Don't count equipped items in inventory weight
			totalWeight += item.Weight
		}
	}
	return totalWeight
}

// CalculateEquipmentWeight calculates the total weight of all equipped items
func (c *Character) CalculateEquipmentWeight() float64 {
	totalWeight := 0.0
	if c.Equipment.Weapon != nil {
		totalWeight += c.Equipment.Weapon.Weight
	}
	if c.Equipment.Armor != nil {
		totalWeight += c.Equipment.Armor.Weight
	}
	if c.Equipment.Accessory != nil {
		totalWeight += c.Equipment.Accessory.Weight
	}
	return totalWeight
}

// CalculateTotalWeight calculates the total weight of inventory and equipped items
func (c *Character) CalculateTotalWeight() float64 {
	return c.CalculateInventoryWeight() + c.CalculateEquipmentWeight()
}

// CalculateWeightLimit returns the maximum weight the character can carry based on strength
func (c *Character) CalculateWeightLimit() float64 {
	// Base weight limit is 50 pounds
	baseLimit := 50.0

	// Each point of strength above 10 adds 10 pounds to the limit
	strengthModifier := float64(c.Attributes.Strength - 10)
	if strengthModifier < 0 {
		// For strength below 10, each point reduces the limit by 0.5 pounds
		return baseLimit + (strengthModifier * 0.5 * 10)
	}

	// For strength above 10, each point adds 10 pounds
	return baseLimit + (strengthModifier * 10)
}

// IsOverEncumbered returns true if the character is carrying more than their weight limit
func (c *Character) IsOverEncumbered() bool {
	return c.CalculateTotalWeight() > c.CalculateWeightLimit()
}

// GetEncumbranceLevel returns the character's encumbrance level
// 0 = Not encumbered, 1 = Lightly encumbered, 2 = Heavily encumbered, 3 = Over encumbered
func (c *Character) GetEncumbranceLevel() int {
	totalWeight := c.CalculateTotalWeight()
	weightLimit := c.CalculateWeightLimit()

	if totalWeight <= weightLimit*0.5 {
		return 0 // Not encumbered
	} else if totalWeight <= weightLimit*0.75 {
		return 1 // Lightly encumbered
	} else if totalWeight <= weightLimit {
		return 2 // Heavily encumbered
	} else {
		return 3 // Over encumbered
	}
}

// CanAddItem checks if an item can be added to the inventory without exceeding weight limit
func (c *Character) CanAddItem(item *Item) bool {
	return (c.CalculateTotalWeight() + item.Weight) <= c.CalculateWeightLimit()
}

// AddToInventory adds an item to the character's inventory if weight limit allows
// Returns true if successful, false if weight limit would be exceeded
func (c *Character) AddToInventory(item *Item) bool {
	if !c.CanAddItem(item) {
		return false
	}

	c.Inventory = append(c.Inventory, item)
	return true
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

// CalculateBaseAC calculates the base armor class without equipment
func (c *Character) CalculateBaseAC() int {
	// Base AC is 10
	baseAC := 10

	// Add dexterity modifier
	dexModifier := GetModifier(c.Attributes.Dexterity)

	return baseAC + dexModifier
}

// CalculateArmorAC calculates the armor class provided by equipped armor
func (c *Character) CalculateArmorAC() int {
	armorAC := 0

	// Add AC from equipped armor
	if c.Equipment.Armor != nil {
		armorAC += c.Equipment.Armor.Power
	}

	// Add AC from equipped accessory if it provides armor
	if c.Equipment.Accessory != nil && c.Equipment.Accessory.Type == ItemArmor {
		armorAC += c.Equipment.Accessory.Power
	}

	return armorAC
}

// CalculateTotalAC calculates the total armor class of the character
func (c *Character) CalculateTotalAC() int {
	// Start with base AC
	totalAC := c.CalculateBaseAC()

	// Add armor AC
	totalAC += c.CalculateArmorAC()

	// Apply class-specific bonuses
	switch c.Class {
	case Monk:
		// Monks get additional AC from wisdom when unarmored
		if c.Equipment.Armor == nil {
			wisModifier := GetModifier(c.Attributes.Wisdom)
			if wisModifier > 0 {
				totalAC += wisModifier
			}
		}
	case Barbarian:
		// Barbarians get additional AC from constitution when unarmored
		if c.Equipment.Armor == nil {
			conModifier := GetModifier(c.Attributes.Constitution)
			if conModifier > 0 {
				totalAC += conModifier
			}
		}
	}

	return totalAC
}

// CalculateHitChance calculates the chance to hit a target with the given AC
func (c *Character) CalculateHitChance(targetAC int) float64 {
	// Base hit chance is 50%
	baseHitChance := 0.5

	// Calculate attack bonus based on strength or dexterity (whichever is higher)
	strModifier := GetModifier(c.Attributes.Strength)
	dexModifier := GetModifier(c.Attributes.Dexterity)
	attackBonus := strModifier
	if dexModifier > strModifier {
		attackBonus = dexModifier
	}

	// Add level-based bonus
	attackBonus += c.Level / 2

	// Add weapon bonus if equipped
	if c.Equipment.Weapon != nil {
		// For simplicity, we'll say 20% of weapon power contributes to hit chance
		attackBonus += c.Equipment.Weapon.Power / 5
	}

	// Calculate hit chance: base + (attack bonus - (targetAC - 10)) * 0.05
	// This means each point of difference changes hit chance by 5%
	// We subtract 10 from targetAC because 10 is the base AC
	hitChance := baseHitChance + float64(attackBonus-(targetAC-10))*0.05

	// Clamp hit chance between 0.05 (5%) and 0.95 (95%)
	if hitChance < 0.05 {
		hitChance = 0.05 // Always at least 5% chance to hit
	} else if hitChance > 0.95 {
		hitChance = 0.95 // Always at least 5% chance to miss
	}

	return hitChance
}

// PerformSkillCheck performs a skill check for the character
func (c *Character) PerformSkillCheck(skillType SkillType, difficultyClass int) bool {
	if c.Skills == nil {
		return false
	}
	return c.Skills.PerformSkillCheck(skillType, c.Attributes, difficultyClass)
}

// AddSkillExperience adds experience to a skill and returns true if the skill leveled up
func (c *Character) AddSkillExperience(skillType SkillType, exp int) bool {
	if c.Skills == nil {
		return false
	}
	return c.Skills.AddSkillExperience(skillType, exp)
}

// GetSkillLevel returns the level of a specific skill
func (c *Character) GetSkillLevel(skillType SkillType) int {
	if c.Skills == nil {
		return 0
	}
	return c.Skills.GetSkillLevel(skillType)
}

// GetSkillBonus returns the bonus for a specific skill
func (c *Character) GetSkillBonus(skillType SkillType) int {
	if c.Skills == nil {
		return 0
	}
	return c.Skills.GetSkillBonus(skillType)
}
