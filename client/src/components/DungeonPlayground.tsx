import React, { useState } from 'react';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  Select,
  VStack,
  HStack,
  Heading,
  Text,
  useToast,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  Badge,
  SimpleGrid
} from '@chakra-ui/react';
import { createDungeon } from '../services/api';
import { Dungeon } from '../types';
import DungeonRoomRenderer from './DungeonRoomRenderer';

interface DungeonPlaygroundProps {
  showCode?: boolean;
}

const DungeonPlayground: React.FC<DungeonPlaygroundProps> = ({ showCode = false }) => {
  // State for dungeon creation
  const [dungeonName, setDungeonName] = useState<string>('Test Dungeon');
  const [floorCount, setFloorCount] = useState<number>(3);
  const [difficulty, setDifficulty] = useState<string>('easy');
  const [isCreating, setIsCreating] = useState<boolean>(false);
  
  // State for dungeon display
  const [currentDungeon, setCurrentDungeon] = useState<Dungeon | null>(null);
  const [error, setError] = useState<string | null>(null);
  
  const toast = useToast();

  // Create a new dungeon
  const handleCreateDungeon = async () => {
    if (!dungeonName.trim()) {
      toast({
        title: 'Error',
        description: 'Please enter a dungeon name',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    try {
      setIsCreating(true);
      const newDungeon = await createDungeon({
        name: dungeonName,
        floors: floorCount,
        difficulty: difficulty,
      });
      
      setCurrentDungeon(newDungeon);
      setError(null);
      
      toast({
        title: 'Dungeon created',
        description: `${newDungeon.name} has been created successfully`,
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (err) {
      console.error('Failed to create dungeon:', err);
      setError('Failed to create dungeon. Please try again.');
      toast({
        title: 'Error',
        description: 'Failed to create dungeon. Please try again.',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    } finally {
      setIsCreating(false);
    }
  };

  return (
    <VStack spacing={6} align="stretch">
      <Heading size="md">Dungeon Playground</Heading>
      <Text>Create and explore procedurally generated dungeons with configurable parameters.</Text>
      
      <SimpleGrid columns={[1, null, 2]} spacing={4}>
        <FormControl>
          <FormLabel>Dungeon Name</FormLabel>
          <Input 
            value={dungeonName} 
            onChange={(e) => setDungeonName(e.target.value)} 
            placeholder="Enter dungeon name"
          />
        </FormControl>

        <SimpleGrid columns={2} spacing={4}>
          <FormControl>
            <FormLabel>Floors</FormLabel>
            <NumberInput 
              min={1} 
              max={10} 
              value={floorCount} 
              onChange={(_, value) => setFloorCount(value)}
            >
              <NumberInputField />
              <NumberInputStepper>
                <NumberIncrementStepper />
                <NumberDecrementStepper />
              </NumberInputStepper>
            </NumberInput>
          </FormControl>
          
          <FormControl>
            <FormLabel>Difficulty</FormLabel>
            <Select 
              value={difficulty} 
              onChange={(e) => setDifficulty(e.target.value)}
            >
              <option value="easy">Easy</option>
              <option value="medium">Medium</option>
              <option value="hard">Hard</option>
              <option value="extreme">Extreme</option>
            </Select>
          </FormControl>
        </SimpleGrid>
      </SimpleGrid>

      <Button 
        colorScheme="blue" 
        onClick={handleCreateDungeon} 
        isLoading={isCreating}
        loadingText="Creating"
        width={["100%", "auto"]}
        alignSelf={["center", "flex-start"]}
        mb={2}
      >
        Create Dungeon
      </Button>
      
      {error && (
        <Box textAlign="center" py={2} bg="red.50" color="red.500" borderRadius="md" mb={4}>
          <Text>{error}</Text>
        </Box>
      )}
      
      {currentDungeon && (
        <Box mt={2}>
          <HStack mb={4} wrap="wrap" spacing={2}>
            <Text fontWeight="bold" fontSize="lg">{currentDungeon.name}</Text>
            <Badge colorScheme="blue">{currentDungeon.difficulty}</Badge>
            <Badge colorScheme="green">{currentDungeon.floors} floors</Badge>
            <Badge colorScheme="purple">{currentDungeon.playerCount} players</Badge>
          </HStack>
          
          <Tabs 
            variant="enclosed" 
            colorScheme="blue"
            isLazy
          >
            <TabList>
              {Array.from({ length: currentDungeon.floors }, (_, i) => (
                <Tab key={i + 1}>Floor {i + 1}</Tab>
              ))}
            </TabList>
            <TabPanels>
              {Array.from({ length: currentDungeon.floors }, (_, i) => (
                <TabPanel key={i + 1}>
                  <Box width="100%" height="400px" position="relative" border="1px solid" borderColor="gray.200" borderRadius="md" overflow="hidden">
                    <DungeonRoomRenderer
                      dungeonId={currentDungeon.id}
                      floorNumber={i + 1}
                      width={20}
                      height={20}
                      debug={true}
                    />
                  </Box>
                </TabPanel>
              ))}
            </TabPanels>
          </Tabs>
        </Box>
      )}

      {!currentDungeon && (
        <Box 
          width="100%" 
          height="200px" 
          bg="gray.50" 
          borderRadius="md" 
          display="flex" 
          alignItems="center" 
          justifyContent="center"
          border="1px dashed"
          borderColor="gray.200"
        >
          <Text color="gray.500">Create a dungeon to start exploring</Text>
        </Box>
      )}
    </VStack>
  );
};

export default DungeonPlayground; 