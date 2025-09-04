import { ref, watch, type Ref } from 'vue'
import { fetchPacksEndpoint } from '@/api/packs'
import type { PackResponse } from '@/types/pack'
import { usePackEvents } from './use-packs-events'

export function usePacksEndpoint(
  endpoint: string,
  pageSize: Ref<number>,
  enabled = ref(true)
) {
  const page = ref(1)
  const maxPages = ref(1)
  const packs = ref<PackResponse[] | null>()
  const error = ref<string | null>()

  const { onPackEvent } = usePackEvents()

  const fetchData = async () => {
    if (!enabled.value) {
      packs.value = null
      error.value = null
      page.value = 1
      maxPages.value = 1
      return
    }

    error.value = null
    packs.value = null
    try {
      const { packs: packsResponse, total } = await fetchPacksEndpoint(
        endpoint,
        page.value - 1,
        pageSize.value
      )
      packs.value = packsResponse
      maxPages.value = Math.ceil(total / pageSize.value)
    } catch (e: unknown) {
      error.value = String(e)
      maxPages.value = 1
    }
  }

  watch([page, pageSize, enabled], fetchData, { immediate: true })

  onPackEvent((event) => {
    if (event.type === 'deleted' && packs.value) {
      packs.value = packs.value.filter(pack => pack.name !== event.packName)
    }
  })

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

function refetch() {
    return fetchData()
  }

  return { packs, error, page, maxPages, next, prev, refetch }
}
