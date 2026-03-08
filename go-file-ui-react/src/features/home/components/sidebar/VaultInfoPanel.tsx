import type { SetStateAction } from "react";
import { Card } from "../../../../components/ui/card";
import { MaximizedSpinner } from "../../../../components/ui/maximizedSpinner";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "../../../../components/ui/tabs";
import { useAuth } from "../../../../hooks/useAuth";
import { useVaults } from "../../../vaults/api/getVaults";
import { VaultUserRole, type Vault } from "../../../vaults/types";
import { DirectoriesTab } from "./DirectoriesTab";
import { UsersTab } from "./UsersTab";

export const VaultInfoPanel = ({
  vaultId,
  setCurrentDir,
}: {
  vaultId: number;
  setCurrentDir: (value: SetStateAction<string>) => void;
}) => {
  const { data } = useVaults({});
  const { userId } = useAuth();

  if (!data) return <MaximizedSpinner />;

  const vault = data.find((v) => v.id === vaultId) as Vault;
  const userPermissions = vault.users.filter((vu) => vu.id === userId);
  const isAdminOrOwner = userPermissions.some(
    (vu) => vu.role === VaultUserRole.ADMIN || vu.role === VaultUserRole.OWNER
  );

  return (
    <Card className="md:col-span-2 p-4 gap-2 flex">
      <Tabs className="h-full">
        <TabsList className="w-full flex justify-start">
          <TabsTrigger value="directories">Your Directories</TabsTrigger>
          {isAdminOrOwner && <TabsTrigger value="users">Users</TabsTrigger>}
        </TabsList>
        <TabsContent value="directories" className="h-full">
          <DirectoriesTab
            permissions={userPermissions}
            setCurrentDir={setCurrentDir}
          />
        </TabsContent>
        <TabsContent value="users" className="h-full">
          <UsersTab vault={vault} />
        </TabsContent>
      </Tabs>
    </Card>
  );
};
