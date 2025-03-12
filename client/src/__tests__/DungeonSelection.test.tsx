import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { ChakraProvider } from '@chakra-ui/react';
import * as api from './mocks/api';
import { Character } from '../types';

// Import the actual component
import DungeonSelectionComponent from '../pages/DungeonSelection';

// Define a mock character
const mockCharacter: Character = {
  id: 'test-id',
  name: 'Test Character',
  class: 'warrior',
  level: 1,
  experience: 0,
  attributes: {
    strength: 10,
    dexterity: 10,
    constitution: 10,
    intelligence: 10,
    wisdom: 10,
    charisma: 10
  },
  maxHp: 20,
  currentHp: 20,
  maxMana: 10,
  currentMana: 10,
  gold: 0,
  currentFloor: 1,
  position: { x: 0, y: 0 },
  inventory: [],
  equipment: {},
  skills: {}
};

// Mock the API service
jest.mock('./mocks/api');

// Mock react-router-dom
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate,
  useLocation: () => ({
    state: {
      character: mockCharacter
    }
  }),
  Routes: ({ children }: { children: React.ReactNode }) => children,
  Route: ({ element }: { element: React.ReactNode }) => element,
  MemoryRouter: ({ children }: { children: React.ReactNode }) => children
}));

// Mock console.error to track error messages
const originalConsoleError = console.error;
const mockConsoleError = jest.fn();

describe('DungeonSelection Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    console.error = mockConsoleError;
    
    // Mock API responses
    api.getCharacter.mockResolvedValue(mockCharacter);
    api.getDungeons.mockResolvedValue([
      {
        id: 'test-dungeon-id',
        name: 'Test Dungeon',
        floors: 5,
        difficulty: 'easy',
        createdAt: new Date().toISOString(),
        playerCount: 0
      }
    ]);
    api.joinDungeon.mockResolvedValue(undefined);
  });

  afterEach(() => {
    console.error = originalConsoleError;
  });

  const renderComponent = () => {
    return render(
      <ChakraProvider>
        <DungeonSelectionComponent />
      </ChakraProvider>
    );
  };

  it('should verify character exists on load', async () => {
    renderComponent();
    
    await waitFor(() => {
      expect(api.getCharacter).toHaveBeenCalledWith(mockCharacter.id);
    });
  });

  it('should navigate back to character selection if character verification fails', async () => {
    // Mock character verification failure
    api.getCharacter.mockRejectedValueOnce(new Error('Character not found'));
    
    renderComponent();
    
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });

  it('should show error when trying to join dungeon with non-existent character', async () => {
    // First call to getCharacter (on component mount) succeeds
    api.getCharacter.mockResolvedValueOnce(mockCharacter);
    
    // Second call to getCharacter (during join dungeon) fails
    api.getCharacter.mockRejectedValueOnce(new Error('Character not found'));
    
    const { findByText, findAllByText } = renderComponent();
    
    // Wait for the dungeon list to load
    await findByText('Test Dungeon');
    
    // Find and click the select button in the table row
    const selectButtons = await findAllByText('Select');
    fireEvent.click(selectButtons[0]);
    
    // Find and click the join dungeon button
    const joinButton = await findByText('Join Selected Dungeon');
    fireEvent.click(joinButton);
    
    // Verify the error was logged
    await waitFor(() => {
      expect(mockConsoleError).toHaveBeenCalledWith(
        'Character verification failed:',
        expect.objectContaining({ message: 'Character not found' })
      );
    });
    
    // Verify navigation back to character selection
    expect(mockNavigate).toHaveBeenCalledWith('/');
  });
}); 