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

// Connect to WebSocket
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
      if (onMessageCallback) {
        onMessageCallback(event);
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
    console.error('Error creating WebSocket connection:', error);
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
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    console.warn('WebSocket not connected, queueing message');
    messageQueue.push(message);
    
    // Try to reconnect if not already reconnecting
    if (!isReconnecting && (!ws || ws.readyState === WebSocket.CLOSED)) {
      if (reconnectAttempts < maxReconnectAttempts) {
        handleReconnect();
      }
    }
    
    return false;
  }
  
  try {
    ws.send(JSON.stringify(message));
    return true;
  } catch (error) {
    console.error('Error sending WebSocket message:', error);
    return false;
  }
};

// Check if WebSocket is connected
export const isWebSocketConnected = (): boolean => {
  return ws !== null && ws.readyState === WebSocket.OPEN;
};

// Manually close the WebSocket connection
export const closeWebSocketConnection = (): void => {
  if (ws) {
    ws.close(1000, "Client closing connection");
    ws = null;
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

// Create a new character
export const createCharacter = async (character: CharacterData): Promise<{ success: boolean; characterId?: string; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/character`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(character),
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.message || 'Failed to create character');
    }
    
    return { success: true, characterId: data.id, message: 'Character created successfully' };
  } catch (error) {
    console.error('Error creating character:', error);
    return { 
      success: false, 
      message: error instanceof Error ? error.message : 'An unknown error occurred' 
    };
  }
};

// Load a character
export const loadCharacter = async (characterId: string): Promise<{ success: boolean; character?: CharacterData; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/character/${characterId}`);
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.message || 'Failed to load character');
    }
    
    return { success: true, character: data };
  } catch (error) {
    console.error('Error loading character:', error);
    return { 
      success: false, 
      message: error instanceof Error ? error.message : 'An unknown error occurred' 
    };
  }
};

// Get list of saved characters
export const getSavedCharacters = async (): Promise<{ success: boolean; characters?: { id: string; name: string; characterClass: string }[]; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/characters`);
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.message || 'Failed to get saved characters');
    }
    
    return { success: true, characters: data };
  } catch (error) {
    console.error('Error getting saved characters:', error);
    return { 
      success: false, 
      message: error instanceof Error ? error.message : 'An unknown error occurred' 
    };
  }
};

// Save the current game
export const saveGame = async (characterId: string): Promise<{ success: boolean; message?: string }> => {
  try {
    // Send a WebSocket message to save the game
    const success = sendWebSocketMessage({
      type: 'action',
      action: 'save_game',
      characterId
    });
    
    if (!success) {
      throw new Error('Failed to send save game request');
    }
    
    return { success: true, message: 'Game save request sent' };
  } catch (error) {
    console.error('Error saving game:', error);
    return { 
      success: false, 
      message: error instanceof Error ? error.message : 'An unknown error occurred' 
    };
  }
};

// Create a new dungeon
export const createDungeon = async (name: string, numFloors: number): Promise<{ success: boolean; dungeonId?: string; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/dungeon`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name, numFloors }),
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.message || 'Failed to create dungeon');
    }
    
    return { success: true, dungeonId: data.id, message: 'Dungeon created successfully' };
  } catch (error) {
    console.error('Error creating dungeon:', error);
    return { 
      success: false, 
      message: error instanceof Error ? error.message : 'An unknown error occurred' 
    };
  }
};

// Get list of available dungeons
export const getAvailableDungeons = async (): Promise<{ success: boolean; dungeons?: DungeonData[]; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/dungeons`);
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.message || 'Failed to get available dungeons');
    }
    
    return { success: true, dungeons: data };
  } catch (error) {
    console.error('Error getting available dungeons:', error);
    return { 
      success: false, 
      message: error instanceof Error ? error.message : 'An unknown error occurred' 
    };
  }
};

// Join an existing dungeon
export const joinDungeon = async (dungeonId: string, characterId: string): Promise<{ success: boolean; floorData?: FloorData; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/dungeon/${dungeonId}/join?characterId=${characterId}`, {
      method: 'POST',
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.message || 'Failed to join dungeon');
    }
    
    return { success: true, floorData: data, message: 'Joined dungeon successfully' };
  } catch (error) {
    console.error('Error joining dungeon:', error);
    return { 
      success: false, 
      message: error instanceof Error ? error.message : 'An unknown error occurred' 
    };
  }
};

// Get a specific floor of a dungeon
export const getDungeonFloor = async (dungeonId: string, level: number, characterId?: string): Promise<{ success: boolean; floor?: FloorData; message?: string }> => {
  try {
    let url = `${API_BASE_URL}/dungeon/${dungeonId}/floor/${level}`;
    if (characterId) {
      url += `?characterId=${characterId}`;
    }
    
    const response = await fetch(url);
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.message || 'Failed to get dungeon floor');
    }
    
    return { success: true, floor: data };
  } catch (error) {
    console.error('Error getting dungeon floor:', error);
    return { 
      success: false, 
      message: error instanceof Error ? error.message : 'An unknown error occurred' 
    };
  }
};

// Get the current floor of a character
export const getCharacterFloor = async (characterId: string): Promise<{ success: boolean; floorData?: FloorData; message?: string }> => {
  try {
    const response = await fetch(`${API_BASE_URL}/character/${characterId}/floor`);
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.message || 'Failed to get character floor');
    }
    
    return { success: true, floorData: data };
  } catch (error) {
    console.error('Error getting character floor:', error);
    return { 
      success: false, 
      message: error instanceof Error ? error.message : 'An unknown error occurred' 
    };
  }
};

// WebSocket message types for dungeon management
export const createDungeonWS = (name: string, numFloors: number): boolean => {
  return sendWebSocketMessage({
    type: 'create_dungeon',
    name,
    numFloors
  });
};

export const joinDungeonWS = (dungeonId: string, characterId: string): boolean => {
  return sendWebSocketMessage({
    type: 'join_dungeon',
    dungeonId,
    characterId
  });
};

export const listDungeonsWS = (): boolean => {
  return sendWebSocketMessage({
    type: 'list_dungeons'
  });
};

// WebSocket message types for game actions
export const moveCharacter = (direction: string): boolean => {
  return sendWebSocketMessage({
    type: 'move',
    direction
  });
};

export const performAction = (action: string): boolean => {
  return sendWebSocketMessage({
    type: 'action',
    action
  });
};

export const getFloorData = (): boolean => {
  return sendWebSocketMessage({
    type: 'get_floor'
  });
}; 