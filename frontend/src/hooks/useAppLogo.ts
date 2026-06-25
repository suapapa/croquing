import { useEffect, useState } from 'react'
import { fetchAppConfig } from '../api/configApi'

const DEFAULT_APP_LOGO = '/example_logo.png'

let cachedAppLogo: string | null = null
let loadPromise: Promise<string> | null = null

function loadAppLogo(): Promise<string> {
  if (cachedAppLogo) {
    return Promise.resolve(cachedAppLogo)
  }
  if (!loadPromise) {
    loadPromise = fetchAppConfig()
      .then((config) => config.app_logo.trim() || DEFAULT_APP_LOGO)
      .catch(() => DEFAULT_APP_LOGO)
      .then((logo) => {
        cachedAppLogo = logo
        return logo
      })
  }
  return loadPromise
}

/** App logo URL from server `APP_LOGO` (falls back to "/example_logo.png"). */
export function useAppLogo(): string {
  const [appLogo, setAppLogo] = useState(cachedAppLogo ?? DEFAULT_APP_LOGO)

  useEffect(() => {
    void loadAppLogo().then(setAppLogo)
  }, [])

  return appLogo
}
