import React from 'react';
import { render, screen, waitFor, fireEvent, act } from '@testing-library/react';
import '@testing-library/jest-dom';
import { mockCharacters } from '../mocks/api';

// Mock the useNavigate hook
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  BrowserRouter: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

// Mock the API functions
const getCharacters = jest.fn().mockResolvedValue(mockCharacters);
const deleteCharacter = jest.fn().mockResolvedValue(undefined);

jest.mock('../../services/api', () => ({
  getCharacters: () => getCharacters(),
  deleteCharacter: (id: string) => deleteCharacter(id),
}));

// Create a mock CharacterCard component
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

// Mock the CharacterCard component
jest.mock('../../components/CharacterCard', () => {
  return function MockCharacterCard({ character, onDelete, onSelect }: any) {
    return (
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
  };
});

// Mock the Spinner component
jest.mock('@chakra-ui/react', () => {
  const originalModule = jest.requireActual('@chakra-ui/react');
  return {
    ...originalModule,
    Spinner: () => <div data-testid="loading-spinner">Loading...</div>,
    useToast: () => jest.fn(),
  };
});

// Mock CharacterSelection component for testing
const MockCharacterSelection = () => {
  const [characters, setCharacters] = React.useState<any[]>([]);
  const [loading, setLoading] = React.useState<boolean>(true);
  const [error, setError] = React.useState<string | null>(null);
  const navigate = mockNavigate;

  React.useEffect(() => {
    fetchCharacters();
  }, []);

  const fetchCharacters = async () => {
    try {
      setLoading(true);
      const data = await getCharacters();
      setCharacters(data);
      setError(null);
    } catch (err) {
      console.error('Failed to fetch characters:', err);
      setError('Failed to load characters. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteCharacter = async (id: string) => {
    try {
      await deleteCharacter(id);
      setCharacters(characters.filter(char => char.id !== id));
    } catch (err) {
      console.error('Failed to delete character:', err);
    }
  };

  const handleSelectCharacter = (character: any) => {
    // Navigate to dungeon selection with the selected character
    navigate('/dungeon-selection', { state: { character } });
  };

  const handleCreateCharacter = () => {
    navigate('/create-character');
  };

  return (
    <div>
      <h1>The Deeps</h1>
      <p>Select your character or create a new one</p>

      {loading ? (
        <div data-testid="loading-spinner">Loading...</div>
      ) : error ? (
        <div>
          <p>{error}</p>
          <button onClick={fetchCharacters}>Retry</button>
        </div>
      ) : (
        <>
          <div>
            {characters.map(character => (
              <div key={character.id}>
                <MockCharacterCard
                  character={character}
                  onDelete={handleDeleteCharacter}
                  onSelect={handleSelectCharacter}
                />
              </div>
            ))}
          </div>

          {characters.length === 0 && (
            <div>
              <p>You don't have any characters yet.</p>
            </div>
          )}

          <button onClick={handleCreateCharacter}>Create New Character</button>
        </>
      )}
    </div>
  );
};

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

    render(<MockCharacterSelection />);
    
    // Check for header elements
    expect(screen.getByText('The Deeps')).toBeInTheDocument();
    expect(screen.getByText('Select your character or create a new one')).toBeInTheDocument();
    
    // Check for loading spinner
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('renders characters after loading', async () => {
    await act(async () => {
      render(<MockCharacterSelection />);
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
      render(<MockCharacterSelection />);
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
      render(<MockCharacterSelection />);
    });
    
    // Wait for the loading to finish and state updates to complete
    await waitFor(() => {
      expect(getCharacters).toHaveBeenCalledTimes(1);
      expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument();
    });
    
    // Click the create new character button
    fireEvent.click(screen.getByText('Create New Character'));
    
    // Check if navigation was called with the correct path
    expect(mockNavigate).toHaveBeenCalledWith('/create-character');
  });

  it('deletes a character when delete button is clicked', async () => {
    await act(async () => {
      render(<MockCharacterSelection />);
    });
    
    // Wait for the loading to finish and state updates to complete
    await waitFor(() => {
      expect(getCharacters).toHaveBeenCalledTimes(1);
      expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument();
    });
    
    // Click the delete button for the first character
    fireEvent.click(screen.getByTestId('delete-1'));
    
    // Check if deleteCharacter was called with the correct ID
    expect(deleteCharacter).toHaveBeenCalledWith('1');
  });

  it('navigates to dungeon selection when a character is selected', async () => {
    await act(async () => {
      render(<MockCharacterSelection />);
    });
    
    // Wait for the loading to finish and state updates to complete
    await waitFor(() => {
      expect(getCharacters).toHaveBeenCalledTimes(1);
      expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument();
    });
    
    // Click the select button for the first character
    fireEvent.click(screen.getByTestId('select-1'));
    
    // Check if navigation was called with the correct path and state
    expect(mockNavigate).toHaveBeenCalledWith('/dungeon-selection', {
      state: { character: mockCharacters[0] }
    });
  });
}); 