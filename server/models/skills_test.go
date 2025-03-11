package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSkills(t *testing.T) {
	// Test skill initialization for different classes
	testCases := []struct {
		class          CharacterClass
		expectedSkills []SkillType
	}{
		{Warrior, []SkillType{SkillMelee, SkillBlock}},
		{Mage, []SkillType{SkillArcana, SkillElemental}},
		{Rogue, []SkillType{SkillStealth, SkillLockpicking}},
	}

	for _, tc := range testCases {
		skills := NewSkills(tc.class)

		// Check that all skills are initialized
		for skillType := range SkillDescriptions {
			skill, exists := skills.SkillList[skillType]
			assert.True(t, exists, "Skill %s not initialized for class %s", skillType, tc.class)
			if !exists {
				continue
			}

			// Check that the skill has the correct description
			assert.Equal(t, SkillDescriptions[skillType], skill.Description,
				"Skill %s has incorrect description for class %s", skillType, tc.class)

			// Check that class skills have higher levels
			isClassSkill := false
			for _, classSkill := range tc.expectedSkills {
				if skillType == classSkill {
					isClassSkill = true
					break
				}
			}

			if isClassSkill {
				assert.Equal(t, 3, skill.Level,
					"Class skill %s should be level 3 for class %s, got %d", skillType, tc.class, skill.Level)
			} else {
				assert.Equal(t, 1, skill.Level,
					"Non-class skill %s should be level 1 for class %s, got %d", skillType, tc.class, skill.Level)
			}
		}
	}
}

func TestSkillChecks(t *testing.T) {
	// Create a character with known attributes
	attrs := Attributes{
		Strength:     14, // +2 modifier
		Dexterity:    16, // +3 modifier
		Constitution: 12, // +1 modifier
		Intelligence: 10, // +0 modifier
		Wisdom:       8,  // -1 modifier
		Charisma:     10, // +0 modifier
	}

	// Create skills with known levels
	skills := NewSkills(Warrior)

	// Set specific skills to higher levels
	skills.SkillList[SkillMelee].Level = 5      // +2 bonus
	skills.SkillList[SkillStealth].Level = 3    // +1 bonus
	skills.SkillList[SkillPerception].Level = 1 // +0 bonus

	// Test cases with fixed dice rolls for deterministic testing
	testCases := []struct {
		name          string
		skillType     SkillType
		dc            int
		fixedRoll     int
		expectedXP    int
		expectSuccess bool
	}{
		{"Melee success", SkillMelee, 15, 10, 5, true},        // 10 (roll) + 2 (skill) + 2 (STR) + 1 (DEX/2) = 15 >= 15
		{"Stealth failure", SkillStealth, 15, 8, 5, false},    // 8 (roll) + 1 (skill) + 3 (DEX) + 0 (WIS/2) = 12 < 15
		{"Perception crit", SkillPerception, 20, 20, 5, true}, // Natural 20 always succeeds
		{"Arcana crit fail", SkillArcana, 5, 1, 5, false},     // Natural 1 always fails
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save the original random function
			origRoll := randomRoll

			// Create a fixed roll for testing
			fixedRoll := tc.fixedRoll - 1 // -1 because randomRoll returns 0-19, but we want 1-20

			// Override the random function for this test
			randomRoll = func(n int) int {
				return fixedRoll
			}

			// Make sure we restore the original function when done
			defer func() {
				randomRoll = origRoll
			}()

			// Perform the skill check
			result := skills.PerformSkillCheck(tc.skillType, attrs, tc.dc)

			assert.Equal(t, tc.expectSuccess, result, "Expected success=%v, got %v", tc.expectSuccess, result)

			// Check that XP was added
			skill := skills.SkillList[tc.skillType]
			assert.GreaterOrEqual(t, skill.Experience, tc.expectedXP,
				"Expected at least %d XP, got %d", tc.expectedXP, skill.Experience)
		})
	}
}

