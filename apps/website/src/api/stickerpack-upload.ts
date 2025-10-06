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
  message?: string
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

  // enqueue job
  const enqueueResp = await fetch(`${API_URL}/user/packs`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(parsedRequest),
  })

  if (!enqueueResp.ok) {
    const text = await enqueueResp.text()
    let errorMessage: string
    try {
      const errJson = JSON.parse(text)
      errorMessage = errJson.message || JSON.stringify(errJson)
    } catch {
      errorMessage = text
    }
    throw new Error(`Failed to create pack: ${errorMessage}`)
  }

  // POST returns JSON { job_id: "<uuid>", status: "queued" }
  const enqueueJson = await enqueueResp.json()
  const jobID: string | undefined =
    enqueueJson.job_id ?? enqueueJson.jobId ?? enqueueJson.id
  if (!jobID) {
    throw new Error('server did not return job_id')
  }

  // connect to job SSE endpoint (GET /job/{jobID}) and parse SSE stream
  const sseResp = await fetch(`${API_URL}/job/${jobID}`, {
    method: 'GET',
    credentials: 'include',
    headers: {
      Accept: 'text/event-stream',
    },
  })

  if (!sseResp.ok) {
    const text = await sseResp.text().catch(() => '')
    throw new Error(`Failed to open job stream: ${text || sseResp.statusText}`)
  }

  if (!sseResp.body) {
    throw new Error('job stream has no body')
  }

  return new Promise<CreatePackResponse>(async (resolve, reject) => {
    const reader = sseResp.body!.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    function findBoundary(buf: string) {
      const rn = buf.indexOf('\r\n\r\n')
      const n = buf.indexOf('\n\n')
      if (rn === -1) return n === -1 ? -1 : { pos: n, len: 2 }
      if (n === -1) return { pos: rn, len: 4 }
      return rn < n ? { pos: rn, len: 4 } : { pos: n, len: 2 }
    }

    function handleEvent(eventType: string, dataStr: string) {
      try {
        if (eventType === 'progress') {
          const p: ProgressEvent = JSON.parse(dataStr)
          if (onProgress) onProgress(p)
        } else if (eventType === 'status') {
          // const s = JSON.parse(dataStr)
        } else if (eventType === 'result') {
          const jobResult = JSON.parse(dataStr)
          if (jobResult.status === 'completed') {
            resolve(jobResult.data as CreatePackResponse)
          } else {
            const errMsg =
              jobResult.error ??
              `job failed with status ${String(jobResult.status)}`
            reject(new Error(errMsg))
          }
        } else if (eventType === 'error') {
          const errMsg = JSON.parse(dataStr)
          reject(new Error(String(errMsg)))
        } else {
          try {
            const maybe = JSON.parse(dataStr)
            if (maybe && maybe.done !== undefined && maybe.total !== undefined) {
              if (onProgress) onProgress(maybe as ProgressEvent)
            }
          } catch {
          }
        }
      } catch (err) {
        reject(new Error(`Failed to parse event ${eventType}: ${String(err)}`))
      }
    }

    function parseBuffer() {
      while (true) {
        const find = findBoundary(buffer)
        if (find === -1) break
        const { pos, len } = find as { pos: number; len: number }
        const rawEvent = buffer.slice(0, pos)
        buffer = buffer.slice(pos + len)

        const lines = rawEvent.split(/\r?\n/)
        let eventType = 'message'
        const dataLines: string[] = []
        for (const line of lines) {
          if (line.startsWith('event:')) {
            eventType = line.slice('event:'.length).trim()
          } else if (line.startsWith('data:')) {
            dataLines.push(line.slice('data:'.length).trim())
          }
        }
        const dataStr = dataLines.join('\n')
        if (dataStr.length > 0) handleEvent(eventType, dataStr)
      }
    }

    try {
      while (true) {
        const { done, value } = await reader.read()
        if (done) {
          reject(new Error('SSE stream ended without result'))
          return
        }
        buffer += decoder.decode(value, { stream: true })
        parseBuffer()
      }
    } catch (err) {
      reject(new Error(`SSE reader error: ${String(err)}`))
    } finally {
      try {
        reader.cancel().catch(() => {})
      } catch {}
    }
  })
}

function stickersToEmoteInput(stickers: Sticker[]): EmoteInput[] {
  return stickers.map(sticker => ({
    source: String(sticker.source),
    id: sticker.id,
    emoji_list: sticker.emoji_list,
  }))
}
