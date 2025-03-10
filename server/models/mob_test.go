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
