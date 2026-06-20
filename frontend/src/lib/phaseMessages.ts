import type { LobbyPhase } from '../types/lobby'

export interface PhaseMessage {
  title: string
  description: string
}

const PHASE_MESSAGES: Record<LobbyPhase, PhaseMessage> = {
  WAITING: {
    title: 'Waiting for the host',
    description: 'The admin is getting things ready. Stay on this page.',
  },
  SELECTING: {
    title: 'Choosing reference photos',
    description: 'The admin is searching PixaBay and picking photos for the session.',
  },
  READY: {
    title: 'Ready to start',
    description: 'Photos are shuffled and hidden. Waiting for the admin to begin.',
  },
  DRAWING: {
    title: 'Draw time',
    description: 'Focus on the reference photo. The timer is server-controlled.',
  },
  BETWEEN_ROUNDS: {
    title: 'Round break',
    description: 'Take a short break before the next pose.',
  },
  FINISHED: {
    title: 'Session complete',
    description: 'Thanks for drawing together. See you next week.',
  },
}

export function getPhaseMessage(phase: LobbyPhase): PhaseMessage {
  return PHASE_MESSAGES[phase]
}

export function getConnectionLabel(
  status: 'connecting' | 'connected' | 'reconnecting' | 'disconnected',
): string {
  switch (status) {
    case 'connecting':
      return 'Connecting…'
    case 'connected':
      return 'Live'
    case 'reconnecting':
      return 'Reconnecting…'
    case 'disconnected':
      return 'Disconnected'
  }
}
