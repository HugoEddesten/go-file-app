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

export const FileMinimized = ({
  file,
  selected = false,
  onClick,
}: {
  file: FileData;
  selected?: boolean;
  onClick?: (e: React.MouseEvent) => void;
}) => {
  const handleDownload = () => {
    window.location.href = `${import.meta.env.VITE_API_URL}files/download${
      file.Key
    }?action=download`;
  };

  return (
    <div
      className="w-18 h-18"
      onClick={(e) => {
        e.stopPropagation();
        onClick?.(e);
      }}
    >
      <ContextMenu>
        <ContextMenuTrigger className={cn(selected && "text-primary")}>
          <div className="flex justify-center w-full">
            <File />
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
            Copy file
            <ContextMenuShortcut>
              <CopyCheck className="text-primary" />
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
