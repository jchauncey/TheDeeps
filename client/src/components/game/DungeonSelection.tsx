import { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Flex,
  Heading,
  Text,
  VStack,
  HStack,
  Input,
  FormControl,
  FormLabel,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  useToast,
  Spinner,
  Divider,
  Badge,
  Card,
  CardBody,
  CardHeader,
  SimpleGrid,
} from '@chakra-ui/react';
import { createDungeon, getAvailableDungeons, joinDungeon, sendWebSocketMessage, joinDungeonWS } from '../../services/api';
import { DungeonData, FloorData } from '../../types/game';

interface DungeonSelectionProps {
  characterId: string;
  onDungeonSelected: (dungeonId: string, floorData: FloorData) => void;
  onBack: () => void;
}

export const DungeonSelection = ({ characterId, onDungeonSelected, onBack }: DungeonSelectionProps) => {
  const [dungeonName, setDungeonName] = useState('');
  const [numFloors, setNumFloors] = useState(5);
  const [availableDungeons, setAvailableDungeons] = useState<DungeonData[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [isJoining, setIsJoining] = useState(false);
  const [selectedDungeonId, setSelectedDungeonId] = useState<string | null>(null);
  const toast = useToast();

  // Load available dungeons on component mount
  useEffect(() => {
    console.log('DungeonSelection component mounted, loading dungeons...');
    console.log('Character ID:', characterId);
    loadDungeons();
  }, [characterId]);

  const loadDungeons = async () => {
    console.log('Loading dungeons...');
    setIsLoading(true);
    try {
      const result = await getAvailableDungeons();
      console.log('Available dungeons result:', result);
      if (result.success && result.dungeons) {
        console.log('Setting available dungeons:', result.dungeons);
        setAvailableDungeons(result.dungeons);
      } else {
        console.error('Failed to load dungeons:', result.message);
        toast({
          title: 'Error',
          description: result.message || 'Failed to load dungeons',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    } catch (error) {
      console.error('Error loading dungeons:', error);
      toast({
        title: 'Error',
        description: 'An unexpected error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateDungeon = async () => {
    if (!dungeonName.trim()) {
      toast({
        title: 'Error',
        description: 'Please enter a dungeon name',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
      return;
    }

    setIsCreating(true);
    try {
      // Use the WebSocket API instead of HTTP API
      const success = sendWebSocketMessage({
        type: 'create_dungeon',
        name: dungeonName,
        numFloors: numFloors
      });
      
      if (success) {
        toast({
          title: 'Creating Dungeon',
          description: 'Dungeon creation request sent',
          status: 'info',
          duration: 3000,
          isClosable: true,
        });
        
        // Wait for the dungeon to be created
        // The server will send a dungeon_created message
        // which will be handled by the WebSocket message handler
        
        // Refresh the dungeon list after a short delay
        setTimeout(() => {
          loadDungeons();
        }, 1000);
      } else {
        toast({
          title: 'Error',
          description: 'Failed to send dungeon creation request',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    } catch (error) {
      toast({
        title: 'Error',
        description: 'An unexpected error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsCreating(false);
    }
  };

  const handleJoinDungeon = async (dungeonId: string) => {
    setIsJoining(true);
    try {
      // Use the WebSocket API instead of HTTP API
      const success = sendWebSocketMessage({
        type: 'join_dungeon',
        dungeonId: dungeonId,
        characterId: characterId
      });
      
      if (success) {
        toast({
          title: 'Joining Dungeon',
          description: 'Dungeon join request sent',
          status: 'info',
          duration: 3000,
          isClosable: true,
        });
        
        // Wait for the dungeon to be joined
        // The server will send a dungeon_joined message
        // which will be handled by the WebSocket message handler
        
        // For now, we'll simulate a successful join after a short delay
        // In a real implementation, we would wait for the server's response
        setTimeout(() => {
          // Create a simple floor data object for demonstration
          const dummyFloorData: FloorData = {
            type: 'floor_data',
            floor: {
              width: 20,
              height: 20,
              tiles: [],
              entities: [],
              items: []
            },
            playerPosition: { x: 10, y: 10 },
            currentFloor: 0,
            playerData: {
              id: characterId,
              name: 'Player',
              characterClass: 'warrior',
              health: 100,
              maxHealth: 100
            }
          };
          
          // Notify parent component that a dungeon has been selected
          onDungeonSelected(dungeonId, dummyFloorData);
        }, 1000);
      } else {
        toast({
          title: 'Error',
          description: 'Failed to send dungeon join request',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    } catch (error) {
      toast({
        title: 'Error',
        description: 'An unexpected error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsJoining(false);
    }
  };

  const handleSelectDungeon = (dungeonId: string) => {
    setSelectedDungeonId(dungeonId);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  return (
    <Box
      w="100%"
      h="100%"
      p={6}
      bg="gray.800"
      color="white"
      overflow="auto"
    >
      <Heading mb={6}>Dungeon Selection</Heading>
      
      <Flex direction={{ base: 'column', md: 'row' }} gap={8}>
        {/* Create Dungeon Section */}
        <Box flex="1" p={4} bg="gray.700" borderRadius="md">
          <Heading size="md" mb={4}>Create New Dungeon</Heading>
          <VStack spacing={4} align="stretch">
            <FormControl>
              <FormLabel>Dungeon Name</FormLabel>
              <Input 
                value={dungeonName}
                onChange={(e) => setDungeonName(e.target.value)}
                placeholder="Enter dungeon name"
              />
            </FormControl>
            
            <FormControl>
              <FormLabel>Number of Floors</FormLabel>
              <NumberInput 
                min={1} 
                max={20} 
                value={numFloors}
                onChange={(_, value) => setNumFloors(value)}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
            </FormControl>
            
            <Button 
              colorScheme="purple" 
              onClick={handleCreateDungeon}
              isLoading={isCreating}
              loadingText="Creating..."
            >
              Create Dungeon
            </Button>
          </VStack>
        </Box>
        
        {/* Available Dungeons Section */}
        <Box flex="2" p={4} bg="gray.700" borderRadius="md">
          <Flex justify="space-between" align="center" mb={4}>
            <Heading size="md">Available Dungeons</Heading>
            <Button 
              size="sm" 
              onClick={loadDungeons} 
              isLoading={isLoading}
              loadingText="Refreshing..."
            >
              Refresh
            </Button>
          </Flex>
          
          {isLoading ? (
            <Flex justify="center" align="center" h="200px">
              <Spinner size="xl" />
            </Flex>
          ) : availableDungeons.length === 0 ? (
            <Text textAlign="center" py={8}>No dungeons available. Create one to get started!</Text>
          ) : (
            <VStack spacing={4} align="stretch" maxH="400px" overflowY="auto">
              <SimpleGrid columns={{ base: 1, lg: 2 }} spacing={4}>
                {availableDungeons.map((dungeon) => (
                  <Card 
                    key={dungeon.id} 
                    bg={selectedDungeonId === dungeon.id ? "purple.700" : "gray.600"}
                    cursor="pointer"
                    onClick={() => handleSelectDungeon(dungeon.id)}
                    _hover={{ bg: "purple.600" }}
                    transition="all 0.2s"
                  >
                    <CardHeader pb={2}>
                      <Flex justify="space-between" align="center">
                        <Heading size="sm">{dungeon.name}</Heading>
                        <Badge colorScheme={dungeon.playerCount > 0 ? "green" : "gray"}>
                          {dungeon.playerCount} {dungeon.playerCount === 1 ? 'player' : 'players'}
                        </Badge>
                      </Flex>
                    </CardHeader>
                    <CardBody pt={0}>
                      <Text fontSize="sm">Floors: {dungeon.numFloors}</Text>
                      <Text fontSize="sm">Created: {formatDate(dungeon.createdAt)}</Text>
                    </CardBody>
                  </Card>
                ))}
              </SimpleGrid>
            </VStack>
          )}
          
          <Flex justify="space-between" mt={6}>
            <Button onClick={onBack} variant="outline">
              Back
            </Button>
            <Button 
              colorScheme="green" 
              isDisabled={!selectedDungeonId}
              onClick={() => selectedDungeonId && handleJoinDungeon(selectedDungeonId)}
              isLoading={isJoining}
              loadingText="Joining..."
            >
              Join Selected Dungeon
            </Button>
          </Flex>
        </Box>
      </Flex>
    </Box>
  );
}; 