import { Card } from "../../../components/ui/card";
import { getVaultUserRole, VaultUserRole, type Vault } from "../types";
import { Database } from "lucide-react";
import { useVaultStore } from "../../../contexts/FileLibraryContext";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../../hooks/useAuth";
import { MaximizedSpinner } from "../../../components/ui/maximizedSpinner";

export const VaultItem = ({ vault }: { vault: Vault }) => {
  const { userId } = useAuth();
  const setVaultId = useVaultStore((state) => state.setVaultId);
  const navigate = useNavigate();

  const handleNavigation = (vaultId: number, path: string) => {
    setVaultId(vaultId);
    navigate(`/vault`, { replace: true, state: { dir: path } });
  };

  const owner = vault.users.find((u) => u.role === VaultUserRole.OWNER);
  const me = vault.users.find((u) => u.id == userId);

  if (!owner || !me) {
    return <MaximizedSpinner />;
  }

  const isOwner = me.id === owner.id;

  return (
    <Card
      className="w-full sm:w-64 p-4 cursor-pointer hover:shadow-md hover:bg-accent/50 transition-all"
      onClick={() => handleNavigation(vault.id, me.path)}
    >
      <div className="flex items-start gap-3">
        <div className="rounded-md bg-primary/10 p-2 text-primary shrink-0">
          <Database size={18} />
        </div>
        <div className="flex-1 min-w-0">
          <p className="font-semibold text-sm truncate">{vault.name}</p>
          <p className="text-xs text-muted-foreground truncate mt-0.5">{owner.email}</p>
        </div>
      </div>
      {!isOwner && (
        <div className="pt-3 border-t">
          <span className="text-xs bg-secondary text-secondary-foreground rounded px-2 py-0.5">
            {getVaultUserRole(me.role)}
          </span>
        </div>
      )}
    </Card>
  );
};
