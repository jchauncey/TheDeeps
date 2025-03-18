import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import DungeonPlayground from '../components/DungeonPlayground';
import * as api from '../services/api';

// Mock the API calls
jest.mock('../services/api', () => ({
  getDungeons: jest.fn(),
  createDungeon: jest.fn(),
}));

// Mock the RoomRenderer component since we're not testing its functionality
jest.mock('../components/RoomRenderer', () => ({
  __esModule: true,
  default: () => <div data-testid="room-renderer">Room Renderer Mock</div>,
}));

describe('DungeonPlayground Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders loading state initially', () => {
    // Mock the API to return a promise that doesn't resolve immediately
    (api.getDungeons as jest.Mock).mockImplementation(() => new Promise(() => {}));
    
    render(<DungeonPlayground />);
    
    expect(screen.getByText('Loading dungeons...')).toBeInTheDocument();
  });

  test('renders empty state when no dungeons are available', async () => {
    // Mock the API to return an empty array
    (api.getDungeons as jest.Mock).mockResolvedValue([]);
    
    render(<DungeonPlayground />);
    
    await waitFor(() => {
      expect(screen.getByText('No dungeons available. Create one to get started.')).toBeInTheDocument();
    });
  });

  test('renders dungeons when they are available', async () => {
    // Mock the API to return some dungeons
    const mockDungeons = [
      { id: '1', name: 'Test Dungeon 1', floors: 3, difficulty: 'easy', createdAt: '2023-01-01', playerCount: 0 },
      { id: '2', name: 'Test Dungeon 2', floors: 5, difficulty: 'hard', createdAt: '2023-01-02', playerCount: 2 },
    ];
    
    (api.getDungeons as jest.Mock).mockResolvedValue(mockDungeons);
    
    render(<DungeonPlayground />);
    
    await waitFor(() => {
      expect(screen.getByText('easy')).toBeInTheDocument();
      expect(screen.getByText('3 floors')).toBeInTheDocument();
    });
  });

  test('creates a new dungeon successfully', async () => {
    // Mock the API calls
    (api.getDungeons as jest.Mock).mockResolvedValue([]);
    
    const newDungeon = { 
      id: '3', 
      name: 'New Test Dungeon', 
      floors: 3, 
      difficulty: 'medium', 
      createdAt: '2023-01-03', 
      playerCount: 0 
    };
    
    (api.createDungeon as jest.Mock).mockResolvedValue(newDungeon);
    
    render(<DungeonPlayground />);
    
    // Wait for the initial loading to complete
    await waitFor(() => {
      expect(screen.getByText('No dungeons available. Create one to get started.')).toBeInTheDocument();
    });
    
    // Fill in the form
    fireEvent.change(screen.getByPlaceholderText('Enter dungeon name'), { target: { value: 'New Test Dungeon' } });
    
    // Select medium difficulty
    fireEvent.change(screen.getByLabelText('Difficulty'), { target: { value: 'medium' } });
    
    // Click the create button
    fireEvent.click(screen.getByText('Create Dungeon'));
    
    // Verify the API was called with the correct parameters
    expect(api.createDungeon).toHaveBeenCalledWith({
      name: 'New Test Dungeon',
      floors: 3, // Default value
      difficulty: 'medium',
    });
    
    // Wait for the dungeon to be created and displayed
    await waitFor(() => {
      expect(screen.getByText('medium')).toBeInTheDocument();
      expect(screen.getByText('3 floors')).toBeInTheDocument();
    });
  });

  test('allows switching between floors', async () => {
    // Mock the API to return a dungeon with multiple floors
    const mockDungeon = { 
      id: '1', 
      name: 'Test Dungeon', 
      floors: 3, 
      difficulty: 'easy', 
      createdAt: '2023-01-01', 
      playerCount: 0 
    };
    
    (api.getDungeons as jest.Mock).mockResolvedValue([mockDungeon]);
    
    render(<DungeonPlayground />);
    
    // Wait for the dungeon to be displayed
    await waitFor(() => {
      expect(screen.getByText('Floor 1')).toBeInTheDocument();
      expect(screen.getByText('Floor 2')).toBeInTheDocument();
      expect(screen.getByText('Floor 3')).toBeInTheDocument();
    });
    
    // Click on Floor 2
    fireEvent.click(screen.getByText('Floor 2'));
    
    // Verify the room renderer is displayed with the correct seed
    expect(screen.getByTestId('room-renderer')).toBeInTheDocument();
  });
}); 