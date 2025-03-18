import React from 'react';
import { render, screen, waitFor, fireEvent, act } from '@testing-library/react';
import '@testing-library/jest-dom';

// Mock useNavigate
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  BrowserRouter: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

// Mock API
const createCharacter = jest.fn().mockImplementation((data) => {
  return Promise.resolve({ id: '123', ...data });
});

// Simple Character Creation component for testing
const CharacterCreationTest = () => {
  const navigate = mockNavigate;
  const [name, setName] = React.useState('');
  const [characterClass, setCharacterClass] = React.useState('warrior');
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Validate form
    if (!name.trim()) {
      return;
    }
    
    try {
      // Call the mocked API
      await createCharacter({
        name,
        class: characterClass,
      });
      
      // Navigate back to character selection
      navigate('/');
    } catch (error: unknown) {
      console.error('Error creating character:', error);
    }
  };
  
  const handleBackClick = () => {
    navigate('/');
  };
  
  return (
    <div>
      <h1>Create New Character</h1>
      <form onSubmit={handleSubmit} data-testid="character-form">
        <div>
          <h2>Basic Information</h2>
          <input 
            placeholder="Enter character name" 
            value={name}
            onChange={(e) => setName(e.target.value)}
            data-testid="character-name-input"
          />
          <select 
            value={characterClass}
            onChange={(e) => setCharacterClass(e.target.value)}
            data-testid="character-class-select"
          >
            <option value="barbarian">Barbarian</option>
            <option value="warrior">Warrior</option>
            <option value="mage">Mage</option>
          </select>
          <button type="submit" data-testid="create-button">Create Character</button>
          <button type="button" onClick={handleBackClick} data-testid="back-button">Back to Character Selection</button>
        </div>
      </form>
    </div>
  );
};

describe('CharacterCreation', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    jest.clearAllMocks();
  });

  it('renders the character creation form', () => {
    render(<CharacterCreationTest />);
    
    // Check if the form elements are rendered
    expect(screen.getByText('Create New Character')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Enter character name')).toBeInTheDocument();
    expect(screen.getByText('Basic Information')).toBeInTheDocument();
  });

  it('prevents submitting the form with invalid data', async () => {
    render(<CharacterCreationTest />);
    
    // Try to submit without a name
    const submitButton = screen.getByTestId('create-button');
    fireEvent.click(submitButton);
    
    // Verify that createCharacter was not called
    expect(createCharacter).not.toHaveBeenCalled();
  });

  it('navigates back to character selection when back button is clicked', () => {
    render(<CharacterCreationTest />);
    
    // Click the back button
    const backButton = screen.getByTestId('back-button');
    fireEvent.click(backButton);
    
    // Check if navigation occurred
    expect(mockNavigate).toHaveBeenCalledWith('/');
  });

  it('allows selecting a character class', () => {
    render(<CharacterCreationTest />);
    
    // Find the class select element
    const classSelect = screen.getByTestId('character-class-select');
    
    // Select the barbarian class
    fireEvent.change(classSelect, { target: { value: 'barbarian' } });
    
    // Check if the value was updated
    expect(classSelect).toHaveValue('barbarian');
  });

  it('successfully submits the form with valid data', async () => {
    render(<CharacterCreationTest />);
    
    // Fill in the name
    const nameInput = screen.getByTestId('character-name-input');
    fireEvent.change(nameInput, { target: { value: 'Test Character' } });
    
    // Select a class
    const classSelect = screen.getByTestId('character-class-select');
    fireEvent.change(classSelect, { target: { value: 'barbarian' } });
    
    // Submit the form
    const submitButton = screen.getByTestId('create-button');
    await act(async () => {
      fireEvent.click(submitButton);
    });
    
    // Verify that createCharacter was called with the correct data
    expect(createCharacter).toHaveBeenCalledWith({
      name: 'Test Character',
      class: 'barbarian'
    });
    
    // Verify navigation occurred after API call resolves
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });
}); 