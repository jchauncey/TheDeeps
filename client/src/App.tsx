import { ChakraProvider, Box, Flex, extendTheme, Text, Button } from '@chakra-ui/react'
import { useState, useEffect, useCallback, useRef } from 'react'
import { 
  StartScreen, 
  CharacterCreation, 
  CharacterSelection,
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
  closeWebSocketConnection,
  getSavedCharacters,
  createCharacter,
  createDungeon,
  joinDungeon,
  identifyCharacter
} from './services/api'
import { CharacterData, FloorData } from './types/game'
import { useClickableToast } from './components/ui/ClickableToast'

// Define the screens we can navigate to
type Screen = 'loading' | 'characterSelection' | 'characterCreation' | 'dungeonSelection' | 'game'

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
  const [currentScreen, setCurrentScreen] = useState<Screen>('loading')
  const [character, setCharacter] = useState<CharacterData | null>(null)
  const [dungeonId, setDungeonId] = useState<string | null>(null)
  const [floorData, setFloorData] = useState<FloorData | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [characterRefreshTrigger, setCharacterRefreshTrigger] = useState(0)
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
  }, [character, floorData, dungeonId])
  
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
      console.log('App: WebSocket message received:', data)
      
      if (data.type === 'floor_data') {
        console.log('Received floor_data message:', data)
        
        // Update floor data
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
            'dungeonId =', data.dungeonId || 'null')
        }
      } else if (data.type === 'welcome') {
        console.log('Received welcome message:', data.message)
      } else if (data.type === 'error') {
        console.error('Received error message:', data.message)
        toast({
          title: 'Error',
          description: data.message,
          status: 'error',
          duration: 5000,
          isClosable: true,
        })
      } else if (data.type === 'character_created') {
        console.log('Received character_created message:', data);
        
        // Update character with ID from server
        if (data.character && data.character.id) {
          console.log('Setting character ID:', data.character.id);
          
          // Create a complete character object from the server response
          const updatedCharacter = {
            ...data.character,
            // Ensure all required fields are present
            name: data.character.name,
            characterClass: data.character.characterClass,
            stats: data.character.stats,
            id: data.character.id
          };
          
          console.log('Updated character object:', updatedCharacter);
          
          // Update the character state
          setCharacter(updatedCharacter);
          
          toast({
            title: 'Character Created',
            description: data.message || 'Character created successfully',
            status: 'success',
          });
          
          // Move to dungeon selection screen
          console.log('Moving to dungeon selection screen');
          setCurrentScreen('dungeonSelection');
          
          // If we already have floor data, transition to the game screen
          if (floorData && floorData.dungeonId) {
            console.log('We have floor data, transitioning to game screen with dungeonId:', floorData.dungeonId);
            setDungeonId(floorData.dungeonId);
            setCurrentScreen('game');
          }
        } else {
          console.error('Character created but no ID received:', data);
          toast({
            title: 'Error',
            description: 'Character created but no ID received',
            status: 'error',
          });
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
      } else if (data.type === 'debug') {
        // Handle debug messages
        console.log(`Debug [${data.level}]: ${data.message}`)
      }
    } catch (error) {
      console.error('Error handling WebSocket message:', error)
    }
  }, [currentScreen, character, toast])
  
  // Initialize WebSocket connection
  const initializeWebSocket = useCallback(() => {
    console.log('Initializing WebSocket connection...')
    
    // Only connect if not already connected
    if (!isWebSocketConnected()) {
      const ws = connectWebSocket(handleWebSocketMessage)
      
      if (ws) {
        console.log('WebSocket connection initialized')
        setIsConnected(true)
        
        // Set up callbacks for connection events
        setWebSocketCallbacks(
          // onDisconnect
          () => {
            console.log('WebSocket disconnected')
            setIsConnected(false)
          },
          // onReconnectFailed
          () => {
            console.log('WebSocket reconnection failed')
            setIsConnected(false)
            toast({
              title: 'Connection Lost',
              description: 'Failed to reconnect to the game server. Please refresh the page.',
              status: 'error',
              duration: null,
              isClosable: true,
            })
          },
          // onConnected
          () => {
            console.log('WebSocket connected')
            setIsConnected(true)
            
            // If we have a character, send an identify_character message
            if (character && character.id) {
              console.log('Sending identify_character message for character:', character.id)
              identifyCharacter(character.id)
            }
          }
        )
      } else {
        console.log('Failed to initialize WebSocket connection')
        setIsConnected(false)
      }
    } else {
      console.log('WebSocket already connected')
      setIsConnected(true)
      
      // Even if already connected, make sure the character is identified
      if (character && character.id) {
        console.log('WebSocket already connected, sending identify_character message for character:', character.id)
        identifyCharacter(character.id)
      }
    }
  }, [handleWebSocketMessage, toast, character])
  
  // Effect to initialize WebSocket when needed
  useEffect(() => {
    if (currentScreen === 'game' && !isConnected) {
      initializeWebSocket()
    }
  }, [currentScreen, isConnected, initializeWebSocket])
  
  // Effect to load saved characters when the app starts
  useEffect(() => {
    const loadSavedCharacters = async () => {
      try {
        const result = await getSavedCharacters();
        if (result.success && result.characters) {
          console.log('Loaded saved characters:', result.characters);
        }
      } catch (error) {
        console.error('Error loading saved characters:', error);
      }
    };

    // We'll load characters in the background, but won't wait for them
    // The CharacterSelection screen will handle loading and displaying them
    loadSavedCharacters();
  }, []);

  // Handle character creation
  const handleCreateCharacter = async (characterData: CharacterData) => {
    try {
      const result = await createCharacter(characterData);
      if (result.success && result.characterId) {
        console.log('Character created successfully:', result.characterId);
        
        // Load the created character
        const loadResult = await loadCharacter(result.characterId);
        if (loadResult.success && loadResult.character) {
          setCharacter(loadResult.character);
          // Navigate to character selection screen
          setCurrentScreen('characterSelection');
          // Trigger a refresh of the character list
          setCharacterRefreshTrigger(prev => prev + 1);
        } else {
          console.error('Error loading created character:', loadResult.message);
          toast({
            title: 'Error',
            description: loadResult.message || 'Failed to load character',
            status: 'error',
          });
        }
      } else {
        console.error('Error creating character:', result.message);
        toast({
          title: 'Error',
          description: result.message || 'Failed to create character',
          status: 'error',
        });
      }
    } catch (error) {
      console.error('Error in character creation:', error);
      toast({
        title: 'Error',
        description: 'An unexpected error occurred',
        status: 'error',
      });
    }
  };

  // Handle character selection
  const handleSelectCharacter = async (characterId: string) => {
    setIsLoading(true);
    try {
      const result = await loadCharacter(characterId);
      if (result.success && result.character) {
        setCharacter(result.character);
        // Navigate to dungeon selection screen
        setCurrentScreen('dungeonSelection');
      } else {
        console.error('Error loading character:', result.message);
        toast({
          title: 'Error',
          description: result.message || 'Failed to load character',
          status: 'error',
        });
      }
    } catch (error) {
      console.error('Error selecting character:', error);
      toast({
        title: 'Error',
        description: 'An unexpected error occurred',
        status: 'error',
      });
    } finally {
      setIsLoading(false);
    }
  };

  // Handle dungeon creation
  const handleCreateDungeon = async (name: string, numFloors: number) => {
    try {
      const result = await createDungeon(name, numFloors);
      if (result.success && result.dungeonId) {
        console.log('Dungeon created successfully:', result.dungeonId);
        
        // Join the created dungeon
        await handleJoinDungeon(result.dungeonId);
      } else {
        console.error('Error creating dungeon:', result.message);
        toast({
          title: 'Error',
          description: result.message || 'Failed to create dungeon',
          status: 'error',
        });
      }
    } catch (error) {
      console.error('Error in dungeon creation:', error);
      toast({
        title: 'Error',
        description: 'An unexpected error occurred',
        status: 'error',
      });
    }
  };

  // Handle joining a dungeon
  const handleJoinDungeon = async (dungeonId: string) => {
    if (!character || !character.id) {
      console.error('No character selected or character has no ID');
      toast({
        title: 'Error',
        description: 'No valid character selected',
        status: 'error',
      });
      return;
    }

    try {
      const result = await joinDungeon(dungeonId, character.id);
      if (result.success && result.floorData) {
        console.log('Joined dungeon successfully:', result.floorData);
        setFloorData(result.floorData);
        setDungeonId(dungeonId);
        
        // Initialize WebSocket for real-time updates
        initializeWebSocket();
        
        // Explicitly identify the character with the WebSocket connection
        // This is crucial for the server to know which character is controlled by this connection
        identifyCharacter(character.id);
        
        // The game screen transition will happen automatically when floorData and dungeonId are set
      } else {
        console.error('Error joining dungeon:', result.message);
        toast({
          title: 'Error',
          description: result.message || 'Failed to join dungeon',
          status: 'error',
        });
      }
    } catch (error) {
      console.error('Error joining dungeon:', error);
      toast({
        title: 'Error',
        description: 'An unexpected error occurred',
        status: 'error',
      });
    }
  };

  // Handle dungeon selection
  const handleDungeonSelected = (dungeonId: string, dungeonFloorData: FloorData) => {
    setDungeonId(dungeonId)
    setFloorData(dungeonFloorData)
    
    // Initialize WebSocket for real-time updates
    initializeWebSocket();
    
    // If we have a character, identify it with the WebSocket connection
    if (character && character.id) {
      identifyCharacter(character.id);
    }
    
    // The game screen transition will happen automatically
  }

  // Handle back to character creation
  const handleBackToCharacterSelection = () => {
    setCurrentScreen('characterSelection')
  }

  // Handle back to start screen
  const handleBackToStartScreen = () => {
    setCurrentScreen('loading')
  }

  // Handle transition from loading to character selection
  const handleCharactersLoaded = () => {
    setCurrentScreen('characterSelection')
  }

  // Clean up WebSocket connection when component unmounts
  useEffect(() => {
    return () => {
      console.log('App unmounting, closing WebSocket connection')
      closeWebSocketConnection()
    }
  }, [])

  return (
    <ChakraProvider theme={theme}>
      {currentScreen === 'loading' && (
        <StartScreen 
          onCharactersLoaded={handleCharactersLoaded}
          isLoading={isLoading}
        />
      )}
      
      {currentScreen === 'characterSelection' && (
        <CharacterSelection
          onSelectCharacter={handleSelectCharacter}
          onCreateNewCharacter={() => setCurrentScreen('characterCreation')}
          onBack={handleBackToStartScreen}
          refreshTrigger={characterRefreshTrigger}
        />
      )}
      
      {currentScreen === 'characterCreation' && (
        <CharacterCreation 
          onCreateCharacter={handleCreateCharacter}
          onBack={handleBackToCharacterSelection}
        />
      )}
      
      {currentScreen === 'dungeonSelection' && character && (
        <DungeonSelection
          characterId={character.id || ''}
          onDungeonSelected={handleDungeonSelected}
          onBack={handleBackToCharacterSelection}
        />
      )}
      
      {currentScreen === 'game' && floorData && (
        <Box position="relative" width="100vw" height="100vh" overflow="hidden">
          <GameBoard floorData={floorData} />
          <GameControls 
            character={character}
            onNewGame={() => {
              // Reset game state and go back to character selection
              setCharacter(null);
              setFloorData(null);
              setDungeonId(null);
              closeWebSocketConnection();
              setIsConnected(false);
              setCurrentScreen('characterSelection');
            }}
            onLoadGame={() => {
              // For now, just show a message that this isn't implemented
              toast({
                title: "Load Game",
                description: "This feature is not yet implemented.",
                status: "info",
              });
            }}
          />
          <GameStatusSimple 
            character={character}
            currentFloor={floorData.currentFloor}
            dungeonName={floorData.dungeonName || 'The Deeps'}
          />
        </Box>
      )}
    </ChakraProvider>
  )
}

export default App
