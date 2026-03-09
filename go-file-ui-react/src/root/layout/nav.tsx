import { NavLink, useNavigate } from "react-router-dom";
import { Separator } from "../../components/ui/separator";
import { Database } from "lucide-react";
import { useAuthQuery } from "../router/api/useAuth";
import { Button } from "../../components/ui/button";
import { logout } from "../../features/auth/api/logout";

const navLinkClass = ({ isActive }: { isActive: boolean }) =>
  `text-sm px-3 py-1.5 rounded-md transition-colors ${
    isActive
      ? "bg-accent text-foreground font-medium"
      : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
  }`;

export const Nav = () => {
  const { data: user } = useAuthQuery();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate("/login");
  };

  return (
    <div className="w-full h-full flex flex-col">
      <div className="flex items-center gap-1 px-4 h-12">
        <NavLink
          to={"/"}
          className="rounded-md bg-primary/10 p-2 text-primary shrink-0 mr-3"
        >
          <Database size={18} />
        </NavLink>
        <NavLink to={"/"} end hidden={!user?.userId} className={navLinkClass}>
          Home
        </NavLink>
        <NavLink
          to={"/profile"}
          hidden={!user?.userId}
          className={navLinkClass}
        >
          Profile
        </NavLink>
        {user?.userId && (
          <Button
            variant="ghost"
            className="ml-auto text-muted-foreground"
            onClick={handleLogout}
          >
            Log out
          </Button>
        )}
      </div>
      <Separator />
    </div>
  );
};
