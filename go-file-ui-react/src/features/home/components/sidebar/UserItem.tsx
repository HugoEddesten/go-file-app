import { ChevronDown, ChevronRight, EllipsisVertical, Pencil, Trash2 } from "lucide-react";
import { useState } from "react";
import { ConfirmDialog } from "../../../../components/ConfirmDialog";
import { Card } from "../../../../components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../../../../components/ui/dropdown-menu";
import { VaultUserIcon } from "../VaultUserIcon";
import { getVaultUserRole, type VaultUser } from "../../../vaults/types";
import { useRemoveUserFromVault } from "../../../vaults/api/removeUserFromVault";
import { EditVaultUserModal } from "./EditVaultUserModal";

export const UserItem = ({
  permissions,
  vaultId,
}: {
  permissions: VaultUser[];
  vaultId: number;
}) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const hasSinglePermission = permissions.length === 1;
  const userId = permissions[0].id;
  const email = permissions[0].email;

  const { mutate: removeUser } = useRemoveUserFromVault({ vaultId });

  return (
    <Card className="flex flex-row items-center gap-2 font-semibold border rounded-md p-2">
      <div
        className="flex flex-1 gap-2 cursor-pointer"
        onClick={() => setIsExpanded((prev) => !prev)}
      >
        {hasSinglePermission ? (
          <>
            <VaultUserIcon role={permissions[0].role} className="w-4 h-4 shrink-0 mt-0.5" />
            <p className="text-xs">{permissions[0].email}</p>
          </>
        ) : (
          <div className="flex flex-col w-full">
            <div className="flex gap-2 self-start">
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
                  <Card
                    key={vu.path}
                    className="p-2! rounded-md text-xs gap-0 grid grid-cols-2"
                  >
                    <div>{getVaultUserRole(vu.role)}</div>
                    <div>{vu.path}</div>
                  </Card>
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      <DropdownMenu>
        <DropdownMenuTrigger className="shrink-0 p-0.5 rounded hover:bg-accent self-start">
          <EllipsisVertical className="w-4 h-4" />
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem className="gap-2" onClick={() => setEditOpen(true)}>
            <Pencil className="w-4 h-4" />
            Edit permissions
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className="text-destructive focus:text-destructive gap-2"
            onClick={() => setConfirmOpen(true)}
          >
            <Trash2 className="w-4 h-4" />
            Remove user
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <EditVaultUserModal
        open={editOpen}
        onOpenChange={setEditOpen}
        permissions={permissions}
        email={email}
        vaultId={vaultId}
      />

      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="Remove user from vault?"
        description={`${email} will lose all access to this vault. This cannot be undone.`}
        onConfirm={() => removeUser(userId)}
      />
    </Card>
  );
};
