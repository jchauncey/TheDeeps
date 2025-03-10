import { CharacterClass } from '../types/game';

// Character class constants
export const CHARACTER_CLASSES: CharacterClass[] = [
  {
    id: 'warrior',
    name: 'Warrior',
    description: 'A skilled fighter and master of weapons',
    primaryAbility: 'Strength',
    savingThrows: ['Strength', 'Constitution'],
    hitDie: 10,
    abilities: ['Second Wind', 'Action Surge', 'Combat Style'],
    proficiencies: ['All Armor', 'Shields', 'Simple Weapons', 'Martial Weapons'],
    recommendedStats: {
      strength: 15,
      dexterity: 12,
      constitution: 14,
      intelligence: 8,
      wisdom: 10,
      charisma: 8
    }
  },
  {
    id: 'mage',
    name: 'Mage',
    description: 'A scholarly magic-user capable of manipulating arcane forces',
    primaryAbility: 'Intelligence',
    savingThrows: ['Intelligence', 'Wisdom'],
    hitDie: 6,
    abilities: ['Arcane Recovery', 'Spellcasting', 'Arcane Tradition'],
    proficiencies: ['Daggers', 'Darts', 'Slings', 'Quarterstaffs', 'Light Crossbows'],
    recommendedStats: {
      strength: 8,
      dexterity: 14,
      constitution: 12,
      intelligence: 15,
      wisdom: 10,
      charisma: 8
    }
  },
  {
    id: 'rogue',
    name: 'Rogue',
    description: 'A scoundrel who uses stealth and trickery to overcome obstacles and enemies',
    primaryAbility: 'Dexterity',
    savingThrows: ['Dexterity', 'Intelligence'],
    hitDie: 8,
    abilities: ['Expertise', 'Sneak Attack', 'Thieves\' Cant', 'Cunning Action'],
    proficiencies: ['Light Armor', 'Simple Weapons', 'Hand Crossbows', 'Longswords', 'Rapiers', 'Shortswords'],
    recommendedStats: {
      strength: 8,
      dexterity: 15,
      constitution: 12,
      intelligence: 12,
      wisdom: 10,
      charisma: 10
    }
  }
];

// Attribute descriptions
export const ATTRIBUTE_DESCRIPTIONS = {
  strength: 'Physical power, affects melee damage and carrying capacity',
  dexterity: 'Agility and reflexes, affects AC, ranged attacks, and initiative',
  constitution: 'Endurance and stamina, affects hit points and resistance to poison',
  intelligence: 'Knowledge and reasoning, affects arcane magic and investigation',
  wisdom: 'Perception and intuition, affects divine magic and insight',
  charisma: 'Force of personality, affects social interactions and certain magic'
}; 