import { useCallback, useEffect, useRef, useState } from 'react'
import { setLobbyDrawDuration } from '../api/lobbyApi'
import { ApiError } from '../api/client'
import { t } from '../lib/i18n'

const MIN_MINUTES = 1
const MAX_MINUTES = 60

interface RoundDurationPickerProps {
  lobbyId: string
  minutes: number
  isAdmin: boolean
}

export function RoundDurationPicker({
  lobbyId,
  minutes,
  isAdmin,
}: RoundDurationPickerProps) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const groupRef = useRef<HTMLDivElement>(null)

  const applyMinutes = useCallback(
    async (next: number) => {
      if (next < MIN_MINUTES || next > MAX_MINUTES || next === minutes) {
        return
      }

      setLoading(true)
      setError(null)
      try {
        await setLobbyDrawDuration(lobbyId, next)
      } catch (err) {
        const message =
          err instanceof ApiError ? err.message : t('duration.updateFailed')
        setError(message)
      } finally {
        setLoading(false)
      }
    },
    [lobbyId, minutes],
  )

  useEffect(() => {
    if (!isAdmin) {
      return
    }

    const el = groupRef.current
    if (!el) {
      return
    }

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'ArrowLeft') {
        event.preventDefault()
        void applyMinutes(minutes - 1)
      } else if (event.key === 'ArrowRight') {
        event.preventDefault()
        void applyMinutes(minutes + 1)
      }
    }

    el.addEventListener('keydown', onKeyDown)
    return () => {
      el.removeEventListener('keydown', onKeyDown)
    }
  }, [isAdmin, minutes, applyMinutes])

  return (
    <div className="round-duration">
      <p className="round-duration__label">{t('duration.label')}</p>
      {isAdmin ? (
        <div
          ref={groupRef}
          className="round-duration__control"
          role="group"
          aria-label={t('duration.ariaGroup')}
          tabIndex={0}
        >
          <button
            type="button"
            className="button button--secondary button--icon-only round-duration__step"
            disabled={loading || minutes <= MIN_MINUTES}
            onClick={() => void applyMinutes(minutes - 1)}
            aria-label={t('duration.decrease')}
          >
            &lt;
          </button>
          <span className="round-duration__value" aria-live="polite">
            {minutes}
          </span>
          <button
            type="button"
            className="button button--secondary button--icon-only round-duration__step"
            disabled={loading || minutes >= MAX_MINUTES}
            onClick={() => void applyMinutes(minutes + 1)}
            aria-label={t('duration.increase')}
          >
            &gt;
          </button>
        </div>
      ) : (
        <p className="round-duration__value round-duration__value--readonly">
          {t('duration.minutes', { count: minutes })}
        </p>
      )}
      {error ? (
        <p className="round-duration__error" role="alert">
          {error}
        </p>
      ) : null}
      {isAdmin ? (
        <p className="round-duration__hint">{t('duration.hint')}</p>
      ) : null}
    </div>
  )
}
