import { CharacterData, DungeonData, FloorData } from '../types/game';

// Base URL for API calls
const API_BASE_URL = 'http://localhost:8080';

// WebSocket connection
let ws: WebSocket | null = null;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const reconnectDelay = 2000; // 2 seconds
let reconnectTimer: number | null = null;
let isReconnecting = false;
let messageQueue: any[] = [];

// Event callbacks
let onMessageCallback: ((event: Event) => void) | null = null;
let onDisconnectCallback: (() => void) | null = null;
let onReconnectFailedCallback: (() => void) | null = null;
let onConnectedCallback: (() => void) | null = null;

// Floor data cache to prevent unnecessary re-renders
let lastFloorData: any = null;
let lastMovementTime = 0;
const MOVEMENT_THROTTLE = 30; // milliseconds - reduced for more responsive movement

// Track last movement request time for throttling
let lastMovementRequestTime = 0;
const MOVEMENT_REQUEST_THROTTLE = 50; // milliseconds - reduced for more responsive movement

// Track what type of update we're receiving
let isProcessingMovementUpdate = false;

// Connect to WebSocket - only used for in-dungeon gameplay
export const connectWebSocket = (onMessage: (event: Event) => void): WebSocket | null => {
  // Store the callback
  onMessageCallback = onMessage;
  
  // If already connected, return the existing connection
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
    return ws;
  }
  
  try {
    // Create WebSocket connection
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.hostname}:8080/ws`;
    
    console.log(`Connecting to WebSocket at ${wsUrl}`);
    ws = new WebSocket(wsUrl);
    
    // Set up event handlers
    ws.onopen = () => {
      console.log('WebSocket connection established');
      reconnectAttempts = 0;
      isReconnecting = false;
      
      // Process any queued messages
      if (messageQueue.length > 0) {
        console.log(`Processing ${messageQueue.length} queued messages`);
        messageQueue.forEach(msg => {
          if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify(msg));
          }
        });
        messageQueue = [];
      }
      
      // Call the connected callback if provided
      if (onConnectedCallback) {
        onConnectedCallback();
      }
    };
    
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        
        // Process floor data to ensure smooth movement
        if (data.type === 'floor_data') {
          // Determine if this is just a movement update or a full floor update
          isProcessingMovementUpdate = false;
          
          if (lastFloorData) {
            // Check if only the player position changed
            const playerPosChanged = 
              lastFloorData.playerPosition.x !== data.playerPosition.x || 
              lastFloorData.playerPosition.y !== data.playerPosition.y;
              
            // Check if entities changed
            const entitiesChanged = JSON.stringify(lastFloorData.floor.entities) !== 
                                   JSON.stringify(data.floor.entities);
                                   
            // Check if items changed
            const itemsChanged = JSON.stringify(lastFloorData.floor.items) !== 
                                JSON.stringify(data.floor.items);
                                
            // Check if floor level changed
            const floorLevelChanged = lastFloorData.currentFloor !== data.currentFloor;
            
            // If only player position changed, mark as movement update
            if (playerPosChanged && !entitiesChanged && !itemsChanged && !floorLevelChanged) {
              isProcessingMovementUpdate = true;
              
              // For movement updates, we only need to update the player position
              // Create a simplified update with just the necessary data
              const simplifiedUpdate = {
                type: 'floor_data',
                playerPosition: data.playerPosition,
                currentFloor: data.currentFloor,
                playerData: data.playerData,
                dungeonId: data.dungeonId,
                isMovementUpdate: true
              };
              
              // Update the last floor data with the new player position
              lastFloorData.playerPosition = data.playerPosition;
              
              // Call the message callback with the simplified update
              if (onMessageCallback) {
                // Create a new event with the simplified data
                const simplifiedEvent = new MessageEvent('message', {
                  data: JSON.stringify(simplifiedUpdate)
                });
                onMessageCallback(simplifiedEvent);
              }
              
              return;
            }
          }
          
          // For full floor updates or first load
          console.log(`Received ${isProcessingMovementUpdate ? 'movement update' : 'full floor data'} for floor ${data.currentFloor}`);
          lastFloorData = data;
          
          // Call the message callback
          if (onMessageCallback) {
            onMessageCallback(event);
          }
        } else {
          // For non-floor data messages, always process
          if (onMessageCallback) {
            onMessageCallback(event);
          }
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
    
    ws.onclose = (event) => {
      console.log(`WebSocket connection closed: ${event.code} - ${event.reason}`);
      
      // Don't attempt to reconnect on normal closure (code 1000) or if going away (code 1001)
      // unless it's a server-side issue
      const shouldReconnect = !(event.wasClean && (event.code === 1000 || event.code === 1001));
      
      if (shouldReconnect && reconnectAttempts < maxReconnectAttempts) {
        handleReconnect();
      } else if (shouldReconnect) {
        console.error('Max reconnection attempts reached');
        if (onReconnectFailedCallback) {
          onReconnectFailedCallback();
        }
      } else if (onDisconnectCallback) {
        // Normal closure, call disconnect callback
        onDisconnectCallback();
      }
    };
    
    return ws;
  } catch (error) {
    console.error('Error connecting to WebSocket:', error);
    return null;
  }
};

// Handle reconnection logic
const handleReconnect = () => {
  if (isReconnecting) return;
  
  isReconnecting = true;
  reconnectAttempts++;
  
  console.log(`Attempting to reconnect (${reconnectAttempts}/${maxReconnectAttempts}) in ${reconnectDelay}ms`);
  
  if (reconnectTimer) {
    window.clearTimeout(reconnectTimer);
  }
  
  reconnectTimer = window.setTimeout(() => {
    console.log(`Reconnecting... Attempt ${reconnectAttempts}`);
    connectWebSocket(onMessageCallback!);
    isReconnecting = false;
  }, reconnectDelay * reconnectAttempts); // Exponential backoff
};

// Set callbacks for connection events
export const setWebSocketCallbacks = (
  onDisconnect?: () => void,
  onReconnectFailed?: () => void,
  onConnected?: () => void
) => {
  if (onDisconnect) onDisconnectCallback = onDisconnect;
  if (onReconnectFailed) onReconnectFailedCallback = onReconnectFailed;
  if (onConnected) onConnectedCallback = onConnected;
};

// Send a message through the WebSocket
export const sendWebSocketMessage = (message: any): boolean => {
  // Throttle movement requests
  if (message.type === 'move') {
    const now = performance.now();
    if (now - lastMovementRequestTime < MOVEMENT_REQUEST_THROTTLE) {
      // Skip this request if it's too soon after the last one
      return false;
    }
    lastMovementRequestTime = now;
  }

  if (!ws) {
    console.log('WebSocket not connected, queueing message:', message);
    messageQueue.push(message);
    return false;
  }
  
  if (ws.readyState !== WebSocket.OPEN) {
    console.log('WebSocket not open, queueing message:', message);
    messageQueue.push(message);
    return false;
  }
  
  try {
    console.log('Sending WebSocket message:', message);
    ws.send(JSON.stringify(message));
    return true;
  } catch (error) {
    console.error('Error sending WebSocket message:', error);
    return false;
  }
};

// Check if WebSocket is connected
export const isWebSocketConnected = (): boolean => {
  const connected = ws !== null && ws.readyState === WebSocket.OPEN;
  console.log('WebSocket connection status:', connected, 'readyState:', ws ? ws.readyState : 'null');
  return connected;
};

// Manually close the WebSocket connection
export const closeWebSocketConnection = (): void => {
  if (ws) {
    try {
      // Set a short timeout to ensure the connection is closed
      ws.onclose = () => {
        console.log('WebSocket connection closed by client');
        ws = null;
      };
      
      // Close the connection with a normal closure code
      ws.close(1000, "Client closing connection");
      
      // If the close event doesn't fire within 100ms, force it
      setTimeout(() => {
        if (ws) {
          console.log('Forcing WebSocket cleanup');
          ws = null;
        }
      }, 100);
    } catch (error) {
      console.error('Error closing WebSocket connection:', error);
      ws = null;
    }
  }
  
  // Clear any pending reconnect timers
  if (reconnectTimer) {
    window.clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }
  
  reconnectAttempts = 0;
  isReconnecting = false;
  messageQueue = [];
};

// REST API functions for static initialization

/**
 * Get all saved characters for the player
 */
export const getSavedCharacters = async (): Promise<{ success: boolean; characters?: { id: string; name: string; characterClass: string }[]; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/characters`);
    
    if (!response.ok) {
      const errorData = await response.json();
      return { 
        success: false, 
        message: errorData.message || `Error: ${response.status} ${response.statusText}` 
      };
    }
    
    const data = await response.json();
    return { success: true, characters: data };
  } catch (error) {
    console.error('Error fetching saved characters:', error);
    return { success: false, message: 'Failed to fetch saved characters' };
  }
};

