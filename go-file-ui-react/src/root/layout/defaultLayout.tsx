import type { ReactNode } from "react";


export const DefaultLayout = ({ children }: { children: ReactNode }) => {
  return (
    <div className="w-full h-screen p-2">
      {children}
    </div>
  );
}