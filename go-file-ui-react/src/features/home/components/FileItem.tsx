import { CopyCheck, Delete, DownloadCloud, File } from "lucide-react";
import type { FileData } from "../Home";
import { cn } from "../../../lib/utils";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuShortcut,
  ContextMenuTrigger,
} from "../../../components/ui/context-menu";
import { useState } from "react";
import { Input } from "../../../components/ui/input";
import { useRenameFile } from "../api/renameFile";
import { useDeleteFile } from "../api/deleteFile";
import { DragSource } from "../../../contexts/DragSource";

export const FileItem = ({
  file,
  selected = false,
  onClick,
  vaultId,
}: {
  file: FileData;
  selected?: boolean;
  onClick?: (e: React.MouseEvent) => void;
  vaultId: number;
}) => {
  const [isRenaming, setIsRenaming] = useState(false);
  const extension = file.Name.split(".")[1];
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
    }files/${vaultId}/download/${file.Key}?action=download`;
  };

  return (
    <DragSource<FileData>
      payload={file}
      className="w-18 h-18"
      onClick={(e) => {
        e.stopPropagation();
        onClick?.(e);
      }}
    >
      <ContextMenu>
        <ContextMenuTrigger className={cn(selected && "text-primary")}>
          <div className="flex justify-center w-full bg-none">
            <File />
          </div>
          <div
            className={cn(
              "text-center text-sm wrap-anywhere overflow-hidden truncate hover:overflow-visible hover:whitespace-break-spaces",
            )}
            onDoubleClick={() => setIsRenaming(true)}
          >
            {isRenaming ? (
              <Input
                className="p-0 h-6 focus:ring-0!"
                defaultValue={file.Name.split(".")[0]}
                autoFocus
                onFocus={(e) => e.target.select()}
                onBlur={(e) => {
                  setIsRenaming(false);
                  if (e.target.value !== file.Name) {
                    renameFileAsync(e.target.value.concat(".", extension));
                  }
                }}
              />
            ) : (
              file.Name
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
          <ContextMenuItem>
            Copy file
            <ContextMenuShortcut>
              <CopyCheck className="text-primary" />
            </ContextMenuShortcut>
          </ContextMenuItem>
          <ContextMenuItem onClick={() => deleteFileAsync()}>
            Delete file
            <ContextMenuShortcut>
              <Delete className="text-primary" />
            </ContextMenuShortcut>
          </ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
    </DragSource>
  );
};
