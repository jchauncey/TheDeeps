package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWeapon(t *testing.T) {
	tests := []struct {
		name         string
		weaponName   string
		damage       int
		value        int
		levelReq     int
		classReq     []CharacterClass
		expectedType ItemType
	}{
		{
			name:         "Basic Sword",
			weaponName:   "Sword",
			damage:       5,
			value:        10,
			levelReq:     1,
			classReq:     nil,
			expectedType: ItemWeapon,
		},
		{
			name:         "Warrior Axe",
			weaponName:   "Battle Axe",
			damage:       10,
			value:        50,
			levelReq:     5,
			classReq:     []CharacterClass{Warrior, Barbarian},
			expectedType: ItemWeapon,
		},
		{
			name:         "Magic Staff",
			weaponName:   "Staff of Fire",
			damage:       15,
			value:        100,
			levelReq:     10,
			classReq:     []CharacterClass{Mage, Sorcerer, Warlock},
			expectedType: ItemWeapon,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewWeapon(tt.weaponName, tt.damage, tt.value, tt.levelReq, tt.classReq)

			// Check basic properties
			assert.Equal(t, tt.weaponName, item.Name, "Weapon name should match")
			assert.Equal(t, tt.expectedType, item.Type, "Item type should be weapon")
			assert.Equal(t, tt.damage, item.Power, "Power (damage) should match")
			assert.Equal(t, tt.value, item.Value, "Value should match")
			assert.Equal(t, tt.levelReq, item.LevelReq, "Level requirement should match")

			// Check class requirements
			assert.Equal(t, len(tt.classReq), len(item.ClassReq), "Class requirement count should match")
			for i, class := range tt.classReq {
				assert.Equal(t, class, item.ClassReq[i], "Class requirement should match")
			}

			// Check that ID is generated
			assert.NotEmpty(t, item.ID, "Item ID should be generated")

			// Check that symbol and color are set
			assert.NotEmpty(t, item.Symbol, "Symbol should be set")
			assert.NotEmpty(t, item.Color, "Color should be set")

			// Check that position is initialized to origin
			assert.Equal(t, 0, item.Position.X, "X position should start at 0")
			assert.Equal(t, 0, item.Position.Y, "Y position should start at 0")

			// Check that equipped is initialized to false
			assert.False(t, item.Equipped, "Equipped should start as false")
		})
	}
}

func TestNewArmor(t *testing.T) {
	tests := []struct {
		name         string
		armorName    string
		defense      int
		value        int
		levelReq     int
		classReq     []CharacterClass
		expectedType ItemType
	}{
		{
			name:         "Leather Armor",
			armorName:    "Leather Armor",
			defense:      3,
			value:        15,
			levelReq:     1,
			classReq:     nil,
			expectedType: ItemArmor,
		},
		{
			name:         "Chain Mail",
			armorName:    "Chain Mail",
			defense:      8,
			value:        75,
			levelReq:     5,
			classReq:     []CharacterClass{Warrior, Paladin, Cleric},
			expectedType: ItemArmor,
		},
		{
			name:         "Mage Robes",
			armorName:    "Arcane Robes",
			defense:      5,
			value:        120,
			levelReq:     8,
			classReq:     []CharacterClass{Mage, Sorcerer, Warlock},
			expectedType: ItemArmor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewArmor(tt.armorName, tt.defense, tt.value, tt.levelReq, tt.classReq)

			// Check basic properties
			assert.Equal(t, tt.armorName, item.Name, "Armor name should match")
			assert.Equal(t, tt.expectedType, item.Type, "Item type should be armor")
			assert.Equal(t, tt.defense, item.Power, "Power (defense) should match")
			assert.Equal(t, tt.value, item.Value, "Value should match")
			assert.Equal(t, tt.levelReq, item.LevelReq, "Level requirement should match")

			// Check class requirements
			assert.Equal(t, len(tt.classReq), len(item.ClassReq), "Class requirement count should match")
			for i, class := range tt.classReq {
				assert.Equal(t, class, item.ClassReq[i], "Class requirement should match")
			}

			// Check that ID is generated
			assert.NotEmpty(t, item.ID, "Item ID should be generated")

			// Check that symbol and color are set
			assert.NotEmpty(t, item.Symbol, "Symbol should be set")
			assert.NotEmpty(t, item.Color, "Color should be set")
		})
	}
}

