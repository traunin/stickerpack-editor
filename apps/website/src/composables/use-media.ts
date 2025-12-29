import { useQuery } from '@tanstack/vue-query'
import { computed, ref, watch } from 'vue'
import { fetchMedia } from '@/api/media'

export function useMedia(fileId: string, retries = 3) {
  const query = useQuery({
    queryKey: ['media', fileId],
    queryFn: () => fetchMedia(fileId),
    retry: retries,
    staleTime: 60 * 60 * 1000,
    gcTime: 24 * 60 * 60 * 1000,
    refetchOnWindowFocus: false,
    refetchOnMount: false,
  })

  const objectUrl = ref<string | undefined>()

  watch(
    () => query.data.value,
    (newData, _, onCleanup) => {
      if (!newData) {
        objectUrl.value = undefined
        return
      }

      const newUrl = URL.createObjectURL(newData.blob)
      objectUrl.value = newUrl

      onCleanup(() => {
        URL.revokeObjectURL(newUrl)
      })
    },
    { immediate: true },
  )

  const displayData = computed(() => {
    if (!query.data.value || !objectUrl.value)
      return undefined
    return {
      url: objectUrl.value,
      isVideo: query.data.value.isVideo,
    }
  })

  return {
    ...query,
    data: displayData,
  }
}
