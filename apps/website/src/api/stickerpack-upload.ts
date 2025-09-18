import { API_URL } from './config'
import type { PackResponse } from '@/types/pack'
import type { Sticker } from '@/types/sticker'

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

export interface ProgressEvent {
  done: number
  total: number
}

interface EmoteInput {
  source: string
  id: string
  emoji_list: string[]
}

export async function uploadPack(
  request: StickerpackRequest,
  onProgress?: (progress: ProgressEvent) => void
): Promise<CreatePackResponse> {
  const parsedRequest = {
    ...request,
    emotes: stickersToEmoteInput(request.emotes),
  }

  const response = await fetch(`${API_URL}/user/packs`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(parsedRequest),
  })

  if (!response.ok) {
    const text = await response.text()
    let errorMessage: string
    try {
      const errJson = JSON.parse(text)
      errorMessage = errJson.message || JSON.stringify(errJson)
    } catch {
      errorMessage = text
    }
    throw new Error(`Failed to create pack: ${errorMessage}`)
  }

  if (!response.body) {
    throw new Error('Response body is null')
  }

  return new Promise<CreatePackResponse>((resolve, reject) => {
    const reader = response.body!.getReader()
    const decoder = new TextDecoder()

    function processChunk() {
      reader.read().then(({ done, value }) => {
        if (done) {
          reject(new Error('Stream ended without completion'))
          return
        }

        try {
          const chunk = decoder.decode(value, { stream: true })
          const lines = chunk.split('\n')

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              const data = line.slice(6) // remove 'data: '

              const eventMatch = chunk.match(/event:\s*(\w+)/)
              const eventType = eventMatch ? eventMatch[1] : 'data'

              if (eventType === 'error') {
                const errorMessage = JSON.parse(data)
                reject(new Error(errorMessage))
                return
              } else if (eventType === 'done') {
                const result: CreatePackResponse = JSON.parse(data)
                resolve(result)
                return
              } else {
                try {
                  const progressData: ProgressEvent = JSON.parse(data)
                  if (onProgress) {
                    onProgress(progressData)
                  }
                } catch (e) {
                  console.warn('Failed to parse progress data:', data)
                }
              }
            }
          }

          processChunk()
        } catch (error) {
          reject(new Error(`Failed to process SSE chunk: ${error}`))
        }
      }).catch(reject)
    }

    processChunk()
  })
}

function stickersToEmoteInput(stickers: Sticker[]): EmoteInput[] {
  return stickers.map(sticker => ({
    source: String(sticker.source),
    id: sticker.id,
    // keywords: [sticker.name],
    emoji_list: sticker.emoji_list,
  }))
}
