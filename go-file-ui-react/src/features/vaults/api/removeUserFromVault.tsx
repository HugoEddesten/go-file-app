import { useMutation } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import { queryClient } from "../../../lib/queryClient";

const removeUserFromVault = async (vaultId: number, userId: number) => {
  await api.delete(`vault/remove-user/${vaultId}`, { data: { userId } });
};

export const useRemoveUserFromVault = ({ vaultId }: { vaultId: number }) => {
  return useMutation({
    mutationFn: (userId: number) => removeUserFromVault(vaultId, userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaults"] });
    },
  });
};
