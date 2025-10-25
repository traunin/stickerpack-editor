import { useWindowSize } from '@vueuse/core'
import { computed } from 'vue'

const { width } = useWindowSize()

export function useIsDesktop(desktopSize = 1000) {
  return computed(() => width.value > desktopSize)
}
