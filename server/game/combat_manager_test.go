package game

import (
	"testing"

	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/stretchr/testify/assert"
)

func TestAttackMob(t *testing.T) {
	// Create a combat manager
	combatManager := NewCombatManager()

	// Create a character with very high attributes to ensure hit and damage
	character := models.NewCharacter("TestWarrior", models.Warrior)
	character.Level = 1
	character.Attributes.Strength = 18  // Maximum strength for consistent damage
	character.Attributes.Dexterity = 18 // Maximum dexterity for consistent hits

	// Create a mob
	mob := models.NewMob(models.MobSkeleton, models.VariantNormal, 1)
	mob.HP = 10
	mob.MaxHP = 10
	mob.Damage = 2

	// Test attack - run multiple times to account for randomness
	var successfulAttack bool
	var damageDealt bool
	var hpReduced bool

	// Try up to 5 times to get a successful attack
	for i := 0; i < 5; i++ {
		// Reset mob HP for each attempt
		mob.HP = mob.MaxHP

		result := combatManager.AttackMob(character, mob)

		if result.Success {
			successfulAttack = true
			if result.DamageDealt > 0 {
				damageDealt = true
			}
			if mob.HP < mob.MaxHP {
				hpReduced = true
			}

			// If all conditions are met, break the loop
			if successfulAttack && damageDealt && hpReduced {
				break
			}
		}
	}

	// Check results
	assert.True(t, successfulAttack, "At least one attack should succeed")
	assert.True(t, damageDealt, "At least one attack should deal damage")
	assert.True(t, hpReduced, "At least one attack should reduce mob HP")

	// Test killing a mob
	mob.HP = 1 // Set HP low enough to be killed in one hit

	// Try up to 5 times to get a killing blow
	var mobKilled bool
	for i := 0; i < 5; i++ {
		// Reset mob HP for each attempt
		mob.HP = 1

		result := combatManager.AttackMob(character, mob)

		if result.Success && result.Killed {
			mobKilled = true
			break
		}
	}

	assert.True(t, mobKilled, "Mob should be killed with 1 HP remaining")
}

func TestUseItem(t *testing.T) {
	// Create a combat manager
	combatManager := NewCombatManager()

	// Create a character
	character := models.NewCharacter("TestWarrior", models.Warrior)
	character.CurrentHP = 5
	character.MaxHP = 10

	// Create a healing potion
	potion := models.NewPotion("Health Potion", 5, 10)

	// Test using the potion
	result := combatManager.UseItem(character, *potion)

	// Check result
	assert.True(t, result.Success, "Using item should succeed")
	assert.NotEmpty(t, result.Message, "Message should not be empty")

	// Check character HP
	assert.Equal(t, 10, character.CurrentHP, "Character HP should be fully restored")

	// Test using a potion at full health
	result = combatManager.UseItem(character, *potion)

	// Check result
	assert.False(t, result.Success, "Using item should fail when at full health")
	assert.NotEmpty(t, result.Message, "Message should not be empty")
}

func TestFlee(t *testing.T) {
	// Create a combat manager
	combatManager := NewCombatManager()

	// Create a character
	character := models.NewCharacter("TestWarrior", models.Warrior)
	character.Attributes.Dexterity = 18 // High dexterity for consistent flee success

	// Create a mob
	mob := models.NewMob(models.MobSkeleton, models.VariantNormal, 1)

	// Test flee attempt (may succeed or fail based on RNG)
	result := combatManager.Flee(character, mob)

	// Check result
	assert.NotEmpty(t, result.Message, "Message should not be empty")

	// If flee failed, check damage taken
	if !result.Success {
		assert.Greater(t, result.DamageTaken, 0, "Damage taken should be positive if flee failed")
	}
}
