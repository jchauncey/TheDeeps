package repositories

import (
	"errors"
	"sync"
	"time"

	"github.com/jchauncey/TheDeeps/server/models"
)

// DungeonRepository errors
var (
	ErrDungeonNotFound = errors.New("dungeon not found")
)

// DungeonRepository manages dungeon instances
type DungeonRepository struct {
	dungeons map[string]*models.DungeonInstance
	mu       sync.RWMutex
}

// NewDungeonRepository creates a new dungeon repository
func NewDungeonRepository() *DungeonRepository {
	return &DungeonRepository{
		dungeons: make(map[string]*models.DungeonInstance),
	}
}

// Create creates a new dungeon instance
func (r *DungeonRepository) Create(name string, numFloors int) *models.DungeonInstance {
	r.mu.Lock()
	defer r.mu.Unlock()

	dungeon := models.NewDungeonInstance(name, numFloors)
	r.dungeons[dungeon.ID] = dungeon
	return dungeon
}

// GetByID gets a dungeon instance by ID
func (r *DungeonRepository) GetByID(id string) (*models.DungeonInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dungeon, exists := r.dungeons[id]
	if !exists {
		return nil, ErrDungeonNotFound
	}
	return dungeon, nil
}

// GetAll gets all dungeon instances
func (r *DungeonRepository) GetAll() []*models.DungeonInstance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dungeons := make([]*models.DungeonInstance, 0, len(r.dungeons))
	for _, dungeon := range r.dungeons {
		dungeons = append(dungeons, dungeon)
	}
	return dungeons
}

// Delete deletes a dungeon instance
func (r *DungeonRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.dungeons[id]; !exists {
		return ErrDungeonNotFound
	}
	delete(r.dungeons, id)
	return nil
}

// CleanupInactive removes inactive dungeon instances
func (r *DungeonRepository) CleanupInactive(maxInactivity time.Duration) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	count := 0
	for id, dungeon := range r.dungeons {
		if !dungeon.IsActive(maxInactivity) && !dungeon.HasPlayers() {
			delete(r.dungeons, id)
			count++
		}
	}
	return count
}

// AddPlayerToDungeon adds a player to a dungeon instance
func (r *DungeonRepository) AddPlayerToDungeon(dungeonID, characterID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	dungeon, exists := r.dungeons[dungeonID]
	if !exists {
		return ErrDungeonNotFound
	}

	dungeon.AddPlayer(characterID)
	return nil
}

// RemovePlayerFromDungeon removes a player from a dungeon instance
func (r *DungeonRepository) RemovePlayerFromDungeon(dungeonID, characterID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	dungeon, exists := r.dungeons[dungeonID]
	if !exists {
		return ErrDungeonNotFound
	}

	dungeon.RemovePlayer(characterID)
	return nil
}

// GetPlayerDungeon gets the dungeon instance a player is in
func (r *DungeonRepository) GetPlayerDungeon(characterID string) *models.DungeonInstance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, dungeon := range r.dungeons {
		if _, exists := dungeon.Players[characterID]; exists {
			return dungeon
		}
	}
	return nil
}
