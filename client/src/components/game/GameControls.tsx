import { useEffect, useState } from 'react';
import { 
  Box, 
  Text, 
  Grid, 
  Badge, 
  Modal, 
  ModalOverlay, 
  ModalContent, 
  ModalHeader, 
  ModalBody, 
  ModalCloseButton,
  useDisclosure
} from '@chakra-ui/react';
import { sendWebSocketMessage } from '../../services/api';

export const GameControls = () => {
  const [lastAction, setLastAction] = useState<string>('');
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [keyMap, setKeyMap] = useState<Record<string, string>>({
    'ArrowUp': 'Move Up',
    'ArrowDown': 'Move Down',
    'ArrowLeft': 'Move Left',
    'ArrowRight': 'Move Right',
    'k': 'Move Up',
    'j': 'Move Down',
    'h': 'Move Left',
    'l': 'Move Right',
    '.': 'Wait',
    'g': 'Pick Up',
    'i': 'Inventory',
    'a': 'Attack',
    'u': 'Use Item',
    '>': 'Descend Stairs',
    '<': 'Ascend Stairs',
    '?': 'Help',
    'Escape': 'Menu'
  });

  // Handle keyboard input
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Prevent default behavior for game keys
      if (Object.keys(keyMap).includes(e.key)) {
        e.preventDefault();
      }

      // Handle movement
      if (['ArrowUp', 'k'].includes(e.key)) {
        handleMove('up');
      } else if (['ArrowDown', 'j'].includes(e.key)) {
        handleMove('down');
      } else if (['ArrowLeft', 'h'].includes(e.key)) {
        handleMove('left');
      } else if (['ArrowRight', 'l'].includes(e.key)) {
        handleMove('right');
      }
      // Handle actions
      else if (e.key === '.') {
        handleAction('wait');
      } else if (e.key === 'g') {
        handleAction('pickup');
      } else if (e.key === 'i') {
        handleAction('inventory');
      } else if (e.key === 'a') {
        handleAction('attack');
      } else if (e.key === 'u') {
        handleAction('use');
      } else if (e.key === '>') {
        handleAction('descend');
      } else if (e.key === '<') {
        handleAction('ascend');
      } else if (e.key === '?') {
        // Toggle help modal
        if (isOpen) {
          onClose();
        } else {
          onOpen();
        }
      } else if (e.key === 'Escape') {
        // Close modal if open, otherwise show menu
        if (isOpen) {
          onClose();
        } else {
          handleAction('menu');
        }
      }
    };

    // Add event listener
    window.addEventListener('keydown', handleKeyDown);

    // Clean up
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [isOpen, onOpen, onClose, keyMap]);

  // Handle movement
  const handleMove = (direction: 'up' | 'down' | 'left' | 'right') => {
    sendWebSocketMessage({
      type: 'move',
      direction
    });
    setLastAction(`Move ${direction}`);
  };

  // Handle actions
  const handleAction = (action: string) => {
    sendWebSocketMessage({
      type: 'action',
      action
    });
    setLastAction(action);
  };

  return (
    <>
      {/* Last action indicator (small and unobtrusive) */}
      {lastAction && (
        <Box 
          position="absolute" 
          bottom={4} 
          left={4} 
          bg="rgba(0, 0, 0, 0.7)" 
          p={2} 
          borderRadius="md"
          zIndex={1}
        >
          <Text fontSize="sm" color="white">Last action: <Badge colorScheme="purple">{lastAction}</Badge></Text>
        </Box>
      )}
      
      {/* Help Modal */}
      <Modal isOpen={isOpen} onClose={onClose} size="lg">
        <ModalOverlay />
        <ModalContent bg="#291326" color="white">
          <ModalHeader>Game Controls</ModalHeader>
          <ModalCloseButton />
          <ModalBody pb={6}>
            <Grid templateColumns="repeat(2, 1fr)" gap={6}>
              <Box>
                <Text fontSize="md" fontWeight="semibold" mb={3}>Movement</Text>
                <Text fontSize="sm">↑/k: Move up</Text>
                <Text fontSize="sm">↓/j: Move down</Text>
                <Text fontSize="sm">←/h: Move left</Text>
                <Text fontSize="sm">→/l: Move right</Text>
                <Text fontSize="sm">.: Wait</Text>
              </Box>
              
              <Box>
                <Text fontSize="md" fontWeight="semibold" mb={3}>Actions</Text>
                <Text fontSize="sm">g: Pick up item</Text>
                <Text fontSize="sm">i: Inventory</Text>
                <Text fontSize="sm">a: Attack</Text>
                <Text fontSize="sm">u: Use item</Text>
                <Text fontSize="sm">&lt;/&gt;: Stairs</Text>
                <Text fontSize="sm">?: Help (this screen)</Text>
                <Text fontSize="sm">Esc: Menu/Close</Text>
              </Box>
            </Grid>
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
}; 