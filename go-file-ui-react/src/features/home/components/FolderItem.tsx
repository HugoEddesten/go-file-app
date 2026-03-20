import { Delete, DownloadCloud, Folder, Share } from "lucide-react";
import type { FileData } from "../Home";
import { useDragPayload } from "../../../contexts/DragContext";
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
import { useVaults } from "../../vaults/api/getVaults";
import { MaximizedSpinner } from "../../../components/ui/maximizedSpinner";
import { FolderUser } from "../../../components/ui/FolderUser";
import { useRenameFile } from "../api/renameFile";
import { useDeleteFile } from "../api/deleteFile";
import { Input } from "../../../components/ui/input";
import { useVault } from "../api/getVault";

export const FolderItem = ({
  file,
  onDoubleClick,
  vaultId,
}: {
  file: FileData;
  onDoubleClick: () => void;
  vaultId: number;
}) => {
  const [isOver, setIsOver] = useState(false);
  const [isRenaming, setIsRenaming] = useState(false);
  const [shareModalOpen, setShareModalOpen] = useState(false);

  const dragPayload = useDragPayload<FileData>();
  const { data: vault } = useVault(vaultId);

  const { mutateAsync: renameFileAsync } = useRenameFile({
    path: file.Key,
    vaultId,
  });
  const { mutateAsync: deleteFileAsync } = useDeleteFile({
    path: file.Key,
    vaultId,
  });

  const handleDownload = () => {
    window.location.href = `${
      import.meta.env.VITE_API_URL
    }files/${vaultId}/download${file.Key}?action=download`;
  };

  const handleDrop = async (e: DragEvent) => {
    setIsOver(false);

    if (dragPayload) {
      await api.put(`/files/${vaultId}/move/${dragPayload.Key}`, {
        destinationKey: file.Key,
      });
      queryClient.invalidateQueries({ queryKey: ["files"] });
      return;
    }

    const droppedFiles = Array.from(e.dataTransfer?.files ?? []);
    if (!droppedFiles.length) return;

    const formData = new FormData();
    formData.append("file", droppedFiles[0]);

    await api.post(`/files/${vaultId}/upload/${file.Key}`, formData);
    queryClient.invalidateQueries({ queryKey: ["files", file.Key] });
  };

  if (!vault) {
    return <MaximizedSpinner />;
  }

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
              {vault?.users.some((vu) => vu.path === file.Key) ? (
                <FolderUser />
              ) : (
                <Folder />
              )}
            </div>
            <div
              onDoubleClick={(e) => {
                setIsRenaming(true);
                e.stopPropagation();
                e.preventDefault();
              }}
            >
              {isRenaming ? (
                <Input
                  className="p-0 h-6 focus:ring-0!"
                  defaultValue={file.Name}
                  autoFocus
                  onFocus={(e) => e.target.select()}
                  onBlur={(e) => {
                    setIsRenaming(false);
                    if (e.target.value !== file.Name) {
                      renameFileAsync(e.target.value);
                    }
                  }}
                />
              ) : (
                <p
                  className={cn(
                    "text-center text-sm wrap-anywhere overflow-hidden truncate hover:overflow-visible hover:whitespace-break-spaces",
                  )}
                >
                  {file.Name}
                </p>
              )}
            </div>
          </ContextMenuTrigger>
          <ContextMenuContent className="w-52">
            <ContextMenuItem onClick={() => handleDownload()}>
              Download file
              <ContextMenuShortcut>
                <DownloadCloud className="text-primary" />
              </ContextMenuShortcut>
            </ContextMenuItem>
            <ContextMenuItem onClick={() => deleteFileAsync()}>
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
