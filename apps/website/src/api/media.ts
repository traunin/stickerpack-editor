interface MediaData {
  url: string
  isVideo: boolean
}

export async function fetchMedia(url: string): Promise<MediaData> {
  const res = await fetch(url)

  if (!res.ok) {
    throw new Error(`Failed to fetch media: ${res.status}`)
  }

  const blob = await res.blob()
  const isVideo = blob.type === 'video/webm'
  const objectUrl = URL.createObjectURL(blob)

  return { url: objectUrl, isVideo }
}
