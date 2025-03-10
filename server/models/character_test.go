package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCharacter(t *testing.T) {
	tests := []struct {
		name           string
		characterName  string
		class          CharacterClass
		expectedLevel  int
		expectedHP     int
		expectedMana   int
		expectedStrMin int
		expectedStrMax int
	}{
		{
			name:           "Warrior Character",
			characterName:  "TestWarrior",
			class:          Warrior,
			expectedLevel:  1,
			expectedHP:     27,
			expectedMana:   0,
			expectedStrMin: 12,
			expectedStrMax: 12,
		},
		{
			name:           "Mage Character",
			characterName:  "TestMage",
			class:          Mage,
			expectedLevel:  1,
			expectedHP:     18,
			expectedMana:   23,
			expectedStrMin: 9,
			expectedStrMax: 9,
		},
		{
			name:           "Rogue Character",
			characterName:  "TestRogue",
			class:          Rogue,
			expectedLevel:  1,
			expectedHP:     19,
			expectedMana:   10,
			expectedStrMin: 10,
			expectedStrMax: 10,
		},
		{
			name:           "Cleric Character",
			characterName:  "TestCleric",
			class:          Cleric,
			expectedLevel:  1,
			expectedHP:     20,
			expectedMana:   23,
			expectedStrMin: 10,
			expectedStrMax: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			character := NewCharacter(tt.characterName, tt.class)

			// Check basic properties
			assert.Equal(t, tt.characterName, character.Name, "Character name should match")
			assert.Equal(t, tt.class, character.Class, "Character class should match")
			assert.Equal(t, tt.expectedLevel, character.Level, "Character level should be 1")
			assert.NotEmpty(t, character.ID, "Character ID should be generated")

			// Check HP and Mana
			assert.Equal(t, tt.expectedHP, character.MaxHP, "Max HP should match expected value")
			assert.Equal(t, character.MaxHP, character.CurrentHP, "Current HP should equal Max HP")
			assert.Equal(t, tt.expectedMana, character.MaxMana, "Max Mana should match expected value")
			assert.Equal(t, character.MaxMana, character.CurrentMana, "Current Mana should equal Max Mana")

			// Check attributes
			assert.GreaterOrEqual(t, character.Attributes.Strength, tt.expectedStrMin, "Strength should be at least minimum value")
			assert.LessOrEqual(t, character.Attributes.Strength, tt.expectedStrMax, "Strength should be at most maximum value")
			assert.Greater(t, character.Attributes.Dexterity, 0, "Dexterity should be positive")
			assert.Greater(t, character.Attributes.Constitution, 0, "Constitution should be positive")
			assert.Greater(t, character.Attributes.Intelligence, 0, "Intelligence should be positive")
			assert.Greater(t, character.Attributes.Wisdom, 0, "Wisdom should be positive")
			assert.Greater(t, character.Attributes.Charisma, 0, "Charisma should be positive")

			// Check other properties
			assert.Equal(t, 0, character.Experience, "Experience should start at 0")
			assert.Equal(t, 0, character.Gold, "Gold should start at 0")
		})
	}
}

