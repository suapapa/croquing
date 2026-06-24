import type { LobbySnapshot } from '../../types/lobby'
import { RoundDurationPicker } from '../RoundDurationPicker'
import { t } from '../../lib/i18n'

interface BetweenPanelProps {
  lobbyId: string
  snapshot: LobbySnapshot
  isAdmin: boolean
}

export function BetweenPanel({ lobbyId, snapshot, isAdmin }: BetweenPanelProps) {
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
        <RoundDurationPicker
          lobbyId={lobbyId}
          minutes={snapshot.draw_duration_minutes}
          isAdmin={isAdmin}
        />
      </div>
    </section>
  )
}

interface FinishedPanelProps {
  lobbyId: string
  snapshot: LobbySnapshot
}

export function FinishedPanel({ lobbyId, snapshot }: FinishedPanelProps) {
  const downloadHref = `/api/lobbies/${encodeURIComponent(lobbyId)}/photos/download`

  return (
    <section className="phase-panel phase-panel--finished" aria-live="polite">
      <div className="finished-panel__card">
        <h2>{t('break.sessionFinished')}</h2>
        <p>
          {snapshot.total_rounds === 1
            ? t('break.completedRoundsDesc', { count: snapshot.total_rounds })
            : t('break.completedRoundsDescPlural', {
                count: snapshot.total_rounds,
              })}
        </p>
        {snapshot.total_rounds > 0 ? (
          <a className="finished-panel__download" href={downloadHref} download>
            {t('break.downloadPhotos')}
          </a>
        ) : null}
      </div>
    </section>
  )
}
