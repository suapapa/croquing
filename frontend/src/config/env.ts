const DEFAULT_API_BASE = 'http://localhost:8080'

/** Backend HTTP origin from `VITE_API_BASE` (falls back to localhost:8080). */
export function getApiBase(): string {
  const base = import.meta.env.VITE_API_BASE?.trim()
  return base || DEFAULT_API_BASE
}

/** WebSocket origin derived from the API base URL. */
export function getWsBase(): string {
  const api = getApiBase()
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
