import { Box, Flex, Text, Progress, Stat, StatLabel, StatNumber, StatHelpText, Divider, Tooltip, Badge } from '@chakra-ui/react';
import { CharacterData, CHARACTER_CLASSES } from '../../types/game';

interface GameStatusProps {
  character: CharacterData | null;
}

export const GameStatus = ({ character }: GameStatusProps) => {
  if (!character) {
    return null;
  }

  // Calculate modifier for a stat
  const calculateModifier = (score: number): number => {
    return Math.floor((score - 10) / 2);
  };

  // Format modifier with + sign for positive values
  const formatModifier = (modifier: number): string => {
    return modifier >= 0 ? `+${modifier}` : `${modifier}`;
  };

  // Calculate max health based on constitution and class hit die
  const calculateMaxHealth = (constitution: number): number => {
    const conModifier = calculateModifier(constitution);
    const classInfo = CHARACTER_CLASSES.find(c => c.id === character.characterClass);
    const hitDie = classInfo?.hitDie || 8;
    
    // Base health + constitution modifier * level
    return hitDie + conModifier;
  };

  // Calculate max mana based on primary spellcasting ability
  const calculateMaxMana = (): number => {
    const classInfo = CHARACTER_CLASSES.find(c => c.id === character.characterClass);
    
    if (!classInfo) return 0;
    
    // For wizard, sorcerer, warlock
    if (['wizard', 'sorcerer'].includes(character.characterClass)) {
      const intModifier = calculateModifier(character.stats.intelligence);
      return 10 + intModifier * 3;
    }
    // For cleric, druid
    else if (['cleric', 'druid'].includes(character.characterClass)) {
      const wisModifier = calculateModifier(character.stats.wisdom);
      return 10 + wisModifier * 3;
    }
    // For bard, paladin, warlock
    else if (['bard', 'paladin', 'warlock'].includes(character.characterClass)) {
      const chaModifier = calculateModifier(character.stats.charisma);
      return 10 + chaModifier * 3;
    }
    // For other classes
    else {
      const bestModifier = Math.max(
        calculateModifier(character.stats.intelligence),
        calculateModifier(character.stats.wisdom),
        calculateModifier(character.stats.charisma)
      );
      return 5 + bestModifier * 2;
    }
  };

  // Calculate armor class based on dexterity
  const calculateArmorClass = (): number => {
    const dexModifier = calculateModifier(character.stats.dexterity);
    
    // Unarmored AC is 10 + DEX modifier
    let ac = 10 + dexModifier;
    
    // Monk and Barbarian have special unarmored defense
    if (character.characterClass === 'monk') {
      const wisModifier = calculateModifier(character.stats.wisdom);
      ac += wisModifier;
    } else if (character.characterClass === 'barbarian') {
      const conModifier = calculateModifier(character.stats.constitution);
      ac += conModifier;
    }
    
    return ac;
  };

  const maxHealth = calculateMaxHealth(character.stats.constitution);
  const maxMana = calculateMaxMana();
  const armorClass = calculateArmorClass();
  
  // Current values (in a real game, these would come from the game state)
  const currentHealth = maxHealth; // For now, assume full health
  const currentMana = maxMana; // For now, assume full mana

  // Get class info
  const classInfo = CHARACTER_CLASSES.find(c => c.id === character.characterClass);

  return (
    <Box
      height="100%"
      width="100%"
      bg="rgba(0, 0, 0, 0.7)"
      p={4}
      borderRadius="md"
      color="white"
      overflowY="auto"
    >
      <Text fontSize="xl" fontWeight="bold" mb={1}>
        {character.name}
      </Text>
      <Text fontSize="md" color="gray.300" mb={3}>
        Level 1 {classInfo?.name || character.characterClass}
      </Text>

      <Divider mb={3} />

      {/* Health, Mana, and AC */}
      <Box mb={4}>
        <Tooltip label={`Health: ${currentHealth}/${maxHealth}`}>
          <Box>
            <Flex justify="space-between" mb={1}>
              <Text fontSize="sm">Health</Text>
              <Text fontSize="sm">{currentHealth}/{maxHealth}</Text>
            </Flex>
            <Progress 
              value={(currentHealth / maxHealth) * 100} 
              colorScheme="red" 
              size="sm" 
              mb={2}
              borderRadius="md"
            />
          </Box>
        </Tooltip>

        <Tooltip label={`Mana: ${currentMana}/${maxMana}`}>
          <Box>
            <Flex justify="space-between" mb={1}>
              <Text fontSize="sm">Mana</Text>
              <Text fontSize="sm">{currentMana}/{maxMana}</Text>
            </Flex>
            <Progress 
              value={(currentMana / maxMana) * 100} 
              colorScheme="blue" 
              size="sm" 
              mb={2}
              borderRadius="md"
            />
          </Box>
        </Tooltip>

        <Flex justify="space-between" mb={1}>
          <Text fontSize="sm">Experience</Text>
          <Text fontSize="sm">0/100</Text>
        </Flex>
        <Progress 
          value={0} 
          colorScheme="purple" 
          size="sm"
          borderRadius="md"
        />
        
        <Flex justify="space-between" mt={3}>
          <Text fontSize="sm" fontWeight="bold">Armor Class</Text>
          <Badge colorScheme="green" fontSize="sm">{armorClass}</Badge>
        </Flex>
      </Box>

      <Divider mb={3} />

      {/* Stats */}
      <Flex justify="space-between" align="center" mb={2}>
        <Text fontSize="md" fontWeight="bold">Stats</Text>
        <Text fontSize="xs" color="gray.400">Value (Modifier)</Text>
      </Flex>
      <Flex flexWrap="wrap" justifyContent="space-between">
        {Object.entries(character.stats).map(([statName, value]) => {
          const modifier = calculateModifier(value);
          const formattedModifier = formatModifier(modifier);
          const modifierColor = modifier >= 0 ? "green.300" : "red.300";
          
          // Check if this is a primary ability for the class
          const isPrimary = classInfo?.primaryAbility.toLowerCase().includes(statName.toLowerCase());
          
          return (
            <Stat 
              key={statName} 
              flex="0 0 48%" 
              mb={2}
              borderLeft={isPrimary ? "2px solid" : "none"}
              borderColor="purple.400"
              pl={isPrimary ? 2 : 0}
            >
              <StatLabel textTransform="capitalize" fontSize="xs">{statName}</StatLabel>
              <Flex align="baseline">
                <StatNumber fontSize="md" mr={1}>{value}</StatNumber>
                <StatHelpText fontSize="xs" color={modifierColor} mt={0}>
                  ({formattedModifier})
                </StatHelpText>
              </Flex>
            </Stat>
          );
        })}
      </Flex>

      <Divider my={3} />

      {/* Class Abilities */}
      <Text fontSize="md" fontWeight="bold" mb={2}>
        Class Abilities
      </Text>
      <Box fontSize="sm" color="gray.300" mb={3}>
        {character.abilities.map((ability, index) => (
          <Badge key={index} colorScheme="purple" mr={2} mb={2}>
            {ability}
          </Badge>
        ))}
      </Box>

      {/* Equipment */}
      <Text fontSize="md" fontWeight="bold" mb={2}>
        Equipment
      </Text>
      <Box fontSize="sm" color="gray.300">
        <Flex mb={1}>
          <Text width="80px" color="gray.400">Weapon:</Text>
          <Text>None</Text>
        </Flex>
        <Flex mb={1}>
          <Text width="80px" color="gray.400">Armor:</Text>
          <Text>None</Text>
        </Flex>
        <Flex mb={1}>
          <Text width="80px" color="gray.400">Shield:</Text>
          <Text>None</Text>
        </Flex>
        <Flex>
          <Text width="80px" color="gray.400">Accessory:</Text>
          <Text>None</Text>
        </Flex>
      </Box>

      <Divider my={3} />

      {/* Proficiencies */}
      <Text fontSize="md" fontWeight="bold" mb={2}>
        Proficiencies
      </Text>
      <Box fontSize="sm" color="gray.300" mb={3}>
        {character.proficiencies.map((prof, index) => (
          <Badge key={index} colorScheme="blue" mr={2} mb={2} variant="outline">
            {prof}
          </Badge>
        ))}
      </Box>

      {/* Inventory */}
      <Text fontSize="md" fontWeight="bold" mb={2}>
        Inventory
      </Text>
      <Box fontSize="sm" color="gray.300">
        <Flex mb={1} align="center">
          <Text width="80px" color="gray.400">Gold:</Text>
          <Flex align="center">
            <Text mr={1}>{character.gold || 0}</Text>
            <Box as="span" color="yellow.400" fontSize="xs">●</Box>
          </Flex>
        </Flex>
        <Text>Empty</Text>
      </Box>

      <Divider my={3} />

      {/* Status Effects */}
      <Text fontSize="md" fontWeight="bold" mb={2}>
        Status Effects
      </Text>
      <Box fontSize="sm" color="gray.300">
        <Text>None</Text>
      </Box>
    </Box>
  );
}; 