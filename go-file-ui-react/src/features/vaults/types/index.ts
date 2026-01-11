export type Vault = {
  id: number;
  name: string;
  users: VaultUser[];
};

export type VaultUser = {
  id: number;
  email: string;
  role: VaultUserRole;
  path: string;
};

export const VaultUserRole = {
  OWNER: 1,
  ADMIN: 2,
  EDITOR: 3,
  VIEWER: 4,
} as const;

export type VaultUserRole = (typeof VaultUserRole)[keyof typeof VaultUserRole];

export const getVaultUserRole = (role: VaultUserRole) => {
  switch (role) {
    case VaultUserRole.OWNER:
      return "Owner";
    case VaultUserRole.ADMIN:
      return "Admin";
    case VaultUserRole.EDITOR:
      return "Editor";
    case VaultUserRole.VIEWER:
      return "Viewer";
  }
}
