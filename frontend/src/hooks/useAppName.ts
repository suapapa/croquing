import { useEffect, useState } from 'react'
import { fetchAppConfig } from '../api/configApi'

const DEFAULT_APP_NAME = ''

let cachedAppName: string | null = null
let loadPromise: Promise<string> | null = null

function loadAppName(): Promise<string> {
  if (cachedAppName) {
    return Promise.resolve(cachedAppName)
  }
  if (!loadPromise) {
    loadPromise = fetchAppConfig()
      .then((config) => config.app_name.trim() || DEFAULT_APP_NAME)
      .catch(() => DEFAULT_APP_NAME)
      .then((name) => {
        cachedAppName = name
        return name
      })
  }
  return loadPromise
}

/** App display name from server `APP_NAME` (falls back to ""). */
export function useAppName(): string {
  const [appName, setAppName] = useState(cachedAppName ?? DEFAULT_APP_NAME)

  useEffect(() => {
    void loadAppName().then(setAppName)
  }, [])

  return appName
}
