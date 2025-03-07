import { useEffect, useState } from 'react';
import { 
  Modal, 
  ModalOverlay, 
  ModalContent, 
  ModalHeader, 
  ModalBody, 
  ModalCloseButton,
  useDisclosure,
  Button,
  Box,
  Tooltip
} from '@chakra-ui/react';
import { sendWebSocketMessage } from '../../services/api';

export const GameControls = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [debugMode, setDebugMode] = useState(false);
  // We define the keyMap but don't need to update it
  const [keyMap] = useState<Record<string, string>>({
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
    'Escape': 'Menu',
    'ctrl+d': 'Toggle Debug Mode'
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
      } else if (e.ctrlKey && e.key === 'd') {
        // Toggle debug mode (F12 or Ctrl+D)
        e.preventDefault(); // Prevent browser's default behavior
        toggleDebugMode();
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
  };

  // Handle actions
  const handleAction = (action: string) => {
    sendWebSocketMessage({
      type: 'action',
      action
    });
  };

  // Toggle debug mode
  const toggleDebugMode = () => {
    setDebugMode(!debugMode);
    sendWebSocketMessage({
      type: 'action',
      action: 'toggle_debug'
    });
  };

  return (
    <>
      {/* Debug Button */}
      <Tooltip label="Toggle Debug Mode (F12 or Ctrl+D)" placement="left">
        <Box position="fixed" top="20px" right="20px" zIndex={1000}>
          <Button
            size="sm"
            colorScheme={debugMode ? "red" : "gray"}
            onClick={toggleDebugMode}
          >
            {debugMode ? "Debug: ON" : "Debug: OFF"}
          </Button>
        </Box>
      </Tooltip>

      {/* Help Modal */}
      <Modal isOpen={isOpen} onClose={onClose} size="lg">
        <ModalOverlay />
        <ModalContent bg="#291326" color="white">
          <ModalHeader>Game Controls</ModalHeader>
          <ModalCloseButton />
          <ModalBody pb={6}>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '1.5rem' }}>
              <div>
                <h3 style={{ fontSize: '1rem', fontWeight: 600, marginBottom: '0.75rem' }}>Movement</h3>
                <p style={{ fontSize: '0.875rem' }}>↑/k: Move up</p>
                <p style={{ fontSize: '0.875rem' }}>↓/j: Move down</p>
                <p style={{ fontSize: '0.875rem' }}>←/h: Move left</p>
                <p style={{ fontSize: '0.875rem' }}>→/l: Move right</p>
                <p style={{ fontSize: '0.875rem' }}>.: Wait</p>
              </div>
              
              <div>
                <h3 style={{ fontSize: '1rem', fontWeight: 600, marginBottom: '0.75rem' }}>Actions</h3>
                <p style={{ fontSize: '0.875rem' }}>g: Pick up item</p>
                <p style={{ fontSize: '0.875rem' }}>i: Inventory</p>
                <p style={{ fontSize: '0.875rem' }}>a: Attack</p>
                <p style={{ fontSize: '0.875rem' }}>u: Use item</p>
                <p style={{ fontSize: '0.875rem' }}>&lt;/&gt;: Stairs</p>
                <p style={{ fontSize: '0.875rem' }}>d: Toggle debug mode</p>
                <p style={{ fontSize: '0.875rem' }}>?: Help (this screen)</p>
                <p style={{ fontSize: '0.875rem' }}>Esc: Menu/Close</p>
              </div>
            </div>
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
}; 