func TestNewPotion(t *testing.T) {
	tests := []struct {
		name         string
		potionName   string
		power        int
		value        int
		expectedType ItemType
	}{
		{
			name:         "Health Potion",
			potionName:   "Health Potion",
			power:        20,
			value:        25,
			expectedType: ItemPotion,
		},
		{
			name:         "Mana Potion",
			potionName:   "Mana Potion",
			power:        30,
			value:        35,
			expectedType: ItemPotion,
		},
		{
			name:         "Strength Potion",
			potionName:   "Strength Potion",
			power:        5,
			value:        50,
			expectedType: ItemPotion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewPotion(tt.potionName, tt.power, tt.value)

			// Check basic properties
			assert.Equal(t, tt.potionName, item.Name, "Potion name should match")
			assert.Equal(t, tt.expectedType, item.Type, "Item type should be potion")
			assert.Equal(t, tt.power, item.Power, "Power should match")
			assert.Equal(t, tt.value, item.Value, "Value should match")

			// Check that ID is generated
			assert.NotEmpty(t, item.ID, "Item ID should be generated")

			// Check that symbol and color are set
			assert.NotEmpty(t, item.Symbol, "Symbol should be set")
			assert.NotEmpty(t, item.Color, "Color should be set")

			// Check that class and level requirements are not set
			assert.Nil(t, item.ClassReq, "Class requirements should be nil")
			assert.Zero(t, item.LevelReq, "Level requirement should be 0")
		})
	}
}

func TestNewGold(t *testing.T) {
	tests := []struct {
		name         string
		amount       int
		expectedType ItemType
	}{
		{
			name:         "Small Gold Pile",
			amount:       10,
			expectedType: ItemGold,
		},
		{
			name:         "Medium Gold Pile",
			amount:       50,
			expectedType: ItemGold,
		},
		{
			name:         "Large Gold Pile",
			amount:       100,
			expectedType: ItemGold,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewGold(tt.amount)

			// Check basic properties
			assert.Equal(t, "Gold", item.Name, "Name should be Gold")
			assert.Equal(t, tt.expectedType, item.Type, "Item type should be gold")
			assert.Equal(t, tt.amount, item.Value, "Value should match amount")

			// Check that ID is generated
			assert.NotEmpty(t, item.ID, "Item ID should be generated")

			// Check that symbol and color are set
			assert.NotEmpty(t, item.Symbol, "Symbol should be set")
			assert.NotEmpty(t, item.Color, "Color should be set")

			// Check that position is initialized to origin
			assert.Equal(t, 0, item.Position.X, "X position should start at 0")
			assert.Equal(t, 0, item.Position.Y, "Y position should start at 0")
		})
	}
}

func TestGenerateRandomItem(t *testing.T) {
	tests := []struct {
		name       string
		floorLevel int
	}{
		{
			name:       "Floor 1 Item",
			floorLevel: 1,
		},
		{
			name:       "Floor 5 Item",
			floorLevel: 5,
		},
		{
			name:       "Floor 10 Item",
			floorLevel: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := GenerateRandomItem(tt.floorLevel)

			// Check that item is generated
			assert.NotNil(t, item, "Item should be generated")

			// Check that ID is generated
			assert.NotEmpty(t, item.ID, "Item ID should be generated")

			// Check that power scales with floor level
			assert.GreaterOrEqual(t, item.Power, tt.floorLevel, "Power should scale with floor level")

			// Check that value scales with floor level
			assert.GreaterOrEqual(t, item.Value, tt.floorLevel*10, "Value should scale with floor level")
		})
	}
}
