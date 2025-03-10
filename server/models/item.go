package models

import (
	"math/rand"
)

// ItemRarity represents the rarity of an item
type ItemRarity string

const (
	// RarityCommon represents a common item
	RarityCommon ItemRarity = "common"
	// RarityUncommon represents an uncommon item
	RarityUncommon ItemRarity = "uncommon"
	// RarityRare represents a rare item
	RarityRare ItemRarity = "rare"
	// RarityEpic represents an epic item
	RarityEpic ItemRarity = "epic"
	// RarityLegendary represents a legendary item
	RarityLegendary ItemRarity = "legendary"
)

// ItemType represents the type of item
type ItemType string

const (
	// ItemTypeWeapon represents a weapon
	ItemTypeWeapon ItemType = "weapon"
	// ItemTypeArmor represents armor
	ItemTypeArmor ItemType = "armor"
	// ItemTypePotion represents a potion
	ItemTypePotion ItemType = "potion"
	// ItemTypeScroll represents a scroll
	ItemTypeScroll ItemType = "scroll"
	// ItemTypeGold represents gold
	ItemTypeGold ItemType = "gold"
	// ItemTypeGem represents a gem
	ItemTypeGem ItemType = "gem"
	// ItemTypeKey represents a key
	ItemTypeKey ItemType = "key"
)

// ItemDefinition defines the base properties of an item type
type ItemDefinition struct {
	Type        ItemType   `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Rarity      ItemRarity `json:"rarity"`
	Value       int        `json:"value"`
	// Add more item properties as needed
}

// ItemInstance represents an instance of an item in the game
type ItemInstance struct {
	ID          string     `json:"id"`
	Type        ItemType   `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Rarity      ItemRarity `json:"rarity"`
	Value       int        `json:"value"`
	Position    Position   `json:"position"`
	// Add more item properties as needed
}

// GetRarityChance returns the chance of getting an item of the specified rarity
// based on the mob difficulty and floor level
func GetRarityChance(difficulty MobDifficulty, floorLevel int) map[ItemRarity]float64 {
	// Base chances
	chances := map[ItemRarity]float64{
		RarityCommon:    0.70,
		RarityUncommon:  0.20,
		RarityRare:      0.08,
		RarityEpic:      0.015,
		RarityLegendary: 0.005,
	}

	// Adjust based on difficulty
	switch difficulty {
	case DifficultyEasy:
		// No adjustment for easy mobs
	case DifficultyNormal:
		chances[RarityCommon] -= 0.05
		chances[RarityUncommon] += 0.04
		chances[RarityRare] += 0.01
	case DifficultyHard:
		chances[RarityCommon] -= 0.15
		chances[RarityUncommon] += 0.05
		chances[RarityRare] += 0.08
		chances[RarityEpic] += 0.015
		chances[RarityLegendary] += 0.005
	case DifficultyElite:
		chances[RarityCommon] -= 0.30
		chances[RarityUncommon] += 0.10
		chances[RarityRare] += 0.15
		chances[RarityEpic] += 0.04
		chances[RarityLegendary] += 0.01
	case DifficultyBoss:
		chances[RarityCommon] -= 0.50
		chances[RarityUncommon] += 0.15
		chances[RarityRare] += 0.25
		chances[RarityEpic] += 0.08
		chances[RarityLegendary] += 0.02
	}

	// Adjust based on floor level (deeper floors have better loot)
	floorMultiplier := float64(floorLevel) * 0.01
	chances[RarityCommon] -= floorMultiplier * 0.5
	chances[RarityUncommon] += floorMultiplier * 0.2
	chances[RarityRare] += floorMultiplier * 0.2
	chances[RarityEpic] += floorMultiplier * 0.07
	chances[RarityLegendary] += floorMultiplier * 0.03

	// Ensure chances are within valid range
	for rarity, chance := range chances {
		if chance < 0 {
			chances[rarity] = 0
		}
	}

	// Normalize chances to ensure they sum to 1.0
	total := 0.0
	for _, chance := range chances {
		total += chance
	}
	for rarity, chance := range chances {
		chances[rarity] = chance / total
	}

	return chances
}

// GetRandomRarity returns a random rarity based on the provided chances
func GetRandomRarity(chances map[ItemRarity]float64) ItemRarity {
	roll := rand.Float64()
	cumulativeChance := 0.0

	for rarity, chance := range chances {
		cumulativeChance += chance
		if roll < cumulativeChance {
			return rarity
		}
	}

	// Default to common if something goes wrong
	return RarityCommon
}

