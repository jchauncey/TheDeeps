import React from 'react';

// Add jest type
declare const jest: any;

// Mock the react-router-dom components
export const Routes = ({ children }: { children: React.ReactNode }) => 
  React.createElement(React.Fragment, null, children);

export const Route = ({ element }: { element: React.ReactNode }) => 
  React.createElement(React.Fragment, null, element);

export const MemoryRouter = ({ children }: { children: React.ReactNode }) => 
  React.createElement(React.Fragment, null, children);

export const useNavigate = jest.fn();
export const useLocation = jest.fn();
export const useParams = jest.fn();
export const Link = ({ to, children }: { to: string, children: React.ReactNode }) => 
  React.createElement('a', { href: to }, children);

// Export the mock as default and named exports
export default {
  Routes,
  Route,
  MemoryRouter,
  useNavigate,
  useLocation,
  useParams,
  Link
}; 