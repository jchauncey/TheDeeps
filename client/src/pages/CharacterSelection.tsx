import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Container,
  Flex,
  Grid,
  Heading,
  Text,
  useToast,
  Spinner,
  Center,
} from '@chakra-ui/react';
import { AddIcon } from '@chakra-ui/icons';
import CharacterCard from '../components/CharacterCard';
import { getCharacters, deleteCharacter } from '../services/api';
import { Character } from '../types';

const CharacterSelection: React.FC = () => {
  const [characters, setCharacters] = useState<Character[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const toast = useToast();

  useEffect(() => {
    fetchCharacters();
  }, []);

  const fetchCharacters = async () => {
    try {
      setLoading(true);
      const data = await getCharacters();
      setCharacters(data);
      setError(null);
    } catch (err) {
      console.error('Failed to fetch characters:', err);
      setError('Failed to load characters. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteCharacter = async (id: string) => {
    try {
      await deleteCharacter(id);
      setCharacters(characters.filter(char => char.id !== id));
      toast({
        title: 'Character deleted',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (err) {
      console.error('Failed to delete character:', err);
      toast({
        title: 'Error',
        description: 'Failed to delete character. Please try again.',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  const handleSelectCharacter = (character: Character) => {
    // Navigate to dungeon selection with the selected character
    navigate('/dungeon-selection', { state: { character } });
  };

  const handleCreateCharacter = () => {
    navigate('/create-character');
  };

  return (
    <Container maxW="container.xl" py={8}>
      <Flex direction="column" align="center" mb={8}>
        <Heading as="h1" size="2xl" mb={2}>
          The Deeps
        </Heading>
        <Text fontSize="xl" color="gray.400">
          Select your character or create a new one
        </Text>
      </Flex>

      {loading ? (
        <Center h="300px">
          <Spinner size="xl" thickness="4px" speed="0.65s" color="blue.500" />
        </Center>
      ) : error ? (
        <Box textAlign="center" p={5} color="red.500">
          <Text>{error}</Text>
          <Button mt={4} onClick={fetchCharacters}>
            Retry
          </Button>
        </Box>
      ) : (
        <>
          <Grid
            templateColumns={{
              base: '1fr',
              md: 'repeat(2, 1fr)',
              lg: 'repeat(3, 1fr)',
              xl: 'repeat(4, 1fr)',
            }}
            gap={6}
            mb={8}
          >
            {characters.map(character => (
              <CharacterCard
                key={character.id}
                character={character}
                onDelete={handleDeleteCharacter}
                onSelect={handleSelectCharacter}
              />
            ))}
          </Grid>

          {characters.length === 0 && (
            <Box textAlign="center" p={10} bg="gray.800" borderRadius="lg">
              <Text fontSize="lg" mb={4}>
                You don't have any characters yet.
              </Text>
            </Box>
          )}

          <Flex justify="center" mt={6}>
            <Button
              leftIcon={<AddIcon />}
              colorScheme="blue"
              size="lg"
              onClick={handleCreateCharacter}
              isDisabled={characters.length >= 10}
            >
              Create New Character
            </Button>
          </Flex>

          {characters.length >= 10 && (
            <Text textAlign="center" color="red.400" mt={2}>
              Maximum number of characters reached (10)
            </Text>
          )}

          <Flex justify="center" mt={8}>
            <Button
              variant="outline"
              colorScheme="teal"
              size="md"
              onClick={() => navigate('/component-playground')}
            >
              Component Playground
            </Button>
          </Flex>
          <Text textAlign="center" fontSize="sm" color="gray.500" mt={2}>
            View and test individual components in isolation
          </Text>
        </>
      )}
    </Container>
  );
};

export default CharacterSelection; 