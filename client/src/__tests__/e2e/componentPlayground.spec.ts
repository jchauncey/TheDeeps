import { test, expect } from '@playwright/test';

test.describe('Component Playground', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the component playground
    await page.goto('/component-playground');
    // Wait for the page to load
    await page.waitForSelector('h1:has-text("Component Playground")');
  });

  test('should load and display the component playground', async ({ page }) => {
    // Check that the page title is correct
    await expect(page.locator('h1')).toContainText('Component Playground');
    
    // Check that the component selector is present
    await expect(page.getByText('Select Component:')).toBeVisible();
    
    // Check that the default component is loaded (CharacterCard)
    await expect(page.getByRole('heading', { name: 'Character Cards' })).toBeVisible();
  });

  test('should switch between components', async ({ page }) => {
    // Select the Room Renderer component
    await page.selectOption('select#component-select', 'RoomRenderer');
    
    // Check that the Room Renderer component is displayed
    await expect(page.getByRole('heading', { name: 'Room Renderer' })).toBeVisible();
    
    // Select the Map Symbols component
    await page.selectOption('select#component-select', 'SymbolRenderer');
    
    // Check that the Map Symbols component is displayed
    await expect(page.getByRole('heading', { name: 'Map Symbols' })).toBeVisible();
    
    // Select the Room Symbol Demo component
    await page.selectOption('select#component-select', 'RoomSymbolDemo');
    
    // Check that the Room Symbol Demo component is displayed
    await expect(page.getByRole('heading', { name: 'Room Symbol Demo' })).toBeVisible();
  });

  test('should toggle debug mode in RoomRenderer', async ({ page }) => {
    // Select the Room Renderer component
    await page.selectOption('select#component-select', 'RoomRenderer');
    
    // Wait for the component to load
    await page.waitForSelector('h2:has-text("Room Renderer")');
    
    // Find the debug mode switch by looking for the switch near the "Show Symbols" text
    const showSymbolsLabel = page.locator('label:has-text("Show Symbols:")');
    await expect(showSymbolsLabel).toBeVisible();
    
    // Get the parent form control
    const formControl = showSymbolsLabel.locator('xpath=ancestor::div[contains(@class, "chakra-form-control")]');
    
    // Find the switch input within the form control
    const switchInput = formControl.locator('input[type="checkbox"]');
    
    // Check that debug mode is enabled by default
    await expect(switchInput).toBeChecked();
    
    // Toggle debug mode off by clicking the switch label
    await showSymbolsLabel.click();
    
    // Check that debug mode is now disabled
    await expect(switchInput).not.toBeChecked();
    
    // Toggle debug mode back on
    await showSymbolsLabel.click();
    
    // Check that debug mode is now enabled again
    await expect(switchInput).toBeChecked();
  });

  test('should change room type in RoomRenderer', async ({ page }) => {
    // Select the Room Renderer component
    await page.selectOption('select#component-select', 'RoomRenderer');
    
    // Check that the default room type is 'entrance'
    await expect(page.locator('select#room-type')).toHaveValue('entrance');
    
    // Change the room type to 'boss'
    await page.selectOption('select#room-type', 'boss');
    
    // Check that the room type has changed
    await expect(page.locator('select#room-type')).toHaveValue('boss');
    
    // Wait for the room to load
    await page.waitForTimeout(1000);
    
    // Check that the room information shows the boss room
    await expect(page.locator('text=Test Room: Boss')).toBeVisible();
  });

  test('should toggle room types in RoomSymbolDemo', async ({ page }) => {
    // Select the Room Symbol Demo component
    await page.selectOption('select#component-select', 'RoomSymbolDemo');
    
    // Wait for the component to load
    await page.waitForSelector('h2:has-text("Room Symbol Demo")');
    
    // Check that the default room type is visible
    await expect(page.locator('label:has-text("Room Type:")')).toBeVisible();
    
    // Find the debug mode switch by looking for the switch near the "Show Symbols" text
    const showSymbolsLabel = page.locator('label:has-text("Show Symbols:")');
    await expect(showSymbolsLabel).toBeVisible();
    
    // Get the parent form control
    const formControl = showSymbolsLabel.locator('xpath=ancestor::div[contains(@class, "chakra-form-control")]');
    
    // Find the switch input within the form control
    const switchInput = formControl.locator('input[type="checkbox"]');
    
    // Check that the Show Symbols toggle is on by default
    await expect(switchInput).toBeChecked();
    
    // Change the room type to 'boss'
    await page.selectOption('select[id="demo-type"]', 'boss');
    
    // Check that the room heading has changed
    await expect(page.getByRole('heading', { name: 'Boss Room' })).toBeVisible();
    
    // Toggle Show Symbols off by clicking the label
    await showSymbolsLabel.click();
    
    // Check that Show Symbols is now off
    await expect(switchInput).not.toBeChecked();
  });
}); 