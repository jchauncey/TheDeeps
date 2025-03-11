const { spawn } = require('child_process');
const path = require('path');

// Get the absolute path to the react-scripts binary
const reactScriptsPath = path.resolve(
  __dirname,
  'node_modules',
  '.bin',
  process.platform === 'win32' ? 'react-scripts.cmd' : 'react-scripts'
);

// Spawn the react-scripts process
const child = spawn(reactScriptsPath, ['start'], {
  stdio: 'inherit',
  env: { ...process.env, BROWSER: 'none' }
});

// Handle process exit
child.on('close', (code) => {
  console.log(`Child process exited with code ${code}`);
  process.exit(code);
});

// Handle process errors
child.on('error', (err) => {
  console.error('Failed to start child process:', err);
  process.exit(1);
}); 