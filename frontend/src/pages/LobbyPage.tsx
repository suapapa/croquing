import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getLobby } from '../api/lobbyApi'
import { ApiError } from '../api/client'
import { isLobbyAdmin } from '../lib/adminStorage'
import type { LobbySnapshot } from '../types/lobby'

export function LobbyPage() {
  const { id } = useParams<{ id: string }>()
  const [snapshot, setSnapshot] = useState<LobbySnapshot | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) {
      return
    }

    let cancelled = false

    async function loadLobby() {
      setLoading(true)
      setError(null)

      const lobbyId = id
      if (!lobbyId) {
        return
      }

      try {
        const data = await getLobby(lobbyId)
        if (!cancelled) {
          setSnapshot(data)
        }
      } catch (err) {
        if (!cancelled) {
          const message =
            err instanceof ApiError ? err.message : 'Failed to load lobby'
          setError(message)
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    void loadLobby()

    return () => {
      cancelled = true
    }
  }, [id])

  if (!id) {
    return (
      <main className="lobby-page">
        <p>Invalid lobby link.</p>
        <Link to="/">Back home</Link>
      </main>
    )
  }

  return (
    <main className="lobby-page">
      <header className="lobby-page__header">
        <Link to="/" className="lobby-page__brand">
          Croquis King
        </Link>
        <span className="lobby-page__role">
          {isLobbyAdmin(id) ? 'Admin' : 'Participant'}
        </span>
      </header>

      {loading ? <p className="lobby-page__status">Loading lobby…</p> : null}

      {error ? (
        <p className="lobby-page__error" role="alert">
          {error}
        </p>
      ) : null}

      {snapshot ? (
        <section className="lobby-page__summary" aria-live="polite">
          <h1>Lobby {snapshot.id.slice(0, 8)}</h1>
          <p>Phase: {snapshot.phase}</p>
          <p>Participants: {snapshot.participant_count}</p>
        </section>
      ) : null}
    </main>
  )
}
