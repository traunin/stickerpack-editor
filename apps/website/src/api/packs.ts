import { API_URL } from './config'

export interface PublicPack {
  id: number
  title: string
  name: string
  thumbnail_id: string
}

interface PacksResponse {
  packs: PublicPack[]
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
