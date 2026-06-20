import { useEffect, useRef, useState } from 'react'
import { getLobbyWsUrl } from '../config/env'
import type { LobbySnapshot, WsEnvelope } from '../types/lobby'

export type ConnectionStatus =
  | 'connecting'
  | 'connected'
  | 'reconnecting'
  | 'disconnected'

export interface UseLobbySocketResult {
  snapshot: LobbySnapshot | null
  status: ConnectionStatus
  serverTimeOffsetMs: number
  error: string | null
}

const INITIAL_RETRY_MS = 500
const MAX_RETRY_MS = 8000

function computeServerOffset(snapshot: LobbySnapshot): number {
  const serverMs = Date.parse(snapshot.server_time)
  if (Number.isNaN(serverMs)) {
    return 0
  }
  return serverMs - Date.now()
}

export function useLobbySocket(
  lobbyId: string | undefined,
): UseLobbySocketResult {
  const [snapshot, setSnapshot] = useState<LobbySnapshot | null>(null)
  const [status, setStatus] = useState<ConnectionStatus>('connecting')
  const [serverTimeOffsetMs, setServerTimeOffsetMs] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const retryDelayRef = useRef(INITIAL_RETRY_MS)
  const reconnectTimerRef = useRef<number | null>(null)

  useEffect(() => {
    if (!lobbyId) {
      return
    }

    let active = true
    let socket: WebSocket | null = null

    const clearReconnectTimer = () => {
      if (reconnectTimerRef.current !== null) {
        window.clearTimeout(reconnectTimerRef.current)
        reconnectTimerRef.current = null
      }
    }

    const scheduleReconnect = () => {
      if (!active) {
        return
      }

      setStatus('reconnecting')
      clearReconnectTimer()
      reconnectTimerRef.current = window.setTimeout(() => {
        connect()
      }, retryDelayRef.current)
      retryDelayRef.current = Math.min(retryDelayRef.current * 2, MAX_RETRY_MS)
    }

    const connect = () => {
      if (!active) {
        return
      }

      setStatus((current) =>
        current === 'reconnecting' ? 'reconnecting' : 'connecting',
      )
      setError(null)

      socket = new WebSocket(getLobbyWsUrl(lobbyId))

      socket.addEventListener('open', () => {
        if (!active) {
          return
        }
        retryDelayRef.current = INITIAL_RETRY_MS
        setStatus('connected')
      })

      socket.addEventListener('message', (event) => {
        if (!active) {
          return
        }

        try {
          const envelope = JSON.parse(String(event.data)) as WsEnvelope
          if (envelope.type !== 'snapshot' || !envelope.payload) {
            return
          }

          setSnapshot(envelope.payload)
          setServerTimeOffsetMs(computeServerOffset(envelope.payload))
          setError(null)
        } catch {
          setError('Received invalid lobby update')
        }
      })

      socket.addEventListener('close', () => {
        if (!active) {
          return
        }
        scheduleReconnect()
      })

      socket.addEventListener('error', () => {
        if (!active) {
          return
        }
        setError('Connection error')
      })
    }

    connect()

    return () => {
      active = false
      clearReconnectTimer()
      socket?.close()
    }
  }, [lobbyId])

  return { snapshot, status, serverTimeOffsetMs, error }
}
