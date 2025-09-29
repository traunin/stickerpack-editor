import { computed, MaybeRef, ref, unref, watch } from 'vue'
import { useDebounce } from '@vueuse/core'
import { API_URL } from '@/api/config'

function validateName(name: string): string | null {
    if (name.length == 0) {
      return 'Name is empty'
    }
    if (!/^[a-zA-Z0-9_]+$/.test(name)) {
      return 'Only letters, numbers, and underscores allowed'
    }
    if (/__/.test(name)) {
      return 'No consecutive underscores'
    }
    if (/_$/.test(name)) {
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
    const validationError = validateName(unwrapped.value)
    if (validationError) {
      error.value = validationError
      return null
    }
    error.value = null
    loading.value = true

    return unwrapped.value
  })
  
  const debounced = useDebounce(validated, 300)

  watch(debounced, async (newName) => {
    if (!newName) return

    available.value = null
    try {
      const res = await fetch(`${API_URL}/user/packs/${encodeURIComponent(newName)}`, {
        method: 'HEAD',
      })
      if (res.status == 200) {
        error.value = "Name taken"
      } 
    } catch {
      error.value = 'Network error'
    } finally {
      loading.value = false
    }
  }, { immediate: true })

  return { error, loading }
}
