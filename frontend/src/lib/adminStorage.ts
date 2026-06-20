const STORAGE_PREFIX = 'croquis-admin:'

/** Persist admin token for a lobby (session-scoped). */
export function saveAdminToken(lobbyId: string, token: string): void {
  sessionStorage.setItem(`${STORAGE_PREFIX}${lobbyId}`, token)
}

/** Read stored admin token, if any. */
export function getAdminToken(lobbyId: string): string | null {
  return sessionStorage.getItem(`${STORAGE_PREFIX}${lobbyId}`)
}

/** Whether this browser session has admin credentials for the lobby. */
export function isLobbyAdmin(lobbyId: string): boolean {
  return getAdminToken(lobbyId) !== null
}
