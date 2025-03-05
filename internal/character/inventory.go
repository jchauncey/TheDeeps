package character

// ItemType represents the type of an item
type ItemType int

const (
	ItemWeapon ItemType = iota
	ItemArmor
	ItemConsumable
	ItemKey
	ItemTreasure
)

// Item represents an item in the game
type Item struct {
	Name        string
	Description string
	Type        ItemType
	Equipment   *Equipment // If the item is equippable, this will be non-nil
	Value       int
}

// Inventory represents a container for items
type Inventory struct {
	Name     string
	Capacity int
	Items    []Item
}

// NewInventory creates a new inventory with the given name and capacity
func NewInventory(name string, capacity int) Inventory {
	return Inventory{
		Name:     name,
		Capacity: capacity,
		Items:    []Item{},
	}
}

// AddItem adds an item to the inventory if there's space
// Returns true if successful, false if inventory is full
func (i *Inventory) AddItem(item Item) bool {
	if len(i.Items) >= i.Capacity {
		return false
	}
	i.Items = append(i.Items, item)
	return true
}

// RemoveItem removes an item from inventory at the given index
// Returns the removed item or an empty item if index is invalid
func (i *Inventory) RemoveItem(index int) Item {
	if index < 0 || index >= len(i.Items) {
		return Item{}
	}

	item := i.Items[index]
	i.Items = append(i.Items[:index], i.Items[index+1:]...)
	return item
}

// GetItem returns the item at the given index
// Returns an empty item if index is invalid
func (i *Inventory) GetItem(index int) Item {
	if index < 0 || index >= len(i.Items) {
		return Item{}
	}
	return i.Items[index]
}

// IsFull returns true if the inventory is full
func (i *Inventory) IsFull() bool {
	return len(i.Items) >= i.Capacity
}

// IsEmpty returns true if the inventory is empty
func (i *Inventory) IsEmpty() bool {
	return len(i.Items) == 0
}
