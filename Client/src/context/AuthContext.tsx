import { createContext, useContext, useEffect, useState, type ReactNode, useCallback } from "react";
import type { User } from "../types";
import { setOn401 } from "../api/axiosInstance";
import { useNavigate } from "react-router";

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps){
  const [user, setUser] = useState<User | null >(null)
  const token = localStorage.getItem("token")
  const navigate = useNavigate();

  const logout = useCallback(() => {
    localStorage.removeItem("token")
    setUser(null)
  }, [])

  useEffect(() => {
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split(".")[1])) as User;
        setUser(payload)
      } catch (error) {
        logout()
      }
    }
  }, [token, logout])

  useEffect(() => {
    setOn401(() => {
      logout();
      navigate("/login");
    });
  }, [logout, navigate]);

  return (
    <AuthContext.Provider value={{
      user,
      logout,
      isAuthenticated: !!user
    }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }

  return context;
}
