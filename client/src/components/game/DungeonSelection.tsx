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
import { createDungeon, getAvailableDungeons, joinDungeon, sendWebSocketMessage } from '../../services/api';
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
    
    if (!characterId) {
      console.error('Character ID is missing in DungeonSelection component');
      toast({
        title: 'Error',
        description: 'Character ID is missing. Please create a character first.',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
      onBack(); // Go back to character creation
      return;
    }
    
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
      // Use the REST API to create a dungeon
      const result = await createDungeon(dungeonName, numFloors);
      
      if (result.success) {
        toast({
          title: 'Success',
          description: 'Dungeon created successfully!',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
        
        // Refresh the dungeon list
        loadDungeons();
      } else {
        toast({
          title: 'Error',
          description: result.message || 'Failed to create dungeon',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    } catch (error) {
      console.error('Error creating dungeon:', error);
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
      // Set the selected dungeon ID
      setSelectedDungeonId(dungeonId);
      
      // Log character ID for debugging
      console.log('Joining dungeon with characterId:', characterId);
      
      if (!characterId) {
        toast({
          title: 'Error',
          description: 'Character ID is missing. Please create a character first.',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
        setIsJoining(false);
        return;
      }
      
      // Use the REST API to join the dungeon
      const result = await joinDungeon(dungeonId, characterId);
      
      if (result.success && result.floorData) {
        toast({
          title: 'Success',
          description: 'Joined dungeon successfully!',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
        
        // Call the parent component's callback with the dungeon ID and floor data
        onDungeonSelected(dungeonId, result.floorData);
      } else {
        toast({
          title: 'Error',
          description: result.message || 'Failed to join dungeon',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    } catch (error) {
      console.error('Error joining dungeon:', error);
      toast({
        title: 'Error',
        description: error instanceof Error ? error.message : 'An unexpected error occurred',
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
      position="fixed"
      top={0}
      left={0}
      right={0}
      bottom={0}
      bg="#291326"
      display="flex"
      alignItems="center"
      justifyContent="center"
      p={6}
    >
      <Flex 
        direction="column" 
        bg="rgba(0, 0, 0, 0.7)" 
        p={6} 
        borderRadius="md" 
        width="100%" 
        maxW="1000px"
        height="80vh"
        overflow="auto"
        color="white"
      >
        <Heading mb={6}>Dungeon Selection</Heading>
        
        <Flex direction={{ base: 'column', md: 'row' }} gap={8} flex="1">
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
      </Flex>
    </Box>
  );
}; 