import { Navigate, Outlet } from "react-router-dom"
import { useAuthQuery } from "../api/useAuth"
import { MaximizedSpinner } from "../../../components/ui/maximizedSpinner"
import { useVaultStore } from "../../../contexts/FileLibraryContext"


export const ProtectedRoute = ({ requireVaultId = false }: { requireVaultId?: boolean }) => {
  const { isLoading, isError } = useAuthQuery()
  const vaultId = useVaultStore(state => state.vaultId);
  
  const blockedByMissingVaultId = requireVaultId && !vaultId


  if (isLoading) {
    return <MaximizedSpinner />
  }

  if (isError) {
    return <Navigate to="/login" replace />
  }

  if (blockedByMissingVaultId) {
    return <Navigate to={"/"} replace />
  }

  return <Outlet context={requireVaultId && vaultId}/>
}