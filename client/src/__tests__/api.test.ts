import { getCharacter, joinDungeon, mockCharacters } from './mocks/api';
import { Character } from '../types';

describe('API Service', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('getCharacter', () => {
    it('should fetch a character by ID', async () => {
      // Call the function with an existing ID
      const result = await getCharacter('1');
      
      // Assertions
      expect(result).toEqual(mockCharacters[0]);
    });

    it('should throw an error when character fetch fails', async () => {
      // Call and expect rejection
      await expect(getCharacter('non-existent-id')).rejects.toThrow('Character not found');
    });
  });

  describe('joinDungeon', () => {
    it('should successfully join a dungeon with valid IDs', async () => {
      // Call the function with valid IDs
      await expect(joinDungeon('1', 'dungeon-1')).resolves.not.toThrow();
    });

    it('should throw a specific error when character is not found', async () => {
      // Call and expect specific error
      try {
        await joinDungeon('non-existent-id', 'dungeon-1');
        fail('Should have thrown an error');
      } catch (error: any) {
        expect(error.isAxiosError).toBe(true);
        expect(error.response.status).toBe(404);
        expect(error.response.data).toBe('character not found');
      }
    });

    it('should throw a specific error when dungeon is not found', async () => {
      // Call and expect specific error
      try {
        await joinDungeon('1', 'non-existent-id');
        fail('Should have thrown an error');
      } catch (error: any) {
        expect(error.isAxiosError).toBe(true);
        expect(error.response.status).toBe(404);
        expect(error.response.data).toBe('dungeon not found');
      }
    });
  });
}); 