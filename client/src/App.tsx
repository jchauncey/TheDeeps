import { ChakraProvider, useToast, Box, Flex, extendTheme } from '@chakra-ui/react'
import { useState, useEffect, useCallback } from 'react'
import { StartScreen } from './components/game/StartScreen'
import { CharacterCreation } from './components/game/CharacterCreation'
import { GameBoard } from './components/game/GameBoard'
import { GameControls } from './components/game/GameControls'
import { GameStatus } from './components/game/GameStatus'
import { connectWebSocket, sendWebSocketMessage } from './services/api'
import { CharacterData, DebugMessage } from './types/game'
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
  
  // Store debug messages
  const [debugMessages, setDebugMessages] = useState<DebugMessage[]>([])
  
  // WebSocket connection
  const [ws, setWs] = useState<WebSocket | null>(null)
  
  // Store floor data
  const [floorData, setFloorData] = useState<FloorData | null>(null)
  
  // Use our custom clickable toast
  const toast = useClickableToast();

  // Handle WebSocket messages
  const handleWebSocketMessage = useCallback((data: any) => {
    console.log('App received WebSocket message:', data);
    
    // Handle debug messages
    if (data.level) {
      setDebugMessages(prev => [...prev, data as DebugMessage]);
    }
    
    // Handle floor data
    if (data.type === 'floor_data') {
      console.log('App received floor data:', data);
      setFloorData(data);
    }
    
    // Dispatch a custom event for other components to listen to
    const event = new CustomEvent('websocket_message', { detail: data });
    window.dispatchEvent(event);
  }, []);

  // Connect to WebSocket when component mounts
  useEffect(() => {
    console.log('Connecting to WebSocket...');
    const socket = connectWebSocket(handleWebSocketMessage);
    setWs(socket);
    
    // Clean up WebSocket connection when component unmounts
    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, [handleWebSocketMessage]);

  // Request floor data when entering game screen
  useEffect(() => {
    if (currentScreen === 'game' && character) {
      console.log('Requesting floor data...');
      const success = sendWebSocketMessage({ 
        type: 'get_floor',
        characterId: character.name // Use character name as ID for now
      });
      console.log('Floor data request sent:', success);
    }
  }, [currentScreen, character]);

  // Handle keyboard controls at the App level
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (currentScreen !== 'game') return;
      
      let direction = '';
      
      // Handle both key and code
      if (e.key === 'ArrowUp' || e.key === 'w' || e.key === 'W') {
        direction = 'up';
      } else if (e.key === 'ArrowDown' || e.key === 's' || e.key === 'S') {
        direction = 'down';
      } else if (e.key === 'ArrowLeft' || e.key === 'a' || e.key === 'A') {
        direction = 'left';
      } else if (e.key === 'ArrowRight' || e.key === 'd' || e.key === 'D') {
        direction = 'right';
      } else {
        return; // Not a movement key
      }
      
      // Prevent default behavior (scrolling)
      e.preventDefault();
      
      // Send move command to server
      console.log(`Sending move command: ${direction}`);
      sendWebSocketMessage({
        type: 'move',
        direction
      });
    };
    
    // Add event listener
    window.addEventListener('keydown', handleKeyDown);
    
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
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
    
    // Send character data to server
    sendWebSocketMessage({
      type: 'create_character',
      character: characterData
    });
    
    toast({
      title: "Character Created",
      description: `${characterData.name} the ${characterData.characterClass} is ready for adventure!`,
      status: "success",
    })
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
              
              {/* Right side - Character Status */}
              <Box 
                w="280px" 
                h="100%" 
                position="relative"
                flexShrink={0} // Prevent shrinking
              >
                <GameStatus character={character} />
              </Box>
            </Flex>
            
            {/* Controls (overlay) */}
            <GameControls />
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
