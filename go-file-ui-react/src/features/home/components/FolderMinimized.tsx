import { Delete, DownloadCloud, Folder } from "lucide-react";
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

export const FolderMinimized = ({
  file,
  onDoubleClick,
}: {
  file: FileData;
  onDoubleClick: () => void;
}) => {
  const [isOver, setIsOver] = useState(false);

  const handleDownload = () => {
    window.location.href = `${import.meta.env.VITE_API_URL}files/download${
      file.Key
    }?action=download`;
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
            <div className="flex justify-center w-full">
              <Folder />
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
          </ContextMenuContent>
        </ContextMenu>
      </DragProvider>
    </div>
  );
};
