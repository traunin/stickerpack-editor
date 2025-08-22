import { API_URL } from './config'
import type { PackResponse } from '@/types/pack'

interface PacksResponse {
  packs: PackResponse[]
  total: number
}

export async function fetchPacksEndpoint(endpoint: string, page = 0, pageSize = 10) {
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
  return { packs: json.packs, total: json.total }
}
