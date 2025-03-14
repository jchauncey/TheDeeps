import { test, expect } from '@playwright/test';

test.describe('RoomSymbolDemo Component', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the component playground
    await page.goto('/component-playground');
    // Wait for the page to load
    await page.waitForSelector('h1:has-text("Component Playground")');
    // Select the Room Symbol Demo component
    await page.selectOption('select#component-select', 'RoomSymbolDemo');
    // Wait for the component to load
    await page.waitForSelector('h2:has-text("Room Symbol Demo")');
  });

  test('should display the component with default settings', async ({ page }) => {
    // Check that the component is visible (we don't need to check the exact text)
    await expect(page.getByRole('heading', { name: 'Room Symbol Demo' })).toBeVisible();
    
    // Check that the room type selector is visible
    await expect(page.locator('label:has-text("Room Type:")')).toBeVisible();
    
    // Find the debug mode switch by looking for the switch near the "Show Symbols" text
    const showSymbolsLabel = page.locator('label:has-text("Show Symbols:")');
    await expect(showSymbolsLabel).toBeVisible();
    
    // Get the parent form control
    const formControl = showSymbolsLabel.locator('xpath=ancestor::div[contains(@class, "chakra-form-control")]');
    
    // Find the switch input within the form control
    const switchInput = formControl.locator('input[type="checkbox"]');
    
    // Check that the Show Symbols toggle is enabled by default
    await expect(switchInput).toBeChecked();
    
    // Take a screenshot of the default view
    await page.screenshot({ 
      path: `./test-results/room-symbol-demo-default.png`,
      fullPage: false 
    });
  });

  test('should switch between different room types', async ({ page }) => {
    // Test each room type
    const roomTypes = [
      { value: 'standard', label: 'Standard' },
      { value: 'entrance', label: 'Entrance' },
      { value: 'boss', label: 'Boss' },
      { value: 'treasure', label: 'Treasure' },
      { value: 'shop', label: 'Shop' },
      { value: 'mixed', label: 'Mixed' },
    ];

    for (const roomType of roomTypes) {
      // Select the room type
      await page.selectOption('select[id="demo-type"]', roomType.value);
      
      // Check that the room heading has changed
      await expect(page.getByRole('heading', { name: `${roomType.label} Room` })).toBeVisible();
      
      // Take a screenshot of each room type
      await page.screenshot({ 
        path: `./test-results/room-symbol-demo-${roomType.value}.png`,
        fullPage: false 
      });
    }
  });

  test('should toggle symbol display', async ({ page }) => {
    // Find the debug mode switch by looking for the switch near the "Show Symbols" text
    const showSymbolsLabel = page.locator('label:has-text("Show Symbols:")');
    await expect(showSymbolsLabel).toBeVisible();
    
    // Get the parent form control
    const formControl = showSymbolsLabel.locator('xpath=ancestor::div[contains(@class, "chakra-form-control")]');
    
    // Find the switch input within the form control
    const switchInput = formControl.locator('input[type="checkbox"]');
    
    // Check that symbols are shown by default
    await expect(switchInput).toBeChecked();
    
    // Take a screenshot with symbols shown
    await page.screenshot({ 
      path: `./test-results/room-symbol-demo-symbols-on.png`,
      fullPage: false 
    });
    
    // Toggle symbols off by clicking the switch label
    await showSymbolsLabel.click();
    
    // Check that symbols are now hidden
    await expect(switchInput).not.toBeChecked();
    
    // Take a screenshot with symbols hidden
    await page.screenshot({ 
      path: `./test-results/room-symbol-demo-symbols-off.png`,
      fullPage: false 
    });
  });

  test('should show correct symbols for different entities', async ({ page }) => {
    // Find the debug mode switch by looking for the switch near the "Show Symbols" text
    const showSymbolsLabel = page.locator('label:has-text("Show Symbols:")');
    
    // Get the parent form control
    const formControl = showSymbolsLabel.locator('xpath=ancestor::div[contains(@class, "chakra-form-control")]');
    
    // Find the switch input within the form control
    const switchInput = formControl.locator('input[type="checkbox"]');
    
    // Make sure symbols are visible
    if (!await switchInput.isChecked()) {
      await showSymbolsLabel.click();
    }
    
    // Select the boss room to see different entity types
    await page.selectOption('select[id="demo-type"]', 'boss');
    
    // Check that the boss room is displayed
    await expect(page.getByRole('heading', { name: 'Boss Room' })).toBeVisible();
    
    // Take a screenshot to capture the symbols
    await page.screenshot({ 
      path: `./test-results/room-symbol-demo-boss-entities.png`,
      fullPage: false 
    });
    
    // Select the treasure room to see item symbols
    await page.selectOption('select[id="demo-type"]', 'treasure');
    
    // Check that the treasure room is displayed
    await expect(page.getByRole('heading', { name: 'Treasure Room' })).toBeVisible();
    
    // Take a screenshot to capture the symbols
    await page.screenshot({ 
      path: `./test-results/room-symbol-demo-treasure-entities.png`,
      fullPage: false 
    });
  });
}); 