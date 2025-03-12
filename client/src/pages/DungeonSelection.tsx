import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Container,
  Flex,
  Heading,
  Text,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  FormControl,
  FormLabel,
  Input,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  useToast,
  Spinner,
  Center,
  VStack,
  HStack,
  useDisclosure,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  RadioGroup,
  Radio,
  Stack,
} from '@chakra-ui/react';
import { AddIcon, ArrowBackIcon, ArrowForwardIcon } from '@chakra-ui/icons';
import { getDungeons, createDungeon, joinDungeon, getCharacter } from '../services/api';
import { Character, Dungeon } from '../types';

interface LocationState {
  character: Character;
}

const DungeonSelection: React.FC = () => {
  const [dungeons, setDungeons] = useState<Dungeon[]>([]);
  const [selectedDungeon, setSelectedDungeon] = useState<Dungeon | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [newDungeonName, setNewDungeonName] = useState<string>('');
  const [newDungeonFloorCount, setNewDungeonFloorCount] = useState<number>(5);
  const [creatingDungeon, setCreatingDungeon] = useState<boolean>(false);
  const [joiningDungeon, setJoiningDungeon] = useState<boolean>(false);
  const [difficulty, setDifficulty] = useState<string>('easy');
  
  const { isOpen, onOpen, onClose } = useDisclosure();
  const navigate = useNavigate();
  const location = useLocation();
  const toast = useToast();
  
  const { character } = (location.state as LocationState) || {};

  useEffect(() => {
    if (!character) {
      navigate('/');
      toast({
        title: 'No character selected',
        description: 'Please select a character first',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
      return;
    }
    
    // Refresh character data to ensure it exists
    const refreshCharacterData = async () => {
      try {
        await getCharacter(character.id);
        fetchDungeons();
      } catch (err) {
        console.error('Failed to verify character:', err);
        toast({
          title: 'Character not found',
          description: 'The selected character no longer exists. Please select another character.',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
        navigate('/');
      }
    };
    
    refreshCharacterData();
  }, [character, navigate, toast]);

  const fetchDungeons = async () => {
    try {
      setLoading(true);
      const data = await getDungeons();
      setDungeons(data);
      setError(null);
    } catch (err) {
      console.error('Failed to fetch dungeons:', err);
      setError('Failed to load dungeons. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateDungeon = async () => {
    if (!newDungeonName.trim()) {
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
      setCreatingDungeon(true);
      const newDungeon = await createDungeon({
        name: newDungeonName,
        floors: newDungeonFloorCount,
        difficulty: difficulty,
      });
      
      setDungeons([...dungeons, newDungeon]);
      setSelectedDungeon(newDungeon);
      setNewDungeonName('');
      setNewDungeonFloorCount(5);
      onClose();
      
      toast({
        title: 'Dungeon created',
        description: `${newDungeon.name} has been created successfully`,
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (err) {
      console.error('Failed to create dungeon:', err);
      toast({
        title: 'Error',
        description: 'Failed to create dungeon. Please try again.',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    } finally {
      setCreatingDungeon(false);
    }
  };

  const handleSelectDungeon = (dungeon: Dungeon) => {
    setSelectedDungeon(dungeon);
  };

  const handleJoinDungeon = async () => {
    if (!selectedDungeon) {
      toast({
        title: 'No dungeon selected',
        description: 'Please select a dungeon first',
        status: 'warning',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    try {
      setJoiningDungeon(true);
      
      // First verify the character exists
      try {
        await getCharacter(character.id);
      } catch (characterErr) {
        console.error('Character verification failed:', characterErr);
        toast({
          title: 'Character not found',
          description: 'The selected character no longer exists. Please select another character.',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
        setJoiningDungeon(false);
        navigate('/'); // Return to character selection
        return;
      }
      
      // Now join the dungeon
      await joinDungeon(character.id, selectedDungeon.id);
      
      // Navigate to game screen with character and dungeon info
      navigate('/game', { 
        state: { 
          character,
          dungeon: selectedDungeon
        } 
      });
      
      toast({
        title: 'Joined dungeon',
        description: `You have entered ${selectedDungeon.name}`,
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (err) {
      console.error('Failed to join dungeon:', err);
      setJoiningDungeon(false);
      toast({
        title: 'Error',
        description: 'Failed to join dungeon. Please try again.',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  const handleBack = () => {
    navigate('/');
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  return (
    <Container maxW="container.xl" py={8}>
      <Flex direction="column" align="center" mb={8}>
        <Heading as="h1" size="2xl" mb={2}>
          The Deeps
        </Heading>
        <Text fontSize="xl" color="gray.400">
          Select a dungeon to explore
        </Text>
      </Flex>

      {loading ? (
        <Center h="300px">
          <Spinner size="xl" thickness="4px" speed="0.65s" color="blue.500" />
        </Center>
      ) : error ? (
        <Box textAlign="center" p={5} color="red.500">
          <Text>{error}</Text>
          <Button mt={4} onClick={fetchDungeons}>
            Retry
          </Button>
        </Box>
      ) : (
        <>
          <Box mb={6}>
            <Heading as="h2" size="lg" mb={4}>
              Available Dungeons
            </Heading>
            
            {dungeons.length === 0 ? (
              <Box textAlign="center" p={10} bg="gray.800" borderRadius="lg">
                <Text fontSize="lg" mb={4}>
                  No dungeons available. Create a new one!
                </Text>
              </Box>
            ) : (
              <Box overflowX="auto">
                <Table variant="simple">
                  <Thead>
                    <Tr>
                      <Th>Name</Th>
                      <Th>Floors</Th>
                      <Th>Difficulty</Th>
                      <Th>Created</Th>
                      <Th>Players</Th>
                      <Th>Action</Th>
                    </Tr>
                  </Thead>
                  <Tbody>
                    {dungeons.map(dungeon => (
                      <Tr 
                        key={dungeon.id}
                        bg={selectedDungeon?.id === dungeon.id ? 'blue.900' : undefined}
                        _hover={{ bg: 'gray.700', cursor: 'pointer' }}
                        onClick={() => handleSelectDungeon(dungeon)}
                      >
                        <Td>{dungeon.name}</Td>
                        <Td>{dungeon.floors}</Td>
                        <Td>
                          <Text
                            color={
                              !dungeon.difficulty ? 'gray.400' :
                              dungeon.difficulty === 'easy' ? 'green.400' :
                              dungeon.difficulty === 'medium' ? 'yellow.400' : 'red.400'
                            }
                            fontWeight="bold"
                          >
                            {dungeon.difficulty ? 
                              dungeon.difficulty.charAt(0).toUpperCase() + dungeon.difficulty.slice(1) : 
                              'Unknown'}
                          </Text>
                        </Td>
                        <Td>{formatDate(dungeon.createdAt)}</Td>
                        <Td>{dungeon.playerCount}</Td>
                        <Td>
                          <Button
                            size="sm"
                            colorScheme="blue"
                            onClick={(e) => {
                              e.stopPropagation();
                              setSelectedDungeon(dungeon);
                              handleJoinDungeon();
                            }}
                          >
                            Select
                          </Button>
                        </Td>
                      </Tr>
                    ))}
                  </Tbody>
                </Table>
              </Box>
            )}
          </Box>

          <Flex justify="space-between" mt={8}>
            <Button
              leftIcon={<ArrowBackIcon />}
              variant="solid" 
              onClick={handleBack}
              bg="gray.600"
              color="cyan.300"
              borderColor="cyan.500"
              borderWidth="1px"
              _hover={{ bg: "gray.700", color: "cyan.200" }}
              _active={{ bg: "gray.800", color: "cyan.100" }}
              boxShadow="sm"
            >
              Back to Characters
            </Button>
            
            <HStack spacing={4}>
              <Button
                leftIcon={<AddIcon />}
                colorScheme="green"
                onClick={onOpen}
              >
                Create New Dungeon
              </Button>
              
              <Button
                rightIcon={<ArrowForwardIcon />}
                colorScheme="blue"
                isDisabled={!selectedDungeon}
                isLoading={joiningDungeon}
                loadingText="Joining..."
                onClick={handleJoinDungeon}
              >
                Join Selected Dungeon
              </Button>
            </HStack>
          </Flex>
        </>
      )}

      {/* Create Dungeon Modal */}
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent bg="gray.800">
          <ModalHeader>Create New Dungeon</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <VStack spacing={4}>
              <FormControl isRequired>
                <FormLabel>Dungeon Name</FormLabel>
                <Input
                  value={newDungeonName}
                  onChange={(e) => setNewDungeonName(e.target.value)}
                  placeholder="Enter dungeon name"
                  bg="gray.700"
                  borderColor="gray.600"
                  _hover={{ borderColor: "blue.400" }}
                  _focus={{ borderColor: "blue.400", boxShadow: "0 0 0 1px var(--chakra-colors-blue-400)" }}
                />
              </FormControl>
              
              <FormControl isRequired>
                <FormLabel>Number of Floors (5-20)</FormLabel>
                <NumberInput
                  min={5}
                  max={20}
                  value={newDungeonFloorCount}
                  onChange={(_, value) => setNewDungeonFloorCount(value)}
                  bg="gray.700"
                  borderRadius="md"
                >
                  <NumberInputField borderColor="gray.600" _hover={{ borderColor: "blue.400" }} />
                  <NumberInputStepper borderColor="gray.600">
                    <NumberIncrementStepper 
                      borderColor="gray.600" 
                      color="gray.200"
                      _hover={{ bg: "blue.500", color: "white" }}
                      _active={{ bg: "blue.600" }}
                    />
                    <NumberDecrementStepper 
                      borderColor="gray.600" 
                      color="gray.200"
                      _hover={{ bg: "blue.500", color: "white" }}
                      _active={{ bg: "blue.600" }}
                    />
                  </NumberInputStepper>
                </NumberInput>
              </FormControl>
              
              <FormControl isRequired>
                <FormLabel>Difficulty</FormLabel>
                <RadioGroup value={difficulty} onChange={setDifficulty}>
                  <Stack direction="row" spacing={5}>
                    <Radio value="easy" colorScheme="green">Easy</Radio>
                    <Radio value="medium" colorScheme="yellow">Medium</Radio>
                    <Radio value="hard" colorScheme="red">Hard</Radio>
                  </Stack>
                </RadioGroup>
              </FormControl>
            </VStack>
          </ModalBody>

          <ModalFooter>
            <Button variant="ghost" mr={3} onClick={onClose}>
              Cancel
            </Button>
            <Button
              colorScheme="blue"
              isLoading={creatingDungeon}
              loadingText="Creating..."
              onClick={handleCreateDungeon}
            >
              Create
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Container>
  );
};

export default DungeonSelection; 