import React, { useState } from 'react';
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
} from '@chakra-ui/react';
import { ArrowBackIcon } from '@chakra-ui/icons';
import { createCharacter } from '../services/api';
import { CharacterClass } from '../types';

interface ClassInfo {
  name: CharacterClass;
  description: string;
  primaryAttributes: string[];
}

const classInfo: ClassInfo[] = [
  {
    name: 'warrior',
    description: 'A skilled fighter and weapon master with high strength and durability.',
    primaryAttributes: ['Strength', 'Constitution'],
  },
  {
    name: 'mage',
    description: 'A powerful spellcaster who harnesses arcane magic through intelligence.',
    primaryAttributes: ['Intelligence', 'Wisdom'],
  },
  {
    name: 'rogue',
    description: 'A stealthy character who excels at precision attacks and evasion.',
    primaryAttributes: ['Dexterity', 'Charisma'],
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
    name: 'warlock',
    description: 'A spellcaster who derives power from a pact with an otherworldly entity.',
    primaryAttributes: ['Charisma', 'Constitution'],
  },
  {
    name: 'bard',
    description: 'A versatile character who uses music and performance to cast spells.',
    primaryAttributes: ['Charisma', 'Dexterity'],
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
    name: 'monk',
    description: 'A martial artist who harnesses the power of their body as a weapon.',
    primaryAttributes: ['Dexterity', 'Wisdom'],
  },
  {
    name: 'barbarian',
    description: 'A fierce warrior who can enter a powerful rage in battle.',
    primaryAttributes: ['Strength', 'Constitution'],
  },
  {
    name: 'sorcerer',
    description: 'A spellcaster with innate magical abilities from their bloodline.',
    primaryAttributes: ['Charisma', 'Constitution'],
  },
];

const CharacterCreation: React.FC = () => {
  const [name, setName] = useState<string>('');
  const [selectedClass, setSelectedClass] = useState<CharacterClass>('warrior');
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const navigate = useNavigate();
  const toast = useToast();

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

  return (
    <Container maxW="container.xl" py={8}>
      <Button 
        leftIcon={<ArrowBackIcon />} 
        variant="outline" 
        mb={8} 
        onClick={() => navigate('/')}
      >
        Back to Character Selection
      </Button>

      <Heading as="h1" size="xl" mb={6} textAlign="center">
        Create New Character
      </Heading>

      <form onSubmit={handleSubmit}>
        <SimpleGrid columns={{ base: 1, md: 2 }} spacing={10}>
          <Box>
            <FormControl isRequired mb={6}>
              <FormLabel>Character Name</FormLabel>
              <Input 
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Enter character name"
                bg="gray.700"
              />
            </FormControl>

            <FormControl isRequired mb={6}>
              <FormLabel>Character Class</FormLabel>
              <Select 
                value={selectedClass}
                onChange={(e) => setSelectedClass(e.target.value as CharacterClass)}
                bg="gray.700"
              >
                {classInfo.map((c) => (
                  <option key={c.name} value={c.name}>
                    {c.name.charAt(0).toUpperCase() + c.name.slice(1)}
                  </option>
                ))}
              </Select>
            </FormControl>

            <Button 
              type="submit" 
              colorScheme="blue" 
              size="lg" 
              width="full"
              mt={4}
              isLoading={isSubmitting}
            >
              Create Character
            </Button>
          </Box>

          <Box>
            <Card bg="gray.800" mb={6}>
              <CardBody>
                <Heading size="md" mb={2} textTransform="capitalize">
                  {selectedClass}
                </Heading>
                <Text mb={4}>{getSelectedClassInfo().description}</Text>
                <Divider mb={4} />
                <Text fontWeight="bold" mb={2}>Primary Attributes:</Text>
                <Stack direction="row" mb={4}>
                  {getSelectedClassInfo().primaryAttributes.map((attr) => (
                    <Badge key={attr} colorScheme="blue" fontSize="0.8em" p={1}>
                      {attr}
                    </Badge>
                  ))}
                </Stack>
              </CardBody>
            </Card>
          </Box>
        </SimpleGrid>
      </form>
    </Container>
  );
};

export default CharacterCreation; 