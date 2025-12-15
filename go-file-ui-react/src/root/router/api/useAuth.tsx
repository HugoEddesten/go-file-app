import { useQuery } from "@tanstack/react-query"
import { api } from "../../../lib/api"

export type UseMeResponse = {
  email?: string,
  userId?: number,
} 

export const useAuthQuery = () => {
  return useQuery({
    queryKey: ["auth", "me"],
    queryFn: async () => {
      const response = await api.get<UseMeResponse>("/auth/me")
      return response.data
    },
  })
}