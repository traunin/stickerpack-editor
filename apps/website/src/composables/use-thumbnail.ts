import { useQuery } from '@tanstack/vue-query'
import { fetchThumbnail } from '@/api/thumbnail'

export function useThumbnail(thumbnailId: string) {
  return useQuery({
    queryKey: ['thumbnail', thumbnailId],
    queryFn: () => fetchThumbnail(thumbnailId),
    staleTime: 60 * 60 * 1000,
    gcTime: 24 * 60 * 60 * 1000,
    refetchOnWindowFocus: false,
    refetchOnMount: false,
  })
}
