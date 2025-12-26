import { QueryClientProvider } from "@tanstack/react-query";
import { AppRouter } from "./root/router/router";
import { queryClient } from "./lib/queryClient";
import { Toaster } from "./components/ui/sonner";

function App() {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <AppRouter />
        <Toaster />
      </QueryClientProvider>
    </>
  );
}

export default App;
