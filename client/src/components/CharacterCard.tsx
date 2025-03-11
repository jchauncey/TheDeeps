import React from 'react';
import {
  Box,
  Heading,
  Text,
  Stack,
  Badge,
  Progress,
  IconButton,
  useDisclosure,
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  Button,
} from '@chakra-ui/react';
import { DeleteIcon } from '@chakra-ui/icons';
import { Character } from '../types';

interface CharacterCardProps {
  character: Character;
  onDelete: (id: string) => void;
  onSelect: (character: Character) => void;
}

const CharacterCard: React.FC<CharacterCardProps> = ({ character, onDelete, onSelect }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = React.useRef<HTMLButtonElement>(null);

  // Get class color
  const getClassColor = (characterClass: string): string => {
    const classColors: Record<string, string> = {
      warrior: 'red.500',
      mage: 'blue.500',
      rogue: 'purple.500',
      cleric: 'yellow.500',
      druid: 'green.500',
      warlock: 'pink.500',
      bard: 'teal.500',
      paladin: 'orange.500',
      ranger: 'cyan.500',
      monk: 'gray.500',
      barbarian: 'red.700',
      sorcerer: 'blue.700',
    };
    return classColors[characterClass] || 'gray.500';
  };

  const handleDelete = () => {
    onDelete(character.id);
    onClose();
  };

  return (
    <>
      <Box
        p={5}
        shadow="md"
        borderWidth="1px"
        borderRadius="lg"
        bg="gray.800"
        _hover={{ shadow: "lg", transform: "translateY(-2px)", transition: "all 0.2s" }}
        cursor="pointer"
        onClick={() => onSelect(character)}
        position="relative"
      >
        <IconButton
          aria-label="Delete character"
          icon={<DeleteIcon />}
          size="sm"
          colorScheme="red"
          variant="ghost"
          position="absolute"
          top={2}
          right={2}
          onClick={(e) => {
            e.stopPropagation();
            onOpen();
          }}
        />
        
        <Heading fontSize="xl" mb={2}>{character.name}</Heading>
        
        <Badge colorScheme={getClassColor(character.class)} mb={2} textTransform="capitalize">
          {character.class}
        </Badge>
        
        <Text fontSize="lg" fontWeight="bold" mb={1}>
          Level {character.level}
        </Text>
        
        <Stack spacing={2} mt={4}>
          <Box>
            <Text fontSize="sm" mb={1}>HP: {character.currentHp}/{character.maxHp}</Text>
            <Progress 
              value={(character.currentHp / character.maxHp) * 100} 
              colorScheme="red" 
              size="sm" 
              borderRadius="md"
            />
          </Box>
          
          {character.maxMana > 0 && (
            <Box>
              <Text fontSize="sm" mb={1}>Mana: {character.currentMana}/{character.maxMana}</Text>
              <Progress 
                value={(character.currentMana / character.maxMana) * 100} 
                colorScheme="blue" 
                size="sm" 
                borderRadius="md"
              />
            </Box>
          )}
          
          <Text fontSize="sm">XP: {character.experience}</Text>
        </Stack>
      </Box>

      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent bg="gray.800">
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Character
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to delete {character.name}? This action cannot be undone.
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              <Button colorScheme="red" onClick={handleDelete} ml={3}>
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
};

export default CharacterCard; 