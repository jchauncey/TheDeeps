import React from 'react';
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalCloseButton,
  Box,
  Flex,
  Text,
  Divider,
  Grid,
  GridItem,
} from '@chakra-ui/react';

// Import the color constants from GameBoard
import { ENTITY_COLORS, TILE_COLORS, ITEM_COLORS, DIFFICULTY_COLORS } from './GameBoard';

interface MapLegendProps {
  isOpen: boolean;
  onClose: () => void;
}

export const MapLegend: React.FC<MapLegendProps> = ({ isOpen, onClose }) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose} size="lg">
      <ModalOverlay />
      <ModalContent bg="gray.800" color="white">
        <ModalHeader>Map Legend</ModalHeader>
        <ModalCloseButton />
        <ModalBody pb={6}>
          <Text mb={4}>
            This legend shows all the symbols and colors used in the game map.
          </Text>

          {/* Tiles Section */}
          <Text fontWeight="bold" fontSize="lg" mb={2}>
            Tiles
          </Text>
          <Grid templateColumns="repeat(3, 1fr)" gap={4} mb={4}>
            {Object.entries(TILE_COLORS).map(([type, color]) => (
              <GridItem key={type}>
                <Flex alignItems="center">
                  <Box
                    w="20px"
                    h="20px"
                    bg={color as string}
                    mr={2}
                    border="1px solid"
                    borderColor="gray.500"
                  />
                  <Text>{formatName(type)}</Text>
                </Flex>
              </GridItem>
            ))}
          </Grid>

          <Divider my={4} />

          {/* Entities Section */}
          <Text fontWeight="bold" fontSize="lg" mb={2}>
            Entities
          </Text>
          <Grid templateColumns="repeat(3, 1fr)" gap={4} mb={4}>
            {Object.entries(ENTITY_COLORS).map(([type, color]) => (
              <GridItem key={type}>
                <Flex alignItems="center">
                  <Box
                    w="20px"
                    h="20px"
                    bg={color as string}
                    mr={2}
                    borderRadius="full"
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    border="1px solid"
                    borderColor="gray.500"
                  >
                    <Text fontSize="xs" fontWeight="bold" color="black">
                      {type.charAt(0).toUpperCase()}
                    </Text>
                  </Box>
                  <Text>{formatName(type)}</Text>
                </Flex>
              </GridItem>
            ))}
          </Grid>

          <Divider my={4} />

          {/* Items Section */}
          <Text fontWeight="bold" fontSize="lg" mb={2}>
            Items
          </Text>
          <Grid templateColumns="repeat(3, 1fr)" gap={4} mb={4}>
            {Object.entries(ITEM_COLORS).map(([type, color]) => (
              <GridItem key={type}>
                <Flex alignItems="center">
                  <Box
                    w="20px"
                    h="20px"
                    bg={color as string}
                    mr={2}
                    borderRadius="full"
                    border="1px solid"
                    borderColor="gray.500"
                  />
                  <Text>{formatName(type)}</Text>
                </Flex>
              </GridItem>
            ))}
          </Grid>

          <Divider my={4} />

          {/* Difficulty Section */}
          <Text fontWeight="bold" fontSize="lg" mb={2}>
            Monster Difficulty
          </Text>
          <Grid templateColumns="repeat(3, 1fr)" gap={4} mb={4}>
            {Object.entries(DIFFICULTY_COLORS).map(([difficulty, color]) => (
              <GridItem key={difficulty}>
                <Flex alignItems="center">
                  <Box
                    w="20px"
                    h="20px"
                    bg="gray.600"
                    mr={2}
                    borderRadius="full"
                    border="2px solid"
                    borderColor={color as string}
                  />
                  <Text>{formatName(difficulty)}</Text>
                </Flex>
              </GridItem>
            ))}
          </Grid>

          <Divider my={4} />

          {/* Controls Section */}
          <Text fontWeight="bold" fontSize="lg" mb={2}>
            Controls
          </Text>
          <Grid templateColumns="repeat(2, 1fr)" gap={4}>
            <GridItem>
              <Flex>
                <Text fontWeight="bold" minWidth="80px">Arrow Keys:</Text>
                <Text>Move</Text>
              </Flex>
            </GridItem>
            <GridItem>
              <Flex>
                <Text fontWeight="bold" minWidth="80px">Space:</Text>
                <Text>Attack/Interact</Text>
              </Flex>
            </GridItem>
            <GridItem>
              <Flex>
                <Text fontWeight="bold" minWidth="80px">P:</Text>
                <Text>Pick up item</Text>
              </Flex>
            </GridItem>
            <GridItem>
              <Flex>
                <Text fontWeight="bold" minWidth="80px">L:</Text>
                <Text>Toggle Legend</Text>
              </Flex>
            </GridItem>
          </Grid>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};

// Helper function to format names
const formatName = (name: string): string => {
  return name
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
}; 