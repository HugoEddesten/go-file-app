import { useMutation } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import { queryClient } from "../../../lib/queryClient";

const removeVaultUser = async (vaultId: number, vaultUserId: number) => {
  await api.delete(`vault/remove-vault-user/${vaultId}`, { data: { vaultUserId } });
};

export const useRemoveVaultUser = ({ vaultId }: { vaultId: number }) => {
  return useMutation({
    mutationFn: (vaultUserId: number) => removeVaultUser(vaultId, vaultUserId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaults"] });
    },
  });
};
