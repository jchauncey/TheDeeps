# The Deeps - Client

This is the client application for The Deeps, a roguelike dungeon crawler game.

## Getting Started

### Prerequisites

- Node.js (v14 or higher)
- npm (v6 or higher)

### Installation

1. Navigate to the client directory:
   ```
   cd client
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Start the development server:
   ```
   npm start
   ```
   
   Or use the run script:
   ```
   ./run.sh
   ```

The application will be available at http://localhost:3000.

## Troubleshooting

### ENOENT: no such file or directory, uv_cwd

If you encounter this error when running `npm run dev` or `npm start`, try one of these alternatives:

1. Use the run script which tries multiple approaches:
   ```
   ./run.sh
   ```

2. Use npx directly:
   ```
   npx react-scripts start
   ```

3. Use node to run the start.js file:
   ```
   node start.js
   ```

4. Use the simple static server:
   ```
   node server.js
   ```
   This will serve a simple test page at http://localhost:3000/

5. Try specifying the current directory with --prefix:
   ```
   npm --prefix $(pwd) start
   ```

6. Try downgrading Node.js to a more stable version (e.g., v16 or v18)

### Missing @chakra-ui/icons

If you encounter errors about missing `@chakra-ui/icons`, install it separately:
```
npm install @chakra-ui/icons
```

### Dependency Conflicts

If you encounter dependency conflicts during installation, you can use:
```
npm install --legacy-peer-deps
```

## Features

### Character Selection

- View all your characters in a grid layout
- See character details including name, class, level, HP, mana, and XP
- Delete characters you no longer want
- Create new characters (up to a maximum of 10)

### Character Creation

- Choose from 12 different character classes
- Each class has unique attributes and abilities
- Simple creation process with class descriptions

## Project Structure

- `src/components`: Reusable UI components
- `src/pages`: Page components for different routes
- `src/services`: API services for communicating with the server
- `src/types`: TypeScript type definitions

## Development

The client is built with:

- React
- TypeScript
- Chakra UI for styling
- React Router for navigation
- Axios for API requests

## Proxy Configuration

The client is configured to proxy API requests to the server running on http://localhost:8080. 