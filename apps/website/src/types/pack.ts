import type { Sticker } from './sticker'

export interface PackPreview {
  id: number
  title: string
  name: string
  thumbnail_id: string
}

export interface PackParameters {
  name?: string
  title: string
  hasWatermark?: boolean
  isPublic: boolean
}

export interface Pack {
  title: string
  isPublic: boolean
  stickers: Sticker[]
}

export interface TelegramStickerSet {
  name: string
  title: string
  is_public: boolean
  stickers: TelegramSticker[]
}

export interface TelegramSticker {
  file_id: string
  fileUniqueID: string
  emoji: string
}
