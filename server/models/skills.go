package models

import (
	"math/rand"
)

// randomRoll is a variable that holds the function used for random rolls
// This allows for easier testing by replacing this function in tests
var randomRoll = func(n int) int {
	return rand.Intn(n)
}

// SkillType represents the type of skill
type SkillType string

const (
	// Combat skills
	SkillMelee  SkillType = "melee"
	SkillRanged SkillType = "ranged"
	SkillDodge  SkillType = "dodge"
	SkillBlock  SkillType = "block"

	// Exploration skills
	SkillStealth    SkillType = "stealth"
	SkillPerception SkillType = "perception"
	SkillSurvival   SkillType = "survival"
	SkillTraps      SkillType = "traps"

	// Interaction skills
	SkillLockpicking  SkillType = "lockpicking"
	SkillPersuasion   SkillType = "persuasion"
	SkillIntimidation SkillType = "intimidation"
	SkillDeception    SkillType = "deception"

	// Magic skills
	SkillArcana     SkillType = "arcana"
	SkillDivination SkillType = "divination"
	SkillElemental  SkillType = "elemental"
	SkillNecromancy SkillType = "necromancy"
)

// Skill represents a character skill with its level and experience
type Skill struct {
	Type        SkillType `json:"type"`
	Level       int       `json:"level"`
	Experience  int       `json:"experience"`
	Description string    `json:"description"`
}

// Skills represents all skills a character has
type Skills struct {
	SkillList map[SkillType]*Skill `json:"skillList"`
}

// SkillAttribute maps skills to their primary and secondary attributes
var SkillAttribute = map[SkillType]struct {
	Primary   string
	Secondary string
}{
	SkillMelee:        {"Strength", "Dexterity"},
	SkillRanged:       {"Dexterity", "Wisdom"},
	SkillDodge:        {"Dexterity", "Constitution"},
	SkillBlock:        {"Strength", "Constitution"},
	SkillStealth:      {"Dexterity", "Wisdom"},
	SkillPerception:   {"Wisdom", "Intelligence"},
	SkillSurvival:     {"Wisdom", "Constitution"},
	SkillTraps:        {"Intelligence", "Dexterity"},
	SkillLockpicking:  {"Dexterity", "Intelligence"},
	SkillPersuasion:   {"Charisma", "Wisdom"},
	SkillIntimidation: {"Charisma", "Strength"},
	SkillDeception:    {"Charisma", "Intelligence"},
	SkillArcana:       {"Intelligence", "Wisdom"},
	SkillDivination:   {"Wisdom", "Intelligence"},
	SkillElemental:    {"Intelligence", "Constitution"},
	SkillNecromancy:   {"Intelligence", "Constitution"},
}

// ClassSkillBonuses defines which skills get bonuses for each character class
var ClassSkillBonuses = map[CharacterClass][]SkillType{
	Warrior:   {SkillMelee, SkillBlock},
	Mage:      {SkillArcana, SkillElemental},
	Rogue:     {SkillStealth, SkillLockpicking},
	Cleric:    {SkillDivination, SkillPersuasion},
	Druid:     {SkillSurvival, SkillElemental},
	Warlock:   {SkillArcana, SkillNecromancy},
	Bard:      {SkillPersuasion, SkillDeception},
	Paladin:   {SkillMelee, SkillPersuasion},
	Ranger:    {SkillRanged, SkillSurvival},
	Monk:      {SkillDodge, SkillPerception},
	Barbarian: {SkillMelee, SkillIntimidation},
	Sorcerer:  {SkillElemental, SkillArcana},
}

// SkillDescriptions provides descriptions for each skill
var SkillDescriptions = map[SkillType]string{
	SkillMelee:        "Proficiency with melee weapons and close combat",
	SkillRanged:       "Accuracy with ranged weapons and thrown objects",
	SkillDodge:        "Ability to avoid incoming attacks",
	SkillBlock:        "Skill at blocking or parrying attacks",
	SkillStealth:      "Ability to move quietly and remain undetected",
	SkillPerception:   "Awareness of surroundings and ability to notice details",
	SkillSurvival:     "Knowledge of wilderness survival and navigation",
	SkillTraps:        "Ability to detect and disarm traps",
	SkillLockpicking:  "Skill at picking locks and disabling security mechanisms",
	SkillPersuasion:   "Ability to convince others through logical argument",
	SkillIntimidation: "Ability to influence others through threats or fear",
	SkillDeception:    "Skill at lying and misleading others",
	SkillArcana:       "Knowledge of magical theory and artifacts",
	SkillDivination:   "Ability to perceive hidden information through magic",
	SkillElemental:    "Control over elemental forces (fire, water, etc.)",
	SkillNecromancy:   "Understanding of death magic and undead creatures",
}

