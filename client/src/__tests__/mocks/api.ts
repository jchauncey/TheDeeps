import { Character, CharacterClass, Attributes } from '../../types';

// Add jest type
declare const jest: any;

// Mock character data
export const mockCharacters: Character[] = [
  {
    id: '1',
    name: 'Test Warrior',
    class: 'warrior',
    level: 1,
    experience: 0,
    attributes: {
      strength: 16,
      dexterity: 12,
      constitution: 14,
      intelligence: 8,
      wisdom: 10,
      charisma: 10
    },
    skills: {},
    maxHp: 20,
    currentHp: 20,
    maxMana: 10,
    currentMana: 10,
    gold: 100,
    currentFloor: 1,
    position: { x: 0, y: 0 },
    inventory: [],
    equipment: {}
  },
  {
    id: '2',
    name: 'Test Mage',
    class: 'mage',
    level: 1,
    experience: 0,
    attributes: {
      strength: 8,
      dexterity: 12,
      constitution: 10,
      intelligence: 16,
      wisdom: 14,
      charisma: 10
    },
    skills: {},
    maxHp: 12,
    currentHp: 12,
    maxMana: 20,
    currentMana: 20,
    gold: 100,
    currentFloor: 1,
    position: { x: 0, y: 0 },
    inventory: [],
    equipment: {}
  }
];

// Mock API functions
export const getCharacters = jest.fn().mockResolvedValue(mockCharacters);
export const getCharacter = jest.fn().mockImplementation((id: string) => {
  const character = mockCharacters.find(char => char.id === id);
  if (character) {
    return Promise.resolve(character);
  }
  return Promise.reject(new Error('Character not found'));
});

export const createCharacter = jest.fn().mockImplementation((characterData: any) => {
  const newCharacter: Character = {
    id: '3', // Mock ID generation
    name: characterData.name,
    class: characterData.class,
    level: 1,
    experience: 0,
    attributes: characterData.attributes || {
      strength: 10,
      dexterity: 10,
      constitution: 10,
      intelligence: 10,
      wisdom: 10,
      charisma: 10
    },
    skills: {},
    maxHp: 15,
    currentHp: 15,
    maxMana: 15,
    currentMana: 15,
    gold: 100,
    currentFloor: 1,
    position: { x: 0, y: 0 },
    inventory: [],
    equipment: {}
  };
  return Promise.resolve(newCharacter);
});

export const deleteCharacter = jest.fn().mockResolvedValue(undefined);

// Mock dungeon data
export const mockDungeons = [
  {
    id: 'dungeon-1',
    name: 'The Dark Cave',
    floors: 5,
    difficulty: 'easy',
    createdAt: new Date().toISOString(),
    playerCount: 0
  },
  {
    id: 'dungeon-2',
    name: 'The Forgotten Temple',
    floors: 10,
    difficulty: 'medium',
    createdAt: new Date().toISOString(),
    playerCount: 2
  }
];

// Add dungeon API mocks
export const getDungeons = jest.fn().mockResolvedValue(mockDungeons);

export const createDungeon = jest.fn().mockImplementation((dungeonData: any) => {
  const newDungeon = {
    id: `dungeon-${Date.now()}`, // Mock ID generation
    name: dungeonData.name,
    floors: dungeonData.floors,
    difficulty: dungeonData.difficulty || 'easy',
    createdAt: new Date().toISOString(),
    playerCount: 0
  };
  return Promise.resolve(newDungeon);
});

// Add joinDungeon mock
export const joinDungeon = jest.fn().mockImplementation((characterId: string, dungeonId: string) => {
  // Check if character exists
  const character = mockCharacters.find(char => char.id === characterId);
  if (!character) {
    return Promise.reject({
      isAxiosError: true,
      response: {
        status: 404,
        data: 'character not found'
      }
    });
  }
  
  // Mock dungeon check (just check if dungeonId is valid-looking)
  if (!dungeonId || dungeonId === 'non-existent-id') {
    return Promise.reject({
      isAxiosError: true,
      response: {
        status: 404,
        data: 'dungeon not found'
      }
    });
  }
  
  return Promise.resolve();
}); 