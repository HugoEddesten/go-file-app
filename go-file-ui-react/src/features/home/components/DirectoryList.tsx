import { Edit2, Folder } from "lucide-react";
import { Card } from "../../../components/ui/card";
import { MaximizedSpinner } from "../../../components/ui/maximizedSpinner";
import { useAuth } from "../../../hooks/useAuth";
import { useVaults } from "../../vaults/api/getVaults";
import type { Vault } from "../../vaults/types";
import { Button } from "../../../components/ui/button";
import type { SetStateAction } from "react";
import { VaultUserIcon } from "./VaultUserIcon";

export const DirectoryList = ({
  vaultId,
  setCurrentDir,
}: {
  vaultId: number;
  setCurrentDir: (value: SetStateAction<string>) => void;
}) => {
  const { data } = useVaults({});
  const { userId } = useAuth();

  if (!data) {
    return <MaximizedSpinner />;
  }
  const vault = data.find((v) => v.id === vaultId) as Vault;
  const usersVaultPermissions = vault.users.filter((vu) => vu.id === userId);

  return (
    <Card className="md:col-span-2 p-2 pt-2 gap-2 flex">
      <h5>Your Directories</h5>
      <Card className="p-2 h-full shadow-inset-md gap-2">
        {usersVaultPermissions.map((vu) => (
          <Button
            onClick={() => setCurrentDir(vu.path)}
            variant={"ghost"}
            className="flex items-center justify-between gap-2"
          >
            <div className="flex gap-2 font-semibold">
              <Folder className="w-4 h-4" />
              <p className="text-xs">{vu.path}</p>
            </div>

            <VaultUserIcon role={vu.role} className="w-4 h-4" />
          </Button>
        ))}
      </Card>
    </Card>
  );
};
