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
    if (tile.character) return '@';
    if (tile.mobId) {
      const mob = mobs[tile.mobId];
      if (!mob) return 'M';
      
      switch (mob.type) {
        case 'rat':
          return 'r';
        case 'goblin':
          return 'g';
        case 'orc':
          return 'o';
        case 'troll':
          return 'T';
        case 'dragon':
          return 'D';
        case 'shopkeeper':
          return '$';
        case 'boss':
          return 'B';
        default:
          return 'M';
      }
    }
    
    if (tile.itemId) {
      const item = items[tile.itemId];
      if (!item) return 'i';
      
      switch (item.type) {
        case 'weapon':
          return 'w';
        case 'armor':
          return 'a';
        case 'potion':
          return 'p';
        case 'gold':
          return '$';
        case 'scroll':
          return '?';
        default:
          return 'i';
      }
    }
    
    switch (tile.type) {
      case 'wall':
        return '#';
      case 'door':
        return '+';
      case 'upStairs':
        return '<';
      case 'downStairs':
        return '>';
      case 'floor':
        return '.';
      default:
        return ' ';
    }
  };

  // Mock tile, mob, and item data for demo
  const mockTiles: Tile[][] = [
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'door', walkable: true, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
    ],
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'wall', walkable: false, explored: true },
    ],
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'wall', walkable: false, explored: true },
    ],
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'wall', walkable: false, explored: true },
    ],
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: demoType === 'entrance' ? 'downStairs' : demoType === 'safe' ? 'upStairs' : 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'wall', walkable: false, explored: true },
    ],
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'wall', walkable: false, explored: true },
    ],
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true, character: demoType === 'entrance' ? 'player1' : undefined },
      { type: 'wall', walkable: false, explored: true },
    ],
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'wall', walkable: false, explored: true },
    ],
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'floor', walkable: true, explored: true },
      { type: 'wall', walkable: false, explored: true },
    ],
    [
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
      { type: 'wall', walkable: false, explored: true },
    ],
  ];

  // Mock mob and item data based on demo type
  const mockMobs: { [key: string]: Mob } = {};
  const mockItems: { [key: string]: Item } = {};

  if (demoType === 'shop') {
    mockMobs['shopkeeper-1'] = {
      id: 'shopkeeper-1',
      type: 'shopkeeper',
      name: 'Shopkeeper',
      health: 50,
      maxHealth: 50,
      position: { x: 5, y: 5 },
    };
    mockTiles[5][5] = {
      ...mockTiles[5][5],
      mobId: 'shopkeeper-1',
    };
  } else if (demoType === 'boss') {
    mockMobs['boss-1'] = {
      id: 'boss-1',
      type: 'boss',
      name: 'Dungeon Boss',
      health: 100,
      maxHealth: 100,
      position: { x: 5, y: 5 },
    };
    mockTiles[5][5] = {
      ...mockTiles[5][5],
      mobId: 'boss-1',
    };
  } else if (demoType === 'treasure') {
    mockItems['gold-1'] = {
      id: 'gold-1',
      type: 'gold',
      name: 'Gold Pile',
      position: { x: 5, y: 5 },
    };
    mockTiles[5][5] = {
      ...mockTiles[5][5],
      itemId: 'gold-1',
    };
  } else if (demoType === 'standard') {
    mockMobs['mob-1'] = {
      id: 'mob-1',
      type: 'goblin',
      name: 'Goblin',
      health: 20,
      maxHealth: 20,
      position: { x: 5, y: 5 },
    };
    mockTiles[5][5] = {
      ...mockTiles[5][5],
      mobId: 'mob-1',
    };
  }

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
            gridTemplateColumns={`repeat(${mockTiles[0].length}, 20px)`}
            gridTemplateRows={`repeat(${mockTiles.length}, 20px)`}
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
            {mockTiles.flatMap((row, y) => 
              row.map((tile, x) => {
                const tileColor = getTileColor(tile);
                const entityColor = getEntityColor(tile, mockMobs, mockItems);
                const tileContent = getTileContent(tile, mockMobs, mockItems);
                
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
                    title={`${tile.type}${tile.mobId ? ' - ' + (mockMobs?.[tile.mobId]?.name || 'Monster') : ''}${tile.itemId ? ' - ' + (mockItems?.[tile.itemId]?.name || 'Item') : ''}`}
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