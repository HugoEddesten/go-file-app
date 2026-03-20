import { useEffect, useState } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { Card } from "../../components/ui/card";
import { DragProvider } from "../../contexts/DragProvider";
import { DragStateProvider } from "../../contexts/DragStateProvider";
import { cn } from "../../lib/utils";
import { useFiles } from "./api/getFiles";
import { FileItem } from "./components/FileItem";
import { api } from "../../lib/api";
import { queryClient } from "../../lib/queryClient";
import { FolderItem } from "./components/FolderItem";
import { PathBreadcrumb } from "./components/PathBreadcrumb";
import { ChevronLeft } from "lucide-react";
import { Button } from "../../components/ui/button";
import { FileLibraryMenuBar } from "./components/FileLibraryMenuBar";
import { MaximizedSpinner } from "../../components/ui/maximizedSpinner";
import { Sidebar } from "./components/sidebar/Sidebar";
import { useVault } from "./api/getVault";
import { useAuth } from "../../hooks/useAuth";

const formatBytes = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(2)} GB`;
};

export type FileData = {
  Name: string;
  Key: string;
};

export const Home = () => {
  const { vaultId: vaultIdParam } = useParams<{ vaultId: string }>();
  const vaultId = Number(vaultIdParam);
  const [searchParams, setSearchParams] = useSearchParams();
  const { userId } = useAuth();
  const navigate = useNavigate();

  const [isDragging, setIsDragging] = useState(false);
  const [selectedFile, setSelectedFile] = useState<string | undefined>();

  const { data: vault, isLoading: isVaultLoading } = useVault(vaultId);

  const startingDir = vault?.users.find((u) => u.id === userId)?.path ?? null;

  useEffect(() => {
    if (startingDir !== null && !searchParams.get("path")) {
      setSearchParams({ path: startingDir }, { replace: true });
    }
  }, [startingDir]);

  const currentDir = searchParams.get("path") ?? "";

  const { data, isLoading } = useFiles({
    vaultId,
    path: currentDir,
    enabled: !!currentDir,
  });
  const files = data ?? [];

  const handleDrop = async (e: DragEvent) => {
    setIsDragging(false);
    const files = Array.from(e.dataTransfer?.files ?? []);
    if (!files.length) return;

    const formData = new FormData();
    formData.append("file", files?.[0]);

    await api.post(`/files/${vaultId}/upload/${currentDir}`, formData);
    queryClient.invalidateQueries({ queryKey: ["files"] });
    queryClient.invalidateQueries({ queryKey: ["vault", vaultId] });
  };

  const handleGoBack = () => {
    if (!startingDir) return;
    const pathParts = currentDir.replace(startingDir, "").split("/");
    const newPath =
      startingDir + pathParts.slice(0, pathParts.length - 1).join("/");
    setSearchParams({ path: newPath });
  };

  if (isLoading || isVaultLoading) {
    return <MaximizedSpinner />;
  }

  if (!vault) {
    navigate("/");
    return null;
  }

  return (
    <div className="flex flex-col h-full">
      <div className="grid grid-cols-1 md:grid-cols-8 h-full gap-4 p-4">
        <Card
          className="p-4 md:col-span-6 relative flex-wrap overflow-hidden w-full h-full gap-4"
          onClick={() => setSelectedFile(undefined)}
        >
          <DragStateProvider>
            <DragProvider
              detectAll={true}
              onEnter={(e) => {
                if (
                  !!e?.dataTransfer?.items &&
                  e.dataTransfer.items.length > 0
                ) {
                  setIsDragging(true);
                }
              }}
              onLeave={() => {
                setIsDragging(false);
              }}
              onDrop={(e) => handleDrop(e)}
              className="flex flex-col gap-2 h-full"
            >
              <div className="text-xs flex items-center gap-2">
                <Button
                  className=" p-2! border cursor-pointer"
                  onClick={handleGoBack}
                >
                  <ChevronLeft className="w-5! h-5!" />
                </Button>
                <PathBreadcrumb
                  currentDir={currentDir}
                  vaultId={vaultId}
                  onSegmentClick={(segment) =>
                    setSearchParams({ path: segment.path })
                  }
                />
              </div>
              <div className="flex flex-col h-full gap-2">
                <FileLibraryMenuBar currentDir={currentDir} vaultId={vaultId} />

                <Card
                  className={cn(
                    "shadow-inset-md p-4 h-full rounded-md",
                    isDragging && "outline-dashed",
                  )}
                >
                  <div className="flex gap-4 flex-wrap h-full content-start">
                    {files.map((f) => (
                      <div key={f.Key} className="h-fit">
                        {f.Name.includes(".") ? (
                          <FileItem
                            key={f.Key}
                            file={f}
                            selected={selectedFile === f.Key}
                            onClick={() => setSelectedFile(f.Key)}
                            vaultId={vaultId}
                          />
                        ) : (
                          <FolderItem
                            onDoubleClick={() =>
                              setSearchParams({ path: f.Key })
                            }
                            key={f.Key}
                            file={f}
                            vaultId={vaultId}
                          />
                        )}
                      </div>
                    ))}
                  </div>
                </Card>
                <div className="flex flex-col gap-1">
                  <div className="h-1.5 w-full rounded-full bg-muted overflow-hidden">
                    <div
                      className="h-full bg-primary rounded-full transition-all"
                      style={{
                        width: `${Math.min((vault.storageUsedBytes / vault.storageLimitBytes) * 100, 100)}%`,
                      }}
                    />
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {formatBytes(vault.storageUsedBytes)} of{" "}
                    {formatBytes(vault.storageLimitBytes)} used
                  </p>
                </div>
              </div>
            </DragProvider>
          </DragStateProvider>
        </Card>
        <Sidebar
          selectedFile={selectedFile}
          vaultId={vaultId}
          setCurrentDir={(path) => setSearchParams({ path: String(path) })}
        />
      </div>
    </div>
  );
};