func TestSkillLeveling(t *testing.T) {
	// Create skills
	skills := NewSkills(Warrior)

	// Test adding experience and leveling up
	skillType := SkillStealth
	initialLevel := skills.GetSkillLevel(skillType)

	// Verify initial level
	assert.Equal(t, 1, initialLevel, "Expected initial level to be 1, got %d", initialLevel)

	// Add enough XP to level up exactly once
	xpNeeded := CalculateExperienceForNextSkillLevel(initialLevel)
	leveledUp := skills.AddSkillExperience(skillType, xpNeeded)

	assert.True(t, leveledUp, "Expected skill to level up, but it didn't")

	newLevel := skills.GetSkillLevel(skillType)
	expectedLevel := initialLevel + 1
	assert.Equal(t, expectedLevel, newLevel, "Expected level to be %d, got %d", expectedLevel, newLevel)

	// Test adding a small amount of XP (not enough to level up)
	leveledUp = skills.AddSkillExperience(skillType, 10)

	assert.False(t, leveledUp, "Skill shouldn't have leveled up with only 10 XP")

	// Verify level hasn't changed
	assert.Equal(t, newLevel, skills.GetSkillLevel(skillType), "Level shouldn't have changed after adding small XP")

	// Test adding enough XP for one more level up
	xpForNextLevel := CalculateExperienceForNextSkillLevel(newLevel)

	// Add the exact amount needed
	leveledUp = skills.AddSkillExperience(skillType, xpForNextLevel)

	assert.True(t, leveledUp, "Expected skill to level up, but it didn't")

	finalLevel := skills.GetSkillLevel(skillType)
	// The skill should now be at level 3 (started at 1, then +1, then +1 again)
	expectedFinalLevel := 3
	assert.Equal(t, expectedFinalLevel, finalLevel, "Expected final level to be %d, got %d", expectedFinalLevel, finalLevel)
}

func TestGetSkillCheckDifficulty(t *testing.T) {
	testCases := []struct {
		name            string
		difficultyClass int
		expectedLabel   string
	}{
		{"Very Easy", 3, "Very Easy"},
		{"Very Easy Boundary", 5, "Very Easy"},
		{"Easy", 8, "Easy"},
		{"Easy Boundary", 10, "Easy"},
		{"Medium", 12, "Medium"},
		{"Medium Boundary", 15, "Medium"},
		{"Hard", 18, "Hard"},
		{"Hard Boundary", 20, "Hard"},
		{"Very Hard", 22, "Very Hard"},
		{"Very Hard Boundary", 25, "Very Hard"},
		{"Nearly Impossible", 30, "Nearly Impossible"},
		{"Extreme", 50, "Nearly Impossible"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetSkillCheckDifficulty(tc.difficultyClass)
			assert.Equal(t, tc.expectedLabel, result,
				"Expected difficulty class %d to be labeled as '%s', got '%s'",
				tc.difficultyClass, tc.expectedLabel, result)
		})
	}
}

func TestGetSkillsForClass(t *testing.T) {
	testCases := []struct {
		name           string
		class          CharacterClass
		expectedSkills []SkillType
	}{
		{"Warrior", Warrior, []SkillType{SkillMelee, SkillBlock}},
		{"Mage", Mage, []SkillType{SkillArcana, SkillElemental}},
		{"Rogue", Rogue, []SkillType{SkillStealth, SkillLockpicking}},
		{"Cleric", Cleric, []SkillType{SkillDivination, SkillPersuasion}},
		{"Druid", Druid, []SkillType{SkillSurvival, SkillElemental}},
		{"Warlock", Warlock, []SkillType{SkillArcana, SkillNecromancy}},
		{"Bard", Bard, []SkillType{SkillPersuasion, SkillDeception}},
		{"Paladin", Paladin, []SkillType{SkillMelee, SkillPersuasion}},
		{"Ranger", Ranger, []SkillType{SkillRanged, SkillSurvival}},
		{"Monk", Monk, []SkillType{SkillDodge, SkillPerception}},
		{"Barbarian", Barbarian, []SkillType{SkillMelee, SkillIntimidation}},
		{"Sorcerer", Sorcerer, []SkillType{SkillElemental, SkillArcana}},
		{"Invalid Class", CharacterClass("invalid"), []SkillType{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetSkillsForClass(tc.class)

			// Check that the result has the same length as expected
			assert.Equal(t, len(tc.expectedSkills), len(result),
				"Expected %d skills for class %s, got %d",
				len(tc.expectedSkills), tc.class, len(result))

			// Check that all expected skills are in the result
			for _, expectedSkill := range tc.expectedSkills {
				found := false
				for _, resultSkill := range result {
					if expectedSkill == resultSkill {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected skill %s not found for class %s", expectedSkill, tc.class)
			}
		})
	}
}
