import { computed, unref } from 'vue'
import type { MaybeRef } from 'vue'

export function usePackTitleCheck(
  titleRef: MaybeRef<string>,
  hasWatermarkRef: MaybeRef<boolean>,
) {
  const error = computed(() => {
    const title = unref(titleRef)
    const hasWatermark = unref(hasWatermarkRef)

    if (title.length < 1) {
      return 'Title is empty'
    }

    if (hasWatermark) {
      const watermark = ` by @${import.meta.env.VITE_BOT_NAME}`
      const watermarkedName = `${title}${watermark}`
      if (watermarkedName.length > 64) {
        return 'Title with watermark too long'
      }
    } else if (title.length > 64) {
      return 'Title too long'
    }

    return null
  })

  return { error }
}
