import { api } from "../../../lib/api"
import { queryClient } from "../../../lib/queryClient"


type LoginRequest = {
  email: string,
  password: string,
}

export const login = async (request: LoginRequest) => {
  const response = await api.post("auth/login", request)
  queryClient.invalidateQueries({ queryKey: ["auth", "me"] })
  return response.data
}