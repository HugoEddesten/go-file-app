import { api } from "../../../lib/api"

export const sendResetPasswordEmail = async (email: string) => {
  return api.post("auth/reset-password", { email })
}
