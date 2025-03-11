package repositories

import (
	"errors"
	"sync"

	"github.com/jchauncey/TheDeeps/server/models"
)

// CharacterRepository handles storage and retrieval of characters
type CharacterRepository struct {
	characters map[string]*models.Character
	mutex      sync.RWMutex
}

// NewCharacterRepository creates a new character repository
func NewCharacterRepository() *CharacterRepository {
	return &CharacterRepository{
		characters: make(map[string]*models.Character),
	}
}

// GetAll returns all characters
func (r *CharacterRepository) GetAll() []*models.Character {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	characters := make([]*models.Character, 0, len(r.characters))
	for _, character := range r.characters {
		characters = append(characters, character)
	}

	return characters
}

// GetByID returns a character by ID
func (r *CharacterRepository) GetByID(id string) (*models.Character, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	character, exists := r.characters[id]
	if !exists {
		return nil, errors.New("character not found")
	}

	return character, nil
}

// Save saves a character
func (r *CharacterRepository) Save(character *models.Character) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.characters[character.ID] = character
	return nil
}

// Delete deletes a character
func (r *CharacterRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.characters[id]; !exists {
		return errors.New("character not found")
	}

	delete(r.characters, id)
	return nil
}

// Count returns the number of characters
func (r *CharacterRepository) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.characters)
}
