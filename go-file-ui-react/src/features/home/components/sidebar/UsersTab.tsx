import { Card } from "../../../../components/ui/card";
import type { Vault } from "../../../vaults/types";
import { UserItem } from "./UserItem";

export const UsersTab = ({ vault }: { vault: Vault }) => {
  const users = vault.users.reduce<{ id: number; email: string }[]>(
    (prev, current) => {
      if (!prev.some((u) => u.id === current.id)) {
        prev.push({ id: current.id, email: current.email });
      }
      return prev;
    },
    []
  );

  return (
    <Card className="p-2 h-full shadow-inset-md gap-2">
      {users.map((u) => {
        const permissions = vault.users.filter((vu) => vu.id === u.id);
        return (
          <UserItem key={u.id} permissions={permissions} vaultId={vault.id} />
        );
      })}
    </Card>
  );
};
