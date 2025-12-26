import type { ReactNode } from "react";
import { Nav } from "./nav";


export const DefaultLayout = ({ children }: { children: ReactNode }) => {
  return (
    <div className="w-full h-screen" style={{display: "grid", gridTemplateRows: "auto 1fr"}}>
      <div className="w-full">
        <Nav />
      </div>

      <div className="p-2 w-full">
        {children}
      </div>
    </div>
  );
}