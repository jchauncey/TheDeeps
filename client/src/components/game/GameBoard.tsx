import { useEffect, useState, useRef } from 'react';
import { Box, Spinner, Text, useToast, Button } from '@chakra-ui/react';
import { sendWebSocketMessage } from '../../services/api';
import { CLASS_COLORS } from '../../constants/gameConstants';
import { MapLegend } from './MapLegend';

// Define tile types and their colors
export const TILE_COLORS = {
  wall: '#333',
  floor: '#555',
  door: '#855',
  stairs_up: '#55f',
  stairs_down: '#f55',
};

// Define entity types and their colors
export const ENTITY_COLORS = {
  player: '#ff0',
  goblin: '#0f0',
  orc: '#0a0',
  skeleton: '#fff',
  rat: '#a50',
  bat: '#a0a',
  troll: '#0aa',
  ogre: '#a00',
  wraith: '#aaf',
  lich: '#a0f',
  ooze: '#0ff',
  ratman: '#a70',
  drake: '#f70',
  dragon: '#f00',
  elemental: '#7af',
};

// Define item types and their colors
export const ITEM_COLORS = {
  potion: '#f0f',
  scroll: '#ff0',
  weapon: '#aaa',
  armor: '#00f',
  gold: '#ff0',
};

// Define difficulty colors for mobs
export const DIFFICULTY_COLORS = {
  easy: '#aaa',    // Light gray border for easy mobs
  normal: '#fff',  // White border for normal mobs
  hard: '#ff0',    // Yellow border for hard mobs
  elite: '#f0f',   // Purple border for elite mobs
  boss: '#f00',    // Red border for boss mobs
};

// Define character class-specific colors and symbols
const CHARACTER_CLASS_STYLES = {
  warrior: { color: CLASS_COLORS.warrior.primary, symbol: '@', secondaryColor: CLASS_COLORS.warrior.secondary },
  mage: { color: CLASS_COLORS.mage.primary, symbol: '@', secondaryColor: CLASS_COLORS.mage.secondary },
  rogue: { color: CLASS_COLORS.rogue.primary, symbol: '@', secondaryColor: CLASS_COLORS.rogue.secondary },
  cleric: { color: CLASS_COLORS.cleric.primary, symbol: '@', secondaryColor: CLASS_COLORS.cleric.secondary },
  ranger: { color: CLASS_COLORS.ranger.primary, symbol: '@', secondaryColor: CLASS_COLORS.ranger.secondary },
  paladin: { color: CLASS_COLORS.paladin.primary, symbol: '@', secondaryColor: CLASS_COLORS.paladin.secondary },
  bard: { color: CLASS_COLORS.bard.primary, symbol: '@', secondaryColor: CLASS_COLORS.bard.secondary },
  monk: { color: CLASS_COLORS.monk.primary, symbol: '@', secondaryColor: CLASS_COLORS.monk.secondary },
  druid: { color: CLASS_COLORS.druid.primary, symbol: '@', secondaryColor: CLASS_COLORS.druid.secondary },
  barbarian: { color: CLASS_COLORS.barbarian.primary, symbol: '@', secondaryColor: CLASS_COLORS.barbarian.secondary },
  sorcerer: { color: CLASS_COLORS.sorcerer.primary, symbol: '@', secondaryColor: CLASS_COLORS.sorcerer.secondary },
  warlock: { color: CLASS_COLORS.warlock.primary, symbol: '@', secondaryColor: CLASS_COLORS.warlock.secondary },
  // Default for any unspecified class
  default: { color: CLASS_COLORS.default.primary, symbol: '@', secondaryColor: CLASS_COLORS.default.secondary },
};

interface Position {
  x: number;
  y: number;
}

interface Entity {
  id: string;
  type: string;
  name: string;
  position: Position;
  characterClass?: string; // Add character class for player entities
  health?: number; // Add health for status indicators
  maxHealth?: number;
  status?: string[]; // Add status effects array
  damage?: number; // Add damage for mob entities
  defense?: number; // Add defense for mob entities
  speed?: number; // Add speed for mob entities
  difficulty?: string; // Add difficulty for mob entities
}

interface Item {
  id: string;
  type: string;
  name: string;
  position: Position;
}

