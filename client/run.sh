#!/bin/bash
# Set Node options for OpenSSL compatibility with newer Node versions
export NODE_OPTIONS="--openssl-legacy-provider --no-deprecation"

# Start the React app
echo "Starting the React app..."
npm start 