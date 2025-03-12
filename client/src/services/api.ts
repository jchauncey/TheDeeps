import axios from 'axios';
import { Character, CharacterClass, Attributes, Dungeon } from '../types';

const API_URL = process.env.REACT_APP_API_URL || '';

// Add request interceptor for debugging
axios.interceptors.request.use(
  config => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  error => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

// Add response interceptor for debugging
axios.interceptors.response.use(
  response => {
    console.log(`API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  error => {
    console.error('API Response Error:', error.response?.status, error.response?.data || error.message);
    return Promise.reject(error);
  }
);

// Character API
export const getCharacters = async (): Promise<Character[]> => {
  try {
    const response = await axios.get(`${API_URL}/characters`);
    return response.data;
  } catch (error) {
    console.error('Failed to get characters:', error);
    throw error;
  }
};

export const getCharacter = async (id: string): Promise<Character> => {
  try {
    const response = await axios.get(`${API_URL}/characters/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Failed to get character ${id}:`, error);
    throw error;
  }
};

interface CreateCharacterRequest {
  name: string;
  class: CharacterClass;
  attributes?: Attributes;
}

export const createCharacter = async (character: CreateCharacterRequest): Promise<Character> => {
  try {
    const response = await axios.post(`${API_URL}/characters`, character);
    return response.data;
  } catch (error) {
    console.error('Failed to create character:', error);
    throw error;
  }
};

export const deleteCharacter = async (id: string): Promise<void> => {
  try {
    await axios.delete(`${API_URL}/characters/${id}`);
  } catch (error) {
    console.error(`Failed to delete character ${id}:`, error);
    throw error;
  }
};

// Dungeon API
export const getDungeons = async (): Promise<Dungeon[]> => {
  try {
    const response = await axios.get(`${API_URL}/dungeons`);
    return response.data;
  } catch (error) {
    console.error('Failed to get dungeons:', error);
    throw error;
  }
};

interface CreateDungeonRequest {
  name: string;
  floors: number;
  difficulty: string;
}

export const createDungeon = async (dungeon: CreateDungeonRequest): Promise<Dungeon> => {
  try {
    const response = await axios.post(`${API_URL}/dungeons`, dungeon);
    return response.data;
  } catch (error) {
    console.error('Failed to create dungeon:', error);
    throw error;
  }
};

export const joinDungeon = async (characterId: string, dungeonId: string): Promise<void> => {
  try {
    await axios.post(`${API_URL}/dungeons/${dungeonId}/join`, { characterId });
  } catch (error) {
    console.error(`Failed to join dungeon ${dungeonId}:`, error);
    throw error;
  }
}; 