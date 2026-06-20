import { useCallback, useEffect, useState } from 'react'
import { copyLobbyJoinUrl } from '../../lib/lobbyLink'

interface CopyLobbyLinkButtonProps {
  lobbyId: string
  className?: string
  compact?: boolean
}

export function CopyLobbyLinkButton({
  lobbyId,
  className = '',
  compact = false,
}: CopyLobbyLinkButtonProps) {
  const [copied, setCopied] = useState(false)
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    if (!copied && !failed) {
      return
    }

    const timer = window.setTimeout(() => {
      setCopied(false)
      setFailed(false)
    }, 2000)

    return () => window.clearTimeout(timer)
  }, [copied, failed])

  const handleCopy = useCallback(async () => {
    const ok = await copyLobbyJoinUrl(lobbyId)
    if (ok) {
      setCopied(true)
      setFailed(false)
      return
    }

    setFailed(true)
    setCopied(false)
  }, [lobbyId])

  const label = copied ? 'Copied!' : failed ? 'Copy failed' : 'Copy link'
  const classes = [
    'button',
    'button--secondary',
    compact ? 'lobby-layout__copy-link' : '',
    className,
  ]
    .filter(Boolean)
    .join(' ')

  return (
    <button type="button" className={classes} onClick={() => void handleCopy()}>
      {label}
    </button>
  )
}
