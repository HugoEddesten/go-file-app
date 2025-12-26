import { Folder, User, Users } from "lucide-react";

export const FolderUser = () => {
  return (
    <div className="relative">
      <Folder />

      <div className="absolute bottom-0.75 right-1 p-0 m-0">
        <Users className="w-4 h-4 m-0"/>
      </div>
    </div>
  );
};
