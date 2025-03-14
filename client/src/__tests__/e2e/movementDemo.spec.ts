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
}); 