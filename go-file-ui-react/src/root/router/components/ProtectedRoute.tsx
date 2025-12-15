import { Navigate, Outlet } from "react-router-dom"
import { useAuthQuery } from "../api/useAuth"
import { MaximizedSpinner } from "../../../components/ui/maximizedSpinner"


export const ProtectedRoute = () => {
  const { isLoading, isError } = useAuthQuery()
  
  if (isLoading) {
    return <MaximizedSpinner />
  }

  if (isError) {
    return <Navigate to="/login" replace />
  }

  return <Outlet />
}