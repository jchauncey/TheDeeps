import { CharacterData, FloorData } from '../types/game';
import { getFloorDataWS } from './api';
import { Dispatch, SetStateAction } from 'react';

// Define the handler types
type MessageHandlerContext<ScreenType = string> = {
  setCharacter: Dispatch<SetStateAction<CharacterData | null>>;
  setFloorData: Dispatch<SetStateAction<FloorData | null>>;
  setDungeonId: Dispatch<SetStateAction<string | null>>;
  setCurrentScreen: Dispatch<SetStateAction<ScreenType>>;
  character: CharacterData | null;
  floorData: FloorData | null;
  currentScreen: ScreenType;
  showToast: (options: any) => void;
};

// Handler for floor_data messages
export const handleFloorDataMessage = <ScreenType>(data: any, context: MessageHandlerContext<ScreenType>) => {
  const { setFloorData, setDungeonId, setCurrentScreen, character, currentScreen } = context;
  
  console.log('=== FLOOR DATA DEBUG ===');
  console.log('Handling floor_data message:', data);
  console.log('Current character:', character);
  console.log('Current screen:', currentScreen);
  console.log('Dungeon ID in floor data:', data.dungeonId);
  
  setFloorData(data);
  
  // If we're not already on the game screen and we have a character and dungeon ID,
  // transition to the game screen
  if (currentScreen !== 'game' && character && data.dungeonId) {
    console.log('FLOOR DATA: Transitioning to game screen with dungeonId:', data.dungeonId);
    setDungeonId(data.dungeonId);
    setCurrentScreen('game' as unknown as ScreenType);
  } else {
    console.log('FLOOR DATA: Not transitioning to game screen:',
      'currentScreen =', currentScreen,
      'character =', character ? `${character.name} (${character.id})` : 'null',
      'data.dungeonId =', data.dungeonId || 'missing');
    
    // Set the dungeonId anyway, so we can transition later when character is available
    if (data.dungeonId) {
      console.log('FLOOR DATA: Setting dungeonId for later transition:', data.dungeonId);
      setDungeonId(data.dungeonId);
    }
  }
};

// Handler for character_created messages
export const handleCharacterCreatedMessage = <ScreenType>(data: any, context: MessageHandlerContext<ScreenType>) => {
  const { setCharacter, setCurrentScreen, setDungeonId, floorData, showToast } = context;
  
  console.log('Handling character_created message:', data);
  
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
    
    showToast({
      title: 'Character Created',
      description: data.message || 'Character created successfully',
      status: 'success',
    });
    
    // Move to dungeon selection screen
    console.log('Moving to dungeon selection screen');
    setCurrentScreen('dungeonSelection' as unknown as ScreenType);
    
    // If we already have floor data, transition to the game screen
    if (floorData && floorData.dungeonId) {
      console.log('We have floor data, transitioning to game screen with dungeonId:', floorData.dungeonId);
      setDungeonId(floorData.dungeonId);
      setCurrentScreen('game' as unknown as ScreenType);
    }
  } else {
    console.error('Character created but no ID received:', data);
    showToast({
      title: 'Error',
      description: 'Character created but no ID received',
      status: 'error',
    });
  }
};

// Handler for dungeon_joined messages
export const handleDungeonJoinedMessage = <ScreenType>(data: any, context: MessageHandlerContext<ScreenType>) => {
  const { setDungeonId, setFloorData, setCurrentScreen, character, currentScreen, showToast } = context;
  
  console.log('=== DUNGEON JOINED DEBUG ===');
  console.log('Handling dungeon_joined message:', data);
  console.log('Current character:', character);
  console.log('Current screen:', currentScreen);
  
  showToast({
    title: 'Success',
    description: data.message || 'Dungeon joined successfully',
    status: 'success',
  });
  
  // Set the dungeonId from the message
  if (data.dungeonId) {
    console.log('Setting dungeonId from dungeon_joined message:', data.dungeonId);
    setDungeonId(data.dungeonId);
    
    // Request floor data explicitly
    console.log('Requesting floor data after joining dungeon');
    getFloorDataWS(data.characterId, data.dungeonId);
    
    // If we don't receive floor data within 2 seconds, try to transition anyway
    setTimeout(() => {
      console.log('=== DUNGEON JOINED TIMEOUT CHECK ===');
      console.log('Current screen:', context.currentScreen);
      console.log('Character:', context.character);
      console.log('Floor data:', context.floorData);
      
      if (context.currentScreen !== 'game' && context.character && data.dungeonId) {
        console.log('Timeout: Transitioning to game screen with dungeonId:', data.dungeonId);
        
        // Create a minimal floor data object if we don't have one yet
        if (!context.floorData) {
          console.log('Creating minimal floor data as fallback');
          const minimalFloorData: FloorData = {
            floor: {
              width: 50,
              height: 50,
              tiles: [],
              entities: [],
              rooms: [],
              items: [],
              level: 1
            },
            playerPosition: { x: 0, y: 0 },
            currentFloor: 1,
            dungeonId: data.dungeonId,
            dungeonName: 'The Deeps'
          };
          setFloorData(minimalFloorData);
        }
        
        setCurrentScreen('game' as unknown as ScreenType);
      } else {
        console.log('Not transitioning to game screen after timeout:',
          'currentScreen =', context.currentScreen,
          'character =', context.character ? 'exists' : 'null',
          'dungeonId =', data.dungeonId || 'missing',
          'floorData =', context.floorData ? 'exists' : 'null');
      }
    }, 2000);
  } else {
    console.warn('DUNGEON JOINED ERROR: Message missing dungeonId');
  }
};

// Handler for dungeon_created messages
export const handleDungeonCreatedMessage = <ScreenType>(data: any, context: MessageHandlerContext<ScreenType>) => {
  const { showToast } = context;
  
  console.log('Handling dungeon_created message:', data);
  showToast({
    title: 'Success',
    description: data.message || 'Dungeon created successfully',
    status: 'success',
  });
};

// Handler for error messages
export const handleErrorMessage = <ScreenType>(data: any, context: MessageHandlerContext<ScreenType>) => {
  const { showToast } = context;
  
  console.error('Handling error message:', data);
  showToast({
    title: 'Error',
    description: data.message,
    status: 'error',
  });
};

// Main message handler that routes to specific handlers
export const handleWebSocketMessage = <ScreenType>(data: any, context: MessageHandlerContext<ScreenType>) => {
  console.log('WebSocket message received:', data);
  
  try {
    switch (data.type) {
      case 'floor_data':
        handleFloorDataMessage(data, context);
        break;
      case 'character_created':
        handleCharacterCreatedMessage(data, context);
        break;
      case 'dungeon_joined':
        handleDungeonJoinedMessage(data, context);
        break;
      case 'dungeon_created':
        handleDungeonCreatedMessage(data, context);
        break;
      case 'error':
        handleErrorMessage(data, context);
        break;
      case 'welcome':
        console.log('Received welcome message:', data.message);
        break;
      default:
        console.log('Unhandled message type:', data.type);
    }
  } catch (error) {
    console.error('Error handling WebSocket message:', error);
    context.showToast({
      title: 'Error',
      description: 'An error occurred while processing server message',
      status: 'error',
    });
  }
}; 