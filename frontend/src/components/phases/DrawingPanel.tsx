import { useEffect, useState } from 'react'
import {
  formatRemainingTime,
  getDrawPhaseState,
  isDrawCritical,
  isDrawUrgent,
  useServerClock,
} from '../../hooks/useServerClock'
import type { LobbySnapshot } from '../../types/lobby'
import { IconExitFullscreen, IconFullscreen } from '../ui/Icons'

interface DrawingPanelProps {
  snapshot: LobbySnapshot
  serverTimeOffsetMs: number
  drawDurationMs?: number
}

const DEFAULT_DRAW_MS = 5 * 60 * 1000

export function DrawingPanel({
  snapshot,
  serverTimeOffsetMs,
  drawDurationMs = DEFAULT_DRAW_MS,
}: DrawingPanelProps) {
  const serverNow = useServerClock(serverTimeOffsetMs)
  const [isFullscreen, setIsFullscreen] = useState(false)

  useEffect(() => {
    function handleFullscreenChange() {
      setIsFullscreen(!!document.fullscreenElement)
    }
    document.addEventListener('fullscreenchange', handleFullscreenChange)
    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange)
    }
  }, [])

  function toggleFullscreen() {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().catch((err) => {
        console.error('Failed to enter fullscreen:', err)
      })
    } else {
      document.exitFullscreen().catch((err) => {
        console.error('Failed to exit fullscreen:', err)
      })
    }
  }
  const { isCountdown, countdownSeconds, drawRemainingMs } = getDrawPhaseState(
    snapshot.draw_ends_at,
    serverNow,
    drawDurationMs,
  )
  const progress = Math.min(
    1,
    Math.max(0, 1 - drawRemainingMs / drawDurationMs),
  )
  const photo = snapshot.current_photo
  const urgent = !isCountdown && isDrawUrgent(drawRemainingMs)
  const critical = !isCountdown && isDrawCritical(drawRemainingMs)

  const photoWrapClass = [
    'drawing-panel__photo-wrap',
    urgent ? 'drawing-panel__photo-wrap--urgent' : '',
    critical ? 'drawing-panel__photo-wrap--urgent-blink' : '',
  ]
    .filter(Boolean)
    .join(' ')

  return (
    <section className="drawing-panel" aria-live="polite">
      {isCountdown ? (
        <div
          className="drawing-panel__countdown"
          role="timer"
          aria-label={`Round starts in ${countdownSeconds} seconds`}
        >
          <span className="drawing-panel__countdown-number" aria-hidden="true">
            {countdownSeconds}
          </span>
        </div>
      ) : (
        <div className="drawing-panel__timer">
          <div
            className={`drawing-panel__timer-bar${
              urgent || critical ? ' drawing-panel__timer-bar--urgent' : ''
            }`}
            style={{ transform: `scaleX(${progress})` }}
            role="progressbar"
            aria-valuemin={0}
            aria-valuemax={drawDurationMs}
            aria-valuenow={drawDurationMs - drawRemainingMs}
            aria-label="Draw time remaining"
          />
          <span className="drawing-panel__timer-label">
            {formatRemainingTime(drawRemainingMs)}
          </span>
        </div>
      )}

      <div className="drawing-panel__body">
        <div className="drawing-panel__stage">
          <div className="drawing-panel__controls">
            <button
              type="button"
              className="button--zen-control"
              onClick={toggleFullscreen}
              title={isFullscreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
              aria-label={isFullscreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
            >
              {isFullscreen ? (
                <IconExitFullscreen style={{ width: '1.25rem', height: '1.25rem' }} />
              ) : (
                <IconFullscreen style={{ width: '1.25rem', height: '1.25rem' }} />
              )}
            </button>
          </div>
          {photo && !isCountdown ? (
            <div className={photoWrapClass}>
              <img
                className="drawing-panel__photo"
                src={photo.large_image_url}
                alt="Reference photo for this croquis round"
                width={photo.width}
                height={photo.height}
              />
            </div>
          ) : !isCountdown ? (
            <p className="drawing-panel__empty">Waiting for photo…</p>
          ) : null}
        </div>

        {photo && !isCountdown ? (
          <footer className="drawing-panel__footer">
            <p className="drawing-panel__attribution">
              Image from{' '}
              <a href={photo.page_url} target="_blank" rel="noreferrer">
                Pixabay
              </a>
            </p>
            {snapshot.total_rounds > 0 ? (
              <p className="drawing-panel__round">
                Round {snapshot.current_round} / {snapshot.total_rounds}
              </p>
            ) : null}
          </footer>
        ) : null}
      </div>
    </section>
  )
}
