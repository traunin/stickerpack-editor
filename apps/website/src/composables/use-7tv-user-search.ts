import { ref, watch } from 'vue'
import { fetchSearchedUsers } from '@/api/7tv/user-search'
import type { User7TV } from '@/types/user-7tv'
import type { Ref } from 'vue'

export function use7tvUserSearch(query: Ref<string>, results: number) {
  const users = ref<User7TV[] | null>()
  const error = ref<string | null>()
  const isLoading = ref(false)

  watch(query, async (q) => {
    error.value = null
    users.value = []
    if (q === '') {
      return
    }
    isLoading.value = true

    try {
      const items = await fetchSearchedUsers(q, results)
      users.value = items
    } catch (e: unknown) {
      error.value = String(e)
    }
    isLoading.value = false
  })

  return { users, error, loading: isLoading }
}
