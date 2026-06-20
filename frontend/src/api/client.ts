import { getApiBase } from '../config/env'
import { getAdminToken } from '../lib/adminStorage'

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

interface RequestOptions extends Omit<RequestInit, 'body'> {
  lobbyId?: string
  body?: unknown
}

export async function apiRequest<T>(
  path: string,
  options: RequestOptions = {},
): Promise<T> {
  const { lobbyId, body, headers, ...rest } = options
  const requestHeaders = new Headers(headers)

  if (body !== undefined) {
    requestHeaders.set('Content-Type', 'application/json')
  }

  if (lobbyId) {
    const token = getAdminToken(lobbyId)
    if (token) {
      requestHeaders.set('X-Admin-Token', token)
    }
  }

  const response = await fetch(`${getApiBase()}${path}`, {
    ...rest,
    headers: requestHeaders,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })

  if (!response.ok) {
    let message = response.statusText
    try {
      const payload = (await response.json()) as { error?: string }
      if (payload.error) {
        message = payload.error
      }
    } catch {
      // ignore parse errors
    }
    throw new ApiError(message, response.status)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}
