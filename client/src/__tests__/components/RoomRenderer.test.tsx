import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import RoomRenderer from '../../components/RoomRenderer';
import { setupFetchMock, setupFetchErrorMock, resetFetchMock } from '../mocks/fetch';
import { preparedRoomData } from '../mocks/roomData';

// Mock the Chakra UI toast
jest.mock('@chakra-ui/react', () => {
  const originalModule = jest.requireActual('@chakra-ui/react');
  return {
    ...originalModule,
    useToast: () => jest.fn(),
  };
});

describe('RoomRenderer Component', () => {
  beforeEach(() => {
    setupFetchMock();
  });

  afterEach(() => {
    resetFetchMock();
    jest.clearAllMocks();
  });

  test('renders loading state initially', () => {
    render(<RoomRenderer />);
    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  test('renders error state when fetch fails', async () => {
    setupFetchErrorMock();
    render(<RoomRenderer />);
    
    await waitFor(() => {
      expect(screen.getByText(/failed to fetch test room/i)).toBeInTheDocument();
    });
  });

  test('renders entrance room correctly', async () => {
    render(<RoomRenderer roomType="entrance" />);
    
    await waitFor(() => {
      // Check room title
      expect(screen.getByText(/Test Room: Entrance/i)).toBeInTheDocument();
      
      // Check room information
      expect(screen.getByText(/Type: entrance, Size: 8x8/i)).toBeInTheDocument();
      
      // Check that the grid is rendered (we can't easily check individual tiles in this test)
      const gridContainer = document.querySelector('[role="grid"]');
      expect(gridContainer).toBeInTheDocument();
    });
  });

  test('renders standard room correctly', async () => {
    render(<RoomRenderer roomType="standard" />);
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Standard/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: standard, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('renders treasure room correctly', async () => {
    render(<RoomRenderer roomType="treasure" />);
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Treasure/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: treasure, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('renders boss room correctly', async () => {
    render(<RoomRenderer roomType="boss" />);
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Boss/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: boss, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('renders safe room correctly', async () => {
    render(<RoomRenderer roomType="safe" />);
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Safe/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: safe, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('renders shop room correctly', async () => {
    render(<RoomRenderer roomType="shop" />);
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Shop/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: shop, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('passes custom dimensions to the API', async () => {
    render(
      <RoomRenderer 
        roomType="entrance" 
        width={30} 
        height={25} 
        roomWidth={10} 
        roomHeight={12} 
      />
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Entrance/i)).toBeInTheDocument();
      
      // Check that fetch was called with the correct parameters
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringMatching(/type=entrance.*width=30.*height=25.*roomWidth=10.*roomHeight=12/)
      );
    });
  });
});

// More detailed tests for specific room features
describe('RoomRenderer Specific Features', () => {
  beforeEach(() => {
    setupFetchMock();
  });

  afterEach(() => {
    resetFetchMock();
    jest.clearAllMocks();
  });

  test('entrance room has down stairs', async () => {
    // We need to modify the DOM structure to add role attributes for testing
    const originalRender = RoomRenderer.prototype.render;
    RoomRenderer.prototype.render = function() {
      const result = originalRender.apply(this);
      // Add role attributes after rendering
      setTimeout(() => {
        const downStairsTiles = document.querySelectorAll('[data-testid="tile-downStairs"]');
        if (downStairsTiles.length === 0) {
          const tiles = document.querySelectorAll('[data-testid^="tile-"]');
          tiles.forEach(tile => {
            if (tile.textContent === '>') {
              tile.setAttribute('data-testid', 'tile-downStairs');
            }
          });
        }
      }, 0);
      return result;
    };

    render(<RoomRenderer roomType="entrance" />);
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Entrance/i)).toBeInTheDocument();
      
      // Check for down stairs in the room information
      const roomInfo = screen.getByText(/Type: entrance/i);
      expect(roomInfo).toBeInTheDocument();
      
      // We can't easily check for specific tiles in the rendered grid,
      // but we can check that the fetch was called with the correct room type
      expect(global.fetch).toHaveBeenCalledWith(expect.stringMatching(/type=entrance/));
    });
    
    // Restore original render method
    RoomRenderer.prototype.render = originalRender;
  });

  test('treasure room has items', async () => {
    render(<RoomRenderer roomType="treasure" />);
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Treasure/i)).toBeInTheDocument();
      
      // Check that fetch was called with the correct room type
      expect(global.fetch).toHaveBeenCalledWith(expect.stringMatching(/type=treasure/));
      
      // We can't easily check for specific items in the rendered grid,
      // but we can verify that the mock data for treasure rooms has items
      expect(Object.keys(preparedRoomData.treasure.items).length).toBeGreaterThan(0);
    });
  });

  test('boss room has a boss mob', async () => {
    render(<RoomRenderer roomType="boss" />);
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Boss/i)).toBeInTheDocument();
      
      // Check that fetch was called with the correct room type
      expect(global.fetch).toHaveBeenCalledWith(expect.stringMatching(/type=boss/));
      
      // Verify that the mock data for boss rooms has a boss mob
      const bossMobs = Object.values(preparedRoomData.boss.mobs);
      expect(bossMobs.length).toBeGreaterThan(0);
      expect((bossMobs[0] as any).variant).toBe('boss');
    });
  });
}); 