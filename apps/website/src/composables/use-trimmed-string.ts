// useTrimmedString.ts
import { computed, Ref, unref } from "vue"

export function useTrimmedString(
  source: Ref<string> | string,
  trimmedLength: number
) {
  const trimmed = computed(() => {
    const value = unref(source)
    return value.length <= (trimmedLength - 3)
      ? value
      : `${value.substring(0, trimmedLength - 3)}...`
  })

  return { trimmed }
}
