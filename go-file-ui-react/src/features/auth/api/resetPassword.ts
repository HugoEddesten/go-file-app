import { api } from "../../../lib/api"

export const resetPassword = async (token: string, password: string) => {
  return api.post(`auth/reset-password/${token}`, { password })
}
