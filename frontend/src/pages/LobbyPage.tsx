import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { AdminControls } from '../components/admin/AdminControls'
import { LobbyLayout } from '../components/layout/LobbyLayout'
import { LobbyPhaseContent } from '../components/phases/LobbyPhaseContent'
import { useLobbySocket } from '../hooks/useLobbySocket'
import { isLobbyAdmin } from '../lib/adminStorage'
import type { Photo } from '../types/lobby'

export function LobbyPage() {
  const { id } = useParams<{ id: string }>()
  const { snapshot, status, serverTimeOffsetMs, error } = useLobbySocket(id)
  const [selectedPhotos, setSelectedPhotos] = useState<Photo[]>([])
  const isAdmin = id ? isLobbyAdmin(id) : false
  const isDrawing = snapshot?.phase === 'DRAWING'

  if (!id) {
    return (
      <main className="lobby-page">
        <p>Invalid lobby link.</p>
        <Link to="/">Back home</Link>
      </main>
    )
  }

  return (
    <main className={`lobby-page${isDrawing ? ' lobby-page--drawing' : ''}`}>
      <LobbyLayout
        lobbyId={id}
        isAdmin={isAdmin}
        snapshot={snapshot}
        connectionStatus={status}
        connectionError={error}
      >
        {isAdmin && snapshot ? (
          <AdminControls lobbyId={id} snapshot={snapshot} />
        ) : null}

        {snapshot ? (
          <LobbyPhaseContent
            lobbyId={id}
            snapshot={snapshot}
            serverTimeOffsetMs={serverTimeOffsetMs}
            isAdmin={isAdmin}
            selectedPhotos={selectedPhotos}
            onSelectionChange={setSelectedPhotos}
          />
        ) : null}
      </LobbyLayout>
    </main>
  )
}
