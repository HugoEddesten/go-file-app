import { useEffect, useRef, useState, type ReactNode } from "react";

export type DragState = {
  isDragging: boolean,
  data?: Record<string, any>,
}


export const useDragOver = ({
  target = "window",
  onEnter,
  onLeave,
  onDrop,
}: {
  target: string;
  onEnter?: (e: DragEvent) => void;
  onLeave?: (e: DragEvent) => void;
  onDrop?: (e: DragEvent) => void;
}) => {
  const dragCounter = useRef(0);
  const [state, setState] = useState<DragState>({isDragging: false})

  useEffect(() => {
    const handleDragEnter = (e: DragEvent) => {
      dragCounter.current += 1;

      if (dragCounter.current === 1) {
        setState({isDragging: true})
        onEnter?.(e);
      }
    };

    const handleDragLeave = (e: DragEvent) => {
      dragCounter.current -= 1;

      if (dragCounter.current === 0) {
        setState({isDragging: false})
        onLeave?.(e);
      }
    };

    const handleDrop = (e: DragEvent) => {
      dragCounter.current = 0;
      setState({isDragging: false})
      onDrop?.(e)
    };

    window.addEventListener("dragenter", handleDragEnter);
    window.addEventListener("dragleave", handleDragLeave);
    window.addEventListener("drop", handleDrop);

    return () => {
      window.removeEventListener("dragenter", handleDragEnter);
      window.removeEventListener("dragleave", handleDragLeave);
      window.removeEventListener("drop", handleDrop);
    };
  }, [onEnter, onLeave]);


  return state
};
