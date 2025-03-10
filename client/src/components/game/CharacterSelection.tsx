import { 
  Box, 
  Button, 
  Flex, 
  Heading, 
  Text, 
  VStack, 
  HStack, 
  SimpleGrid, 
  Card, 
  CardBody, 
  CardHeader, 
  CardFooter,
  Badge, 
  Spinner, 
  Image,
  IconButton,
  useToast,
  useDisclosure,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Progress,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  Divider
} from '@chakra-ui/react';
import { useState, useEffect } from 'react';
import { DeleteIcon } from '@chakra-ui/icons';
import { getSavedCharacters, deleteCharacter, loadCharacter } from '../../services/api';
import { CharacterData } from '../../types/game';

interface CharacterSelectionProps {
  onSelectCharacter: (characterId: string) => void;
  onCreateNewCharacter: () => void;
  onBack: () => void;
  refreshTrigger?: number; // Optional prop to trigger a refresh
}

// Ensure all characters have required fields
type SafeCharacterData = CharacterData & {
  id: string; // Ensure id is always present and not undefined
};

export const CharacterSelection = ({ 
  onSelectCharacter, 
  onCreateNewCharacter, 
  onBack,
  refreshTrigger = 0
}: CharacterSelectionProps) => {
  const [characters, setCharacters] = useState<SafeCharacterData[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedCharacterId, setSelectedCharacterId] = useState<string | null>(null);
  const [characterToDelete, setCharacterToDelete] = useState<{ id: string; name: string } | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();

  // Function to load characters
  const loadCharacters = async () => {
    setIsLoading(true);
    try {
      const result = await getSavedCharacters();
      if (result.success && result.characters) {
        // Get full character data for each character
        const fullCharacters = await Promise.all(
          result.characters.map(async (char) => {
            const charResult = await loadCharacter(char.id);
            if (charResult.success && charResult.character) {
              // Ensure id is present
              return {
                ...charResult.character,
                id: char.id
              } as SafeCharacterData;
            } else {
              // Create a valid CharacterData object with default values
              return {
                id: char.id,
                name: char.name,
                characterClass: char.characterClass,
                stats: { strength: 10, dexterity: 10, constitution: 10, intelligence: 10, wisdom: 10, charisma: 10 },
                abilities: [],
                proficiencies: [],
                gold: 0,
                level: 1,
                health: 100,
                maxHealth: 100,
                mana: 50,
                maxMana: 50,
                experience: 0
              } as SafeCharacterData;
            }
          })
        );
        
        // Sort characters by name
        const sortedCharacters = [...fullCharacters].sort((a, b) => 
          a.name.localeCompare(b.name)
        );
        setCharacters(sortedCharacters);
      } else {
        toast({
          title: 'Error',
          description: result.message || 'Failed to load characters',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    } catch (error) {
      console.error('Error loading characters:', error);
      toast({
        title: 'Error',
        description: 'An unexpected error occurred while loading characters',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  // Load saved characters on component mount and when refreshTrigger changes
  useEffect(() => {
    loadCharacters();
  }, [refreshTrigger]);

  const handleSelectCharacter = (characterId: string) => {
    setSelectedCharacterId(characterId);
  };

  const handleConfirmSelection = () => {
    if (selectedCharacterId) {
      onSelectCharacter(selectedCharacterId);
    }
  };

  // Handle character deletion
  const handleDeleteClick = (e: React.MouseEvent, character: { id: string; name: string }) => {
    e.stopPropagation(); // Prevent card selection when clicking delete
    setCharacterToDelete(character);
    onOpen();
  };

  const handleConfirmDelete = async () => {
    if (!characterToDelete) return;
    
    setIsDeleting(true);
    try {
      const result = await deleteCharacter(characterToDelete.id);
      if (result.success) {
        toast({
          title: 'Character Deleted',
          description: `${characterToDelete.name} has been deleted.`,
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
        
        // If the deleted character was selected, clear the selection
        if (selectedCharacterId === characterToDelete.id) {
          setSelectedCharacterId(null);
        }
        
        // Reload the character list
        loadCharacters();
      } else {
        toast({
          title: 'Error',
          description: result.message || 'Failed to delete character',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    } catch (error) {
      console.error('Error deleting character:', error);
      toast({
        title: 'Error',
        description: 'An unexpected error occurred while deleting the character',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsDeleting(false);
      onClose();
    }
  };

  // Get class color for badges
  const getClassColor = (characterClass: string): string => {
    const classColors: Record<string, string> = {
      warrior: 'red',
      mage: 'blue',
      rogue: 'green',
      cleric: 'yellow',
      ranger: 'teal',
      paladin: 'purple',
      druid: 'orange',
      bard: 'pink',
      monk: 'cyan',
      warlock: 'gray',
    };
    
    return classColors[characterClass.toLowerCase()] || 'gray';
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
      >
        <Flex justify="space-between" align="center" mb={6}>
          <Heading color="white" size="lg">Character Selection</Heading>
          <Button variant="outline" colorScheme="purple" onClick={onBack}>
            Back to Title
          </Button>
        </Flex>

        {isLoading ? (
          <Flex justify="center" align="center" flex={1}>
            <Spinner size="xl" color="purple.500" />
          </Flex>
        ) : (
          <>
            {characters.length === 0 ? (
              <Flex direction="column" justify="center" align="center" flex={1} gap={6}>
                <Text color="white" fontSize="xl">You don't have any characters yet.</Text>
                <Button 
                  colorScheme="purple" 
                  size="lg" 
                  onClick={onCreateNewCharacter}
                >
                  Create Your First Character
                </Button>
              </Flex>
            ) : (
              <>
                <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={3} mb={6} flex={1} overflowY="auto">
                  {characters.map((character) => (
                    <Card 
                      key={character.id} 
                      cursor="pointer"
                      bg={selectedCharacterId === character.id ? "purple.800" : "gray.800"}
                      color="white"
                      borderWidth={2}
                      borderColor={selectedCharacterId === character.id ? "purple.400" : "transparent"}
                      onClick={() => handleSelectCharacter(character.id)}
                      _hover={{ 
                        transform: "translateY(-2px)", 
                        boxShadow: "lg",
                        borderColor: selectedCharacterId === character.id ? "purple.400" : "purple.200"
                      }}
                      transition="all 0.2s"
                      position="relative"
                      maxW="250px"
                      maxH="220px"
                      fontSize="sm"
                    >
                      <CardHeader pb={1} pt={2} px={3}>
                        <Flex justify="space-between" align="center">
                          <Heading size="sm" noOfLines={1}>{character.name}</Heading>
                          <Badge colorScheme={getClassColor(character.characterClass)} fontSize="xs">
                            {character.characterClass}
                          </Badge>
                        </Flex>
                        <Text fontSize="md" fontWeight="bold">
                          Level {character.level || 1}
                        </Text>
                      </CardHeader>
                      <CardBody py={1} px={3}>
                        <VStack spacing={1} align="stretch">
                          <HStack justify="space-between" fontSize="xs">
                            <Text fontWeight="semibold">HP:</Text>
                            <Text>{character.health || 0}/{character.maxHealth || 100}</Text>
                          </HStack>
                          <Progress 
                            value={(character.health || 0) / (character.maxHealth || 100) * 100} 
                            colorScheme="red" 
                            size="xs" 
                            borderRadius="md"
                          />
                          
                          <HStack justify="space-between" fontSize="xs">
                            <Text fontWeight="semibold">Mana:</Text>
                            <Text>{character.mana || 0}/{character.maxMana || 50}</Text>
                          </HStack>
                          <Progress 
                            value={(character.mana || 0) / (character.maxMana || 50) * 100} 
                            colorScheme="blue" 
                            size="xs" 
                            borderRadius="md"
                          />
                          
                          <HStack justify="space-between" fontSize="xs">
                            <Text fontWeight="semibold">XP:</Text>
                            <Text>{character.experience || 0}</Text>
                          </HStack>
                          <Progress 
                            value={((character.experience || 0) % 1000) / 1000 * 100} 
                            colorScheme="green" 
                            size="xs" 
                            borderRadius="md"
                          />
                        </VStack>
                      </CardBody>
                      <CardFooter pt={0} pb={1} px={3} justifyContent="flex-end">
                        <IconButton
                          aria-label="Delete character"
                          icon={<DeleteIcon />}
                          size="xs"
                          colorScheme="red"
                          variant="ghost"
                          onClick={(e) => handleDeleteClick(e, { id: character.id, name: character.name })}
                          _hover={{ bg: 'rgba(229, 62, 62, 0.3)' }}
                        />
                      </CardFooter>
                    </Card>
                  ))}
                </SimpleGrid>

                <Flex justify="space-between" mt={4}>
                  {characters.length < 10 && (
                    <Button 
                      colorScheme="teal" 
                      onClick={onCreateNewCharacter}
                    >
                      Create New Character
                    </Button>
                  )}
                  <Button 
                    colorScheme="purple" 
                    isDisabled={!selectedCharacterId}
                    onClick={handleConfirmSelection}
                    ml="auto"
                  >
                    Select Character
                  </Button>
                </Flex>
              </>
            )}
          </>
        )}
      </Flex>

      {/* Confirmation Modal for Character Deletion */}
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent bg="gray.800" color="white">
          <ModalHeader>Confirm Deletion</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            {characterToDelete && (
              <Text>
                Are you sure you want to delete <strong>{characterToDelete.name}</strong>? 
                This action cannot be undone.
              </Text>
            )}
          </ModalBody>

          <ModalFooter>
            <Button variant="solid" colorScheme="gray" mr={3} onClick={onClose}>
              Cancel
            </Button>
            <Button 
              colorScheme="red" 
              onClick={handleConfirmDelete}
              isLoading={isDeleting}
              loadingText="Deleting..."
            >
              Delete
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Box>
  );
}; 