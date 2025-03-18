# Testing Strategy

This directory contains tests for the client components of The Deeps. The tests are organized into different categories:

## Test Types

### Unit Tests (`/components`)
- Located in the `/components` directory
- Test individual components in isolation
- Use mocked API responses and don't require the server to be running
- Run with Jest and React Testing Library
- Example: `RoomRenderer.test.tsx` tests the RoomRenderer component with mocked API responses

### End-to-End Tests (`/e2e`)
- Located in the `/e2e` directory
- Test complete user flows and component integration
- **Require the server to be running to work properly**
- Run with Playwright
- Example: `roomRenderer.spec.ts` tests the RoomRenderer component in the context of the Component Playground

## Testing RoomRenderer

The RoomRenderer component is special because it directly interacts with the server's `/test/room` endpoint to generate room data. There are two approaches to testing it:

1. **Unit Tests** (`/components/RoomRenderer.test.tsx`):
   - Use mocked responses from the `/test/room` endpoint
   - Test rendering, error handling, and component behavior in isolation
   - Don't test actual API integration

2. **E2E Tests** (`/e2e/roomRenderer.spec.ts`):
   - Test the RoomRenderer inside the Component Playground
   - Require the actual server to be running
   - Test real API integration and response handling
   - Verify visual appearance and user interactions

## Running Tests

- Unit tests: `npm test`
- E2E tests: 
  1. Start the server: `cd server && go run main.go`
  2. Run tests: `npm run test:e2e`

## Test File Naming Conventions

- Unit tests: `ComponentName.test.tsx`
- E2E tests: `componentName.spec.ts`
