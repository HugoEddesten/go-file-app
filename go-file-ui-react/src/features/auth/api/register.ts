import { api } from "../../../lib/api"
import { queryClient } from "../../../lib/queryClient"


export type AuthRequest = {
  email: string,
  password: string,
}

export const register = async (request: AuthRequest) => {
  const response = await api.post("auth/register", request)
  queryClient.invalidateQueries({ queryKey: ["auth", "me"] })
  return response
}