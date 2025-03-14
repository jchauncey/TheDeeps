import React, { useState, useRef, useEffect } from 'react';
import { 
  Box, 
  Container, 
  Heading, 
  Select, 
  VStack, 
  Divider, 
  Text,
  FormControl,
  FormLabel,
  Switch,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  Code,
  Button,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  SimpleGrid,
  HStack,
  Tooltip
} from '@chakra-ui/react';
import { useNavigate, useLocation } from 'react-router-dom';
import CharacterCard from '../components/CharacterCard';
import RoomRenderer from '../components/RoomRenderer';
import SymbolRenderer from '../components/SymbolRenderer';
import RoomSymbolDemo from '../components/RoomSymbolDemo';
import MovementDemo from '../components/MovementDemo';
import { ArrowBackIcon, RepeatIcon } from '@chakra-ui/icons';
import { Character, CharacterClass } from '../types';

// Component options
const components = [
  { value: 'CharacterCard', label: 'Character Card' },
  { value: 'RoomRenderer', label: 'Room Renderer' },
  { value: 'SymbolRenderer', label: 'Map Symbols' },
  { value: 'RoomSymbolDemo', label: 'Room Symbol Demo' },
  { value: 'MovementDemo', label: 'Movement Demo' },
  // Add more components as needed
];

// Create mock character data directly in this file
const mockCharacters: Character[] = [
  {
    id: '1',
    name: 'Test Warrior',
    class: 'warrior' as CharacterClass,
    level: 1,
    experience: 0,
    attributes: {
      strength: 16,
      dexterity: 12,
      constitution: 14,
      intelligence: 8,
      wisdom: 10,
      charisma: 10
    },
    skills: {},
    maxHp: 20,
    currentHp: 20,
    maxMana: 10,
    currentMana: 10,
    gold: 100,
    currentFloor: 1,
    position: { x: 0, y: 0 },
    inventory: [],
    equipment: {}
  },
  {
    id: '2',
    name: 'Test Mage',
    class: 'mage' as CharacterClass,
    level: 1,
    experience: 0,
    attributes: {
      strength: 8,
      dexterity: 12,
      constitution: 10,
      intelligence: 16,
      wisdom: 14,
      charisma: 10
    },
    skills: {},
    maxHp: 12,
    currentHp: 12,
    maxMana: 20,
    currentMana: 20,
    gold: 100,
    currentFloor: 1,
    position: { x: 0, y: 0 },
    inventory: [],
    equipment: {}
  }
];

// Mock functions for component props
const handleDelete = (id: string) => {
  alert(`Delete character with ID: ${id}`);
};

const handleSelect = (character: Character) => {
  alert(`Selected character: ${character.name}`);
};

