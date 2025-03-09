import { ChakraProvider, Box, Flex, extendTheme } from '@chakra-ui/react'
import { useState, useEffect, useCallback, useRef } from 'react'
import { 
  StartScreen, 
  CharacterCreation, 
  GameBoard, 
  GameControls, 
  GameStatusSimple,
  DungeonSelection
} from './components/game'
import { 
  connectWebSocket, 
  sendWebSocketMessage, 
  isWebSocketConnected, 
  setWebSocketCallbacks 
} from './services/api'
import { CharacterData, FloorData } from './types/game'
import { useClickableToast } from './components/ui/ClickableToast'

// Define the screens we can navigate to
type Screen = 'start' | 'characterCreation' | 'dungeonSelection' | 'game'

// Create a custom theme with toast configuration
const theme = extendTheme({
  // Configure toast defaults
  toast: {
    defaultOptions: {
      position: 'top',
      duration: 5000,
      isClosable: true,
    },
  },
})

function App() {
  const [currentScreen, setCurrentScreen] = useState<Screen>('start')
  const [character, setCharacter] = useState<CharacterData | null>(null)
  const [dungeonId, setDungeonId] = useState<string | null>(null)
  const [floorData, setFloorData] = useState<FloorData | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [connectionAttempted, setConnectionAttempted] = useState(false)
  const toast = useClickableToast()
  
  // Handle WebSocket messages
  const handleWebSocketMessage = useCallback((event: Event) => {
    try {
      const messageEvent = event as MessageEvent
      const data = JSON.parse(messageEvent.data)
      console.log('WebSocket message received:', data)
      
      if (data.type === 'floor_data') {
        setFloorData(data)
      } else if (data.type === 'welcome') {
        toast({
          title: 'Connected',
          description: data.message,
          status: 'success',
        })
      } else if (data.type === 'error') {
        toast({
          title: 'Error',
          description: data.message,
          status: 'error',
        })
      } else if (data.type === 'character_created') {
        // Update character with ID from server
        if (data.character && data.character.id) {
          console.log('Setting character ID:', data.character.id)
          
          // Create a complete character object from the server response
          const updatedCharacter = {
            ...data.character,
            // Ensure all required fields are present
            name: data.character.name,
            characterClass: data.character.characterClass,
            stats: data.character.stats,
            id: data.character.id
          };
          
          // Update the character state
          setCharacter(updatedCharacter);
          
          toast({
            title: 'Character Created',
            description: data.message || 'Character created successfully',
            status: 'success',
          })
          
          // Move to dungeon selection screen
          setCurrentScreen('dungeonSelection')
        } else {
          console.error('Character created but no ID received');
          toast({
            title: 'Error',
            description: 'Character created but no ID received',
            status: 'error',
          })
        }
      } else if (data.type === 'dungeon_created') {
        toast({
          title: 'Success',
          description: data.message || 'Dungeon created successfully',
          status: 'success',
        })
      } else if (data.type === 'dungeon_joined') {
        toast({
          title: 'Success',
          description: data.message || 'Dungeon joined successfully',
          status: 'success',
        })
        
        // Request floor data
        if (character && character.id) {
          sendWebSocketMessage({ 
            type: 'get_floor'
          });
        }
      }
    } catch (error) {
      console.error('Error handling WebSocket message:', error)
    }
  }, [toast, character, setCharacter])
  
  // Handle WebSocket connection
  const initializeWebSocket = useCallback(() => {
    if (!connectionAttempted) {
      setConnectionAttempted(true)
      
      const ws = connectWebSocket(handleWebSocketMessage)
      
      // Set up WebSocket event callbacks
      setWebSocketCallbacks(
        // onDisconnect
        () => {
          setIsConnected(false)
          toast({
            title: 'Disconnected',
            description: 'Connection to server closed',
            status: 'warning',
          })
        },
        // onReconnectFailed
        () => {
          setIsConnected(false)
          toast({
            title: 'Connection Failed',
            description: 'Unable to connect to game server. Please refresh the page.',
            status: 'error',
            duration: null,
            isClosable: true,
          })
        },
        // onConnected
        () => {
          setIsConnected(true)
          toast({
            title: 'Connected',
            description: 'Connected to game server',
            status: 'success',
          })
        }
      )
      
      if (ws) {
        setIsConnected(true)
      } else {
        setIsConnected(false)
        toast({
          title: 'Connection Failed',
          description: 'Unable to connect to game server',
          status: 'error',
        })
      }
    }
  }, [connectionAttempted, handleWebSocketMessage, toast])
  
  // Initialize WebSocket on component mount
  useEffect(() => {
    initializeWebSocket()
    
    // Check connection status periodically
    const connectionCheck = setInterval(() => {
      const connected = isWebSocketConnected()
      setIsConnected(connected)
    }, 5000)
    
    return () => {
      clearInterval(connectionCheck)
    }
  }, [initializeWebSocket])
  
  // Handle character creation
  const handleCreateCharacter = (characterData: CharacterData) => {
    if (!isConnected) {
      toast({
        title: 'Not Connected',
        description: 'Cannot create character: not connected to server',
        status: 'error',
      })
      return
    }
    
    setCharacter(characterData)
    
    // Send character creation message
    const success = sendWebSocketMessage({
      type: 'create_character',
      character: {
        name: characterData.name,
        characterClass: characterData.characterClass,
        stats: characterData.stats
      }
    })
    
    if (!success) {
      toast({
        title: 'Error',
        description: 'Failed to send character data to server',
        status: 'error',
      })
    } else {
      toast({
        title: 'Creating Character',
        description: 'Character creation request sent',
        status: 'info',
        duration: 3000,
      })
      // We'll wait for the character_created message before navigating
    }
  }
  
  // Handle dungeon selection
  const handleDungeonSelected = (dungeonId: string, dungeonFloorData: FloorData) => {
    setDungeonId(dungeonId)
    setFloorData(dungeonFloorData)
    setCurrentScreen('game')
  }
  
  // Handle back to character creation
  const handleBackToCharacterCreation = () => {
    setCurrentScreen('characterCreation')
  }
  
  // Handle new game
  const handleNewGame = () => {
    setCurrentScreen('characterCreation')
  }
  
  // Handle load game
  const handleLoadGame = () => {
    toast({
      title: "Load Game",
      description: "This feature is not yet implemented.",
      status: "info",
    })
  }
  
  return (
    <ChakraProvider theme={theme}>
      <Flex 
        direction="column" 
        h="100vh" 
        w="100vw"
        bg="gray.900" 
        color="white"
        overflow="hidden"
      >
        {currentScreen === 'start' && (
          <StartScreen 
            onNewGame={handleNewGame} 
            onLoadGame={handleLoadGame}
          />
        )}
        
        {currentScreen === 'characterCreation' && (
          <CharacterCreation 
            onCreateCharacter={handleCreateCharacter} 
            onBack={() => setCurrentScreen('start')}
          />
        )}
        
        {currentScreen === 'dungeonSelection' && character && (
          <Box flex="1" overflow="hidden" position="relative" w="100%" h="100%">
            <DungeonSelection 
              onDungeonSelected={handleDungeonSelected}
              onBack={handleBackToCharacterCreation}
              characterId={character.id || ''}
            />
          </Box>
        )}
        
        {currentScreen === 'game' && floorData && (
          <Flex 
            flex="1" 
            overflow="hidden" 
            position="relative" 
            width="100%" 
            height="100%"
          >
            {/* Map window anchored to the left */}
            <Box 
              flex="1" 
              height="100%" 
              overflow="hidden"
            >
              <GameBoard floorData={floorData} />
            </Box>
            
            {/* Character status panel anchored to the right */}
            <Box 
              width="300px" 
              height="100%" 
              position="relative"
            >
              <GameStatusSimple 
                character={character}
                currentFloor={floorData.currentFloor}
                dungeonName={floorData.dungeonName || 'The Deeps'}
              />
            </Box>
            
            {/* Game controls positioned at the bottom */}
            <Box 
              position="absolute" 
              bottom="0" 
              left="0" 
              width="100%" 
              zIndex="10"
            >
              <GameControls 
                character={character}
                onNewGame={handleNewGame}
                onLoadGame={handleLoadGame}
              />
            </Box>
          </Flex>
        )}
      </Flex>
    </ChakraProvider>
  )
}

export default App
