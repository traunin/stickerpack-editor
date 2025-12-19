import { useQuery } from '@tanstack/vue-query'
import { fetchPack } from '@/api/packs'

export function useLoadPack(name: string) {
  return useQuery({
    queryKey: ['pack', name],
    queryFn: () => fetchPack(name),
    enabled: name != null,
    staleTime: 1000 * 60 * 60,
  })
}
