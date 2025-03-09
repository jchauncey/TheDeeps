import { Box, Flex, Text, Progress, Icon, Button, Tooltip, Divider } from '@chakra-ui/react';
import { CharacterData, CHARACTER_CLASSES } from '../../types/game';
import { FaHeart, FaFlask, FaStar, FaUser, FaDungeon } from 'react-icons/fa';
import { useState, useEffect } from 'react';
import { CharacterProfileModal } from '../../components/game/CharacterProfileModal';
import { CLASS_COLORS } from '../../constants/gameConstants';
// Import Game Icons for class-specific icons
import { 
  GiSwordman, 
  GiWizardFace, 
  GiDaggers,
  GiHolySymbol, 
  GiOakLeaf,
  GiDevilMask,
  GiLyre,
  GiShield,
  GiBowArrow,
  GiMonkFace,
  GiAxeSwing,
  GiFireball,
  GiDungeonGate
} from 'react-icons/gi';
import { OPEN_CHARACTER_PROFILE_EVENT } from './GameControls';

interface GameStatusSimpleProps {
  character: CharacterData | null;
  currentFloor?: number;
  dungeonName?: string;
}

export const GameStatusSimple = ({ character, currentFloor = 1, dungeonName = 'The Deeps' }: GameStatusSimpleProps) => {
  const [isProfileOpen, setIsProfileOpen] = useState(false);

  // Listen for the custom event to open the profile
  useEffect(() => {
    const handleOpenProfile = () => {
      if (character) {
        setIsProfileOpen(true);
      }
    };

    // Add event listener for the custom event
    window.addEventListener(OPEN_CHARACTER_PROFILE_EVENT, handleOpenProfile);

    // Clean up
    return () => {
      window.removeEventListener(OPEN_CHARACTER_PROFILE_EVENT, handleOpenProfile);
    };
  }, [character]);

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
  
  // Get the appropriate icon component based on class
  const getClassIcon = () => {
    switch(classColors.iconType) {
      case 'GiSwordman': return GiSwordman;
      case 'GiWizardFace': return GiWizardFace;
      case 'GiDaggers': return GiDaggers;
      case 'GiHolySymbol': return GiHolySymbol;
      case 'GiOakLeaf': return GiOakLeaf;
      case 'GiDevilMask': return GiDevilMask;
      case 'GiLyre': return GiLyre;
      case 'GiShield': return GiShield;
      case 'GiBowArrow': return GiBowArrow;
      case 'GiMonkFace': return GiMonkFace;
      case 'GiAxeSwing': return GiAxeSwing;
      case 'GiFireball': return GiFireball;
      case 'GiDungeonGate': return GiDungeonGate;
      default: return FaUser;
    }
  };
  
  const ClassIcon = getClassIcon();
  
  // Calculate derived stats
  const calculateModifier = (stat: number) => Math.floor((stat - 10) / 2);
  
  // Calculate health, mana, and AC
  const constitutionMod = calculateModifier(character.stats.constitution);
  const intelligenceMod = calculateModifier(character.stats.intelligence);
  
  // Get character level with default value of 1
  const characterLevel = character.level || 1;
  
  const maxHealth = classInfo ? 
    classInfo.hitDie + constitutionMod + (characterLevel - 1) * (Math.floor(classInfo.hitDie / 2) + 1 + constitutionMod) : 
    10 + constitutionMod;
  
  const maxMana = classInfo?.primaryAbility.toLowerCase().includes('intelligence') ? 
    10 + intelligenceMod * 2 + (characterLevel - 1) * (4 + intelligenceMod) : 
    5 + intelligenceMod + (characterLevel - 1) * (2 + Math.floor(intelligenceMod / 2));
  
  // Current values (in a real game, these would come from the character state)
  const currentHealth = character.health || maxHealth;
  const currentMana = character.mana || maxMana;
  const currentXP = character.experience || 0;
  const nextLevelXP = characterLevel * 100;

  return (
    <>
      <Box
        height="100%"
        width="100%"
        bg="rgba(0, 0, 0, 0.7)"
        p={4}
        borderRadius="md"
        color="white"
        borderLeft="4px solid"
        borderColor={classColors.primary}
        display="flex"
        flexDirection="column"
        justifyContent="flex-start"
      >
        {/* Character Name and Level */}
        <Flex justify="space-between" align="center" mb={4}>
          <Flex align="center">
            <Icon as={ClassIcon} color={classColors.primary} mr={2} boxSize="20px" />
            <Text fontSize="lg" fontWeight="bold">
              {character.name}
            </Text>
            <Text fontSize="sm" ml={2} color="gray.300">
              Lvl {characterLevel}
            </Text>
          </Flex>
          <Tooltip label="Open Character Profile (Press 'C')">
            <Button 
              size="sm" 
              variant="outline"
              onClick={() => setIsProfileOpen(true)}
              _hover={{ bg: classColors.primary + '30' }}
              borderColor={classColors.primary}
              color={classColors.primary}
            >
              Profile
            </Button>
          </Tooltip>
        </Flex>

        {/* Dungeon Information */}
        <Box mb={3} p={2} bg="rgba(0, 0, 0, 0.3)" borderRadius="md">
          <Flex align="center" mb={1}>
            <Icon as={FaDungeon} color="purple.400" mr={2} />
            <Text fontSize="sm" fontWeight="bold">{dungeonName}</Text>
          </Flex>
          <Flex align="center">
            <Text fontSize="sm" color="gray.300">Floor {currentFloor}</Text>
          </Flex>
        </Box>

        <Divider mb={3} borderColor="whiteAlpha.300" />

        {/* Health, Mana, and XP with Icons */}
        <Box>
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
        </Box>
      </Box>

      {/* Character Profile Modal */}
      <CharacterProfileModal 
        character={character} 
        isOpen={isProfileOpen} 
        onClose={() => setIsProfileOpen(false)} 
        currentFloor={currentFloor}
        dungeonName={dungeonName}
      />
    </>
  );
}; 