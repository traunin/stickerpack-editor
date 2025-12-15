import type { PackPreview } from '@/types/pack'
import { stickersToEmoteInput } from '@/types/sticker'
import type { Sticker } from '@/types/sticker'
import { enqueueJob, handleJobSSE } from './job'
import type { ProgressEvent } from './job'

export interface CreatePackRequest {
  pack_name: string
  title: string
  emotes: Sticker[]
  has_watermark: boolean
  is_public: boolean
}

export interface CreatePackResponse {
  pack_url: string
  pack: PackPreview
}

export async function createPack(
  request: CreatePackRequest,
  onProgress?: (progress: ProgressEvent) => void,
): Promise<CreatePackResponse> {
  const parsedRequest = {
    ...request,
    emotes: stickersToEmoteInput(request.emotes),
  }

  const jobID = await enqueueJob('/user/packs', 'POST', parsedRequest)
  return handleJobSSE<CreatePackResponse>(jobID, onProgress)
}
