import { Box, Flex, Text, Progress, Stat, StatLabel, StatNumber, StatHelpText, Divider, Tooltip, Badge, Avatar, Icon, Grid, GridItem } from '@chakra-ui/react';
import { CharacterData, CHARACTER_CLASSES } from '../../types/game';
import { FaHeart, FaFlask, FaBrain, FaShieldAlt, FaRunning, FaStar } from 'react-icons/fa';

interface GameStatusProps {
  character: CharacterData | null;
}

// Define character class-specific colors for styling
const CLASS_COLORS = {
  warrior: { primary: '#f55', secondary: '#922', icon: '⚔️' },
  mage: { primary: '#55f', secondary: '#229', icon: '🔮' },
  rogue: { primary: '#5c5', secondary: '#292', icon: '🗡️' },
  cleric: { primary: '#ff5', secondary: '#992', icon: '✨' },
  ranger: { primary: '#5f5', secondary: '#292', icon: '🏹' },
  paladin: { primary: '#f5f', secondary: '#929', icon: '🛡️' },
  bard: { primary: '#f95', secondary: '#952', icon: '🎵' },
  monk: { primary: '#5ff', secondary: '#299', icon: '👊' },
  druid: { primary: '#9f5', secondary: '#592', icon: '🍃' },
  barbarian: { primary: '#f55', secondary: '#922', icon: '🪓' },
  sorcerer: { primary: '#95f', secondary: '#529', icon: '🌟' },
  warlock: { primary: '#a5f', secondary: '#529', icon: '👁️' },
  // Default for any unspecified class
  default: { primary: '#ff0', secondary: '#990', icon: '🧙' },
};

