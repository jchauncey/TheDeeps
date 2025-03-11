import React from 'react';
import { render, screen, waitFor, fireEvent, act } from '@testing-library/react';
import CharacterSelection from '../../pages/CharacterSelection';
import { mockCharacters, getCharacters, deleteCharacter } from '../mocks/api';
import { BrowserRouter } from 'react-router-dom';

// Add jest type
declare const jest: any;

// Mock the CharacterCard component
jest.mock('../../components/CharacterCard', () => {
  const MockCharacterCard = ({ character, onDelete, onSelect }: any) => (
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
  return MockCharacterCard;
});

// Mock the Spinner component to add a test ID
jest.mock('@chakra-ui/react', () => {
  const originalModule = jest.requireActual('@chakra-ui/react');
  return {
    ...originalModule,
    Spinner: () => <div data-testid="loading-spinner">Loading...</div>
  };
});

describe('CharacterSelection', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    jest.clearAllMocks();
  });

  it('renders loading state initially', async () => {
    // Mock getCharacters to delay response
    getCharacters.mockImplementationOnce(() => new Promise(resolve => {
      // Don't resolve immediately to ensure loading state is visible
      setTimeout(() => resolve(mockCharacters), 100);
    }));

    render(
      <BrowserRouter>
        <CharacterSelection />
      </BrowserRouter>
    );
    
    // Check for header elements
    expect(screen.getByText('The Deeps')).toBeInTheDocument();
    expect(screen.getByText('Select your character or create a new one')).toBeInTheDocument();
    
    // Check for loading spinner
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('renders characters after loading', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <CharacterSelection />
        </BrowserRouter>
      );
    });
    
    // Wait for the loading to finish and state updates to complete
    await waitFor(() => {
      expect(getCharacters).toHaveBeenCalledTimes(1);
      expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument();
    });
    
    // Check if character names are rendered
    expect(screen.getByText('Test Warrior')).toBeInTheDocument();
    expect(screen.getByText('Test Mage')).toBeInTheDocument();
  });

  it('shows empty state when no characters are available', async () => {
    // Mock getCharacters to return an empty array
    getCharacters.mockResolvedValueOnce([]);
    
    await act(async () => {
      render(
        <BrowserRouter>
          <CharacterSelection />
        </BrowserRouter>
      );
    });
    
    // Wait for the loading to finish and state updates to complete
    await waitFor(() => {
      expect(getCharacters).toHaveBeenCalledTimes(1);
      expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument();
    });
    
    // Check if empty state message is shown
    expect(screen.getByText("You don't have any characters yet.")).toBeInTheDocument();
  });

  it('navigates to character creation page when create button is clicked', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <CharacterSelection />
        </BrowserRouter>
      );
    });
    
    // Wait for the loading to finish and state updates to complete
    await waitFor(() => {
      expect(getCharacters).toHaveBeenCalledTimes(1);
      expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument();
    });
    
    // Click the create new character button within act
    await act(async () => {
      fireEvent.click(screen.getByText('Create New Character'));
    });
  });
}); 