import { 
  Box, 
  Flex, 
  VStack, 
  HStack, 
  Text, 
  Input, 
  Button, 
  Select, 
  Heading, 
  Stat, 
  StatLabel, 
  StatNumber, 
  StatHelpText,
  IconButton,
  useToast,
  Spinner
} from '@chakra-ui/react'
import { useState, useEffect } from 'react'
import { AddIcon, MinusIcon } from '@chakra-ui/icons'
import { CharacterData } from '../../types/game'
import { createCharacter } from '../../services/api'

interface CharacterCreationProps {
  onCreateCharacter: (character: CharacterData) => void;
  onBack: () => void;
}

// Base stats for all characters
const BASE_STATS = {
  strength: 8,
  dexterity: 8,
  constitution: 8,
  intelligence: 8,
  wisdom: 8,
  charisma: 8
};

// Class-specific stat allocations
const CLASS_STAT_ALLOCATIONS = {
  warrior: {
    strength: 4,
    constitution: 3,
    dexterity: 2,
    wisdom: 1,
    intelligence: 0,
    charisma: 0
  },
  rogue: {
    dexterity: 4,
    charisma: 2,
    intelligence: 2,
    constitution: 1,
    strength: 1,
    wisdom: 0
  },
  mage: {
    intelligence: 5,
    wisdom: 2,
    dexterity: 2,
    constitution: 1,
    charisma: 0,
    strength: 0
  },
  cleric: {
    wisdom: 4,
    charisma: 3,
    constitution: 2,
    strength: 1,
    intelligence: 0,
    dexterity: 0
  }
};

const CLASSES = [
  { id: 'warrior', name: 'Warrior', description: 'Strong melee fighter with heavy armor' },
  { id: 'rogue', name: 'Rogue', description: 'Stealthy character with high dexterity' },
  { id: 'mage', name: 'Mage', description: 'Powerful spellcaster with arcane knowledge' },
  { id: 'cleric', name: 'Cleric', description: 'Divine spellcaster with healing abilities' }
];

// Calculate the ability score modifier using D&D rules
const calculateModifier = (score: number): number => {
  return Math.floor((score - 10) / 2);
};

// Format the modifier with a + sign for positive values
const formatModifier = (modifier: number): string => {
  return modifier >= 0 ? `+${modifier}` : `${modifier}`;
};

