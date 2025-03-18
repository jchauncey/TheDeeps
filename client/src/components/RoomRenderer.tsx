import React, { useEffect, useState } from 'react';
import {
  Box,
  Text,
  Spinner,
  Center,
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
  onLoad?: () => void;
  onError?: () => void;
  debug?: boolean;
}

const RoomRenderer: React.FC<RoomRendererProps> = ({
  roomType = 'entrance',
  width = 20,
  height = 20,
  roomWidth,
  roomHeight,
  onLoad,
  onError,
  debug = true,
}) => {
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [floor, setFloor] = useState<Floor | null>(null);
  const [rawResponse, setRawResponse] = useState<string>('');

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
        
        const url = `/test/room?${params.toString()}`;
        console.log(`Fetching room data from: ${url}`);
        
        // Fetch the test room
        const response = await fetch(url);
        
        // Clone the response for debugging
        const responseClone = response.clone();
        const rawText = await responseClone.text();
        setRawResponse(rawText);
        
        console.log(`Response status: ${response.status} ${response.statusText}`);
        
        // Log headers in a way that works with older TypeScript targets
        const headers: Record<string, string> = {};
        response.headers.forEach((value, key) => {
          headers[key] = value;
        });
        console.log('Response headers:', headers);
        
        if (!response.ok) {
          throw new Error(`Failed to fetch test room: ${response.statusText}`);
        }
        
        // Parse the JSON from the raw text to avoid double-parsing issues
        let data;
        try {
          data = JSON.parse(rawText);
          console.log('Parsed room data:', data);
        } catch (jsonError) {
          console.error('Failed to parse JSON response:', jsonError);
          console.error('Raw response text:', rawText);
          throw new Error(`Failed to parse room data: ${jsonError instanceof Error ? jsonError.message : 'Invalid JSON'}`);
        }
        
        // Validate the response structure
        if (!data || !data.tiles || !Array.isArray(data.tiles) || data.tiles.length === 0) {
          console.error('Invalid room data structure:', data);
          throw new Error('Invalid room data: missing or empty tiles array');
        }
        
        // Ensure we have a proper 2D array of tiles
        if (!Array.isArray(data.tiles[0])) {
          console.log('Converting flat tiles array to 2D array');
          // If tiles is a flat array, convert it to a 2D array
          const tilesArray: Tile[][] = [];
          for (let y = 0; y < data.height; y++) {
            const row: Tile[] = [];
            for (let x = 0; x < data.width; x++) {
              const index = y * data.width + x;
              if (index < data.tiles.length) {
                row.push(data.tiles[index]);
              } else {
                // Fill with default tile if out of bounds
                row.push({ type: 'wall', walkable: false, explored: true });
              }
            }
            tilesArray.push(row);
          }
          data.tiles = tilesArray;
        }
        
        setFloor(data);
        setLoading(false);
        onLoad?.();
      } catch (err) {
        console.error('Error fetching test room:', err);
        setError(err instanceof Error ? err.message : 'An unknown error occurred');
        setLoading(false);
        onError?.();
      }
    };
    
    fetchTestRoom();
  }, [width, height, roomWidth, roomHeight, roomType, onLoad, onError]);

  // Function to get tile color based on tile type
  const getTileColor = (tile: Tile): string => {
    if (!tile.explored) return '#000'; // Black for unexplored
    
    switch (tile.type) {
      case 'wall':
        return '#555'; // Dark gray for walls
      case 'floor':
        return '#111'; // Very dark for floor
      case 'upStairs':
        return '#00A'; // Blue for up stairs
      case 'downStairs':
        return '#A00'; // Red for down stairs
      case 'door':
        return '#850'; // Brown for doors
      case 'corridor':
        return '#222'; // Slightly lighter than floor for corridors
      default:
        return '#111'; // Default dark
    }
  };

  // Function to get tile background color for entities (mobs, items, player)
  const getEntityColor = (tile: Tile, mobs: { [key: string]: Mob }, items: { [key: string]: Item }): string | null => {
    if (tile.character) return '#FF0'; // Yellow for player
    
    if (tile.mobId && mobs && mobs[tile.mobId]) {
      const mob = mobs[tile.mobId];
      switch (mob.type) {
        case 'dragon':
          return '#F44'; // Red for dragon
        case 'ogre':
          return '#F84'; // Orange for ogre
        case 'boss':
          return '#F22'; // Bright red for boss
        case 'goblin':
          return '#8F4'; // Green for goblin
        case 'skeleton':
          return '#FFF'; // White for skeleton
        case 'shopkeeper':
          return '#F44'; // Red for shopkeeper (matching server-side definition)
        default:
          return '#F88'; // Default mob color
      }
    }
    
    if (tile.itemId && items && items[tile.itemId]) {
      const item = items[tile.itemId];
      switch (item.type) {
        case 'gold':
          return '#FF0'; // Yellow for gold
        case 'potion':
          return '#F0F'; // Magenta for potion
        case 'weapon':
          return '#0FF'; // Cyan for weapon
        case 'armor':
          return '#88F'; // Blue for armor
        default:
          return '#8F8'; // Default item color
      }
    }
    
    return null;
  };

  // Function to get tile content (for ASCII mode or debug)
  const getTileContent = (tile: Tile, mobs: { [key: string]: Mob }, items: { [key: string]: Item }): string => {
    if (!tile.explored) return ' ';
    
    if (tile.character) return '@'; // Player character
    if (tile.mobId && mobs && mobs[tile.mobId]) {
      const mob = mobs[tile.mobId];
      switch (mob.type) {
        case 'dragon':
          return 'D';
        case 'ogre':
          return 'O';
        case 'boss':
          return 'B';
        case 'goblin':
          return 'g';
        case 'skeleton':
          return 's';
        case 'shopkeeper':
          return 'S';
        default:
          return 'M'; // Generic mob
      }
    }
    if (tile.itemId && items && items[tile.itemId]) {
      const item = items[tile.itemId];
      switch (item.type) {
        case 'gold':
          return '$';
        case 'potion':
          return '!';
        case 'weapon':
          return '/';
        case 'armor':
          return '[';
        default:
          return 'i'; // Generic item
      }
    }
    
    switch (tile.type) {
      case 'wall':
        return '#';
      case 'floor':
        return '.';
      case 'upStairs':
        return '<';
      case 'downStairs':
        return '>';
      case 'door':
        return '+';
      case 'corridor':
        return '·'; // Middle dot for corridors
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

  // Check if the tiles array is valid
  const hasTiles = floor.tiles && Array.isArray(floor.tiles) && floor.tiles.length > 0;
  
  // If we don't have valid tiles, render a fallback
  if (!hasTiles) {
    return (
      <Box>
        <Text fontSize="xl" mb={4}>
          Test Room: {roomType.charAt(0).toUpperCase() + roomType.slice(1)} (Fallback)
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
          <Box 
            display="grid" 
            gridTemplateColumns={`repeat(10, 20px)`}
            gridTemplateRows={`repeat(10, 20px)`}
            gap={0}
            role="grid"
            width="fit-content"
            margin="0 auto"
            border="1px solid gray.500"
            position="relative"
            _after={{
              content: '""',
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              backgroundImage: 
                'linear-gradient(to right, gray.600 1px, transparent 1px), linear-gradient(to bottom, gray.600 1px, transparent 1px)',
              backgroundSize: '20px 20px',
              pointerEvents: 'none',
              zIndex: 10
            }}
          >
            {Array(10).fill(0).map((_, y) => 
              Array(10).fill(0).map((_, x) => {
                // Create a simple room layout
                const isWall = x === 0 || x === 9 || y === 0 || y === 9;
                const isBoss = roomType === 'boss' && x === 5 && y === 5;
                const isPlayer = x === 3 && y === 5; // Player is always at position (3,5)
                const isItem = roomType === 'treasure' && x === 5 && y === 5;
                const isStairs = roomType === 'entrance' && x === 5 && y === 5;
                const isDoor = x === 5 && y === 0;
                
                // Determine tile type and entity
                const tileColor = isWall ? '#555' : isDoor ? '#850' : '#111';
                
                // Entity color
                let entityColor = null;
                if (isBoss) entityColor = '#F22';
                if (isPlayer) entityColor = '#FF0';
                if (isItem) entityColor = '#FF0';
                
                return (
                  <Box 
                    key={`${x}-${y}`}
                    bg={tileColor}
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    width="20px"
                    height="20px"
                    border={isWall ? '1px solid #666' : 'none'}
                    position="relative"
                  >
                    {/* Entity overlay */}
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
                    
                    {/* Special tile markers */}
                    {(isStairs || isDoor) && (
                      <Box 
                        position="absolute"
                        top="0"
                        left="0"
                        width="100%"
                        height="100%"
                        display="flex"
                        alignItems="center"
                        justifyContent="center"
                        color={isStairs ? 'red.300' : 'yellow.600'}
                        fontWeight="bold"
                        zIndex={1}
                      >
                        {isStairs ? '↓' : '+'}
                      </Box>
                    )}
                    
                    {/* Show ASCII in debug mode */}
                    {debug && (
                      <Text 
                        fontSize="12px" 
                        color={isPlayer ? 'yellow' : isBoss ? 'red' : isItem ? 'green' : 'white'}
                        fontWeight="bold"
                        zIndex={3}
                        position="absolute"
                        top="50%"
                        left="50%"
                        transform="translate(-50%, -50%)"
                        textShadow={isPlayer || isBoss ? "0px 0px 2px white" : "none"}
                      >
                        {isWall ? '#' : isBoss ? 'B' : isPlayer ? '@' : isItem ? '$' : isStairs ? '>' : isDoor ? '+' : '.'}
                      </Text>
                    )}
                  </Box>
                );
              })
            )}
          </Box>
        </Box>
        
        {/* Legend */}
        <Box mt={4} display="flex" flexWrap="wrap" gap={4}>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#555" border="1px solid #666" mr={2} />
            <Text fontSize="sm">Wall</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#111" mr={2} />
            <Text fontSize="sm">Floor</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#850" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Text color="yellow.600">+</Text>
            </Box>
            <Text fontSize="sm">Door</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#111" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Box width="14px" height="14px" borderRadius="50%" bg="#FF0" />
            </Box>
            <Text fontSize="sm">Player/Item</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#111" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Box width="14px" height="14px" borderRadius="50%" bg="#F22" />
            </Box>
            <Text fontSize="sm">Monster</Text>
          </Box>
        </Box>
        
        <Box mt={4}>
          <Text>Room Information (Fallback):</Text>
          <Text>Type: {roomType}, Size: 10x10, Position: (0, 0)</Text>
        </Box>

        {debug && (
          <Box mt={4} p={4} bg="gray.700" borderRadius="md" overflowX="auto">
            <Text color="red.300" mb={2}>WARNING: Using fallback rendering because server data is invalid</Text>
            <Text color="white" mb={2}>Debug Information:</Text>
            <Text color="white" mb={2}>Room Type: {roomType}</Text>
            <Text color="white" mb={2}>Dimensions: {width}x{height}</Text>
            <Text color="white" mb={2}>Room Dimensions: {roomWidth}x{roomHeight}</Text>
            <Text color="white" mb={2}>Raw Response:</Text>
            <Box 
              as="pre" 
              p={2} 
              bg="black" 
              color="green.300" 
              fontSize="xs" 
              maxH="200px" 
              overflowY="auto"
            >
              {rawResponse.length > 1000 ? rawResponse.substring(0, 1000) + '...' : rawResponse}
            </Box>
          </Box>
        )}
      </Box>
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
        <Box 
          display="grid" 
          gridTemplateColumns={`repeat(${floor.width}, 20px)`}
          gridTemplateRows={`repeat(${floor.height}, 20px)`}
          gap={0}
          role="grid"
          width="fit-content"
          margin="0 auto"
          border="1px solid gray.500"
          position="relative"
          _after={{
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundImage: 
              'linear-gradient(to right, gray.600 1px, transparent 1px), linear-gradient(to bottom, gray.600 1px, transparent 1px)',
            backgroundSize: '20px 20px',
            pointerEvents: 'none',
            zIndex: 10
          }}
        >
          {floor.tiles.flatMap((row, y) => 
            row.map((tile, x) => {
              // Ensure tile is a valid object
              const validTile = typeof tile === 'object' && tile !== null ? tile : { 
                type: 'floor', 
                walkable: true, 
                explored: true 
              };
              
              const tileColor = getTileColor(validTile);
              const entityColor = getEntityColor(validTile, floor.mobs || {}, floor.items || {});
              const tileContent = getTileContent(validTile, floor.mobs || {}, floor.items || {});
              
              // Determine if this is a special tile (stairs, door)
              const isSpecialTile = validTile.type === 'upStairs' || validTile.type === 'downStairs' || validTile.type === 'door';
              
              return (
                <Box 
                  key={`${x}-${y}`}
                  bg={tileColor}
                  data-testid={`tile-${x}-${y}`}
                  data-tile-type={validTile.type}
                  data-tile-content={tileContent}
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  width="20px"
                  height="20px"
                  border={validTile.type === 'wall' ? '1px solid #666' : 'none'}
                  position="relative"
                  title={`${validTile.type}${validTile.mobId ? ' - ' + (floor.mobs?.[validTile.mobId]?.name || 'Monster') : ''}${validTile.itemId ? ' - ' + (floor.items?.[validTile.itemId]?.name || 'Item') : ''}`}
                >
                  {/* Entity overlay (mob, item, player) */}
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
                  
                  {/* Special tile markers */}
                  {isSpecialTile && (
                    <Box 
                      position="absolute"
                      top="0"
                      left="0"
                      width="100%"
                      height="100%"
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                      color={validTile.type === 'upStairs' ? 'blue.300' : validTile.type === 'downStairs' ? 'red.300' : 'yellow.600'}
                      fontWeight="bold"
                      zIndex={1}
                    >
                      {validTile.type === 'upStairs' ? '↑' : validTile.type === 'downStairs' ? '↓' : '+'}
                    </Box>
                  )}
                  
                  {/* Show ASCII in debug mode */}
                  {debug && (
                    <Text 
                      fontSize="12px" 
                      color={validTile.character ? 'black' : validTile.mobId ? 'red' : validTile.itemId ? 'green' : 'white'}
                      fontWeight="bold"
                      zIndex={3}
                      position="absolute"
                      top="50%"
                      left="50%"
                      transform="translate(-50%, -50%)"
                      textShadow={validTile.character ? "0px 0px 2px white" : "none"}
                    >
                      {tileContent}
                    </Text>
                  )}
                </Box>
              );
            })
          )}
        </Box>
      </Box>
      
      {/* Legend */}
      <Box mt={4} display="flex" flexWrap="wrap" gap={4}>
        <Box display="flex" alignItems="center">
          <Box width="20px" height="20px" bg="#555" border="1px solid #666" mr={2} />
          <Text fontSize="sm">Wall</Text>
        </Box>
        <Box display="flex" alignItems="center">
          <Box width="20px" height="20px" bg="#111" mr={2} />
          <Text fontSize="sm">Floor</Text>
        </Box>
        <Box display="flex" alignItems="center">
          <Box width="20px" height="20px" bg="#850" mr={2} display="flex" alignItems="center" justifyContent="center">
            <Text color="yellow.600">+</Text>
          </Box>
          <Text fontSize="sm">Door</Text>
        </Box>
        <Box display="flex" alignItems="center">
          <Box width="20px" height="20px" bg="#111" mr={2} display="flex" alignItems="center" justifyContent="center">
            <Box width="14px" height="14px" borderRadius="50%" bg="#FF0" />
          </Box>
          <Text fontSize="sm">Player</Text>
        </Box>
        <Box display="flex" alignItems="center">
          <Box width="20px" height="20px" bg="#111" mr={2} display="flex" alignItems="center" justifyContent="center">
            <Box width="14px" height="14px" borderRadius="50%" bg="#F22" />
          </Box>
          <Text fontSize="sm">Monster</Text>
        </Box>
        <Box display="flex" alignItems="center">
          <Box width="20px" height="20px" bg="#111" mr={2} display="flex" alignItems="center" justifyContent="center">
            <Box width="14px" height="14px" borderRadius="50%" bg="#FF0" />
          </Box>
          <Text fontSize="sm">Item</Text>
        </Box>
      </Box>
      
      <Box mt={4}>
        <Text>Room Information:</Text>
        {floor.rooms.map(room => (
          <Text key={room.id} data-testid={`room-info-${room.type}`}>
            Type: {room.type}, Size: {room.width}x{room.height}, Position: ({room.x}, {room.y})
          </Text>
        ))}
      </Box>

      {debug && (
        <Box mt={4} p={4} bg="gray.700" borderRadius="md" overflowX="auto">
          <Text color="white" mb={2}>Debug Information:</Text>
          <Text color="white" mb={2}>Room Type: {roomType}</Text>
          <Text color="white" mb={2}>Dimensions: {width}x{height}</Text>
          <Text color="white" mb={2}>Room Dimensions: {roomWidth}x{roomHeight}</Text>
          <Text color="white" mb={2}>Tiles Array Size: {floor.tiles.length} x {floor.tiles[0]?.length || 0}</Text>
          <Text color="white" mb={2}>Number of Rooms: {floor.rooms.length}</Text>
          <Text color="white" mb={2}>Number of Mobs: {Object.keys(floor.mobs || {}).length}</Text>
          <Text color="white" mb={2}>Number of Items: {Object.keys(floor.items || {}).length}</Text>
          
          <Text color="white" mb={2}>Tile Type Counts:</Text>
          {(() => {
            const counts: Record<string, number> = {};
            floor.tiles.forEach(row => {
              row.forEach(tile => {
                const type = tile?.type || 'unknown';
                counts[type] = (counts[type] || 0) + 1;
              });
            });
            
            return (
              <Box pl={4}>
                {Object.entries(counts).map(([type, count]) => (
                  <Text key={type} color="white" fontSize="sm">
                    {type}: {count} ({((count / (floor.width * floor.height)) * 100).toFixed(1)}%)
                  </Text>
                ))}
              </Box>
            );
          })()}
          
          <Text color="white" mb={2} mt={4}>Raw Response:</Text>
          <Box 
            as="pre" 
            p={2} 
            bg="black" 
            color="green.300" 
            fontSize="xs" 
            maxH="200px" 
            overflowY="auto"
          >
            {rawResponse.length > 1000 ? rawResponse.substring(0, 1000) + '...' : rawResponse}
          </Box>
        </Box>
      )}
    </Box>
  );
};

export default RoomRenderer; 