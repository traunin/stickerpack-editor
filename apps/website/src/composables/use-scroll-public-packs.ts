import { useInfiniteQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import type { PacksResponse } from '@/api/packs'
import { fetchPublicPacks } from '@/api/packs'

export function useScrollPublicPacks(pageSize = 10, enabled = true) {
  const queryResult = useInfiniteQuery({
    queryKey: ['packs', 'public', pageSize],
    queryFn: ({ pageParam = 0 }) => fetchPublicPacks(pageParam, pageSize),
    getNextPageParam: (lastPage: PacksResponse, allPages) => {
      const loadedCount = allPages.reduce((sum, page) => sum + page.packs.length, 0)
      if (loadedCount >= lastPage.total) {
        return undefined
      }
      return allPages.length
    },
    initialPageParam: 0,
    enabled,
  })

  const packs = computed(() => {
    return queryResult.data.value?.pages.flatMap(page => page.packs) || []
  })

  const total = computed(() => {
    return queryResult.data.value?.pages[0]?.total || 0
  })

  const hasMore = computed(() => queryResult.hasNextPage.value)

  const loadMore = () => {
    if (hasMore.value && !queryResult.isFetchingNextPage.value) {
      queryResult.fetchNextPage()
    }
  }

  return {
    packs,
    total,
    allPages: queryResult.data,
    isLoading: queryResult.isLoading,
    isFetching: queryResult.isFetching,
    isFetchingNextPage: queryResult.isFetchingNextPage,
    hasMore,
    loadMore,
    fetchNextPage: queryResult.fetchNextPage,
    error: queryResult.error,
    isError: queryResult.isError,
    refetch: queryResult.refetch,
  }
}
