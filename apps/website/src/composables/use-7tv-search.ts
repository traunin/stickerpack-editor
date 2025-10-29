import { ref, watch } from 'vue'
import { fetchSearchedEmotes } from '@/api/7tv/emote-search'
import type { Emote } from '@/types/sticker'
import type { Ref } from 'vue'

export function use7tvSearch(query: Ref<string>, pageSize: number) {
  const page = ref(1)
  const maxPages = ref(1)
  const emotes = ref<Emote[] | null>()
  const error = ref<string | null>()
  let lastQuery = ''

  watch([query, page], async ([q, p]) => {
    if (q !== lastQuery) {
      lastQuery = q
      if (page.value !== 1) {
        page.value = 1
        return // watch triggered by setting page to 1
      }
    }

    error.value = null
    emotes.value = []

    try {
      const { items, pageCount } = await fetchSearchedEmotes(q, p, pageSize)
      emotes.value = items
      maxPages.value = pageCount
    } catch (e: unknown) {
      error.value = String(e)
      maxPages.value = 1
    }
  }, { immediate: true })

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

  return { emotes, error, page, maxPages, next, prev }
}
