package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jchauncey/TheDeeps/server/models"
)

// CombatResult represents the result of a combat action
type CombatResult struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	DamageDealt  int           `json:"damageDealt,omitempty"`
	DamageTaken  int           `json:"damageTaken,omitempty"`
	CriticalHit  bool          `json:"criticalHit,omitempty"`
	Killed       bool          `json:"killed,omitempty"`
	ExpGained    int           `json:"expGained,omitempty"`
	GoldGained   int           `json:"goldGained,omitempty"`
	ItemsDropped []models.Item `json:"itemsDropped,omitempty"`
}

// CombatManager handles combat mechanics
type CombatManager struct {
	rng *rand.Rand
}

// NewCombatManager creates a new combat manager
func NewCombatManager() *CombatManager {
	return &CombatManager{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// AttackMob handles a character attacking a mob
func (cm *CombatManager) AttackMob(character *models.Character, mob *models.Mob) CombatResult {
	result := CombatResult{
		Success: true,
	}

	// Calculate hit chance (base 70% + dexterity modifier)
	hitChance := 70 + models.GetModifier(character.Attributes.Dexterity)*5
	hitRoll := cm.rng.Intn(100) + 1

	// Check if attack hits
	if hitRoll > hitChance {
		result.Success = false
		result.Message = "Attack missed!"
		return result
	}

	// Calculate damage
	baseDamage := 1 // Unarmed damage
	// TODO: Add weapon damage when equipment is implemented

	// Add strength modifier
	damage := baseDamage + models.GetModifier(character.Attributes.Strength)
	if damage < 1 {
		damage = 1 // Minimum damage is 1
	}

	// Check for critical hit (5% base chance + wisdom modifier)
	critChance := 5 + models.GetModifier(character.Attributes.Wisdom)
	critRoll := cm.rng.Intn(100) + 1
	if critRoll <= critChance {
		damage *= 2
		result.CriticalHit = true
		result.Message = "Critical hit!"
	} else {
		result.Message = "Hit!"
	}

	// Apply damage to mob
	result.DamageDealt = damage
	mob.HP -= damage

	// Check if mob is killed
	if mob.HP <= 0 {
		mob.HP = 0
		result.Killed = true
		result.Message = fmt.Sprintf("Killed %s!", mob.Name)

		// Calculate experience gained
		expGained := calculateExpGain(mob, character.Level)
		result.ExpGained = expGained
		character.AddExperience(expGained)

		// Calculate gold gained
		goldGained := mob.GoldValue
		result.GoldGained = goldGained
		character.Gold += goldGained

		// Generate item drops
		// TODO: Implement item drop mechanics
	} else {
		// Mob counterattack
		mobDamage := calculateMobDamage(mob, character)
		result.DamageTaken = mobDamage
		character.CurrentHP -= mobDamage
		if character.CurrentHP < 0 {
			character.CurrentHP = 0
		}
	}

	return result
}

// UseItem handles a character using an item during combat
func (cm *CombatManager) UseItem(character *models.Character, item models.Item) CombatResult {
	result := CombatResult{
		Success: true,
	}

	switch item.Type {
	case models.ItemPotion:
		// Handle potion use
		if character.CurrentHP < character.MaxHP {
			healAmount := item.Power
			character.CurrentHP += healAmount
			if character.CurrentHP > character.MaxHP {
				character.CurrentHP = character.MaxHP
			}
			result.Message = fmt.Sprintf("Used %s and healed %d HP!", item.Name, healAmount)
		} else {
			result.Success = false
			result.Message = "Already at full health!"
		}
	default:
		result.Success = false
		result.Message = "Cannot use this item in combat!"
	}

	return result
}

// Flee handles a character attempting to flee from combat
func (cm *CombatManager) Flee(character *models.Character, mob *models.Mob) CombatResult {
	result := CombatResult{}

	// Calculate flee chance (base 50% + dexterity modifier - mob level)
	fleeChance := 50 + models.GetModifier(character.Attributes.Dexterity)*5 - mob.Level
	if fleeChance < 10 {
		fleeChance = 10 // Minimum 10% chance to flee
	}
	if fleeChance > 90 {
		fleeChance = 90 // Maximum 90% chance to flee
	}

	fleeRoll := cm.rng.Intn(100) + 1
	if fleeRoll <= fleeChance {
		result.Success = true
		result.Message = "Successfully fled from combat!"
	} else {
		result.Success = false
		result.Message = "Failed to flee!"

		// Mob gets a free attack
		mobDamage := calculateMobDamage(mob, character)
		result.DamageTaken = mobDamage
		character.CurrentHP -= mobDamage
		if character.CurrentHP < 0 {
			character.CurrentHP = 0
		}
	}

	return result
}

// Helper functions

// calculateExpGain calculates experience gained from defeating a mob
func calculateExpGain(mob *models.Mob, characterLevel int) int {
	baseExp := 0

	// Base experience based on mob type
	switch mob.Type {
	case models.MobSkeleton:
		baseExp = 10
	case models.MobGoblin:
		baseExp = 15
	case models.MobTroll:
		baseExp = 25
	case models.MobOrc:
		baseExp = 20
	case models.MobOgre:
		baseExp = 30
	case models.MobWraith:
		baseExp = 35
	case models.MobLich:
		baseExp = 50
	case models.MobOoze:
		baseExp = 20
	case models.MobRatman:
		baseExp = 12
	case models.MobDrake:
		baseExp = 40
	case models.MobDragon:
		baseExp = 100
	case models.MobElemental:
		baseExp = 45
	default:
		baseExp = 10
	}

	// Adjust based on variant
	switch mob.Variant {
	case models.VariantEasy:
		baseExp = int(float64(baseExp) * 0.7)
	case models.VariantNormal:
		// No adjustment
	case models.VariantHard:
		baseExp = int(float64(baseExp) * 1.5)
	case models.VariantBoss:
		baseExp = int(float64(baseExp) * 3.0)
	}

	// Level difference adjustment
	levelDiff := mob.Level - characterLevel
	if levelDiff > 5 {
		// Bonus for defeating much stronger mobs
		baseExp = int(float64(baseExp) * 1.5)
	} else if levelDiff < -5 {
		// Penalty for defeating much weaker mobs
		baseExp = int(float64(baseExp) * 0.5)
	}

	return baseExp
}

// calculateMobDamage calculates damage dealt by a mob to a character
func calculateMobDamage(mob *models.Mob, character *models.Character) int {
	// Base damage from mob
	damage := mob.Damage

	// Apply character defense (from constitution modifier)
	defense := models.GetModifier(character.Attributes.Constitution)
	damage -= defense

	// Ensure minimum damage of 1
	if damage < 1 {
		damage = 1
	}

	return damage
}
