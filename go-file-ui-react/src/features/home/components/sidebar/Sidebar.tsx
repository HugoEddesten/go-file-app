import type { SetStateAction } from "react";
import { FilePreview } from "./FilePreview";
import { VaultInfoPanel } from "./VaultInfoPanel";

export const Sidebar = ({
  selectedFile,
  vaultId,
  setCurrentDir,
}: {
  selectedFile: string | undefined;
  vaultId: number;
  setCurrentDir: (value: SetStateAction<string>) => void;
}) => {
  if (selectedFile) {
    return <FilePreview fileKey={selectedFile} vaultId={vaultId} />;
  }

  return <VaultInfoPanel vaultId={vaultId} setCurrentDir={setCurrentDir} />;
};
