export type CharacterClass = 
  | 'warrior'
  | 'mage'
  | 'rogue'
  | 'cleric'
  | 'druid'
  | 'warlock'
  | 'bard'
  | 'paladin'
  | 'ranger'
  | 'monk'
  | 'barbarian'
  | 'sorcerer';

export interface Attributes {
  strength: number;
  dexterity: number;
  constitution: number;
  intelligence: number;
  wisdom: number;
  charisma: number;
}

export interface Position {
  x: number;
  y: number;
}

export interface Item {
  id: string;
  name: string;
  type: string;
  value: number;
  weight: number;
  level: number;
  attributes: Record<string, number>;
}

export interface Equipment {
  weapon?: Item;
  armor?: Item;
  accessory?: Item;
}

export interface Character {
  id: string;
  name: string;
  class: CharacterClass;
  level: number;
  experience: number;
  attributes: Attributes;
  skills: Record<string, any>;
  maxHp: number;
  currentHp: number;
  maxMana: number;
  currentMana: number;
  gold: number;
  currentFloor: number;
  currentDungeon?: string;
  position: Position;
  inventory: Item[];
  equipment: Equipment;
}

export interface Dungeon {
  id: string;
  name: string;
  floors: number;
  difficulty: string;
  createdAt: string;
  playerCount: number;
} 