func TestAddExperience(t *testing.T) {
	tests := []struct {
		name          string
		initialLevel  int
		initialExp    int
		expToAdd      int
		expectedLevel int
		expectedExp   int
		shouldLevelUp bool
	}{
		{
			name:          "No Level Up",
			initialLevel:  1,
			initialExp:    0,
			expToAdd:      50,
			expectedLevel: 1,
			expectedExp:   50,
			shouldLevelUp: false,
		},
		{
			name:          "Level Up Once",
			initialLevel:  1,
			initialExp:    500,
			expToAdd:      600,
			expectedLevel: 2,
			expectedExp:   1100, // Experience is not reset after level up
			shouldLevelUp: true,
		},
		{
			name:          "Level Up Twice",
			initialLevel:  1,
			initialExp:    0,
			expToAdd:      3500,
			expectedLevel: 4,    // Levels up from 1 to 4
			expectedExp:   3500, // Experience is not reset after level up
			shouldLevelUp: true,
		},
		{
			name:          "Level Up With Remainder",
			initialLevel:  2,
			initialExp:    0,
			expToAdd:      2500,
			expectedLevel: 3,
			expectedExp:   2500, // Experience is not reset after level up
			shouldLevelUp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			character := NewCharacter("TestCharacter", Warrior)
			character.Level = tt.initialLevel
			character.Experience = tt.initialExp

			// Save initial attributes for comparison
			initialHP := character.MaxHP
			initialMana := character.MaxMana
			initialStr := character.Attributes.Strength

			// Add experience
			leveledUp := character.AddExperience(tt.expToAdd)

			// Check level and experience
			assert.Equal(t, tt.expectedLevel, character.Level, "Character level should match expected")
			assert.Equal(t, tt.expectedExp, character.Experience, "Character experience should match expected")
			assert.Equal(t, tt.shouldLevelUp, leveledUp, "Level up flag should match expected")

			// If leveled up, check that attributes increased
			if tt.shouldLevelUp {
				assert.Greater(t, character.MaxHP, initialHP, "Max HP should increase on level up")
				if character.MaxMana > 0 {
					assert.GreaterOrEqual(t, character.MaxMana, initialMana, "Max Mana should not decrease on level up")
				}
				assert.GreaterOrEqual(t, character.Attributes.Strength, initialStr, "Strength should not decrease on level up")
			} else {
				assert.Equal(t, initialHP, character.MaxHP, "Max HP should not change without level up")
				assert.Equal(t, initialMana, character.MaxMana, "Max Mana should not change without level up")
				assert.Equal(t, initialStr, character.Attributes.Strength, "Strength should not change without level up")
			}
		})
	}
}

func TestGetModifier(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{
			name:     "Very Low Value",
			value:    1,
			expected: -4,
		},
		{
			name:     "Low Value",
			value:    5,
			expected: -2,
		},
		{
			name:     "Below Average",
			value:    9,
			expected: 0,
		},
		{
			name:     "Average Value",
			value:    10,
			expected: 0,
		},
		{
			name:     "Above Average",
			value:    12,
			expected: 1,
		},
		{
			name:     "High Value",
			value:    16,
			expected: 3,
		},
		{
			name:     "Very High Value",
			value:    20,
			expected: 5,
		},
		{
			name:     "Maximum Value",
			value:    30,
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := GetModifier(tt.value)
			assert.Equal(t, tt.expected, modifier, "Modifier should match expected value")
		})
	}
}

func TestEquipItem(t *testing.T) {
	character := NewCharacter("TestCharacter", Warrior)

	// Create test items
	sword1 := NewWeapon("First Sword", 10, 100, 1, nil)
	sword2 := NewWeapon("Second Sword", 15, 150, 1, nil)
	armor1 := NewArmor("First Armor", 5, 80, 1, nil)
	armor2 := NewArmor("Second Armor", 8, 120, 1, nil)

	// Add items to inventory
	character.AddToInventory(sword1)
	character.AddToInventory(sword2)
	character.AddToInventory(armor1)
	character.AddToInventory(armor2)

	// Test equipping first weapon
	success := character.EquipItem(sword1.ID)
	assert.True(t, success)
	assert.Equal(t, sword1, character.Equipment.Weapon)
	assert.True(t, sword1.Equipped)

	// Test equipping second weapon (should replace first)
	success = character.EquipItem(sword2.ID)
	assert.True(t, success)
	assert.Equal(t, sword2, character.Equipment.Weapon)
	assert.True(t, sword2.Equipped)
	assert.False(t, sword1.Equipped)

	// Test equipping first armor
	success = character.EquipItem(armor1.ID)
	assert.True(t, success)
	assert.Equal(t, armor1, character.Equipment.Armor)
	assert.True(t, armor1.Equipped)

	// Test equipping second armor (should replace first)
	success = character.EquipItem(armor2.ID)
	assert.True(t, success)
	assert.Equal(t, armor2, character.Equipment.Armor)
	assert.True(t, armor2.Equipped)
	assert.False(t, armor1.Equipped)

	// Verify both types of equipment are still equipped
	assert.Equal(t, sword2, character.Equipment.Weapon)
	assert.Equal(t, armor2, character.Equipment.Armor)
}

