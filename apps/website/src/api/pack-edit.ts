import type { PackPreview } from '@/types/pack'
import { stickersToEmoteInput } from '@/types/sticker'
import type { Sticker } from '@/types/sticker'
import { enqueueJob, handleJobSSE } from './job'
import type { ProgressEvent } from './job'

export interface EditPackRequest {
  updated_title?: string
  updated_is_public?: boolean
  deleted_stickers: string[]
  added_stickers: Sticker[]
  emoji_updates: StickerEmojiUpdate[]
  position_updates: StickerPositionUpdate[]
}

export interface EditPackResponse {
  pack: PackPreview
}

export interface StickerEmojiUpdate {
  id: string
  emojis: string[]
}

export interface StickerPositionUpdate {
  id: string
  position: number
}

export async function editPack(
  name: string,
  request: EditPackRequest,
  onProgress?: (progress: ProgressEvent) => void,
): Promise<EditPackResponse> {
  const parsedRequest = {
    ...request,
    added_stickers: stickersToEmoteInput(request.added_stickers),
  }
  console.log(parsedRequest)
  const jobID = await enqueueJob(`/user/packs/${name}`, 'PATCH', parsedRequest)
  return handleJobSSE<EditPackResponse>(jobID, onProgress)
}
