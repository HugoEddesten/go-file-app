import { Delete, DownloadCloud, Folder, Share } from "lucide-react";
import type { FileData } from "../Home";
import { cn } from "../../../lib/utils";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuShortcut,
  ContextMenuTrigger,
} from "../../../components/ui/context-menu";
import { DragProvider } from "../../../contexts/DragProvider";
import { useState } from "react";
import { api } from "../../../lib/api";
import { queryClient } from "../../../lib/queryClient";
import { ShareVaultModal } from "./ShareVaultModal";
import { all } from "axios";
import { useVaults } from "../../vaults/api/getVaults";
import { MaximizedSpinner } from "../../../components/ui/maximizedSpinner";
import { FolderUser } from "../../../components/ui/FolderUser";

export const FolderMinimized = ({
  file,
  onDoubleClick,
  vaultId,
}: {
  file: FileData;
  onDoubleClick: () => void;
  vaultId: number;
}) => {
  const [isOver, setIsOver] = useState(false);
  const [shareModalOpen, setShareModalOpen] = useState(false);

  const { data: vaults } = useVaults({});
  
  if (!vaults) {
    return <MaximizedSpinner />
  }
  
  const currentVault = vaults.find(v => v.id === vaultId)

  const handleDownload = () => {
    window.location.href = `${
      import.meta.env.VITE_API_URL
    }files/${vaultId}/download${file.Key}?action=download`;
  };

  const handleDrop = async (e: DragEvent) => {
    setIsOver(false);
    const files = Array.from(e.dataTransfer?.files ?? []);
    if (!files.length) return;

    const formData = new FormData();

    formData.append("file", files?.[0]);

    await api.post(`/files/upload/${file.Key}`, formData);
    queryClient.invalidateQueries({ queryKey: ["files", file.Key] });
  };

  return (
    <div className="w-18 h-18">
      <DragProvider
        onEnter={() => setIsOver(true)}
        onLeave={() => setIsOver(false)}
        onDrop={handleDrop}
      >
        <ContextMenu>
          <ContextMenuTrigger
            onDoubleClick={onDoubleClick}
            className={cn(isOver && "text-primary")}
          >
            <div className="flex justify-center w-full bg-card">
              {currentVault?.users.some(vu => vu.path === file.Key) ? (
                <FolderUser />
              ) : (
                <Folder />
              )}
            </div>
            <p
              className={cn(
                "text-center text-sm wrap-anywhere overflow-hidden truncate hover:overflow-visible hover:whitespace-break-spaces"
              )}
            >
              {file.Name}
            </p>
          </ContextMenuTrigger>
          <ContextMenuContent className="w-52">
            <ContextMenuItem onClick={() => handleDownload()}>
              Download file
              <ContextMenuShortcut>
                <DownloadCloud className="text-primary" />
              </ContextMenuShortcut>
            </ContextMenuItem>
            <ContextMenuItem>
              Delete file
              <ContextMenuShortcut>
                <Delete className="text-primary" />
              </ContextMenuShortcut>
            </ContextMenuItem>
            <ContextMenuItem onClick={() => setShareModalOpen(true)}>
              <div className="flex w-full justify-between">
                Share folder
                <ContextMenuShortcut>
                  <Share className="text-primary" />
                </ContextMenuShortcut>
              </div>
            </ContextMenuItem>
          </ContextMenuContent>
        </ContextMenu>
      </DragProvider>
      <ShareVaultModal
        defaultPath={file.Key}
        open={shareModalOpen}
        onOpenChange={setShareModalOpen}
      />
    </div>
  );
};
