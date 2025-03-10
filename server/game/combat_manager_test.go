package game

import (
	"testing"

	"github.com/jchauncey/TheDeeps/server/models"
	"github.com/stretchr/testify/assert"
)

func TestAttackMob(t *testing.T) {
	// Create a combat manager
	combatManager := NewCombatManager()

	// Create a character
	character := models.NewCharacter("TestWarrior", models.Warrior)
	character.Level = 1
	character.Attributes.Strength = 16 // High strength for consistent damage

	// Create a mob
	mob := models.NewMob(models.MobSkeleton, models.VariantNormal, 1)
	mob.HP = 10
	mob.MaxHP = 10
	mob.Damage = 2

	// Test attack
	result := combatManager.AttackMob(character, mob)

	// Check result
	assert.True(t, result.Success, "Attack should succeed")
	assert.NotEmpty(t, result.Message, "Message should not be empty")
	assert.Greater(t, result.DamageDealt, 0, "Damage dealt should be positive")

	// Check mob HP
	assert.Less(t, mob.HP, mob.MaxHP, "Mob HP should be reduced")

	// Test killing a mob
	mob.HP = 1 // Set HP low enough to be killed in one hit
	result = combatManager.AttackMob(character, mob)

	// Check result
	assert.True(t, result.Success, "Attack should succeed")
	assert.True(t, result.Killed, "Mob should be killed")
	assert.NotEmpty(t, result.Message, "Message should not be empty")
	assert.Greater(t, result.ExpGained, 0, "Experience gained should be positive")
	assert.Greater(t, result.GoldGained, 0, "Gold gained should be positive")

	// Check mob HP
	assert.Equal(t, 0, mob.HP, "Mob HP should be 0")
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
