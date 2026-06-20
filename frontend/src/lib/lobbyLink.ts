export function getLobbyJoinUrl(lobbyId: string): string {
  return `${window.location.origin}/lobby/${encodeURIComponent(lobbyId)}`
}

export async function copyLobbyJoinUrl(lobbyId: string): Promise<boolean> {
  return copyTextToClipboard(getLobbyJoinUrl(lobbyId))
}

async function copyTextToClipboard(text: string): Promise<boolean> {
  if (navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(text)
      return true
    } catch {
      // Fall back to execCommand below.
    }
  }

  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', '')
  textarea.style.position = 'fixed'
  textarea.style.left = '-9999px'
  document.body.appendChild(textarea)
  textarea.select()

  let copied = false
  try {
    const success = document.execCommand('copy')
    if (success) {
      copied = true
    }
  } catch {
    // Keep copied as false
  } finally {
    document.body.removeChild(textarea)
  }

  return copied
}
