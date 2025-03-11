import React from 'react';
import { Character } from '../../types';

interface CharacterCardProps {
  character: Character;
  onDelete: (id: string) => void;
  onSelect: (character: Character) => void;
}

const MockCharacterCard: React.FC<CharacterCardProps> = ({ character, onDelete, onSelect }) => (
  <div data-testid={`character-card-${character.id}`}>
    <div>{character.name}</div>
    <div>{character.class}</div>
    <button 
      onClick={() => onDelete(character.id)} 
      data-testid={`delete-${character.id}`}
    >
      Delete
    </button>
    <button 
      onClick={() => onSelect(character)} 
      data-testid={`select-${character.id}`}
    >
      Select
    </button>
  </div>
);

export default MockCharacterCard; 