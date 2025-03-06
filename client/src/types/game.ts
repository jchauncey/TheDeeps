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
}

// Game state interface
export interface GameState {
  player: CharacterData | null;
  gameStarted: boolean;
  level: number;
  health: number;
  maxHealth: number;
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