/**
 * Create a new character using REST API
 */
export const createCharacter = async (character: CharacterData): Promise<{ success: boolean; characterId?: string; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/characters`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: character.name,
        characterClass: character.characterClass,
        stats: character.stats,
      }),
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      return { 
        success: false, 
        message: errorData.message || `Error: ${response.status} ${response.statusText}` 
      };
    }
    
    const data = await response.json();
    return { success: true, characterId: data.id };
  } catch (error) {
    console.error('Error creating character:', error);
    return { success: false, message: 'Failed to create character' };
  }
};

/**
 * Load a character by ID using REST API
 */
export const loadCharacter = async (characterId: string): Promise<{ success: boolean; character?: CharacterData; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/characters/${characterId}`);
    
    if (!response.ok) {
      const errorData = await response.json();
      return { 
        success: false, 
        message: errorData.message || `Error: ${response.status} ${response.statusText}` 
      };
    }
    
    const character = await response.json();
    return { success: true, character };
  } catch (error) {
    console.error('Error loading character:', error);
    return { success: false, message: 'Failed to load character' };
  }
};

/**
 * Get available dungeons using REST API
 */
export const getAvailableDungeons = async (): Promise<{ success: boolean; dungeons?: DungeonData[]; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/dungeons`);
    
    if (!response.ok) {
      const errorData = await response.json();
      return { 
        success: false, 
        message: errorData.message || `Error: ${response.status} ${response.statusText}` 
      };
    }
    
    const dungeons = await response.json();
    return { success: true, dungeons };
  } catch (error) {
    console.error('Error fetching available dungeons:', error);
    return { success: false, message: 'Failed to fetch available dungeons' };
  }
};

/**
 * Create a new dungeon using REST API
 */
export const createDungeon = async (name: string, numFloors: number): Promise<{ success: boolean; dungeonId?: string; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/dungeons`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name,
        numFloors,
      }),
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      return { 
        success: false, 
        message: errorData.message || `Error: ${response.status} ${response.statusText}` 
      };
    }
    
    const data = await response.json();
    return { success: true, dungeonId: data.id };
  } catch (error) {
    console.error('Error creating dungeon:', error);
    return { success: false, message: 'Failed to create dungeon' };
  }
};

