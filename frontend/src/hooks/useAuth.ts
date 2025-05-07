import { useState, useEffect } from 'react';

interface AuthState {
  token: string | null;
  isAuthenticated: boolean;
}

export const useAuth = () => {
  const [authState, setAuthState] = useState<AuthState>({
    token: localStorage.getItem('token'),
    isAuthenticated: !!localStorage.getItem('token'),
  });

  useEffect(() => {
    const token = localStorage.getItem('token');
    setAuthState({
      token,
      isAuthenticated: !!token,
    });
  }, []);

  const login = (token: string) => {
    localStorage.setItem('token', token);
    setAuthState({
      token,
      isAuthenticated: true,
    });
  };

  const logout = () => {
    localStorage.removeItem('token');
    setAuthState({
      token: null,
      isAuthenticated: false,
    });
  };

  return {
    token: authState.token,
    isAuthenticated: authState.isAuthenticated,
    login,
    logout,
  };
}; 