// NewSkills creates a new Skills struct with default values
func NewSkills(class CharacterClass) *Skills {
	skills := &Skills{
		SkillList: make(map[SkillType]*Skill),
	}

	// Initialize all skills at level 1
	for skillType, desc := range SkillDescriptions {
		skills.SkillList[skillType] = &Skill{
			Type:        skillType,
			Level:       1,
			Experience:  0,
			Description: desc,
		}
	}

	// Apply class bonuses
	if bonuses, exists := ClassSkillBonuses[class]; exists {
		for _, skillType := range bonuses {
			if skill, ok := skills.SkillList[skillType]; ok {
				skill.Level += 2 // Class skills start at level 3
			}
		}
	}

	return skills
}

// GetSkillLevel returns the level of a specific skill
func (s *Skills) GetSkillLevel(skillType SkillType) int {
	if skill, exists := s.SkillList[skillType]; exists {
		return skill.Level
	}
	return 0
}

// GetSkillBonus returns the bonus for a specific skill
func (s *Skills) GetSkillBonus(skillType SkillType) int {
	level := s.GetSkillLevel(skillType)
	// Skill bonus is (level - 1) / 2, similar to D&D proficiency
	return (level - 1) / 2
}

// CalculateExperienceForNextSkillLevel calculates the experience needed for the next skill level
func CalculateExperienceForNextSkillLevel(level int) int {
	return level * 100 // Simple progression: level * 100 XP needed
}

// AddSkillExperience adds experience to a skill and levels it up if necessary
// Returns true if the skill leveled up
func (s *Skills) AddSkillExperience(skillType SkillType, exp int) bool {
	skill, exists := s.SkillList[skillType]
	if !exists {
		return false
	}

	skill.Experience += exp
	leveledUp := false

	// Check if we have enough XP to level up
	nextLevelExp := CalculateExperienceForNextSkillLevel(skill.Level)
	if skill.Experience >= nextLevelExp && skill.Level < 20 {
		skill.Level++
		leveledUp = true
	}

	return leveledUp
}

// GetAttributeValue returns the value of an attribute by name
func GetAttributeValue(attrs Attributes, attrName string) int {
	switch attrName {
	case "Strength":
		return attrs.Strength
	case "Dexterity":
		return attrs.Dexterity
	case "Constitution":
		return attrs.Constitution
	case "Intelligence":
		return attrs.Intelligence
	case "Wisdom":
		return attrs.Wisdom
	case "Charisma":
		return attrs.Charisma
	default:
		return 0
	}
}

// PerformSkillCheck performs a skill check against a difficulty class (DC)
// Returns true if the check succeeds, false otherwise
func (s *Skills) PerformSkillCheck(skillType SkillType, attrs Attributes, difficultyClass int) bool {
	_, exists := s.SkillList[skillType]
	if !exists {
		return false
	}

	// Get the attribute modifiers
	attrInfo, exists := SkillAttribute[skillType]
	if !exists {
		return false
	}

	primaryAttr := GetAttributeValue(attrs, attrInfo.Primary)
	secondaryAttr := GetAttributeValue(attrs, attrInfo.Secondary)

	// Calculate primary and secondary attribute modifiers
	primaryMod := GetModifier(primaryAttr)
	secondaryMod := GetModifier(secondaryAttr) / 2 // Secondary attribute has half effect

	// Calculate total bonus: skill bonus + primary modifier + half of secondary modifier
	totalBonus := s.GetSkillBonus(skillType) + primaryMod + secondaryMod

	// Roll a d20 + bonus and compare to DC
	roll := randomRoll(20) + 1 // 1-20
	result := roll + totalBonus

	// Add a small amount of XP for using the skill, more if it was challenging
	xpGain := 5
	if result >= difficultyClass {
		// More XP for succeeding at harder checks
		xpGain += (difficultyClass - 10) / 2
		if xpGain < 5 {
			xpGain = 5
		}
	}
	s.AddSkillExperience(skillType, xpGain)

	// Critical success on natural 20
	if roll == 20 {
		return true
	}

	// Critical failure on natural 1
	if roll == 1 {
		return false
	}

	return result >= difficultyClass
}

// GetSkillCheckDifficulty returns a descriptive string for a difficulty class
func GetSkillCheckDifficulty(dc int) string {
	switch {
	case dc <= 5:
		return "Very Easy"
	case dc <= 10:
		return "Easy"
	case dc <= 15:
		return "Medium"
	case dc <= 20:
		return "Hard"
	case dc <= 25:
		return "Very Hard"
	default:
		return "Nearly Impossible"
	}
}

// GetSkillsForClass returns the recommended skills for a character class
func GetSkillsForClass(class CharacterClass) []SkillType {
	if skills, exists := ClassSkillBonuses[class]; exists {
		return skills
	}
	return []SkillType{}
}
