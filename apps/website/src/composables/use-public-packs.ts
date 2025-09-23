import { type Ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { fetchPublicPacks } from '@/api/packs'

export function usePublicPacks(page: Ref<number>, pageSize: Ref<number>) {
    return useQuery({
        queryKey: ['packs', 'public', page, pageSize],
        queryFn: () => fetchPublicPacks(page.value - 1, pageSize.value),
        staleTime: 5 * 60 * 1000,
        refetchOnWindowFocus: false
    })
}
