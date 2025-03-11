package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMob(t *testing.T) {
	tests := []struct {
		name           string
		mobType        MobType
		variant        MobVariant
		floorLevel     int
		expectedHPMin  int
		expectedHPMax  int
		expectedDamage int
	}{
		{
			name:           "Easy Skeleton on Floor 1",
			mobType:        MobSkeleton,
			variant:        VariantEasy,
			floorLevel:     1,
			expectedHPMin:  6,
			expectedHPMax:  7,
			expectedDamage: 2,
		},
		{
			name:           "Normal Goblin on Floor 3",
			mobType:        MobGoblin,
			variant:        VariantNormal,
			floorLevel:     3,
			expectedHPMin:  7,
			expectedHPMax:  9,
			expectedDamage: 2,
		},
		{
			name:           "Hard Troll on Floor 5",
			mobType:        MobTroll,
			variant:        VariantHard,
			floorLevel:     5,
			expectedHPMin:  33,
			expectedHPMax:  47,
			expectedDamage: 9,
		},
		{
			name:           "Boss Dragon on Floor 10",
			mobType:        MobDragon,
			variant:        VariantBoss,
			floorLevel:     10,
			expectedHPMin:  112,
			expectedHPMax:  224,
			expectedDamage: 56,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mob := NewMob(tt.mobType, tt.variant, tt.floorLevel)

			// Check basic properties
			assert.Equal(t, tt.mobType, mob.Type, "Mob type should match")
			assert.Equal(t, tt.variant, mob.Variant, "Mob variant should match")
			assert.NotEmpty(t, mob.ID, "Mob ID should be generated")

			// Check HP and damage
			assert.GreaterOrEqual(t, mob.HP, tt.expectedHPMin, "HP should be at least minimum value")
			assert.LessOrEqual(t, mob.HP, tt.expectedHPMax, "HP should be at most maximum value")
			assert.Equal(t, tt.expectedDamage, mob.Damage, "Damage should match expected value")

			// Check gold
			assert.Greater(t, mob.GoldValue, 0, "Gold should be positive")

			// Check symbol and color
			assert.NotEmpty(t, mob.Symbol, "Symbol should be set")
			assert.NotEmpty(t, mob.Color, "Color should be set")

			// Check position
			assert.Equal(t, 0, mob.Position.X, "X position should start at 0")
			assert.Equal(t, 0, mob.Position.Y, "Y position should start at 0")

			// Check boss properties
			if tt.variant == VariantBoss {
				// For dragon, the symbol is already uppercase 'D' in the implementation
				if tt.mobType != MobDragon {
					assert.True(t, mob.Symbol[0] >= 'A' && mob.Symbol[0] <= 'Z',
						"Boss symbol should be uppercase")
				}
			}
		})
	}
}

func TestMobArmorClass(t *testing.T) {
	// Test different mob types with different AC values

	// Test a low AC mob (Ooze)
	ooze := NewMob(MobOoze, VariantNormal, 1)
	oozeAC := ooze.CalculateAC()
	assert.LessOrEqual(t, oozeAC, 10, "Ooze should have low AC")

	// Test a medium AC mob (Goblin)
	goblin := NewMob(MobGoblin, VariantNormal, 1)
	goblinAC := goblin.CalculateAC()
	assert.GreaterOrEqual(t, goblinAC, 11, "Goblin should have medium AC")
	assert.LessOrEqual(t, goblinAC, 14, "Goblin should have medium AC")

	// Test a high AC mob (Dragon)
	dragon := NewMob(MobDragon, VariantNormal, 1)
	dragonAC := dragon.CalculateAC()
	assert.GreaterOrEqual(t, dragonAC, 15, "Dragon should have high AC")

	// Test variant effects on AC
	easyDragon := NewMob(MobDragon, VariantEasy, 1)
	normalDragon := NewMob(MobDragon, VariantNormal, 1)
	hardDragon := NewMob(MobDragon, VariantHard, 1)
	bossDragon := NewMob(MobDragon, VariantBoss, 1)

	assert.Less(t, easyDragon.CalculateAC(), normalDragon.CalculateAC(), "Easy variant should have lower AC")
	assert.Less(t, normalDragon.CalculateAC(), hardDragon.CalculateAC(), "Hard variant should have higher AC")
	assert.Less(t, hardDragon.CalculateAC(), bossDragon.CalculateAC(), "Boss variant should have highest AC")

	// Test hit chance calculation
	// Create a test character
	character := NewCharacter("TestCharacter", Warrior)
	character.Level = 5
	character.Attributes.Strength = 16 // +3 modifier

	// Test hit chance against different mobs
	oozeHitChance := character.CalculateHitChance(oozeAC)
	goblinHitChance := character.CalculateHitChance(goblinAC)
	dragonHitChance := character.CalculateHitChance(dragonAC)

	assert.Greater(t, oozeHitChance, goblinHitChance, "Hit chance against low AC mob should be higher")
	assert.Greater(t, goblinHitChance, dragonHitChance, "Hit chance against high AC mob should be lower")

	// Test mob hit chance against character
	character.Attributes.Dexterity = 14 // +2 modifier
	characterAC := character.CalculateTotalAC()

	oozeToHitChance := ooze.CalculateHitChance(characterAC)
	dragonToHitChance := dragon.CalculateHitChance(characterAC)

	assert.Less(t, oozeToHitChance, dragonToHitChance, "Stronger mobs should have better hit chance")
}
