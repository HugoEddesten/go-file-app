import {
  CircleCheckIcon,
  InfoIcon,
  Loader2Icon,
  OctagonXIcon,
  TriangleAlertIcon,
} from "lucide-react";
import { useTheme } from "next-themes";
import { Toaster as Sonner, type ToasterProps } from "sonner";

const Toaster = ({ ...props }: ToasterProps) => {
  const { theme = "system" } = useTheme();

  return (
    <Sonner
      theme={theme as ToasterProps["theme"]}
      className="toaster group"
      icons={{
        success: <CircleCheckIcon className="size-4" />,
        info: <InfoIcon className="size-4" />,
        warning: <TriangleAlertIcon className="size-4" />,
        error: <OctagonXIcon className="size-4" />,
        loading: <Loader2Icon className="size-4 animate-spin" />,
      }}
      style={
        {
          "--normal-bg": "var(--popover)",
          "--normal-text": "var(--popover-foreground)",
          "--normal-border": "var(--border)",
          
          "--success-bg": "hsl(143, 85%, 96%)",
          "--success-text": "hsl(140, 100%, 27%)",
          "--success-border": "hsl(145, 92%, 91%)",
          
          "--error-bg": "hsl(359, 100%, 97%)",
          "--error-text": "hsl(360, 100%, 45%)",
          "--error-border": "hsl(359, 100%, 94%)",
          
          "--warning-bg": "hsl(48, 100%, 96%)",
          "--warning-text": "hsl(25, 95%, 53%)",
          "--warning-border": "hsl(48, 96%, 89%)",
          
          "--info-bg": "hsl(208, 100%, 97%)",
          "--info-text": "hsl(210, 92%, 45%)",
          "--info-border": "hsl(221, 91%, 91%)",
          
          "--border-radius": "var(--radius)",
        } as React.CSSProperties
      }
      {...props}
    />
  );
};

export { Toaster };
