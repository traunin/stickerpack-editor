export interface Emote {
  id: string
  name: string
  preview: string
  full: string
}

export type Source = '7tv' | 'tenor'

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
