import { apiRequest } from './client'
import type { CreateLobbyResponse, LobbySnapshot, Photo } from '../types/lobby'

export function createLobby(): Promise<CreateLobbyResponse> {
  return apiRequest<CreateLobbyResponse>('/api/lobbies', { method: 'POST' })
}

export function getLobby(lobbyId: string): Promise<LobbySnapshot> {
  return apiRequest<LobbySnapshot>(
    `/api/lobbies/${encodeURIComponent(lobbyId)}`,
  )
}

export function setLobbyPhotos(
  lobbyId: string,
  photos: Photo[],
): Promise<LobbySnapshot> {
  return apiRequest<LobbySnapshot>(
    `/api/lobbies/${encodeURIComponent(lobbyId)}/photos`,
    {
      method: 'PUT',
      lobbyId,
      body: { photos },
    },
  )
}

export function reopenLobbyPhotos(lobbyId: string): Promise<LobbySnapshot> {
  return apiRequest<LobbySnapshot>(
    `/api/lobbies/${encodeURIComponent(lobbyId)}/photos/reopen`,
    { method: 'POST', lobbyId },
  )
}

export function markLobbyReady(lobbyId: string): Promise<LobbySnapshot> {
  return apiRequest<LobbySnapshot>(
    `/api/lobbies/${encodeURIComponent(lobbyId)}/ready`,
    { method: 'POST', lobbyId },
  )
}

export function startLobbySession(lobbyId: string): Promise<LobbySnapshot> {
  return apiRequest<LobbySnapshot>(
    `/api/lobbies/${encodeURIComponent(lobbyId)}/start`,
    { method: 'POST', lobbyId },
  )
}

export function nextLobbyRound(lobbyId: string): Promise<LobbySnapshot> {
  return apiRequest<LobbySnapshot>(
    `/api/lobbies/${encodeURIComponent(lobbyId)}/next`,
    { method: 'POST', lobbyId },
  )
}

export function finishLobbySession(lobbyId: string): Promise<LobbySnapshot> {
  return apiRequest<LobbySnapshot>(
    `/api/lobbies/${encodeURIComponent(lobbyId)}/finish`,
    { method: 'POST', lobbyId },
  )
}

export function setLobbyDrawDuration(
  lobbyId: string,
  minutes: number,
): Promise<LobbySnapshot> {
  return apiRequest<LobbySnapshot>(
    `/api/lobbies/${encodeURIComponent(lobbyId)}/draw-duration`,
    {
      method: 'PUT',
      lobbyId,
      body: { minutes },
    },
  )
}
