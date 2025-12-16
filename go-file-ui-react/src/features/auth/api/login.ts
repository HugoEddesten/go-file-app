import { api } from "../../../lib/api"
import { queryClient } from "../../../lib/queryClient"
import type { AuthRequest } from "./register"

export const login = async (request: AuthRequest) => {
  const response = await api.post("auth/login", request)
  queryClient.invalidateQueries({ queryKey: ["auth", "me"] })
  return response.data
}