import React, { useEffect, useState } from 'react';
import { Box, Spinner, Text, Center, useColorModeValue, Button } from '@chakra-ui/react';
import RoomRenderer from './RoomRenderer';

// Define types for floor data
interface Tile {
  type: string;
  walkable?: boolean;
  explored?: boolean;
  character?: string;
  mobId?: string;
  itemId?: string;
  roomId?: string;
}

interface FloorData {
  level: number;
  width: number;
  height: number;
  tiles: Tile[][];
}

interface FloorVisualizerProps {
  floorData: FloorData;
}

// Define the room visualization component
const FloorVisualizer: React.FC<FloorVisualizerProps> = ({ floorData }) => {
  // Colors matching the RoomRenderer component
  const bgColor = 'gray.900';
  const floorColor = '#222';
  const wallColor = '#333';
  const playerColor = 'blue.500';
  const itemColor = 'yellow.500';
  const mobColor = 'red.500';
  const upStairsColor = 'blue.300';
  const downStairsColor = 'red.300';
  
  // Default size if not provided
  const width = floorData.width || 20;
  const height = floorData.height || 20;
  
  if (!floorData || !floorData.tiles) {
    return <Text>Invalid floor data</Text>;
  }
  
  return (
    <Box 
      width="100%" 
      height="100%" 
      bg={bgColor} 
      p={2} 
      borderRadius="md" 
      overflow="auto"
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="flex-start"
      color="white"
      position="relative"
      maxHeight="100%"
    >
      <Box 
        display="inline-block"
        overflowY="auto"
        overflowX="auto"
        maxHeight="calc(100% - 40px)"
        padding={2}
        position="relative"
      >
        {Array.isArray(floorData.tiles) && floorData.tiles.map((row: Tile[], y: number) => (
          <Box key={`row-${y}`} display="flex" height="20px">
            {Array.isArray(row) && row.map((tile: Tile, x: number) => {
              let tileColor = floorColor;
              let entityColor = null;
              let symbol = '';
              
              // Determine tile color (background)
              switch(tile.type) {
                case '#':
                case 'wall':
                  tileColor = wallColor;
                  break;
                case 'upStairs':
                case '<':
                  tileColor = floorColor;
                  symbol = '<';
                  break;
                case 'downStairs':
                case '>':
                  tileColor = floorColor;
                  symbol = '>';
                  break;
                default:
                  tileColor = floorColor;
              }
              
              // Determine entity color (for overlay)
              if (tile.character) {
                entityColor = playerColor;
              } else if (tile.mobId) {
                entityColor = mobColor;
              } else if (tile.itemId) {
                entityColor = itemColor;
              } else if (tile.type === 'upStairs' || tile.type === '<') {
                entityColor = upStairsColor;
              } else if (tile.type === 'downStairs' || tile.type === '>') {
                entityColor = downStairsColor;
              }
              
              // Render in a style similar to RoomRenderer
              return (
                <Box 
                  key={`tile-${x}-${y}`} 
                  width="20px" 
                  height="20px" 
                  minWidth="20px"
                  minHeight="20px"
                  bg={tileColor}
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  border={tile.type === 'wall' || tile.type === '#' ? '1px solid #666' : 'none'}
                  position="relative"
                  data-tile-type={tile.type}
                >
                  {/* Entity overlay (mob, item, player, stairs) */}
                  {entityColor && (
                    <Box 
                      position="absolute"
                      top="3px"
                      left="3px"
                      width="14px"
                      height="14px"
                      borderRadius="50%"
                      bg={entityColor}
                      zIndex={2}
                    />
                  )}
                  
                  {/* Special symbols for stairs */}
                  {symbol && (
                    <Box 
                      position="absolute"
                      top="0"
                      left="0"
                      width="100%"
                      height="100%"
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                      color={tile.type === 'upStairs' || tile.type === '<' ? 'blue.300' : 'red.300'}
                      fontWeight="bold"
                      zIndex={1}
                      fontSize="xs"
                    >
                      {symbol}
                    </Box>
                  )}
                </Box>
              );
            })}
          </Box>
        ))}
      </Box>
      
      <Box mt={2} fontSize="sm" position="absolute" bottom={2} left={0} right={0} textAlign="center">
        <Text fontWeight="bold">Floor {floorData.level}</Text>
        <Text>Width: {width}, Height: {height}</Text>
      </Box>
    </Box>
  );
};

interface DungeonRoomRendererProps {
  dungeonId: string;
  floorNumber: number;
  width?: number;
  height?: number;
  debug?: boolean;
}

const DungeonRoomRenderer: React.FC<DungeonRoomRendererProps> = ({
  dungeonId,
  floorNumber,
  width = 20,
  height = 20,
  debug = true
}) => {
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [floorData, setFloorData] = useState<FloorData | null>(null);
  
  const loadFloorData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      // Fetch floor data from the server API using the correct endpoint path
      console.log(`Fetching floor ${floorNumber} data from server...`);
      // Update the API endpoint path to match what's registered on the server
      const response = await fetch(`/dungeons/${dungeonId}/floor/${floorNumber}`);
      
      if (!response.ok) {
        throw new Error(`Server returned ${response.status}: ${response.statusText}`);
      }
      
      const serverFloorData = await response.json();
      console.log(`Successfully fetched floor ${floorNumber} data from server:`, serverFloorData);
      
      setFloorData(serverFloorData);
    } catch (err) {
      console.error(`Error fetching floor ${floorNumber} data:`, err);
      setError(err instanceof Error ? err.message : 'Failed to fetch floor data from server');
    } finally {
      setLoading(false);
    }
  };
  
  useEffect(() => {
    loadFloorData();
  }, [dungeonId, floorNumber]);
  
  if (loading) {
    return (
      <Center height="100%" width="100%" bg="gray.50" borderRadius="md" border="1px solid" borderColor="gray.200">
        <Spinner size="xl" color="blue.500" thickness="4px" />
        <Text ml={4} fontWeight="medium">Loading floor data from server...</Text>
      </Center>
    );
  }
  
  if (error) {
    return (
      <Center 
        height="100%" 
        width="100%" 
        bg="red.50" 
        borderRadius="md" 
        border="1px solid" 
        borderColor="red.200"
        flexDirection="column"
        p={4}
      >
        <Text fontWeight="bold" mb={2} color="red.500">Error Loading Floor Data</Text>
        <Text>{error}</Text>
        <Button mt={4} colorScheme="blue" size="sm" onClick={loadFloorData}>
          Try Again
        </Button>
      </Center>
    );
  }
  
  if (!floorData) {
    return (
      <Center height="100%" width="100%" bg="gray.50" borderRadius="md" border="1px solid" borderColor="gray.200">
        <Text>No floor data available</Text>
      </Center>
    );
  }
  
  return (
    <Box 
      position="relative" 
      height="100%" 
      width="100%" 
      border="1px solid" 
      borderColor="gray.200" 
      borderRadius="md"
      overflow="auto"
    >
      <FloorVisualizer floorData={floorData} />
    </Box>
  );
};

export default DungeonRoomRenderer;