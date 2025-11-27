import { useQuery } from '@tanstack/vue-query'
import { fetchMedia } from '@/api/media'

export function useMedia(fileId: string, retries = 3) {
  return useQuery({
    queryKey: ['media', fileId],
    queryFn: () => fetchMedia(fileId),
    retry: retries,
    staleTime: 60 * 60 * 1000,
    gcTime: 24 * 60 * 60 * 1000,
    refetchOnWindowFocus: false,
    refetchOnMount: false,
  })
}
