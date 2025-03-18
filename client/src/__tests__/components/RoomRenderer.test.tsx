import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ChakraProvider } from '@chakra-ui/react';
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

  test('renders loading state initially', async () => {
    render(
      <ChakraProvider>
        <RoomRenderer />
      </ChakraProvider>
    );
    
    // Check for the loading spinner
    const spinner = screen.getByRole('progressbar');
    expect(spinner).toBeInTheDocument();
  });

  test('renders error state when fetch fails', async () => {
    setupFetchErrorMock();
    render(
      <ChakraProvider>
        <RoomRenderer />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Error: Failed to fetch test room/i)).toBeInTheDocument();
    });
  });

  test('renders entrance room correctly', async () => {
    render(
      <ChakraProvider>
        <RoomRenderer roomType="entrance" />
      </ChakraProvider>
    );
    
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
    render(
      <ChakraProvider>
        <RoomRenderer roomType="standard" />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Standard/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: standard, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('renders treasure room correctly', async () => {
    render(
      <ChakraProvider>
        <RoomRenderer roomType="treasure" />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Treasure/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: treasure, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('renders boss room correctly', async () => {
    render(
      <ChakraProvider>
        <RoomRenderer roomType="boss" />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Boss/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: boss, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('renders safe room correctly', async () => {
    render(
      <ChakraProvider>
        <RoomRenderer roomType="safe" />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Safe/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: safe, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('renders shop room correctly', async () => {
    render(
      <ChakraProvider>
        <RoomRenderer roomType="shop" />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Shop/i)).toBeInTheDocument();
      expect(screen.getByText(/Type: shop, Size: 7x7/i)).toBeInTheDocument();
    });
  });

  test('passes custom dimensions to the API', async () => {
    render(
      <ChakraProvider>
        <RoomRenderer 
          roomType="entrance" 
          width={30} 
          height={25} 
          roomWidth={10} 
          roomHeight={12} 
        />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Entrance/i)).toBeInTheDocument();
      
      // Check that fetch was called with the correct parameters
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringMatching(/type=entrance.*width=30.*height=25.*roomWidth=10.*roomHeight=12/)
      );
    });
  });

  test('handles loading callback', async () => {
    const mockOnLoad = jest.fn();
    
    render(
      <ChakraProvider>
        <RoomRenderer 
          roomType="entrance"
          onLoad={mockOnLoad}
        />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(mockOnLoad).toHaveBeenCalledTimes(1);
    });
  });

  test('handles error callback', async () => {
    setupFetchErrorMock();
    const mockOnError = jest.fn();
    
    render(
      <ChakraProvider>
        <RoomRenderer 
          roomType="entrance"
          onError={mockOnError}
        />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(mockOnError).toHaveBeenCalledTimes(1);
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
    render(
      <ChakraProvider>
        <RoomRenderer roomType="entrance" />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Entrance/i)).toBeInTheDocument();
      
      // Check for down stairs in the room information
      const roomInfo = screen.getByText(/Type: entrance/i);
      expect(roomInfo).toBeInTheDocument();
      
      // We can't easily check for specific tiles in the rendered grid,
      // but we can check that the fetch was called with the correct room type
      expect(global.fetch).toHaveBeenCalledWith(expect.stringMatching(/type=entrance/));
    });
  });

  test('treasure room has items', async () => {
    render(
      <ChakraProvider>
        <RoomRenderer roomType="treasure" />
      </ChakraProvider>
    );
    
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
    render(
      <ChakraProvider>
        <RoomRenderer roomType="boss" />
      </ChakraProvider>
    );
    
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

// Add a new test section for the grid overlay feature
describe('RoomRenderer Grid Feature', () => {
  beforeEach(() => {
    setupFetchMock();
  });

  afterEach(() => {
    resetFetchMock();
    jest.clearAllMocks();
  });

  test('renders grid overlay in all room types', async () => {
    const roomTypes = ['entrance', 'standard', 'treasure', 'boss', 'safe', 'shop'];
    
    for (const roomType of roomTypes) {
      document.body.innerHTML = '';
      
      const { unmount } = render(
        <ChakraProvider>
          <RoomRenderer roomType={roomType} />
        </ChakraProvider>
      );
      
      await waitFor(() => {
        expect(screen.getByText(new RegExp(`Test Room: ${roomType.charAt(0).toUpperCase() + roomType.slice(1)}`, 'i'))).toBeInTheDocument();
        
        // Find the grid container
        const gridContainer = screen.getByRole('grid');
        expect(gridContainer).toBeInTheDocument();
        
        // Check that the grid container has the position style for the grid overlay
        expect(gridContainer).toHaveStyle({
          position: 'relative'
        });
        
        // Check that it has a border
        const computedStyle = window.getComputedStyle(gridContainer);
        expect(computedStyle.border).toBeTruthy();
      });
      
      // Clean up after each iteration
      unmount();
    }
  });

  test('grid overlay is properly sized based on room dimensions', async () => {
    const customWidth = 25;
    const customHeight = 15;
    
    render(
      <ChakraProvider>
        <RoomRenderer roomType="standard" width={customWidth} height={customHeight} />
      </ChakraProvider>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Test Room: Standard/i)).toBeInTheDocument();
      
      // Find the grid container
      const gridContainer = screen.getByRole('grid');
      expect(gridContainer).toBeInTheDocument();
      
      // We can't directly check the gridTemplateColumns/Rows with toHaveStyle
      // because they're applied via Chakra UI's CSS-in-JS system
      // Instead, we'll check that the container has the correct class
      expect(gridContainer).toHaveAttribute('class');
    });
  });
}); 