import { CharacterData } from '../types/game';

// Base URL for API calls
const API_BASE_URL = 'http://localhost:8080';

// WebSocket connection
let ws: WebSocket | null = null;
let isConnecting = false;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const reconnectDelay = 2000; // 2 seconds
let pendingMessages: any[] = []; // Store messages that couldn't be sent due to connection issues

// Connect to WebSocket
export const connectWebSocket = (onMessage: (event: Event) => void): WebSocket => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    console.log('WebSocket already connected');
    return ws;
  }
  
  if (isConnecting) {
    console.log('WebSocket connection already in progress');
    return ws as WebSocket;
  }
  
  isConnecting = true;
  
  console.log('Connecting to WebSocket...');
  ws = new WebSocket('ws://localhost:8080/ws');
  
  ws.onopen = () => {
    console.log('WebSocket connection established');
    isConnecting = false;
    reconnectAttempts = 0;
    
    // Send any pending messages
    if (pendingMessages.length > 0) {
      console.log(`Sending ${pendingMessages.length} pending messages`);
      pendingMessages.forEach(message => {
        try {
          ws?.send(JSON.stringify(message));
        } catch (error) {
          console.error('Error sending pending message:', error);
        }
      });
      pendingMessages = [];
    }
    
    // Create a custom event to notify the app that the connection is established
    const event = new CustomEvent('websocket_connected');
    window.dispatchEvent(event);
  };
  
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      console.log('WebSocket message received:', data);
      
      // Create a custom event with the parsed data
      const customEvent = new CustomEvent('websocket_message', { detail: data });
      
      // Dispatch the event to the window
      window.dispatchEvent(customEvent);
      
      // Also call the provided callback
      onMessage(customEvent);
    } catch (error) {
      console.error('Error parsing WebSocket message:', error);
    }
  };
  
  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
    isConnecting = false;
  };
  
  ws.onclose = (event) => {
    console.log(`WebSocket connection closed: ${event.code} ${event.reason}`);
    isConnecting = false;
    ws = null;
    
    // Attempt to reconnect if not a normal closure
    if (reconnectAttempts < maxReconnectAttempts) {
      reconnectAttempts++;
      console.log(`Attempting to reconnect (${reconnectAttempts}/${maxReconnectAttempts})...`);
      setTimeout(() => connectWebSocket(onMessage), reconnectDelay);
    } else if (reconnectAttempts >= maxReconnectAttempts) {
      console.error('Maximum reconnect attempts reached. Please refresh the page.');
      
      // Create a custom event to notify the app that reconnection failed
      const event = new CustomEvent('websocket_reconnect_failed');
      window.dispatchEvent(event);
    }
  };
  
  return ws;
};

// Send message through WebSocket
export const sendWebSocketMessage = (message: any): boolean => {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    console.log('WebSocket not connected, storing message and attempting to connect...');
    // Store the message to send when connection is established
    pendingMessages.push(message);
    
    ws = connectWebSocket((event) => {
      console.log('Reconnected, pending messages will be sent automatically');
    });
    return false;
  }
  
  try {
    console.log('Sending WebSocket message:', message);
    ws.send(JSON.stringify(message));
    return true;
  } catch (error) {
    console.error('Error sending WebSocket message:', error);
    // Store the message to try again later
    pendingMessages.push(message);
    return false;
  }
};

// Create a new character
export const createCharacter = async (character: CharacterData): Promise<{ success: boolean; message?: string }> => {
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
    
    return { success: true, message: 'Character created successfully' };
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
export const getSavedCharacters = async (): Promise<{ success: boolean; characters?: { id: string; name: string }[]; message?: string }> => {
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