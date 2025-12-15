import { Navigate, Outlet } from "react-router-dom"
import { useAuthQuery } from "../api/useAuth"
import { MaximizedSpinner } from "../../../components/ui/maximizedSpinner"

export const PublicRoute = () => {
  const { isLoading, isSuccess } = useAuthQuery()

  if (isLoading) return <MaximizedSpinner />

  if (isSuccess) {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}