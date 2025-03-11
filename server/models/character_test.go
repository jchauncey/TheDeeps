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

func TestArmorClass(t *testing.T) {
	// Create a character with known attributes
	character := NewCharacter("TestCharacter", Warrior)

	// Set attributes to predictable values
	character.Attributes.Dexterity = 14 // +2 modifier

	// Test base AC calculation (10 + dex modifier)
	baseAC := character.CalculateBaseAC()
	assert.Equal(t, 12, baseAC, "Base AC should be 10 + dex modifier")

	// Test total AC with no equipment
	totalAC := character.CalculateTotalAC()
	assert.Equal(t, 12, totalAC, "Total AC with no equipment should equal base AC")

	// Add armor to inventory and equip it
	armor := NewArmor("Test Armor", 5, 100, 1, nil)
	character.AddToInventory(armor)
	character.EquipItem(armor.ID)

	// Test armor AC calculation
	armorAC := character.CalculateArmorAC()
	assert.Equal(t, 5, armorAC, "Armor AC should equal armor power")

	// Test total AC with armor equipped
	totalAC = character.CalculateTotalAC()
	assert.Equal(t, 17, totalAC, "Total AC should be base AC + armor AC")

	// Test hit chance calculation
	// Against low AC target (AC 10)
	lowACHitChance := character.CalculateHitChance(10)
	assert.InDelta(t, 0.6, lowACHitChance, 0.01, "Hit chance against low AC should be high")

	// Against high AC target (AC 20)
	highACHitChance := character.CalculateHitChance(20)
	assert.InDelta(t, 0.1, highACHitChance, 0.01, "Hit chance against high AC should be low")

	// Test monk AC bonus (wisdom modifier when unarmored)
	monkCharacter := NewCharacter("TestMonk", Monk)
	monkCharacter.Attributes.Dexterity = 14 // +2 modifier
	monkCharacter.Attributes.Wisdom = 16    // +3 modifier

	// Unequip any armor
	monkCharacter.Equipment.Armor = nil

	// Test monk's total AC (should include wisdom bonus)
	monkAC := monkCharacter.CalculateTotalAC()
	assert.Equal(t, 15, monkAC, "Monk AC should be 10 + dex modifier + wisdom modifier")

	// Equip armor on monk
	monkCharacter.AddToInventory(armor)
	monkCharacter.EquipItem(armor.ID)

	// Test monk's total AC with armor (should not include wisdom bonus)
	monkArmoredAC := monkCharacter.CalculateTotalAC()
	assert.Equal(t, 17, monkArmoredAC, "Monk AC with armor should be base AC + armor AC")
}

// TestRemoveFromInventory tests the RemoveFromInventory function
func TestRemoveFromInventory(t *testing.T) {
	tests := []struct {
		name           string
		inventory      []*Item
		itemIDToRemove string
		expectedFound  bool
		expectedItem   *Item
	}{
		{
			name: "Remove Existing Item",
			inventory: []*Item{
				{ID: "item1", Name: "Sword", Type: ItemWeapon},
				{ID: "item2", Name: "Shield", Type: ItemArmor},
			},
			itemIDToRemove: "item1",
			expectedFound:  true,
			expectedItem: &Item{
				ID:   "item1",
				Name: "Sword",
				Type: ItemWeapon,
			},
		},
		{
			name: "Remove Non-Existent Item",
			inventory: []*Item{
				{ID: "item1", Name: "Sword", Type: ItemWeapon},
			},
			itemIDToRemove: "item2",
			expectedFound:  false,
			expectedItem:   nil,
		},
		{
			name:           "Remove From Empty Inventory",
			inventory:      []*Item{},
			itemIDToRemove: "item1",
			expectedFound:  false,
			expectedItem:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a character with the test inventory
			character := NewCharacter("TestChar", Warrior)
			character.Inventory = tt.inventory

			// Call RemoveFromInventory
			item, found := character.RemoveFromInventory(tt.itemIDToRemove)

			// Check results
			assert.Equal(t, tt.expectedFound, found, "Found status should match expected")
			if tt.expectedFound {
				assert.Equal(t, tt.expectedItem.ID, item.ID, "Removed item ID should match expected")
				assert.Equal(t, tt.expectedItem.Name, item.Name, "Removed item name should match expected")
				assert.Equal(t, tt.expectedItem.Type, item.Type, "Removed item type should match expected")

				// Verify item is no longer in inventory
				for _, invItem := range character.Inventory {
					assert.NotEqual(t, tt.itemIDToRemove, invItem.ID, "Item should be removed from inventory")
				}
			} else {
				assert.Nil(t, item, "Item should be nil when not found")
			}
		})
	}
}

