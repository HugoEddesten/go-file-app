import { ChevronDown, ChevronRight, Folder } from "lucide-react";
import { Card, CardContent } from "../../../components/ui/card";
import { MaximizedSpinner } from "../../../components/ui/maximizedSpinner";
import { useAuth } from "../../../hooks/useAuth";
import { useVaults } from "../../vaults/api/getVaults";
import {
  getVaultUserRole,
  VaultUserRole,
  type Vault,
  type VaultUser,
} from "../../vaults/types";
import { Button } from "../../../components/ui/button";
import { useState, type SetStateAction } from "react";
import { VaultUserIcon } from "./VaultUserIcon";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "../../../components/ui/tabs";

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
  const users = vault.users.reduce<{ id: number; email: string }[]>(
    (prev, current) => {
      if (!prev.some((vu) => vu.id === current.id)) {
        prev.push({ id: current.id, email: current.email });
      }
      return prev;
    },
    []
  );
  return (
    <Card className="md:col-span-2 p-4 gap-2 flex">
      <Tabs className="h-full">
        <TabsList className="w-full flex justify-start">
          <TabsTrigger value="directories" className="w-4!">
            Your Directories
          </TabsTrigger>
          {usersVaultPermissions.some(
            (vu) =>
              vu.role === VaultUserRole.ADMIN || vu.role === VaultUserRole.OWNER
          ) && <TabsTrigger value="users">Users</TabsTrigger>}
        </TabsList>
        <TabsContent value="directories" className="h-full">
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
        </TabsContent>
        <TabsContent value="users" className="h-full">
          <Card className="p-2 h-full shadow-inset-md gap-2">
            {users.map((u) => {
              const permissions = vault.users.filter((vu) => vu.id === u.id);

              return <UserItem permissions={permissions} />;
            })}
          </Card>
        </TabsContent>
      </Tabs>
    </Card>
  );
};

const UserItem = ({ permissions }: { permissions: VaultUser[] }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const hasSinglePermission = permissions.length === 1;

  return (
    <Card
      className="flex flex-row gap-2 font-semibold border rounded-md p-2"
      onClick={() => setIsExpanded((prev) => !prev)}
    >
      {hasSinglePermission ? (
        <>
          <VaultUserIcon role={permissions[0].role} className="w-4 h-4" />
          <p className="text-xs">{permissions[0].email}</p>
        </>
      ) : (
        <div className="flex flex-col w-full">
          <div className="flex gap-2">
            {isExpanded ? (
              <ChevronDown className="w-4 h-4" />
            ) : (
              <ChevronRight className="w-4 h-4" />
            )}
            <p className="text-xs">{permissions[0].email}</p>
          </div>

          {isExpanded && (
            <div className="flex flex-col w-full gap-2 p-2">
              <div className="grid px-2 grid-cols-2 text-xs">
                <div>role:</div>
                <div>path:</div>
              </div>
              {permissions.map((vu) => (
                <Card className="p-2! rounded-md text-xs gap-0 grid grid-cols-2">
                  <div>{getVaultUserRole(vu.role)}</div>
                  <div>{vu.path}</div>
                </Card>
              ))}
            </div>
          )}
        </div>
      )}
    </Card>
  );
};