/**
 * Join a dungeon using REST API
 */
export const joinDungeon = async (dungeonId: string, characterId: string): Promise<{ success: boolean; floorData?: FloorData; message?: string }> => {
  try {
    console.log(`Joining dungeon ${dungeonId} with character ${characterId}`);
    
    const response = await fetch(`${API_BASE_URL}/dungeons/${dungeonId}/join`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        characterId,
      }),
    });
    
    const responseText = await response.text();
    console.log(`Join dungeon response: ${response.status} ${response.statusText}`, responseText);
    
    if (!response.ok) {
      let errorMessage = `Error: ${response.status} ${response.statusText}`;
      try {
        const errorData = JSON.parse(responseText);
        errorMessage = errorData.message || errorMessage;
      } catch (e) {
        // If parsing fails, use the raw response text
        errorMessage = responseText || errorMessage;
      }
      
      console.error('Error joining dungeon:', errorMessage);
      return { 
        success: false, 
        message: errorMessage
      };
    }
    
    let floorData: FloorData;
    try {
      floorData = JSON.parse(responseText);
      console.log('Parsed floor data:', floorData);
      
      // Validate the floor data
      if (!floorData.floor || !floorData.playerPosition) {
        throw new Error('Invalid floor data received from server');
      }
    } catch (e) {
      console.error('Error parsing floor data:', e);
      return { 
        success: false, 
        message: 'Invalid floor data received from server' 
      };
    }
    
    // After joining the dungeon via REST, connect to WebSocket for real-time updates
    if (!isWebSocketConnected()) {
      console.log('Connecting to WebSocket for real-time updates');
      connectWebSocket(onMessageCallback || (() => {}));
    }
    
    // Send a WebSocket message to associate the character with the WebSocket connection
    // This is crucial for the server to know which character is controlled by this connection
    console.log('Sending identify_character WebSocket message with characterId:', characterId);
    identifyCharacter(characterId);
    
    return { success: true, floorData };
  } catch (error) {
    console.error('Error joining dungeon:', error);
    return { success: false, message: error instanceof Error ? error.message : 'Failed to join dungeon' };
  }
};

/**
 * Get a specific floor of a dungeon
 */
