import { api } from "../../../lib/api"
import { queryClient } from "../../../lib/queryClient"

export const logout = async () => {
  await api.post("auth/logout")
  queryClient.removeQueries();
}
