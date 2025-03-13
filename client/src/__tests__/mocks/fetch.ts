import { preparedRoomData } from './roomData';

// Mock implementation of fetch for tests
export const mockFetch = (url: string) => {
  // Extract room type from URL
  const urlObj = new URL(url, 'http://localhost');
  const roomType = urlObj.searchParams.get('type') || 'entrance';
  
  // Get the appropriate room data based on the room type
  const roomData = preparedRoomData[roomType as keyof typeof preparedRoomData] || preparedRoomData.entrance;
  
  return Promise.resolve({
    ok: true,
    status: 200,
    json: () => Promise.resolve(roomData)
  });
};

// Mock implementation for invalid requests
export const mockFetchError = () => {
  return Promise.resolve({
    ok: false,
    status: 500,
    statusText: 'Internal Server Error',
    json: () => Promise.reject(new Error('Failed to fetch'))
  });
};

// Setup global fetch mock
export const setupFetchMock = () => {
  global.fetch = jest.fn().mockImplementation(mockFetch);
};

// Setup global fetch mock that returns an error
export const setupFetchErrorMock = () => {
  global.fetch = jest.fn().mockImplementation(mockFetchError);
};

// Reset fetch mock
export const resetFetchMock = () => {
  if (global.fetch && typeof jest !== 'undefined') {
    (global.fetch as jest.Mock).mockReset();
  }
}; 