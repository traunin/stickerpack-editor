export interface Emote {
  id: string
  name: string
  preview: string
  full: string
}

export type Source = '7tv' | 'tenor' | 'telegram'

export interface Sticker extends Emote {
  uuid: string
  emoji_list: string[]
  source: Source
}

export function createSticker(emote: Emote, source: Source): Sticker {
  return {
    uuid: crypto.randomUUID(),
    source,
    emoji_list: ['😀'],
    ...emote,
  }
}

export function stickersToEmoteInput(stickers: Sticker[]): EmoteInput[] {
  return stickers.map(sticker => ({
    source: String(sticker.source),
    id: sticker.id,
    emoji_list: sticker.emoji_list,
  }))
}

interface EmoteInput {
  source: string
  id: string
  emoji_list: string[]
}
