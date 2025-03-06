import { useEffect, useState, useRef } from 'react';
import { Box, Spinner, Text, useToast } from '@chakra-ui/react';
import { sendWebSocketMessage } from '../../services/api';

// Define tile types and their colors
const TILE_COLORS = {
  wall: '#333',
  floor: '#555',
  door: '#855',
  stairs_up: '#55f',
  stairs_down: '#f55',
};

// Define entity types and their colors
const ENTITY_COLORS = {
  player: '#ff0',
  goblin: '#0f0',
  orc: '#0a0',
  skeleton: '#fff',
  rat: '#a50',
  bat: '#a0a',
};

// Define item types and their colors
const ITEM_COLORS = {
  potion: '#f0f',
  scroll: '#ff0',
  weapon: '#aaa',
  armor: '#00f',
  gold: '#ff0',
};

interface Position {
  x: number;
  y: number;
}

interface Entity {
  id: string;
  type: string;
  name: string;
  position: Position;
}

interface Item {
  id: string;
  type: string;
  name: string;
  position: Position;
}

interface Tile {
  type: string;
  explored: boolean;
  visible: boolean;
  entity?: Entity;
  item?: Item;
}

interface Room {
  x: number;
  y: number;
  width: number;
  height: number;
}

interface Floor {
  level: number;
  width: number;
  height: number;
  tiles: Tile[][];
  rooms: Room[];
  entities: Entity[];
  items: Item[];
}

interface FloorData {
  type: string;
  floor: Floor;
  playerPosition: Position;
  currentFloor: number;
}

