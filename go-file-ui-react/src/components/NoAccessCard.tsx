import { LockKeyhole } from "lucide-react";
import { Card } from "./ui/card";
import { cn } from "../lib/utils";

export const NoAccessCard = (props: React.HTMLAttributes<HTMLDivElement>) => {
  return (
    <Card {...props} className={cn("flex flex-col items-center justify-center gap-3 w-full h-full text-destructive", props?.className)}>
      <div className="rounded-full bg-destructive/10 p-4">
        <LockKeyhole className="w-8 h-8" />
      </div>
      <div className="flex flex-col items-center gap-1">
        <p className="font-semibold text-base">Access Denied</p>
        <p className="text-sm text-destructive/70">You don't have permission to view this vault.</p>
      </div>
    </Card>
  );
};
