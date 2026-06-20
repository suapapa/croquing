import { apiRequest } from './client'

export interface PixabaySearchHit {
  pixabay_id: number
  page_url: string
  preview_url: string
  webformat_url: string
  large_image_url: string
  width: number
  height: number
  views: number
  downloads: number
  likes: number
}

export interface PixabaySearchResponse {
  total: number
  total_hits: number
  hits: PixabaySearchHit[]
  rate_limit: {
    limit: number
    remaining: number
    reset: number
  }
}

export interface PixabaySearchParams {
  lobbyId: string
  query: string
  order: 'popular' | 'latest'
  page: number
  perPage?: number
}

export function searchPixabay({
  lobbyId,
  query,
  order,
  page,
  perPage = 20,
}: PixabaySearchParams): Promise<PixabaySearchResponse> {
  const params = new URLSearchParams({
    q: query,
    order,
    page: String(page),
    per_page: String(perPage),
    lobby_id: lobbyId,
  })

  return apiRequest<PixabaySearchResponse>(`/api/pixabay/search?${params}`, {
    lobbyId,
  })
}

export function hitToPhoto(hit: PixabaySearchHit) {
  return {
    pixabay_id: hit.pixabay_id,
    preview_url: hit.preview_url,
    large_image_url: hit.large_image_url,
    page_url: hit.page_url,
    width: hit.width,
    height: hit.height,
  }
}