export const GameBoard = () => {
  const [loading, setLoading] = useState(true);
  const [floor, setFloor] = useState<Floor | null>(null);
  const [playerPos, setPlayerPos] = useState<Position | null>(null);
  const [currentFloor, setCurrentFloor] = useState(1);
  const [viewportSize, setViewportSize] = useState({ width: 0, height: 0 });
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const toast = useToast();

  // Request floor data when component mounts
  useEffect(() => {
    requestFloorData();
  }, []);

  // Calculate viewport size on mount and window resize
  useEffect(() => {
    const updateViewportSize = () => {
      // Calculate available space (accounting for padding and other elements)
      const availableHeight = window.innerHeight - 120; // Account for padding, header, etc.
      const availableWidth = window.innerWidth - 40;
      
      setViewportSize({
        width: availableWidth,
        height: availableHeight
      });
    };
    
    // Initial calculation
    updateViewportSize();
    
    // Add resize listener
    window.addEventListener('resize', updateViewportSize);
    
    return () => {
      window.removeEventListener('resize', updateViewportSize);
    };
  }, []);

  // Draw the floor when data changes or viewport size changes
  useEffect(() => {
    if (floor && playerPos && viewportSize.width > 0 && viewportSize.height > 0) {
      drawFloor();
    }
  }, [floor, playerPos, viewportSize]);

  // Handle WebSocket messages
  useEffect(() => {
    // This function would be called by the parent component when a WebSocket message is received
    const handleWebSocketMessage = (data: any) => {
      if (data.type === 'floor_data') {
        const floorData = data as FloorData;
        setFloor(floorData.floor);
        setPlayerPos(floorData.playerPosition);
        setCurrentFloor(floorData.currentFloor);
        setLoading(false);
      }
    };

    // Register the handler with the parent component
    window.addEventListener('websocket_message', (e: any) => handleWebSocketMessage(e.detail));

    return () => {
      window.removeEventListener('websocket_message', (e: any) => handleWebSocketMessage(e.detail));
    };
  }, []);

  // Ensure the game board has focus for keyboard controls
  useEffect(() => {
    const focusContainer = () => {
      if (containerRef.current) {
        containerRef.current.focus();
      }
    };

    // Focus when component mounts
    focusContainer();

    // Add click listener to focus when clicked
    if (containerRef.current) {
      containerRef.current.addEventListener('click', focusContainer);
    }

    return () => {
      if (containerRef.current) {
        containerRef.current.removeEventListener('click', focusContainer);
      }
    };
  }, []);

  // Request floor data from the server
  const requestFloorData = () => {
    sendWebSocketMessage({ type: 'get_floor' });
  };

  // Draw the floor on the canvas
  const drawFloor = () => {
    if (!floor || !canvasRef.current || !playerPos) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Calculate the visible area (viewport)
    const visibleTiles = {
      width: Math.min(floor.width, 40), // Limit to 40 tiles wide
      height: Math.min(floor.height, 25) // Limit to 25 tiles high
    };
    
    // Calculate tile size based on available space and visible area
    const tileSize = Math.min(
      Math.floor(viewportSize.width / visibleTiles.width),
      Math.floor(viewportSize.height / visibleTiles.height)
    );
    
    // Set canvas size
    canvas.width = visibleTiles.width * tileSize;
    canvas.height = visibleTiles.height * tileSize;

    // Calculate viewport center (player position)
    const viewportCenterX = Math.floor(visibleTiles.width / 2);
    const viewportCenterY = Math.floor(visibleTiles.height / 2);
    
    // Calculate top-left corner of viewport in dungeon coordinates
    const viewportStartX = Math.max(0, playerPos.x - viewportCenterX);
    const viewportStartY = Math.max(0, playerPos.y - viewportCenterY);
    
    // Adjust if we're near the edge of the map
    const maxStartX = Math.max(0, floor.width - visibleTiles.width);
    const maxStartY = Math.max(0, floor.height - visibleTiles.height);
    
    const adjustedStartX = Math.min(viewportStartX, maxStartX);
    const adjustedStartY = Math.min(viewportStartY, maxStartY);

    // Clear canvas
    ctx.fillStyle = '#000';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Draw tiles within viewport
    for (let y = 0; y < visibleTiles.height; y++) {
      for (let x = 0; x < visibleTiles.width; x++) {
        // Convert viewport coordinates to dungeon coordinates
        const dungeonX = adjustedStartX + x;
        const dungeonY = adjustedStartY + y;
        
        // Skip if out of bounds
        if (dungeonX >= floor.width || dungeonY >= floor.height) continue;
        
        const tile = floor.tiles[dungeonY][dungeonX];
        
        // Skip unexplored tiles
        if (!tile.explored && !tile.visible) continue;
        
        // Draw tile with different opacity based on visibility
        const baseColor = TILE_COLORS[tile.type as keyof typeof TILE_COLORS] || '#000';
        
        // Parse the hex color to RGB
        const r = parseInt(baseColor.slice(1, 3), 16);
        const g = parseInt(baseColor.slice(3, 5), 16);
        const b = parseInt(baseColor.slice(5, 7), 16);
        
        // Set color with opacity based on visibility
        ctx.fillStyle = tile.visible 
          ? baseColor 
          : `rgba(${r}, ${g}, ${b}, 0.5)`;
        
        ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
        
        // Draw grid lines
        ctx.strokeStyle = '#222';
        ctx.strokeRect(x * tileSize, y * tileSize, tileSize, tileSize);
      }
    }

    // Draw items within viewport
    floor.items.forEach(item => {
      // Convert dungeon coordinates to viewport coordinates
      const viewportX = item.position.x - adjustedStartX;
      const viewportY = item.position.y - adjustedStartY;
      
      // Skip if outside viewport
      if (viewportX < 0 || viewportX >= visibleTiles.width || 
          viewportY < 0 || viewportY >= visibleTiles.height) return;
      
      const tile = floor.tiles[item.position.y][item.position.x];
      
      // Skip items on unexplored tiles
      if (!tile.explored && !tile.visible) return;
      
      // Draw with different opacity based on visibility
      const baseColor = ITEM_COLORS[item.type as keyof typeof ITEM_COLORS] || '#fff';
      
      // Parse the hex color to RGB
      const r = parseInt(baseColor.slice(1, 3), 16);
      const g = parseInt(baseColor.slice(3, 5), 16);
      const b = parseInt(baseColor.slice(5, 7), 16);
      
      ctx.fillStyle = tile.visible 
        ? baseColor 
        : `rgba(${r}, ${g}, ${b}, 0.5)`;
      
      ctx.beginPath();
      ctx.arc(
        viewportX * tileSize + tileSize / 2,
        viewportY * tileSize + tileSize / 2,
        tileSize / 4,
        0,
        Math.PI * 2
      );
      ctx.fill();
    });

    // Draw entities within viewport
    floor.entities.forEach(entity => {
      // Convert dungeon coordinates to viewport coordinates
      const viewportX = entity.position.x - adjustedStartX;
      const viewportY = entity.position.y - adjustedStartY;
      
      // Skip if outside viewport
      if (viewportX < 0 || viewportX >= visibleTiles.width || 
          viewportY < 0 || viewportY >= visibleTiles.height) return;
      
      const tile = floor.tiles[entity.position.y][entity.position.x];
      
      // Only draw entities on visible tiles
      if (!tile.visible) return;
      
      ctx.fillStyle = ENTITY_COLORS[entity.type as keyof typeof ENTITY_COLORS] || '#f00';
      ctx.fillRect(
        viewportX * tileSize + tileSize / 4,
        viewportY * tileSize + tileSize / 4,
        tileSize / 2,
        tileSize / 2
      );
    });

    // Draw player
    // Convert dungeon coordinates to viewport coordinates
    const playerViewportX = playerPos.x - adjustedStartX;
    const playerViewportY = playerPos.y - adjustedStartY;
    
    ctx.fillStyle = ENTITY_COLORS.player;
    ctx.fillRect(
      playerViewportX * tileSize + tileSize / 4,
      playerViewportY * tileSize + tileSize / 4,
      tileSize / 2,
      tileSize / 2
    );
  };

  if (loading) {
    return (
      <Box
        width="100%"
        height="100%"
        display="flex"
        alignItems="center"
        justifyContent="center"
        bg="#291326"
      >
        <Spinner size="xl" color="purple.500" thickness="4px" />
        <Text ml={4} color="white" fontSize="xl">
          Loading dungeon...
        </Text>
      </Box>
    );
  }

  return (
    <Box
      ref={containerRef}
      width="100%"
      height="100%"
      bg="#291326"
      p={4}
      borderRadius="md"
      overflow="hidden"
      position="relative"
      tabIndex={0} // Make the container focusable
      outline="none" // Remove the focus outline
      _focus={{ boxShadow: "none" }} // Remove focus shadow
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
    >
      <Text color="white" mb={2} fontSize="lg">
        Floor {currentFloor}
      </Text>
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
      >
        <canvas
          ref={canvasRef}
          style={{
            imageRendering: 'pixelated',
            border: '1px solid #555',
          }}
        />
      </Box>
    </Box>
  );
}; 