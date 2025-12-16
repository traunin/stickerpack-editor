import { API_URL } from '@/api/config'
import type { Pack, PackPreview, TelegramSticker, TelegramStickerSet } from '@/types/pack'
import type { Source, Sticker } from '@/types/sticker'
import { splitEmojis } from './emoji'

export interface PacksResponse {
  packs: PackPreview[]
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

export async function fetchPack(name: string) {
  const res = await fetch(`${API_URL}/user/packs/${name}`, {
    credentials: 'include',
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch pack: ${res.status}`)
  }

  const json: TelegramStickerSet = await res.json()
  return stickerSetToPack(json)
}

function stickerSetToPack(set: TelegramStickerSet): Pack {
  const stickers: Sticker[] = set.stickers.map((s: TelegramSticker) => {
    const source: Source = 'telegram'
    return {
      id: s.file_id,
      name: s.emoji,
      preview: `${API_URL}/media?file_id=${s.file_id}`,
      full: `${API_URL}/media?file_id=${s.file_id}`,
      uuid: crypto.randomUUID().toString(),
      emoji_list: splitEmojis(s.emoji),
      source,
    }
  })

  return {
    title: set.title,
    isPublic: set.is_public,
    stickers,
  }
}
