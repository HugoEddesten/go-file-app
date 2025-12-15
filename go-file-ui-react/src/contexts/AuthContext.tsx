import { createContext } from "react"

type AuthContextType = {
  isAuthenticated: boolean
  isLoading: boolean
}

export const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  isLoading: true
})