import { API_URL } from '@/api/config'

interface ThumbnailData {
  url: string
  isVideo: boolean
}

export async function fetchThumbnail(thumbnailId: string): Promise<ThumbnailData> {
  const url = `${API_URL}/thumbnail?thumbnail_id=${thumbnailId}`
  const res = await fetch(url)

  if (!res.ok) {
    throw new Error(`Failed to fetch thumbnail: ${res.status}`)
  }

  const blob = await res.blob()
  const isVideo = blob.type === 'video/webm'
  const objectURL = URL.createObjectURL(blob)

  return { url: objectURL, isVideo }
}
