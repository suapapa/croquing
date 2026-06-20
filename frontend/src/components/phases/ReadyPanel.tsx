import type { LobbySnapshot } from '../../types/lobby'

interface ReadyPanelProps {
  snapshot: LobbySnapshot
}

export function ReadyPanel({ snapshot }: ReadyPanelProps) {
  return (
    <section className="phase-panel phase-panel--ready" aria-live="polite">
      <div className="ready-panel__card">
        <p className="ready-panel__count">{snapshot.total_rounds}</p>
        <h2>photos ready</h2>
        <p>
          The order is shuffled and hidden. Thumbnails stay off until each draw
          round begins.
        </p>
        <p className="ready-panel__hint">Waiting for the admin to start…</p>
      </div>
    </section>
  )
}
