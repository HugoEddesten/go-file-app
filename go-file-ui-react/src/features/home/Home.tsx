import { useMemo, useState } from "react";
import { Card } from "../../components/ui/card";
import { Spinner } from "../../components/ui/spinner";
import { DragContext } from "../../contexts/DragContext";
import { DragProvider } from "../../contexts/DragProvider";
import { cn } from "../../lib/utils";
import { useFiles } from "./api/getFiles";
import { FileMinimized } from "./components/FileMinimized";
import { api } from "../../lib/api";
import { queryClient } from "../../lib/queryClient";
import { Input } from "../../components/ui/input";
import { useAuth } from "../../hooks/useAuth";
import { FolderMinimized } from "./components/FolderMinimized";
import { ChevronLeft } from "lucide-react";
import { Button } from "../../components/ui/button";

export type FileData = {
  Name: string;
  Key: string;
};

export const Home = () => {
  const [isDragging, setIsDragging] = useState(false);
  const [currentDir, setCurrentDir] = useState("/");
  const { data, isLoading } = useFiles({ path: currentDir });
  console.log(data)
  const handleDrop = async (e: DragEvent) => {
    setIsDragging(false);
    const files = Array.from(e.dataTransfer?.files ?? []);
    if (!files.length) return;

    const formData = new FormData();

    formData.append("file", files?.[0]);

    await api.post(`/files/upload/`, formData);
    queryClient.invalidateQueries({ queryKey: ["files"] });
  };

  const handleGoBack = () => {
    const pathParts = currentDir.slice(1).split("/");
    const newPath = "/" + pathParts.slice(0, pathParts.length - 1).join("/");
    setCurrentDir(newPath);
  };

  const files = data ?? [];

  if (isLoading) {
    return <Spinner />;
  }

  return (
    <div className="flex flex-col h-full">
      <h1 className="col-span-8">My files</h1>
      <div className="grid grid-cols-8 h-full gap-4 p-4">
        <Card className="p-6 col-span-6 relative flex-wrap overflow-hidden w-full h-full gap-4">
          <DragContext value={{ isOver: false }}>
            <DragProvider
              detectAll={true}
              onEnter={() => setIsDragging(true)}
              onLeave={() => setIsDragging(false)}
              onDrop={(e) => handleDrop(e)}
              className="flex flex-col gap-2"
            >
              <div className="text-xs flex items-center gap-2">
                <Button className="p-0! cursor-pointer" onClick={handleGoBack}>
                  <ChevronLeft className="w-6! h-6!" />
                </Button>
                <Input readOnly value={currentDir} />
              </div>
              <div>
                <div
                  className={cn(
                    "hidden absolute top-0 right-0 bottom-0 rounded-[inherit] left-0 bg-accent justify-center items-center font-bold border-2 border-dashed opacity-0 transition-opacity transition-2",
                    isDragging && "flex opacity-90"
                  )}
                  onDrop={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                  }}
                >
                  Drop here
                </div>
                <div className="flex gap-4 flex-wrap">
                  {files.map((f) => (
                    <div key={f.Key}>
                      {f.Name.includes(".") ? (
                        <FileMinimized key={f.Key} file={f} />
                      ) : (
                        <FolderMinimized
                          onDoubleClick={() => setCurrentDir(f.Key)}
                          key={f.Key}
                          file={f}
                        />
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </DragProvider>
          </DragContext>
        </Card>
        <Card className="col-span-2"></Card>
      </div>
    </div>
  );
};
