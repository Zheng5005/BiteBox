import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import App from './App';
import { useAuth } from './context/AuthContext';
import { vi, expect, it, describe } from 'vitest';
import React from 'react';

// Mock the AuthContext hook
vi.mock('./context/AuthContext', () => ({
  useAuth: vi.fn(),
  AuthProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

describe('APP-001: App Routing and Authentication', () => {
  it('should redirect non-authenticated user from /profile to /login', () => {
    (useAuth as any).mockReturnValue({
      user: null,
      isAuthenticated: false,
    });

    render(
      <MemoryRouter initialEntries={['/profile']}>
        <App />
      </MemoryRouter>
    );

    // Check if Login page is visible
    expect(screen.getAllByText('Log In').length).toBeGreaterThan(0);
    // Ensure Profile page content is NOT visible
    expect(screen.queryByText('My Recipes')).not.toBeInTheDocument();
  });
});
