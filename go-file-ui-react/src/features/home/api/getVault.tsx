import { useQuery } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { Vault } from "../../vaults/types";

const getVault = async (vaultId: number): Promise<Vault> => {
  const response = await api.get(`vault/get-vault/${vaultId}`);
  return response.data;
};

export const useVault = (vaultId: number) => {
  return useQuery({
    queryKey: ["vault", vaultId],
    queryFn: () => getVault(vaultId),
  });
};
