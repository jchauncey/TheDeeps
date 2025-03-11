import axios from 'axios';
import { Character, CharacterClass, Attributes } from '../types';

const API_URL = process.env.REACT_APP_API_URL || '';

// Character API
export const getCharacters = async (): Promise<Character[]> => {
  const response = await axios.get(`${API_URL}/characters`);
  return response.data;
};

export const getCharacter = async (id: string): Promise<Character> => {
  const response = await axios.get(`${API_URL}/characters/${id}`);
  return response.data;
};

interface CreateCharacterRequest {
  name: string;
  class: CharacterClass;
  attributes?: Attributes;
}

export const createCharacter = async (character: CreateCharacterRequest): Promise<Character> => {
  const response = await axios.post(`${API_URL}/characters`, character);
  return response.data;
};

export const deleteCharacter = async (id: string): Promise<void> => {
  await axios.delete(`${API_URL}/characters/${id}`);
}; 