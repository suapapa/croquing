import type { LobbySnapshot } from '../../types/lobby'

interface BetweenPanelProps {
  snapshot: LobbySnapshot
}

export function BetweenPanel({ snapshot }: BetweenPanelProps) {
  return (
    <section className="phase-panel phase-panel--between" aria-live="polite">
      <div className="between-panel__card">
        <h2>Take a breather</h2>
        <p>The reference photo is hidden until the next round starts.</p>
        {snapshot.total_rounds > 0 ? (
          <p className="between-panel__round">
            Completed round {snapshot.current_round} of {snapshot.total_rounds}
          </p>
        ) : null}
      </div>
    </section>
  )
}

interface FinishedPanelProps {
  snapshot: LobbySnapshot
}

export function FinishedPanel({ snapshot }: FinishedPanelProps) {
  return (
    <section className="phase-panel phase-panel--finished" aria-live="polite">
      <div className="finished-panel__card">
        <h2>Session finished</h2>
        <p>
          You completed {snapshot.total_rounds}{' '}
          {snapshot.total_rounds === 1 ? 'round' : 'rounds'}. Great work
          everyone.
        </p>
      </div>
    </section>
  )
}
