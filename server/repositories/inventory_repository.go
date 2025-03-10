package repositories

import (
	"sync"

	"github.com/jchauncey/TheDeeps/server/models"
)

// InventoryRepository handles storage and retrieval of inventory items
type InventoryRepository struct {
	items     map[string]*models.Item
	itemsLock sync.RWMutex
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{
		items: make(map[string]*models.Item),
	}
}

// SaveItem saves an item to the repository
func (r *InventoryRepository) SaveItem(item *models.Item) error {
	r.itemsLock.Lock()
	defer r.itemsLock.Unlock()

	r.items[item.ID] = item
	return nil
}

// GetItem retrieves an item from the repository by ID
func (r *InventoryRepository) GetItem(itemID string) (*models.Item, bool) {
	r.itemsLock.RLock()
	defer r.itemsLock.RUnlock()

	item, exists := r.items[itemID]
	return item, exists
}

// DeleteItem removes an item from the repository
func (r *InventoryRepository) DeleteItem(itemID string) bool {
	r.itemsLock.Lock()
	defer r.itemsLock.Unlock()

	if _, exists := r.items[itemID]; exists {
		delete(r.items, itemID)
		return true
	}
	return false
}

// GetAllItems returns all items in the repository
func (r *InventoryRepository) GetAllItems() []*models.Item {
	r.itemsLock.RLock()
	defer r.itemsLock.RUnlock()

	items := make([]*models.Item, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	return items
}

// GenerateRandomItems generates a specified number of random items based on floor level
func (r *InventoryRepository) GenerateRandomItems(count int, floorLevel int) []*models.Item {
	items := make([]*models.Item, 0, count)
	for i := 0; i < count; i++ {
		item := models.GenerateRandomItem(floorLevel)
		r.SaveItem(item)
		items = append(items, item)
	}
	return items
}
