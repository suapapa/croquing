import { useEffect, useState } from 'react'
import { fetchAppConfig } from '../api/configApi'

const DEFAULT_APP_LOGO_LINK = 'https://homin.dev'

let cachedAppLogoLink: string | null = null
let loadPromise: Promise<string> | null = null

function loadAppLogoLink(): Promise<string> {
  if (cachedAppLogoLink) {
    return Promise.resolve(cachedAppLogoLink)
  }
  if (!loadPromise) {
    loadPromise = fetchAppConfig()
      .then((config) => config.app_logo_link.trim() || DEFAULT_APP_LOGO_LINK)
      .catch(() => DEFAULT_APP_LOGO_LINK)
      .then((link) => {
        cachedAppLogoLink = link
        return link
      })
  }
  return loadPromise
}

/** App logo link from server `APP_LOGO_LINK` (falls back to "https://homin.dev"). */
export function useAppLogoLink(): string {
  const [appLogoLink, setAppLogoLink] = useState(
    cachedAppLogoLink ?? DEFAULT_APP_LOGO_LINK,
  )

  useEffect(() => {
    void loadAppLogoLink().then(setAppLogoLink)
  }, [])

  return appLogoLink
}
