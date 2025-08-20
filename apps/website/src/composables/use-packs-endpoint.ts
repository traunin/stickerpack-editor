import { ref, watch, type Ref } from 'vue'
import { fetchPacksEndpoint, type PublicPack } from '@/api/packs'

export function usePacksEndpoint(
  endpoint: string,
  pageSize: Ref<number>,
  enabled = ref(true)
) {
  const page = ref(1)
  const maxPages = ref(1)
  const publicPacks = ref<PublicPack[] | null>()
  const error = ref<string | null>()

  const fetchData = async () => {
    if (!enabled.value) {
      publicPacks.value = null
      error.value = null
      page.value = 1
      maxPages.value = 1
      return
    }

    error.value = null
    publicPacks.value = []
    try {
      const { packs, total } = await fetchPacksEndpoint(endpoint, page.value - 1, pageSize.value)
      publicPacks.value = packs
      maxPages.value = Math.ceil(total / pageSize.value)
    } catch (e: unknown) {
      error.value = String(e)
      maxPages.value = 1
    }
  }

  watch([page, pageSize, enabled], fetchData, { immediate: true })

  function next() {
    if (page.value < maxPages.value) {
      page.value++
    }
  }

  function prev() {
    if (page.value > 1) {
      page.value--
    }
  }

  return { publicPacks, error, page, maxPages, next, prev }
}