import { MaximizedSpinner } from "../../components/ui/maximizedSpinner";
import { Separator } from "../../components/ui/separator";
import { useAuth } from "../../hooks/useAuth";
import { useVaults } from "./api/getVaults";
import { VaultItem } from "./components/VaultItem";

export const Vaults = () => {
  const { data, isLoading } = useVaults({});
  const { userId } = useAuth();

  if (isLoading) return <MaximizedSpinner />;

  const userVaults = data?.filter((v) =>
    v.users.some((u) => u.id == userId && u.role === 1),
  );
  const sharedVaults = data?.filter((v) =>
    v.users.some((u) => u.id == userId && u.role !== 1),
  );

  return (
    <div className="h-full w-full flex justify-center p-4">
      <div className="flex flex-col gap-8 w-full max-w-3xl">
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <h3>My vaults</h3>
          </div>
          <Separator />
          <div className="flex flex-wrap gap-4">
            {userVaults?.length ?? 0 > 0
              ?
                userVaults?.map((v) => (
                  <VaultItem key={v.id} vault={v} />
                ))
              :
                <p className="text-xs text-muted-foreground truncate mt-0.5">No vaults found</p>
            }
          </div>
        </div>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <h3>Shared with me</h3>
          </div>
          <Separator />
          <div className="flex flex-wrap gap-4">
            {sharedVaults?.length ?? 0 > 0
              ?
                sharedVaults?.map((v) => (
                  <VaultItem key={v.id} vault={v} />
                ))
              :
                <p className="text-xs text-muted-foreground truncate mt-0.5">Vaults other people have shared with you will show up here</p>
            }
          </div>
        </div>
      </div>
    </div>
  );
};
