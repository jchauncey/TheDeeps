package models

import (
	"fmt"
	"math/rand"
)

// MobDifficulty represents the difficulty level of a mob
type MobDifficulty string

const (
	// DifficultyEasy represents an easy mob
	DifficultyEasy MobDifficulty = "easy"
	// DifficultyNormal represents a normal mob
	DifficultyNormal MobDifficulty = "normal"
	// DifficultyHard represents a hard mob
	DifficultyHard MobDifficulty = "hard"
	// DifficultyElite represents an elite mob
	DifficultyElite MobDifficulty = "elite"
	// DifficultyBoss represents a boss mob
	DifficultyBoss MobDifficulty = "boss"
)

// MobType represents the type of mob
type MobType string

const (
	// Mob types
	MobSkeleton  MobType = "skeleton"
	MobGoblin    MobType = "goblin"
	MobTroll     MobType = "troll"
	MobOrc       MobType = "orc"
	MobOgre      MobType = "ogre"
	MobWraith    MobType = "wraith"
	MobLich      MobType = "lich"
	MobOoze      MobType = "ooze"
	MobRatman    MobType = "ratman"
	MobDrake     MobType = "drake"
	MobDragon    MobType = "dragon"
	MobElemental MobType = "elemental"
)

// MobDefinition defines the base stats and properties of a mob type
type MobDefinition struct {
	Type        MobType `json:"type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MinLevel    int     `json:"minLevel"` // Minimum dungeon level where this mob can appear
	BaseHealth  int     `json:"baseHealth"`
	BaseDamage  int     `json:"baseDamage"`
	BaseDefense int     `json:"baseDefense"`
	BaseSpeed   int     `json:"baseSpeed"`
	GoldRange   [2]int  `json:"goldRange"` // [min, max] gold drop range for normal difficulty
	// Special abilities or traits could be added here
}

// MobInstance represents an instance of a mob with specific stats based on difficulty
type MobInstance struct {
	ID         string        `json:"id"`
	Type       MobType       `json:"type"`
	Name       string        `json:"name"`
	Difficulty MobDifficulty `json:"difficulty"`
	Health     int           `json:"health"`
	MaxHealth  int           `json:"maxHealth"`
	Damage     int           `json:"damage"`
	Defense    int           `json:"defense"`
	Speed      int           `json:"speed"`
	GoldDrop   int           `json:"goldDrop"`
	Position   Position      `json:"position"`
	Status     []string      `json:"status,omitempty"`
}

// MobDefinitions is a map of mob types to their definitions
var MobDefinitions = map[MobType]MobDefinition{
	MobSkeleton: {
		Type:        MobSkeleton,
		Name:        "Skeleton",
		Description: "A reanimated skeleton wielding rusty weapons",
		MinLevel:    1,
		BaseHealth:  20,
		BaseDamage:  5,
		BaseDefense: 2,
		BaseSpeed:   3,
		GoldRange:   [2]int{2, 8},
	},
	MobGoblin: {
		Type:        MobGoblin,
		Name:        "Goblin",
		Description: "A small, cunning creature with a love for shiny things",
		MinLevel:    1,
		BaseHealth:  15,
		BaseDamage:  4,
		BaseDefense: 1,
		BaseSpeed:   5,
		GoldRange:   [2]int{5, 12},
	},
	MobTroll: {
		Type:        MobTroll,
		Name:        "Troll",
		Description: "A large, regenerating brute with immense strength",
		MinLevel:    3,
		BaseHealth:  50,
		BaseDamage:  10,
		BaseDefense: 5,
		BaseSpeed:   2,
		GoldRange:   [2]int{10, 25},
	},
	MobOrc: {
		Type:        MobOrc,
		Name:        "Orc",
		Description: "A fierce warrior with a thirst for battle",
		MinLevel:    2,
		BaseHealth:  30,
		BaseDamage:  8,
		BaseDefense: 3,
		BaseSpeed:   3,
		GoldRange:   [2]int{8, 15},
	},
	MobOgre: {
		Type:        MobOgre,
		Name:        "Ogre",
		Description: "A massive, dim-witted creature with devastating strength",
		MinLevel:    4,
		BaseHealth:  70,
		BaseDamage:  15,
		BaseDefense: 7,
		BaseSpeed:   1,
		GoldRange:   [2]int{15, 30},
	},
	MobWraith: {
		Type:        MobWraith,
		Name:        "Wraith",
		Description: "A spectral entity that drains life force",
		MinLevel:    5,
		BaseHealth:  35,
		BaseDamage:  12,
		BaseDefense: 4,
		BaseSpeed:   4,
		GoldRange:   [2]int{12, 20},
	},
	MobLich: {
		Type:        MobLich,
		Name:        "Lich",
		Description: "An undead sorcerer with powerful magic",
		MinLevel:    7,
		BaseHealth:  60,
		BaseDamage:  18,
		BaseDefense: 6,
		BaseSpeed:   3,
		GoldRange:   [2]int{20, 40},
	},
	MobOoze: {
		Type:        MobOoze,
		Name:        "Ooze",
		Description: "A corrosive slime that dissolves everything it touches",
		MinLevel:    2,
		BaseHealth:  25,
		BaseDamage:  7,
		BaseDefense: 8,
		BaseSpeed:   1,
		GoldRange:   [2]int{5, 15},
	},
	MobRatman: {
		Type:        MobRatman,
		Name:        "Ratman",
		Description: "A humanoid rat with disease-carrying weapons",
		MinLevel:    1,
		BaseHealth:  18,
		BaseDamage:  6,
		BaseDefense: 2,
		BaseSpeed:   6,
		GoldRange:   [2]int{4, 10},
	},
	MobDrake: {
		Type:        MobDrake,
		Name:        "Drake",
		Description: "A smaller cousin of dragons with elemental breath",
		MinLevel:    6,
		BaseHealth:  65,
		BaseDamage:  14,
		BaseDefense: 9,
		BaseSpeed:   5,
		GoldRange:   [2]int{18, 35},
	},
	MobDragon: {
		Type:        MobDragon,
		Name:        "Dragon",
		Description: "An ancient, powerful dragon with devastating attacks",
		MinLevel:    8,
		BaseHealth:  100,
		BaseDamage:  25,
		BaseDefense: 15,
		BaseSpeed:   4,
		GoldRange:   [2]int{50, 100},
	},
	MobElemental: {
		Type:        MobElemental,
		Name:        "Elemental",
		Description: "A being of pure elemental energy",
		MinLevel:    5,
		BaseHealth:  45,
		BaseDamage:  16,
		BaseDefense: 10,
		BaseSpeed:   3,
		GoldRange:   [2]int{15, 25},
	},
}

// GetMobsForFloorLevel returns a list of mob types that can appear on the given floor level
func GetMobsForFloorLevel(level int) []MobType {
	var availableMobs []MobType

	for mobType, definition := range MobDefinitions {
		if definition.MinLevel <= level {
			availableMobs = append(availableMobs, mobType)
		}
	}

	return availableMobs
}

// GetDifficultyMultiplier returns stat multipliers based on mob difficulty
func GetDifficultyMultiplier(difficulty MobDifficulty) (float64, float64) {
	// Returns (statMultiplier, goldMultiplier)
	switch difficulty {
	case DifficultyEasy:
		return 0.75, 0.5
	case DifficultyNormal:
		return 1.0, 1.0
	case DifficultyHard:
		return 1.5, 2.0
	case DifficultyElite:
		return 2.5, 3.0
	case DifficultyBoss:
		return 5.0, 5.0
	default:
		return 1.0, 1.0
	}
}

// GetRandomDifficulty returns a random difficulty based on floor level
func GetRandomDifficulty(floorLevel int) MobDifficulty {
	// Higher floor levels have higher chances of harder mobs
	roll := rand.Intn(100)

	// Base chances adjusted by floor level
	easyChance := Max(70-floorLevel*5, 10)
	normalChance := Min(20+floorLevel*2, 40)
	hardChance := Min(8+floorLevel*2, 30)
	eliteChance := Min(2+floorLevel, 15)
	// Boss chance is whatever remains (up to 5%)

	if roll < easyChance {
		return DifficultyEasy
	} else if roll < easyChance+normalChance {
		return DifficultyNormal
	} else if roll < easyChance+normalChance+hardChance {
		return DifficultyHard
	} else if roll < easyChance+normalChance+hardChance+eliteChance {
		return DifficultyElite
	} else {
		return DifficultyBoss
	}
}

// CreateMobInstance creates a new mob instance with stats based on type, difficulty, and floor level
func CreateMobInstance(mobType MobType, difficulty MobDifficulty, floorLevel int, position Position) *MobInstance {
	definition, exists := MobDefinitions[mobType]
	if !exists {
		// Default to goblin if type doesn't exist
		definition = MobDefinitions[MobGoblin]
	}

	// Get multipliers for the difficulty
	statMult, goldMult := GetDifficultyMultiplier(difficulty)

	// Scale stats based on floor level and difficulty
	levelScale := 1.0 + float64(floorLevel-1)*0.2

	// Calculate final stats
	health := int(float64(definition.BaseHealth) * statMult * levelScale)
	damage := int(float64(definition.BaseDamage) * statMult * levelScale)
	defense := int(float64(definition.BaseDefense) * statMult * levelScale)
	speed := int(float64(definition.BaseSpeed) * statMult * levelScale)

	// Calculate gold drop
	minGold := definition.GoldRange[0]
	maxGold := definition.GoldRange[1]
	goldRange := maxGold - minGold
	goldDrop := minGold + rand.Intn(goldRange+1)
	goldDrop = int(float64(goldDrop) * goldMult * levelScale)

	// Create name with difficulty prefix
	name := fmt.Sprintf("%s %s", string(difficulty), definition.Name)
	if difficulty == DifficultyNormal {
		name = definition.Name // No prefix for normal difficulty
	}

	return &MobInstance{
		ID:         generateID(),
		Type:       mobType,
		Name:       name,
		Difficulty: difficulty,
		Health:     health,
		MaxHealth:  health,
		Damage:     damage,
		Defense:    defense,
		Speed:      speed,
		GoldDrop:   goldDrop,
		Position:   position,
		Status:     []string{},
	}
}

// Helper function to get max of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper function to get min of two integers
func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
