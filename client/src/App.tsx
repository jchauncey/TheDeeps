import { ChakraProvider, useToast, Box, Flex } from '@chakra-ui/react'
import { useState, useEffect } from 'react'
import { StartScreen } from './components/game/StartScreen'
import { CharacterCreation } from './components/game/CharacterCreation'
import { GameBoard } from './components/game/GameBoard'
import { GameControls } from './components/game/GameControls'
import { GameStatus } from './components/game/GameStatus'
import { connectWebSocket } from './services/api'
import { CharacterData, DebugMessage } from './types/game'

// Define the screens we can navigate to
type Screen = 'start' | 'characterCreation' | 'game'

function App() {
  // Track which screen we're on
  const [currentScreen, setCurrentScreen] = useState<Screen>('start')
  
  // Store character data
  const [character, setCharacter] = useState<CharacterData | null>(null)
  
  // Store debug messages
  const [debugMessages, setDebugMessages] = useState<DebugMessage[]>([])
  
  // WebSocket connection
  const [ws, setWs] = useState<WebSocket | null>(null)
  
  const toast = useToast()

  // Connect to WebSocket when component mounts
  useEffect(() => {
    const socket = connectWebSocket((data) => {
      // Handle incoming WebSocket messages
      if (data.level) {
        // It's a debug message
        setDebugMessages(prev => [...prev, data as DebugMessage])
      }
      
      // Dispatch a custom event for other components to listen to
      const event = new CustomEvent('websocket_message', { detail: data });
      window.dispatchEvent(event);
    })
    
    setWs(socket)
    
    // Clean up WebSocket connection when component unmounts
    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close()
      }
    }
  }, [])

  const handleNewGame = () => {
    setCurrentScreen('characterCreation')
  }

  const handleLoadGame = () => {
    // TODO: Implement load game functionality
    toast({
      title: "Load Game",
      description: "This feature is not yet implemented.",
      status: "info",
      duration: 3000,
      isClosable: true,
    })
  }

  const handleBackToStart = () => {
    setCurrentScreen('start')
  }

  const handleCreateCharacter = (characterData: CharacterData) => {
    setCharacter(characterData)
    setCurrentScreen('game')
    
    toast({
      title: "Character Created",
      description: `${characterData.name} the ${characterData.characterClass} is ready for adventure!`,
      status: "success",
      duration: 5000,
      isClosable: true,
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
                <GameBoard />
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
    <ChakraProvider>
      {renderScreen()}
    </ChakraProvider>
  )
}

export default App
