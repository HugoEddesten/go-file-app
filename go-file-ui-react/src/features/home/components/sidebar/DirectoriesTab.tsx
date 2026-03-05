import { Folder } from "lucide-react";
import type { SetStateAction } from "react";
import { Card } from "../../../../components/ui/card";
import { Button } from "../../../../components/ui/button";
import { VaultUserIcon } from "../VaultUserIcon";
import type { VaultUser } from "../../../vaults/types";

export const DirectoriesTab = ({
  permissions,
  setCurrentDir,
}: {
  permissions: VaultUser[];
  setCurrentDir: (value: SetStateAction<string>) => void;
}) => {
  return (
    <Card className="p-2 h-full shadow-inset-md gap-2">
      {permissions.map((vu) => (
        <Button
          key={vu.path}
          onClick={() => setCurrentDir(vu.path)}
          variant="ghost"
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
  );
};
