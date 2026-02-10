import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Routes, Route, useNavigate } from 'react-router';
import Profile from '../pages/Profile';
import { useAuth } from '../context/AuthContext';
import { setOn401 } from './axiosInstance';
import { http, HttpResponse } from 'msw';
import { server } from '../test/setup';
import { vi, expect, it, describe } from 'vitest';
import React, { useEffect } from 'react';

// Mock the AuthContext hook
vi.mock('../context/AuthContext', () => ({
  useAuth: vi.fn(),
}));

describe('API-001: axiosInstance 401 error handling', () => {
  it('should call logout and redirect to login when a 401 error occurs', async () => {
    const logoutMock = vi.fn();
    (useAuth as any).mockReturnValue({
      user: { id: 1, name: 'Test User' },
      isAuthenticated: true,
      logout: logoutMock,
    });

    server.use(
      http.get('http://localhost:8080/api/users', () => {
        return new HttpResponse(null, { status: 401 });
      })
    );

    // This wrapper simulates the behavior in AuthProvider
    const AuthProviderSimulator = ({ children }: { children: React.ReactNode }) => {
      const navigate = useNavigate();
      useEffect(() => {
        setOn401(() => {
          logoutMock();
          navigate('/login');
        });
      }, [navigate]);
      return <>{children}</>;
    };

    render(
      <MemoryRouter initialEntries={['/profile']}>
        <AuthProviderSimulator>
          <Routes>
            <Route path="/profile" element={<Profile />} />
            <Route path="/login" element={<div>Login Page</div>} />
          </Routes>
        </AuthProviderSimulator>
      </MemoryRouter>
    );

    // Assert that logout is called
    await waitFor(() => {
      expect(logoutMock).toHaveBeenCalled();
    }, { timeout: 3000 });

    // Check if redirected to login
    expect(screen.getByText('Login Page')).toBeInTheDocument();
  });
});
