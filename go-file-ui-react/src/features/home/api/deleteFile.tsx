import { useMutation } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import { queryClient } from "../../../lib/queryClient";

const deleteFile = async (vaultId: number, path: string): Promise<void> => {
  await api.delete(`files/${vaultId}/delete${path}`);
};

export const useDeleteFile = ({
  path,
  vaultId,
}: {
  path: string;
  vaultId: number;
}) => {
  return useMutation({
    mutationKey: ["files", vaultId, path, "delete"],
    mutationFn: () => deleteFile(vaultId, path),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["files", vaultId] });
      queryClient.invalidateQueries({ queryKey: ["vault", vaultId] });
    },
  });
};
