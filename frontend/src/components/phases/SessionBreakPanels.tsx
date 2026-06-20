import type { LobbySnapshot } from '../../types/lobby'
import { t } from '../../lib/i18n'

interface BetweenPanelProps {
  snapshot: LobbySnapshot
}

export function BetweenPanel({ snapshot }: BetweenPanelProps) {
  return (
    <section className="phase-panel phase-panel--between" aria-live="polite">
      <div className="between-panel__card">
        <h2>{t('break.takeBreather')}</h2>
        <p>{t('break.hiddenDesc')}</p>
        {snapshot.total_rounds > 0 ? (
          <p className="between-panel__round">
            {t('break.completedRound', {
              current: snapshot.current_round,
              total: snapshot.total_rounds,
            })}
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
        <h2>{t('break.sessionFinished')}</h2>
        <p>
          {snapshot.total_rounds === 1
            ? t('break.completedRoundsDesc', { count: snapshot.total_rounds })
            : t('break.completedRoundsDescPlural', { count: snapshot.total_rounds })}
        </p>
      </div>
    </section>
  )
}
