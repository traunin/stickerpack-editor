import type { Emote } from '@/types/sticker'

interface TenorSearchResponse {
  results: {
    id: string
    title: string
    content_description: string
    media_formats: {
      tinygif?: { url: string }
      gif: { url: string }
      mediumgif?: { url: string }
      nanogif?: { url: string }
    }
  }[]
  next: string
}

const TENOR_API_KEY = import.meta.env.VITE_TENOR_API_KEY

export async function searchTenor(query: string, pos = '', limit = 10) {
  if (!query) {
    return { items: [], next: '' }
  }
  const url = new URL('https://tenor.googleapis.com/v2/search')
  url.searchParams.set('q', query)
  url.searchParams.set('key', TENOR_API_KEY)
  url.searchParams.set('limit', String(limit))
  if (pos)
    url.searchParams.set('pos', pos)

  const res = await fetch(url.toString())
  if (!res.ok) {
    throw new Error(`Failed to fetch from Tenor: ${res.status}`)
  }

  const json: TenorSearchResponse = await res.json()
  const items = json.results ?? []

  return {
    items: items.map((e): Emote => ({
      id: e.id,
      name: e.title || e.content_description,
      preview: e.media_formats.tinygif?.url ?? e.media_formats.nanogif?.url ?? '',
      full: e.media_formats.gif?.url ?? e.media_formats.mediumgif?.url ?? '',
    })),
    next: json.next ?? '',
  }
}