// TestUnequipItem tests the UnequipItem function
func TestUnequipItem(t *testing.T) {
	tests := []struct {
		name              string
		equippedItems     map[ItemType]*Item
		itemTypeToUnequip ItemType
		expectedSuccess   bool
	}{
		{
			name: "Unequip Weapon",
			equippedItems: map[ItemType]*Item{
				ItemWeapon: {ID: "weapon1", Name: "Sword", Type: ItemWeapon},
			},
			itemTypeToUnequip: ItemWeapon,
			expectedSuccess:   true,
		},
		{
			name: "Unequip Armor",
			equippedItems: map[ItemType]*Item{
				ItemArmor: {ID: "armor1", Name: "Shield", Type: ItemArmor},
			},
			itemTypeToUnequip: ItemArmor,
			expectedSuccess:   true,
		},
		{
			name: "Unequip Accessory",
			equippedItems: map[ItemType]*Item{
				ItemArtifact: {ID: "acc1", Name: "Ring", Type: ItemArtifact},
			},
			itemTypeToUnequip: ItemArtifact,
			expectedSuccess:   true,
		},
		{
			name: "Unequip Non-Equipped Item Type",
			equippedItems: map[ItemType]*Item{
				ItemWeapon: {ID: "weapon1", Name: "Sword", Type: ItemWeapon},
			},
			itemTypeToUnequip: ItemArmor,
			expectedSuccess:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a character
			character := NewCharacter("TestChar", Warrior)

			// Set up equipment
			if weapon, exists := tt.equippedItems[ItemWeapon]; exists {
				character.Equipment.Weapon = weapon
			}
			if armor, exists := tt.equippedItems[ItemArmor]; exists {
				character.Equipment.Armor = armor
			}
			if accessory, exists := tt.equippedItems[ItemArtifact]; exists {
				character.Equipment.Accessory = accessory
			}

			// Call UnequipItem
			success := character.UnequipItem(tt.itemTypeToUnequip)

			// Check results
			assert.Equal(t, tt.expectedSuccess, success, "Success status should match expected")

			if tt.expectedSuccess {
				// Check that equipment slot is now empty
				switch tt.itemTypeToUnequip {
				case ItemWeapon:
					assert.Nil(t, character.Equipment.Weapon, "Weapon slot should be empty")
				case ItemArmor:
					assert.Nil(t, character.Equipment.Armor, "Armor slot should be empty")
				case ItemArtifact:
					assert.Nil(t, character.Equipment.Accessory, "Accessory slot should be empty")
				}
			}
		})
	}
}

