import { apiRequest } from './client'

export interface AppConfig {
  app_name: string
  app_logo: string
  app_logo_link: string
}

export function fetchAppConfig(): Promise<AppConfig> {
  return apiRequest<AppConfig>('/api/config')
}