// GenerateLootFromMob generates loot from a defeated mob
func GenerateLootFromMob(mob *MobInstance, floorLevel int) []ItemInstance {
	var loot []ItemInstance

	// Always drop gold
	goldItem := ItemInstance{
		ID:          generateID(),
		Type:        ItemTypeGold,
		Name:        "Gold",
		Description: "A pile of gold coins",
		Rarity:      RarityCommon,
		Value:       mob.GoldDrop,
		Position:    mob.Position,
	}
	loot = append(loot, goldItem)

	// Determine if the mob drops an item based on difficulty
	dropChance := 0.0
	switch mob.Difficulty {
	case DifficultyEasy:
		dropChance = 0.1
	case DifficultyNormal:
		dropChance = 0.2
	case DifficultyHard:
		dropChance = 0.4
	case DifficultyElite:
		dropChance = 0.7
	case DifficultyBoss:
		dropChance = 1.0 // Bosses always drop an item
	}

	// Roll for item drop
	if rand.Float64() < dropChance {
		// Get rarity chances based on mob difficulty and floor level
		rarityChances := GetRarityChance(mob.Difficulty, floorLevel)

		// Get a random rarity
		rarity := GetRandomRarity(rarityChances)

		// Generate a random item of the determined rarity
		item := generateRandomItem(rarity, floorLevel, mob.Position)
		loot = append(loot, item)
	}

	return loot
}

// generateRandomItem generates a random item of the specified rarity
func generateRandomItem(rarity ItemRarity, floorLevel int, position Position) ItemInstance {
	// This is a simplified implementation - in a real game, you would have
	// a database of item templates to choose from based on rarity, level, etc.

	// For now, we'll just generate a basic weapon or armor
	itemTypes := []ItemType{ItemTypeWeapon, ItemTypeArmor, ItemTypePotion, ItemTypeScroll}
	itemType := itemTypes[rand.Intn(len(itemTypes))]

	var name, description string
	var value int

	// Base value based on rarity
	baseValue := 0
	switch rarity {
	case RarityCommon:
		baseValue = 5
	case RarityUncommon:
		baseValue = 15
	case RarityRare:
		baseValue = 50
	case RarityEpic:
		baseValue = 200
	case RarityLegendary:
		baseValue = 1000
	}

	// Scale value based on floor level
	value = baseValue * (1 + floorLevel/2)

	// Generate name and description based on type and rarity
	switch itemType {
	case ItemTypeWeapon:
		weapons := []string{"Sword", "Axe", "Mace", "Dagger", "Bow", "Staff", "Wand"}
		weaponType := weapons[rand.Intn(len(weapons))]
		name = rarityPrefix(rarity) + " " + weaponType
		description = "A " + string(rarity) + " " + weaponType + " that deals damage to enemies."

	case ItemTypeArmor:
		armors := []string{"Helmet", "Chestplate", "Leggings", "Boots", "Gloves", "Shield"}
		armorType := armors[rand.Intn(len(armors))]
		name = rarityPrefix(rarity) + " " + armorType
		description = "A " + string(rarity) + " " + armorType + " that protects from damage."

	case ItemTypePotion:
		potions := []string{"Healing", "Mana", "Strength", "Dexterity", "Intelligence"}
		potionType := potions[rand.Intn(len(potions))]
		name = rarityPrefix(rarity) + " Potion of " + potionType
		description = "A " + string(rarity) + " potion that enhances " + potionType + "."

	case ItemTypeScroll:
		scrolls := []string{"Fireball", "Ice Storm", "Lightning", "Teleport", "Identify"}
		scrollType := scrolls[rand.Intn(len(scrolls))]
		name = rarityPrefix(rarity) + " Scroll of " + scrollType
		description = "A " + string(rarity) + " scroll that casts " + scrollType + "."
	}

	return ItemInstance{
		ID:          generateID(),
		Type:        itemType,
		Name:        name,
		Description: description,
		Rarity:      rarity,
		Value:       value,
		Position:    position,
	}
}

// rarityPrefix returns a prefix for item names based on rarity
func rarityPrefix(rarity ItemRarity) string {
	switch rarity {
	case RarityCommon:
		return "Common"
	case RarityUncommon:
		return "Fine"
	case RarityRare:
		return "Rare"
	case RarityEpic:
		return "Epic"
	case RarityLegendary:
		return "Legendary"
	default:
		return "Common"
	}
}
