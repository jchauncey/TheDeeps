import React, { useState } from 'react';
import {
  Box,
  Flex,
  Text,
  Select,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  FormControl,
  FormLabel,
  Button,
  Heading,
  Container,
  VStack,
  HStack,
  Divider,
} from '@chakra-ui/react';
import RoomRenderer from '../components/RoomRenderer';

const TestRoom: React.FC = () => {
  const [roomType, setRoomType] = useState<string>('entrance');
  const [width, setWidth] = useState<number>(20);
  const [height, setHeight] = useState<number>(20);
  const [roomWidth, setRoomWidth] = useState<number | undefined>(undefined);
  const [roomHeight, setRoomHeight] = useState<number | undefined>(undefined);
  const [key, setKey] = useState<number>(0); // Used to force re-render

  const handleApply = () => {
    // Force re-render of the RoomRenderer component
    setKey(prevKey => prevKey + 1);
  };

  return (
    <Container maxW="container.xl" py={8}>
      <Heading as="h1" mb={6}>Room Renderer Test</Heading>
      <Text mb={4}>
        This page allows you to test the rendering of different room types and configurations.
        Use the controls below to customize the test room.
      </Text>
      
      <Flex direction={{ base: 'column', md: 'row' }} gap={6}>
        <Box flex="1" p={4} borderWidth="1px" borderRadius="md">
          <Heading as="h2" size="md" mb={4}>Room Controls</Heading>
          
          <VStack spacing={4} align="stretch">
            <FormControl>
              <FormLabel>Room Type</FormLabel>
              <Select 
                value={roomType} 
                onChange={(e) => setRoomType(e.target.value)}
              >
                <option value="entrance">Entrance</option>
                <option value="standard">Standard</option>
                <option value="treasure">Treasure</option>
                <option value="boss">Boss</option>
                <option value="safe">Safe</option>
                <option value="shop">Shop</option>
              </Select>
            </FormControl>
            
            <FormControl>
              <FormLabel>Floor Width</FormLabel>
              <NumberInput 
                min={10} 
                max={100} 
                value={width}
                onChange={(_, val) => setWidth(val)}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
            </FormControl>
            
            <FormControl>
              <FormLabel>Floor Height</FormLabel>
              <NumberInput 
                min={10} 
                max={100} 
                value={height}
                onChange={(_, val) => setHeight(val)}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
            </FormControl>
            
            <FormControl>
              <FormLabel>Room Width (optional)</FormLabel>
              <NumberInput 
                min={5} 
                max={20} 
                value={roomWidth || ''}
                onChange={(valueString, valueNumber) => setRoomWidth(valueString === '' ? undefined : valueNumber)}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
            </FormControl>
            
            <FormControl>
              <FormLabel>Room Height (optional)</FormLabel>
              <NumberInput 
                min={5} 
                max={20} 
                value={roomHeight || ''}
                onChange={(valueString, valueNumber) => setRoomHeight(valueString === '' ? undefined : valueNumber)}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
            </FormControl>
            
            <Button colorScheme="blue" onClick={handleApply}>
              Apply Changes
            </Button>
          </VStack>
        </Box>
        
        <Box flex="2" p={4} borderWidth="1px" borderRadius="md">
          <Heading as="h2" size="md" mb={4}>Room Preview</Heading>
          <Box h="500px" overflowY="auto">
            <RoomRenderer 
              key={key}
              roomType={roomType}
              width={width}
              height={height}
              roomWidth={roomWidth}
              roomHeight={roomHeight}
            />
          </Box>
        </Box>
      </Flex>
      
      <Box mt={8}>
        <Heading as="h2" size="md" mb={4}>Legend</Heading>
        <Flex wrap="wrap" gap={4}>
          <HStack>
            <Box bg="#555" w="20px" h="20px" />
            <Text>Wall (#)</Text>
          </HStack>
          <HStack>
            <Box bg="#222" w="20px" h="20px" />
            <Text>Floor (.)</Text>
          </HStack>
          <HStack>
            <Box bg="#F00" w="20px" h="20px" />
            <Text>Down Stairs (&gt;)</Text>
          </HStack>
          <HStack>
            <Box bg="#00F" w="20px" h="20px" />
            <Text>Up Stairs (&lt;)</Text>
          </HStack>
          <HStack>
            <Box bg="#222" w="20px" h="20px" color="yellow" display="flex" alignItems="center" justifyContent="center">
              @
            </Box>
            <Text>Character</Text>
          </HStack>
          <HStack>
            <Box bg="#222" w="20px" h="20px" color="red" display="flex" alignItems="center" justifyContent="center">
              M
            </Box>
            <Text>Monster</Text>
          </HStack>
          <HStack>
            <Box bg="#222" w="20px" h="20px" color="green" display="flex" alignItems="center" justifyContent="center">
              i
            </Box>
            <Text>Item</Text>
          </HStack>
        </Flex>
      </Box>
    </Container>
  );
};

export default TestRoom; 