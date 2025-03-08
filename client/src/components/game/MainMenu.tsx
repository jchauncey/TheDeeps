import { 
  Modal, 
  ModalOverlay, 
  ModalContent, 
  ModalHeader, 
  ModalBody, 
  ModalCloseButton,
  Button,
  VStack,
  useToast,
  Text
} from '@chakra-ui/react';
import { CharacterData } from '../../types/game';
import { saveGame } from '../../services/api';

interface MainMenuProps {
  isOpen: boolean;
  onClose: () => void;
  onNewGame: () => void;
  onLoadGame: () => void;
  character: CharacterData | null;
}

export const MainMenu = ({ isOpen, onClose, onNewGame, onLoadGame, character }: MainMenuProps) => {
  const toast = useToast();

  const handleSaveGame = async () => {
    if (!character) {
      toast({
        title: "Error",
        description: "No character data to save",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    // Use the saveGame function from the API
    const result = await saveGame(character.name);

    if (result.success) {
      toast({
        title: "Game Saved",
        description: `${character.name}'s progress has been saved.`,
        status: "success",
        duration: 3000,
        isClosable: true,
      });
    } else {
      toast({
        title: "Save Failed",
        description: result.message || "Could not save game. Please try again.",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
    }
    
    onClose();
  };

  const handleNewGame = () => {
    onNewGame();
    onClose();
  };

  const handleLoadGame = () => {
    onLoadGame();
    onClose();
  };

  const handleReturnToGame = () => {
    onClose();
  };

  const handleQuitGame = () => {
    // Return to start screen
    onNewGame(); // This will take us to character creation, which has a back button to start screen
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay backdropFilter="blur(2px)" />
      <ModalContent bg="#291326" color="white" borderLeft="4px solid #6B46C1">
        <ModalHeader>Game Menu</ModalHeader>
        <ModalCloseButton />
        <ModalBody pb={6}>
          <VStack spacing={4} align="stretch">
            <Button 
              onClick={handleReturnToGame}
              bg="#6B46C1"
              color="white"
              _hover={{ bg: "#805AD5" }}
            >
              Return to Game
            </Button>
            
            <Button 
              onClick={handleSaveGame}
              bg="transparent"
              color="white"
              border="1px solid"
              borderColor="#6B46C1"
              _hover={{ bg: "rgba(107, 70, 193, 0.2)" }}
            >
              Save Game
            </Button>
            
            <Button 
              onClick={handleLoadGame}
              bg="transparent"
              color="white"
              border="1px solid"
              borderColor="#6B46C1"
              _hover={{ bg: "rgba(107, 70, 193, 0.2)" }}
            >
              Load Game
            </Button>
            
            <Button 
              onClick={handleNewGame}
              bg="transparent"
              color="white"
              border="1px solid"
              borderColor="#6B46C1"
              _hover={{ bg: "rgba(107, 70, 193, 0.2)" }}
            >
              New Game
            </Button>
            
            <Button 
              onClick={handleQuitGame}
              bg="transparent"
              color="red.400"
              border="1px solid"
              borderColor="red.400"
              _hover={{ bg: "rgba(245, 101, 101, 0.2)" }}
              mt={2}
            >
              Quit to Main Menu
            </Button>
            
            {character && (
              <Text fontSize="sm" color="gray.400" textAlign="center" mt={2}>
                Playing as {character.name}, Level {character.level || 1} {character.characterClass.charAt(0).toUpperCase() + character.characterClass.slice(1)}
              </Text>
            )}
          </VStack>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
}; 