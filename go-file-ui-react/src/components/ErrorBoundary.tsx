import { Component, type ErrorInfo, type ReactNode } from "react";
import { Button } from "./ui/button";

type Props = { children: ReactNode };
type State = { error: Error | null };

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    // if (import.meta.env.DEV) {
    //   // Let Vite's error overlay take over in development
    //   throw error;
    // }
    console.error("Uncaught error:", error, info);
  }

  render() {
    if (this.state.error) { //&& import.meta.env.PROD) {
      return (
        <div className="flex h-screen w-screen flex-col items-center justify-center gap-4 p-8 text-center">
          <h1 className="text-2xl font-semibold">Something went wrong</h1>
          <p className="text-muted-foreground max-w-sm text-sm">
            An unexpected error occurred. Try reloading the page. If the
            problem persists, contact support.
          </p>
          <Button onClick={() => window.location.reload()}>Reload page</Button>
        </div>
      );
    }

    return this.props.children;
  }
}
