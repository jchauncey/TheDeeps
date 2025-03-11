#!/bin/bash
# Try multiple approaches to start the React app

echo "Attempting to start the React app using node start.js..."
node start.js

# If that fails, try npx directly
if [ $? -ne 0 ]; then
  echo "Failed to start with node start.js, trying npx..."
  npx react-scripts start
fi

# If that fails too, try with --prefix
if [ $? -ne 0 ]; then
  echo "Failed to start with npx, trying npm with --prefix..."
  npm --prefix $(pwd) start
fi

# If all else fails, use the simple server
if [ $? -ne 0 ]; then
  echo "All React app start methods failed. Starting simple static server..."
  node server.js
fi 