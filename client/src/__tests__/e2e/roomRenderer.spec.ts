import { test, expect } from '@playwright/test';

test.describe('RoomRenderer Component', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the component playground
    await page.goto('/component-playground');
    // Wait for the page to load
    await page.waitForSelector('h1:has-text("Component Playground")');
    // Select the Room Renderer component
    await page.selectOption('select#component-select', 'RoomRenderer');
    // Wait for the component to load
    await page.waitForSelector('h2:has-text("Room Renderer")');
  });

  test('should render different room types', async ({ page }) => {
    // Test each room type
    const roomTypes = [
      { value: 'entrance', label: 'Entrance' },
      { value: 'standard', label: 'Standard' },
      { value: 'treasure', label: 'Treasure' },
      { value: 'boss', label: 'Boss' },
      { value: 'safe', label: 'Safe' },
      { value: 'shop', label: 'Shop' },
    ];

    for (const roomType of roomTypes) {
      // Select the room type
      await page.selectOption('select#room-type', roomType.value);
      
      // Wait for the room to load
      await page.waitForTimeout(1000);
      
      // Check that the room title is correct
      await expect(page.getByText(`Test Room: ${roomType.label}`, { exact: true })).toBeVisible();
      
      // Take a screenshot of the room
      await page.screenshot({ 
        path: `./test-results/room-${roomType.value}.png`,
        fullPage: false 
      });
    }
  });

  test('should adjust room dimensions', async ({ page }) => {
    // Set room dimensions
    await page.fill('input#room-width', '30');
    await page.fill('input#room-height', '25');
    await page.fill('input#room-dim-width', '15');
    await page.fill('input#room-dim-height', '12');
    
    // Click the refresh button to apply changes
    await page.getByRole('button', { name: 'Refresh' }).click();
    
    // Wait for the room to load
    await page.waitForTimeout(1000);
    
    // Check that the room has loaded with new dimensions
    // This is a visual test, so we'll take a screenshot
    await page.screenshot({ 
      path: `./test-results/room-custom-dimensions.png`,
      fullPage: false 
    });
  });

  test('should toggle debug mode (symbols)', async ({ page }) => {
    // Find the debug mode switch by looking for the switch near the "Show Symbols" text
    const showSymbolsLabel = page.getByText('Show Symbols:', { exact: true });
    await expect(showSymbolsLabel).toBeVisible();
    
    // Get the parent form control
    const formControl = showSymbolsLabel.locator('xpath=ancestor::div[contains(@class, "chakra-form-control")]');
    
    // Find the switch input within the form control
    const switchInput = formControl.locator('input[type="checkbox"]');
    
    // Check that debug mode is enabled by default
    await expect(switchInput).toBeChecked();
    
    // Take a screenshot with debug mode on
    await page.screenshot({ 
      path: `./test-results/room-debug-on.png`,
      fullPage: false 
    });
    
    // Toggle debug mode off by clicking the switch label
    await showSymbolsLabel.click();
    
    // Check that debug mode is now disabled
    await expect(switchInput).not.toBeChecked();
    
    // Take a screenshot with debug mode off
    await page.screenshot({ 
      path: `./test-results/room-debug-off.png`,
      fullPage: false 
    });
  });

  test('should refresh room data', async ({ page }) => {
    // Click the refresh button
    await page.getByRole('button', { name: 'Refresh' }).click();
    
    // Wait for the room to load
    await page.waitForTimeout(1000);
    
    // Check that the room has loaded
    await expect(page.getByText(/Test Room:/, { exact: false })).toBeVisible();
  });
}); 