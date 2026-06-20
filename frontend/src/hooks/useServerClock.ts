import { useSyncExternalStore } from 'react'

function subscribe(onStoreChange: () => void) {
  const timer = window.setInterval(onStoreChange, 1000)
  return () => window.clearInterval(timer)
}

/** Returns server-adjusted current time, refreshed every second. */
export function useServerClock(serverTimeOffsetMs: number): number {
  return useSyncExternalStore(
    subscribe,
    () => Date.now() + serverTimeOffsetMs,
    () => Date.now() + serverTimeOffsetMs,
  )
}

/** Remaining milliseconds until drawEndsAt, using server-adjusted clock. */
export function getDrawRemainingMs(
  drawEndsAt: string | undefined,
  serverNowMs: number,
): number {
  if (!drawEndsAt) {
    return 0
  }
  const endsAt = Date.parse(drawEndsAt)
  if (Number.isNaN(endsAt)) {
    return 0
  }
  return Math.max(0, endsAt - serverNowMs)
}

export function formatRemainingTime(remainingMs: number): string {
  const totalSeconds = Math.ceil(remainingMs / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}
