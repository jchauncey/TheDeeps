import { 
  Modal, 
  ModalOverlay, 
  ModalContent, 
  ModalHeader, 
  ModalFooter, 
  ModalBody, 
  ModalCloseButton,
  Box,
  Flex,
  Text,
  Divider,
  Badge,
  Icon,
  Grid,
  GridItem,
  Button
} from '@chakra-ui/react';
import { CharacterData, CHARACTER_CLASSES } from '../../types/game';
import { FaHeart, FaFlask, FaBrain, FaShieldAlt, FaRunning, FaStar, FaUser } from 'react-icons/fa';

// Import CLASS_COLORS from the same place GameStatus uses it
import { CLASS_COLORS } from '../../constants/gameConstants';

interface CharacterProfileModalProps {
  character: CharacterData | null;
  isOpen: boolean;
  onClose: () => void;
}

export const CharacterProfileModal = ({ character, isOpen, onClose }: CharacterProfileModalProps) => {
  if (!character) return null;

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
    <Modal isOpen={isOpen} onClose={onClose} size="lg" scrollBehavior="inside">
      <ModalOverlay backdropFilter="blur(2px)" />
      <ModalContent bg="gray.900" color="white" borderLeft={`4px solid ${classColors.primary}`}>
        <ModalHeader>Character Profile</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          {/* Character Header */}
          <Flex align="center" mb={3}>
            <Box 
              bg={classColors.primary} 
              p={3} 
              borderRadius="md" 
              mr={3}
              display="flex"
              alignItems="center"
              justifyContent="center"
              width="60px"
              height="60px"
            >
              <Icon as={FaUser} color={classColors.textColor} boxSize="1.5em" />
            </Box>
            <Box>
              <Text fontSize="2xl" fontWeight="bold" mb={0}>
                {character.name}
              </Text>
              <Flex align="center">
                <Badge 
                  mr={2}
                  px={2}
                  py={0.5}
                  borderRadius="full"
                  bg={classColors.primary}
                  color={classColors.textColor}
                >
                  Level {characterLevel}
                </Badge>
                <Text fontSize="md" color="gray.300">
                  {classInfo?.name || character.characterClass}
                </Text>
              </Flex>
            </Box>
          </Flex>

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
                  color={classColors.textColor}
                >
                  {ability}
                </Badge>
              ))}
            </Flex>
          </Box>

          <Divider mb={4} />

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

          <Divider mb={4} />

          {/* Inventory Section */}
          <Box>
            <Text fontSize="md" fontWeight="bold" mb={2}>
              Inventory
            </Text>
            <Flex justify="space-between" mb={2}>
              <Badge bg={classColors.primary} color={classColors.textColor}>Gold: {character.gold || 0}</Badge>
              <Badge bg="red.500" color="white">Potions: {character.potions || 0}</Badge>
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
        </ModalBody>
        <ModalFooter>
          <Button colorScheme="blue" mr={3} onClick={onClose}>
            Close
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}; 