export type LobbyPhase =
  | 'WAITING'
  | 'SELECTING'
  | 'READY'
  | 'DRAWING'
  | 'BETWEEN_ROUNDS'
  | 'FINISHED'

export interface Photo {
  pixabay_id: number
  preview_url: string
  large_image_url: string
  page_url: string
  width: number
  height: number
}

export interface LobbySnapshot {
  id: string
  phase: LobbyPhase
  participant_count: number
  selected_count: number
  current_round: number
  total_rounds: number
  draw_ends_at?: string
  current_photo?: Photo
  server_time: string
}

export interface CreateLobbyResponse {
  id: string
  admin_token: string
  join_url: string
}

export interface WsEnvelope {
  type: string
  payload: LobbySnapshot
}
