import { UserPen, UserRoundCog, UserSearch, UserStar } from "lucide-react";
import { VaultUserRole } from "../../vaults/types";
import type { ClassNameValue } from "tailwind-merge";

export const VaultUserIcon = ({
  role,
  className,
}: {
  role: VaultUserRole;
  className: ClassNameValue;
}) => {
  switch (role) {
    case VaultUserRole.OWNER:
      return <UserStar className={className as string} />;
    case VaultUserRole.ADMIN:
      return <UserRoundCog className={className as string} />;
    case VaultUserRole.EDITOR:
      return <UserPen className={className as string} />;
    case VaultUserRole.VIEWER:
      return <UserSearch className={className as string} />;
  }
};