export const CharacterCreation = ({ onCreateCharacter, onBack }: CharacterCreationProps) => {
  const [name, setName] = useState('');
  const [characterClass, setCharacterClass] = useState('');
  const [stats, setStats] = useState({...BASE_STATS});
  const [pointsRemaining, setPointsRemaining] = useState(10);
  const [autoAllocated, setAutoAllocated] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const toast = useToast();

  // Auto-allocate points when class changes
  useEffect(() => {
    if (characterClass && !autoAllocated) {
      const allocation = CLASS_STAT_ALLOCATIONS[characterClass as keyof typeof CLASS_STAT_ALLOCATIONS];
      if (allocation) {
        // Reset stats to base
        const newStats = {...BASE_STATS};
        
        // Apply class allocation
        Object.entries(allocation).forEach(([stat, points]) => {
          newStats[stat as keyof typeof newStats] += points;
        });
        
        setStats(newStats);
        setPointsRemaining(0);
        setAutoAllocated(true);
        
        toast({
          title: "Stats auto-allocated",
          description: `Points have been allocated based on the ${characterClass} class. You can still adjust them manually.`,
          status: "info",
          duration: 5000,
          isClosable: true,
        });
      }
    }
  }, [characterClass, autoAllocated, toast]);

  const handleClassChange = (newClass: string) => {
    setCharacterClass(newClass);
    setAutoAllocated(false);
  };

  const resetStats = () => {
    setStats({...BASE_STATS});
    setPointsRemaining(10);
    setAutoAllocated(false);
  };

  const handleStatChange = (stat: keyof typeof stats, change: number) => {
    // Don't allow stats below 8 or above 18
    const newValue = stats[stat] + change;
    if (newValue < 8 || newValue > 18) return;
    
    // Check if we have enough points
    if (change > 0 && pointsRemaining < 1) return;
    
    setStats({
      ...stats,
      [stat]: newValue
    });
    
    setPointsRemaining(pointsRemaining - change);
  };

  const handleCreateCharacter = async () => {
    if (!name || !characterClass) return;
    
    const characterData: CharacterData = {
      name,
      characterClass,
      stats
    };
    
    setIsSubmitting(true);
    
    try {
      // Send character data to server
      const result = await createCharacter(characterData);
      
      if (result.success) {
        toast({
          title: "Character created",
          description: "Your character has been saved successfully.",
          status: "success",
          duration: 3000,
          isClosable: true,
        });
        
        // Notify parent component
        onCreateCharacter(characterData);
      } else {
        toast({
          title: "Error creating character",
          description: result.message || "There was a problem saving your character.",
          status: "error",
          duration: 5000,
          isClosable: true,
        });
      }
    } catch (error) {
      toast({
        title: "Error creating character",
        description: "There was a problem connecting to the server.",
        status: "error",
        duration: 5000,
        isClosable: true,
      });
      console.error("Error creating character:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const getClassDescription = () => {
    const selectedClass = CLASSES.find(c => c.id === characterClass);
    return selectedClass ? selectedClass.description : 'Select a class to see description';
  };

  return (
    <Box
      position="fixed"
      top={0}
      left={0}
      right={0}
      bottom={0}
      bg="#291326"
      color="white"
      p={8}
      overflowY="auto"
    >
      <Flex direction="column" maxW="800px" mx="auto">
        <Heading size="xl" mb={6} textAlign="center">Character Creation</Heading>
        
        <VStack spacing={8} align="stretch">
          {/* Name Input */}
          <Box>
            <Text mb={2} fontSize="lg">Character Name</Text>
            <Input 
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter character name"
              size="lg"
              bg="whiteAlpha.200"
              borderColor="purple.300"
              _hover={{ borderColor: 'purple.400' }}
              _focus={{ borderColor: 'purple.500' }}
            />
          </Box>
          
          {/* Class Selection */}
          <Box>
            <Text mb={2} fontSize="lg">Character Class</Text>
            <Select 
              value={characterClass}
              onChange={(e) => handleClassChange(e.target.value)}
              placeholder="Select class"
              size="lg"
              bg="whiteAlpha.200"
              borderColor="purple.300"
              _hover={{ borderColor: 'purple.400' }}
              _focus={{ borderColor: 'purple.500' }}
            >
              {CLASSES.map(c => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </Select>
            <Text mt={2} fontSize="md" color="gray.300">{getClassDescription()}</Text>
          </Box>
          
          {/* Stats Allocation */}
          <Box>
            <Flex justify="space-between" align="center" mb={4}>
              <Text fontSize="lg">Character Stats</Text>
              <HStack>
                <Text color={pointsRemaining > 0 ? "green.300" : "gray.400"}>
                  Points remaining: {pointsRemaining}
                </Text>
                <Button 
                  size="sm" 
                  colorScheme="purple" 
                  variant="outline" 
                  onClick={resetStats}
                >
                  Reset Stats
                </Button>
              </HStack>
            </Flex>
            
            <VStack spacing={3} align="stretch">
              {Object.entries(stats).map(([statName, value]) => {
                const modifier = calculateModifier(value);
                const formattedModifier = formatModifier(modifier);
                const modifierColor = modifier >= 0 ? "green.300" : "red.300";
                
                return (
                  <Flex key={statName} justify="space-between" align="center" p={2} bg="whiteAlpha.100" borderRadius="md">
                    <Stat>
                      <StatLabel textTransform="capitalize">{statName}</StatLabel>
                      <StatNumber>{value}</StatNumber>
                      <StatHelpText color={modifierColor}>
                        Modifier: {formattedModifier}
                      </StatHelpText>
                    </Stat>
                    <HStack>
                      <IconButton
                        aria-label={`Decrease ${statName}`}
                        icon={<MinusIcon />}
                        size="sm"
                        colorScheme="purple"
                        variant="outline"
                        isDisabled={stats[statName as keyof typeof stats] <= 8 || isSubmitting}
                        onClick={() => handleStatChange(statName as keyof typeof stats, -1)}
                      />
                      <IconButton
                        aria-label={`Increase ${statName}`}
                        icon={<AddIcon />}
                        size="sm"
                        colorScheme="purple"
                        isDisabled={pointsRemaining <= 0 || stats[statName as keyof typeof stats] >= 18 || isSubmitting}
                        onClick={() => handleStatChange(statName as keyof typeof stats, 1)}
                      />
                    </HStack>
                  </Flex>
                );
              })}
            </VStack>
          </Box>
          
          {/* Buttons */}
          <HStack spacing={4} justify="center" mt={6}>
            <Button
              onClick={onBack}
              size="md"
              variant="outline"
              colorScheme="purple"
              isDisabled={isSubmitting}
            >
              Back
            </Button>
            <Button
              onClick={handleCreateCharacter}
              size="md"
              colorScheme="purple"
              isDisabled={!name || !characterClass || isSubmitting}
              leftIcon={isSubmitting ? <Spinner size="sm" /> : undefined}
            >
              {isSubmitting ? 'Creating...' : 'Create Character'}
            </Button>
          </HStack>
        </VStack>
      </Flex>
    </Box>
  );
}; 