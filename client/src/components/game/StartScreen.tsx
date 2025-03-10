import { Box, Image, HStack, Button, Text, Flex, Spinner, Tooltip } from '@chakra-ui/react'
import { useState, useEffect } from 'react'
import { isWebSocketConnected } from '../../services/api'

interface StartScreenProps {
  onNewGame: () => void;
  onLoadGame: () => void;
}

export const StartScreen = ({ onNewGame, onLoadGame }: StartScreenProps) => {
  const [imageLoaded, setImageLoaded] = useState(false);
  const [isConnected, setIsConnected] = useState(false);
  const [connectionChecked, setConnectionChecked] = useState(false);

  // Check WebSocket connection status
  useEffect(() => {
    const checkConnection = () => {
      const connected = isWebSocketConnected();
      setIsConnected(connected);
      setConnectionChecked(true);
    };

    // Initial check
    checkConnection();

    // Set up periodic connection check
    const intervalId = setInterval(checkConnection, 1000);

    return () => clearInterval(intervalId);
  }, []);

  return (
    <Box
      position="fixed"
      top={0}
      left={0}
      right={0}
      bottom={0}
      bg="#291326"
      display="flex"
      alignItems="center"
      justifyContent="center"
    >
      <Flex 
        direction="column" 
        align="center" 
        justify="center" 
        width="100%" 
        height="100%"
        px={4}
      >
        <Box 
          width="100%" 
          display="flex" 
          justifyContent="center" 
          alignItems="center"
          mb={12}
        >
          <Image
            src="/logo.png"
            alt="The Deeps Logo"
            maxW="600px"
            w="100%"
            objectFit="contain"
            onLoad={() => setImageLoaded(true)}
            onError={(e) => console.error('Logo failed to load:', e)}
          />
          {!imageLoaded && (
            <Text color="white" fontSize="2xl" position="absolute">Loading logo...</Text>
          )}
        </Box>
        
        <HStack spacing={6}>
          <Tooltip 
            label={!isConnected && connectionChecked ? "Waiting for server connection..." : ""}
            isDisabled={isConnected}
          >
            <Button
              size="md"
              bg="#6B46C1"
              color="white"
              width="200px"
              height="50px"
              fontSize="xl"
              onClick={onNewGame}
              _hover={{
                transform: isConnected ? 'scale(1.05)' : 'none',
                bg: isConnected ? '#805AD5' : '#6B46C1'
              }}
              _active={{
                bg: '#553C9A'
              }}
              border="2px solid"
              borderColor="purple.200"
              isDisabled={!isConnected}
              opacity={isConnected ? 1 : 0.7}
            >
              {!isConnected && connectionChecked ? (
                <Flex align="center">
                  <Spinner size="sm" mr={2} />
                  Connecting...
                </Flex>
              ) : "New Game"}
            </Button>
          </Tooltip>
          
          <Tooltip 
            label={!isConnected && connectionChecked ? "Waiting for server connection..." : ""}
            isDisabled={isConnected}
          >
            <Button
              size="md"
              bg="transparent"
              color="white"
              width="200px"
              height="50px"
              fontSize="xl"
              onClick={onLoadGame}
              border="2px solid"
              borderColor="purple.200"
              _hover={{
                transform: isConnected ? 'scale(1.05)' : 'none',
                bg: isConnected ? 'rgba(107, 70, 193, 0.2)' : 'transparent'
              }}
              _active={{
                bg: 'rgba(107, 70, 193, 0.4)'
              }}
              isDisabled={!isConnected}
              opacity={isConnected ? 1 : 0.7}
            >
              {!isConnected && connectionChecked ? (
                <Flex align="center">
                  <Spinner size="sm" mr={2} />
                  Connecting...
                </Flex>
              ) : "Load Game"}
            </Button>
          </Tooltip>
        </HStack>
        
        {!isConnected && connectionChecked && (
          <Text color="red.300" mt={4} fontSize="sm">
            Connecting to game server...
          </Text>
        )}
      </Flex>
    </Box>
  )
} 