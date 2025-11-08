import { useInfiniteQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { fetch7TVUserEmotes } from '@/api/7tv/user-emotes'

export function useScroll7TVUserEmotes(
  userID: string,
  pageSize = 10,
  enabled = true,
) {
  const queryResult = useInfiniteQuery({
    queryKey: ['7TVuserEmotes', userID, pageSize],
    queryFn: ({ pageParam = 1 }) => fetch7TVUserEmotes(userID, pageParam, pageSize),
    getNextPageParam: (lastPage, allPages) => {
      const loadedCount = allPages.reduce((sum, page) => sum + page.items.length, 0)
      if (loadedCount >= lastPage.totalCount) {
        return undefined
      }
      return allPages.length + 1
    },
    initialPageParam: 1,
    enabled,
  })

  const emotes = computed(() => queryResult.data.value?.pages.flatMap(p => p.items) || [])
  const total = computed(() => queryResult.data.value?.pages[0]?.totalCount || 0)
  const hasMore = computed(() => queryResult.hasNextPage.value)

  const loadMore = () => {
    if (hasMore.value && !queryResult.isFetchingNextPage.value) {
      queryResult.fetchNextPage()
    }
  }

  return {
    emotes,
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
