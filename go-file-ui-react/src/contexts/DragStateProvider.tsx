import { useState } from "react"
import { DragContext } from "./DragContext"

export const DragStateProvider = ({ children }: { children: React.ReactNode }) => {
  const [payload, setPayload] = useState<unknown>(null)

  return (
    <DragContext.Provider value={{ isOver: false, payload, setPayload }}>
      {children}
    </DragContext.Provider>
  )
}
