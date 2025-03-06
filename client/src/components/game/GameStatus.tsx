import { Box, Flex, Text, Progress, Stat, StatLabel, StatNumber, StatHelpText } from '@chakra-ui/react';
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
      position="absolute"
      top={4}
      right={4}
      width="250px"
      bg="rgba(0, 0, 0, 0.7)"
      p={4}
      borderRadius="md"
      color="white"
    >
      <Text fontSize="xl" fontWeight="bold" mb={2}>
        {character.name}
      </Text>
      <Text fontSize="md" color="gray.300" mb={4}>
        Level 1 {character.characterClass}
      </Text>

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
    </Box>
  );
}; 