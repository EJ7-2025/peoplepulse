'use client';

import { createContext, useState, useContext, useEffect, ReactNode } from 'react';

interface AuthContextType { token: string | null; logout: () => void; }
const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    // SIMULAÇÃO DE LOGIN: Gera um token falso para o usuário 'Edgard' (ID 1)
    // Este token não é válido, mas permite que nossa lógica de frontend funcione.
    const fakeToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInJvbGUiOiJkaXJldG9yaWEifQ.fake_signature";
    setToken(fakeToken);
    localStorage.setItem('authToken', fakeToken);
  }, []);

  const logout = () => {
    setToken(null);
    localStorage.removeItem('authToken');
  };

  // A função login fica vazia, pois não é mais necessária por enquanto
  const login = (token: string) => {};

  return (
    <AuthContext.Provider value={{ token, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) { throw new Error('useAuth must be used within an AuthProvider'); }
  return context;
}