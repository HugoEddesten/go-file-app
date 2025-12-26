import { useState } from "react";
import { Card } from "../../components/ui/card";
import { DragContext } from "../../contexts/DragContext";
import { DragProvider } from "../../contexts/DragProvider";
import { cn } from "../../lib/utils";
import { useFiles } from "./api/getFiles";
import { FileMinimized } from "./components/FileMinimized";
import { api } from "../../lib/api";
import { queryClient } from "../../lib/queryClient";
import { Input } from "../../components/ui/input";
import { FolderMinimized } from "./components/FolderMinimized";
import { ChevronLeft } from "lucide-react";
import { Button } from "../../components/ui/button";
import { FileLibraryMenuBar } from "./components/FileLibraryMenuBar";
import { MaximizedSpinner } from "../../components/ui/maximizedSpinner";
import { FilePreview } from "./components/FilePreview";
import {
  useLocation,
  useOutletContext,
} from "react-router-dom";
import { DirectoryList } from "./components/DirectoryList";

export type FileData = {
  Name: string;
  Key: string;
};

export const Home = () => {
  const { dir } = useLocation().state;

  const [isDragging, setIsDragging] = useState(false);
  const [currentDir, setCurrentDir] = useState<string>(dir);
  const [selectedFile, setSelectedFile] = useState<string | undefined>();

  const vaultId = useOutletContext() as number;
  const { data, isLoading } = useFiles({ vaultId: vaultId, path: currentDir });
  console.log(data)
  const handleDrop = async (e: DragEvent) => {
    setIsDragging(false);
    const files = Array.from(e.dataTransfer?.files ?? []);
    if (!files.length) return;

    const formData = new FormData();

    formData.append("file", files?.[0]);

    await api.post(`/files/${vaultId}/upload/${currentDir}`, formData);
    queryClient.invalidateQueries({ queryKey: ["files"] });
  };

  const handleGoBack = () => {
    const pathParts = currentDir.replace(dir, "").split("/");

    const newPath = dir + pathParts.slice(0, pathParts.length - 1).join("/");
    setCurrentDir(newPath);
  };

  const files = data ?? [];
  console.log(files)
  return (
    <div className="flex flex-col h-full">
      <div className="grid grid-cols-1 md:grid-cols-8 h-full gap-4 p-4">
        <Card
          className="p-4 md:col-span-6 relative flex-wrap overflow-hidden w-full h-full gap-4"
          onClick={() => setSelectedFile(undefined)}
        >
          <DragContext value={{ isOver: false }}>
            <DragProvider
              detectAll={true}
              onEnter={() => {
                console.log("parent in");
                setIsDragging(true);
              }}
              onLeave={() => {
                console.log("parent out");
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
                <Input readOnly value={currentDir} />
              </div>
              <div className="flex flex-col h-full gap-2">
                <FileLibraryMenuBar currentDir={currentDir} vaultId={vaultId}/>

                <Card
                  className={cn(
                    "shadow-inset-md p-4 h-full rounded-md",
                    isDragging && "outline-dashed"
                  )}
                >
                  <div className="flex gap-4 flex-wrap h-full content-start">
                    {isLoading ? (
                      <MaximizedSpinner />
                    ) : (
                      files.map((f) => (
                        <div key={f.Key} className="h-fit">
                          {f.Name.includes(".") ? (
                            <FileMinimized
                              key={f.Key}
                              file={f}
                              selected={selectedFile === f.Key}
                              onClick={() => setSelectedFile(f.Key)}
                              vaultId={vaultId}
                            />
                          ) : (
                            <FolderMinimized
                              onDoubleClick={() => setCurrentDir(f.Key)}
                              key={f.Key}
                              file={f}
                              vaultId={vaultId}
                            />
                          )}
                        </div>
                      ))
                    )}
                  </div>
                </Card>
              </div>
            </DragProvider>
          </DragContext>
        </Card>
        {selectedFile ? (
          <FilePreview fileKey={selectedFile} vaultId={vaultId} />
        ) : (
          <DirectoryList vaultId={vaultId} setCurrentDir={setCurrentDir}/>
        )}
      </div>
    </div>
  );
};
