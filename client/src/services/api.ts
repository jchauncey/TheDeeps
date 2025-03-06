import { CharacterData } from '../types/game';

// Base URL for API calls
const API_BASE_URL = 'http://localhost:8080';

// WebSocket connection
let ws: WebSocket | null = null;

// Connect to WebSocket
export const connectWebSocket = (onMessage: (data: any) => void): WebSocket => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    return ws;
  }

  ws = new WebSocket('ws://localhost:8080/ws');
  
  ws.onopen = () => {
    console.log('WebSocket connection established');
  };
  
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      onMessage(data);
    } catch (error) {
      console.error('Error parsing WebSocket message:', error);
    }
  };
  
  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
  };
  
  ws.onclose = () => {
    console.log('WebSocket connection closed');
    ws = null;
  };
  
  return ws;
};

// Send message through WebSocket
export const sendWebSocketMessage = (message: any): boolean => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(message));
    return true;
  }
  return false;
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