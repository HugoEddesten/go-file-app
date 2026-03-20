import { useQuery } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { FileData } from "../Home";

const getFiles = async (vaultId: number, path: string): Promise<FileData[]> => {
  const response = await api.get(`files/${vaultId}/list${path}`);
  return response.data;
};

export const useFiles = ({ vaultId, path, enabled = true }: { vaultId: number, path: string, enabled?: boolean }) => {
  return useQuery({
    queryKey: ["files", vaultId, path],
    queryFn: () => getFiles(vaultId, path),
    enabled,
  });
};