export const getDungeonFloor = async (dungeonId: string, level: number, characterId?: string): Promise<{ success: boolean; floor?: FloorData; message?: string }> => {
  try {
    let url = `${API_BASE_URL}/dungeons/${dungeonId}/floor/${level}`;
    if (characterId) {
      url += `?characterId=${characterId}`;
    }
    
    const response = await fetch(url);
    
    if (!response.ok) {
      const errorData = await response.json();
      return { 
        success: false, 
        message: errorData.message || `Error: ${response.status} ${response.statusText}` 
      };
    }
    
    const floor = await response.json();
    return { success: true, floor };
  } catch (error) {
    console.error('Error getting dungeon floor:', error);
    return { success: false, message: 'Failed to get dungeon floor' };
  }
};

/**
 * Get the current floor of a character
 */
export const getCharacterFloor = async (characterId: string): Promise<{ success: boolean; floorData?: FloorData; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/characters/${characterId}/floor`);
    
    if (!response.ok) {
      const errorData = await response.json();
      return { 
        success: false, 
        message: errorData.message || `Error: ${response.status} ${response.statusText}` 
      };
    }
    
    const floorData = await response.json();
    return { success: true, floorData };
  } catch (error) {
    console.error('Error getting character floor:', error);
    return { success: false, message: 'Failed to get character floor' };
  }
};

/**
 * Save the current game
 */
export const saveGame = async (characterId: string): Promise<{ success: boolean; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/characters/${characterId}/save`, {
      method: 'POST',
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      return { 
        success: false, 
        message: errorData.message || `Error: ${response.status} ${response.statusText}` 
      };
    }
    
    return { success: true, message: 'Game saved successfully' };
  } catch (error) {
    console.error('Error saving game:', error);
    return { success: false, message: 'Failed to save game' };
  }
};

// WebSocket-specific game actions (only for in-dungeon gameplay)

/**
 * Move character in a direction
 */
export const moveCharacter = (direction: string): boolean => {
  return sendWebSocketMessage({ type: 'move', direction });
};

/**
 * Perform an action
 */
export const performAction = (action: string): boolean => {
  return sendWebSocketMessage({ type: 'action', action });
};

/**
 * Request floor data
 */
export const getFloorData = (): boolean => {
  return sendWebSocketMessage({ type: 'get_floor_data' });
};

// Deprecated WebSocket methods - kept for backward compatibility
// These should be replaced with REST API calls in new code

/**
 * @deprecated Use createCharacter REST API instead
 */
export const createCharacterWS = (character: {
  name: string;
  characterClass: string;
  stats: any;
  abilities: string[];
}): boolean => {
  console.warn('createCharacterWS is deprecated. Use createCharacter REST API instead.');
  return sendWebSocketMessage({ type: 'create_character', ...character });
};

/**
 * @deprecated Use createDungeon REST API instead
 */
export const createDungeonWS = (name: string, numFloors: number): boolean => {
  console.warn('createDungeonWS is deprecated. Use createDungeon REST API instead.');
  return sendWebSocketMessage({ type: 'create_dungeon', name, numFloors });
};

/**
 * @deprecated Use joinDungeon REST API instead
 */
export const joinDungeonWS = (dungeonId: string, characterId: string): boolean => {
  console.warn('joinDungeonWS is deprecated. Use joinDungeon REST API instead.');
  return sendWebSocketMessage({ type: 'join_dungeon', dungeonId, characterId });
};

/**
 * @deprecated Use getAvailableDungeons REST API instead
 */
export const listDungeonsWS = (): boolean => {
  console.warn('listDungeonsWS is deprecated. Use getAvailableDungeons REST API instead.');
  return sendWebSocketMessage({ type: 'list_dungeons' });
};

/**
 * Delete a character by ID
 */
export const deleteCharacter = async (characterId: string): Promise<{ success: boolean; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/characters/${characterId}`, {
      method: 'DELETE',
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      return { 
        success: false, 
        message: errorData.message || `Error: ${response.status} ${response.statusText}` 
      };
    }
    
    return { success: true };
  } catch (error) {
    console.error('Error deleting character:', error);
    return { success: false, message: 'Failed to delete character' };
  }
};

/**
 * Identify the character with the WebSocket connection
 * This is crucial for the server to know which character is controlled by this connection
 */
export const identifyCharacter = (characterId: string): boolean => {
  console.log('Identifying character with WebSocket connection:', characterId);
  return sendWebSocketMessage({ 
    type: 'identify_character', 
    characterId 
  });
}; 