// TestUseItem tests the UseItem function
func TestUseItem(t *testing.T) {
	tests := []struct {
		name            string
		inventory       []*Item
		itemIDToUse     string
		initialHP       int
		initialMana     int
		expectedSuccess bool
		expectedHP      int
		expectedMana    int
	}{
		{
			name: "Use Health Potion",
			inventory: []*Item{
				{ID: "potion1", Name: "Health Potion", Type: ItemPotion, Power: 10},
			},
			itemIDToUse:     "potion1",
			initialHP:       10,
			initialMana:     20,
			expectedSuccess: true,
			expectedHP:      20,
			expectedMana:    20,
		},
		{
			name: "Use Mana Potion (Scroll)",
			inventory: []*Item{
				{ID: "potion2", Name: "Mana Potion", Type: ItemScroll, Power: 10},
			},
			itemIDToUse:     "potion2",
			initialHP:       20,
			initialMana:     10,
			expectedSuccess: true,
			expectedHP:      20,
			expectedMana:    20,
		},
		{
			name: "Use Non-Existent Item",
			inventory: []*Item{
				{ID: "potion1", Name: "Health Potion", Type: ItemPotion, Power: 10},
			},
			itemIDToUse:     "potion2",
			initialHP:       20,
			initialMana:     20,
			expectedSuccess: false,
			expectedHP:      20,
			expectedMana:    20,
		},
		{
			name: "Use Non-Consumable Item",
			inventory: []*Item{
				{ID: "sword1", Name: "Sword", Type: ItemWeapon},
			},
			itemIDToUse:     "sword1",
			initialHP:       20,
			initialMana:     20,
			expectedSuccess: false,
			expectedHP:      20,
			expectedMana:    20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a character with the test inventory
			character := NewCharacter("TestChar", Warrior)
			character.Inventory = tt.inventory
			character.MaxHP = 30
			character.CurrentHP = tt.initialHP
			character.MaxMana = 30
			character.CurrentMana = tt.initialMana

			// Initial inventory count
			initialInventoryCount := len(character.Inventory)

			// Call UseItem
			success := character.UseItem(tt.itemIDToUse)

			// Check results
			assert.Equal(t, tt.expectedSuccess, success, "Success status should match expected")
			assert.Equal(t, tt.expectedHP, character.CurrentHP, "Current HP should match expected")
			assert.Equal(t, tt.expectedMana, character.CurrentMana, "Current Mana should match expected")

			if tt.expectedSuccess {
				// Check that item was removed from inventory
				assert.Equal(t, initialInventoryCount-1, len(character.Inventory), "Inventory should have one less item")

				// Verify item is no longer in inventory
				for _, invItem := range character.Inventory {
					assert.NotEqual(t, tt.itemIDToUse, invItem.ID, "Item should be removed from inventory")
				}
			} else {
				// Check that inventory didn't change for non-consumable or non-existent items
				assert.Equal(t, initialInventoryCount, len(character.Inventory), "Inventory should remain unchanged")
			}
		})
	}
}

// TestCalculateAttackPower tests the CalculateAttackPower function
func TestCalculateAttackPower(t *testing.T) {
	tests := []struct {
		name          string
		strength      int
		weaponPower   int
		expectedPower int
	}{
		{
			name:          "No Weapon",
			strength:      14,
			weaponPower:   0,
			expectedPower: 3, // Strength modifier of +2 + level 1
		},
		{
			name:          "With Weapon",
			strength:      16,
			weaponPower:   5,
			expectedPower: 9, // Strength modifier of +3 + weapon power of 5 + level 1
		},
		{
			name:          "Low Strength",
			strength:      8,
			weaponPower:   3,
			expectedPower: 3, // Strength modifier of -1 + weapon power of 3 + level 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a character
			character := NewCharacter("TestChar", Warrior)
			character.Attributes.Strength = tt.strength

			// Set up weapon if needed
			if tt.weaponPower > 0 {
				character.Equipment.Weapon = &Item{
					ID:    "weapon1",
					Name:  "Test Weapon",
					Type:  ItemWeapon,
					Power: tt.weaponPower,
				}
			}

			// Call CalculateAttackPower
			power := character.CalculateAttackPower()

			// Check result
			assert.Equal(t, tt.expectedPower, power, "Attack power should match expected")
		})
	}
}

// TestCalculateDefensePower tests the CalculateDefensePower function
func TestCalculateDefensePower(t *testing.T) {
	tests := []struct {
		name          string
		constitution  int
		armorPower    int
		expectedPower int
	}{
		{
			name:          "No Armor",
			constitution:  14,
			armorPower:    0,
			expectedPower: 2, // Constitution modifier of +2
		},
		{
			name:          "With Armor",
			constitution:  16,
			armorPower:    5,
			expectedPower: 8, // Constitution modifier of +3 + armor power of 5
		},
		{
			name:          "Low Constitution",
			constitution:  8,
			armorPower:    3,
			expectedPower: 2, // Constitution modifier of -1 + armor power of 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a character
			character := NewCharacter("TestChar", Warrior)
			character.Attributes.Constitution = tt.constitution

			// Set up armor if needed
			if tt.armorPower > 0 {
				character.Equipment.Armor = &Item{
					ID:    "armor1",
					Name:  "Test Armor",
					Type:  ItemArmor,
					Power: tt.armorPower,
				}
			}

			// Call CalculateDefensePower
			power := character.CalculateDefensePower()

			// Check result
			assert.Equal(t, tt.expectedPower, power, "Defense power should match expected")
		})
	}
}

