import React, { useState, useEffect, useCallback, useRef } from 'react';
import {
  Box,
  Grid,
  GridItem,
  Text,
  VStack,
  HStack,
  Button,
  useToast,
  Tooltip,
  Badge,
  Flex,
  Kbd,
  Divider,
} from '@chakra-ui/react';
import { RepeatIcon, InfoIcon } from '@chakra-ui/icons';

// Define tile types
type TileType = 'floor' | 'wall' | 'door' | 'stairs' | 'target';

// Define entity types
type EntityType = 'player' | 'path' | 'history';

// Interface for grid position
interface Position {
  x: number;
  y: number;
}

// Interface for grid cell
interface Cell {
  type: TileType;
  entities: EntityType[];
  walkable: boolean;
}

// Interface for component props
interface MovementDemoProps {
  gridSize?: number;
  wallDensity?: number;
  debug?: boolean;
}

// Movement modes
type MovementMode = 'cardinal' | 'diagonal';

// Direction vectors for movement
const DIRECTIONS = {
  cardinal: [
    { x: 0, y: -1 }, // Up
    { x: 1, y: 0 },  // Right
    { x: 0, y: 1 },  // Down
    { x: -1, y: 0 }, // Left
  ],
  diagonal: [
    { x: 0, y: -1 },  // Up
    { x: 1, y: -1 },  // Up-Right
    { x: 1, y: 0 },   // Right
    { x: 1, y: 1 },   // Down-Right
    { x: 0, y: 1 },   // Down
    { x: -1, y: 1 },  // Down-Left
    { x: -1, y: 0 },  // Left
    { x: -1, y: -1 }, // Up-Left
  ],
};

// Key mappings
const KEY_MAPPINGS: Record<string, Position> = {
  'ArrowUp': { x: 0, y: -1 },
  'ArrowRight': { x: 1, y: 0 },
  'ArrowDown': { x: 0, y: 1 },
  'ArrowLeft': { x: -1, y: 0 },
  'w': { x: 0, y: -1 },
  'd': { x: 1, y: 0 },
  's': { x: 0, y: 1 },
  'a': { x: -1, y: 0 },
};

// Diagonal key mappings
const DIAGONAL_KEY_MAPPINGS: Record<string, Position> = {
  ...KEY_MAPPINGS,
  'q': { x: -1, y: -1 }, // Up-Left
  'e': { x: 1, y: -1 },  // Up-Right
  'z': { x: -1, y: 1 },  // Down-Left
  'c': { x: 1, y: 1 },   // Down-Right
};

