import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  Box,
  Flex,
  Text,
  useToast,
  Spinner,
  Center,
} from '@chakra-ui/react';
import { Character, Dungeon } from '../types';

interface LocationState {
  character: Character;
  dungeon: Dungeon;
}

const Game: React.FC = () => {
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  
  const navigate = useNavigate();
  const location = useLocation();
  const toast = useToast();
  
  const { character, dungeon } = (location.state as LocationState) || {};

  useEffect(() => {
    if (!character || !dungeon) {
      navigate('/');
      toast({
        title: 'Error',
        description: 'Character or dungeon information missing',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
      return;
    }
    
    // Initialize game
    const initGame = async () => {
      try {
        setLoading(true);
        
        // TODO: Initialize game logic here
        
        // Show welcome message
        toast({
          title: 'Entered Dungeon',
          description: `You have entered ${dungeon.name}, floor 1`,
          status: 'info',
          duration: 3000,
          isClosable: true,
        });
        
        setLoading(false);
      } catch (err) {
        console.error('Failed to initialize game:', err);
        setError('Failed to initialize game. Please try again.');
        setLoading(false);
      }
    };
    
    initGame();
  }, [character, dungeon, navigate, toast]);

  if (loading) {
    return (
      <Center h="100vh">
        <Spinner size="xl" thickness="4px" speed="0.65s" color="blue.500" />
      </Center>
    );
  }

  if (error) {
    return (
      <Center h="100vh">
        <Text color="red.500">{error}</Text>
      </Center>
    );
  }

  return (
    <Flex h="100vh">
      {/* Left side - Map window */}
      <Box flex="3" bg="gray.900" p={4}>
        <Text fontSize="xl" mb={4}>
          {dungeon.name} - Floor {character.currentFloor}
        </Text>
        <Box 
          bg="black" 
          color="white" 
          fontFamily="monospace" 
          p={4} 
          borderRadius="md"
          h="calc(100% - 60px)"
          overflow="auto"
        >
          {/* Placeholder for the game map */}
          <Text>Game map will be displayed here</Text>
          <Text mt={4}>Character: {character.name} ({character.class})</Text>
          <Text>Position: ({character.position.x}, {character.position.y})</Text>
        </Box>
      </Box>
      
      {/* Right side - Character status panel */}
      <Box flex="1" bg="gray.800" p={4}>
        <Text fontSize="xl" mb={4}>
          {character.name}
        </Text>
        <Text>Level: {character.level}</Text>
        <Text>Class: {character.class}</Text>
        <Text mt={2}>HP: {character.currentHp}/{character.maxHp}</Text>
        <Text>Mana: {character.currentMana}/{character.maxMana}</Text>
        <Text mt={2}>Floor: {character.currentFloor}</Text>
      </Box>
    </Flex>
  );
};

export default Game; 