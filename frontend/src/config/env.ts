/** Backend HTTP origin from `VITE_API_BASE` (empty = same origin, via Vite proxy in dev). */
export function getApiBase(): string {
  const base = import.meta.env.VITE_API_BASE?.trim()
  return base ?? ''
}

/** WebSocket origin derived from the API base URL or the current page origin. */
export function getWsBase(): string {
  const api = getApiBase()
  if (!api) {
    if (typeof window === 'undefined') {
      return 'ws://localhost:8080'
    }
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    return `${protocol}//${window.location.host}`
  }
  if (api.startsWith('https://')) {
    return `wss://${api.slice('https://'.length)}`
  }
  if (api.startsWith('http://')) {
    return `ws://${api.slice('http://'.length)}`
  }
  return api.replace(/^http/, 'ws')
}

/** Build the lobby WebSocket URL for a given lobby id. */
export function getLobbyWsUrl(lobbyId: string): string {
  return `${getWsBase()}/ws/lobby/${encodeURIComponent(lobbyId)}`
}
