import { useMutation } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { FileData } from "../Home";
import { queryClient } from "../../../lib/queryClient";

const createFile = async (path: string, vaultId: number, ext?: string): Promise<FileData> => {
  const response = await api.post(`files/${vaultId}/create${path}?ext=${ext ?? ""}`);
  return response.data;
};

export const useCreateFile = ({ path, vaultId }: { path: string, vaultId: number }) => {
  return useMutation({
    mutationKey: ["files", vaultId, path],
    mutationFn: (ext?: string) => createFile(path, vaultId, ext),
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ["files", vaultId, path]})
    }
  });
};
