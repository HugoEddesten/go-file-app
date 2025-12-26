import { FileText, Folder } from "lucide-react";
import {
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarTrigger,
} from "../../../components/ui/menubar";
import { useCreateFile } from "../api/addFile";
import { toast } from "sonner";

export const FileLibraryMenuBar = ({ currentDir, vaultId }: { currentDir: string, vaultId: number }) => {
  const { mutateAsync } = useCreateFile({ path: currentDir, vaultId: vaultId });

  const handleNewFolderClick = async () => {
    try {
      await mutateAsync(undefined);
    } catch (err: any) {
      const message = err?.response?.data;
      toast.error(message ?? "Something went wrong");
    }
  };

  const handleNewTextfileClick = async () => {
    try {
      await mutateAsync(".txt");
    } catch (err: any) {
      const message = err?.response?.data;
      toast.error(message ?? "Something went wrong");
    }
  };

  return (
    <Menubar>
      <MenubarMenu>
        <MenubarTrigger>
          <span className="text-xs font-semibold">New</span>
        </MenubarTrigger>
        <MenubarContent className="flex flex-col gap-1">
          <MenubarItem
            className="flex justify-between hover:bg-accent"
            onClick={handleNewTextfileClick}
          >
            Textfile
            <FileText className="text-primary" />
          </MenubarItem>
          <MenubarItem
            className="flex justify-between hover:bg-accent"
            onClick={handleNewFolderClick}
          >
            Folder
            <Folder className="text-primary" />
          </MenubarItem>
        </MenubarContent>
      </MenubarMenu>
      <MenubarMenu>
        <MenubarTrigger>
          <span className="text-xs font-semibold">View</span>
        </MenubarTrigger>
        <MenubarContent className="p-0">
          <MenubarItem className="flex justify-between hover:bg-accent">
            Textfile
            <FileText />
          </MenubarItem>
          <MenubarItem className="flex justify-between hover:bg-accent">
            Folder
            <Folder />
          </MenubarItem>
        </MenubarContent>
      </MenubarMenu>
    </Menubar>
  );
};
