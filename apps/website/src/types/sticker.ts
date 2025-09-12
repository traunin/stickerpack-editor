export interface Emote {
  id: string
  name: string
  preview: string
  full: string
}

export type Source = '7tv' | 'tenor'

export interface Sticker extends Emote {
  emoji_list: string[]
  source: Source
}
