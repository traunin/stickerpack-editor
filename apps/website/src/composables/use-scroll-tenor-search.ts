import { useInfiniteQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { searchTenor } from '@/api/tenor/search'
import type { Ref } from 'vue'

export function useScrollTenorSearch(query: Ref<string>, pageSize = 10, enabled = true) {
  const queryResult = useInfiniteQuery({
    queryKey: ['tenorSearch', query, pageSize],
    queryFn: ({ pageParam = '' }) => searchTenor(query.value, pageParam, pageSize),
    getNextPageParam: (lastPage) => {
      return lastPage.next && lastPage.items.length > 0 ? lastPage.next : undefined
    },
    initialPageParam: '',
    enabled: enabled && !!query,
  })

  const emotes = computed(() => queryResult.data.value?.pages.flatMap(p => p.items) || [])
  const loaded = computed(() => emotes.value.length)
  const hasMore = computed(() => queryResult.hasNextPage.value)

  const loadMore = () => {
    if (hasMore.value && !queryResult.isFetchingNextPage.value) {
      queryResult.fetchNextPage()
    }
  }

  return {
    emotes,
    loaded,
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
