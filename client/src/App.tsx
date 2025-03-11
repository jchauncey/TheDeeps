import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import CharacterSelection from './pages/CharacterSelection';
import CharacterCreation from './pages/CharacterCreation';

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<CharacterSelection />} />
        <Route path="/create-character" element={<CharacterCreation />} />
      </Routes>
    </Router>
  );
};

export default App; 