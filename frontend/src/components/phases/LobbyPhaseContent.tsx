import type { LobbySnapshot } from '../../types/lobby'
import type { Photo } from '../../types/lobby'
import { DrawingPanel } from './DrawingPanel'
import { PhotoSelectionPanel } from './PhotoSelectionPanel'
import { ReadyPanel } from './ReadyPanel'
import { BetweenPanel, FinishedPanel } from './SessionBreakPanels'

interface LobbyPhaseContentProps {
  lobbyId: string
  snapshot: LobbySnapshot
  serverTimeOffsetMs: number
  isAdmin: boolean
  selectedPhotos: Photo[]
  onSelectionChange: (photos: Photo[]) => void
}

export function LobbyPhaseContent({
  lobbyId,
  snapshot,
  serverTimeOffsetMs,
  isAdmin,
  selectedPhotos,
  onSelectionChange,
}: LobbyPhaseContentProps) {
  switch (snapshot.phase) {
    case 'WAITING':
    case 'SELECTING':
      return (
        <PhotoSelectionPanel
          lobbyId={lobbyId}
          snapshot={snapshot}
          selectedPhotos={selectedPhotos}
          onSelectionChange={onSelectionChange}
          isAdmin={isAdmin}
        />
      )
    case 'READY':
      return <ReadyPanel snapshot={snapshot} />
    case 'DRAWING':
      return (
        <DrawingPanel
          snapshot={snapshot}
          serverTimeOffsetMs={serverTimeOffsetMs}
        />
      )
    case 'BETWEEN_ROUNDS':
      return <BetweenPanel snapshot={snapshot} />
    case 'FINISHED':
      return <FinishedPanel lobbyId={lobbyId} snapshot={snapshot} />
    default:
      return null
  }
}
