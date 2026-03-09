import { createContext, useContext } from "react"

type DragContextValue = {
  isOver: boolean
  payload: unknown | undefined
  setPayload: (payload: unknown) => void
}

export const DragContext = createContext<DragContextValue | null>(null)

export const useDragContext = () => {
  const ctx = useContext(DragContext)
  if (!ctx) throw new Error("useDragContext must be used inside DragStateProvider")
  return ctx
}

export const useDragPayload = <T,>(): T | null => {
  const { payload } = useDragContext()
  return payload as T | null
}
