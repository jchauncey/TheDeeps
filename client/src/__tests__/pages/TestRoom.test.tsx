import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import TestRoom from '../../pages/TestRoom';
import { setupFetchMock, resetFetchMock } from '../mocks/fetch';

// Mock the RoomRenderer component to simplify testing
jest.mock('../../components/RoomRenderer', () => {
  return function MockRoomRenderer(props: any) {
    return (
      <div data-testid="room-renderer" data-props={JSON.stringify(props)}>
        <div>Room Renderer Mock</div>
        <div>Room Type: {props.roomType}</div>
        <div>Width: {props.width}</div>
        <div>Height: {props.height}</div>
        {props.roomWidth && <div>Room Width: {props.roomWidth}</div>}
        {props.roomHeight && <div>Room Height: {props.roomHeight}</div>}
      </div>
    );
  };
});

describe('TestRoom Page', () => {
  beforeEach(() => {
    setupFetchMock();
  });

  afterEach(() => {
    resetFetchMock();
    jest.clearAllMocks();
  });

  test('renders the page with default values', () => {
    render(<TestRoom />);
    
    // Check page title
    expect(screen.getByText(/Room Renderer Test/i)).toBeInTheDocument();
    
    // Check that the RoomRenderer is rendered with default props
    expect(screen.getByTestId('room-renderer')).toBeInTheDocument();
    expect(screen.getByText(/Room Type: entrance/i)).toBeInTheDocument();
    expect(screen.getByText(/Width: 20/i)).toBeInTheDocument();
    expect(screen.getByText(/Height: 20/i)).toBeInTheDocument();
    
    // Check that the controls are rendered
    expect(screen.getByText(/Room Controls/i)).toBeInTheDocument();
    
    // Use getByRole for form elements instead of getByText
    expect(screen.getByLabelText(/Room Type/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Floor Width/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Floor Height/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Room Width \(optional\)/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Room Height \(optional\)/i)).toBeInTheDocument();
    
    // Check that the legend is rendered
    expect(screen.getByText(/Legend/i)).toBeInTheDocument();
    expect(screen.getByText(/Wall \(#\)/i)).toBeInTheDocument();
    expect(screen.getByText(/Floor \(\.\)/i)).toBeInTheDocument();
    expect(screen.getByText(/Down Stairs/i)).toBeInTheDocument();
    expect(screen.getByText(/Up Stairs/i)).toBeInTheDocument();
    expect(screen.getByText(/Character/i)).toBeInTheDocument();
    expect(screen.getByText(/Monster/i)).toBeInTheDocument();
    expect(screen.getByText(/Item/i)).toBeInTheDocument();
  });

  test('changes room type when dropdown is changed', () => {
    render(<TestRoom />);
    
    // Get the room type select element
    const roomTypeSelect = screen.getByLabelText(/Room Type/i);
    
    // Change the room type to "boss"
    fireEvent.change(roomTypeSelect, { target: { value: 'boss' } });
    
    // Click the Apply Changes button
    fireEvent.click(screen.getByText(/Apply Changes/i));
    
    // Check that the RoomRenderer is re-rendered with the new room type
    expect(screen.getByText(/Room Type: boss/i)).toBeInTheDocument();
  });

  test('changes floor dimensions when inputs are changed', () => {
    render(<TestRoom />);
    
    // Get the width and height inputs using getByRole
    const widthInput = screen.getByRole('spinbutton', { name: /Floor Width/i });
    const heightInput = screen.getByRole('spinbutton', { name: /Floor Height/i });
    
    // Change the width and height
    fireEvent.change(widthInput, { target: { value: '30' } });
    fireEvent.change(heightInput, { target: { value: '25' } });
    
    // Click the Apply Changes button
    fireEvent.click(screen.getByText(/Apply Changes/i));
    
    // Check that the RoomRenderer is re-rendered with the new dimensions
    expect(screen.getByText(/Width: 30/i)).toBeInTheDocument();
    expect(screen.getByText(/Height: 25/i)).toBeInTheDocument();
  });

  test('changes room dimensions when inputs are changed', () => {
    render(<TestRoom />);
    
    // Get the room width and height inputs using getByRole
    const roomWidthInput = screen.getByRole('spinbutton', { name: /Room Width/i });
    const roomHeightInput = screen.getByRole('spinbutton', { name: /Room Height/i });
    
    // Change the room width and height
    fireEvent.change(roomWidthInput, { target: { value: '10' } });
    fireEvent.change(roomHeightInput, { target: { value: '12' } });
    
    // Click the Apply Changes button
    fireEvent.click(screen.getByText(/Apply Changes/i));
    
    // Check that the RoomRenderer is re-rendered with the new room dimensions
    expect(screen.getByText(/Room Width: 10/i)).toBeInTheDocument();
    expect(screen.getByText(/Room Height: 12/i)).toBeInTheDocument();
  });

  test('applies all changes together', () => {
    render(<TestRoom />);
    
    // Get all the inputs
    const roomTypeSelect = screen.getByLabelText(/Room Type/i);
    const widthInput = screen.getByRole('spinbutton', { name: /Floor Width/i });
    const heightInput = screen.getByRole('spinbutton', { name: /Floor Height/i });
    const roomWidthInput = screen.getByRole('spinbutton', { name: /Room Width/i });
    const roomHeightInput = screen.getByRole('spinbutton', { name: /Room Height/i });
    
    // Change all the values
    fireEvent.change(roomTypeSelect, { target: { value: 'treasure' } });
    fireEvent.change(widthInput, { target: { value: '40' } });
    fireEvent.change(heightInput, { target: { value: '35' } });
    fireEvent.change(roomWidthInput, { target: { value: '15' } });
    fireEvent.change(roomHeightInput, { target: { value: '18' } });
    
    // Click the Apply Changes button
    fireEvent.click(screen.getByText(/Apply Changes/i));
    
    // Check that the RoomRenderer is re-rendered with all the new values
    expect(screen.getByText(/Room Type: treasure/i)).toBeInTheDocument();
    expect(screen.getByText(/Width: 40/i)).toBeInTheDocument();
    expect(screen.getByText(/Height: 35/i)).toBeInTheDocument();
    expect(screen.getByText(/Room Width: 15/i)).toBeInTheDocument();
    expect(screen.getByText(/Room Height: 18/i)).toBeInTheDocument();
  });
}); 