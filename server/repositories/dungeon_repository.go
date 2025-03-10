package repositories

import (
	"errors"
	"sync"

	"github.com/jchauncey/TheDeeps/server/models"
)

// DungeonRepository handles storage and retrieval of dungeons
type DungeonRepository struct {
	dungeons map[string]*models.Dungeon
	mutex    sync.RWMutex
}

// NewDungeonRepository creates a new dungeon repository
func NewDungeonRepository() *DungeonRepository {
	return &DungeonRepository{
		dungeons: make(map[string]*models.Dungeon),
	}
}

// GetAll returns all dungeons
func (r *DungeonRepository) GetAll() []*models.Dungeon {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	dungeons := make([]*models.Dungeon, 0, len(r.dungeons))
	for _, dungeon := range r.dungeons {
		dungeons = append(dungeons, dungeon)
	}

	return dungeons
}

// GetByID returns a dungeon by ID
func (r *DungeonRepository) GetByID(id string) (*models.Dungeon, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	dungeon, exists := r.dungeons[id]
	if !exists {
		return nil, errors.New("dungeon not found")
	}

	return dungeon, nil
}

// Save saves a dungeon
func (r *DungeonRepository) Save(dungeon *models.Dungeon) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.dungeons[dungeon.ID] = dungeon
	return nil
}

// Delete deletes a dungeon
func (r *DungeonRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.dungeons[id]; !exists {
		return errors.New("dungeon not found")
	}

	delete(r.dungeons, id)
	return nil
}

// GetFloor returns a specific floor of a dungeon
func (r *DungeonRepository) GetFloor(dungeonID string, level int) (*models.Floor, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	dungeon, exists := r.dungeons[dungeonID]
	if !exists {
		return nil, errors.New("dungeon not found")
	}

	floor, exists := dungeon.FloorData[level]
	if !exists {
		// Generate the floor if it doesn't exist
		floor = dungeon.GenerateFloor(level)
	}

	return floor, nil
}

// AddCharacterToDungeon adds a character to a dungeon
func (r *DungeonRepository) AddCharacterToDungeon(dungeonID string, characterID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	dungeon, exists := r.dungeons[dungeonID]
	if !exists {
		return errors.New("dungeon not found")
	}

	dungeon.AddCharacter(characterID)
	return nil
}

// RemoveCharacterFromDungeon removes a character from a dungeon
func (r *DungeonRepository) RemoveCharacterFromDungeon(dungeonID string, characterID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	dungeon, exists := r.dungeons[dungeonID]
	if !exists {
		return errors.New("dungeon not found")
	}

	dungeon.RemoveCharacter(characterID)
	return nil
}

// GetCharacterFloor gets the floor level for a character in a dungeon
func (r *DungeonRepository) GetCharacterFloor(dungeonID string, characterID string) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	dungeon, exists := r.dungeons[dungeonID]
	if !exists {
		return 0, errors.New("dungeon not found")
	}

	floor := dungeon.GetCharacterFloor(characterID)
	if floor == 0 {
		return 0, errors.New("character not found in dungeon")
	}

	return floor, nil
}

// SetCharacterFloor sets the floor level for a character in a dungeon
func (r *DungeonRepository) SetCharacterFloor(dungeonID string, characterID string, floor int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	dungeon, exists := r.dungeons[dungeonID]
	if !exists {
		return errors.New("dungeon not found")
	}

	dungeon.SetCharacterFloor(characterID, floor)
	return nil
}

// SaveFloor saves a floor for a dungeon
func (r *DungeonRepository) SaveFloor(dungeonID string, floorLevel int, floor *models.Floor) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	dungeon, exists := r.dungeons[dungeonID]
	if !exists {
		return errors.New("dungeon not found")
	}

	dungeon.FloorData[floorLevel] = floor
	return nil
}
