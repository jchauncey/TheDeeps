import React from 'react';
import { render, screen, waitFor, fireEvent, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CharacterCreation from '../../pages/CharacterCreation';
import { createCharacter } from '../mocks/api';
import { BrowserRouter } from 'react-router-dom';

// Add jest type
declare const jest: any;

// Mock useNavigate
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => {
  const originalModule = jest.requireActual('react-router-dom');
  return {
    __esModule: true,
    ...originalModule,
    useNavigate: () => mockNavigate,
  };
});

// Mock useToast
const mockToast = jest.fn();
jest.mock('@chakra-ui/react', () => {
  const originalModule = jest.requireActual('@chakra-ui/react');
  return {
    __esModule: true,
    ...originalModule,
    useToast: () => mockToast,
  };
});

describe('CharacterCreation', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    jest.clearAllMocks();
  });

  it('renders the character creation form', () => {
    render(
      <BrowserRouter>
        <CharacterCreation />
      </BrowserRouter>
    );
    
    // Check if the form elements are rendered
    expect(screen.getByText('Create New Character')).toBeInTheDocument();
    
    // Check if the form has input fields
    expect(screen.getByPlaceholderText('Enter character name')).toBeInTheDocument();
    
    // Check if basic sections are rendered
    expect(screen.getByText('Basic Information')).toBeInTheDocument();
  });

  it('prevents submitting the form with invalid data', async () => {
    render(
      <BrowserRouter>
        <CharacterCreation />
      </BrowserRouter>
    );
    
    // Try to submit without a name
    const submitButton = screen.getByRole('button', { name: 'Create Character' });
    fireEvent.click(submitButton);
    
    // Verify that createCharacter was not called
    expect(createCharacter).not.toHaveBeenCalled();
  });

  it('navigates back to character selection when back button is clicked', () => {
    render(
      <BrowserRouter>
        <CharacterCreation />
      </BrowserRouter>
    );
    
    // Click the back button
    const backButton = screen.getByText('Back to Character Selection');
    fireEvent.click(backButton);
    
    // Check if navigation occurred
    expect(mockNavigate).toHaveBeenCalledWith('/');
  });

  it('allows selecting a character class', async () => {
    render(
      <BrowserRouter>
        <CharacterCreation />
      </BrowserRouter>
    );
    
    // Find the class select element
    const classSelect = screen.getByRole('combobox');
    
    // Select the barbarian class (which should be available in the options)
    fireEvent.change(classSelect, { target: { value: 'barbarian' } });
    
    // Check if barbarian is in the document
    expect(screen.getAllByText(/barbarian/i)[0]).toBeInTheDocument();
  });

  it('successfully submits the form with valid data', async () => {
    // Mock successful character creation
    createCharacter.mockResolvedValueOnce({ id: '123', name: 'Test Character' });
    
    render(
      <BrowserRouter>
        <CharacterCreation />
      </BrowserRouter>
    );
    
    // Fill in the name
    const nameInput = screen.getByPlaceholderText('Enter character name');
    fireEvent.change(nameInput, { target: { value: 'Test Character' } });
    
    // Select a class
    const classSelect = screen.getByRole('combobox');
    fireEvent.change(classSelect, { target: { value: 'barbarian' } });
    
    // Submit the form
    const submitButton = screen.getByRole('button', { name: 'Create Character' });
    await act(async () => {
      fireEvent.click(submitButton);
    });
    
    // Verify that createCharacter was called with the correct data
    expect(createCharacter).toHaveBeenCalledWith(expect.objectContaining({
      name: 'Test Character',
      class: expect.any(String)
    }));
    
    // Verify navigation occurred
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });
}); 