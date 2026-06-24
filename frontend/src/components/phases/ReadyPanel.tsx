import type { LobbySnapshot } from '../../types/lobby'
import { RoundDurationPicker } from '../RoundDurationPicker'
import { t } from '../../lib/i18n'

interface ReadyPanelProps {
  lobbyId: string
  snapshot: LobbySnapshot
  isAdmin: boolean
}

export function ReadyPanel({ lobbyId, snapshot, isAdmin }: ReadyPanelProps) {
  return (
    <section className="phase-panel phase-panel--ready" aria-live="polite">
      <div className="ready-panel__card">
        <p className="ready-panel__count">{snapshot.total_rounds}</p>
        <h2>{t('ready.photosReady')}</h2>
        <p>{t('ready.desc')}</p>
        <RoundDurationPicker
          lobbyId={lobbyId}
          minutes={snapshot.draw_duration_minutes}
          isAdmin={isAdmin}
        />
        <p className="ready-panel__hint">{t('ready.hint')}</p>
      </div>
    </section>
  )
}
