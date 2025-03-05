package character

// ClassType represents a character class
type ClassType int

const (
	Warrior ClassType = iota
	Wizard
	Rogue
	Ranger
	Cleric
)

// Class represents a character class with its attributes
type Class struct {
	Type        ClassType
	Name        string
	Description string

	// Class bonuses
	StrBonus int
	WisBonus int
	ConBonus int
	DexBonus int
	ChaBonus int

	// Starting equipment bonuses
	StartingWeaponDmg int
	StartingArmorVal  int

	// Special abilities
	SpecialAbility string
}

// Player represents the player character
type Player struct {
	X, Y   int
	HP     int
	MaxHP  int
	Name   string
	Symbol rune

	// Character class
	Class *Class

	// Base Stats
	Strength     int
	Wisdom       int
	Constitution int
	Dexterity    int
	Charisma     int

	// Equipment
	Weapon Equipment
	Armor  Equipment

	// Inventory
	Bag  *Bag
	Gold int
}

// Equipment represents equippable items
type Equipment struct {
	Name      string
	ArmorVal  int
	DamageVal int
}

// Bag represents the player's inventory
type Bag struct {
	Name     string
	Capacity int
	Items    []*Item
}

// GetClassName returns the string name of a class type
func GetClassName(classType ClassType) string {
	switch classType {
	case Warrior:
		return "Warrior"
	case Wizard:
		return "Wizard"
	case Rogue:
		return "Rogue"
	case Ranger:
		return "Ranger"
	case Cleric:
		return "Cleric"
	default:
		return "Unknown"
	}
}

// GetClassByType returns a class configuration for the given class type
func GetClassByType(classType ClassType) *Class {
	switch classType {
	case Warrior:
		return &Class{
			Type:              Warrior,
			Name:              "Warrior",
			Description:       "A strong fighter skilled in melee combat",
			StrBonus:          3,
			WisBonus:          0,
			ConBonus:          2,
			DexBonus:          1,
			ChaBonus:          0,
			StartingWeaponDmg: 5, // Starts with a better weapon
			StartingArmorVal:  3, // Starts with better armor
			SpecialAbility:    "Berserk: Can enter a rage to deal more damage",
		}
	case Wizard:
		return &Class{
			Type:              Wizard,
			Name:              "Wizard",
			Description:       "A master of arcane magic",
			StrBonus:          0,
			WisBonus:          3,
			ConBonus:          0,
			DexBonus:          1,
			ChaBonus:          2,
			StartingWeaponDmg: 2, // Staff
			StartingArmorVal:  1, // Robes
			SpecialAbility:    "Fireball: Can cast a powerful fire spell",
		}
	case Rogue:
		return &Class{
			Type:              Rogue,
			Name:              "Rogue",
			Description:       "A stealthy thief skilled in deception",
			StrBonus:          1,
			WisBonus:          1,
			ConBonus:          0,
			DexBonus:          3,
			ChaBonus:          1,
			StartingWeaponDmg: 4, // Dagger with bonus
			StartingArmorVal:  2, // Light armor
			SpecialAbility:    "Backstab: Deal extra damage when attacking from behind",
		}
	case Ranger:
		return &Class{
			Type:              Ranger,
			Name:              "Ranger",
			Description:       "A skilled hunter and wilderness expert",
			StrBonus:          1,
			WisBonus:          1,
			ConBonus:          1,
			DexBonus:          2,
			ChaBonus:          1,
			StartingWeaponDmg: 4, // Bow
			StartingArmorVal:  2, // Light armor
			SpecialAbility:    "Eagle Eye: Can spot enemies from a distance",
		}
	case Cleric:
		return &Class{
			Type:              Cleric,
			Name:              "Cleric",
			Description:       "A holy warrior with healing powers",
			StrBonus:          1,
			WisBonus:          2,
			ConBonus:          2,
			DexBonus:          0,
			ChaBonus:          1,
			StartingWeaponDmg: 3, // Mace
			StartingArmorVal:  3, // Medium armor
			SpecialAbility:    "Heal: Can restore health to self or allies",
		}
	default:
		return nil
	}
}

// Calculate total armor value based on equipment and stats
func (p *Player) GetArmorValue() int {
	baseArmor := p.Armor.ArmorVal
	dexBonus := (p.Dexterity - 10) / 2 // Each 2 points above 10 gives +1
	return baseArmor + dexBonus
}

// Calculate damage value based on equipment and stats
func (p *Player) GetDamageValue() int {
	baseDamage := p.Weapon.DamageVal
	strBonus := (p.Strength - 10) / 2 // Each 2 points above 10 gives +1
	return baseDamage + strBonus
}

// UseSpecialAbility uses the character's class special ability
// Returns a description of what happened and success status
func (p *Player) UseSpecialAbility() (string, bool) {
	if p.Class == nil {
		return "You have no special abilities.", false
	}

	switch p.Class.Type {
	case Warrior:
		// Berserk increases damage temporarily
		p.Weapon.DamageVal += 2
		return "You enter a berserker rage, increasing your damage!", true
	case Wizard:
		// Fireball does area damage (would need to be implemented in game logic)
		return "You cast a powerful fireball!", true
	case Rogue:
		// Backstab does extra damage (would need position checking in game logic)
		return "You prepare to strike from the shadows!", true
	case Ranger:
		// Eagle Eye reveals the map (would need to be implemented in game logic)
		return "Your keen eyes spot hidden details in your surroundings!", true
	case Cleric:
		// Heal restores health
		healAmount := 5 + (p.Wisdom-10)/2
		p.HP += healAmount
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
		return "Divine energy flows through you, healing your wounds!", true
	default:
		return "Your special ability fails.", false
	}
}

