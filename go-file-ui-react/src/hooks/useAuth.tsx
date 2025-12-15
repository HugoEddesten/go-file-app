import { queryClient } from "../lib/queryClient"
import type { UseMeResponse } from "../root/router/api/useAuth";


export const useAuth = (): UseMeResponse => {
  const authData = queryClient.getQueryData<UseMeResponse>(["auth", "me"]) ?? {};
  return authData;
}