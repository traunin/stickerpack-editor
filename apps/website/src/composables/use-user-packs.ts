import { type Ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { fetchUserPacks } from '@/api/packs'

export function useUserPacks(page: Ref<number>, pageSize: Ref<number>, enabled: Ref<boolean>) {
    return useQuery({
        queryKey: ['packs', 'user', page, pageSize],
        queryFn: () => fetchUserPacks(page.value - 1, pageSize.value),
        enabled,
        staleTime: 5 * 60 * 1000,
        refetchOnWindowFocus: false
    })
}
