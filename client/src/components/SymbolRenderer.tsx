import React from 'react';
import {
  Box,
  SimpleGrid,
  Text,
  Heading,
  VStack,
  HStack,
  Divider,
} from '@chakra-ui/react';

interface SymbolCategory {
  name: string;
  symbols: {
    name: string;
    symbol: string;
    color: string;
    bgColor: string;
    description: string;
  }[];
}

const SymbolRenderer: React.FC = () => {
  // Define all the symbols used in the game
  const symbolCategories: SymbolCategory[] = [
    {
      name: 'Tile Types',
      symbols: [
        { name: 'Wall', symbol: '#', color: 'white', bgColor: '#555', description: 'Impassable barrier' },
        { name: 'Floor', symbol: '.', color: 'white', bgColor: '#111', description: 'Walkable surface' },
        { name: 'Corridor', symbol: '·', color: 'white', bgColor: '#222', description: 'Connects rooms' },
        { name: 'Door', symbol: '+', color: 'yellow.600', bgColor: '#850', description: 'Room entrance/exit' },
        { name: 'Up Stairs', symbol: '↑', color: 'blue.300', bgColor: '#00A', description: 'Go up a level' },
        { name: 'Down Stairs', symbol: '↓', color: 'red.300', bgColor: '#A00', description: 'Go down a level' },
        { name: 'Unexplored', symbol: ' ', color: 'white', bgColor: '#000', description: 'Not yet seen' },
      ],
    },
    {
      name: 'Characters',
      symbols: [
        { name: 'Player', symbol: '@', color: 'yellow', bgColor: '#111', description: 'Player character' },
      ],
    },
    {
      name: 'Monsters',
      symbols: [
        { name: 'Dragon', symbol: 'D', color: 'red', bgColor: '#111', description: 'Powerful fire-breathing monster' },
        { name: 'Ogre', symbol: 'O', color: 'orange', bgColor: '#111', description: 'Large, strong monster' },
        { name: 'Boss', symbol: 'B', color: 'red', bgColor: '#111', description: 'Powerful unique monster' },
        { name: 'Goblin', symbol: 'g', color: 'green', bgColor: '#111', description: 'Small, weak monster' },
        { name: 'Skeleton', symbol: 's', color: 'white', bgColor: '#111', description: 'Undead monster' },
        { name: 'Shopkeeper', symbol: 'S', color: 'red', bgColor: '#111', description: 'Sells items' },
        { name: 'Generic Monster', symbol: 'M', color: 'red', bgColor: '#111', description: 'Unspecified monster' },
      ],
    },
    {
      name: 'Items',
      symbols: [
        { name: 'Gold', symbol: '$', color: 'yellow', bgColor: '#111', description: 'Currency' },
        { name: 'Potion', symbol: '!', color: 'magenta', bgColor: '#111', description: 'Consumable item' },
        { name: 'Weapon', symbol: '/', color: 'cyan', bgColor: '#111', description: 'Increases attack' },
        { name: 'Armor', symbol: '[', color: 'blue', bgColor: '#111', description: 'Increases defense' },
        { name: 'Generic Item', symbol: 'i', color: 'green', bgColor: '#111', description: 'Unspecified item' },
      ],
    },
  ];

  // Render a single symbol
  const renderSymbol = (symbol: SymbolCategory['symbols'][0]) => (
    <Box 
      key={symbol.name} 
      borderWidth="1px" 
      borderColor="gray.700" 
      borderRadius="md" 
      overflow="hidden"
    >
      <HStack spacing={0}>
        {/* Symbol display */}
        <Box 
          bg={symbol.bgColor} 
          width="40px" 
          height="40px" 
          display="flex" 
          alignItems="center" 
          justifyContent="center"
          fontFamily="monospace"
          fontSize="xl"
          fontWeight="bold"
          color={symbol.color}
        >
          {symbol.symbol}
        </Box>
        
        {/* Symbol info */}
        <Box p={2} bg="gray.800" flex="1">
          <Text fontWeight="bold" fontSize="sm">{symbol.name}</Text>
          <Text fontSize="xs" color="gray.400">{symbol.description}</Text>
        </Box>
      </HStack>
    </Box>
  );

  // Render a category of symbols
  const renderCategory = (category: SymbolCategory) => (
    <Box key={category.name} mb={6}>
      <Heading size="sm" mb={3} color="blue.300">{category.name}</Heading>
      <SimpleGrid columns={[1, 2, 3]} spacing={3}>
        {category.symbols.map(renderSymbol)}
      </SimpleGrid>
    </Box>
  );

  // Render a visual example of entity representation
  const renderEntityExample = () => (
    <Box mb={6}>
      <Heading size="sm" mb={3} color="blue.300">Entity Representation</Heading>
      <SimpleGrid columns={[1, 2, 3]} spacing={3}>
        {/* Player */}
        <Box borderWidth="1px" borderColor="gray.700" borderRadius="md" overflow="hidden">
          <HStack spacing={0}>
            <Box 
              bg="#111" 
              width="40px" 
              height="40px" 
              display="flex" 
              alignItems="center" 
              justifyContent="center"
              position="relative"
            >
              <Box 
                position="absolute"
                width="24px"
                height="24px"
                borderRadius="50%"
                bg="#FF0"
              />
            </Box>
            <Box p={2} bg="gray.800" flex="1">
              <Text fontWeight="bold" fontSize="sm">Player</Text>
              <Text fontSize="xs" color="gray.400">Yellow circle</Text>
            </Box>
          </HStack>
        </Box>
        
        {/* Monster */}
        <Box borderWidth="1px" borderColor="gray.700" borderRadius="md" overflow="hidden">
          <HStack spacing={0}>
            <Box 
              bg="#111" 
              width="40px" 
              height="40px" 
              display="flex" 
              alignItems="center" 
              justifyContent="center"
              position="relative"
            >
              <Box 
                position="absolute"
                width="24px"
                height="24px"
                borderRadius="50%"
                bg="#F22"
              />
            </Box>
            <Box p={2} bg="gray.800" flex="1">
              <Text fontWeight="bold" fontSize="sm">Monster</Text>
              <Text fontSize="xs" color="gray.400">Red circle</Text>
            </Box>
          </HStack>
        </Box>
        
        {/* Item */}
        <Box borderWidth="1px" borderColor="gray.700" borderRadius="md" overflow="hidden">
          <HStack spacing={0}>
            <Box 
              bg="#111" 
              width="40px" 
              height="40px" 
              display="flex" 
              alignItems="center" 
              justifyContent="center"
              position="relative"
            >
              <Box 
                position="absolute"
                width="24px"
                height="24px"
                borderRadius="50%"
                bg="#0FF"
              />
            </Box>
            <Box p={2} bg="gray.800" flex="1">
              <Text fontWeight="bold" fontSize="sm">Item</Text>
              <Text fontSize="xs" color="gray.400">Colored circle</Text>
            </Box>
          </HStack>
        </Box>
      </SimpleGrid>
    </Box>
  );

  return (
    <VStack spacing={4} align="stretch">
      <Text>
        This component displays all the symbols used in the map rendering system. 
        In normal mode, entities are shown as colored circles, while in debug mode, 
        ASCII symbols are displayed.
      </Text>
      
      <Divider />
      
      {renderEntityExample()}
      
      {symbolCategories.map(renderCategory)}
    </VStack>
  );
};

export default SymbolRenderer; 