// DragProvider.tsx
import { useEffect, useRef, useState } from "react";
import { DragContext } from "./DragContext";
import type { ClassNameValue } from "tailwind-merge";

type DragProviderProps = {
  children: React.ReactNode;
  className?: ClassNameValue;
  detectAll?: boolean;
  onEnter?: (e: DragEvent) => void;
  onLeave?: (e: DragEvent) => void;
  onDrop?: (e: DragEvent) => void;
};

export const DragProvider = ({
  children,
  className,
  detectAll = false,
  onEnter,
  onLeave,
  onDrop,
}: DragProviderProps) => {
  const [isOver, setIsOver] = useState(false);
  const dragCounter = useRef(0);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;

    const handleDragEnter = (e: DragEvent) => {
      dragCounter.current++;

      if (dragCounter.current === 1) {
        setIsOver(true);
        onEnter?.(e);
      }
    };

    const handleDragLeave = (e: DragEvent) => {
      dragCounter.current--;

      if (dragCounter.current === 0) {
        setIsOver(false);
        onLeave?.(e);
      }
    };

    const handleDrop = (e: DragEvent) => {
      dragCounter.current = 0;
      setIsOver(false);
      onDrop?.(e);
    };

    if (detectAll) {
      window.addEventListener("dragenter", handleDragEnter);
      window.addEventListener("dragleave", handleDragLeave);
      window.addEventListener("dragover", (e) => e.preventDefault());
    } else {
      el.addEventListener("dragenter", handleDragEnter);
      el.addEventListener("dragleave", handleDragLeave);
      el.addEventListener("dragover", (e) => e.preventDefault());
    }

    el.addEventListener("drop", handleDrop);

    return () => {
      if (detectAll) {
        window.removeEventListener("dragenter", handleDragEnter);
        window.removeEventListener("dragleave", handleDragLeave);
        window.removeEventListener("dragover", (e) => e.preventDefault());
      } else {
        el.removeEventListener("dragenter", handleDragEnter);
        el.removeEventListener("dragleave", handleDragLeave);
        el.removeEventListener("dragover", (e) => e.preventDefault());
      }
      el.removeEventListener("drop", handleDrop);
    };
  }, [onEnter, onLeave, onDrop]);

  return (
    <DragContext.Provider value={{ isOver }}>
      <div className={`${className}`} ref={ref}>
        {children}
      </div>
    </DragContext.Provider>
  );
};
