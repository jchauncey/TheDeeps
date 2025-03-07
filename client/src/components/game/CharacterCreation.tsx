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
  Spinner,
  Divider,
  Tooltip
} from '@chakra-ui/react'
import { useState, useEffect } from 'react'
import { AddIcon, MinusIcon, InfoIcon } from '@chakra-ui/icons'
import { CharacterData, CHARACTER_CLASSES } from '../../types/game'
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
  const [pointsRemaining, setPointsRemaining] = useState(27); // Using D&D 5e point buy system
  const [autoAllocated, setAutoAllocated] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const toast = useToast();

  // Auto-allocate points when class changes
  useEffect(() => {
    if (characterClass && !autoAllocated) {
      const selectedClass = CHARACTER_CLASSES.find(c => c.id === characterClass);
      if (selectedClass) {
        // Reset stats to base
        const newStats = {...BASE_STATS};
        
        // Apply class recommended stats
        Object.entries(selectedClass.recommendedStats).forEach(([stat, value]) => {
          newStats[stat as keyof typeof newStats] = value;
        });
        
        setStats(newStats);
        
        // Calculate points used
        let pointsUsed = 0;
        Object.values(newStats).forEach(value => {
          pointsUsed += calculatePointCost(value);
        });
        
        setPointsRemaining(27 - pointsUsed);
        setAutoAllocated(true);
        
        toast({
          title: "Stats auto-allocated",
          description: `Points have been allocated based on the ${selectedClass.name} class. You can still adjust them manually.`,
          status: "info",
          position: "top",
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
    setPointsRemaining(27);
    setAutoAllocated(false);
  };

  // Calculate point cost for a stat value (D&D 5e point buy system)
  const calculatePointCost = (value: number): number => {
    if (value <= 8) return 0;
    if (value === 9) return 1;
    if (value === 10) return 2;
    if (value === 11) return 3;
    if (value === 12) return 4;
    if (value === 13) return 5;
    if (value === 14) return 7;
    if (value === 15) return 9;
    return 0; // Default case
  };

  const handleStatChange = (stat: keyof typeof stats, change: number) => {
    // Don't allow stats below 8 or above 15 (D&D 5e point buy limits)
    const newValue = stats[stat] + change;
    if (newValue < 8 || newValue > 15) return;
    
    // Calculate point cost difference
    const currentCost = calculatePointCost(stats[stat]);
    const newCost = calculatePointCost(newValue);
    const pointDifference = newCost - currentCost;
    
    // Check if we have enough points
    if (pointsRemaining < pointDifference) return;
    
    setStats({
      ...stats,
      [stat]: newValue
    });
    
    setPointsRemaining(pointsRemaining - pointDifference);
  };

  const handleCreateCharacter = async () => {
    if (!name || !characterClass) return;
    
    const selectedClass = CHARACTER_CLASSES.find(c => c.id === characterClass);
    if (!selectedClass) return;
    
    const characterData: CharacterData = {
      name,
      characterClass,
      stats,
      abilities: selectedClass.abilities,
      proficiencies: selectedClass.proficiencies
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
          position: "top",
        });
        
        // Notify parent component
        onCreateCharacter(characterData);
      } else {
        toast({
          title: "Error creating character",
          description: result.message || "There was a problem saving your character.",
          status: "error",
          position: "top",
        });
      }
    } catch (error) {
      toast({
        title: "Error creating character",
        description: "There was a problem connecting to the server.",
        status: "error",
        position: "top",
      });
      console.error("Error creating character:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const getSelectedClass = () => {
    return CHARACTER_CLASSES.find(c => c.id === characterClass);
  };

  // Calculate health and mana based on stats and class
  const calculateHealth = () => {
    const selectedClass = getSelectedClass();
    if (!selectedClass) return 0;
    
    const conModifier = calculateModifier(stats.constitution);
    return selectedClass.hitDie + conModifier;
  };

  const calculateMana = () => {
    const selectedClass = getSelectedClass();
    if (!selectedClass) return 0;
    
    // Spellcasting classes
    if (['wizard', 'sorcerer', 'warlock'].includes(characterClass)) {
      return 10 + calculateModifier(stats.intelligence) * 2;
    } else if (['cleric', 'druid'].includes(characterClass)) {
      return 10 + calculateModifier(stats.wisdom) * 2;
    } else if (['bard', 'paladin'].includes(characterClass)) {
      return 10 + calculateModifier(stats.charisma) * 2;
    } else {
      // Non-spellcasting classes get less mana
      return 5 + Math.max(
        calculateModifier(stats.intelligence),
        calculateModifier(stats.wisdom),
        calculateModifier(stats.charisma)
      );
    }
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
              {CHARACTER_CLASSES.map(c => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </Select>
            
            {characterClass && (
              <Box mt={3} p={3} bg="whiteAlpha.100" borderRadius="md">
                <Text fontSize="md" fontWeight="bold" mb={1}>{getSelectedClass()?.name}</Text>
                <Text fontSize="sm" color="gray.300" mb={2}>{getSelectedClass()?.description}</Text>
                
                <Divider my={2} />
                
                <Text fontSize="sm" mb={1}>
                  <Text as="span" fontWeight="bold">Primary Ability:</Text> {getSelectedClass()?.primaryAbility}
                </Text>
                <Text fontSize="sm" mb={1}>
                  <Text as="span" fontWeight="bold">Hit Die:</Text> d{getSelectedClass()?.hitDie}
                </Text>
                <Text fontSize="sm" mb={1}>
                  <Text as="span" fontWeight="bold">Saving Throws:</Text> {getSelectedClass()?.savingThrows.join(', ')}
                </Text>
                
                <Divider my={2} />
                
                <Text fontSize="sm" fontWeight="bold" mb={1}>Starting Abilities:</Text>
                <Text fontSize="sm" color="gray.300">
                  {getSelectedClass()?.abilities.join(', ')}
                </Text>
              </Box>
            )}
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
                
                // Determine if this is a primary stat for the selected class
                const selectedClass = getSelectedClass();
                const isPrimary = selectedClass?.primaryAbility.toLowerCase().includes(statName.toLowerCase());
                
                return (
                  <Flex 
                    key={statName} 
                    justify="space-between" 
                    align="center" 
                    p={2} 
                    bg={isPrimary ? "whiteAlpha.200" : "whiteAlpha.100"} 
                    borderRadius="md"
                    borderLeft={isPrimary ? "3px solid" : "none"}
                    borderColor="purple.400"
                  >
                    <Stat>
                      <Flex align="center">
                        <StatLabel textTransform="capitalize" fontSize="sm">{statName}</StatLabel>
                        {isPrimary && (
                          <Tooltip label="Primary ability for this class">
                            <InfoIcon ml={1} color="purple.300" fontSize="xs" />
                          </Tooltip>
                        )}
                      </Flex>
                      <Flex align="baseline">
                        <StatNumber fontSize="md">{value}</StatNumber>
                        <StatHelpText fontSize="xs" color={modifierColor} mt={0} ml={1}>
                          ({formattedModifier})
                        </StatHelpText>
                      </Flex>
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
                        isDisabled={pointsRemaining <= 0 || stats[statName as keyof typeof stats] >= 15 || isSubmitting}
                        onClick={() => handleStatChange(statName as keyof typeof stats, 1)}
                      />
                    </HStack>
                  </Flex>
                );
              })}
            </VStack>
          </Box>
          
          {/* Character Preview */}
          {characterClass && (
            <Box p={4} bg="whiteAlpha.100" borderRadius="md">
              <Text fontSize="md" fontWeight="bold" mb={3}>Character Preview</Text>
              
              <Flex justify="space-between" wrap="wrap">
                <Box flex="1" minW="200px" mr={4}>
                  <Text fontSize="sm" mb={1}>
                    <Text as="span" fontWeight="bold">Health:</Text> {calculateHealth()}
                  </Text>
                  <Text fontSize="sm" mb={1}>
                    <Text as="span" fontWeight="bold">Mana:</Text> {calculateMana()}
                  </Text>
                  <Text fontSize="sm" mb={1}>
                    <Text as="span" fontWeight="bold">Armor Class:</Text> {10 + calculateModifier(stats.dexterity)}
                  </Text>
                </Box>
                
                <Box flex="1" minW="200px">
                  <Text fontSize="sm" fontWeight="bold" mb={1}>Proficiencies:</Text>
                  <Text fontSize="xs" color="gray.300">
                    {getSelectedClass()?.proficiencies.join(', ')}
                  </Text>
                </Box>
              </Flex>
            </Box>
          )}
          
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