// Mock data for different room types
export const mockRoomData = {
  entrance: {
    level: 1,
    width: 20,
    height: 20,
    tiles: Array(20).fill(null).map(() => 
      Array(20).fill(null).map(() => ({
        type: 'wall',
        walkable: false,
        explored: false
      }))
    ),
    rooms: [
      {
        id: 'test-room',
        type: 'entrance',
        x: 6,
        y: 6,
        width: 8,
        height: 8,
        explored: true
      }
    ],
    upStairs: [],
    downStairs: [{ x: 12, y: 12 }],
    mobs: {},
    items: {}
  },
  
  standard: {
    level: 1,
    width: 20,
    height: 20,
    tiles: Array(20).fill(null).map(() => 
      Array(20).fill(null).map(() => ({
        type: 'wall',
        walkable: false,
        explored: false
      }))
    ),
    rooms: [
      {
        id: 'test-room',
        type: 'standard',
        x: 6,
        y: 6,
        width: 7,
        height: 7,
        explored: true
      }
    ],
    upStairs: [],
    downStairs: [],
    mobs: {},
    items: {}
  },
  
  treasure: {
    level: 1,
    width: 20,
    height: 20,
    tiles: Array(20).fill(null).map(() => 
      Array(20).fill(null).map(() => ({
        type: 'wall',
        walkable: false,
        explored: false
      }))
    ),
    rooms: [
      {
        id: 'test-room',
        type: 'treasure',
        x: 6,
        y: 6,
        width: 7,
        height: 7,
        explored: true
      }
    ],
    upStairs: [],
    downStairs: [],
    mobs: {},
    items: {
      'test-item-0': {
        id: 'test-item-0',
        type: 'weapon',
        name: 'Test Item 0',
        position: { x: 8, y: 8 }
      },
      'test-item-1': {
        id: 'test-item-1',
        type: 'weapon',
        name: 'Test Item 1',
        position: { x: 10, y: 8 }
      }
    }
  },
  
  boss: {
    level: 1,
    width: 20,
    height: 20,
    tiles: Array(20).fill(null).map(() => 
      Array(20).fill(null).map(() => ({
        type: 'wall',
        walkable: false,
        explored: false
      }))
    ),
    rooms: [
      {
        id: 'test-room',
        type: 'boss',
        x: 6,
        y: 6,
        width: 7,
        height: 7,
        explored: true
      }
    ],
    upStairs: [],
    downStairs: [],
    mobs: {
      'test-boss': {
        id: 'test-boss',
        type: 'dragon',
        name: 'Dragon Boss',
        health: 100,
        maxHealth: 100,
        position: { x: 9, y: 9 },
        variant: 'boss'
      }
    },
    items: {}
  },
  
  safe: {
    level: 1,
    width: 20,
    height: 20,
    tiles: Array(20).fill(null).map(() => 
      Array(20).fill(null).map(() => ({
        type: 'wall',
        walkable: false,
        explored: false
      }))
    ),
    rooms: [
      {
        id: 'test-room',
        type: 'safe',
        x: 6,
        y: 6,
        width: 7,
        height: 7,
        explored: true
      }
    ],
    upStairs: [],
    downStairs: [],
    mobs: {},
    items: {}
  },
  
  shop: {
    level: 1,
    width: 20,
    height: 20,
    tiles: Array(20).fill(null).map(() => 
      Array(20).fill(null).map(() => ({
        type: 'wall',
        walkable: false,
        explored: false
      }))
    ),
    rooms: [
      {
        id: 'test-room',
        type: 'shop',
        x: 6,
        y: 6,
        width: 7,
        height: 7,
        explored: true
      }
    ],
    upStairs: [],
    downStairs: [],
    mobs: {},
    items: {
      'shop-item-0': {
        id: 'shop-item-0',
        type: 'weapon',
        name: 'Shop Item 0',
        position: { x: 8, y: 8 }
      },
      'shop-item-1': {
        id: 'shop-item-1',
        type: 'armor',
        name: 'Shop Item 1',
        position: { x: 10, y: 8 }
      }
    }
  }
};

// Helper function to set up room tiles based on room data
export function setupRoomTiles(roomData: any) {
  const room = roomData.rooms[0];
  
  // Set floor tiles for the room
  for (let y = 0; y < room.height; y++) {
    for (let x = 0; x < room.width; x++) {
      roomData.tiles[room.y + y][room.x + x] = {
        type: 'floor',
        walkable: true,
        explored: true,
        roomId: room.id
      };
    }
  }
  
  // Add character in the center of the room
  const centerX = room.x + Math.floor(room.width / 2);
  const centerY = room.y + Math.floor(room.height / 2);
  roomData.tiles[centerY][centerX].character = 'test-character';
  
  // Add down stairs for entrance rooms
  if (room.type === 'entrance' && roomData.downStairs.length > 0) {
    const stairsPos = roomData.downStairs[0];
    roomData.tiles[stairsPos.y][stairsPos.x].type = 'downStairs';
  }
  
  // Add items
  Object.values(roomData.items).forEach((item: any) => {
    const { x, y } = item.position;
    if (x >= 0 && x < roomData.width && y >= 0 && y < roomData.height) {
      roomData.tiles[y][x].itemId = item.id;
    }
  });
  
  // Add mobs
  Object.values(roomData.mobs).forEach((mob: any) => {
    const { x, y } = mob.position;
    if (x >= 0 && x < roomData.width && y >= 0 && y < roomData.height) {
      roomData.tiles[y][x].mobId = mob.id;
    }
  });
  
  return roomData;
}

// Prepare all room data with tiles set up
export const preparedRoomData = {
  entrance: setupRoomTiles(JSON.parse(JSON.stringify(mockRoomData.entrance))),
  standard: setupRoomTiles(JSON.parse(JSON.stringify(mockRoomData.standard))),
  treasure: setupRoomTiles(JSON.parse(JSON.stringify(mockRoomData.treasure))),
  boss: setupRoomTiles(JSON.parse(JSON.stringify(mockRoomData.boss))),
  safe: setupRoomTiles(JSON.parse(JSON.stringify(mockRoomData.safe))),
  shop: setupRoomTiles(JSON.parse(JSON.stringify(mockRoomData.shop)))
}; 