const ComponentPlayground: React.FC = () => {
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const componentParam = queryParams.get('component');
  
  const [selectedComponent, setSelectedComponent] = useState<string>(
    componentParam && components.some(c => c.value === componentParam) 
      ? componentParam 
      : 'CharacterCard'
  );
  const [roomType, setRoomType] = useState<string>('entrance');
  const [showCode, setShowCode] = useState<boolean>(false);
  const [roomWidth, setRoomWidth] = useState<number>(20);
  const [roomHeight, setRoomHeight] = useState<number>(20);
  const [roomDimWidth, setRoomDimWidth] = useState<number>(8);
  const [roomDimHeight, setRoomDimHeight] = useState<number>(8);
  const [refreshKey, setRefreshKey] = useState<number>(0);
  const [isRoomLoading, setIsRoomLoading] = useState<boolean>(false);
  const [debugMode, setDebugMode] = useState<boolean>(true);
  const [gridSize, setGridSize] = useState<number>(15);
  const [wallDensity, setWallDensity] = useState<number>(0.2);
  const navigate = useNavigate();

  // Room type options
  const roomTypes = [
    { value: 'entrance', label: 'Entrance Room' },
    { value: 'standard', label: 'Standard Room' },
    { value: 'treasure', label: 'Treasure Room' },
    { value: 'boss', label: 'Boss Room' },
    { value: 'safe', label: 'Safe Room' },
    { value: 'shop', label: 'Shop Room' },
  ];

  // Function to refresh the room data
  const handleRefreshRoom = () => {
    setIsRoomLoading(true);
    setRefreshKey(prevKey => prevKey + 1);
  };

  // Reset loading state when room type or dimensions change
  useEffect(() => {
    setIsRoomLoading(true);
  }, [roomType, roomWidth, roomHeight, roomDimWidth, roomDimHeight]);

  // Render the selected component
  const renderComponent = () => {
    switch (selectedComponent) {
      case 'CharacterCard':
        return (
          <VStack spacing={6}>
            <Heading size="md">Character Cards</Heading>
            <Text>These cards display character information and provide actions.</Text>
            
            <Tabs variant="enclosed" width="100%">
              <TabList>
                {mockCharacters.map(character => (
                  <Tab key={character.id}>{character.name}</Tab>
                ))}
              </TabList>
              <TabPanels>
                {mockCharacters.map(character => (
                  <TabPanel key={character.id}>
                    <Box maxW="sm" mx="auto">
                      <CharacterCard 
                        character={character} 
                        onDelete={handleDelete} 
                        onSelect={handleSelect} 
                      />
                    </Box>
                    
                    {showCode && (
                      <Box mt={6} p={4} bg="gray.700" borderRadius="md" overflowX="auto">
                        <Code colorScheme="gray" whiteSpace="pre">
{`<CharacterCard 
  character={{
    id: "${character.id}",
    name: "${character.name}",
    class: "${character.class}",
    level: ${character.level},
    // ... other properties
  }} 
  onDelete={handleDelete} 
  onSelect={handleSelect} 
/>`}
                        </Code>
                      </Box>
                    )}
                  </TabPanel>
                ))}
              </TabPanels>
            </Tabs>
          </VStack>
        );
      
      case 'RoomRenderer':
        return (
          <VStack spacing={6}>
            <Heading size="md">Room Renderer</Heading>
            <Text>This component renders different types of dungeon rooms.</Text>
            
            <HStack width="100%" justifyContent="space-between">
              <FormControl display="flex" alignItems="center">
                <FormLabel htmlFor="room-type" mb={0}>Room Type:</FormLabel>
                <Select 
                  id="room-type"
                  value={roomType}
                  onChange={(e) => setRoomType(e.target.value)}
                  width="auto"
                  ml={2}
                >
                  {roomTypes.map(type => (
                    <option key={type.value} value={type.value}>{type.label}</option>
                  ))}
                </Select>
              </FormControl>
              
              <HStack>
                <FormControl display="flex" alignItems="center">
                  <FormLabel htmlFor="debug-mode" mb={0} fontSize="sm">Show Symbols:</FormLabel>
                  <Switch 
                    id="debug-mode" 
                    isChecked={debugMode} 
                    onChange={() => setDebugMode(!debugMode)} 
                    colorScheme="red"
                    size="sm"
                  />
                </FormControl>
                
                <Tooltip label="Refresh room data">
                  <Button 
                    leftIcon={<RepeatIcon />} 
                    onClick={handleRefreshRoom}
                    size="sm"
                    colorScheme="blue"
                    isLoading={isRoomLoading}
                    loadingText="Loading"
                  >
                    Refresh
                  </Button>
                </Tooltip>
              </HStack>
            </HStack>
            
            <SimpleGrid columns={2} spacing={4} width="100%">
              <FormControl>
                <FormLabel htmlFor="room-width">Map Width:</FormLabel>
                <NumberInput 
                  id="room-width" 
                  value={roomWidth} 
                  onChange={(_, val) => setRoomWidth(val)}
                  min={10} 
                  max={50}
                >
                  <NumberInputField />
                  <NumberInputStepper>
                    <NumberIncrementStepper />
                    <NumberDecrementStepper />
                  </NumberInputStepper>
                </NumberInput>
              </FormControl>
              
              <FormControl>
                <FormLabel htmlFor="room-height">Map Height:</FormLabel>
                <NumberInput 
                  id="room-height" 
                  value={roomHeight} 
                  onChange={(_, val) => setRoomHeight(val)}
                  min={10} 
                  max={50}
                >
                  <NumberInputField />
                  <NumberInputStepper>
                    <NumberIncrementStepper />
                    <NumberDecrementStepper />
                  </NumberInputStepper>
                </NumberInput>
              </FormControl>
              
              <FormControl>
                <FormLabel htmlFor="room-dim-width">Room Width:</FormLabel>
                <NumberInput 
                  id="room-dim-width" 
                  value={roomDimWidth} 
                  onChange={(_, val) => setRoomDimWidth(val)}
                  min={5} 
                  max={20}
                >
                  <NumberInputField />
                  <NumberInputStepper>
                    <NumberIncrementStepper />
                    <NumberDecrementStepper />
                  </NumberInputStepper>
                </NumberInput>
              </FormControl>
              
              <FormControl>
                <FormLabel htmlFor="room-dim-height">Room Height:</FormLabel>
                <NumberInput 
                  id="room-dim-height" 
                  value={roomDimHeight} 
                  onChange={(_, val) => setRoomDimHeight(val)}
                  min={5} 
                  max={20}
                >
                  <NumberInputField />
                  <NumberInputStepper>
                    <NumberIncrementStepper />
                    <NumberDecrementStepper />
                  </NumberInputStepper>
                </NumberInput>
              </FormControl>
            </SimpleGrid>
            
            <Box width="100%" p={4} bg="gray.700" borderRadius="md">
              <Box 
                bg="black" 
                color="white" 
                fontFamily="monospace" 
                p={4} 
                borderRadius="md"
                overflow="auto"
                fontSize="16px"
                lineHeight="1"
                position="relative"
              >
                <RoomRenderer 
                  key={refreshKey}
                  roomType={roomType} 
                  width={roomWidth} 
                  height={roomHeight} 
                  roomWidth={roomDimWidth}
                  roomHeight={roomDimHeight}
                  onLoad={() => setIsRoomLoading(false)}
                  onError={() => setIsRoomLoading(false)}
                  debug={debugMode}
                />
              </Box>
              <Text color="yellow.300" mt={4} fontSize="sm">
                Note: This component is fetching real room data from the server. If you see an error, the server endpoint might not be available.
              </Text>
            </Box>
            
            {showCode && (
              <Box mt={6} p={4} bg="gray.700" borderRadius="md" overflowX="auto">
                <Code colorScheme="gray" whiteSpace="pre">
{`<RoomRenderer 
  roomType="${roomType}" 
  width={${roomWidth}} 
  height={${roomHeight}} 
  roomWidth={${roomDimWidth}}
  roomHeight={${roomDimHeight}}
/>`}
                </Code>
              </Box>
            )}
          </VStack>
        );
      
      case 'SymbolRenderer':
        return (
          <VStack spacing={6}>
            <Heading size="md">Map Symbols</Heading>
            <Text>This component displays all the symbols used in the map rendering system.</Text>
            
            <Box width="100%" p={4} bg="gray.700" borderRadius="md">
              <SymbolRenderer />
            </Box>
            
            {showCode && (
              <Box mt={6} p={4} bg="gray.700" borderRadius="md" overflowX="auto">
                <Code colorScheme="gray" whiteSpace="pre">
{`<SymbolRenderer />`}
                </Code>
              </Box>
            )}
          </VStack>
        );
      
      case 'RoomSymbolDemo':
        return (
          <VStack spacing={6}>
            <Heading size="md">Room Symbol Demo</Heading>
            <Text>This component demonstrates the symbols in the context of different room types.</Text>
            
            <Box width="100%" p={4} bg="gray.700" borderRadius="md">
              <RoomSymbolDemo />
            </Box>
            
            {showCode && (
              <Box mt={6} p={4} bg="gray.700" borderRadius="md" overflowX="auto">
                <Code colorScheme="gray" whiteSpace="pre">
{`<RoomSymbolDemo />`}
                </Code>
              </Box>
            )}
          </VStack>
        );
      
      case 'MovementDemo':
        return (
          <VStack spacing={6}>
            <Heading size="md">Movement Demo</Heading>
            <Text>This component demonstrates character movement in a grid-based environment.</Text>
            
            <SimpleGrid columns={2} spacing={4} width="100%">
              <FormControl>
                <FormLabel htmlFor="grid-size">Grid Size:</FormLabel>
                <NumberInput 
                  id="grid-size" 
                  value={gridSize} 
                  onChange={(_, val) => setGridSize(val)}
                  min={10} 
                  max={30}
                >
                  <NumberInputField />
                  <NumberInputStepper>
                    <NumberIncrementStepper />
                    <NumberDecrementStepper />
                  </NumberInputStepper>
                </NumberInput>
              </FormControl>
              
              <FormControl>
                <FormLabel htmlFor="wall-density">Wall Density:</FormLabel>
                <NumberInput 
                  id="wall-density" 
                  value={wallDensity} 
                  onChange={(_, val) => setWallDensity(val)}
                  min={0.1} 
                  max={0.4}
                  step={0.05}
                  precision={2}
                >
                  <NumberInputField />
                  <NumberInputStepper>
                    <NumberIncrementStepper />
                    <NumberDecrementStepper />
                  </NumberInputStepper>
                </NumberInput>
              </FormControl>
            </SimpleGrid>
            
            <Box width="100%" p={4} bg="gray.700" borderRadius="md">
              <MovementDemo 
                key={`${gridSize}-${wallDensity}`}
                gridSize={gridSize}
                wallDensity={wallDensity}
                debug={debugMode}
              />
            </Box>
            
            {showCode && (
              <Box mt={6} p={4} bg="gray.700" borderRadius="md" overflowX="auto">
                <Code colorScheme="gray" whiteSpace="pre">
{`<MovementDemo 
  gridSize={${gridSize}}
  wallDensity={${wallDensity}}
  debug={${debugMode}}
/>`}
                </Code>
              </Box>
            )}
          </VStack>
        );
      
      default:
        return <Text>Select a component to view</Text>;
    }
  };

  return (
    <Container maxW="container.xl" py={8}>
      <Button 
        leftIcon={<ArrowBackIcon />} 
        variant="solid" 
        mb={8} 
        onClick={() => navigate('/')}
        bg="gray.600"
        color="cyan.300"
        borderColor="cyan.500"
        borderWidth="1px"
        _hover={{ bg: "gray.700", color: "cyan.200" }}
      >
        Back to Home
      </Button>

      <Heading as="h1" size="xl" mb={6} textAlign="center" color="blue.300">
        Component Playground
      </Heading>
      <Text textAlign="center" mb={8} color="gray.400">
        View and test individual components in isolation
      </Text>

      <Box mb={8}>
        <FormControl display="flex" alignItems="center" mb={4}>
          <FormLabel htmlFor="component-select" mb={0}>Select Component:</FormLabel>
          <Select 
            id="component-select"
            value={selectedComponent}
            onChange={(e) => setSelectedComponent(e.target.value)}
            width="auto"
            ml={2}
          >
            {components.map(comp => (
              <option key={comp.value} value={comp.value}>{comp.label}</option>
            ))}
          </Select>
        </FormControl>

        <FormControl display="flex" alignItems="center">
          <FormLabel htmlFor="show-code" mb={0}>Show Code:</FormLabel>
          <Switch 
            id="show-code" 
            isChecked={showCode} 
            onChange={() => setShowCode(!showCode)} 
            colorScheme="blue"
          />
        </FormControl>
      </Box>

      <Divider mb={8} />

      <Box>
        {renderComponent()}
      </Box>
    </Container>
  );
};

export default ComponentPlayground; 