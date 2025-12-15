import { createContext, useContext } from "react"

type DragContextValue = {
  isOver: boolean
}

export const DragContext = createContext<DragContextValue | null>(null)

export const useDragContext = () => {
  const ctx = useContext(DragContext)
  if (!ctx) {
    throw new Error("useDragContext must be used inside DragProvider")
  }
  return ctx
}