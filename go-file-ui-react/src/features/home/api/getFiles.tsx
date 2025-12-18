import { useQuery } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { FileData } from "../Home";

const getFiles = async (path: string): Promise<FileData[]> => {
  const response = await api.get(`files/list/${path}`);
  return response.data;
};

export const useFiles = ({ path }: { path: string }) => {
  return useQuery({
    queryKey: ["files", path],
    queryFn: () => getFiles(path),
  });
};
