import { useMutation } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { FileData } from "../Home";
import { queryClient } from "../../../lib/queryClient";

const renameFile = async (
  vaultId: number,
  path: string,
  newName: string,
): Promise<FileData> => {
  const response = await api.put(`files/${vaultId}/rename${path}`, {
    newName: newName,
  });
  return response.data;
};

export const useRenameFile = ({
  path,
  vaultId,
}: {
  path: string;
  vaultId: number;
}) => {
  return useMutation({
    mutationKey: ["files", vaultId, path],
    mutationFn: (newName: string) => renameFile(vaultId, path, newName),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["files", vaultId] });
    },
  });
};
