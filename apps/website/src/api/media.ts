interface MediaData {
  blob: Blob
  isVideo: boolean
}

export async function fetchMedia(url: string): Promise<MediaData> {
  const res = await fetch(url)

  if (!res.ok) {
    throw new Error(`Failed to fetch media: ${res.status}`)
  }

  const blob = await res.blob()
  const isVideo = blob.type.startsWith('video/')

  return { blob, isVideo }
}
