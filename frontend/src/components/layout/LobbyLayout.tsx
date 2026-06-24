import { Link } from 'react-router-dom'
import type { ReactNode } from 'react'
import type { LobbySnapshot } from '../../types/lobby'
import type { ConnectionStatus } from '../../hooks/useLobbySocket'
import { getConnectionLabel, getPhaseMessage, t } from '../../lib/i18n'
import { CopyLobbyLinkButton } from '../lobby/CopyLobbyLinkButton'
import { IconLogo } from '../ui/Icons'

interface LobbyLayoutProps {
  lobbyId: string
  isAdmin: boolean
  snapshot: LobbySnapshot | null
  connectionStatus: ConnectionStatus
  connectionError: string | null
  children?: ReactNode
}

export function LobbyLayout({
  lobbyId,
  isAdmin,
  snapshot,
  connectionStatus,
  connectionError,
  children,
}: LobbyLayoutProps) {
  const phaseMessage = snapshot ? getPhaseMessage(snapshot.phase) : null
  const isLive = connectionStatus === 'connected'
  const isDisconnected =
    connectionStatus === 'disconnected' || connectionStatus === 'reconnecting'

  return (
    <div className="lobby-layout">
      <header className="lobby-layout__header">
        <div className="lobby-layout__brand-row">
          <Link to="/" className="lobby-layout__brand">
            <IconLogo className="lobby-layout__logo" aria-hidden="true" />
            <span>Croquing</span>
          </Link>
          <span
            className={`lobby-layout__badge lobby-layout__badge--${
              isAdmin ? 'admin' : 'participant'
            }`}
          >
            {isAdmin ? t('lobby.badge.admin') : t('lobby.badge.participant')}
          </span>
        </div>

        <div className="lobby-layout__meta">
          <span
            className={`lobby-layout__connection lobby-layout__connection--${
              isLive ? 'live' : 'offline'
            }`}
            aria-live="polite"
          >
            <span className="lobby-layout__connection-dot" aria-hidden="true" />
            {getConnectionLabel(connectionStatus)}
          </span>
          {snapshot ? (
            <span className="lobby-layout__participants">
              {snapshot.participant_count === 1
                ? t('lobby.participantCount', {
                    count: snapshot.participant_count,
                  })
                : t('lobby.participantCountPlural', {
                    count: snapshot.participant_count,
                  })}
            </span>
          ) : null}
          {snapshot ? <CopyLobbyLinkButton lobbyId={lobbyId} compact /> : null}
        </div>
      </header>

      {isDisconnected ? (
        <div
          className="lobby-layout__banner lobby-layout__banner--warning"
          role="status"
        >
          {t('lobby.connection.lost')}
        </div>
      ) : null}

      {connectionError ? (
        <div
          className="lobby-layout__banner lobby-layout__banner--error"
          role="alert"
        >
          {connectionError}
        </div>
      ) : null}

      {!snapshot && connectionStatus !== 'disconnected' ? (
        <div className="lobby-layout__loading" role="status">
          <div className="lobby-layout__spinner" aria-hidden="true" />
          <p>{t('lobby.loadingState')}</p>
        </div>
      ) : null}

      {snapshot && phaseMessage ? (
        <section className="lobby-layout__intro" aria-live="polite">
          <p className="lobby-layout__lobby-id">
            {t('lobby.lobbyId', { id: lobbyId.slice(0, 8) })}
          </p>
          <h1>{phaseMessage.title}</h1>
          <p>{phaseMessage.description}</p>
        </section>
      ) : null}

      {children}
    </div>
  )
}
