package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCharacterSkillIntegration(t *testing.T) {
	// Create a character with the new skills system
	character := NewCharacterWithSkills("TestCharacter", Rogue)

	// Check that the character has skills
	assert.NotNil(t, character.Skills, "Character skills should not be nil")

	// Check that class skills have higher levels
	rogueSkills := []SkillType{SkillStealth, SkillLockpicking}
	for _, skillType := range rogueSkills {
		level := character.GetSkillLevel(skillType)
		assert.Equal(t, 3, level, "Rogue should have level 3 in %s, got %d", skillType, level)
	}

	// Check that non-class skills are at level 1
	nonRogueSkills := []SkillType{SkillMelee, SkillBlock, SkillArcana}
	for _, skillType := range nonRogueSkills {
		level := character.GetSkillLevel(skillType)
		assert.Equal(t, 1, level, "Rogue should have level 1 in %s, got %d", skillType, level)
	}

	// Test skill check with character attributes
	// Rogue has high dexterity, so stealth checks should be easier
	skillType := SkillStealth
	dc := 15 // Medium difficulty

	// Save the original random function
	origRoll := randomRoll

	// Override the random function for this test
	randomRoll = func(n int) int {
		return 10 // Roll an 11 (10+1)
	}

	// Make sure we restore the original function when done
	defer func() {
		randomRoll = origRoll
	}()

	// Perform skill check
	// Rogue with Dex 13 (+1), Stealth level 3 (+1), and roll of 11 = 13
	// Should fail against DC 15
	result := character.PerformSkillCheck(skillType, dc)
	assert.False(t, result, "Skill check should have failed with roll 11 + modifiers against DC %d", dc)

	// Add experience to level up the skill
	initialLevel := character.GetSkillLevel(skillType)
	xpNeeded := CalculateExperienceForNextSkillLevel(initialLevel)
	leveledUp := character.AddSkillExperience(skillType, xpNeeded)

	assert.True(t, leveledUp, "Skill should have leveled up")

	newLevel := character.GetSkillLevel(skillType)
	assert.Equal(t, initialLevel+1, newLevel, "Skill level should have increased by 1, got %d (was %d)", newLevel, initialLevel)

	// Try the skill check again with the higher skill level
	// Now with Stealth level 4 (+1), should still fail but be closer
	result = character.PerformSkillCheck(skillType, dc)
	assert.False(t, result, "Skill check should still fail even with higher skill level")

	// Increase the character's dexterity
	character.Attributes.Dexterity = 18 // +4 modifier

	// Try the skill check again with higher dexterity
	// Roll 11 + Stealth bonus 1 + Dex mod 4 = 16, should succeed against DC 15
	result = character.PerformSkillCheck(skillType, dc)
	assert.True(t, result, "Skill check should succeed with higher dexterity")
}
