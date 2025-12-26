import { Card } from "../../../components/ui/card";
import { getVaultUserRole, VaultUserRole, type Vault } from "../types";
import { ExternalLink } from "lucide-react";
import { Button } from "../../../components/ui/button";
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

  return (
    <Card className="w-fit p-4">
      <div className="flex justify-between items-center">
        <span>{vault.name}</span>
        <Button onClick={() => handleNavigation(vault.id, me.path)}>
          <ExternalLink />
        </Button>
      </div>
      <div>
        <p className="text-xs">Owner: {owner.email}</p>
        {me.id !== owner.id && (
          <p className="text-xs">Access: {getVaultUserRole(me.role)}</p>
        )}
      </div>
    </Card>
  );
};
