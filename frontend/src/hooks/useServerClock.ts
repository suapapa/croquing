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

/** Must match lobby.RoundCountdown on the server. */
export const ROUND_COUNTDOWN_MS = 5_000

const ONE_MINUTE_MS = 60_000
const TEN_SECONDS_MS = 10_000

export interface DrawPhaseState {
  isCountdown: boolean
  countdownSeconds: number
  drawRemainingMs: number
}

/** Splits total remaining time into pre-draw countdown vs active draw timer. */
export function getDrawPhaseState(
  drawEndsAt: string | undefined,
  serverNowMs: number,
  drawDurationMs: number,
): DrawPhaseState {
  const totalRemainingMs = getDrawRemainingMs(drawEndsAt, serverNowMs)
  const countdownRemainingMs = totalRemainingMs - drawDurationMs

  if (countdownRemainingMs > 0) {
    return {
      isCountdown: true,
      countdownSeconds: Math.ceil(countdownRemainingMs / 1000),
      drawRemainingMs: drawDurationMs,
    }
  }

  return {
    isCountdown: false,
    countdownSeconds: 0,
    drawRemainingMs: totalRemainingMs,
  }
}

export function isDrawUrgent(drawRemainingMs: number): boolean {
  return drawRemainingMs > 0 && drawRemainingMs <= ONE_MINUTE_MS
}

export function isDrawCritical(drawRemainingMs: number): boolean {
  return drawRemainingMs > 0 && drawRemainingMs <= TEN_SECONDS_MS
}
