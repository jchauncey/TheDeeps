import { preparedRoomData } from './roomData';

// Mock implementation of fetch for tests
export const mockFetch = (url: string) => {
  // Extract room type from URL
  const urlObj = new URL(url, 'http://localhost');
  const roomType = urlObj.searchParams.get('type') || 'entrance';
  
  // Get the appropriate room data based on the room type
  const roomData = preparedRoomData[roomType as keyof typeof preparedRoomData] || preparedRoomData.entrance;
  
  // Create a proper response object with clone method
  const responseText = JSON.stringify(roomData);
  
  return Promise.resolve({
    ok: true,
    status: 200,
    statusText: 'OK',
    json: () => Promise.resolve(roomData),
    text: () => Promise.resolve(responseText),
    headers: new Headers(),
    clone: function() {
      return {
        ok: this.ok,
        status: this.status,
        statusText: this.statusText,
        json: this.json,
        text: () => Promise.resolve(responseText),
        headers: this.headers,
        clone: this.clone
      };
    }
  });
};

// Mock implementation for invalid requests
export const mockFetchError = () => {
  return Promise.resolve({
    ok: false,
    status: 500,
    statusText: 'Internal Server Error',
    json: () => Promise.reject(new Error('Failed to fetch')),
    text: () => Promise.resolve('Internal Server Error'),
    headers: new Headers(),
    clone: function() {
      return {
        ok: this.ok,
        status: this.status,
        statusText: this.statusText,
        json: this.json,
        text: () => Promise.resolve('Internal Server Error'),
        headers: this.headers,
        clone: this.clone
      };
    }
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