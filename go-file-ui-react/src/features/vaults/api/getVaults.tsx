import { useQuery } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import type { Vault } from "../types";

const getVaults = async (): Promise<Vault[]> => {
  const response = await api.get(`vault/get-user-vaults`);
  return response.data;
};

export const useVaults = ({ }: { }) => {
  return useQuery({
    queryKey: ["vaults"],
    queryFn: () => getVaults(),
  });
};
