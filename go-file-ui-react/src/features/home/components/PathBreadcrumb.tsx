import { useState } from "react";
import { House } from "lucide-react";
import { cn } from "../../../lib/utils";
import { useDragPayload } from "../../../contexts/DragContext";
import type { FileData } from "../Home";
import { api } from "../../../lib/api";
import { queryClient } from "../../../lib/queryClient";
import { DragProvider } from "../../../contexts/DragProvider";

type PathSegmentProps = {
  children: React.ReactNode;
  isLast: boolean;
  isDroppable: boolean;
  onDrop: () => void;
  onClick: () => void;
};

type Segment = {
  label: string | null;
  path: string;
};

const PathSegment = ({
  children,
  isLast,
  isDroppable,
  onDrop,
  onClick,
}: PathSegmentProps) => {
  const [isOver, setIsOver] = useState(false);

  return (
    <DragProvider
      onDrop={onDrop}
      onEnter={() => {
        setIsOver(true);
      }}
      onLeave={() => {
        setIsOver(false);
      }}
      className="h-9 flex items-center"
    >
      <span
        onClick={onClick}
        className={cn(
          "rounded px-1 py-0.5 transition-colors select-none",
          isDroppable && "cursor-copy",
          isDroppable && isOver && "bg-primary text-primary-foreground",
          isLast && "text-muted-foreground",
          "hover:bg-primary hover:text-primary-foreground hover:cursor-pointer",
        )}
      >
        {children}
      </span>
    </DragProvider>
  );
};

export const PathBreadcrumb = ({
  currentDir,
  vaultId,
  onSegmentClick,
}: {
  currentDir: string;
  vaultId: number;
  onSegmentClick: (segment: Segment) => void;
}) => {
  const dragPayload = useDragPayload<FileData>();
  const parts = currentDir.split("/").filter(Boolean);

  const segments: Segment[] = [
    { label: null, path: "/" },
    ...parts.map((part, i) => ({
      label: part,
      path: "/" + parts.slice(0, i + 1).join("/"),
    })),
  ];

  const handleMove = async (destPath: string) => {
    if (!dragPayload) return;
    await api.put(`/files/${vaultId}/move/${dragPayload.Key}`, {
      destinationKey: destPath,
    });
    queryClient.invalidateQueries({ queryKey: ["files"] });
  };

  return (
    <div className="flex h-9 w-full rounded-md border border-input bg-background px-3 text-sm items-center gap-0.5">
      {segments.map((seg, i) => {
        const isLast = i === segments.length - 1;
        const isDroppable = !isLast && !!dragPayload;
        return (
          <div key={seg.path} className="inline-flex items-center gap-0.5">
            {i > 0 && (
              <span className="text-muted-foreground select-none">/</span>
            )}
            <PathSegment
              isLast={isLast}
              isDroppable={isDroppable}
              onDrop={() => handleMove(seg.path)}
              onClick={() => onSegmentClick(seg)}
            >
              {seg.label === null ? (
                <House className="w-3.5 h-3.5" />
              ) : (
                seg.label
              )}
            </PathSegment>
          </div>
        );
      })}
    </div>
  );
};
