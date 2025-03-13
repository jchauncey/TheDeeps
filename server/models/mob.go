package models

import (
	"github.com/google/uuid"
)

// MobType represents the type of mob
type MobType string

const (
	MobSkeleton   MobType = "skeleton"
	MobGoblin     MobType = "goblin"
	MobTroll      MobType = "troll"
	MobOrc        MobType = "orc"
	MobOgre       MobType = "ogre"
	MobWraith     MobType = "wraith"
	MobLich       MobType = "lich"
	MobOoze       MobType = "ooze"
	MobRatman     MobType = "ratman"
	MobDrake      MobType = "drake"
	MobDragon     MobType = "dragon"
	MobElemental  MobType = "elemental"
	MobShopkeeper MobType = "shopkeeper"
)

// MobVariant represents the difficulty variant of a mob
type MobVariant string

const (
	VariantEasy   MobVariant = "easy"
	VariantNormal MobVariant = "normal"
	VariantHard   MobVariant = "hard"
	VariantBoss   MobVariant = "boss"
)

// Mob represents a monster in the dungeon
type Mob struct {
	ID        string     `json:"id"`
	Type      MobType    `json:"type"`
	Variant   MobVariant `json:"variant"`
	Name      string     `json:"name"`
	Level     int        `json:"level"`
	HP        int        `json:"hp"`
	MaxHP     int        `json:"maxHp"`
	Damage    int        `json:"damage"`
	Defense   int        `json:"defense"`
	AC        int        `json:"ac"`        // Armor Class
	Dexterity int        `json:"dexterity"` // Dexterity attribute for AC calculation
	GoldValue int        `json:"goldValue"`
	Position  Position   `json:"position"`
	Symbol    string     `json:"symbol"`
	Color     string     `json:"color"`
}

