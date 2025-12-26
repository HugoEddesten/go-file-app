import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'

interface VaultState {
  vaultId: number | null
  setVaultId: (vaultId: number | null) => void
  clearVaultId: () => void
}

export const useVaultStore = create<VaultState>()(
  persist(
    (set) => ({
      vaultId: null,
      setVaultId: (vaultId) => set({ vaultId }),
      clearVaultId: () => set({ vaultId: null }),
    }),
    {
      name: 'vault-storage',
      storage: createJSONStorage(() => localStorage),
    }
  )
)