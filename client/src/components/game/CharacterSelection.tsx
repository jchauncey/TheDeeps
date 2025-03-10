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
  Badge, 
  Spinner, 
  Image,
  useToast
} from '@chakra-ui/react';
import { useState, useEffect } from 'react';
import { getSavedCharacters } from '../../services/api';
import { CharacterData } from '../../types/game';

interface CharacterSelectionProps {
  onSelectCharacter: (characterId: string) => void;
  onCreateNewCharacter: () => void;
  onBack: () => void;
  refreshTrigger?: number; // Optional prop to trigger a refresh
}

export const CharacterSelection = ({ 
  onSelectCharacter, 
  onCreateNewCharacter, 
  onBack,
  refreshTrigger = 0
}: CharacterSelectionProps) => {
  const [characters, setCharacters] = useState<{ id: string; name: string; characterClass: string }[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedCharacterId, setSelectedCharacterId] = useState<string | null>(null);
  const toast = useToast();

  // Function to load characters
  const loadCharacters = async () => {
    setIsLoading(true);
    try {
      const result = await getSavedCharacters();
      if (result.success && result.characters) {
        // Sort characters by name
        const sortedCharacters = [...result.characters].sort((a, b) => 
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
                <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={4} mb={6} flex={1} overflowY="auto">
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
                    >
                      <CardHeader pb={2}>
                        <Flex justify="space-between" align="center">
                          <Heading size="md">{character.name}</Heading>
                          <Badge colorScheme={getClassColor(character.characterClass)}>
                            {character.characterClass}
                          </Badge>
                        </Flex>
                      </CardHeader>
                      <CardBody>
                        {/* We could add more character details here in the future */}
                        <Text>Ready for adventure!</Text>
                      </CardBody>
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
    </Box>
  );
}; 