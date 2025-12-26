

import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { FileMetadata } from "../types";

const getFile = async (path: string): Promise<FileMetadata> => {
  const response = await api.get(`files/1/metadata/${path}`);
  return response.data;
};

export const useFileDetails = ({ path }: { path: string; }) => {
  return useQuery({
    queryKey: ["file-metadata", 1, path],
    queryFn: () => getFile(path),
  });
};