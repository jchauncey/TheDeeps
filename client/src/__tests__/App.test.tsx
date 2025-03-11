import React from 'react';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import App from '../App';

// Add jest type
declare const jest: any;

// Mock the page components
jest.mock('../pages/CharacterSelection', () => {
  const MockCharacterSelection = () => <div data-testid="character-selection">Character Selection Page</div>;
  return MockCharacterSelection;
});

jest.mock('../pages/CharacterCreation', () => {
  const MockCharacterCreation = () => <div data-testid="character-creation">Character Creation Page</div>;
  return MockCharacterCreation;
});

// Mock the react-router-dom components
jest.mock('react-router-dom', () => {
  const originalModule = jest.requireActual('react-router-dom');
  return {
    ...originalModule,
    BrowserRouter: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
    Routes: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
    Route: ({ path, element }: { path: string; element: React.ReactNode }) => (
      <div data-testid={`route-${path}`}>{element}</div>
    ),
  };
});

describe('App', () => {
  it('renders CharacterSelection at the root path', () => {
    render(
      <MemoryRouter initialEntries={['/']}>
        <App />
      </MemoryRouter>
    );
    
    expect(screen.getByTestId('character-selection')).toBeInTheDocument();
  });

  it('renders CharacterCreation at the /create-character path', () => {
    render(
      <MemoryRouter initialEntries={['/create-character']}>
        <App />
      </MemoryRouter>
    );
    
    expect(screen.getByTestId('character-creation')).toBeInTheDocument();
  });
}); 