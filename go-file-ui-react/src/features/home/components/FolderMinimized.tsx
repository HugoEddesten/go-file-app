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

export const FolderMinimized = ({
  file,
  onDoubleClick,
}: {
  file: FileData;
  onDoubleClick: () => void;
}) => {
  const handleDownload = () => {
    window.location.href = `${import.meta.env.VITE_API_URL}files/download${
      file.Key
    }?action=download`;
  };

  return (
    <div className="w-18 h-18">
      <ContextMenu>
        <ContextMenuTrigger onDoubleClick={onDoubleClick}>
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
    </div>
  );
};
