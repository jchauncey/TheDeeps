import { Box, Flex, Text, Progress, Stat, StatLabel, StatNumber, StatHelpText, Divider } from '@chakra-ui/react';
import { CharacterData } from '../../types/game';

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
        Level 1 {character.characterClass}
      </Text>

      <Divider mb={3} />

      {/* Health and Experience */}
      <Box mb={4}>
        <Flex justify="space-between" mb={1}>
          <Text fontSize="sm">Health</Text>
          <Text fontSize="sm">100/100</Text>
        </Flex>
        <Progress value={100} colorScheme="red" size="sm" mb={2} />

        <Flex justify="space-between" mb={1}>
          <Text fontSize="sm">Experience</Text>
          <Text fontSize="sm">0/100</Text>
        </Flex>
        <Progress value={0} colorScheme="blue" size="sm" />
      </Box>

      <Divider mb={3} />

      {/* Stats */}
      <Text fontSize="md" fontWeight="bold" mb={2}>
        Stats
      </Text>
      <Flex flexWrap="wrap" justifyContent="space-between">
        {Object.entries(character.stats).map(([statName, value]) => {
          const modifier = calculateModifier(value);
          const formattedModifier = formatModifier(modifier);
          const modifierColor = modifier >= 0 ? "green.300" : "red.300";
          
          return (
            <Stat key={statName} flex="0 0 48%" mb={2}>
              <StatLabel textTransform="capitalize" fontSize="xs">{statName}</StatLabel>
              <StatNumber fontSize="md">{value}</StatNumber>
              <StatHelpText fontSize="xs" color={modifierColor} mt={0}>
                {formattedModifier}
              </StatHelpText>
            </Stat>
          );
        })}
      </Flex>

      <Divider my={3} />

      {/* Equipment */}
      <Text fontSize="md" fontWeight="bold" mb={2}>
        Equipment
      </Text>
      <Box fontSize="sm" color="gray.300">
        <Text>Weapon: None</Text>
        <Text>Armor: None</Text>
        <Text>Accessory: None</Text>
      </Box>

      <Divider my={3} />

      {/* Inventory */}
      <Text fontSize="md" fontWeight="bold" mb={2}>
        Inventory
      </Text>
      <Box fontSize="sm" color="gray.300">
        <Text>Empty</Text>
      </Box>
    </Box>
  );
}; 