// AddToInventory adds an item to the player's inventory if there's space
// Returns true if successful, false if inventory is full
func (p *Player) AddToInventory(item *Item) bool {
	if len(p.Bag.Items) >= p.Bag.Capacity {
		return false
	}
	p.Bag.Items = append(p.Bag.Items, item)
	return true
}

// RemoveFromInventory removes an item from inventory at the given index
// Returns the removed item or nil if index is invalid
func (p *Player) RemoveFromInventory(index int) *Item {
	if index < 0 || index >= len(p.Bag.Items) {
		return nil
	}

	item := p.Bag.Items[index]
	p.Bag.Items = append(p.Bag.Items[:index], p.Bag.Items[index+1:]...)
	return item
}

// EquipWeapon equips a weapon from inventory
func (p *Player) EquipWeapon(item *Item) bool {
	if item.Type != ItemWeapon || item.Equipment == nil {
		return false
	}
	p.Weapon = *item.Equipment
	return true
}

// EquipArmor equips armor from inventory
func (p *Player) EquipArmor(item *Item) bool {
	if item.Type != ItemArmor || item.Equipment == nil {
		return false
	}
	p.Armor = *item.Equipment
	return true
}

// UseItem uses a consumable item
// Returns true if the item was used successfully
func (p *Player) UseItem(item *Item) bool {
	if item.Type != ItemConsumable {
		return false
	}

	// Handle different consumable effects
	switch item.Name {
	case "Healing Potion":
		p.HP += item.Value
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
	// Add more item effects here
	default:
		return false
	}

	return true
}

// NewPlayer creates a new player with starting equipment and stats
func NewPlayer(startX, startY int, classType ClassType) *Player {
	// Get class configuration
	class := GetClassByType(classType)
	if class == nil {
		// Default to warrior if invalid class
		class = GetClassByType(Warrior)
	}

	// Create starting equipment based on class
	startingWeapon := Equipment{
		Name:      getClassWeaponName(classType),
		DamageVal: class.StartingWeaponDmg,
	}

	startingArmor := Equipment{
		Name:     getClassArmorName(classType),
		ArmorVal: class.StartingArmorVal,
	}

	// Create backpack
	backpack := &Bag{
		Name:     "Backpack",
		Capacity: 10,
		Items:    make([]*Item, 0),
	}

	// Add a healing potion to inventory
	healingPotion := &Item{
		Name:        "Healing Potion",
		Description: "Restores 10 HP when consumed",
		Type:        ItemConsumable,
		Value:       10,
	}

	backpack.Items = append(backpack.Items, healingPotion)

	// Base stats (before class bonuses)
	baseStr := 10
	baseWis := 10
	baseCon := 10
	baseDex := 10
	baseCha := 10

	// Calculate max HP based on class and constitution
	conBonus := (baseCon + class.ConBonus - 10) / 2
	maxHP := 15 + conBonus

	// Add class-specific starting item if applicable
	addClassSpecificItem(backpack, classType)

	return &Player{
		X:      startX,
		Y:      startY,
		HP:     maxHP,
		MaxHP:  maxHP,
		Name:   "Adventurer",
		Symbol: '@',

		// Set class
		Class: class,

		// Starting stats with class bonuses
		Strength:     baseStr + class.StrBonus,
		Wisdom:       baseWis + class.WisBonus,
		Constitution: baseCon + class.ConBonus,
		Dexterity:    baseDex + class.DexBonus,
		Charisma:     baseCha + class.ChaBonus,

		// Starting equipment
		Weapon: startingWeapon,
		Armor:  startingArmor,

		// Inventory
		Bag:  backpack,
		Gold: 10,
	}
}

// Helper functions for class-specific items and equipment

func getClassWeaponName(classType ClassType) string {
	switch classType {
	case Warrior:
		return "Longsword"
	case Wizard:
		return "Staff"
	case Rogue:
		return "Dagger"
	case Ranger:
		return "Bow"
	case Cleric:
		return "Mace"
	default:
		return "Dagger"
	}
}

func getClassArmorName(classType ClassType) string {
	switch classType {
	case Warrior:
		return "Chain Mail"
	case Wizard:
		return "Cloth Robes"
	case Rogue:
		return "Leather Armor"
	case Ranger:
		return "Hide Armor"
	case Cleric:
		return "Scale Mail"
	default:
		return "Leather Armor"
	}
}

func addClassSpecificItem(bag *Bag, classType ClassType) {
	var item *Item

	switch classType {
	case Warrior:
		item = &Item{
			Name:        "Sharpening Stone",
			Description: "Can be used to improve weapon damage",
			Type:        ItemConsumable,
			Value:       2, // Adds 2 damage when used
		}
	case Wizard:
		item = &Item{
			Name:        "Spell Scroll",
			Description: "Contains a powerful spell",
			Type:        ItemConsumable,
			Value:       15, // Does 15 damage when used
		}
	case Rogue:
		item = &Item{
			Name:        "Lockpicks",
			Description: "Used to open locked doors and chests",
			Type:        ItemKey,
			Value:       1,
		}
	case Ranger:
		item = &Item{
			Name:        "Hunting Trap",
			Description: "Can be set to damage enemies",
			Type:        ItemConsumable,
			Value:       8, // Does 8 damage when used
		}
	case Cleric:
		item = &Item{
			Name:        "Holy Symbol",
			Description: "Enhances healing abilities",
			Type:        ItemConsumable,
			Value:       5, // Adds 5 HP to healing when used
		}
	default:
		return
	}

	if item != nil {
		bag.Items = append(bag.Items, item)
	}
}
