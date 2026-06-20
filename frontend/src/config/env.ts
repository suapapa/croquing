const DEFAULT_API_BASE = 'http://localhost:8080'

/** Backend HTTP origin from `VITE_API_BASE` (falls back to localhost:8080). */
export function getApiBase(): string {
  const base = import.meta.env.VITE_API_BASE?.trim()
  return base || DEFAULT_API_BASE
}
