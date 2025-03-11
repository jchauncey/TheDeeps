import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Container,
  FormControl,
  FormLabel,
  Heading,
  Input,
  Select,
  SimpleGrid,
  Text,
  useToast,
  Divider,
  Card,
  CardBody,
  Stack,
  Badge,
  Flex,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  Tooltip,
  Progress,
} from '@chakra-ui/react';
import { ArrowBackIcon, InfoIcon, AddIcon, MinusIcon } from '@chakra-ui/icons';
import { createCharacter } from '../services/api';
import { CharacterClass, Attributes } from '../types';

interface ClassInfo {
  name: CharacterClass;
  description: string;
  primaryAttributes: string[];
}

// Define class information
const classInfo: ClassInfo[] = [
  {
    name: 'barbarian',
    description: 'A fierce warrior who can enter a powerful rage in battle.',
    primaryAttributes: ['Strength', 'Constitution'],
  },
  {
    name: 'bard',
    description: 'A versatile character who uses music and performance to cast spells.',
    primaryAttributes: ['Charisma', 'Dexterity'],
  },
  {
    name: 'cleric',
    description: 'A divine spellcaster who channels the power of their deity.',
    primaryAttributes: ['Wisdom', 'Charisma'],
  },
  {
    name: 'druid',
    description: 'A nature-focused spellcaster with shapeshifting abilities.',
    primaryAttributes: ['Wisdom', 'Constitution'],
  },
  {
    name: 'mage',
    description: 'A powerful spellcaster who harnesses arcane magic through intelligence.',
    primaryAttributes: ['Intelligence', 'Wisdom'],
  },
  {
    name: 'monk',
    description: 'A martial artist who harnesses the power of their body as a weapon.',
    primaryAttributes: ['Dexterity', 'Wisdom'],
  },
  {
    name: 'paladin',
    description: 'A holy warrior who combines martial prowess with divine magic.',
    primaryAttributes: ['Strength', 'Charisma'],
  },
  {
    name: 'ranger',
    description: 'A skilled hunter and wilderness expert with combat and tracking abilities.',
    primaryAttributes: ['Dexterity', 'Wisdom'],
  },
  {
    name: 'rogue',
    description: 'A stealthy character who excels at precision attacks and evasion.',
    primaryAttributes: ['Dexterity', 'Charisma'],
  },
  {
    name: 'sorcerer',
    description: 'A spellcaster with innate magical abilities from their bloodline.',
    primaryAttributes: ['Charisma', 'Constitution'],
  },
  {
    name: 'warlock',
    description: 'A spellcaster who derives power from a pact with an otherworldly entity.',
    primaryAttributes: ['Charisma', 'Constitution'],
  },
  {
    name: 'warrior',
    description: 'A skilled fighter and weapon master with high strength and durability.',
    primaryAttributes: ['Strength', 'Constitution'],
  },
];

// Constants for attribute allocation
const BASE_ATTRIBUTE_VALUE = 8;
const TOTAL_ATTRIBUTE_POINTS = 27;
const MIN_ATTRIBUTE_VALUE = 8;
const MAX_ATTRIBUTE_VALUE = 15;

