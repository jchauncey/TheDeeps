import { ChakraProvider, Box, Flex, extendTheme, Text, Button } from '@chakra-ui/react'
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
  setWebSocketCallbacks,
  loadCharacter,
  closeWebSocketConnection
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
  
  // Function to check if we can transition to the game screen
  const checkAndTransitionToGame = useCallback(() => {
    console.log('Checking if we can transition to game screen:',
      'character =', character ? `${character.name} (${character.id})` : 'null',
      'floorData =', floorData ? 'exists' : 'null',
      'dungeonId =', dungeonId || 'null')
    
    if (character && floorData && dungeonId) {
      console.log('Transitioning to game screen with dungeonId:', dungeonId)
      setCurrentScreen('game')
      return true
    }
    return false
  }, [character, floorData, dungeonId, setCurrentScreen])
  
  // Effect to check if we can transition to the game screen when any of the dependencies change
  useEffect(() => {
    if (currentScreen !== 'game') {
      checkAndTransitionToGame()
    }
  }, [character, floorData, dungeonId, currentScreen, checkAndTransitionToGame])
  
  // Handle WebSocket messages
  const handleWebSocketMessage = useCallback((event: Event) => {
    try {
      const messageEvent = event as MessageEvent
      const data = JSON.parse(messageEvent.data)
      console.log('WebSocket message received:', data)
      
      if (data.type === 'floor_data') {
        console.log('Received floor_data message:', data)
        setFloorData(data)
        
        // Dispatch a custom event to notify other components
        const customEvent = new CustomEvent('websocket_message', {
          detail: data
        });
        window.dispatchEvent(customEvent);
        
        // If we're not already on the game screen and we have a character and dungeon ID,
        // transition to the game screen
        if (currentScreen !== 'game' && character && data.dungeonId) {
          console.log('Transitioning to game screen with dungeonId:', data.dungeonId)
          setDungeonId(data.dungeonId)
          setCurrentScreen('game')
        } else {
          console.log('Not transitioning to game screen:',
            'currentScreen =', currentScreen,
            'character =', character ? 'exists' : 'null',
            'data.dungeonId =', data.dungeonId || 'missing')
          
          // Set the dungeonId anyway, so we can transition later when character is available
          if (data.dungeonId) {
            setDungeonId(data.dungeonId)
          }
        }
      } else if (data.type === 'welcome') {
        toast({
          title: 'Connected',
          description: data.message,
          status: 'success',
        })
      } else if (data.type === 'error') {
        console.error('Received error message:', data.message)
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
          
          // If we already have floor data, transition to the game screen
          if (floorData && floorData.dungeonId) {
            console.log('We have floor data, transitioning to game screen with dungeonId:', floorData.dungeonId)
            setDungeonId(floorData.dungeonId)
            setCurrentScreen('game')
          }
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
        console.log('Received dungeon_joined message:', data)
        toast({
          title: 'Success',
          description: data.message || 'Dungeon joined successfully',
          status: 'success',
        })
        
        // The server should automatically send floor_data after joining a dungeon
        // We'll set the dungeonId here in case we need it
        if (data.dungeonId) {
          console.log('Setting dungeonId from dungeon_joined message:', data.dungeonId)
          setDungeonId(data.dungeonId)
          
          // If we have both character and floor data, transition to the game screen
          if (character && floorData) {
            console.log('We have character and floor data, transitioning to game screen')
            setCurrentScreen('game')
          } else {
            console.log('Not transitioning to game screen yet:',
              'character =', character ? 'exists' : 'null',
              'floorData =', floorData ? 'exists' : 'null')
            
            // Try to transition after a short delay to allow state updates to propagate
            setTimeout(() => {
              checkAndTransitionToGame()
            }, 500)
          }
        } else {
          console.warn('dungeon_joined message missing dungeonId')
        }
      }
    } catch (error) {
      console.error('Error handling WebSocket message:', error)
    }
  }, [toast, character, setCharacter, currentScreen, setDungeonId, setCurrentScreen, floorData, checkAndTransitionToGame])
  
  // Effect to load character from localStorage when component mounts
  useEffect(() => {
    const loadCharacterFromStorage = async () => {
      try {
        const storedCharacterId = localStorage.getItem('currentCharacterId');
        if (storedCharacterId && !character) {
          console.log('Found stored characterId in localStorage:', storedCharacterId);
          
          // Load character data from server
          const result = await loadCharacter(storedCharacterId);
          if (result.success && result.character) {
            console.log('Loaded character from server:', result.character);
            setCharacter(result.character);
          } else {
            console.error('Failed to load character from server:', result.message);
          }
        }
      } catch (error) {
        console.error('Error loading character from localStorage:', error);
      }
    };
    
    loadCharacterFromStorage();
  }, [character, setCharacter]);
  
  // Effect to log state changes for debugging
  useEffect(() => {
    console.log('App state changed:',
      'currentScreen =', currentScreen,
      'character =', character ? `${character.name} (${character.id})` : 'null',
      'dungeonId =', dungeonId || 'null',
      'floorData =', floorData ? 'exists' : 'null')
  }, [currentScreen, character, dungeonId, floorData])
  
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
            description: 'Could not connect to server after multiple attempts',
            status: 'error',
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
    // Reset game state
    setCharacter(null)
    setFloorData(null)
    setDungeonId(null)
    
    // Close WebSocket connection if it exists
    closeWebSocketConnection()
    setIsConnected(false)
    setConnectionAttempted(false)
    
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
        
        {currentScreen === 'dungeonSelection' && character && character.id ? (
          <Box flex="1" overflow="hidden" position="relative" w="100%" h="100%">
            <DungeonSelection 
              onDungeonSelected={handleDungeonSelected}
              onBack={handleBackToCharacterCreation}
              characterId={character.id}
            />
          </Box>
        ) : currentScreen === 'dungeonSelection' && (!character || !character.id) ? (
          // If we're on the dungeon selection screen but don't have a character,
          // show an error and redirect to character creation
          <Box flex="1" display="flex" justifyContent="center" alignItems="center" flexDirection="column">
            <Text fontSize="xl" mb={4}>No character selected. Please create a character first.</Text>
            <Button colorScheme="blue" onClick={handleBackToCharacterCreation}>Create Character</Button>
          </Box>
        ) : null}
        
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