// NewMob creates a new mob based on type, variant, and floor level
func NewMob(mobType MobType, variant MobVariant, floorLevel int) *Mob {
	// Base stats that will be modified by variant and floor level
	baseHP := 10
	baseDamage := 2
	baseDefense := 0
	baseAC := 10
	baseDexterity := 10
	baseGoldValue := 5
	symbol := "m"
	color := "#FF0000" // Default red

	// Adjust base stats based on mob type
	switch mobType {
	case MobSkeleton:
		baseHP = 8
		baseDamage = 3
		baseAC = 12
		baseDexterity = 8
		symbol = "s"
		color = "#FFFFFF" // White
	case MobGoblin:
		baseHP = 6
		baseDamage = 2
		baseAC = 11
		baseDexterity = 14
		symbol = "g"
		color = "#00FF00" // Green
	case MobTroll:
		baseHP = 20
		baseDamage = 4
		baseDefense = 2
		baseAC = 14
		baseDexterity = 8
		symbol = "T"
		color = "#008000" // Dark green
	case MobOrc:
		baseHP = 12
		baseDamage = 3
		baseDefense = 1
		baseAC = 13
		baseDexterity = 10
		symbol = "o"
		color = "#808000" // Olive
	case MobOgre:
		baseHP = 25
		baseDamage = 5
		baseDefense = 1
		baseAC = 15
		baseDexterity = 6
		symbol = "O"
		color = "#800080" // Purple
	case MobWraith:
		baseHP = 15
		baseDamage = 6
		baseAC = 13
		baseDexterity = 16
		symbol = "W"
		color = "#000080" // Navy
	case MobLich:
		baseHP = 30
		baseDamage = 8
		baseDefense = 3
		baseAC = 16
		baseDexterity = 12
		symbol = "L"
		color = "#800000" // Maroon
	case MobOoze:
		baseHP = 18
		baseDamage = 3
		baseDefense = 4
		baseAC = 8
		baseDexterity = 4
		symbol = "j"
		color = "#008080" // Teal
	case MobRatman:
		baseHP = 7
		baseDamage = 2
		baseAC = 12
		baseDexterity = 15
		symbol = "r"
		color = "#808080" // Gray
	case MobDrake:
		baseHP = 22
		baseDamage = 6
		baseDefense = 2
		baseAC = 15
		baseDexterity = 12
		symbol = "d"
		color = "#FF8000" // Orange
	case MobDragon:
		baseHP = 40
		baseDamage = 10
		baseDefense = 5
		baseAC = 18
		baseDexterity = 10
		symbol = "D"
		color = "#FF0000" // Red
	case MobElemental:
		baseHP = 20
		baseDamage = 7
		baseDefense = 2
		baseAC = 14
		baseDexterity = 14
		symbol = "E"
		color = "#0000FF" // Blue
	case MobShopkeeper:
		baseHP = 10
		baseDamage = 2
		baseDefense = 0
		baseAC = 10
		baseDexterity = 10
		baseGoldValue = 5
		symbol = "S"
		color = "#FF0000" // Red
	}

	// Adjust stats based on variant
	variantMultiplier := 1.0
	switch variant {
	case VariantEasy:
		variantMultiplier = 0.8
		baseGoldValue = int(float64(baseGoldValue) * 0.7)
	case VariantNormal:
		variantMultiplier = 1.0
	case VariantHard:
		variantMultiplier = 1.3
		baseGoldValue = int(float64(baseGoldValue) * 1.5)
	case VariantBoss:
		variantMultiplier = 2.0
		baseGoldValue = int(float64(baseGoldValue) * 3.0)
		symbol = string([]rune(symbol)[0] - 32) // Convert to uppercase
	}

	// Adjust stats based on floor level (deeper floors have stronger mobs)
	floorMultiplier := 1.0 + (float64(floorLevel-1) * 0.2)

	// Calculate final stats
	finalHP := int(float64(baseHP) * variantMultiplier * floorMultiplier)
	finalDamage := int(float64(baseDamage) * variantMultiplier * floorMultiplier)
	finalDefense := int(float64(baseDefense) * variantMultiplier * floorMultiplier)
	finalAC := int(float64(baseAC) * variantMultiplier)
	finalDexterity := baseDexterity // Dexterity doesn't scale with level
	finalGoldValue := int(float64(baseGoldValue) * variantMultiplier * floorMultiplier)

	// Generate a name based on type and variant
	name := string(mobType)
	if variant == VariantBoss {
		name = "Boss " + name
	}

	return &Mob{
		ID:        uuid.New().String(),
		Type:      mobType,
		Variant:   variant,
		Name:      name,
		Level:     floorLevel,
		HP:        finalHP,
		MaxHP:     finalHP,
		Damage:    finalDamage,
		Defense:   finalDefense,
		AC:        finalAC,
		Dexterity: finalDexterity,
		GoldValue: finalGoldValue,
		Position:  Position{X: 0, Y: 0}, // Will be set when placed on the map
		Symbol:    symbol,
		Color:     color,
	}
}

// CalculateAC calculates the total armor class of the mob
func (m *Mob) CalculateAC() int {
	// Base AC from the mob's natural armor
	totalAC := m.AC

	// Add dexterity modifier
	dexModifier := (m.Dexterity - 10) / 2
	if dexModifier > 0 {
		totalAC += dexModifier
	}

	return totalAC
}

// CalculateHitChance calculates the chance to hit a target with the given AC
func (m *Mob) CalculateHitChance(targetAC int) float64 {
	// Base hit chance is 50%
	baseHitChance := 0.5

	// Calculate attack bonus based on level and type
	attackBonus := m.Level / 2

	// Add damage as a factor (stronger mobs are more accurate)
	attackBonus += m.Damage / 3

	// Calculate hit chance: base + (attack bonus - (targetAC - 10)) * 0.05
	// This means each point of difference changes hit chance by 5%
	// We subtract 10 from targetAC because 10 is the base AC
	hitChance := baseHitChance + float64(attackBonus-(targetAC-10))*0.05

	// Clamp hit chance between 0.05 (5%) and 0.95 (95%)
	if hitChance < 0.05 {
		hitChance = 0.05 // Always at least 5% chance to hit
	} else if hitChance > 0.95 {
		hitChance = 0.95 // Always at least 5% chance to miss
	}

	return hitChance
}
