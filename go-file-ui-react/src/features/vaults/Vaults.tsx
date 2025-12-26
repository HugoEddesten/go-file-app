import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";
import { MaximizedSpinner } from "../../components/ui/maximizedSpinner";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "../../components/ui/tabs";
import { useAuth } from "../../hooks/useAuth";
import { useVaults } from "./api/getVaults";
import { VaultItem } from "./components/VaultItem";

export const Vaults = () => {
  const { data, isLoading } = useVaults({});
  const { userId } = useAuth();

  if (isLoading) return <MaximizedSpinner />;

  const userVaults = data?.filter((v) =>
    v.users.some((u) => u.id == userId && u.role === 1)
  );
  const sharedVaults = data?.filter((v) =>
    v.users.some((u) => u.id == userId && u.role !== 1)
  );

  return (
    <div className="h-full w-full flex justify-center p-4">
      <Tabs defaultValue="my_vaults">
        <TabsList>
          <TabsTrigger value="my_vaults">My vaults</TabsTrigger>
          <TabsTrigger value="shared_vaults">Shared vaults</TabsTrigger>
        </TabsList>
        <TabsContent value="my_vaults">
          <Card className="w-4xl">
            <CardHeader>
              <CardTitle className="flex flex-col gap-2">
                <h3>My vaults</h3>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex">
                {userVaults?.map((v) => (
                  <VaultItem vault={v} />
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="shared_vaults">
          <Card className="w-4xl">
            <CardHeader>
              <CardTitle>
                <h3>Shared vaults</h3>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex">
                {sharedVaults?.map((v) => (
                  <VaultItem vault={v} />
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
};
