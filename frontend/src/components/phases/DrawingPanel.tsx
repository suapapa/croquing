import {
  formatRemainingTime,
  getDrawRemainingMs,
  useServerClock,
} from '../../hooks/useServerClock'
import type { LobbySnapshot } from '../../types/lobby'

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
  const remainingMs = getDrawRemainingMs(snapshot.draw_ends_at, serverNow)
  const progress = Math.min(
    1,
    Math.max(0, 1 - remainingMs / drawDurationMs),
  )
  const photo = snapshot.current_photo

  return (
    <section className="drawing-panel" aria-live="polite">
      <div className="drawing-panel__timer">
        <div
          className="drawing-panel__timer-bar"
          style={{ transform: `scaleX(${progress})` }}
          role="progressbar"
          aria-valuemin={0}
          aria-valuemax={drawDurationMs}
          aria-valuenow={drawDurationMs - remainingMs}
          aria-label="Draw time remaining"
        />
        <span className="drawing-panel__timer-label">
          {formatRemainingTime(remainingMs)}
        </span>
      </div>

      <div className="drawing-panel__stage">
        {photo ? (
          <>
            <img
              className="drawing-panel__photo"
              src={photo.large_image_url}
              alt="Reference photo for this croquis round"
              width={photo.width}
              height={photo.height}
            />
            <p className="drawing-panel__attribution">
              Image from{' '}
              <a href={photo.page_url} target="_blank" rel="noreferrer">
                PixaBay
              </a>
            </p>
          </>
        ) : (
          <p className="drawing-panel__empty">Waiting for photo…</p>
        )}
      </div>

      {snapshot.total_rounds > 0 ? (
        <p className="drawing-panel__round">
          Round {snapshot.current_round} / {snapshot.total_rounds}
        </p>
      ) : null}
    </section>
  )
}
