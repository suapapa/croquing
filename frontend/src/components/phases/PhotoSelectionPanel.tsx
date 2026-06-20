import { useState } from 'react'
import { ApiError } from '../../api/client'
import { markLobbyReady, setLobbyPhotos } from '../../api/lobbyApi'
import type { LobbySnapshot, Photo } from '../../types/lobby'
import { PixabaySearchPanel } from '../search/PixabaySearchPanel'

const RECOMMENDED_COUNT = 5

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

  if (!isAdmin && (snapshot.phase === 'WAITING' || snapshot.phase === 'SELECTING')) {
    return (
      <section className="phase-panel">
        <p>The admin is choosing reference photos. Photos stay hidden until drawing begins.</p>
      </section>
    )
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
        />
      ) : null}

      {awaitingReady ? (
        <div className="photo-selection__summary">
          <h2>{snapshot.selected_count} photos saved</h2>
          <p>
            Review your selection, then shuffle and lock the order. Recommended
            session size is about {RECOMMENDED_COUNT} photos.
          </p>
          {selectedPhotos.length > 0 ? (
            <ul className="photo-selection__list">
              {selectedPhotos.map((photo) => (
                <li key={photo.pixabay_id}>
                  <img src={photo.preview_url} alt="" loading="lazy" />
                </li>
              ))}
            </ul>
          ) : null}
          <button
            type="button"
            className="button button--primary"
            disabled={loading}
            onClick={() => void handleMarkReady()}
          >
            {loading ? 'Shuffling…' : 'Selection complete'}
          </button>
        </div>
      ) : null}

      {canConfirm ? (
        <div className="photo-selection__actions">
          <button
            type="button"
            className="button button--primary"
            disabled={loading || selectedPhotos.length === 0}
            onClick={() => void handleSaveSelection()}
          >
            {loading ? 'Saving…' : `Save ${selectedPhotos.length} photos`}
          </button>
        </div>
      ) : null}

      {error ? (
        <p className="photo-selection__error" role="alert">
          {error}
        </p>
      ) : null}
    </section>
  )
}