export const GameStatus = ({ character }: GameStatusProps) => {
  if (!character) {
    return (
      <Box
        height="100%"
        width="100%"
        bg="rgba(0, 0, 0, 0.7)"
        p={4}
        borderRadius="md"
        color="white"
        display="flex"
        alignItems="center"
        justifyContent="center"
      >
        <Text>No character data available</Text>
      </Box>
    );
  }

  // Find class info
  const classInfo = CHARACTER_CLASSES.find(c => c.id === character.characterClass);
  
  // Get class colors
  const classColors = CLASS_COLORS[character.characterClass as keyof typeof CLASS_COLORS] || CLASS_COLORS.default;
  
  // Calculate derived stats
  const calculateModifier = (stat: number) => Math.floor((stat - 10) / 2);
  const formatModifier = (mod: number) => mod >= 0 ? `+${mod}` : `${mod}`;
  
  // Calculate health, mana, and AC
  const constitutionMod = calculateModifier(character.stats.constitution);
  const intelligenceMod = calculateModifier(character.stats.intelligence);
  const dexterityMod = calculateModifier(character.stats.dexterity);
  
  // Get character level with default value of 1
  const characterLevel = character.level || 1;
  
  const maxHealth = classInfo ? 
    classInfo.hitDie + constitutionMod + (characterLevel - 1) * (Math.floor(classInfo.hitDie / 2) + 1 + constitutionMod) : 
    10 + constitutionMod;
  
  const maxMana = classInfo?.primaryAbility.toLowerCase().includes('intelligence') ? 
    10 + intelligenceMod * 2 + (characterLevel - 1) * (4 + intelligenceMod) : 
    5 + intelligenceMod + (characterLevel - 1) * (2 + Math.floor(intelligenceMod / 2));
  
  const armorClass = 10 + dexterityMod;
  
  // Current values (in a real game, these would come from the character state)
  const currentHealth = character.health || maxHealth;
  const currentMana = character.mana || maxMana;
  const currentXP = character.experience || 0;
  const nextLevelXP = characterLevel * 100;

  return (
    <Box
      height="100%"
      width="100%"
      bg="rgba(0, 0, 0, 0.7)"
      p={4}
      borderRadius="md"
      color="white"
      overflowY="auto"
      borderLeft={`4px solid ${classColors.primary}`}
    >
      {/* Character Header with Avatar */}
      <Flex align="center" mb={3}>
        <Avatar 
          size="md" 
          name={character.name} 
          bg={classColors.primary}
          color="white"
          icon={<Text fontSize="xl">{classColors.icon}</Text>}
          mr={3}
        />
        <Box>
          <Text fontSize="xl" fontWeight="bold" mb={0}>
            {character.name}
          </Text>
          <Flex align="center">
            <Badge 
              colorScheme="purple" 
              mr={2}
              px={2}
              py={0.5}
              borderRadius="full"
              bg={classColors.primary}
              color="white"
            >
              Level {characterLevel}
            </Badge>
            <Text fontSize="sm" color="gray.300">
              {classInfo?.name || character.characterClass}
            </Text>
          </Flex>
        </Box>
      </Flex>

      <Divider mb={4} />

      {/* Health, Mana, and AC with Icons */}
      <Box mb={4}>
        <Tooltip label={`Health: ${currentHealth}/${maxHealth}`}>
          <Box>
            <Flex justify="space-between" mb={1} align="center">
              <Flex align="center">
                <Icon as={FaHeart} color="red.400" mr={2} />
                <Text fontSize="sm">Health</Text>
              </Flex>
              <Text fontSize="sm">{currentHealth}/{maxHealth}</Text>
            </Flex>
            <Progress 
              value={(currentHealth / maxHealth) * 100} 
              colorScheme="red" 
              size="sm" 
              mb={3}
              borderRadius="md"
              bg="whiteAlpha.200"
            />
          </Box>
        </Tooltip>

        <Tooltip label={`Mana: ${currentMana}/${maxMana}`}>
          <Box>
            <Flex justify="space-between" mb={1} align="center">
              <Flex align="center">
                <Icon as={FaFlask} color="blue.400" mr={2} />
                <Text fontSize="sm">Mana</Text>
              </Flex>
              <Text fontSize="sm">{currentMana}/{maxMana}</Text>
            </Flex>
            <Progress 
              value={(currentMana / maxMana) * 100} 
              colorScheme="blue" 
              size="sm" 
              mb={3}
              borderRadius="md"
              bg="whiteAlpha.200"
            />
          </Box>
        </Tooltip>

        <Tooltip label={`Experience: ${currentXP}/${nextLevelXP}`}>
          <Box>
            <Flex justify="space-between" mb={1} align="center">
              <Flex align="center">
                <Icon as={FaStar} color="purple.400" mr={2} />
                <Text fontSize="sm">Experience</Text>
              </Flex>
              <Text fontSize="sm">{currentXP}/{nextLevelXP}</Text>
            </Flex>
            <Progress 
              value={(currentXP / nextLevelXP) * 100} 
              colorScheme="purple" 
              size="sm"
              mb={3}
              borderRadius="md"
              bg="whiteAlpha.200"
            />
          </Box>
        </Tooltip>
        
        <Flex justify="space-between" mt={2} align="center">
          <Flex align="center">
            <Icon as={FaShieldAlt} color="green.400" mr={2} />
            <Text fontSize="sm" fontWeight="bold">Armor Class</Text>
          </Flex>
          <Badge 
            fontSize="sm" 
            px={2} 
            py={0.5} 
            borderRadius="full"
            bg="green.700"
            color="white"
          >
            {armorClass}
          </Badge>
        </Flex>
      </Box>

      <Divider mb={4} />

      {/* Stats with improved layout */}
      <Box mb={4}>
        <Text fontSize="md" fontWeight="bold" mb={3}>
          Attributes
        </Text>
        
        <Grid templateColumns="repeat(2, 1fr)" gap={3}>
          {Object.entries(character.stats).map(([statName, value]) => {
            const modifier = calculateModifier(value);
            const formattedModifier = formatModifier(modifier);
            const modifierColor = modifier >= 0 ? "green.300" : "red.300";
            
            // Check if this is a primary ability for the class
            const isPrimary = classInfo?.primaryAbility.toLowerCase().includes(statName.toLowerCase());
            
            // Choose icon based on stat
            let StatIcon;
            switch(statName.toLowerCase()) {
              case 'strength': StatIcon = FaRunning; break;
              case 'dexterity': StatIcon = FaRunning; break;
              case 'constitution': StatIcon = FaHeart; break;
              case 'intelligence': StatIcon = FaBrain; break;
              case 'wisdom': StatIcon = FaBrain; break;
              case 'charisma': StatIcon = FaStar; break;
              default: StatIcon = FaStar;
            }
            
            return (
              <GridItem 
                key={statName}
                p={2}
                borderRadius="md"
                bg={isPrimary ? `${classColors.secondary}30` : "whiteAlpha.100"}
                borderLeft={isPrimary ? "3px solid" : "none"}
                borderColor={classColors.primary}
              >
                <Flex justify="space-between" align="center">
                  <Flex align="center">
                    <Icon as={StatIcon} color={isPrimary ? classColors.primary : "gray.400"} mr={2} />
                    <Text textTransform="capitalize" fontSize="sm">{statName}</Text>
                  </Flex>
                  <Flex align="center">
                    <Text fontWeight="bold" mr={1}>{value}</Text>
                    <Text fontSize="xs" color={modifierColor}>
                      ({formattedModifier})
                    </Text>
                  </Flex>
                </Flex>
              </GridItem>
            );
          })}
        </Grid>
      </Box>

      <Divider mb={4} />

      {/* Class Abilities with improved styling */}
      <Box mb={4}>
        <Text fontSize="md" fontWeight="bold" mb={2}>
          Class Abilities
        </Text>
        <Flex flexWrap="wrap" gap={2}>
          {character.abilities.map((ability, index) => (
            <Badge 
              key={index} 
              px={2} 
              py={1}
              borderRadius="full"
              bg={`${classColors.secondary}90`}
              color="white"
            >
              {ability}
            </Badge>
          ))}
        </Flex>
      </Box>

      {/* Equipment Section */}
      <Box mb={4}>
        <Text fontSize="md" fontWeight="bold" mb={2}>
          Equipment
        </Text>
        <Grid templateColumns="repeat(2, 1fr)" gap={2}>
          {character.equipment ? (
            Object.entries(character.equipment).map(([slot, item], index) => (
              <GridItem 
                key={index}
                p={2}
                borderRadius="md"
                bg="whiteAlpha.100"
              >
                <Text fontSize="xs" color="gray.400" textTransform="capitalize">{slot}</Text>
                <Text fontSize="sm">{item?.name || "Empty"}</Text>
              </GridItem>
            ))
          ) : (
            <Text fontSize="sm" color="gray.400">No equipment</Text>
          )}
        </Grid>
      </Box>

      {/* Inventory Section */}
      <Box>
        <Text fontSize="md" fontWeight="bold" mb={2}>
          Inventory
        </Text>
        <Flex justify="space-between" mb={2}>
          <Badge colorScheme="yellow">Gold: {character.gold || 0}</Badge>
          <Badge colorScheme="red">Potions: {character.potions || 0}</Badge>
        </Flex>
        <Box 
          p={2} 
          borderRadius="md" 
          bg="whiteAlpha.100" 
          height="100px"
          overflowY="auto"
        >
          {character.inventory && character.inventory.length > 0 ? (
            character.inventory.map((item, index) => (
              <Text key={index} fontSize="sm" mb={1}>{item.name}</Text>
            ))
          ) : (
            <Text fontSize="sm" color="gray.400">Inventory empty</Text>
          )}
        </Box>
      </Box>
    </Box>
  );
}; 