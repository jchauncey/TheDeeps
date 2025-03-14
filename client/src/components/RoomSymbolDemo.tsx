import React, { useState } from 'react';
import {
  Box,
  SimpleGrid,
  Text,
  Heading,
  VStack,
  HStack,
  Divider,
  Select,
  FormControl,
  FormLabel,
  Switch,
  Button,
} from '@chakra-ui/react';

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
  position: { x: number, y: number };
}

interface Item {
  id: string;
  type: string;
  name: string;
  position: { x: number, y: number };
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

const RoomSymbolDemo: React.FC = () => {
  const [demoType, setDemoType] = useState<string>('standard');
  const [debugMode, setDebugMode] = useState<boolean>(true);
  
  // Define tile types
  const tileTypes = [
    { value: 'wall', label: 'Wall' },
    { value: 'floor', label: 'Floor' },
    { value: 'door', label: 'Door' },
    { value: 'upStairs', label: 'Up Stairs' },
    { value: 'downStairs', label: 'Down Stairs' },
    { value: 'corridor', label: 'Corridor' },
  ];
  
  // Define entity types
  const entityTypes = [
    { value: 'player', label: 'Player' },
    { value: 'dragon', label: 'Dragon' },
    { value: 'ogre', label: 'Ogre' },
    { value: 'boss', label: 'Boss' },
    { value: 'goblin', label: 'Goblin' },
    { value: 'skeleton', label: 'Skeleton' },
    { value: 'shopkeeper', label: 'Shopkeeper' },
    { value: 'gold', label: 'Gold' },
    { value: 'potion', label: 'Potion' },
    { value: 'weapon', label: 'Weapon' },
    { value: 'armor', label: 'Armor' },
  ];
  
  // Define demo types
  const demoTypes = [
    { value: 'standard', label: 'Standard Room' },
    { value: 'entrance', label: 'Entrance Room' },
    { value: 'boss', label: 'Boss Room' },
    { value: 'treasure', label: 'Treasure Room' },
    { value: 'shop', label: 'Shop Room' },
    { value: 'mixed', label: 'Mixed Entities' },
  ];
  
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
          return '#F44'; // Red for shopkeeper
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

  // Generate a demo room based on the selected type
  const generateDemoRoom = () => {
    const width = 10;
    const height = 10;
    const tiles: Tile[][] = [];
    const mobs: { [key: string]: Mob } = {};
    const items: { [key: string]: Item } = {};
    
    // Initialize tiles with walls around the edges and floor in the middle
    for (let y = 0; y < height; y++) {
      tiles[y] = [];
      for (let x = 0; x < width; x++) {
        const isWall = x === 0 || x === width - 1 || y === 0 || y === height - 1;
        tiles[y][x] = {
          type: isWall ? 'wall' : 'floor',
          walkable: !isWall,
          explored: true,
        };
      }
    }
    
    // Add a door at the top
    tiles[0][width / 2] = {
      type: 'door',
      walkable: true,
      explored: true,
    };
    
    // Add room-specific features
    switch (demoType) {
      case 'entrance':
        // Add down stairs in the center
        tiles[height / 2][width / 2] = {
          type: 'downStairs',
          walkable: true,
          explored: true,
        };
        // Add player character
        tiles[height / 2 + 2][width / 2] = {
          type: 'floor',
          walkable: true,
          explored: true,
          character: 'player1',
        };
        break;
        
      case 'boss':
        // Add boss in the center
        const bossId = 'boss1';
        mobs[bossId] = {
          id: bossId,
          type: 'boss',
          name: 'Dragon Boss',
          health: 100,
          maxHealth: 100,
          position: { x: width / 2, y: height / 2 },
        };
        tiles[height / 2][width / 2] = {
          type: 'floor',
          walkable: true,
          explored: true,
          mobId: bossId,
        };
        break;
        
      case 'treasure':
        // Add various items
        for (let i = 0; i < 3; i++) {
          const itemId = `item${i}`;
          const itemX = width / 2 - 1 + i;
          const itemY = height / 2;
          const itemType = i === 0 ? 'gold' : i === 1 ? 'weapon' : 'armor';
          
          items[itemId] = {
            id: itemId,
            type: itemType,
            name: `${itemType.charAt(0).toUpperCase() + itemType.slice(1)} Item`,
            position: { x: itemX, y: itemY },
          };
          
          tiles[itemY][itemX] = {
            type: 'floor',
            walkable: true,
            explored: true,
            itemId: itemId,
          };
        }
        break;
        
      case 'shop':
        // Add shopkeeper
        const shopkeeperId = 'shopkeeper1';
        mobs[shopkeeperId] = {
          id: shopkeeperId,
          type: 'shopkeeper',
          name: 'Shopkeeper',
          health: 50,
          maxHealth: 50,
          position: { x: width / 2, y: height / 2 - 1 },
        };
        tiles[height / 2 - 1][width / 2] = {
          type: 'floor',
          walkable: true,
          explored: true,
          mobId: shopkeeperId,
        };
        
        // Add shop items
        for (let i = 0; i < 3; i++) {
          const itemId = `shop${i}`;
          const itemX = width / 2 - 1 + i;
          const itemY = height / 2 + 1;
          const itemType = i === 0 ? 'potion' : i === 1 ? 'weapon' : 'armor';
          
          items[itemId] = {
            id: itemId,
            type: itemType,
            name: `Shop ${itemType.charAt(0).toUpperCase() + itemType.slice(1)}`,
            position: { x: itemX, y: itemY },
          };
          
          tiles[itemY][itemX] = {
            type: 'floor',
            walkable: true,
            explored: true,
            itemId: itemId,
          };
        }
        break;
        
      case 'mixed':
        // Add a variety of entities
        
        // Player
        tiles[3][3] = {
          type: 'floor',
          walkable: true,
          explored: true,
          character: 'player1',
        };
        
        // Monsters
        const monsterTypes = ['dragon', 'ogre', 'goblin', 'skeleton'];
        for (let i = 0; i < monsterTypes.length; i++) {
          const mobId = `mob${i}`;
          const mobType = monsterTypes[i];
          const mobX = 3 + i;
          const mobY = 5;
          
          mobs[mobId] = {
            id: mobId,
            type: mobType,
            name: `${mobType.charAt(0).toUpperCase() + mobType.slice(1)}`,
            health: 50,
            maxHealth: 50,
            position: { x: mobX, y: mobY },
          };
          
          tiles[mobY][mobX] = {
            type: 'floor',
            walkable: true,
            explored: true,
            mobId: mobId,
          };
        }
        
        // Items
        const itemTypes = ['gold', 'potion', 'weapon', 'armor'];
        for (let i = 0; i < itemTypes.length; i++) {
          const itemId = `item${i}`;
          const itemType = itemTypes[i];
          const itemX = 3 + i;
          const itemY = 7;
          
          items[itemId] = {
            id: itemId,
            type: itemType,
            name: `${itemType.charAt(0).toUpperCase() + itemType.slice(1)}`,
            position: { x: itemX, y: itemY },
          };
          
          tiles[itemY][itemX] = {
            type: 'floor',
            walkable: true,
            explored: true,
            itemId: itemId,
          };
        }
        
        // Special tiles
        tiles[2][7] = {
          type: 'upStairs',
          walkable: true,
          explored: true,
        };
        
        tiles[7][2] = {
          type: 'downStairs',
          walkable: true,
          explored: true,
        };
        break;
        
      default: // standard room
        // Add a few random monsters
        const mobId = 'mob1';
        mobs[mobId] = {
          id: mobId,
          type: 'goblin',
          name: 'Goblin',
          health: 20,
          maxHealth: 20,
          position: { x: width / 2, y: height / 2 },
        };
        tiles[height / 2][width / 2] = {
          type: 'floor',
          walkable: true,
          explored: true,
          mobId: mobId,
        };
        break;
    }
    
    return { tiles, mobs, items };
  };
  
  const { tiles, mobs, items } = generateDemoRoom();
  
  return (
    <VStack spacing={4} align="stretch">
      <Text>
        This component demonstrates the different symbols used in the map rendering system
        in the context of different room types.
      </Text>
      
      <HStack spacing={4} wrap="wrap">
        <FormControl display="flex" alignItems="center" width="auto">
          <FormLabel htmlFor="demo-type" mb={0}>Room Type:</FormLabel>
          <Select 
            id="demo-type"
            value={demoType}
            onChange={(e) => setDemoType(e.target.value)}
            width="auto"
            ml={2}
          >
            {demoTypes.map(type => (
              <option key={type.value} value={type.value}>{type.label}</option>
            ))}
          </Select>
        </FormControl>
        
        <FormControl display="flex" alignItems="center" width="auto">
          <FormLabel htmlFor="debug-mode" mb={0}>Show Symbols:</FormLabel>
          <Switch 
            id="debug-mode" 
            isChecked={debugMode} 
            onChange={() => setDebugMode(!debugMode)} 
            colorScheme="red"
          />
        </FormControl>
      </HStack>
      
      <Divider />
      
      <Box>
        <Heading size="sm" mb={3} color="blue.300">
          {demoType.charAt(0).toUpperCase() + demoType.slice(1)} Room
        </Heading>
        
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
            gridTemplateColumns={`repeat(${tiles[0].length}, 20px)`}
            gridTemplateRows={`repeat(${tiles.length}, 20px)`}
            gap={0}
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
            {tiles.flatMap((row, y) => 
              row.map((tile, x) => {
                const tileColor = getTileColor(tile);
                const entityColor = getEntityColor(tile, mobs, items);
                const tileContent = getTileContent(tile, mobs, items);
                
                // Determine if this is a special tile (stairs, door)
                const isSpecialTile = tile.type === 'upStairs' || tile.type === 'downStairs' || tile.type === 'door';
                
                return (
                  <Box 
                    key={`${x}-${y}`}
                    bg={tileColor}
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    width="20px"
                    height="20px"
                    border={tile.type === 'wall' ? '1px solid #666' : 'none'}
                    position="relative"
                    title={`${tile.type}${tile.mobId ? ' - ' + (mobs?.[tile.mobId]?.name || 'Monster') : ''}${tile.itemId ? ' - ' + (items?.[tile.itemId]?.name || 'Item') : ''}`}
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
                        color={tile.type === 'upStairs' ? 'blue.300' : tile.type === 'downStairs' ? 'red.300' : 'yellow.600'}
                        fontWeight="bold"
                        zIndex={1}
                      >
                        {tile.type === 'upStairs' ? '↑' : tile.type === 'downStairs' ? '↓' : '+'}
                      </Box>
                    )}
                    
                    {/* Show ASCII in debug mode */}
                    {debugMode && (
                      <Text 
                        fontSize="12px" 
                        color={tile.character ? 'black' : tile.mobId ? 'red' : tile.itemId ? 'green' : 'white'}
                        fontWeight="bold"
                        zIndex={3}
                        position="absolute"
                        top="50%"
                        left="50%"
                        transform="translate(-50%, -50%)"
                        textShadow={tile.character ? "0px 0px 2px white" : "none"}
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
      </Box>
      
      <Box mt={4}>
        <Heading size="sm" mb={2} color="blue.300">Legend</Heading>
        <SimpleGrid columns={[2, 3, 4]} spacing={3}>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#555" border="1px solid #666" mr={2} />
            <Text fontSize="sm">Wall (#)</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#111" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Text fontSize="10px">.</Text>
            </Box>
            <Text fontSize="sm">Floor (.)</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#850" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Text color="yellow.600">+</Text>
            </Box>
            <Text fontSize="sm">Door (+)</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#00A" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Text color="blue.300">↑</Text>
            </Box>
            <Text fontSize="sm">Up Stairs (&lt;)</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#A00" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Text color="red.300">↓</Text>
            </Box>
            <Text fontSize="sm">Down Stairs (&gt;)</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#111" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Box width="14px" height="14px" borderRadius="50%" bg="#FF0" />
            </Box>
            <Text fontSize="sm">Player (@)</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#111" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Box width="14px" height="14px" borderRadius="50%" bg="#F22" />
            </Box>
            <Text fontSize="sm">Monster (M)</Text>
          </Box>
          <Box display="flex" alignItems="center">
            <Box width="20px" height="20px" bg="#111" mr={2} display="flex" alignItems="center" justifyContent="center">
              <Box width="14px" height="14px" borderRadius="50%" bg="#0FF" />
            </Box>
            <Text fontSize="sm">Item (i)</Text>
          </Box>
        </SimpleGrid>
      </Box>
    </VStack>
  );
};

export default RoomSymbolDemo; 