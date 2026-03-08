import { useMutation } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { VaultUserRole } from "../types";
import { queryClient } from "../../../lib/queryClient";

type AddVaultUserPayload = {
  role: VaultUserRole;
  path: string;
  email: string;
};

const addVaultUser = async (vaultId: number, payload: AddVaultUserPayload) => {
  await api.post(`vault/assign-user/${vaultId}`, payload);
};

export const useAddVaultUser = ({ vaultId }: { vaultId: number }) => {
  return useMutation({
    mutationFn: (payload: AddVaultUserPayload) => addVaultUser(vaultId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaults"] });
    },
  });
};
