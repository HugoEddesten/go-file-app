import { useMutation } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import { queryClient } from "../../../lib/queryClient";
import type { VaultUserRole } from "../types";

type UpdateVaultUserPayload = {
  vaultUserId: number;
  role: VaultUserRole;
  path: string;
};

const updateVaultUser = async (vaultId: number, payload: UpdateVaultUserPayload) => {
  await api.put(`vault/update-vault-user/${vaultId}`, payload);
};

export const useUpdateVaultUser = ({ vaultId }: { vaultId: number }) => {
  return useMutation({
    mutationFn: (payload: UpdateVaultUserPayload) => updateVaultUser(vaultId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaults"] });
    },
  });
};
