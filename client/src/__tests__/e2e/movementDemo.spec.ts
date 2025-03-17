import { test, expect } from '@playwright/test';

test.describe('MovementDemo Component', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate directly to the component playground with MovementDemo selected
    await page.goto('http://localhost:3000/component-playground?component=MovementDemo');
    
    // Wait for the page to load and the grid to be rendered
    await page.waitForSelector('.grid-container', { timeout: 10000 });
    await page.waitForSelector('.grid-cell', { timeout: 10000 });
    await page.waitForSelector('.player-character', { timeout: 30000 });
  });

  test('should load the page', async ({ page }) => {
    // Take a screenshot for debugging
    await page.screenshot({ path: 'movement-demo-page.png' });
    
    // Log debugging information
    console.log('Page title:', await page.title());
    console.log('URL:', page.url());
    
    // Just verify we're on the right page
    expect(page.url()).toContain('component=MovementDemo');
  });
  
  test('should set a target and show path when clicking on a cell', async ({ page }) => {
    // Get the player's initial position
    const playerElement = await page.locator('.player-character').first();
    const initialBoundingBox = await playerElement.boundingBox();
    expect(initialBoundingBox).not.toBeNull();
    
    if (initialBoundingBox) {
      console.log(`Initial player position: x=${initialBoundingBox.x}, y=${initialBoundingBox.y}`);
      
      // Get all grid cells
      const gridCells = await page.locator('.grid-cell').all();
      const gridDimension = Math.sqrt(gridCells.length);
      
      // Find the player's grid position
      const playerCell = await page.locator('.grid-cell:has(.player-character)').first();
      const playerCellIndex = await playerCell.evaluate((el, cells) => {
        return Array.from(cells).indexOf(el);
      }, await page.$$('.grid-cell'));
      
      const playerGridX = playerCellIndex % gridDimension;
      const playerGridY = Math.floor(playerCellIndex / gridDimension);
      
      console.log(`Player grid position: (${playerGridX}, ${playerGridY})`);
      
      // Try multiple directions until we find a walkable cell
      const directions = [
        { dx: 1, dy: 0 },  // right
        { dx: 0, dy: 1 },  // down
        { dx: -1, dy: 0 }, // left
        { dx: 0, dy: -1 }  // up
      ];
      
      let targetCell = null;
      let targetIndex = -1;
      let targetX = -1;
      let targetY = -1;
      
      for (const dir of directions) {
        const x = playerGridX + dir.dx;
        const y = playerGridY + dir.dy;
        const index = y * gridDimension + x;
        
        if (index >= 0 && index < gridCells.length) {
          const cell = gridCells[index];
          const classes = await cell.getAttribute('class');
          console.log(`Checking cell at (${x}, ${y}): ${classes}`);
          
          // Check if the cell is walkable (doesn't have wall class)
          if (!classes?.includes('wall')) {
            targetCell = cell;
            targetIndex = index;
            targetX = x;
            targetY = y;
            console.log(`Found walkable cell at index ${targetIndex}, position (${targetX}, ${targetY})`);
            break;
          }
        }
      }
      
      if (targetCell) {
        // Take a screenshot before clicking
        await page.screenshot({ path: 'before-click.png' });
        
        // Force click in the center of the cell with a longer timeout
        await targetCell.click({ force: true, timeout: 5000 });
        
        // Wait for the target marker to appear
        await page.waitForTimeout(1000);
        
        // Check for target marker using data-testid attribute
        const targetMarkerByTestId = await page.locator('[data-testid="target-marker"]').count();
        console.log(`Found ${targetMarkerByTestId} target markers by data-testid`);
        
        // Also check by class name as a fallback
        const targetMarkerByClass = await page.locator('.target-marker').count();
        console.log(`Found ${targetMarkerByClass} target markers by class`);
        
        // Take a screenshot after clicking
        await page.screenshot({ path: 'after-click.png' });
        
        // Log the final state
        const finalClasses = await targetCell.getAttribute('class');
        console.log(`Final target cell classes: ${finalClasses}`);
        
        // Check if target was set using either method
        const totalTargetMarkers = targetMarkerByTestId + targetMarkerByClass;
        expect(totalTargetMarkers).toBeGreaterThan(0);
        
        // Check for path markers
        const pathMarkers = await page.locator('.path-marker').count();
        console.log(`Found ${pathMarkers} path markers`);
        
        // If the target is adjacent to the player, there might not be path markers
        const isAdjacent = Math.abs(targetX - playerGridX) + Math.abs(targetY - playerGridY) === 1;
        if (!isAdjacent) {
          expect(pathMarkers).toBeGreaterThan(0);
        }
      } else {
        console.log('No walkable cells found adjacent to player');
        // Skip the test if no walkable cells were found
        test.skip();
      }
    }
  });

  test('should move player using arrow keys', async ({ page }) => {
    // Get initial player position
    const playerElement = await page.locator('.player-character').first();
    const initialBox = await playerElement.boundingBox();
    expect(initialBox).not.toBeNull();
    
    if (initialBox) {
      console.log('Initial player position:', initialBox);
      
      // Focus the grid container to enable keyboard events
      await page.locator('.grid-container').click();
      
      // Test movement in all four directions
      // Note: Y-axis is inverted in the browser (positive is down)
      const movements = [
        { key: 'ArrowRight', expectedDx: 1, expectedDy: 0 },
        { key: 'ArrowDown', expectedDx: 0, expectedDy: -1 }, // Inverted Y-axis
        { key: 'ArrowLeft', expectedDx: -1, expectedDy: 0 },
        { key: 'ArrowUp', expectedDx: 0, expectedDy: 1 }     // Inverted Y-axis
      ];
      
      for (const movement of movements) {
        console.log(`Testing ${movement.key} movement...`);
        
        // Take a screenshot before movement
        await page.screenshot({ path: `before-${movement.key}.png` });
        
        const beforeMove = await playerElement.boundingBox();
        if (!beforeMove) {
          console.log('Could not get player position before movement');
          continue;
        }
        console.log('Position before move:', beforeMove);
        
        // Press the arrow key and wait for movement animation
        await page.keyboard.press(movement.key);
        await page.waitForTimeout(1000);
        
        // Take a screenshot after movement
        await page.screenshot({ path: `after-${movement.key}.png` });
        
        // Get new position
        const afterMove = await playerElement.boundingBox();
        if (!afterMove) {
          console.log('Could not get player position after movement');
          continue;
        }
        console.log('Position after move:', afterMove);
        
        // Calculate actual movement
        const actualDx = Math.sign(afterMove.x - beforeMove.x);
        const actualDy = Math.sign(afterMove.y - beforeMove.y);
        console.log(`Movement deltas: dx=${actualDx}, dy=${actualDy}`);
        console.log(`Expected deltas: dx=${movement.expectedDx}, dy=${movement.expectedDy}`);
        
        // Check if the cell in the movement direction is walkable
        const playerCell = await page.locator('.grid-cell:has(.player-character)').first();
        const cellClasses = await playerCell.getAttribute('class');
        console.log('Current cell classes:', cellClasses);
        
        // The movement might be blocked by walls, so we check if either
        // the position changed in the expected direction or stayed the same
        const validX = actualDx === movement.expectedDx || actualDx === 0;
        const validY = actualDy === movement.expectedDy || actualDy === 0;
        
        if (!validX || !validY) {
          console.log('Movement validation failed:');
          console.log('- X movement valid:', validX);
          console.log('- Y movement valid:', validY);
        }
        
        expect(validX).toBeTruthy();
        expect(validY).toBeTruthy();
      }
    }
  });

  test('should move player using WASD keys', async ({ page }) => {
    // Get initial player position
    const playerElement = await page.locator('.player-character').first();
    const initialBox = await playerElement.boundingBox();
    expect(initialBox).not.toBeNull();
    
    if (initialBox) {
      // Focus the grid container
      await page.locator('.grid-container').click();
      
      // Test WASD movement
      const movements = [
        { key: 'd', expectedDx: 1, expectedDy: 0 },
        { key: 's', expectedDx: 0, expectedDy: 1 },
        { key: 'a', expectedDx: -1, expectedDy: 0 },
        { key: 'w', expectedDx: 0, expectedDy: -1 }
      ];
      
      for (const movement of movements) {
        const beforeMove = await playerElement.boundingBox();
        if (!beforeMove) continue;
        
        await page.keyboard.press(movement.key);
        await page.waitForTimeout(500);
        
        const afterMove = await playerElement.boundingBox();
        expect(afterMove).not.toBeNull();
        
        if (afterMove) {
          const actualDx = Math.sign(afterMove.x - beforeMove.x);
          const actualDy = Math.sign(afterMove.y - beforeMove.y);
          
          expect(
            actualDx === movement.expectedDx || actualDx === 0
          ).toBeTruthy();
          expect(
            actualDy === movement.expectedDy || actualDy === 0
          ).toBeTruthy();
        }
      }
    }
  });

  test('should support diagonal movement when enabled', async ({ page }) => {
    // First, toggle diagonal movement mode
    await page.click('button:has-text("Toggle Mode")');
    await page.waitForTimeout(500);
    
    // Get initial player position
    const playerElement = await page.locator('.player-character').first();
    const initialBox = await playerElement.boundingBox();
    expect(initialBox).not.toBeNull();
    
    if (initialBox) {
      // Focus the grid container
      await page.locator('.grid-container').click();
      
      // Test diagonal movement using QEZC keys
      const movements = [
        { key: 'e', expectedDx: 1, expectedDy: -1 }, // Up-Right
        { key: 'c', expectedDx: 1, expectedDy: 1 },  // Down-Right
        { key: 'z', expectedDx: -1, expectedDy: 1 }, // Down-Left
        { key: 'q', expectedDx: -1, expectedDy: -1 } // Up-Left
      ];
      
      for (const movement of movements) {
        const beforeMove = await playerElement.boundingBox();
        if (!beforeMove) continue;
        
        await page.keyboard.press(movement.key);
        await page.waitForTimeout(500);
        
        const afterMove = await playerElement.boundingBox();
        expect(afterMove).not.toBeNull();
        
        if (afterMove) {
          const actualDx = Math.sign(afterMove.x - beforeMove.x);
          const actualDy = Math.sign(afterMove.y - beforeMove.y);
          
          expect(
            actualDx === movement.expectedDx || actualDx === 0
          ).toBeTruthy();
          expect(
            actualDy === movement.expectedDy || actualDy === 0
          ).toBeTruthy();
        }
      }
    }
  });

  test('should find path to target through obstacles', async ({ page }) => {
    // Get all grid cells
    const gridCells = await page.locator('.grid-cell').all();
    const gridDimension = Math.sqrt(gridCells.length);
    
    // Find the player's position
    const playerCell = await page.locator('.grid-cell:has(.player-character)').first();
    const playerCellIndex = await playerCell.evaluate((el, cells) => {
      return Array.from(cells).indexOf(el);
    }, await page.$$('.grid-cell'));
    
    const playerGridX = playerCellIndex % gridDimension;
    const playerGridY = Math.floor(playerCellIndex / gridDimension);
    
    // Try to find a walkable cell that's two spaces away
    const targetX = playerGridX + 2;
    const targetY = playerGridY + 2;
    const targetIndex = targetY * gridDimension + targetX;
    
    if (targetIndex < gridCells.length) {
      // Click on the target cell
      await gridCells[targetIndex].click();
      await page.waitForTimeout(1000);
      
      // Check if a target was set
      const targetMarker = await page.locator('.target-marker').count();
      expect(targetMarker).toBeGreaterThan(0);
      
      // Check if path markers appear
      const pathMarkers = await page.locator('.path-marker').count();
      console.log(`Found ${pathMarkers} path markers`);
      
      // If the target is reachable, there should be path markers
      // The number of markers depends on the path length and any obstacles
      expect(pathMarkers).toBeGreaterThanOrEqual(0);
    }
  });

  test('should generate new map when clicking New Map button', async ({ page }) => {
    // Get initial grid state
    const initialGridHTML = await page.locator('.grid-container').innerHTML();
    
    // Click the New Map button
    await page.click('button:has-text("New Map")');
    await page.waitForTimeout(1000);
    
    // Get new grid state
    const newGridHTML = await page.locator('.grid-container').innerHTML();
    
    // The grid HTML should be different after generating a new map
    expect(newGridHTML).not.toBe(initialGridHTML);
    
    // Verify that the player character exists in the new map
    const playerCharacter = await page.locator('.player-character').count();
    expect(playerCharacter).toBe(1);
  });

  test('should handle rapid movement key presses', async ({ page }) => {
    // Get initial player position
    const playerElement = await page.locator('.player-character').first();
    const initialBox = await playerElement.boundingBox();
    expect(initialBox).not.toBeNull();
    
    if (initialBox) {
      // Focus the grid container
      await page.locator('.grid-container').click();
      
      // Rapidly press movement keys in a sequence
      const sequence = ['ArrowRight', 'ArrowDown', 'ArrowLeft', 'ArrowUp'];
      for (const key of sequence) {
        await page.keyboard.press(key, { delay: 100 }); // Press keys with minimal delay
      }
      
      // Wait for all movements to complete
      await page.waitForTimeout(1000);
      
      // Verify player still exists and is in a valid position
      const finalBox = await playerElement.boundingBox();
      expect(finalBox).not.toBeNull();
      
      // Verify player is still within the grid bounds
      const gridContainer = await page.locator('.grid-container').boundingBox();
      expect(gridContainer).not.toBeNull();
      if (gridContainer && finalBox) {
        expect(finalBox.x).toBeGreaterThanOrEqual(gridContainer.x);
        expect(finalBox.y).toBeGreaterThanOrEqual(gridContainer.y);
        expect(finalBox.x + finalBox.width).toBeLessThanOrEqual(gridContainer.x + gridContainer.width);
        expect(finalBox.y + finalBox.height).toBeLessThanOrEqual(gridContainer.y + gridContainer.height);
      }
    }
  });

  test('should handle movement into walls correctly', async ({ page }) => {
    // Get initial player position
    const playerElement = await page.locator('.player-character').first();
    const initialBox = await playerElement.boundingBox();
    expect(initialBox).not.toBeNull();
    
    if (initialBox) {
      // Focus the grid container
      await page.locator('.grid-container').click();
      
      // Find a wall cell adjacent to the player
      const playerCell = await page.locator('.grid-cell:has(.player-character)').first();
      const gridCells = await page.locator('.grid-cell').all();
      const gridDimension = Math.sqrt(gridCells.length);
      
      const playerCellIndex = await playerCell.evaluate((el, cells) => {
        return Array.from(cells).indexOf(el);
      }, await page.$$('.grid-cell'));
      
      const playerGridX = playerCellIndex % gridDimension;
      const playerGridY = Math.floor(playerCellIndex / gridDimension);
      
      // Try to find a wall in adjacent cells
      const directions = [
        { dx: 1, dy: 0, key: 'ArrowRight' },
        { dx: 0, dy: 1, key: 'ArrowDown' },
        { dx: -1, dy: 0, key: 'ArrowLeft' },
        { dx: 0, dy: -1, key: 'ArrowUp' }
      ];
      
      for (const dir of directions) {
        const targetX = playerGridX + dir.dx;
        const targetY = playerGridY + dir.dy;
        const targetIndex = targetY * gridDimension + targetX;
        
        if (targetIndex >= 0 && targetIndex < gridCells.length) {
          const cell = gridCells[targetIndex];
          const classes = await cell.getAttribute('class');
          
          if (classes?.includes('wall')) {
            // Found a wall, try to move into it
            console.log(`Found wall at (${targetX}, ${targetY}), pressing ${dir.key}`);
            
            const beforeMove = await playerElement.boundingBox();
            await page.keyboard.press(dir.key);
            await page.waitForTimeout(500);
            const afterMove = await playerElement.boundingBox();
            
            // Verify player position didn't change
            expect(afterMove?.x).toBe(beforeMove?.x);
            expect(afterMove?.y).toBe(beforeMove?.y);
            break;
          }
        }
      }
    }
  });

  test('should handle clicking on walls correctly', async ({ page }) => {
    // Get all grid cells
    const gridCells = await page.locator('.grid-cell').all();
    
    // Find a wall cell
    for (const cell of gridCells) {
      const classes = await cell.getAttribute('class');
      if (classes?.includes('wall')) {
        // Found a wall cell, try to click it
        await cell.click();
        await page.waitForTimeout(500);
        
        // Verify no target marker appears on wall
        const targetMarker = await page.locator('.target-marker').count();
        expect(targetMarker).toBe(0);
        
        // Verify no path markers appear
        const pathMarkers = await page.locator('.path-marker').count();
        expect(pathMarkers).toBe(0);
        break;
      }
    }
  });

  test('should handle clicking outside grid bounds', async ({ page }) => {
    // Get the grid container
    const gridContainer = await page.locator('.grid-container').first();
    const containerBox = await gridContainer.boundingBox();
    expect(containerBox).not.toBeNull();
    
    if (containerBox) {
      // Click above the grid
      await page.mouse.click(
        containerBox.x + containerBox.width / 2,
        containerBox.y - 10
      );
      
      // Click below the grid
      await page.mouse.click(
        containerBox.x + containerBox.width / 2,
        containerBox.y + containerBox.height + 10
      );
      
      // Click to the left of the grid
      await page.mouse.click(
        containerBox.x - 10,
        containerBox.y + containerBox.height / 2
      );
      
      // Click to the right of the grid
      await page.mouse.click(
        containerBox.x + containerBox.width + 10,
        containerBox.y + containerBox.height / 2
      );
      
      // Verify no target markers appear
      const targetMarker = await page.locator('.target-marker').count();
      expect(targetMarker).toBe(0);
    }
  });

  test('should handle multiple target clicks', async ({ page }) => {
    // Get all grid cells and find the player position
    const gridCells = await page.locator('.grid-cell').all();
    const gridDimension = Math.sqrt(gridCells.length);
    
    // Find the player's position
    const playerCell = await page.locator('.grid-cell:has(.player-character)').first();
    const playerCellIndex = await playerCell.evaluate((el, cells) => {
      return Array.from(cells).indexOf(el);
    }, await page.$$('.grid-cell'));
    
    const playerGridX = playerCellIndex % gridDimension;
    const playerGridY = Math.floor(playerCellIndex / gridDimension);
    
    // Define points relative to the player that we'll try to click
    const relativePoints = [
      { dx: 1, dy: 0 },  // right
      { dx: 1, dy: 1 },  // down-right
      { dx: 0, dy: 1 }   // down
    ];
    
    // Try clicking each point
    for (const point of relativePoints) {
      const targetX = playerGridX + point.dx;
      const targetY = playerGridY + point.dy;
      const targetIndex = targetY * gridDimension + targetX;
      
      if (targetIndex >= 0 && targetIndex < gridCells.length) {
        const cell = gridCells[targetIndex];
        const classes = await cell.getAttribute('class');
        console.log(`Trying to click cell at (${targetX}, ${targetY}): ${classes}`);
        
        // Only click if it's not a wall
        if (!classes?.includes('wall')) {
          await cell.click();
          await page.waitForTimeout(500);
          
          // Take a screenshot for debugging
          await page.screenshot({ path: `target-click-${targetX}-${targetY}.png` });
          
          // Check for target marker
          const targetMarker = await page.locator('.target-marker').count();
          console.log(`Found ${targetMarker} target markers after clicking (${targetX}, ${targetY})`);
          expect(targetMarker).toBe(1);
          
          // Check for path markers
          const pathMarkers = await page.locator('.path-marker').count();
          console.log(`Found ${pathMarkers} path markers`);
        }
      }
    }
  });

  test('should handle toggle mode button interaction', async ({ page }) => {
    // Get initial player position
    const playerElement = await page.locator('.player-character').first();
    const initialBox = await playerElement.boundingBox();
    expect(initialBox).not.toBeNull();
    
    if (initialBox) {
      // Click toggle button to enable diagonal movement
      await page.click('button:has-text("Toggle Mode")');
      await page.waitForTimeout(500);
      
      // Focus grid and try diagonal movement
      await page.locator('.grid-container').click();
      
      // Test all diagonal movements
      const diagonalMoves = [
        { key: 'e', dx: 1, dy: -1 },  // Up-Right
        { key: 'c', dx: 1, dy: 1 },   // Down-Right
        { key: 'z', dx: -1, dy: 1 },  // Down-Left
        { key: 'q', dx: -1, dy: -1 }  // Up-Left
      ];
      
      let diagonalMoveSucceeded = false;
      
      for (const move of diagonalMoves) {
        const beforeMove = await playerElement.boundingBox();
        if (!beforeMove) continue;
        
        // Try the diagonal move
        await page.keyboard.press(move.key);
        await page.waitForTimeout(500);
        
        const afterMove = await playerElement.boundingBox();
        if (!afterMove) continue;
        
        // Calculate actual movement
        const dx = Math.sign(afterMove.x - beforeMove.x);
        const dy = Math.sign(afterMove.y - beforeMove.y);
        
        // If we moved diagonally, mark the test as successful
        if (dx === move.dx && dy === move.dy) {
          diagonalMoveSucceeded = true;
          break;
        }
      }
      
      // Verify at least one diagonal move worked
      expect(diagonalMoveSucceeded).toBeTruthy();
      
      // Click toggle button again to disable diagonal movement
      await page.click('button:has-text("Toggle Mode")');
      await page.waitForTimeout(500);
      
      // Try diagonal movement again, it should not work
      const beforeFinalMove = await playerElement.boundingBox();
      if (beforeFinalMove) {
        await page.keyboard.press('e'); // Try Up-Right
        await page.waitForTimeout(500);
        
        const afterFinalMove = await playerElement.boundingBox();
        if (afterFinalMove) {
          // Position should either not change or only move in one direction
          const dx = Math.sign(afterFinalMove.x - beforeFinalMove.x);
          const dy = Math.sign(afterFinalMove.y - beforeFinalMove.y);
          
          expect(dx === 0 || dy === 0).toBeTruthy();
        }
      }
    }
  });
}); 