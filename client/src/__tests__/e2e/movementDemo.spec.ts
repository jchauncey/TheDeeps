import { test, expect } from '@playwright/test';

test.describe('MovementDemo Component', () => {
  test('should render the MovementDemo component correctly', async ({ page }) => {
    // Navigate directly to the component playground with MovementDemo selected
    await page.goto('http://localhost:3000/playground?component=MovementDemo');
    
    // Wait for the page to load
    await page.waitForTimeout(2000);
    
    // Take a screenshot to debug
    await page.screenshot({ path: 'screenshot.png' });
    
    // Check if the component title is visible (using a more general selector)
    await expect(page.locator('h2:has-text("Movement Demo")')).toBeVisible();
    
    // Check if the description is visible (using a more general selector)
    await expect(page.locator('text="This component demonstrates character movement"')).toBeVisible();
    
    // Check if the movement controls are visible
    await expect(page.locator('text="Movement Mode:"')).toBeVisible();
    
    // Check if the grid controls are visible
    await expect(page.locator('text="Grid Size:"')).toBeVisible();
    await expect(page.locator('text="Wall Density:"')).toBeVisible();
  });
}); 