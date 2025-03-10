import { useEffect, useState, useRef } from 'react';
import { Box, Spinner, Text, useToast, Button, Flex } from '@chakra-ui/react';
import { sendWebSocketMessage } from '../../services/api';
import { CLASS_COLORS } from '../../constants/gameConstants';
import { MapLegend } from './MapLegend';
import { FloorData, Position, Entity, Tile, Room, Floor, DungeonItem } from '../../types/game';

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

// Update the Entity interface to include description
interface EnhancedEntity extends Entity {
  description?: string;
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
  const [hoveredEntity, setHoveredEntity] = useState<EnhancedEntity | null>(null);
  const [tooltipPosition, setTooltipPosition] = useState<Position | null>(null);
  const [isLegendOpen, setIsLegendOpen] = useState(false);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const toast = useToast();

  // Process floor data when it changes
  useEffect(() => {
    if (floorData) {
      console.log('Processing floor data:', JSON.stringify(floorData, null, 2));
      
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
        
        console.log('Floor data processed successfully', {
          floorLevel: floorData.currentFloor,
          playerPosition: floorData.playerPosition,
          tilesCount: floorData.floor.tiles.length,
          entitiesCount: floorData.floor.entities?.length || 0,
          itemsCount: floorData.floor.items?.length || 0
        });
        
        // Force a redraw after a short delay to ensure the canvas is ready
        setTimeout(() => {
          drawFloor();
        }, 100);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error processing floor data';
        console.error('Error processing floor data:', errorMessage, floorData);
        setError(errorMessage);
        setLoading(false);
        
        // Show error toast
        toast({
          title: 'Error',
          description: `Failed to load dungeon: ${errorMessage}`,
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    } else {
      console.log('No floor data provided to GameBoard component');
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
        
        // Set the canvas dimensions to match the container
        if (canvasRef.current) {
          canvasRef.current.width = rect.width;
          canvasRef.current.height = rect.height;
        }
        
        setViewportSize({
          width: Math.floor(rect.width),
          height: Math.floor(rect.height)
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

  // Effect to update the game board when floorData changes
  useEffect(() => {
    if (floorData) {
      setFloor(floorData.floor);
      setPlayerPos(floorData.playerPosition);
      setCurrentFloor(floorData.currentFloor);
      setLoading(false);
      setError(null);
      
      // Dispatch a custom event to notify other components
      const customEvent = new CustomEvent('websocket_message', {
        detail: floorData
      });
      window.dispatchEvent(customEvent);
    }
  }, [floorData]);

  // Effect to listen for WebSocket messages
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
    console.log('drawFloor called with:', {
      floorExists: !!floor,
      playerPosExists: !!playerPos,
      canvasExists: !!canvasRef.current
    });
    
    if (!floor || !playerPos || !canvasRef.current) {
      console.error('Cannot draw floor - missing required data:', {
        floor: !!floor,
        playerPos: !!playerPos,
        canvas: !!canvasRef.current
      });
      return;
    }
    
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) {
      console.error('Failed to get canvas context');
      return;
    }
    
    console.log('Drawing floor with viewport size:', viewportSize);
    console.log('Player position for drawing:', playerPos);
    
    // Clear the canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    // Calculate the ideal number of tiles to show
    const idealVisibleTilesX = 30; // Show more tiles horizontally
    const idealVisibleTilesY = 20; // Show more tiles vertically
    
    // Calculate tile size based on viewport and ideal number of tiles
    const tileSize = Math.min(
      canvas.width / idealVisibleTilesX,
      canvas.height / idealVisibleTilesY
    );
    
    console.log('Tile size calculated:', tileSize);
    
    // Calculate viewport boundaries to only render what's visible
    const visibleTilesX = Math.ceil(canvas.width / tileSize);
    const visibleTilesY = Math.ceil(canvas.height / tileSize);
    
    // Center the viewport on the player with improved edge handling
    let startX = Math.max(0, playerPos.x - Math.floor(visibleTilesX / 2));
    let startY = Math.max(0, playerPos.y - Math.floor(visibleTilesY / 2));
    
    // Ensure we always show the same number of tiles when possible
    // This prevents the map from shrinking when reaching the edges
    if (startX + visibleTilesX > floor.width) {
      startX = Math.max(0, floor.width - visibleTilesX);
    }
    
    if (startY + visibleTilesY > floor.height) {
      startY = Math.max(0, floor.height - visibleTilesY);
    }
    
    const endX = Math.min(floor.width, startX + visibleTilesX);
    const endY = Math.min(floor.height, startY + visibleTilesY);
    
    console.log('Viewport calculation:', {
      startX, startY, endX, endY,
      visibleTilesX, visibleTilesY,
      floorWidth: floor.width,
      floorHeight: floor.height
    });
    
    // Calculate offset to center the map in the viewport
    // Use fixed tile counts to maintain consistent sizing
    const tilesShownX = endX - startX;
    const tilesShownY = endY - startY;
    const offsetX = (canvas.width - tilesShownX * tileSize) / 2;
    const offsetY = (canvas.height - tilesShownY * tileSize) / 2;
    
    // Draw tiles
    try {
      for (let y = startY; y < endY; y++) {
        for (let x = startX; x < endX; x++) {
          if (!floor.tiles[y] || !floor.tiles[y][x]) {
            console.error(`Missing tile data at position [${y}][${x}]`);
            continue;
          }
          
          const tile = floor.tiles[y][x];
          const screenX = (x - startX) * tileSize + offsetX;
          const screenY = (y - startY) * tileSize + offsetY;
          
          // Draw tile background
          ctx.fillStyle = TILE_COLORS[tile.type as keyof typeof TILE_COLORS] || '#000';
          ctx.fillRect(screenX, screenY, tileSize, tileSize);
          
          // Draw tile border
          ctx.strokeStyle = '#222';
          ctx.lineWidth = 1;
          ctx.strokeRect(screenX, screenY, tileSize, tileSize);
          
          // Draw item if present
          if (tile.item) {
            drawItem(ctx, tile.item, screenX, screenY, tileSize);
          }
        }
      }
      
      // Draw entities
      if (floor.entities && floor.entities.length > 0) {
        console.log(`Drawing ${floor.entities.length} entities`);
        for (const entity of floor.entities) {
          // Only draw entities within the viewport
          if (
            entity.position.x >= startX && 
            entity.position.x < endX && 
            entity.position.y >= startY && 
            entity.position.y < endY
          ) {
            const screenX = (entity.position.x - startX) * tileSize + offsetX;
            const screenY = (entity.position.y - startY) * tileSize + offsetY;
            drawEntity(ctx, entity, screenX, screenY, tileSize);
          }
        }
      } else {
        console.log('No entities to draw');
      }
      
      console.log('Floor drawing completed successfully');
    } catch (error) {
      console.error('Error drawing floor:', error);
      setError('Error drawing floor: ' + (error instanceof Error ? error.message : String(error)));
      
      // Show error toast
      toast({
        title: 'Rendering Error',
        description: 'Failed to render dungeon. Please try refreshing the page.',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  // Update the drawItem function to use DungeonItem
  const drawItem = (
    ctx: CanvasRenderingContext2D, 
    item: DungeonItem, 
    x: number, 
    y: number, 
    size: number
  ) => {
    // Draw item
    const itemColor = ITEM_COLORS[item.type as keyof typeof ITEM_COLORS] || '#fff';
    ctx.fillStyle = itemColor;
    
    // Draw a smaller rectangle for the item
    const padding = size * 0.25;
    ctx.fillRect(x + padding, y + padding, size - padding * 2, size - padding * 2);
    
    // Draw item border
    ctx.strokeStyle = '#fff';
    ctx.lineWidth = 1;
    ctx.strokeRect(x + padding, y + padding, size - padding * 2, size - padding * 2);
  };

  // Add the drawEntity function
  const drawEntity = (
    ctx: CanvasRenderingContext2D, 
    entity: Entity, 
    x: number, 
    y: number, 
    size: number
  ) => {
    // Check if this is a player entity
    if (entity.type === 'player') {
      // Get class-specific styling
      const characterClass = entity.characterClass || 'default';
      const classStyle = CHARACTER_CLASS_STYLES[characterClass as keyof typeof CHARACTER_CLASS_STYLES] || CHARACTER_CLASS_STYLES.default;
      
      // Draw player with class-specific color
      ctx.fillStyle = classStyle.color;
      
      // Draw background for better visibility
      ctx.beginPath();
      ctx.arc(
        x + size / 2,
        y + size / 2,
        size / 3,
        0,
        Math.PI * 2
      );
      ctx.fill();
      
      // Draw the @ symbol
      ctx.fillStyle = '#000'; // Black text for contrast
      ctx.font = `bold ${Math.max(12, size / 1.5)}px monospace`;
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText(
        '@',
        x + size / 2,
        y + size / 2
      );
      
      // Add a border with class color
      ctx.strokeStyle = classStyle.secondaryColor;
      ctx.lineWidth = 2;
      ctx.stroke();
    } else {
      // Draw non-player entity
      const entityColor = ENTITY_COLORS[entity.type as keyof typeof ENTITY_COLORS] || '#f00';
      ctx.fillStyle = entityColor;
      
      // Draw entity as a circle
      ctx.beginPath();
      ctx.arc(
        x + size / 2,
        y + size / 2,
        size / 3,
        0,
        Math.PI * 2
      );
      ctx.fill();
      
      // Determine difficulty from entity name or use the difficulty property
      let difficulty = entity.difficulty || 'normal';
      if (!entity.difficulty) {
        if (entity.name.startsWith('easy')) {
          difficulty = 'easy';
        } else if (entity.name.startsWith('hard')) {
          difficulty = 'hard';
        } else if (entity.name.startsWith('elite')) {
          difficulty = 'elite';
        } else if (entity.name.startsWith('boss')) {
          difficulty = 'boss';
        }
      }
      
      // Add a border with difficulty color
      ctx.strokeStyle = DIFFICULTY_COLORS[difficulty as keyof typeof DIFFICULTY_COLORS] || '#000';
      ctx.lineWidth = difficulty === 'boss' ? 3 : difficulty === 'elite' ? 2 : 1;
      ctx.stroke();
      
      // Add a letter indicator for entity type
      ctx.fillStyle = '#000';
      ctx.font = `${Math.max(8, size / 2)}px monospace`;
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText(
        entity.type.charAt(0).toUpperCase(),
        x + size / 2,
        y + size / 2
      );
    }
    
    // Draw health bar if health is available
    if (entity.health !== undefined && entity.maxHealth !== undefined && entity.health < entity.maxHealth) {
      const healthPercentage = entity.health / entity.maxHealth;
      
      // Health bar background
      ctx.fillStyle = '#500';
      ctx.fillRect(
        x,
        y - size / 5,
        size,
        size / 10
      );
      
      // Health bar fill
      ctx.fillStyle = healthPercentage > 0.5 ? '#0f0' : healthPercentage > 0.25 ? '#ff0' : '#f00';
      ctx.fillRect(
        x,
        y - size / 5,
        size * healthPercentage,
        size / 10
      );
    }
  };

  // Update the handleMouseMove function to account for the new offset
  const handleMouseMove = (e: React.MouseEvent<HTMLCanvasElement>) => {
    if (!floor || !playerPos || !canvasRef.current) return;
    
    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const mouseX = e.clientX - rect.left;
    const mouseY = e.clientY - rect.top;
    
    // Calculate the ideal number of tiles to show
    const idealVisibleTilesX = 30;
    const idealVisibleTilesY = 20;
    
    // Calculate tile size
    const tileSize = Math.min(
      canvas.width / idealVisibleTilesX,
      canvas.height / idealVisibleTilesY
    );
    
    // Calculate viewport boundaries
    const visibleTilesX = Math.ceil(canvas.width / tileSize);
    const visibleTilesY = Math.ceil(canvas.height / tileSize);
    
    // Center the viewport on the player
    const startX = Math.max(0, playerPos.x - Math.floor(visibleTilesX / 2));
    const startY = Math.max(0, playerPos.y - Math.floor(visibleTilesY / 2));
    const endX = Math.min(floor.width, startX + visibleTilesX);
    const endY = Math.min(floor.height, startY + visibleTilesY);
    
    // Calculate offset to center the map in the viewport
    const offsetX = (canvas.width - (endX - startX) * tileSize) / 2;
    const offsetY = (canvas.height - (endY - startY) * tileSize) / 2;
    
    // Convert mouse position to tile coordinates
    const tileX = Math.floor((mouseX - offsetX) / tileSize) + startX;
    const tileY = Math.floor((mouseY - offsetY) / tileSize) + startY;
    
    // Check if mouse is over a valid tile
    if (
      tileX >= startX && 
      tileX < endX && 
      tileY >= startY && 
      tileY < endY
    ) {
      // Check if there's an entity at this position
      const entityAtPosition = floor.entities.find(
        entity => entity.position.x === tileX && entity.position.y === tileY
      );
      
      if (entityAtPosition) {
        setHoveredEntity(entityAtPosition as EnhancedEntity);
        setTooltipPosition({ x: mouseX, y: mouseY });
      } else {
        setHoveredEntity(null);
      }
    } else {
      setHoveredEntity(null);
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

  // Update the MapLegend component to work without props
  const MapLegendComponent = () => (
    <Box>
      <Text fontWeight="bold" fontSize="lg" mb={2}>
        Map Legend
      </Text>
      
      {/* Tiles Section */}
      <Text fontWeight="bold" fontSize="md" mb={2}>
        Tiles
      </Text>
      <Flex flexWrap="wrap" mb={4}>
        {Object.entries(TILE_COLORS).map(([type, color]) => (
          <Flex key={type} alignItems="center" mr={4} mb={2}>
            <Box
              w="20px"
              h="20px"
              bg={color as string}
              mr={2}
              border="1px solid"
              borderColor="gray.500"
            />
            <Text fontSize="sm">{type.replace('_', ' ')}</Text>
          </Flex>
        ))}
      </Flex>
      
      {/* Player Characters Section */}
      <Text fontWeight="bold" fontSize="md" mb={2}>
        Player Characters
      </Text>
      <Flex flexWrap="wrap" mb={4}>
        {Object.entries(CHARACTER_CLASS_STYLES).slice(0, 6).map(([className, style]) => (
          <Flex key={className} alignItems="center" mr={4} mb={2}>
            <Box
              w="20px"
              h="20px"
              bg={style.color}
              mr={2}
              borderRadius="full"
              display="flex"
              alignItems="center"
              justifyContent="center"
              border="2px solid"
              borderColor={style.secondaryColor}
            >
              <Text fontSize="xs" fontWeight="bold" color="black">
                @
              </Text>
            </Box>
            <Text fontSize="sm">{className}</Text>
          </Flex>
        ))}
      </Flex>
      
      {/* Entities Section */}
      <Text fontWeight="bold" fontSize="md" mb={2}>
        Monsters
      </Text>
      <Flex flexWrap="wrap" mb={4}>
        {Object.entries(ENTITY_COLORS).map(([type, color]) => (
          <Flex key={type} alignItems="center" mr={4} mb={2}>
            <Box
              w="20px"
              h="20px"
              bg={color as string}
              mr={2}
              borderRadius="full"
              display="flex"
              alignItems="center"
              justifyContent="center"
              border="1px solid"
              borderColor="gray.500"
            >
              <Text fontSize="xs" fontWeight="bold" color="black">
                {type.charAt(0).toUpperCase()}
              </Text>
            </Box>
            <Text fontSize="sm">{type}</Text>
          </Flex>
        ))}
      </Flex>
      
      {/* Items Section */}
      <Text fontWeight="bold" fontSize="md" mb={2}>
        Items
      </Text>
      <Flex flexWrap="wrap" mb={4}>
        {Object.entries(ITEM_COLORS).map(([type, color]) => (
          <Flex key={type} alignItems="center" mr={4} mb={2}>
            <Box
              w="20px"
              h="20px"
              bg={color as string}
              mr={2}
              border="1px solid"
              borderColor="gray.500"
            />
            <Text fontSize="sm">{type}</Text>
          </Flex>
        ))}
      </Flex>
      
      {/* Controls Section */}
      <Text fontWeight="bold" fontSize="md" mb={2}>
        Controls
      </Text>
      <Flex flexDirection="column">
        <Flex mb={1}>
          <Text fontWeight="bold" minWidth="80px" fontSize="sm">WASD:</Text>
          <Text fontSize="sm">Move</Text>
        </Flex>
        <Flex mb={1}>
          <Text fontWeight="bold" minWidth="80px" fontSize="sm">Space:</Text>
          <Text fontSize="sm">Attack/Interact</Text>
        </Flex>
        <Flex mb={1}>
          <Text fontWeight="bold" minWidth="80px" fontSize="sm">P:</Text>
          <Text fontSize="sm">Pick up item</Text>
        </Flex>
        <Flex mb={1}>
          <Text fontWeight="bold" minWidth="80px" fontSize="sm">U:</Text>
          <Text fontSize="sm">Ascend stairs</Text>
        </Flex>
        <Flex mb={1}>
          <Text fontWeight="bold" minWidth="80px" fontSize="sm">D:</Text>
          <Text fontSize="sm">Descend stairs</Text>
        </Flex>
      </Flex>
    </Box>
  );

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
      tabIndex={0}
      onKeyDown={handleKeyDown}
      outline="none"
      overflow="hidden"
      backgroundColor="#111"
    >
      {loading ? (
        <Flex 
          position="absolute" 
          top="0" 
          left="0" 
          right="0" 
          bottom="0" 
          alignItems="center" 
          justifyContent="center"
          flexDirection="column"
          bg="rgba(0,0,0,0.8)"
          zIndex="10"
        >
          <Spinner size="xl" color="blue.500" thickness="4px" speed="0.65s" mb={4} />
          <Text color="white" fontSize="xl">Loading Dungeon...</Text>
        </Flex>
      ) : error ? (
        <Flex 
          position="absolute" 
          top="0" 
          left="0" 
          right="0" 
          bottom="0" 
          alignItems="center" 
          justifyContent="center"
          flexDirection="column"
          bg="rgba(0,0,0,0.8)"
          zIndex="10"
          p={6}
        >
          <Text color="red.500" fontSize="xl" mb={4}>Error Loading Dungeon</Text>
          <Text color="white" textAlign="center">{error}</Text>
          <Button 
            mt={6} 
            colorScheme="blue" 
            onClick={() => {
              setLoading(true);
              setError(null);
              requestFloorData();
            }}
          >
            Retry
          </Button>
        </Flex>
      ) : null}

      <canvas
        ref={canvasRef}
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
        style={{ display: 'block' }}
      />

      {/* Tooltip for hovered entity */}
      {hoveredEntity && tooltipPosition && (
        <Box
          position="absolute"
          top={`${tooltipPosition.y + 20}px`}
          left={`${tooltipPosition.x + 20}px`}
          bg="rgba(0, 0, 0, 0.8)"
          color="white"
          p={2}
          borderRadius="md"
          zIndex={10}
          maxWidth="250px"
        >
          <Text fontWeight="bold">{hoveredEntity.name}</Text>
          {hoveredEntity.description && (
            <Text fontSize="sm" mt={1}>{hoveredEntity.description}</Text>
          )}
          {hoveredEntity.health !== undefined && hoveredEntity.maxHealth !== undefined && (
            <Text fontSize="sm" mt={1}>
              Health: {hoveredEntity.health}/{hoveredEntity.maxHealth}
            </Text>
          )}
        </Box>
      )}

      {/* Map Legend Button */}
      <Button
        position="absolute"
        bottom="20px"
        right="20px"
        size="sm"
        colorScheme="blue"
        onClick={() => setIsLegendOpen(!isLegendOpen)}
      >
        {isLegendOpen ? 'Hide Legend' : 'Show Legend'}
      </Button>

      {/* Map Legend */}
      {isLegendOpen && <MapLegendComponent />}
    </Box>
  );
}; 