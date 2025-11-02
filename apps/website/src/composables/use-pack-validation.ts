import { computed } from 'vue'
import { useTgAuthStore } from '@/stores/use-tg-auth'
import type { Ref } from 'vue'

export function usePackValidation(
  nameError: Ref<string | null>,
  titleError: Ref<string | null>,
  stickerCount: Ref<number>,
  maxStickers: number,
) {
  const authStore = useTgAuthStore()

  return computed(() => {
    if (!authStore.isLoggedIn)
      return 'You are not logged in'
    if (nameError.value)
      return nameError.value
    if (titleError.value)
      return titleError.value
    if (stickerCount.value === 0)
      return 'No emotes selected'
    if (stickerCount.value > maxStickers)
      return `Too many emotes (max ${maxStickers})`
    return null
  })
}
