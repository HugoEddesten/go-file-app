

import { useQuery } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { FileMetadata } from "../types";

const getFile = async (path: string, vaultId: number): Promise<FileMetadata> => {
  const response = await api.get(`files/${vaultId}/metadata${path}`);
  return response.data;
};

export const useFileDetails = ({ path, vaultId }: { path: string; vaultId: number }) => {
  return useQuery({
    queryKey: ["file-metadata", vaultId, path],
    queryFn: () => getFile(path, vaultId),
  });
};