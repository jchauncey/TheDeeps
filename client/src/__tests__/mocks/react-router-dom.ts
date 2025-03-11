import React from 'react';

// Add jest type
declare const jest: any;

const mockNavigate = jest.fn();

module.exports = {
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  BrowserRouter: ({ children }: { children: React.ReactNode }) => React.createElement(React.Fragment, null, children),
  MemoryRouter: ({ children }: { children: React.ReactNode }) => React.createElement(React.Fragment, null, children),
}; 