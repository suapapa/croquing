import { useState } from 'react'
import { ApiError } from '../../api/client'
import { markLobbyReady, reopenLobbyPhotos, setLobbyPhotos } from '../../api/lobbyApi'
import type { LobbySnapshot, Photo } from '../../types/lobby'
import { PixabaySearchPanel } from '../search/PixabaySearchPanel'
import { PhotoReviewPanel } from './PhotoReviewPanel'
import { ParticipantWaitPanel } from './ParticipantWaitPanel'

interface PhotoSelectionPanelProps {
  lobbyId: string
  snapshot: LobbySnapshot
  selectedPhotos: Photo[]
  onSelectionChange: (photos: Photo[]) => void
  isAdmin: boolean
}

export function PhotoSelectionPanel({
  lobbyId,
  snapshot,
  selectedPhotos,
  onSelectionChange,
  isAdmin,
}: PhotoSelectionPanelProps) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const canSearch = isAdmin && snapshot.phase === 'WAITING'
  const canConfirm =
    isAdmin &&
    snapshot.phase === 'WAITING' &&
    selectedPhotos.length > 0
  const awaitingReady = isAdmin && snapshot.phase === 'SELECTING'

  async function handleSaveSelection() {
    setLoading(true)
    setError(null)
    try {
      await setLobbyPhotos(lobbyId, selectedPhotos)
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : 'Failed to save selection'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  async function handleMarkReady() {
    setLoading(true)
    setError(null)
    try {
      await markLobbyReady(lobbyId)
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : 'Failed to confirm selection'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  async function handleReopenSelection() {
    setLoading(true)
    setError(null)
    try {
      await reopenLobbyPhotos(lobbyId)
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : 'Failed to reopen photo selection'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  if (!isAdmin && (snapshot.phase === 'WAITING' || snapshot.phase === 'SELECTING')) {
    return <ParticipantWaitPanel />
  }

  if (!isAdmin) {
    return null
  }

  return (
    <section className="phase-panel">
      {canSearch ? (
        <PixabaySearchPanel
          lobbyId={lobbyId}
          selectedPhotos={selectedPhotos}
          onSelectionChange={onSelectionChange}
          footerStart={
            canConfirm ? (
              <button
                type="button"
                className="button button--primary"
                disabled={loading || selectedPhotos.length === 0}
                onClick={() => void handleSaveSelection()}
              >
                {loading ? 'Saving…' : `Save ${selectedPhotos.length} photos`}
              </button>
            ) : null
          }
        />
      ) : null}

      {awaitingReady ? (
        <PhotoReviewPanel
          photos={selectedPhotos}
          loading={loading}
          onEdit={() => void handleReopenSelection()}
          onConfirm={() => void handleMarkReady()}
        />
      ) : null}

      {error ? (
        <p className="photo-selection__error" role="alert">
          {error}
        </p>
      ) : null}
    </section>
  )
}
