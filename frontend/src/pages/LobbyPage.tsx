import { Link, useParams } from 'react-router-dom'
import { useLobbySocket } from '../hooks/useLobbySocket'
import { isLobbyAdmin } from '../lib/adminStorage'

export function LobbyPage() {
  const { id } = useParams<{ id: string }>()
  const { snapshot, status, serverTimeOffsetMs, error } = useLobbySocket(id)

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

      <p className="lobby-page__status">Connection: {status}</p>

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
          <p>Server offset: {serverTimeOffsetMs} ms</p>
        </section>
      ) : status === 'connected' ? (
        <p className="lobby-page__status">Waiting for lobby state…</p>
      ) : null}
    </main>
  )
}