// TestCalculateArmorClass tests the CalculateBaseAC, CalculateArmorAC, and CalculateTotalAC functions
func TestCalculateArmorClass(t *testing.T) {
	tests := []struct {
		name            string
		dexterity       int
		armorPower      int
		expectedBaseAC  int
		expectedArmorAC int
		expectedTotalAC int
	}{
		{
			name:            "No Armor, Average Dexterity",
			dexterity:       10,
			armorPower:      0,
			expectedBaseAC:  10,
			expectedArmorAC: 0,
			expectedTotalAC: 10,
		},
		{
			name:            "No Armor, High Dexterity",
			dexterity:       16,
			armorPower:      0,
			expectedBaseAC:  13,
			expectedArmorAC: 0,
			expectedTotalAC: 13,
		},
		{
			name:            "With Armor, Average Dexterity",
			dexterity:       10,
			armorPower:      5,
			expectedBaseAC:  10,
			expectedArmorAC: 5,
			expectedTotalAC: 15,
		},
		{
			name:            "With Armor, High Dexterity",
			dexterity:       18,
			armorPower:      7,
			expectedBaseAC:  14,
			expectedArmorAC: 7,
			expectedTotalAC: 21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a character
			character := NewCharacter("TestChar", Warrior)
			character.Attributes.Dexterity = tt.dexterity

			// Set up armor if needed
			if tt.armorPower > 0 {
				character.Equipment.Armor = &Item{
					ID:    "armor1",
					Name:  "Test Armor",
					Type:  ItemArmor,
					Power: tt.armorPower,
				}
			}

			// Call AC calculation functions
			baseAC := character.CalculateBaseAC()
			armorAC := character.CalculateArmorAC()
			totalAC := character.CalculateTotalAC()

			// Check results
			assert.Equal(t, tt.expectedBaseAC, baseAC, "Base AC should match expected")
			assert.Equal(t, tt.expectedArmorAC, armorAC, "Armor AC should match expected")
			assert.Equal(t, tt.expectedTotalAC, totalAC, "Total AC should match expected")
		})
	}
}

// TestCalculateHitChance tests the CalculateHitChance function
func TestCalculateHitChance(t *testing.T) {
	tests := []struct {
		name           string
		strength       int
		targetAC       int
		expectedChance float64
	}{
		{
			name:           "Easy Hit",
			strength:       16,
			targetAC:       10,
			expectedChance: 0.65, // 65% chance to hit
		},
		{
			name:           "Moderate Hit",
			strength:       12,
			targetAC:       15,
			expectedChance: 0.30, // 30% chance to hit
		},
		{
			name:           "Difficult Hit",
			strength:       10,
			targetAC:       20,
			expectedChance: 0.05, // 5% minimum chance to hit
		},
		{
			name:           "Very Difficult Hit",
			strength:       8,
			targetAC:       25,
			expectedChance: 0.05, // 5% minimum chance to hit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a character
			character := NewCharacter("TestChar", Warrior)
			character.Attributes.Strength = tt.strength

			// Call CalculateHitChance
			chance := character.CalculateHitChance(tt.targetAC)

			// Check result with a small delta for floating point comparison
			assert.InDelta(t, tt.expectedChance, chance, 0.01, "Hit chance should match expected")
		})
	}
}

// TestGetSkillBonus tests the GetSkillBonus function
func TestGetSkillBonus(t *testing.T) {
	tests := []struct {
		name          string
		skillType     SkillType
		skillLevel    int
		expectedBonus int
	}{
		{
			name:          "Melee Skill",
			skillType:     SkillMelee,
			skillLevel:    3,
			expectedBonus: 1,
		},
		{
			name:          "Stealth Skill",
			skillType:     SkillStealth,
			skillLevel:    5,
			expectedBonus: 2,
		},
		{
			name:          "Perception Skill",
			skillType:     SkillPerception,
			skillLevel:    1,
			expectedBonus: 0,
		},
		{
			name:          "Arcana Skill",
			skillType:     SkillArcana,
			skillLevel:    0,
			expectedBonus: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a character with skills
			character := NewCharacterWithSkills("TestChar", Warrior)

			// Set the skill level directly in the skill list
			if skill, exists := character.Skills.SkillList[tt.skillType]; exists {
				skill.Level = tt.skillLevel
			} else {
				// Create the skill if it doesn't exist
				character.Skills.SkillList[tt.skillType] = &Skill{
					Type:  tt.skillType,
					Level: tt.skillLevel,
				}
			}

			// Call GetSkillBonus
			bonus := character.GetSkillBonus(tt.skillType)

			// Check result
			assert.Equal(t, tt.expectedBonus, bonus, "Skill bonus should match expected")
		})
	}
}

