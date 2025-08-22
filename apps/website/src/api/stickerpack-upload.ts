import type { Emote } from '@/composables/use-emote-search'
import { API_URL } from './config'
import type { PackResponse } from '@/types/pack'

export type Source = '7tv'

export interface Sticker extends Emote {
  emoji_list: string[]
  source: Source
}

export interface StickerpackRequest {
  pack_name: string
  title: string
  emotes: Sticker[]
  has_watermark: boolean
  is_public: boolean
}

export interface CreatePackResponse {
  pack_url: string
  pack: PackResponse
}

interface EmoteInput {
  source: string
  id: string
  // keywords: string[]
  emoji_list: string[]
}

export async function uploadPack(request: StickerpackRequest) {
  const parsedRequest = {
    ...request,
    emotes: stickersToEmoteInput(request.emotes),
  }

  try {
    const res = await fetch(`${API_URL}/user/packs`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(parsedRequest),
    })

    if (!res.ok) {
      return await res.text()
    }

    const data: CreatePackResponse = await res.json()
    return data
  } catch (err) {
    console.error('Failed to create pack:', err)
  }
}

function stickersToEmoteInput(stickers: Sticker[]): EmoteInput[] {
  return stickers.map(sticker => ({
    source: String(sticker.source),
    id: sticker.id,
    // keywords: [sticker.name],
    emoji_list: sticker.emoji_list,
  }))
}
