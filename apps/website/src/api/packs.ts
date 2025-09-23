import { API_URL } from '@/api/config'
import type { PackResponse } from '@/types/pack'

export interface PacksResponse {
  packs: PackResponse[]
  total: number
}

async function fetchPacksEndpoint(endpoint: string, page = 0, pageSize = 10) {
  const query = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  })

  const res = await fetch(`${API_URL}/${endpoint}?${query}`, {
    credentials: 'include',
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch packs: ${res.status}`)
  }

  const json: PacksResponse = await res.json()
  return json
}


export async function fetchUserPacks(page = 0, pageSize = 10) {
  return fetchPacksEndpoint('user/packs', page, pageSize)
}

export async function fetchPublicPacks(page = 0, pageSize = 10) {
  return fetchPacksEndpoint('public/packs', page, pageSize)
}
