import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createLobby } from '../api/lobbyApi'
import { saveAdminToken } from '../lib/adminStorage'

export function HomePage() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleCreateLobby() {
    setLoading(true)
    setError(null)

    try {
      const created = await createLobby()
      saveAdminToken(created.id, created.admin_token)
      navigate(`/lobby/${created.id}`)
    } catch (err) {
      const message =
        err instanceof Error ? err.message : 'Failed to create lobby'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="home-page">
      <div className="home-page__content">
        <p className="home-page__eyebrow">Real-time croquis meetups</p>
        <h1>Croquis King</h1>
        <p className="home-page__lead">
          Create a lobby, share the link, and draw together with synchronized
          photos and timers — no screen sharing required.
        </p>

        <button
          type="button"
          className="button button--primary"
          onClick={() => void handleCreateLobby()}
          disabled={loading}
        >
          {loading ? 'Creating lobby…' : 'Create lobby'}
        </button>

        {error ? (
          <p className="home-page__error" role="alert">
            {error}
          </p>
        ) : null}
      </div>
    </main>
  )
}