// TestAddExperienceManaIncrease tests the mana increase logic for different character classes when leveling up
func TestAddExperienceManaIncrease(t *testing.T) {
	tests := []struct {
		name                 string
		class                CharacterClass
		attributeToIncrease  string
		attributeValue       int
		expectedManaIncrease int
	}{
		{
			name:                 "Mage Intelligence Bonus",
			class:                Mage,
			attributeToIncrease:  "Intelligence",
			attributeValue:       16, // +3 modifier
			expectedManaIncrease: 4,  // +3 from Intelligence modifier, +1 base
		},
		{
			name:                 "Sorcerer Intelligence Bonus",
			class:                Sorcerer,
			attributeToIncrease:  "Intelligence",
			attributeValue:       14, // +2 modifier
			expectedManaIncrease: 3,  // +2 from Intelligence modifier, +1 base
		},
		{
			name:                 "Warlock Intelligence Bonus",
			class:                Warlock,
			attributeToIncrease:  "Intelligence",
			attributeValue:       12, // +1 modifier
			expectedManaIncrease: 2,  // +1 from Intelligence modifier, +1 base
		},
		{
			name:                 "Cleric Wisdom Bonus",
			class:                Cleric,
			attributeToIncrease:  "Wisdom",
			attributeValue:       18, // +4 modifier
			expectedManaIncrease: 5,  // +4 from Wisdom modifier, +1 base
		},
		{
			name:                 "Druid Wisdom Bonus",
			class:                Druid,
			attributeToIncrease:  "Wisdom",
			attributeValue:       16, // +3 modifier
			expectedManaIncrease: 4,  // +3 from Wisdom modifier, +1 base
		},
		{
			name:                 "Bard Charisma Bonus",
			class:                Bard,
			attributeToIncrease:  "Charisma",
			attributeValue:       20, // +5 modifier
			expectedManaIncrease: 6,  // +5 from Charisma modifier, +1 base
		},
		{
			name:                 "Paladin Charisma Bonus",
			class:                Paladin,
			attributeToIncrease:  "Charisma",
			attributeValue:       14, // +2 modifier
			expectedManaIncrease: 3,  // +2 from Charisma modifier, +1 base
		},
		{
			name:                 "Ranger Default Bonus",
			class:                Ranger,
			attributeToIncrease:  "",
			attributeValue:       0,
			expectedManaIncrease: 1, // Just the base +1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a character of the specified class
			character := NewCharacter("TestCharacter", tt.class)

			// Set the relevant attribute if specified
			if tt.attributeToIncrease != "" {
				switch tt.attributeToIncrease {
				case "Intelligence":
					character.Attributes.Intelligence = tt.attributeValue
				case "Wisdom":
					character.Attributes.Wisdom = tt.attributeValue
				case "Charisma":
					character.Attributes.Charisma = tt.attributeValue
				}
			}

			// Skip test if character has no mana
			if character.MaxMana == 0 {
				t.Skip("Character class does not use mana")
				return
			}

			// Record initial mana
			initialMana := character.MaxMana

			// Add enough experience to level up
			expNeeded := CalculateExperienceForNextLevel(character.Level)
			character.AddExperience(expNeeded)

			// Check that mana increased by the expected amount
			manaIncrease := character.MaxMana - initialMana
			assert.Equal(t, tt.expectedManaIncrease, manaIncrease,
				"Mana should increase by the expected amount on level up")

			// Verify that current mana equals max mana after level up
			assert.Equal(t, character.MaxMana, character.CurrentMana,
				"Current mana should equal max mana after level up")
		})
	}
}