interface Tile {
  type: string;
  explored: boolean;
  visible: boolean;
  entity?: Entity;
  item?: Item;
}

interface Room {
  x: number;
  y: number;
  width: number;
  height: number;
}

interface Floor {
  level: number;
  width: number;
  height: number;
  tiles: Tile[][];
  rooms: Room[];
  entities: Entity[];
  items: Item[];
}

interface FloorData {
  type: string;
  floor: Floor;
  playerPosition: Position;
  currentFloor: number;
  playerData?: {
    characterClass?: string;
    health?: number;
    maxHealth?: number;
    status?: string[];
  };
}

interface GameBoardProps {
  floorData: FloorData | null;
}

export const GameBoard = ({ floorData }: GameBoardProps) => {
  const [loading, setLoading] = useState(true);
  const [floor, setFloor] = useState<Floor | null>(null);
  const [playerPos, setPlayerPos] = useState<Position | null>(null);
  const [currentFloor, setCurrentFloor] = useState(1);
  const [viewportSize, setViewportSize] = useState({ width: 800, height: 600 });
  const [error, setError] = useState<string | null>(null);
  const [hoveredEntity, setHoveredEntity] = useState<Entity | null>(null);
  const [tooltipPosition, setTooltipPosition] = useState<Position | null>(null);
  const [isLegendOpen, setIsLegendOpen] = useState(false);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const toast = useToast();

  // Process floor data when it changes
  useEffect(() => {
    if (floorData) {
      console.log('Processing floor data:', floorData);
      
      try {
        // Validate floor data
        if (!floorData.floor) {
          throw new Error('Floor data is missing floor property');
        }
        
        if (!floorData.playerPosition) {
          throw new Error('Floor data is missing playerPosition property');
        }
        
        if (!floorData.floor.tiles || !Array.isArray(floorData.floor.tiles)) {
          throw new Error('Floor data has invalid tiles property');
        }
        
        // Set floor data
        setFloor(floorData.floor);
        setPlayerPos(floorData.playerPosition);
        setCurrentFloor(floorData.currentFloor);
        setLoading(false);
        setError(null);
        
        console.log('Floor data processed successfully');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error processing floor data';
        console.error(errorMessage, floorData);
        setError(errorMessage);
        setLoading(false);
      }
    }
  }, [floorData]);

  // Request floor data when component mounts if not provided
  useEffect(() => {
    if (!floorData) {
      console.log('No floor data provided, requesting from server...');
      requestFloorData();
    }
  }, [floorData]);

  // Calculate viewport size on mount and when container size changes
  useEffect(() => {
    const updateViewportSize = () => {
      if (containerRef.current) {
        const rect = containerRef.current.getBoundingClientRect();
        console.log('Container size:', rect.width, rect.height);
        setViewportSize({
          width: Math.floor(rect.width - 30), // Account for padding
          height: Math.floor(rect.height - 60) // Account for header and padding
        });
      }
    };
    
    // Initial calculation after a short delay to ensure container is rendered
    const timeoutId = setTimeout(updateViewportSize, 100);
    
    // Create a ResizeObserver to watch for container size changes
    const resizeObserver = new ResizeObserver(() => {
      console.log('Container resized');
      updateViewportSize();
    });
    
    if (containerRef.current) {
      resizeObserver.observe(containerRef.current);
    }
    
    // Also listen for window resize events
    window.addEventListener('resize', updateViewportSize);
    
    return () => {
      clearTimeout(timeoutId);
      if (containerRef.current) {
        resizeObserver.unobserve(containerRef.current);
      }
      window.removeEventListener('resize', updateViewportSize);
    };
  }, []);

  // Draw the floor when data changes or viewport size changes
  useEffect(() => {
    if (floor && playerPos && viewportSize.width > 0 && viewportSize.height > 0) {
      console.log('Drawing floor with viewport size:', viewportSize);
      console.log('Player position for drawing:', playerPos);
      drawFloor();
    }
  }, [floor, playerPos, viewportSize]);

  // Handle WebSocket messages for additional updates
  useEffect(() => {
    const handleWebSocketMessage = (e: Event) => {
      const customEvent = e as CustomEvent;
      const data = customEvent.detail;
      
      if (data.type === 'floor_data') {
        try {
          const floorData = data as FloorData;
          setFloor(floorData.floor);
          setPlayerPos(floorData.playerPosition);
          setCurrentFloor(floorData.currentFloor);
          setLoading(false);
          setError(null);
        } catch (err) {
          const errorMessage = err instanceof Error ? err.message : 'Unknown error processing WebSocket message';
          console.error(errorMessage, data);
          setError(errorMessage);
        }
      }
    };

    window.addEventListener('websocket_message', handleWebSocketMessage);

    return () => {
      window.removeEventListener('websocket_message', handleWebSocketMessage);
    };
  }, []);

  // Ensure the game board has focus for keyboard controls
  useEffect(() => {
    const focusContainer = () => {
      if (containerRef.current) {
        containerRef.current.focus();
      }
    };

    // Focus when component mounts
    focusContainer();

    // Add click listener to focus when clicked
    if (containerRef.current) {
      containerRef.current.addEventListener('click', focusContainer);
    }

    return () => {
      if (containerRef.current) {
        containerRef.current.removeEventListener('click', focusContainer);
      }
    };
  }, []);

  // Request floor data from the server
  const requestFloorData = () => {
    console.log('Requesting floor data from server...');
    const success = sendWebSocketMessage({ type: 'get_floor' });
    console.log('Request sent:', success);
  };

  // Draw the floor on the canvas
  const drawFloor = () => {
    if (!floor || !canvasRef.current || !playerPos) {
      console.log('Missing data for drawing:', { 
        hasFloor: !!floor, 
        hasCanvas: !!canvasRef.current, 
        hasPlayerPos: !!playerPos 
      });
      return;
    }

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) {
      console.log('Could not get canvas context');
      return;
    }

    try {
      // Calculate the visible area (viewport)
      const visibleTiles = {
        width: Math.min(floor.width, 40), // Limit to 40 tiles wide
        height: Math.min(floor.height, 25) // Limit to 25 tiles high
      };
      
      // Calculate tile size based on available space and visible area
      const tileSize = Math.max(
        1,
        Math.min(
          Math.floor(viewportSize.width / visibleTiles.width),
          Math.floor(viewportSize.height / visibleTiles.height)
        )
      );
      
      console.log('Tile size calculated:', tileSize);
      
      // Set canvas size
      canvas.width = visibleTiles.width * tileSize;
      canvas.height = visibleTiles.height * tileSize;
      
      console.log('Canvas size set to:', canvas.width, canvas.height);

      // Calculate viewport center (player position)
      const viewportCenterX = Math.floor(visibleTiles.width / 2);
      const viewportCenterY = Math.floor(visibleTiles.height / 2);
      
      // Calculate top-left corner of viewport in dungeon coordinates
      const viewportStartX = Math.max(0, playerPos.x - viewportCenterX);
      const viewportStartY = Math.max(0, playerPos.y - viewportCenterY);
      
      // Adjust if we're near the edge of the map
      const maxStartX = Math.max(0, floor.width - visibleTiles.width);
      const maxStartY = Math.max(0, floor.height - visibleTiles.height);
      
      const adjustedStartX = Math.min(viewportStartX, maxStartX);
      const adjustedStartY = Math.min(viewportStartY, maxStartY);
      
      console.log('Viewport start:', adjustedStartX, adjustedStartY);

      // Clear canvas
      ctx.fillStyle = '#000';
      ctx.fillRect(0, 0, canvas.width, canvas.height);

      // Draw tiles within viewport
      for (let y = 0; y < visibleTiles.height; y++) {
        for (let x = 0; x < visibleTiles.width; x++) {
          // Convert viewport coordinates to dungeon coordinates
          const dungeonX = adjustedStartX + x;
          const dungeonY = adjustedStartY + y;
          
          // Skip if out of bounds
          if (dungeonX >= floor.width || dungeonY >= floor.height) continue;
          
          // Check if tiles array is properly structured
          if (!floor.tiles[dungeonY] || !floor.tiles[dungeonY][dungeonX]) {
            console.error('Invalid tile data at', dungeonX, dungeonY);
            continue;
          }
          
          const tile = floor.tiles[dungeonY][dungeonX];
          
          // Draw all tiles
          const baseColor = TILE_COLORS[tile.type as keyof typeof TILE_COLORS] || '#000';
          
          // Parse the hex color to RGB
          let r = 0, g = 0, b = 0;
          try {
            // Ensure baseColor is a valid hex color
            if (baseColor && baseColor.startsWith('#') && baseColor.length >= 7) {
              r = parseInt(baseColor.slice(1, 3), 16);
              g = parseInt(baseColor.slice(3, 5), 16);
              b = parseInt(baseColor.slice(5, 7), 16);
            }
            
            // Validate the parsed values
            r = isNaN(r) ? 0 : r;
            g = isNaN(g) ? 0 : g;
            b = isNaN(b) ? 0 : b;
          } catch (error) {
            console.error('Error parsing color:', baseColor, error);
            // Default to gray if parsing fails
            r = g = b = 128;
          }
          
          // Set color for the tile
          ctx.fillStyle = baseColor;
          ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
          
          // Draw grid lines
          ctx.strokeStyle = 'rgba(0, 0, 0, 0.2)';
          ctx.lineWidth = 0.5;
          ctx.strokeRect(x * tileSize, y * tileSize, tileSize, tileSize);
        }
      }

      // Draw items within viewport
      if (floor.items && Array.isArray(floor.items)) {
        floor.items.forEach(item => {
          // Convert dungeon coordinates to viewport coordinates
          const viewportX = item.position.x - adjustedStartX;
          const viewportY = item.position.y - adjustedStartY;
          
          // Skip if outside viewport
          if (viewportX < 0 || viewportX >= visibleTiles.width || 
              viewportY < 0 || viewportY >= visibleTiles.height) return;
          
          // Check if position is valid
          if (!floor.tiles[item.position.y] || !floor.tiles[item.position.y][item.position.x]) {
            console.error('Invalid item position:', item.position);
            return;
          }
          
          // Tile reference is kept for future visibility checks
          // @ts-ignore
          const tile = floor.tiles[item.position.y][item.position.x];
          
          // Draw with different opacity based on visibility
          const baseColor = ITEM_COLORS[item.type as keyof typeof ITEM_COLORS] || '#fff';
          
          ctx.fillStyle = baseColor;
          
          ctx.beginPath();
          ctx.arc(
            viewportX * tileSize + tileSize / 2,
            viewportY * tileSize + tileSize / 2,
            tileSize / 4,
            0,
            Math.PI * 2
          );
          ctx.fill();
        });
      }

      // Draw entities within viewport
      if (floor.entities && Array.isArray(floor.entities)) {
        floor.entities.forEach(entity => {
          // Convert dungeon coordinates to viewport coordinates
          const viewportX = entity.position.x - adjustedStartX;
          const viewportY = entity.position.y - adjustedStartY;
          
          // Skip if outside viewport
          if (viewportX < 0 || viewportX >= visibleTiles.width || 
              viewportY < 0 || viewportY >= visibleTiles.height) return;
          
          // Check if position is valid
          if (!floor.tiles[entity.position.y] || !floor.tiles[entity.position.y][entity.position.x]) {
            console.error('Invalid entity position:', entity.position);
            return;
          }
          
          // Draw entity with improved visuals
          const baseColor = ENTITY_COLORS[entity.type as keyof typeof ENTITY_COLORS] || '#f00';
          ctx.fillStyle = baseColor;
          
          // Draw entity as a circle with a border
          ctx.beginPath();
          ctx.arc(
            viewportX * tileSize + tileSize / 2,
            viewportY * tileSize + tileSize / 2,
            tileSize / 3,
            0,
            Math.PI * 2
          );
          ctx.fill();
          
          // Determine difficulty from entity name
          let difficulty = 'normal';
          if (entity.name.startsWith('easy')) {
            difficulty = 'easy';
          } else if (entity.name.startsWith('hard')) {
            difficulty = 'hard';
          } else if (entity.name.startsWith('elite')) {
            difficulty = 'elite';
          } else if (entity.name.startsWith('boss')) {
            difficulty = 'boss';
          }
          
          // Add a border with difficulty color
          ctx.strokeStyle = DIFFICULTY_COLORS[difficulty as keyof typeof DIFFICULTY_COLORS] || '#000';
          ctx.lineWidth = difficulty === 'boss' ? 3 : difficulty === 'elite' ? 2 : 1;
          ctx.stroke();
          
          // Add a letter indicator for entity type
          ctx.fillStyle = '#000';
          ctx.font = `${Math.max(8, tileSize / 2)}px monospace`;
          ctx.textAlign = 'center';
          ctx.textBaseline = 'middle';
          ctx.fillText(
            entity.type.charAt(0).toUpperCase(),
            viewportX * tileSize + tileSize / 2,
            viewportY * tileSize + tileSize / 2
          );
          
          // Draw health bar if health is available
          if (entity.health !== undefined && entity.maxHealth !== undefined && entity.health < entity.maxHealth) {
            const healthPercentage = entity.health / entity.maxHealth;
            
            // Health bar background
            ctx.fillStyle = '#500';
            ctx.fillRect(
              viewportX * tileSize,
              viewportY * tileSize - tileSize / 5,
              tileSize,
              tileSize / 10
            );
            
            // Health bar fill
            ctx.fillStyle = healthPercentage > 0.5 ? '#0f0' : healthPercentage > 0.25 ? '#ff0' : '#f00';
            ctx.fillRect(
              viewportX * tileSize,
              viewportY * tileSize - tileSize / 5,
              tileSize * healthPercentage,
              tileSize / 10
            );
          }
        });
      }

      // Draw player with class-specific styling
      // Convert dungeon coordinates to viewport coordinates
      const playerViewportX = playerPos.x - adjustedStartX;
      const playerViewportY = playerPos.y - adjustedStartY;
      
      console.log(`Drawing player at viewport coordinates: (${playerViewportX}, ${playerViewportY})`);
      
      // Only draw player if within viewport
      if (playerViewportX >= 0 && playerViewportX < visibleTiles.width &&
          playerViewportY >= 0 && playerViewportY < visibleTiles.height) {
        
        // Find player entity to get class and health info
        const playerEntity = floor.entities?.find(e => e.type === 'player');
        const characterClass = playerEntity?.characterClass || floorData?.playerData?.characterClass || 'default';
        const classStyle = CHARACTER_CLASS_STYLES[characterClass as keyof typeof CHARACTER_CLASS_STYLES] || 
                           CHARACTER_CLASS_STYLES.default;
        
        // Calculate health percentage if available
        const healthPercentage = playerEntity?.health && playerEntity?.maxHealth 
          ? playerEntity.health / playerEntity.maxHealth 
          : floorData?.playerData?.health && floorData?.playerData?.maxHealth
            ? floorData.playerData.health / floorData.playerData.maxHealth
            : 1;
        
        // Draw a more visible player marker with class-specific styling
        ctx.fillStyle = classStyle.color;
        
        // Draw player as a circle
        ctx.beginPath();
        ctx.arc(
          playerViewportX * tileSize + tileSize / 2,
          playerViewportY * tileSize + tileSize / 2,
          tileSize / 2.5,
          0,
          Math.PI * 2
        );
        ctx.fill();
        
        // Add a border
        ctx.strokeStyle = classStyle.secondaryColor;
        ctx.lineWidth = 2;
        ctx.stroke();
        
        // Add the class symbol
        ctx.fillStyle = '#000';
        ctx.font = `bold ${Math.max(10, tileSize / 1.5)}px monospace`;
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(
          classStyle.symbol,
          playerViewportX * tileSize + tileSize / 2,
          playerViewportY * tileSize + tileSize / 2
        );
        
        // Draw health indicator if health is less than 100%
        if (healthPercentage < 1) {
          // Health bar background
          ctx.fillStyle = '#500';
          ctx.fillRect(
            playerViewportX * tileSize,
            playerViewportY * tileSize - tileSize / 5,
            tileSize,
            tileSize / 10
          );
          
          // Health bar fill
          ctx.fillStyle = healthPercentage > 0.5 ? '#0f0' : healthPercentage > 0.25 ? '#ff0' : '#f00';
          ctx.fillRect(
            playerViewportX * tileSize,
            playerViewportY * tileSize - tileSize / 5,
            tileSize * healthPercentage,
            tileSize / 10
          );
        }
        
        // Draw status effects if any
        if (playerEntity?.status && playerEntity.status.length > 0) {
          const statusColors = {
            poisoned: '#0f0',
            burning: '#f50',
            frozen: '#0ff',
            blessed: '#ff0',
            cursed: '#f0f',
            default: '#fff'
          };
          
          // Draw status indicator
          playerEntity.status.slice(0, 3).forEach((status, index) => {
            const statusColor = statusColors[status as keyof typeof statusColors] || statusColors.default;
            ctx.fillStyle = statusColor;
            ctx.beginPath();
            ctx.arc(
              playerViewportX * tileSize + tileSize / 4 + (index * tileSize / 4),
              playerViewportY * tileSize + tileSize - tileSize / 6,
              tileSize / 10,
              0,
              Math.PI * 2
            );
            ctx.fill();
          });
        }
      }
      
      console.log('Drawing complete');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error drawing floor';
      console.error(errorMessage);
      setError(errorMessage);
    }
  };

  // Handle mouse move to show entity tooltips
  const handleMouseMove = (e: React.MouseEvent<HTMLCanvasElement>) => {
    if (!floor || !canvasRef.current) return;

    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const mouseX = e.clientX - rect.left;
    const mouseY = e.clientY - rect.top;

    // Calculate the visible area (viewport)
    const visibleTiles = {
      width: Math.min(floor.width, 40), // Limit to 40 tiles wide
      height: Math.min(floor.height, 25) // Limit to 25 tiles high
    };
    
    // Calculate tile size
    const tileSize = Math.max(
      1,
      Math.min(
        Math.floor(viewportSize.width / visibleTiles.width),
        Math.floor(viewportSize.height / visibleTiles.height)
      )
    );

    // Calculate viewport center (player position)
    const viewportCenterX = Math.floor(visibleTiles.width / 2);
    const viewportCenterY = Math.floor(visibleTiles.height / 2);
    
    // Calculate top-left corner of viewport in dungeon coordinates
    const viewportStartX = Math.max(0, playerPos!.x - viewportCenterX);
    const viewportStartY = Math.max(0, playerPos!.y - viewportCenterY);
    
    // Adjust if we're near the edge of the map
    const maxStartX = Math.max(0, floor.width - visibleTiles.width);
    const maxStartY = Math.max(0, floor.height - visibleTiles.height);
    
    const adjustedStartX = Math.min(viewportStartX, maxStartX);
    const adjustedStartY = Math.min(viewportStartY, maxStartY);

    // Convert mouse position to tile coordinates
    const tileX = Math.floor(mouseX / tileSize);
    const tileY = Math.floor(mouseY / tileSize);

    // Convert viewport coordinates to dungeon coordinates
    const dungeonX = adjustedStartX + tileX;
    const dungeonY = adjustedStartY + tileY;

    // Find entity at this position
    const entity = floor.entities.find(e => 
      e.position.x === dungeonX && e.position.y === dungeonY
    );

    if (entity) {
      setHoveredEntity(entity);
      setTooltipPosition({ x: e.clientX, y: e.clientY });
    } else {
      setHoveredEntity(null);
      setTooltipPosition(null);
    }
  };

  // Handle mouse leave to hide tooltips
  const handleMouseLeave = () => {
    setHoveredEntity(null);
    setTooltipPosition(null);
  };

  // Focus the container to capture keyboard events
  const focusContainer = () => {
    if (containerRef.current) {
      containerRef.current.focus();
    }
  };

  // Handle key down events
  const handleKeyDown = (e: React.KeyboardEvent<HTMLDivElement>) => {
    if (!floor) return;

    // Prevent default behavior for arrow keys to avoid scrolling
    if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', ' '].includes(e.key)) {
      e.preventDefault();
    }

    // Toggle legend with 'L' key
    if (e.key.toLowerCase() === 'l') {
      setIsLegendOpen(!isLegendOpen);
      return;
    }

    // Handle movement with arrow keys
    switch (e.key) {
      case 'ArrowUp':
        sendWebSocketMessage({ type: 'move', direction: 'up' });
        break;
      case 'ArrowDown':
        sendWebSocketMessage({ type: 'move', direction: 'down' });
        break;
      case 'ArrowLeft':
        sendWebSocketMessage({ type: 'move', direction: 'left' });
        break;
      case 'ArrowRight':
        sendWebSocketMessage({ type: 'move', direction: 'right' });
        break;
      case ' ': // Space bar for attack
        sendWebSocketMessage({ type: 'action', action: 'attack' });
        break;
      case 'p': // 'p' for pickup
      case 'P':
        sendWebSocketMessage({ type: 'action', action: 'pickup' });
        break;
      default:
        break;
    }
  };

  // Set up the game board when floor data changes
  useEffect(() => {
    if (floorData) {
      setLoading(false);
      setFloor(floorData.floor);
      setPlayerPos(floorData.playerPosition);
      setCurrentFloor(floorData.currentFloor);
      
      // Focus the container to capture keyboard events
      if (containerRef.current) {
        containerRef.current.focus();
      }
    }
  }, [floorData]);

  if (loading) {
    return (
      <Box
        width="100%"
        height="100%"
        display="flex"
        alignItems="center"
        justifyContent="center"
        bg="#291326"
      >
        <Spinner size="xl" color="purple.500" thickness="4px" />
        <Text ml={4} color="white" fontSize="xl">
          Loading dungeon...
        </Text>
      </Box>
    );
  }

  if (error) {
    return (
      <Box
        width="100%"
        height="100%"
        display="flex"
        flexDirection="column"
        alignItems="center"
        justifyContent="center"
        bg="#291326"
        p={4}
      >
        <Text color="red.500" fontSize="xl" mb={4}>
          Error loading dungeon
        </Text>
        <Text color="white" fontSize="md">
          {error}
        </Text>
        <Button 
          mt={4} 
          colorScheme="purple" 
          onClick={() => {
            setLoading(true);
            setError(null);
            requestFloorData();
          }}
        >
          Retry
        </Button>
      </Box>
    );
  }

  return (
    <Box 
      ref={containerRef}
      position="relative" 
      width="100%" 
      height="100%" 
      bg="gray.900"
      tabIndex={0}
      onKeyDown={handleKeyDown}
      onClick={focusContainer}
      overflow="hidden"
    >
      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="100%">
          <Spinner size="xl" color="blue.500" />
          <Text ml={4} color="white">Loading dungeon...</Text>
        </Box>
      ) : (
        <>
          <canvas 
            ref={canvasRef} 
            style={{ 
              display: 'block',
              margin: '0 auto',
              imageRendering: 'pixelated'
            }}
            onMouseMove={handleMouseMove}
            onMouseLeave={handleMouseLeave}
          />
          
          {/* Entity tooltip */}
          {hoveredEntity && tooltipPosition && (
            <Box
              position="fixed"
              left={`${tooltipPosition.x + 10}px`}
              top={`${tooltipPosition.y + 10}px`}
              bg="gray.800"
              color="white"
              p={2}
              borderRadius="md"
              boxShadow="md"
              zIndex={1000}
              maxWidth="250px"
            >
              <Text fontWeight="bold">{hoveredEntity.name}</Text>
              {hoveredEntity.health !== undefined && hoveredEntity.maxHealth !== undefined && (
                <Text>Health: {hoveredEntity.health}/{hoveredEntity.maxHealth}</Text>
              )}
              {hoveredEntity.damage !== undefined && (
                <Text>Damage: {hoveredEntity.damage}</Text>
              )}
              {hoveredEntity.defense !== undefined && (
                <Text>Defense: {hoveredEntity.defense}</Text>
              )}
              {hoveredEntity.speed !== undefined && (
                <Text>Speed: {hoveredEntity.speed}</Text>
              )}
              {hoveredEntity.status && hoveredEntity.status.length > 0 && (
                <Text>Status: {hoveredEntity.status.join(', ')}</Text>
              )}
            </Box>
          )}

          {/* Legend button */}
          <Box position="absolute" top="10px" right="10px">
            <Button 
              colorScheme="blue" 
              size="sm" 
              onClick={() => setIsLegendOpen(true)}
            >
              Legend (L)
            </Button>
          </Box>

          {/* Map Legend Modal */}
          <MapLegend isOpen={isLegendOpen} onClose={() => setIsLegendOpen(false)} />
        </>
      )}
    </Box>
  );
}; 