const CharacterCreation: React.FC = () => {
  const [name, setName] = useState<string>('');
  const [selectedClass, setSelectedClass] = useState<CharacterClass>('warrior');
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [attributes, setAttributes] = useState<Attributes>({
    strength: BASE_ATTRIBUTE_VALUE,
    dexterity: BASE_ATTRIBUTE_VALUE,
    constitution: BASE_ATTRIBUTE_VALUE,
    intelligence: BASE_ATTRIBUTE_VALUE,
    wisdom: BASE_ATTRIBUTE_VALUE,
    charisma: BASE_ATTRIBUTE_VALUE,
  });
  const [pointsRemaining, setPointsRemaining] = useState<number>(TOTAL_ATTRIBUTE_POINTS);
  
  const navigate = useNavigate();
  const toast = useToast();

  // Calculate the cost of increasing an attribute
  const getAttributeCost = (currentValue: number): number => {
    if (currentValue < 13) return 1;
    if (currentValue < 14) return 2;
    return 3;
  };

  // Calculate the refund when decreasing an attribute
  const getAttributeRefund = (currentValue: number): number => {
    if (currentValue <= 13) return 1;
    if (currentValue <= 14) return 2;
    return 3;
  };

  // Handle attribute change
  const handleAttributeChange = (attr: keyof Attributes, newValue: number) => {
    const oldValue = attributes[attr];
    
    // Don't allow values outside the min/max range
    if (newValue < MIN_ATTRIBUTE_VALUE || newValue > MAX_ATTRIBUTE_VALUE) {
      return;
    }
    
    // Calculate points cost/refund
    let pointsDelta = 0;
    if (newValue > oldValue) {
      // Increasing attribute
      for (let i = oldValue; i < newValue; i++) {
        pointsDelta -= getAttributeCost(i);
      }
    } else if (newValue < oldValue) {
      // Decreasing attribute
      for (let i = oldValue; i > newValue; i--) {
        pointsDelta += getAttributeRefund(i - 1);
      }
    } else {
      // No change
      return;
    }
    
    // Check if we have enough points
    if (pointsRemaining + pointsDelta < 0) {
      toast({
        title: 'Not enough points',
        description: 'You do not have enough attribute points remaining',
        status: 'error',
        duration: 2000,
        isClosable: true,
      });
      return;
    }
    
    // Update attributes and remaining points
    setAttributes(prev => ({
      ...prev,
      [attr]: newValue
    }));
    setPointsRemaining(prev => prev + pointsDelta);
  };

  // Reset attributes to base values
  const resetAttributes = () => {
    setAttributes({
      strength: BASE_ATTRIBUTE_VALUE,
      dexterity: BASE_ATTRIBUTE_VALUE,
      constitution: BASE_ATTRIBUTE_VALUE,
      intelligence: BASE_ATTRIBUTE_VALUE,
      wisdom: BASE_ATTRIBUTE_VALUE,
      charisma: BASE_ATTRIBUTE_VALUE,
    });
    setPointsRemaining(TOTAL_ATTRIBUTE_POINTS);
  };

  // Pre-allocate points based on class selection
  useEffect(() => {
    resetAttributes();
    
    // Get the selected class info
    const selectedClassInfo = getSelectedClassInfo();
    
    // Create a new attributes object with base values
    const newAttributes = {
      strength: BASE_ATTRIBUTE_VALUE,
      dexterity: BASE_ATTRIBUTE_VALUE,
      constitution: BASE_ATTRIBUTE_VALUE,
      intelligence: BASE_ATTRIBUTE_VALUE,
      wisdom: BASE_ATTRIBUTE_VALUE,
      charisma: BASE_ATTRIBUTE_VALUE,
    };
    
    // Allocate points to primary attributes (2 points each)
    let remainingPoints = TOTAL_ATTRIBUTE_POINTS;
    selectedClassInfo.primaryAttributes.forEach(attr => {
      const attrKey = attr.toLowerCase() as keyof Attributes;
      newAttributes[attrKey] += 2;
      remainingPoints -= 2;
    });
    
    setAttributes(newAttributes);
    setPointsRemaining(remainingPoints);
  }, [selectedClass]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!name.trim()) {
      toast({
        title: 'Error',
        description: 'Please enter a character name',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    try {
      setIsSubmitting(true);
      await createCharacter({
        name: name.trim(),
        class: selectedClass,
        attributes: attributes,
      });
      
      toast({
        title: 'Character created',
        description: `${name} has been created successfully!`,
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      
      navigate('/');
    } catch (err) {
      console.error('Failed to create character:', err);
      toast({
        title: 'Error',
        description: 'Failed to create character. Please try again.',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const getSelectedClassInfo = (): ClassInfo => {
    return classInfo.find(c => c.name === selectedClass) || classInfo[0];
  };

  // Get attribute modifier (D&D style)
  const getAttributeModifier = (value: number): number => {
    return Math.floor((value - 10) / 2);
  };

  // Format modifier as string with + or -
  const formatModifier = (mod: number): string => {
    return mod >= 0 ? `+${mod}` : `${mod}`;
  };

  // Get color for attribute based on value
  const getAttributeColor = (value: number): string => {
    if (value >= 15) return "green.300";
    if (value >= 13) return "teal.300";
    if (value >= 11) return "blue.300";
    if (value == 10) return "white";
    if (value <= 8) return "red.400";
    return "orange.300";
  };

  // Sort classes alphabetically for the dropdown
  const sortedClasses = [...classInfo].sort((a, b) => 
    a.name.localeCompare(b.name)
  );

  return (
    <Container maxW="container.xl" py={8}>
      <Button 
        leftIcon={<ArrowBackIcon />} 
        variant="solid" 
        mb={8} 
        onClick={() => navigate('/')}
        bg="gray.600"
        color="cyan.300"
        borderColor="cyan.500"
        borderWidth="1px"
        _hover={{ bg: "gray.700", color: "cyan.200" }}
        _active={{ bg: "gray.800", color: "cyan.100" }}
        boxShadow="sm"
      >
        Back to Character Selection
      </Button>

      <Heading as="h1" size="xl" mb={6} textAlign="center" color="blue.300">
        Create New Character
      </Heading>

      <form onSubmit={handleSubmit}>
        <SimpleGrid columns={{ base: 1, lg: 2 }} spacing={10}>
          <Stack spacing={6}>
            <Card bg="gray.700" shadow="md">
              <CardBody>
                <Heading size="md" mb={4} color="cyan.300">Basic Information</Heading>
                <FormControl isRequired mb={4}>
                  <FormLabel color="gray.200">Character Name</FormLabel>
                  <Input 
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Enter character name"
                    bg="gray.600"
                    _hover={{ bg: "gray.600" }}
                    _focus={{ bg: "gray.600", borderColor: "blue.400" }}
                    color="white"
                  />
                </FormControl>

                <FormControl isRequired>
                  <FormLabel color="gray.200">Character Class</FormLabel>
                  <Select 
                    value={selectedClass}
                    onChange={(e) => setSelectedClass(e.target.value as CharacterClass)}
                    bg="gray.600"
                    _hover={{ bg: "gray.600" }}
                    _focus={{ bg: "gray.600", borderColor: "blue.400" }}
                    color="white"
                  >
                    {sortedClasses.map((c) => (
                      <option key={c.name} value={c.name} style={{backgroundColor: '#2D3748', color: 'white'}}>
                        {c.name.charAt(0).toUpperCase() + c.name.slice(1)}
                      </option>
                    ))}
                  </Select>
                </FormControl>
              </CardBody>
            </Card>

            <Card bg="gray.700" shadow="md">
              <CardBody>
                <Flex justify="space-between" align="center" mb={4}>
                  <Heading size="md" color="cyan.300">Attribute Points</Heading>
                  <Text color="gray.200">
                    Remaining: <Text as="span" fontWeight="bold" color={pointsRemaining > 0 ? "green.300" : "gray.400"}>
                      {pointsRemaining}
                    </Text>
                  </Text>
                </Flex>
                <Progress 
                  value={(pointsRemaining / TOTAL_ATTRIBUTE_POINTS) * 100} 
                  colorScheme="green" 
                  size="sm" 
                  mb={4}
                  borderRadius="md"
                />
                <Text fontSize="sm" mb={4} color="gray.300">
                  Allocate your attribute points. Your class's primary attributes have been pre-allocated.
                  <Tooltip label="Higher attributes cost more points. 13-14 costs 2 points per level, 15+ costs 3 points per level.">
                    <InfoIcon ml={2} color="blue.300" cursor="pointer" />
                  </Tooltip>
                </Text>
                <Button 
                  size="sm" 
                  onClick={resetAttributes} 
                  mb={4} 
                  colorScheme="orange" 
                  variant="solid"
                  leftIcon={<MinusIcon />}
                >
                  Reset Attributes
                </Button>
                
                <SimpleGrid columns={1} spacing={4}>
                  {Object.entries(attributes).map(([key, value]) => {
                    const attrKey = key as keyof Attributes;
                    const isPrimary = getSelectedClassInfo().primaryAttributes
                      .map(a => a.toLowerCase())
                      .includes(key.toLowerCase());
                    const modifier = getAttributeModifier(value);
                    
                    return (
                      <Flex 
                        key={key} 
                        justify="space-between" 
                        align="center" 
                        p={3} 
                        borderWidth="1px" 
                        borderRadius="md" 
                        borderColor={isPrimary ? "blue.400" : "gray.500"}
                        bg={isPrimary ? "blue.800" : "gray.700"}
                        _hover={{ bg: isPrimary ? "blue.700" : "gray.650" }}
                        transition="all 0.2s"
                        boxShadow="sm"
                      >
                        <Stat size="sm">
                          <StatLabel textTransform="capitalize" fontSize="md" fontWeight="bold" color="gray.100">
                            {key} {isPrimary && <Badge colorScheme="blue" ml={1}>Primary</Badge>}
                          </StatLabel>
                          <Flex align="center">
                            <StatNumber color={getAttributeColor(value)} fontSize="xl">{value}</StatNumber>
                            <StatHelpText ml={2} mb={0} color={getAttributeColor(modifier)}>
                              ({formatModifier(modifier)})
                            </StatHelpText>
                          </Flex>
                        </Stat>
                        <Flex>
                          <Button 
                            size="sm" 
                            onClick={() => handleAttributeChange(attrKey, value - 1)}
                            isDisabled={value <= MIN_ATTRIBUTE_VALUE}
                            mr={1}
                            colorScheme="red"
                            variant="solid"
                            bg="red.500"
                            _hover={{ bg: "red.600" }}
                            _active={{ bg: "red.700" }}
                          >
                            <MinusIcon />
                          </Button>
                          <Button 
                            size="sm" 
                            onClick={() => handleAttributeChange(attrKey, value + 1)}
                            isDisabled={value >= MAX_ATTRIBUTE_VALUE || pointsRemaining < getAttributeCost(value)}
                            colorScheme="green"
                            variant="solid"
                            bg="green.500"
                            _hover={{ bg: "green.600" }}
                            _active={{ bg: "green.700" }}
                            position="relative"
                          >
                            <AddIcon />
                            {value < MAX_ATTRIBUTE_VALUE && (
                              <Badge 
                                position="absolute" 
                                top="-8px" 
                                right="-8px" 
                                colorScheme="yellow" 
                                fontSize="0.6em"
                                borderRadius="full"
                              >
                                {getAttributeCost(value)}
                              </Badge>
                            )}
                          </Button>
                        </Flex>
                      </Flex>
                    );
                  })}
                </SimpleGrid>
              </CardBody>
            </Card>

            <Button 
              type="submit" 
              colorScheme="cyan"
              size="lg" 
              width="full"
              isLoading={isSubmitting}
              bg="cyan.500"
              _hover={{ bg: "cyan.600" }}
              _active={{ bg: "cyan.700" }}
              boxShadow="md"
            >
              Create Character
            </Button>
          </Stack>

          <Stack spacing={6}>
            <Card bg="gray.700" shadow="md">
              <CardBody>
                <Heading size="md" mb={2} textTransform="capitalize" color="cyan.300">
                  {selectedClass}
                </Heading>
                <Text mb={4} color="gray.300">{getSelectedClassInfo().description}</Text>
                <Divider mb={4} />
                <Text fontWeight="bold" mb={2} color="gray.200">Primary Attributes:</Text>
                <Stack direction="row" mb={4}>
                  {getSelectedClassInfo().primaryAttributes.map((attr) => (
                    <Badge key={attr} colorScheme="blue" fontSize="0.8em" p={1}>
                      {attr}
                    </Badge>
                  ))}
                </Stack>
              </CardBody>
            </Card>

            <Card bg="gray.700" shadow="md">
              <CardBody>
                <Heading size="md" mb={4} color="cyan.300">Character Preview</Heading>
                <SimpleGrid columns={2} spacing={6}>
                  <Box bg="gray.600" p={4} borderRadius="md">
                    <Text fontWeight="bold" color="gray.300" mb={1}>HP</Text>
                    <Text fontSize="xl" fontWeight="bold" color="red.300">
                      {10 + getAttributeModifier(attributes.constitution)}
                    </Text>
                    <Text fontSize="sm" color="gray.400">
                      Base 10 + Constitution Modifier
                    </Text>
                  </Box>
                  
                  {['mage', 'sorcerer', 'warlock', 'cleric', 'druid', 'bard', 'paladin'].includes(selectedClass) && (
                    <Box bg="gray.600" p={4} borderRadius="md">
                      <Text fontWeight="bold" color="gray.300" mb={1}>Mana</Text>
                      <Text fontSize="xl" fontWeight="bold" color="blue.300">
                        {(() => {
                          switch (selectedClass) {
                            case 'mage':
                            case 'sorcerer':
                            case 'warlock':
                              return 10 + getAttributeModifier(attributes.intelligence);
                            case 'cleric':
                            case 'druid':
                              return 10 + getAttributeModifier(attributes.wisdom);
                            case 'bard':
                            case 'paladin':
                              return 5 + getAttributeModifier(attributes.charisma);
                            default:
                              return 0;
                          }
                        })()}
                      </Text>
                      <Text fontSize="sm" color="gray.400">
                        {(() => {
                          switch (selectedClass) {
                            case 'mage':
                            case 'sorcerer':
                            case 'warlock':
                              return 'Base 10 + Intelligence Modifier';
                            case 'cleric':
                            case 'druid':
                              return 'Base 10 + Wisdom Modifier';
                            case 'bard':
                            case 'paladin':
                              return 'Base 5 + Charisma Modifier';
                            default:
                              return '';
                          }
                        })()}
                      </Text>
                    </Box>
                  )}
                  
                  <Box bg="gray.600" p={4} borderRadius="md">
                    <Text fontWeight="bold" color="gray.300" mb={1}>Attack</Text>
                    <Text fontSize="xl" fontWeight="bold" color="orange.300">
                      {(() => {
                        switch (selectedClass) {
                          case 'warrior':
                          case 'barbarian':
                          case 'paladin':
                            return `1d8 ${formatModifier(getAttributeModifier(attributes.strength))}`;
                          case 'rogue':
                          case 'ranger':
                          case 'monk':
                            return `1d6 ${formatModifier(getAttributeModifier(attributes.dexterity))}`;
                          case 'mage':
                          case 'sorcerer':
                            return `1d4 ${formatModifier(getAttributeModifier(attributes.intelligence))}`;
                          case 'cleric':
                          case 'druid':
                            return `1d6 ${formatModifier(getAttributeModifier(attributes.wisdom))}`;
                          case 'warlock':
                          case 'bard':
                            return `1d6 ${formatModifier(getAttributeModifier(attributes.charisma))}`;
                          default:
                            return '1d4';
                        }
                      })()}
                    </Text>
                    <Text fontSize="sm" color="gray.400">
                      {(() => {
                        switch (selectedClass) {
                          case 'warrior':
                          case 'barbarian':
                          case 'paladin':
                            return 'Base 1d8 + Strength Modifier';
                          case 'rogue':
                          case 'ranger':
                          case 'monk':
                            return 'Base 1d6 + Dexterity Modifier';
                          case 'mage':
                          case 'sorcerer':
                            return 'Base 1d4 + Intelligence Modifier';
                          case 'cleric':
                          case 'druid':
                            return 'Base 1d6 + Wisdom Modifier';
                          case 'warlock':
                          case 'bard':
                            return 'Base 1d6 + Charisma Modifier';
                          default:
                            return 'Base 1d4';
                        }
                      })()}
                    </Text>
                  </Box>
                  
                  <Box bg="gray.600" p={4} borderRadius="md">
                    <Text fontWeight="bold" color="gray.300" mb={1}>Defense</Text>
                    <Text fontSize="xl" fontWeight="bold" color="cyan.300">
                      {10 + getAttributeModifier(attributes.dexterity)}
                    </Text>
                    <Text fontSize="sm" color="gray.400">
                      Base 10 + Dexterity Modifier
                    </Text>
                  </Box>
                </SimpleGrid>
              </CardBody>
            </Card>

            <Card bg="gray.700" shadow="md">
              <CardBody>
                <Heading size="md" mb={4} color="cyan.300">Starting Equipment</Heading>
                <SimpleGrid columns={1} spacing={3}>
                  {(() => {
                    switch (selectedClass) {
                      case 'warrior':
                        return (
                          <>
                            <Text color="gray.200">• Longsword <Text as="span" color="orange.300">(1d8 damage)</Text></Text>
                            <Text color="gray.200">• Chain mail armor <Text as="span" color="cyan.300">(AC 16)</Text></Text>
                            <Text color="gray.200">• Shield <Text as="span" color="cyan.300">(+2 AC)</Text></Text>
                            <Text color="gray.200">• Explorer's pack</Text>
                          </>
                        );
                      case 'mage':
                        return (
                          <>
                            <Text color="gray.200">• Quarterstaff <Text as="span" color="orange.300">(1d6 damage)</Text></Text>
                            <Text color="gray.200">• Spellbook</Text>
                            <Text color="gray.200">• Arcane focus</Text>
                            <Text color="gray.200">• Scholar's pack</Text>
                          </>
                        );
                      case 'rogue':
                        return (
                          <>
                            <Text color="gray.200">• Shortsword <Text as="span" color="orange.300">(1d6 damage)</Text></Text>
                            <Text color="gray.200">• Shortbow with 20 arrows <Text as="span" color="orange.300">(1d6 damage)</Text></Text>
                            <Text color="gray.200">• Leather armor <Text as="span" color="cyan.300">(AC 11)</Text></Text>
                            <Text color="gray.200">• Thieves' tools</Text>
                          </>
                        );
                      case 'cleric':
                        return (
                          <>
                            <Text color="gray.200">• Mace <Text as="span" color="orange.300">(1d6 damage)</Text></Text>
                            <Text color="gray.200">• Scale mail <Text as="span" color="cyan.300">(AC 14)</Text></Text>
                            <Text color="gray.200">• Shield <Text as="span" color="cyan.300">(+2 AC)</Text></Text>
                            <Text color="gray.200">• Holy symbol</Text>
                          </>
                        );
                      default:
                        return (
                          <>
                            <Text color="gray.200">• Basic weapon</Text>
                            <Text color="gray.200">• Basic armor</Text>
                            <Text color="gray.200">• Adventurer's pack</Text>
                            <Text color="gray.200">• 10 gold pieces</Text>
                          </>
                        );
                    }
                  })()}
                </SimpleGrid>
              </CardBody>
            </Card>
          </Stack>
        </SimpleGrid>
      </form>
    </Container>
  );
};

export default CharacterCreation; 