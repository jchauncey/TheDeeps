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
  
  // Animation state for smooth movement
  const [animatingPlayerPos, setAnimatingPlayerPos] = useState<Position | null>(null);
  const [animationProgress, setAnimationProgress] = useState(1); // 0 to 1
  const animationRef = useRef<number | null>(null);
  const lastUpdateTimeRef = useRef<number>(0);
  const ANIMATION_DURATION = 100; // milliseconds - reduced for faster response
  
  // Key state tracking for continuous movement
  const [keysPressed, setKeysPressed] = useState<Set<string>>(new Set());
  const keyIntervalRef = useRef<number | null>(null);
  const KEY_REPEAT_DELAY = 150; // milliseconds between key repeats - increased for better animation completion
  const lastMoveTimeRef = useRef<number>(0);
  const movementAllowedRef = useRef<boolean>(true);
  const isMovingRef = useRef<boolean>(false);
  
  // Caching for map rendering
  const mapCacheRef = useRef<HTMLCanvasElement | null>(null);
  const lastMapPositionRef = useRef<Position | null>(null);
  const needsFullRedrawRef = useRef<boolean>(true);
  
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const toast = useToast();

  // Process floor data when it changes
  useEffect(() => {
    if (floorData) {
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
        
        // Check if this is just a player movement update or a full floor update
        const isJustMovement = floor && 
                              playerPos && 
                              floor.level === floorData.floor.level &&
                              JSON.stringify(floor.entities) === JSON.stringify(floorData.floor.entities) &&
                              JSON.stringify(floor.items) === JSON.stringify(floorData.floor.items) &&
                              (playerPos.x !== floorData.playerPosition.x || 
                               playerPos.y !== floorData.playerPosition.y);
        
        // If it's a full floor update or first load, mark for full redraw
        if (!isJustMovement || !floor) {
          needsFullRedrawRef.current = true;
          setFloor(floorData.floor);
        }
        
        // Check if this is a player movement update
        if (playerPos && 
            (playerPos.x !== floorData.playerPosition.x || 
             playerPos.y !== floorData.playerPosition.y)) {
          
          // Start animation from current position to new position
          setAnimatingPlayerPos(playerPos);
          setAnimationProgress(0);
          
          // Cancel any existing animation
          if (animationRef.current) {
            cancelAnimationFrame(animationRef.current);
          }
          
          // Start animation loop
          lastUpdateTimeRef.current = performance.now();
          animationRef.current = requestAnimationFrame(animateMovement);
        }
        
        // Always update player position
        setPlayerPos(floorData.playerPosition);
        setCurrentFloor(floorData.currentFloor);
        setLoading(false);
        setError(null);
        
        // Only force a redraw if we're not animating
        if (animationProgress >= 1) {
          // Force a redraw after a short delay to ensure the canvas is ready
          setTimeout(() => {
            drawFloor();
          }, 100);
        }
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
          
          // Check if this is just a movement update
          if (floorData.isMovementUpdate) {
            // For movement updates, we only need to update the player position
            if (playerPos && 
                (playerPos.x !== floorData.playerPosition.x || 
                 playerPos.y !== floorData.playerPosition.y)) {
              
              // Start animation from current position to new position
              setAnimatingPlayerPos(playerPos);
              setAnimationProgress(0);
              
              // Cancel any existing animation
              if (animationRef.current) {
                cancelAnimationFrame(animationRef.current);
              }
              
              // Start animation loop
              lastUpdateTimeRef.current = performance.now();
              animationRef.current = requestAnimationFrame(animateMovement);
              
              // Update player position
              setPlayerPos(floorData.playerPosition);
            }
          } else {
            // For full floor updates
            needsFullRedrawRef.current = true;
            setFloor(floorData.floor);
            setPlayerPos(floorData.playerPosition);
            setCurrentFloor(floorData.currentFloor);
            setLoading(false);
            setError(null);
          }
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
  }, [playerPos]);

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
    if (!floor || !playerPos || !canvasRef.current) {
      console.error('Cannot draw floor - missing required data');
      return;
    }
    
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) {
      console.error('Failed to get canvas context');
      return;
    }
    
    // Calculate the interpolated player position for smooth animation
    let displayPlayerPos = { ...playerPos };
    
    // If we're animating, interpolate between the start and end positions
    if (animatingPlayerPos && animationProgress < 1) {
      displayPlayerPos = {
        x: animatingPlayerPos.x + (playerPos.x - animatingPlayerPos.x) * animationProgress,
        y: animatingPlayerPos.y + (playerPos.y - animatingPlayerPos.y) * animationProgress
      };
    }
    
    // Calculate the ideal number of tiles to show
    const idealVisibleTilesX = 30; // Show more tiles horizontally
    const idealVisibleTilesY = 20; // Show more tiles vertically
    
    // Calculate tile size based on viewport and ideal number of tiles
    const tileSize = Math.min(
      canvas.width / idealVisibleTilesX,
      canvas.height / idealVisibleTilesY
    );
    
    // Calculate viewport boundaries to only render what's visible
    const visibleTilesX = Math.ceil(canvas.width / tileSize);
    const visibleTilesY = Math.ceil(canvas.height / tileSize);
    
    // Center the viewport on the interpolated player position with improved edge handling
    let startX = Math.max(0, Math.floor(displayPlayerPos.x) - Math.floor(visibleTilesX / 2));
    let startY = Math.max(0, Math.floor(displayPlayerPos.y) - Math.floor(visibleTilesY / 2));
    
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
    
    // Calculate offset to center the map in the viewport
    // Use fixed tile counts to maintain consistent sizing
    const tilesShownX = endX - startX;
    const tilesShownY = endY - startY;
    const offsetX = (canvas.width - tilesShownX * tileSize) / 2;
    const offsetY = (canvas.height - tilesShownY * tileSize) / 2;
    
    // Check if we need to redraw the entire map or just update the player position
    const viewportChanged = !lastMapPositionRef.current || 
                           Math.abs(startX - lastMapPositionRef.current.x) > 0 || 
                           Math.abs(startY - lastMapPositionRef.current.y) > 0;
    
    // If we need a full redraw or the viewport has changed, redraw the entire map
    if (needsFullRedrawRef.current || viewportChanged || !mapCacheRef.current) {
      // Create or reuse the map cache canvas
      if (!mapCacheRef.current) {
        mapCacheRef.current = document.createElement('canvas');
      }
      
      // Set the map cache canvas dimensions
      mapCacheRef.current.width = canvas.width;
      mapCacheRef.current.height = canvas.height;
      
      const mapCtx = mapCacheRef.current.getContext('2d');
      if (!mapCtx) {
        console.error('Failed to get map cache context');
        return;
      }
      
      // Clear the map cache
      mapCtx.clearRect(0, 0, mapCacheRef.current.width, mapCacheRef.current.height);
      
      try {
        // Draw tiles and static entities to the map cache
        for (let y = startY; y < endY; y++) {
          for (let x = startX; x < endX; x++) {
            if (!floor.tiles[y] || !floor.tiles[y][x]) {
              continue;
            }
            
            const tile = floor.tiles[y][x];
            const screenX = (x - startX) * tileSize + offsetX;
            const screenY = (y - startY) * tileSize + offsetY;
            
            // Draw tile background
            mapCtx.fillStyle = TILE_COLORS[tile.type as keyof typeof TILE_COLORS] || '#000';
            mapCtx.fillRect(screenX, screenY, tileSize, tileSize);
            
            // Draw tile border
            mapCtx.strokeStyle = '#222';
            mapCtx.lineWidth = 1;
            mapCtx.strokeRect(screenX, screenY, tileSize, tileSize);
            
            // Draw item if present
            if (tile.item) {
              drawItem(mapCtx, tile.item, screenX, screenY, tileSize);
            }
          }
        }
        
        // Draw non-player entities to the map cache
        if (floor.entities && floor.entities.length > 0) {
          for (const entity of floor.entities) {
            // Skip the player entity - we'll draw it separately with animation
            if (entity.type === 'player' && entity.id === floorData?.playerData?.id) {
              continue;
            }
            
            // Only draw entities within the viewport
            if (
              entity.position.x >= startX && 
              entity.position.x < endX && 
              entity.position.y >= startY && 
              entity.position.y < endY
            ) {
              const screenX = (entity.position.x - startX) * tileSize + offsetX;
              const screenY = (entity.position.y - startY) * tileSize + offsetY;
              drawEntity(mapCtx, entity, screenX, screenY, tileSize);
            }
          }
        }
        
        // Update the last map position and mark that we don't need a full redraw next time
        lastMapPositionRef.current = { x: startX, y: startY };
        needsFullRedrawRef.current = false;
      } catch (error) {
        console.error('Error drawing map cache:', error);
        return;
      }
    }
    
    // Clear the main canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    // Draw the cached map to the main canvas
    if (mapCacheRef.current) {
      ctx.drawImage(mapCacheRef.current, 0, 0);
    }
    
    // Draw the player with interpolated position
    if (displayPlayerPos) {
      // Calculate the screen position with sub-tile precision for smooth animation
      const screenX = (displayPlayerPos.x - startX) * tileSize + offsetX;
      const screenY = (displayPlayerPos.y - startY) * tileSize + offsetY;
      
      // Create a player entity object
      const playerEntity: Entity = {
        id: floorData?.playerData?.id || 'player',
        type: 'player',
        name: floorData?.playerData?.name || 'Player',
        position: displayPlayerPos,
        characterClass: floorData?.playerData?.characterClass || 'default',
        health: floorData?.playerData?.health || 100,
        maxHealth: floorData?.playerData?.maxHealth || 100
      };
      
      // Draw the player
      drawEntity(ctx, playerEntity, screenX, screenY, tileSize);
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

    // Add key to pressed keys set
    setKeysPressed(prev => {
      const newSet = new Set(prev);
      newSet.add(e.key);
      return newSet;
    });
    
    // Process the key press immediately if it's a movement key
    if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(e.key) && 
        !isMovingRef.current && movementAllowedRef.current) {
      processMovementKeys(e.key);
    }
    
    // If this is the first key press, start the interval for continuous movement
    if (!keyIntervalRef.current) {
      // Set up interval for continuous movement while keys are held down
      keyIntervalRef.current = window.setInterval(() => {
        // Only process keys if we're not currently animating
        if (!isMovingRef.current && movementAllowedRef.current && keysPressed.size > 0) {
          // Get the most recently pressed movement key
          const movementKeys = Array.from(keysPressed).filter(key => 
            ['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(key)
          );
          
          if (movementKeys.length > 0) {
            const lastKey = movementKeys[movementKeys.length - 1];
            processMovementKeys(lastKey);
          }
        }
      }, KEY_REPEAT_DELAY);
    }
  };
  
  // Handle key up events
  const handleKeyUp = (e: React.KeyboardEvent<HTMLDivElement>) => {
    // Remove key from pressed keys set
    setKeysPressed(prev => {
      const newSet = new Set(prev);
      newSet.delete(e.key);
      return newSet;
    });
    
    // If no keys are pressed, clear the interval
    if (keysPressed.size === 0 && keyIntervalRef.current) {
      clearInterval(keyIntervalRef.current);
      keyIntervalRef.current = null;
    }
  };
  
  // Animation function for smooth movement
  const animateMovement = (timestamp: number) => {
    if (!animatingPlayerPos || !playerPos) {
      setAnimationProgress(1);
      isMovingRef.current = false;
      return;
    }
    
    isMovingRef.current = true;
    
    const elapsed = timestamp - lastUpdateTimeRef.current;
    const newProgress = Math.min(1, elapsed / ANIMATION_DURATION);
    
    setAnimationProgress(newProgress);
    
    if (newProgress < 1) {
      // Continue animation
      animationRef.current = requestAnimationFrame(animateMovement);
      // Redraw with animation
      drawFloor();
    } else {
      // Animation complete
      setAnimatingPlayerPos(null);
      isMovingRef.current = false;
      // Final redraw
      drawFloor();
      
      // If keys are still pressed, allow the next movement immediately
      if (keysPressed.size > 0) {
        movementAllowedRef.current = true;
        
        // Process the next movement immediately if keys are still pressed
        const movementKeys = Array.from(keysPressed).filter(key => 
          ['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(key)
        );
        
        if (movementKeys.length > 0) {
          const lastKey = movementKeys[movementKeys.length - 1];
          processMovementKeys(lastKey);
        }
      }
    }
  };
  
  // Process movement keys
  const processMovementKeys = (key: string) => {
    // Check if movement is allowed (not too soon after last move or during animation)
    const now = performance.now();
    if (!movementAllowedRef.current || isMovingRef.current) {
      return;
    }
    
    // Temporarily disable movement until animation completes
    movementAllowedRef.current = false;
    
    // Handle movement with arrow keys - send to server for validation
    // The server will validate the move and send back updated floor data if valid
    switch (key) {
      case 'ArrowUp':
        sendWebSocketMessage({ type: 'move', direction: 'up' });
        lastMoveTimeRef.current = now;
        break;
      case 'ArrowDown':
        sendWebSocketMessage({ type: 'move', direction: 'down' });
        lastMoveTimeRef.current = now;
        break;
      case 'ArrowLeft':
        sendWebSocketMessage({ type: 'move', direction: 'left' });
        lastMoveTimeRef.current = now;
        break;
      case 'ArrowRight':
        sendWebSocketMessage({ type: 'move', direction: 'right' });
        lastMoveTimeRef.current = now;
        break;
      case ' ': // Space bar for attack
        sendWebSocketMessage({ type: 'action', action: 'attack' });
        movementAllowedRef.current = true; // Re-enable movement for non-movement actions
        break;
      case 'p': // 'p' for pickup
      case 'P':
        sendWebSocketMessage({ type: 'action', action: 'pickup' });
        movementAllowedRef.current = true; // Re-enable movement for non-movement actions
        break;
      default:
        movementAllowedRef.current = true; // Re-enable movement for unknown keys
        break;
    }
  };
  
  // Clean up interval on unmount
  useEffect(() => {
    return () => {
      if (keyIntervalRef.current) {
        clearInterval(keyIntervalRef.current);
      }
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, []);
  
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

  // Force a full redraw when the floor changes
  useEffect(() => {
    if (floor) {
      needsFullRedrawRef.current = true;
      mapCacheRef.current = null;
      lastMapPositionRef.current = null;
    }
  }, [currentFloor]);

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
      onKeyUp={handleKeyUp}
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
          <Spinner size="xl" color="blue.500" thickness="4px" mb={4} />
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