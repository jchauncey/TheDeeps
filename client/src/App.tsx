import { ChakraProvider, Box, Flex, extendTheme } from '@chakra-ui/react'
import { useState, useEffect, useCallback } from 'react'
import { 
  StartScreen, 
  CharacterCreation, 
  GameBoard, 
  GameControls, 
  GameStatusSimple 
} from './components/game'
import { connectWebSocket, sendWebSocketMessage } from './services/api'
import { CharacterData } from './types/game'
import { useClickableToast } from './components/ui/ClickableToast'

// Define the screens we can navigate to
type Screen = 'start' | 'characterCreation' | 'game'

// Define the floor data interface
interface Position {
  x: number;
  y: number;
}

interface FloorData {
  type: string;
  floor: any; // Using any for simplicity, but should match the Floor interface in GameBoard
  playerPosition: Position;
  currentFloor: number;
}

// Create a custom theme with toast configuration
const theme = extendTheme({
  // Configure toast defaults
  toast: {
    defaultOptions: {
      position: 'top',
      duration: 5000,
      isClosable: true,
      variant: 'solid',
    },
  },
})

function App() {
  // Track which screen we're on
  const [currentScreen, setCurrentScreen] = useState<Screen>('start')
  
  // Store character data
  const [character, setCharacter] = useState<CharacterData | null>(null)
  
  // Store floor data
  const [floorData, setFloorData] = useState<FloorData | null>(null)
  
  // Use our custom clickable toast
  const toast = useClickableToast();

  // Handle WebSocket messages
  const handleWebSocketMessage = useCallback((event: Event) => {
    const customEvent = event as CustomEvent;
    const data = customEvent.detail;
    
    console.log('WebSocket message received:', data);
    
    if (data.type === 'floor_data') {
      setFloorData(data);
      
      // Update character data if playerData is provided
      if (data.playerData && character) {
        // Convert server character data to client format
        const updatedCharacter = {
          ...character,
          health: data.playerData.health,
          maxHealth: data.playerData.maxHealth,
          mana: data.playerData.mana,
          maxMana: data.playerData.maxMana,
          experience: data.playerData.experience,
          level: data.playerData.level,
          gold: data.playerData.gold,
          status: data.playerData.status || [],
          // Add other properties as needed
        };
        
        setCharacter(updatedCharacter);
      }
    } else if (data.type === 'character_created') {
      // Handle character creation success
      console.log('Character created successfully:', data.character);
      
      // Request floor data once after character creation
      sendWebSocketMessage({ 
        type: 'get_floor',
        characterId: data.character.name
      });
      
      toast({
        title: "Character Created",
        description: data.message,
        status: "success",
      });
    } else if (data.type === 'debug') {
      // Handle debug messages
      toast({
        title: "Debug",
        description: data.message,
        status: data.level === "error" ? "error" : "info",
        duration: 3000,
        isClosable: true,
      });
    }
  }, [character, toast]);

  // Connect to WebSocket when component mounts
  useEffect(() => {
    console.log('Connecting to WebSocket...');
    const socket = connectWebSocket(handleWebSocketMessage);
    
    // Add event listener for reconnection failures
    const handleReconnectFailed = () => {
      toast({
        title: "Connection Error",
        description: "Failed to connect to the game server. Please refresh the page.",
        status: "error",
        duration: null, // Don't auto-dismiss
        isClosable: true,
      });
    };
    
    // Add event listener for successful connections
    const handleConnected = () => {
      // If we're in the game screen, request floor data if we don't have it yet
      if (currentScreen === 'game' && character && !floorData) {
        console.log('Connected, no floor data yet, requesting floor data...');
        sendWebSocketMessage({ 
          type: 'get_floor',
          characterId: character.name // Use character name as ID for now
        });
      }
    };
    
    // Add event listeners
    window.addEventListener('websocket_reconnect_failed', handleReconnectFailed);
    window.addEventListener('websocket_connected', handleConnected);
    
    // Clean up WebSocket connection and event listeners when component unmounts
    return () => {
      window.removeEventListener('websocket_reconnect_failed', handleReconnectFailed);
      window.removeEventListener('websocket_connected', handleConnected);
      
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, [handleWebSocketMessage, currentScreen, character, toast]);

  // Request floor data when entering game screen
  useEffect(() => {
    if (currentScreen === 'game' && character && !floorData) {
      console.log('Game screen active, no floor data yet, requesting floor data...');
      // Add a small delay to ensure WebSocket is connected
      const timeoutId = setTimeout(() => {
        const success = sendWebSocketMessage({ 
          type: 'get_floor',
          characterId: character.name // Use character name as ID for now
        });
        console.log('Floor data request sent:', success);
        
        if (!success) {
          toast({
            title: "Connection Issue",
            description: "Having trouble connecting to the game server. Retrying...",
            status: "warning",
          });
        }
      }, 500);
      
      return () => clearTimeout(timeoutId);
    }
  }, [currentScreen, character, floorData, toast]);

  // Handle keyboard controls at the App level
  useEffect(() => {
    // We're removing the keyboard handling from App.tsx since it's already in GameControls.tsx
    // This avoids conflicts between the two handlers
  }, [currentScreen]);

  const handleNewGame = () => {
    setCurrentScreen('characterCreation')
  }

  const handleLoadGame = () => {
    // TODO: Implement load game functionality
    toast({
      title: "Load Game",
      description: "This feature is not yet implemented.",
      status: "info",
    })
  }

  const handleBackToStart = () => {
    setCurrentScreen('start')
  }

  const handleCreateCharacter = (characterData: CharacterData) => {
    setCharacter(characterData)
    setCurrentScreen('game')
    
    // Ensure WebSocket connection is established before sending character data
    console.log('Creating character, ensuring WebSocket connection...');
    
    // Add a small delay to ensure WebSocket is connected
    setTimeout(() => {
      const success = sendWebSocketMessage({
        type: 'create_character',
        character: characterData
      });
      
      if (!success) {
        toast({
          title: "Connection Issue",
          description: "Having trouble connecting to the game server. Your character will be created when the connection is established.",
          status: "warning",
          duration: 5000,
        });
      }
    }, 500);
    
    // Don't show toast here, we'll show it when we receive the character_created message
  }

  // Render the appropriate screen
  const renderScreen = () => {
    switch (currentScreen) {
      case 'start':
        return <StartScreen onNewGame={handleNewGame} onLoadGame={handleLoadGame} />
      case 'characterCreation':
        return <CharacterCreation onCreateCharacter={handleCreateCharacter} onBack={handleBackToStart} />
      case 'game':
        return (
          <Box 
            w="100vw" 
            h="100vh" 
            bg="#291326" 
            overflow="hidden"
            position="relative"
          >
            <Flex 
              w="100%" 
              h="100%" 
              p={4}
              gap={4}
            >
              {/* Left side - Game Board */}
              <Box 
                flex="1" 
                h="100%" 
                position="relative"
                borderRadius="md"
                overflow="hidden"
                minW="0" // Important for flex child to shrink properly
              >
                <GameBoard floorData={floorData} />
              </Box>
              
              {/* Right side - Character Status (simplified) */}
              <Box 
                w="280px" 
                h="100%" 
                position="relative"
                flexShrink={0} // Prevent shrinking
              >
                <GameStatusSimple character={character} />
              </Box>
            </Flex>
            
            {/* Controls (overlay) */}
            <GameControls 
              character={character}
              onNewGame={handleNewGame}
              onLoadGame={handleLoadGame}
            />
          </Box>
        )
      default:
        return <StartScreen onNewGame={handleNewGame} onLoadGame={handleLoadGame} />
    }
  }

  return (
    <ChakraProvider theme={theme}>
      {renderScreen()}
    </ChakraProvider>
  )
}

export default App
