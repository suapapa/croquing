import type { LobbySnapshot } from '../../types/lobby'
import { t } from '../../lib/i18n'

interface ReadyPanelProps {
  snapshot: LobbySnapshot
}

export function ReadyPanel({ snapshot }: ReadyPanelProps) {
  return (
    <section className="phase-panel phase-panel--ready" aria-live="polite">
      <div className="ready-panel__card">
        <p className="ready-panel__count">{snapshot.total_rounds}</p>
        <h2>{t('ready.photosReady')}</h2>
        <p>{t('ready.desc')}</p>
        <p className="ready-panel__hint">{t('ready.hint')}</p>
      </div>
    </section>
  )
}
