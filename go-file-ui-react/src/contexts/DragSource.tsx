import { useDragContext } from "./DragContext";

type DragSourceProps<T> = Omit<
  React.HTMLAttributes<HTMLDivElement>,
  "onDragStart" | "onDragEnd"
> & {
  children?: React.ReactNode;
  payload: T;
  onDragStart?: (e: DragEvent, payload: T) => void;
  onDragEnd?: (e: DragEvent) => void;
};

export const DragSource = <T,>({
  children,
  payload,
  onDragStart,
  onDragEnd,
  ...rest
}: DragSourceProps<T>) => {
  const { setPayload } = useDragContext();

  return (
    <div
      {...rest}
      draggable
      onDragStart={(e) => {
        setPayload(payload);
        onDragStart?.(e.nativeEvent, payload);
      }}
      onDragEnd={(e) => {
        setPayload(null);
        onDragEnd?.(e.nativeEvent);
      }}
    >
      {children}
    </div>
  );
};