func TestInventoryWeight(t *testing.T) {
	// Create a character with a class that doesn't modify strength
	character := NewCharacter("TestCharacter", Rogue)

	// Create test items with different weights
	sword := NewWeaponWithWeight("Test Sword", 10, 100, 5.0, 1, nil)
	armor := NewArmorWithWeight("Test Armor", 5, 80, 15.0, 1, nil)
	potion := NewPotion("Health Potion", 20, 30) // Weight 0.5
	gold := NewGold(1000)                        // Weight 10.0

	// Test initial weight (should be 0)
	assert.Equal(t, 0.0, character.CalculateTotalWeight())
	assert.Equal(t, 0.0, character.CalculateInventoryWeight())
	assert.Equal(t, 0.0, character.CalculateEquipmentWeight())

	// Ensure strength is exactly 10 for consistent testing
	character.Attributes.Strength = 10

	// Test weight limit calculation
	// Character has 10 strength, so limit should be 50
	weightLimit := character.CalculateWeightLimit()
	assert.Equal(t, 50.0, weightLimit)

	// Increase strength and test weight limit
	character.Attributes.Strength = 15
	weightLimit = character.CalculateWeightLimit()
	assert.Equal(t, 100.0, weightLimit) // 50 + (5 * 10)

	// Decrease strength and test weight limit
	character.Attributes.Strength = 8
	weightLimit = character.CalculateWeightLimit()
	assert.Equal(t, 40.0, weightLimit) // 50 + (-2 * 0.5 * 10)

	// Reset strength to 10
	character.Attributes.Strength = 10

	// Add items to inventory and test weight
	success := character.AddToInventory(sword)
	assert.True(t, success)
	assert.Equal(t, 5.0, character.CalculateInventoryWeight())

	success = character.AddToInventory(armor)
	assert.True(t, success)
	assert.Equal(t, 20.0, character.CalculateInventoryWeight())

	success = character.AddToInventory(potion)
	assert.True(t, success)
	assert.Equal(t, 20.5, character.CalculateInventoryWeight())

	// Test equipping items and weight transfer
	character.EquipItem(sword.ID)
	assert.Equal(t, 15.5, character.CalculateInventoryWeight())
	assert.Equal(t, 5.0, character.CalculateEquipmentWeight())
	assert.Equal(t, 20.5, character.CalculateTotalWeight())

	character.EquipItem(armor.ID)
	assert.Equal(t, 0.5, character.CalculateInventoryWeight())
	assert.Equal(t, 20.0, character.CalculateEquipmentWeight())
	assert.Equal(t, 20.5, character.CalculateTotalWeight())

	// Test encumbrance levels
	assert.Equal(t, 0, character.GetEncumbranceLevel()) // Not encumbered (20.5 < 25)

	// Add gold to approach weight limit
	success = character.AddToInventory(gold)
	assert.True(t, success)
	assert.Equal(t, 10.5, character.CalculateInventoryWeight())
	assert.Equal(t, 30.5, character.CalculateTotalWeight())

	// Test encumbrance levels
	assert.Equal(t, 1, character.GetEncumbranceLevel()) // Lightly encumbered (30.5 < 37.5)

	// Create heavy items to exceed weight limit
	heavyItem1 := NewArmorWithWeight("Heavy Armor", 10, 200, 10.0, 1, nil)
	success = character.AddToInventory(heavyItem1)
	assert.True(t, success)
	assert.Equal(t, 40.5, character.CalculateTotalWeight())
	assert.False(t, character.IsOverEncumbered())       // 40.5 < 50
	assert.Equal(t, 2, character.GetEncumbranceLevel()) // Heavily encumbered

	// Try to add another item that would exceed the limit
	heavyItem2 := NewArmorWithWeight("Another Heavy Armor", 10, 200, 15.0, 1, nil)
	success = character.AddToInventory(heavyItem2)
	assert.False(t, success)                                // Should fail due to weight limit (40.5 + 15 > 50)
	assert.Equal(t, 40.5, character.CalculateTotalWeight()) // Weight should not change
}
