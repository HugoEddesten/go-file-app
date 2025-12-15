import { QueryClientProvider } from "@tanstack/react-query";
import { AppRouter } from "./root/router/router";
import { queryClient } from "./lib/queryClient";
import { DefaultLayout } from "./root/layout/defaultLayout";
import { Toaster } from "./components/ui/sonner";

function App() {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <DefaultLayout>
          <AppRouter />
          <Toaster />
        </DefaultLayout>
      </QueryClientProvider>
    </>
  );
}

export default App;
