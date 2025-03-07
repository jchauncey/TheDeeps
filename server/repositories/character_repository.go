package repositories

import (
	"errors"
	"sync"

	"github.com/jchauncey/TheDeeps/server/models"
)

var (
	// ErrCharacterNotFound is returned when a character is not found
	ErrCharacterNotFound = errors.New("character not found")
)

// CharacterRepository handles character storage
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

// Create creates a new character
func (r *CharacterRepository) Create(character *models.Character) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.characters[character.ID] = character
	return nil
}

// GetByID retrieves a character by ID
func (r *CharacterRepository) GetByID(id string) (*models.Character, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	character, ok := r.characters[id]
	if !ok {
		return nil, ErrCharacterNotFound
	}

	return character, nil
}

// Update updates an existing character
func (r *CharacterRepository) Update(character *models.Character) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.characters[character.ID]
	if !ok {
		return ErrCharacterNotFound
	}

	r.characters[character.ID] = character
	return nil
}

// Delete deletes a character by ID
func (r *CharacterRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.characters[id]
	if !ok {
		return ErrCharacterNotFound
	}

	delete(r.characters, id)
	return nil
}

// GetAll retrieves all characters
func (r *CharacterRepository) GetAll() []*models.Character {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	characters := make([]*models.Character, 0, len(r.characters))
	for _, character := range r.characters {
		characters = append(characters, character)
	}

	return characters
}
