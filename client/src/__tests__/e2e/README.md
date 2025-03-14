# Browser Tests for TheDeeps

This directory contains browser-level tests for TheDeeps client application using Playwright.

## Overview

These tests verify that the UI components work correctly in a real browser environment. They focus on:

1. Component Playground functionality
2. RoomRenderer component
3. RoomSymbolDemo component

## Running the Tests

### Prerequisites

- Node.js and npm installed
- TheDeeps client dependencies installed (`npm install` in the client directory)
- Playwright browsers installed (`npx playwright install`)

### Commands

Run all browser tests:

```bash
npm run test:e2e
```

Run tests with UI mode (interactive):

```bash
npm run test:e2e:ui
```

Run a specific test file:

```bash
npx playwright test src/__tests__/e2e/componentPlayground.spec.ts
```

### Test Results

- Test results will be displayed in the console
- Screenshots will be saved in the `test-results` directory
- HTML report will be generated in the `playwright-report` directory

## Test Structure

Each test file follows this structure:

1. `beforeEach` hook to navigate to the component
2. Individual test cases for specific functionality
3. Screenshots captured for visual verification

## Adding New Tests

To add a new test:

1. Create a new file in this directory with the `.spec.ts` extension
2. Import the Playwright test utilities
3. Define your test cases using the `test` function
4. Use page interactions and assertions to verify functionality

Example:

```typescript
import { test, expect } from '@playwright/test';

test.describe('My Component', () => {
  test('should do something', async ({ page }) => {
    await page.goto('/component-playground');
    await expect(page.locator('selector')).toBeVisible();
  });
});
```

## Debugging Tests

- Use `test.only` to run a single test
- Use `page.pause()` to pause execution and inspect the page
- Use the UI mode (`npm run test:e2e:ui`) for interactive debugging
- Check screenshots in the `test-results` directory