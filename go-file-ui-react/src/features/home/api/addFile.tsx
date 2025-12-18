import { useMutation } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { FileData } from "../Home";
import { queryClient } from "../../../lib/queryClient";

const createFile = async (path: string, ext?: string): Promise<FileData> => {
  const response = await api.post(`files/create${path}?ext=${ext ?? ""}`);
  return response.data;
};

export const useCreateFile = ({ path }: { path: string }) => {
  return useMutation({
    mutationKey: ["files", path],
    mutationFn: (ext?: string) => createFile(path, ext),
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ["files", path]})
    }
  });
};
