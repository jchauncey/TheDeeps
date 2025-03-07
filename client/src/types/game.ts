// Character data structure
export interface CharacterData {
  name: string;
  characterClass: string;
  stats: {
    strength: number;
    dexterity: number;
    constitution: number;
    intelligence: number;
    wisdom: number;
    charisma: number;
  };
  abilities: string[];
  proficiencies: string[];
}

// Character class definitions
export interface CharacterClass {
  id: string;
  name: string;
  description: string;
  primaryAbility: string;
  savingThrows: string[];
  hitDie: number;
  abilities: string[];
  proficiencies: string[];
  recommendedStats: {
    strength: number;
    dexterity: number;
    constitution: number;
    intelligence: number;
    wisdom: number;
    charisma: number;
  };
}

// D&D Character Classes
export const CHARACTER_CLASSES: CharacterClass[] = [
  {
    id: 'barbarian',
    name: 'Barbarian',
    description: 'A fierce warrior who can enter a battle rage',
    primaryAbility: 'Strength',
    savingThrows: ['Strength', 'Constitution'],
    hitDie: 12,
    abilities: ['Rage', 'Unarmored Defense', 'Reckless Attack'],
    proficiencies: ['Light Armor', 'Medium Armor', 'Shields', 'Simple Weapons', 'Martial Weapons'],
    recommendedStats: {
      strength: 15,
      dexterity: 14,
      constitution: 14,
      intelligence: 8,
      wisdom: 10,
      charisma: 8
    }
  },
  {
    id: 'bard',
    name: 'Bard',
    description: 'An inspiring magician whose power echoes the music of creation',
    primaryAbility: 'Charisma',
    savingThrows: ['Dexterity', 'Charisma'],
    hitDie: 8,
    abilities: ['Spellcasting', 'Bardic Inspiration', 'Jack of All Trades'],
    proficiencies: ['Light Armor', 'Simple Weapons', 'Hand Crossbows', 'Longswords', 'Rapiers', 'Shortswords'],
    recommendedStats: {
      strength: 8,
      dexterity: 14,
      constitution: 12,
      intelligence: 10,
      wisdom: 10,
      charisma: 15
    }
  },
  {
    id: 'cleric',
    name: 'Cleric',
    description: 'A priestly champion who wields divine magic in service of a higher power',
    primaryAbility: 'Wisdom',
    savingThrows: ['Wisdom', 'Charisma'],
    hitDie: 8,
    abilities: ['Spellcasting', 'Divine Domain', 'Channel Divinity'],
    proficiencies: ['Light Armor', 'Medium Armor', 'Shields', 'Simple Weapons'],
    recommendedStats: {
      strength: 10,
      dexterity: 10,
      constitution: 12,
      intelligence: 8,
      wisdom: 15,
      charisma: 12
    }
  },
  {
    id: 'druid',
    name: 'Druid',
    description: 'A priest of the Old Faith, wielding the powers of nature and adopting animal forms',
    primaryAbility: 'Wisdom',
    savingThrows: ['Intelligence', 'Wisdom'],
    hitDie: 8,
    abilities: ['Spellcasting', 'Wild Shape', 'Druid Circle'],
    proficiencies: ['Light Armor', 'Medium Armor', 'Shields', 'Clubs', 'Daggers', 'Darts', 'Javelins', 'Maces', 'Quarterstaffs', 'Scimitars', 'Sickles', 'Slings', 'Spears'],
    recommendedStats: {
      strength: 8,
      dexterity: 12,
      constitution: 12,
      intelligence: 10,
      wisdom: 15,
      charisma: 10
    }
  },
  {
    id: 'fighter',
    name: 'Fighter',
    description: 'A master of martial combat, skilled with a variety of weapons and armor',
    primaryAbility: 'Strength or Dexterity',
    savingThrows: ['Strength', 'Constitution'],
    hitDie: 10,
    abilities: ['Fighting Style', 'Second Wind', 'Action Surge'],
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
    id: 'monk',
    name: 'Monk',
    description: 'A master of martial arts, harnessing the power of the body in pursuit of physical and spiritual perfection',
    primaryAbility: 'Dexterity & Wisdom',
    savingThrows: ['Strength', 'Dexterity'],
    hitDie: 8,
    abilities: ['Unarmored Defense', 'Martial Arts', 'Ki', 'Unarmored Movement'],
    proficiencies: ['Simple Weapons', 'Shortswords'],
    recommendedStats: {
      strength: 10,
      dexterity: 15,
      constitution: 12,
      intelligence: 8,
      wisdom: 14,
      charisma: 8
    }
  },
  {
    id: 'paladin',
    name: 'Paladin',
    description: 'A holy warrior bound to a sacred oath',
    primaryAbility: 'Strength & Charisma',
    savingThrows: ['Wisdom', 'Charisma'],
    hitDie: 10,
    abilities: ['Divine Sense', 'Lay on Hands', 'Divine Smite', 'Sacred Oath'],
    proficiencies: ['All Armor', 'Shields', 'Simple Weapons', 'Martial Weapons'],
    recommendedStats: {
      strength: 15,
      dexterity: 8,
      constitution: 12,
      intelligence: 8,
      wisdom: 10,
      charisma: 14
    }
  },
  {
    id: 'ranger',
    name: 'Ranger',
    description: 'A warrior who combats threats on the edges of civilization',
    primaryAbility: 'Dexterity & Wisdom',
    savingThrows: ['Strength', 'Dexterity'],
    hitDie: 10,
    abilities: ['Favored Enemy', 'Natural Explorer', 'Ranger Archetype'],
    proficiencies: ['Light Armor', 'Medium Armor', 'Shields', 'Simple Weapons', 'Martial Weapons'],
    recommendedStats: {
      strength: 10,
      dexterity: 15,
      constitution: 12,
      intelligence: 8,
      wisdom: 14,
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
  },
  {
    id: 'sorcerer',
    name: 'Sorcerer',
    description: 'A spellcaster who draws on inherent magic from a gift or bloodline',
    primaryAbility: 'Charisma',
    savingThrows: ['Constitution', 'Charisma'],
    hitDie: 6,
    abilities: ['Spellcasting', 'Sorcerous Origin', 'Font of Magic', 'Metamagic'],
    proficiencies: ['Daggers', 'Darts', 'Slings', 'Quarterstaffs', 'Light Crossbows'],
    recommendedStats: {
      strength: 8,
      dexterity: 12,
      constitution: 14,
      intelligence: 10,
      wisdom: 8,
      charisma: 15
    }
  },
  {
    id: 'warlock',
    name: 'Warlock',
    description: 'A wielder of magic that is derived from a bargain with an extraplanar entity',
    primaryAbility: 'Charisma',
    savingThrows: ['Wisdom', 'Charisma'],
    hitDie: 8,
    abilities: ['Otherworldly Patron', 'Pact Magic', 'Eldritch Invocations', 'Pact Boon'],
    proficiencies: ['Light Armor', 'Simple Weapons'],
    recommendedStats: {
      strength: 8,
      dexterity: 12,
      constitution: 12,
      intelligence: 10,
      wisdom: 10,
      charisma: 15
    }
  },
  {
    id: 'wizard',
    name: 'Wizard',
    description: 'A scholarly magic-user capable of manipulating the structures of reality',
    primaryAbility: 'Intelligence',
    savingThrows: ['Intelligence', 'Wisdom'],
    hitDie: 6,
    abilities: ['Spellcasting', 'Arcane Recovery', 'Arcane Tradition', 'Spell Mastery'],
    proficiencies: ['Daggers', 'Darts', 'Slings', 'Quarterstaffs', 'Light Crossbows'],
    recommendedStats: {
      strength: 8,
      dexterity: 12,
      constitution: 12,
      intelligence: 15,
      wisdom: 10,
      charisma: 8
    }
  }
];

// Game state interface
export interface GameState {
  player: CharacterData | null;
  gameStarted: boolean;
  level: number;
  health: number;
  maxHealth: number;
  mana: number;
  maxMana: number;
  experience: number;
  inventory: Item[];
}

// Item interface
export interface Item {
  id: string;
  name: string;
  type: 'weapon' | 'armor' | 'potion' | 'scroll' | 'misc';
  description: string;
  value: number;
  effects?: {
    [key: string]: number;
  };
}

// Debug message interface
export interface DebugMessage {
  message: string;
  level: 'info' | 'warning' | 'error' | 'debug';
  timestamp: number;
} 