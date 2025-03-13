import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import CharacterSelection from './pages/CharacterSelection';
import CharacterCreation from './pages/CharacterCreation';
import DungeonSelection from './pages/DungeonSelection';
import Game from './pages/Game';
import TestRoom from './pages/TestRoom';
import ComponentPlayground from './pages/ComponentPlayground';
import SplashScreen from './components/SplashScreen';

const App: React.FC = () => {
  // Set isInitializing to false for immediate rendering during development
  const [isInitializing, setIsInitializing] = useState(false);
  
  // Simulate initialization process
  useEffect(() => {
    // You can add actual initialization logic here
    // For example, loading game assets, connecting to server, etc.
    
    // For now, we're just using the splash screen with its minimum display time
    console.log('Game initializing...');
  }, []);

  const handleInitializationComplete = () => {
    console.log('Initialization complete!');
    setIsInitializing(false);
  };

  return (
    <>
      {isInitializing ? (
        <SplashScreen onInitializationComplete={handleInitializationComplete} />
      ) : (
        <Router>
          <Routes>
            <Route path="/" element={<CharacterSelection />} />
            <Route path="/create-character" element={<CharacterCreation />} />
            <Route path="/dungeon-selection" element={<DungeonSelection />} />
            <Route path="/game" element={<Game />} />
            <Route path="/test-room" element={<TestRoom />} />
            <Route path="/component-playground" element={<ComponentPlayground />} />
          </Routes>
        </Router>
      )}
    </>
  );
};

export default App; 