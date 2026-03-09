import { api } from "../../../lib/api"

export type InviteInfo = {
  email: string
  vaultName: string
  token: string
}

export const getInviteInfo = async (token: string): Promise<InviteInfo> => {
  const response = await api.get(`invites/${token}`)
  return response.data
}
