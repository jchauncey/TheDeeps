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
import { OPEN_CHARACTER_PROFILE_EVENT } from '../game';

export const GameControls = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  // We define the keyMap but don't need to update it
  const [keyMap] = useState<Record<string, string>>({
    'w': 'Move Up',
    's': 'Move Down',
    'a': 'Move Left',
    'd': 'Move Right',
    '.': 'Wait',
    'g': 'Pick Up',
    'i': 'Inventory',
    'f': 'Attack',
    'u': 'Use Item',
    '>': 'Descend Stairs',
    '<': 'Ascend Stairs',
    '?': 'Help',
    'Escape': 'Menu',
    'c': 'Character Profile',
    'ctrl+d': 'Toggle Debug Mode'
  });

  // Handle keyboard input
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Prevent default behavior for game keys
      if (Object.keys(keyMap).includes(e.key)) {
        e.preventDefault();
      }

      // Handle movement - only WASD keys
      if (['w', 'W'].includes(e.key)) {
        handleMove('up');
      } else if (['s', 'S'].includes(e.key)) {
        handleMove('down');
      } else if (['a', 'A'].includes(e.key)) {
        handleMove('left');
      } else if (['d', 'D'].includes(e.key)) {
        handleMove('right');
      }
      // Handle actions
      else if (e.key === '.') {
        handleAction('wait');
      } else if (e.key === 'g') {
        handleAction('pickup');
      } else if (e.key === 'i') {
        handleAction('inventory');
      } else if (e.key === 'f') {
        handleAction('attack');
      } else if (e.key === 'u') {
        handleAction('use');
      } else if (e.key === '>') {
        handleAction('descend');
      } else if (e.key === '<') {
        handleAction('ascend');
      } else if (['c', 'C'].includes(e.key)) {
        // Dispatch custom event to open character profile
        window.dispatchEvent(new Event(OPEN_CHARACTER_PROFILE_EVENT));
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
  }, [keyMap, isOpen, onOpen, onClose]);

  // Handle movement
  const handleMove = (direction: 'up' | 'down' | 'left' | 'right') => {
    console.log(`GameControls: Sending move command: ${direction}`);
    sendWebSocketMessage({
      type: 'move',
      direction
    });
  };

  // Handle actions
  const handleAction = (action: string) => {
    console.log(`GameControls: Sending action command: ${action}`);
    sendWebSocketMessage({
      type: 'action',
      action
    });
  };

  return (
    <>
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
                <p style={{ fontSize: '0.875rem' }}>W: Move up</p>
                <p style={{ fontSize: '0.875rem' }}>S: Move down</p>
                <p style={{ fontSize: '0.875rem' }}>A: Move left</p>
                <p style={{ fontSize: '0.875rem' }}>D: Move right</p>
                <p style={{ fontSize: '0.875rem' }}>.: Wait</p>
              </div>
              
              <div>
                <h3 style={{ fontSize: '1rem', fontWeight: 600, marginBottom: '0.75rem' }}>Actions</h3>
                <p style={{ fontSize: '0.875rem' }}>g: Pick up item</p>
                <p style={{ fontSize: '0.875rem' }}>i: Inventory</p>
                <p style={{ fontSize: '0.875rem' }}>f: Attack</p>
                <p style={{ fontSize: '0.875rem' }}>u: Use item</p>
                <p style={{ fontSize: '0.875rem' }}>&lt;/&gt;: Stairs</p>
                <p style={{ fontSize: '0.875rem' }}>c: Character Profile</p>
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