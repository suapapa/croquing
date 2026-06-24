import { apiRequest } from './client'

export interface AppConfig {
  app_name: string
}

export function fetchAppConfig(): Promise<AppConfig> {
  return apiRequest<AppConfig>('/api/config')
}
