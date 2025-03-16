import { test, expect } from '@playwright/test';

test.describe('MovementDemo Component', () => {
  test('should load the page', async ({ page }) => {
    // Navigate directly to the component playground with MovementDemo selected
    await page.goto('http://localhost:3000/playground?component=MovementDemo');
    
    // Wait for the page to load
    await page.waitForTimeout(5000);
    
    // Take a screenshot for debugging
    await page.screenshot({ path: 'movement-demo-page.png' });
    
    // Log debugging information
    console.log('Page title:', await page.title());
    console.log('URL:', page.url());
    
    // Just verify we're on the right page
    expect(page.url()).toContain('component=MovementDemo');
  });
  
  test('should move player one tile to the right when right arrow key is pressed', async ({ page }) => {
    // Navigate to the component playground
    await page.goto('http://localhost:3000/playground?component=MovementDemo');
    
    // Wait for the page to load and the grid to be rendered
    await page.waitForTimeout(5000);
    
    // Find the player character element
    const playerElement = await page.locator('.player-character').first();
    
    // Get the initial position of the player
    const initialBoundingBox = await playerElement.boundingBox();
    
    // Make sure we found the player
    expect(initialBoundingBox).not.toBeNull();
    
    if (initialBoundingBox) {
      // Store the initial position
      const initialX = initialBoundingBox.x;
      const initialY = initialBoundingBox.y;
      
      // Focus on the grid container to ensure it receives keyboard events
      await page.locator('.grid-container').click();
      
      // Press the right arrow key to move the player
      await page.keyboard.press('ArrowRight');
      
      // Wait a moment for the movement to complete
      await page.waitForTimeout(500);
      
      // Get the new position of the player
      const playerAfterMove = await page.locator('.player-character').first();
      const newBoundingBox = await playerAfterMove.boundingBox();
      
      // Make sure we still have the player
      expect(newBoundingBox).not.toBeNull();
      
      if (newBoundingBox) {
        // Check if the player moved to the right (x position increased)
        // The cell width is 30px, so we expect the x position to increase by approximately that amount
        // We use a range check to account for any minor positioning differences
        expect(newBoundingBox.x).toBeGreaterThan(initialX + 20);
        
        // The y position should remain approximately the same
        expect(newBoundingBox.y).toBeCloseTo(initialY, 0);
        
        // Take a screenshot showing the movement
        await page.screenshot({ path: 'player-moved-right.png' });
      }
    }
  });
}); 