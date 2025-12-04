import { refDebounced } from '@vueuse/core'
import { computed, ref, unref, watch } from 'vue'
import { API_URL } from '@/api/config'
import type { MaybeRef } from 'vue'

function validateName(name: string): string | null {
  if (name.length === 0) {
    return 'Name is empty'
  }
  if (!/^[a-z]/i.test(name)) {
    return 'Name must start with a letter'
  }
  if (!/^\w+$/.test(name)) {
    return 'Only letters, numbers, and underscores allowed'
  }
  if (/__/.test(name)) {
    return 'No consecutive underscores'
  }
  if (name.endsWith('_')) {
    return `Can't end with an underscore`
  }
  if (`${name}_by_${import.meta.env.VITE_BOT_NAME}`.length > 64) {
    return 'Name too long'
  }

  return null
}

export function usePackNameCheck(name: MaybeRef<string>) {
  const available = ref<boolean | null>(null)
  const error = ref<string | null>(null)
  const loading = ref(false)

  const unwrapped = computed(() => unref(name))

  const validated = computed(() => {
    available.value = null
    const validationError = validateName(unwrapped.value)
    if (validationError) {
      error.value = validationError
      available.value = false
      loading.value = false
      return null
    }
    error.value = 'Checking if name is taken...'
    loading.value = true

    return unwrapped.value
  })

  const debounced = refDebounced(validated, 300)

  watch(debounced, async (newName) => {
    if (!newName)
      return

    try {
      const res = await fetch(`${API_URL}/user/packs/${encodeURIComponent(newName)}`, {
        method: 'HEAD',
      })
      available.value = res.status !== 200
      if (!available.value) {
        error.value = 'Name taken'
      } else {
        error.value = null
      }
    } catch {
      available.value = false
      error.value = 'Network error'
    } finally {
      loading.value = false
    }
  }, { immediate: true })

  return { available, error, loading }
}