const MovementDemo: React.FC<MovementDemoProps> = ({
  gridSize = 15,
  wallDensity = 0.2,
  debug = false,
}) => {
  // State for the grid
  const [grid, setGrid] = useState<Cell[][]>([]);
  
  // State for player position
  const [playerPos, setPlayerPos] = useState<Position>({ x: 0, y: 0 });
  
  // State for target position
  const [targetPos, setTargetPos] = useState<Position | null>(null);
  
  // State for movement mode
  const [movementMode, setMovementMode] = useState<MovementMode>('cardinal');
  
  // State for movement history
  const [movementHistory, setMovementHistory] = useState<Position[]>([]);
  
  // State for path
  const [path, setPath] = useState<Position[]>([]);
  
  // Toast for notifications
  const toast = useToast();
  
  // Ref for the grid container
  const gridRef = useRef<HTMLDivElement>(null);
  
  // Initialize grid
  const initializeGrid = useCallback(() => {
    const newGrid: Cell[][] = [];
    
    // Create empty grid
    for (let y = 0; y < gridSize; y++) {
      const row: Cell[] = [];
      for (let x = 0; x < gridSize; x++) {
        row.push({
          type: 'floor',
          entities: [],
          walkable: true,
        });
      }
      newGrid.push(row);
    }
    
    // Add walls
    for (let y = 0; y < gridSize; y++) {
      for (let x = 0; x < gridSize; x++) {
        // Add walls around the edges
        if (x === 0 || y === 0 || x === gridSize - 1 || y === gridSize - 1) {
          newGrid[y][x].type = 'wall';
          newGrid[y][x].walkable = false;
          continue;
        }
        
        // Randomly add walls based on density
        if (Math.random() < wallDensity) {
          newGrid[y][x].type = 'wall';
          newGrid[y][x].walkable = false;
        }
      }
    }
    
    // Add doors and stairs
    const doorX = Math.floor(gridSize / 2);
    const doorY = gridSize - 1;
    if (newGrid[doorY][doorX]) {
      newGrid[doorY][doorX].type = 'door';
      newGrid[doorY][doorX].walkable = true;
    }
    
    const stairsX = Math.floor(gridSize / 4);
    const stairsY = Math.floor(gridSize / 4);
    if (newGrid[stairsY][stairsX]) {
      newGrid[stairsY][stairsX].type = 'stairs';
      newGrid[stairsY][stairsX].walkable = true;
    }
    
    // Find a valid starting position for the player
    let startX = 1;
    let startY = 1;
    
    // Make sure the starting position is walkable
    while (!newGrid[startY][startX].walkable) {
      startX = Math.floor(Math.random() * (gridSize - 2)) + 1;
      startY = Math.floor(Math.random() * (gridSize - 2)) + 1;
    }
    
    // Set player position
    setPlayerPos({ x: startX, y: startY });
    
    // Clear target, path, and history
    setTargetPos(null);
    setPath([]);
    setMovementHistory([]);
    
    return newGrid;
  }, [gridSize, wallDensity]);
  
  // Initialize grid on mount and when parameters change
  useEffect(() => {
    setGrid(initializeGrid());
  }, [initializeGrid]);
  
  // Update grid with player, path, and history
  useEffect(() => {
    if (grid.length === 0) return;
    
    // Create a new grid to avoid mutating state directly
    const newGrid = grid.map(row => row.map(cell => ({
      ...cell,
      entities: [] as EntityType[],
    })));
    
    // Add history to grid
    movementHistory.forEach(pos => {
      if (newGrid[pos.y] && newGrid[pos.y][pos.x]) {
        newGrid[pos.y][pos.x].entities.push('history' as EntityType);
      }
    });
    
    // Add path to grid
    path.forEach(pos => {
      if (newGrid[pos.y] && newGrid[pos.y][pos.x]) {
        newGrid[pos.y][pos.x].entities.push('path' as EntityType);
      }
    });
    
    // Add player to grid
    if (newGrid[playerPos.y] && newGrid[playerPos.y][playerPos.x]) {
      newGrid[playerPos.y][playerPos.x].entities.push('player' as EntityType);
    }
    
    // Add target to grid if it exists
    if (targetPos && newGrid[targetPos.y] && newGrid[targetPos.y][targetPos.x]) {
      newGrid[targetPos.y][targetPos.x].type = 'target';
    }
    
    setGrid(newGrid);
  }, [playerPos, targetPos, path, movementHistory]);
  
  // Handle keyboard movement
  const handleKeyDown = useCallback((e: KeyboardEvent) => {
    const keyMappings = movementMode === 'cardinal' ? KEY_MAPPINGS : DIAGONAL_KEY_MAPPINGS;
    const direction = keyMappings[e.key];
    
    if (!direction) return;
    
    const newPos = {
      x: playerPos.x + direction.x,
      y: playerPos.y + direction.y,
    };
    
    // Check if the new position is valid
    if (
      newPos.x >= 0 && 
      newPos.x < gridSize && 
      newPos.y >= 0 && 
      newPos.y < gridSize &&
      grid[newPos.y] && 
      grid[newPos.y][newPos.x] && 
      grid[newPos.y][newPos.x].walkable
    ) {
      // Add current position to history
      setMovementHistory(prev => [...prev, playerPos]);
      
      // Update player position
      setPlayerPos(newPos);
      
      // Clear path if we're not following it
      if (path.length > 0) {
        const nextPathPos = path[0];
        if (nextPathPos.x !== newPos.x || nextPathPos.y !== newPos.y) {
          setPath([]);
        } else {
          setPath(prev => prev.slice(1));
        }
      }
      
      // Check if we reached the target
      if (targetPos && newPos.x === targetPos.x && newPos.y === targetPos.y) {
        toast({
          title: 'Target reached!',
          status: 'success',
          duration: 2000,
          isClosable: true,
        });
        setTargetPos(null);
        setPath([]);
      }
    }
  }, [playerPos, grid, gridSize, movementMode, path, targetPos, toast]);
  
  // Add keyboard event listener
  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [handleKeyDown]);
  
  // Handle cell click
  const handleCellClick = (x: number, y: number) => {
    // If the cell is not walkable, do nothing
    if (!grid[y][x].walkable) return;
    
    // If the cell is the player's position, do nothing
    if (x === playerPos.x && y === playerPos.y) return;
    
    // Set the target position
    setTargetPos({ x, y });
    
    // Find path to target
    const newPath = findPath(playerPos, { x, y });
    setPath(newPath);
    
    // Focus the grid to enable keyboard events
    if (gridRef.current) {
      gridRef.current.focus();
    }
  };
  
  // Find path using A* algorithm
  const findPath = (start: Position, end: Position): Position[] => {
    // If start or end is not walkable, return empty path
    if (!grid[start.y][start.x].walkable || !grid[end.y][end.x].walkable) {
      return [];
    }
    
    // A* algorithm
    const openSet: Position[] = [start];
    const closedSet: Position[] = [];
    const cameFrom: Record<string, Position> = {};
    const gScore: Record<string, number> = {};
    const fScore: Record<string, number> = {};
    
    // Initialize scores
    const posKey = (pos: Position) => `${pos.x},${pos.y}`;
    gScore[posKey(start)] = 0;
    fScore[posKey(start)] = heuristic(start, end);
    
    while (openSet.length > 0) {
      // Find the node with the lowest fScore
      let current = openSet[0];
      let lowestFScore = fScore[posKey(current)];
      let currentIndex = 0;
      
      for (let i = 1; i < openSet.length; i++) {
        const pos = openSet[i];
        const score = fScore[posKey(pos)];
        
        if (score < lowestFScore) {
          lowestFScore = score;
          current = pos;
          currentIndex = i;
        }
      }
      
      // If we reached the end, reconstruct the path
      if (current.x === end.x && current.y === end.y) {
        const path: Position[] = [];
        let currentPos = current;
        
        while (posKey(currentPos) in cameFrom) {
          path.unshift(currentPos);
          currentPos = cameFrom[posKey(currentPos)];
        }
        
        return path;
      }
      
      // Remove current from openSet and add to closedSet
      openSet.splice(currentIndex, 1);
      closedSet.push(current);
      
      // Get neighbors
      const directions = movementMode === 'cardinal' ? DIRECTIONS.cardinal : DIRECTIONS.diagonal;
      
      for (const dir of directions) {
        const neighbor = {
          x: current.x + dir.x,
          y: current.y + dir.y,
        };
        
        // Skip if out of bounds
        if (
          neighbor.x < 0 || 
          neighbor.x >= gridSize || 
          neighbor.y < 0 || 
          neighbor.y >= gridSize
        ) {
          continue;
        }
        
        // Skip if not walkable
        if (!grid[neighbor.y][neighbor.x].walkable) {
          continue;
        }
        
        // Skip if in closedSet
        if (closedSet.some(pos => pos.x === neighbor.x && pos.y === neighbor.y)) {
          continue;
        }
        
        // Calculate tentative gScore
        const tentativeGScore = gScore[posKey(current)] + 1;
        
        // Add to openSet if not already there
        const inOpenSet = openSet.some(pos => pos.x === neighbor.x && pos.y === neighbor.y);
        if (!inOpenSet) {
          openSet.push(neighbor);
        } else if (tentativeGScore >= gScore[posKey(neighbor)]) {
          continue;
        }
        
        // This path is better, record it
        cameFrom[posKey(neighbor)] = current;
        gScore[posKey(neighbor)] = tentativeGScore;
        fScore[posKey(neighbor)] = gScore[posKey(neighbor)] + heuristic(neighbor, end);
      }
    }
    
    // No path found
    return [];
  };
  
  // Heuristic function for A* (Manhattan distance)
  const heuristic = (a: Position, b: Position): number => {
    return Math.abs(a.x - b.x) + Math.abs(a.y - b.y);
  };
  
  // Toggle movement mode
  const toggleMovementMode = () => {
    setMovementMode(prev => prev === 'cardinal' ? 'diagonal' : 'cardinal');
    setPath([]);
  };
  
  // Generate a new map
  const generateNewMap = useCallback(() => {
    const newGrid = initializeGrid();
    setGrid(newGrid);
    
    // Find a random walkable position for the player
    const walkableCells: Position[] = [];
    newGrid.forEach((row, y) => {
      row.forEach((cell, x) => {
        if (cell.walkable) {
          walkableCells.push({ x, y });
        }
      });
    });
    
    if (walkableCells.length > 0) {
      const randomIndex = Math.floor(Math.random() * walkableCells.length);
      setPlayerPos(walkableCells[randomIndex]);
    }
    
    // Clear target and path
    setTargetPos(null);
    setPath([]);
    setMovementHistory([]);
    
    toast({
      title: "New map generated",
      status: "info",
      duration: 2000,
      isClosable: true,
    });
  }, [initializeGrid, toast]);
  
  // Clear history
  const clearHistory = () => {
    setMovementHistory([]);
  };
  
  // Get background color for a cell based on its type
  const getCellBackground = (type: TileType): string => {
    switch (type) {
      case 'floor':
        return 'gray.700';
      case 'wall':
        return 'gray.900';
      case 'door':
        return 'yellow.900';
      case 'stairs':
        return 'blue.900';
      case 'target':
        return 'green.900';
      default:
        return 'gray.700';
    }
  };
  
  // Render tile
  const renderTile = (cell: Cell, x: number, y: number) => {
    // Determine background color based on tile type
    let bgColor = 'gray.700';
    let symbol = ' ';
    
    switch (cell.type) {
      case 'floor':
        bgColor = 'gray.700';
        symbol = '.';
        break;
      case 'wall':
        bgColor = 'gray.900';
        symbol = '#';
        break;
      case 'door':
        bgColor = 'yellow.800';
        symbol = '+';
        break;
      case 'stairs':
        bgColor = 'blue.800';
        symbol = '>';
        break;
      case 'target':
        bgColor = 'red.700';
        symbol = 'X';
        break;
    }
    
    // Determine entity color and symbol
    let entityColor = '';
    let entitySymbol = '';
    
    if (cell.entities.includes('player')) {
      entityColor = 'cyan.300';
      entitySymbol = '@';
    } else if (cell.entities.includes('path')) {
      entityColor = 'green.400';
      entitySymbol = '*';
    } else if (cell.entities.includes('history')) {
      entityColor = 'purple.400';
      entitySymbol = '·';
    }
    
    return (
      <GridItem
        key={`${x}-${y}`}
        w="100%"
        h="100%"
        bg={bgColor}
        border="1px solid"
        borderColor="gray.800"
        display="flex"
        alignItems="center"
        justifyContent="center"
        position="relative"
        onClick={() => handleCellClick(x, y)}
        cursor={cell.walkable ? 'pointer' : 'not-allowed'}
        className="grid-cell"
      >
        {debug && (
          <Text
            fontSize="xs"
            color="gray.500"
            position="absolute"
            top="1px"
            left="1px"
          >
            {x},{y}
          </Text>
        )}
        
        {/* Tile symbol */}
        {debug && (
          <Text
            fontSize="md"
            fontFamily="monospace"
            color="gray.400"
          >
            {symbol}
          </Text>
        )}
        
        {/* Entity */}
        {entitySymbol && (
          <Text
            fontSize="xl"
            fontFamily="monospace"
            color={entityColor}
            fontWeight="bold"
            position="absolute"
            className={cell.entities.includes('player') ? 'player-character' : ''}
          >
            {entitySymbol}
          </Text>
        )}
      </GridItem>
    );
  };
  
  return (
    <VStack spacing={4} align="stretch" width="100%">
      <HStack spacing={4} justify="space-between">
        <HStack>
          <Text>Movement Mode: {movementMode === 'cardinal' ? 'Cardinal' : 'Diagonal'}</Text>
          <Button size="sm" onClick={toggleMovementMode}>Toggle Mode</Button>
        </HStack>
        <Button 
          leftIcon={<RepeatIcon />} 
          size="sm" 
          onClick={generateNewMap}
        >
          New Map
        </Button>
      </HStack>
      
      <Divider />
      
      <Box 
        ref={gridRef} 
        className="grid-container" 
        data-testid="movement-grid"
        position="relative" 
        width="100%" 
        height="400px" 
        overflowX="auto" 
        overflowY="auto" 
        borderWidth="1px" 
        borderColor="gray.600" 
        borderRadius="md"
      >
        <Grid 
          templateColumns={`repeat(${gridSize}, 1fr)`} 
          gap={0} 
          width={`${gridSize * 30}px`} 
          height={`${gridSize * 30}px`}
        >
          {grid.map((row, y) => 
            row.map((cell, x) => (
              <GridItem 
                key={`${x}-${y}`} 
                className="grid-cell"
                data-testid="grid-cell"
                width="30px" 
                height="30px" 
                bg={getCellBackground(cell.type)} 
                border="1px solid" 
                borderColor="gray.700"
                position="relative"
                onClick={() => handleCellClick(x, y)}
                cursor="pointer"
                _hover={{ opacity: 0.8 }}
              >
                {/* Player character */}
                {playerPos.x === x && playerPos.y === y && (
                  <Box 
                    className="player-character"
                    data-testid="player-character"
                    position="absolute" 
                    top="0" 
                    left="0" 
                    width="100%" 
                    height="100%" 
                    display="flex" 
                    alignItems="center" 
                    justifyContent="center"
                  >
                    <Text fontSize="xl" fontWeight="bold" color="yellow.300">@</Text>
                  </Box>
                )}
                
                {/* Target marker */}
                {targetPos && targetPos.x === x && targetPos.y === y && (
                  <Box 
                    className="target-marker"
                    data-testid="target-marker"
                    position="absolute" 
                    top="0" 
                    left="0" 
                    width="100%" 
                    height="100%" 
                    display="flex" 
                    alignItems="center" 
                    justifyContent="center"
                    zIndex={1}
                  >
                    <Text fontSize="xl" fontWeight="bold" color="green.300">X</Text>
                  </Box>
                )}
                
                {/* Path marker */}
                {path.some(pos => pos.x === x && pos.y === y) && !(playerPos.x === x && playerPos.y === y) && !(targetPos && targetPos.x === x && targetPos.y === y) && (
                  <Box 
                    className="path-marker"
                    data-testid="path-marker"
                    position="absolute" 
                    top="0" 
                    left="0" 
                    width="100%" 
                    height="100%" 
                    display="flex" 
                    alignItems="center" 
                    justifyContent="center"
                    zIndex={1}
                  >
                    <Text fontSize="md" color="blue.300">•</Text>
                  </Box>
                )}
                
                {/* Movement history */}
                {movementHistory.some(pos => pos.x === x && pos.y === y) && !(playerPos.x === x && playerPos.y === y) && (
                  <Box 
                    className="history-marker"
                    data-testid="history-marker"
                    position="absolute" 
                    top="0" 
                    left="0" 
                    width="100%" 
                    height="100%" 
                    bg="rgba(255, 255, 255, 0.1)"
                    zIndex={0}
                  />
                )}
                
                {/* Debug symbols */}
                {debug && (
                  <Tooltip label={`(${x},${y}) - ${cell.type}`}>
                    <Text fontSize="xs" color="gray.400" position="absolute" bottom="1px" right="1px">
                      {x},{y}
                    </Text>
                  </Tooltip>
                )}
              </GridItem>
            ))
          )}
        </Grid>
      </Box>
      
      <Divider />
      
      <Box>
        <Text fontSize="sm" color="gray.400">
          <Kbd>Arrow keys</Kbd> to move | <Kbd>Click</Kbd> to set target
        </Text>
      </Box>
    </VStack>
  );
};

export default MovementDemo; 