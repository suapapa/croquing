import { useEffect, useState } from 'react'

/** Returns server-adjusted current time, refreshed every second. */
export function useServerClock(serverTimeOffsetMs: number): number {
  const [now, setNow] = useState(() => Date.now() + serverTimeOffsetMs)

  useEffect(() => {
    setNow(Date.now() + serverTimeOffsetMs)
    const timer = window.setInterval(() => {
      setNow(Date.now() + serverTimeOffsetMs)
    }, 1000)
    return () => window.clearInterval(timer)
  }, [serverTimeOffsetMs])

  return now
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
