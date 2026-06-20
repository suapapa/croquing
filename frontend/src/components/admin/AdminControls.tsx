import {
  finishLobbySession,
  nextLobbyRound,
  startLobbySession,
} from '../../api/lobbyApi'
import { ApiError } from '../../api/client'
import type { LobbySnapshot } from '../../types/lobby'
import { useCallback, useState } from 'react'
import { t } from '../../lib/i18n'

interface AdminControlsProps {
  lobbyId: string
  snapshot: LobbySnapshot
}

export function AdminControls({ lobbyId, snapshot }: AdminControlsProps) {
  const [loadingAction, setLoadingAction] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  const runAction = useCallback(
    async (actionId: string, action: () => Promise<unknown>) => {
      setLoadingAction(actionId)
      setError(null)
      try {
        await action()
      } catch (err) {
        const message =
          err instanceof ApiError ? err.message : t('admin.actionFailed')
        setError(message)
      } finally {
        setLoadingAction(null)
      }
    },
    [],
  )

  const hasControls =
    snapshot.phase === 'READY' || snapshot.phase === 'BETWEEN_ROUNDS'

  if (!hasControls) {
    return null
  }

  return (
    <section className="admin-controls" aria-label="Admin controls">
      {error ? (
        <p className="admin-controls__error" role="alert">
          {error}
        </p>
      ) : null}

      {snapshot.phase === 'READY' ? (
        <button
          type="button"
          className="button button--primary"
          disabled={loadingAction !== null}
          onClick={() => void runAction('start', () => startLobbySession(lobbyId))}
        >
          {loadingAction === 'start' ? t('admin.starting') : t('admin.startSession')}
        </button>
      ) : null}

      {snapshot.phase === 'BETWEEN_ROUNDS' ? (
        <div className="admin-controls__row">
          <button
            type="button"
            className="button button--primary"
            disabled={loadingAction !== null}
            onClick={() => void runAction('next', () => nextLobbyRound(lobbyId))}
          >
            {loadingAction === 'next' ? t('admin.loading') : t('admin.nextPhoto')}
          </button>
          <button
            type="button"
            className="button button--secondary"
            disabled={loadingAction !== null}
            onClick={() =>
              void runAction('finish', () => finishLobbySession(lobbyId))
            }
          >
            {loadingAction === 'finish' ? t('admin.ending') : t('admin.endSession')}
          </button>
        </div>
      ) : null}
    </section>
  )
}
