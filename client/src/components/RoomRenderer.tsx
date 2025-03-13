import React, { useEffect, useState } from 'react';
import {
  Box,
  Flex,
  Text,
  useToast,
  Spinner,
  Center,
  Grid,
} from '@chakra-ui/react';

// Define the types we need
interface Position {
  x: number;
  y: number;
}

interface Room {
  id: string;
  type: string;
  x: number;
  y: number;
  width: number;
  height: number;
  explored: boolean;
}

interface Tile {
  type: string;
  walkable: boolean;
  explored: boolean;
  roomId?: string;
  character?: string;
  mobId?: string;
  itemId?: string;
}

interface Mob {
  id: string;
  type: string;
  name: string;
  health: number;
  maxHealth: number;
  position: Position;
}

interface Item {
  id: string;
  type: string;
  name: string;
  position: Position;
}

interface Floor {
  level: number;
  width: number;
  height: number;
  tiles: Tile[][];
  rooms: Room[];
  upStairs: Position[];
  downStairs: Position[];
  mobs: { [key: string]: Mob };
  items: { [key: string]: Item };
}

interface RoomRendererProps {
  roomType?: string;
  width?: number;
  height?: number;
  roomWidth?: number;
  roomHeight?: number;
}

const RoomRenderer: React.FC<RoomRendererProps> = ({
  roomType = 'entrance',
  width = 20,
  height = 20,
  roomWidth,
  roomHeight,
}) => {
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [floor, setFloor] = useState<Floor | null>(null);
  
  const toast = useToast();

  useEffect(() => {
    const fetchTestRoom = async () => {
      try {
        setLoading(true);
        
        // Build query parameters
        const params = new URLSearchParams();
        if (roomType) params.append('type', roomType);
        if (width) params.append('width', width.toString());
        if (height) params.append('height', height.toString());
        if (roomWidth) params.append('roomWidth', roomWidth.toString());
        if (roomHeight) params.append('roomHeight', roomHeight.toString());
        
        // Fetch the test room
        const response = await fetch(`/test/room?${params.toString()}`);
        
        if (!response.ok) {
          throw new Error(`Failed to fetch test room: ${response.statusText}`);
        }
        
        const data = await response.json();
        setFloor(data);
        setLoading(false);
      } catch (err) {
        console.error('Failed to fetch test room:', err);
        setError(err instanceof Error ? err.message : 'Unknown error');
        setLoading(false);
      }
    };
    
    fetchTestRoom();
  }, [roomType, width, height, roomWidth, roomHeight]);

  // Function to get tile color based on tile type
  const getTileColor = (tile: Tile): string => {
    if (!tile.explored) return '#000'; // Black for unexplored
    
    switch (tile.type) {
      case 'wall':
        return '#555'; // Dark gray for walls
      case 'floor':
        return '#222'; // Dark for floor
      case 'upStairs':
        return '#00F'; // Blue for up stairs
      case 'downStairs':
        return '#F00'; // Red for down stairs
      default:
        return '#222'; // Default dark
    }
  };

  // Function to get tile content
  const getTileContent = (tile: Tile, mobs: { [key: string]: Mob }, items: { [key: string]: Item }): string => {
    if (!tile.explored) return ' ';
    
    if (tile.character) return '@'; // Player character
    if (tile.mobId && mobs[tile.mobId]) {
      const mob = mobs[tile.mobId];
      switch (mob.type) {
        case 'dragon':
          return 'D';
        case 'ogre':
          return 'O';
        default:
          return 'M'; // Generic mob
      }
    }
    if (tile.itemId && items[tile.itemId]) return 'i';
    
    switch (tile.type) {
      case 'wall':
        return '#';
      case 'floor':
        return '.';
      case 'upStairs':
        return '<';
      case 'downStairs':
        return '>';
      default:
        return ' ';
    }
  };

  if (loading) {
    return (
      <Center h="100%">
        <Spinner size="xl" thickness="4px" speed="0.65s" color="blue.500" />
      </Center>
    );
  }

  if (error) {
    return (
      <Center h="100%">
        <Text color="red.500">{error}</Text>
      </Center>
    );
  }

  if (!floor) {
    return (
      <Center h="100%">
        <Text>No floor data available</Text>
      </Center>
    );
  }

  return (
    <Box>
      <Text fontSize="xl" mb={4}>
        Test Room: {roomType.charAt(0).toUpperCase() + roomType.slice(1)}
      </Text>
      
      <Box 
        bg="black" 
        color="white" 
        fontFamily="monospace" 
        p={4} 
        borderRadius="md"
        overflow="auto"
        fontSize="16px"
        lineHeight="1"
      >
        <Grid templateColumns={`repeat(${floor.width}, 1fr)`} gap={0} role="grid">
          {floor.tiles.flatMap((row, y) => 
            row.map((tile, x) => (
              <Box 
                key={`${x}-${y}`}
                bg={getTileColor(tile)}
                color={tile.character ? 'yellow' : tile.mobId ? 'red' : tile.itemId ? 'green' : 'white'}
                w="20px"
                h="20px"
                display="flex"
                alignItems="center"
                justifyContent="center"
                fontWeight={tile.character || tile.mobId ? 'bold' : 'normal'}
                data-testid={`tile-${x}-${y}`}
                data-tile-type={tile.type}
                data-tile-content={getTileContent(tile, floor.mobs, floor.items)}
              >
                {getTileContent(tile, floor.mobs, floor.items)}
              </Box>
            ))
          )}
        </Grid>
      </Box>
      
      <Box mt={4}>
        <Text>Room Information:</Text>
        {floor.rooms.map(room => (
          <Text key={room.id} data-testid={`room-info-${room.type}`}>
            Type: {room.type}, Size: {room.width}x{room.height}, Position: ({room.x}, {room.y})
          </Text>
        ))}
      </Box>
    </Box>
  );
};

export default RoomRenderer; 