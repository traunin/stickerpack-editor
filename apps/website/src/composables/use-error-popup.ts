import { ref } from 'vue'

export function useErrorPopup() {
  const message = ref<string | null>(null)

  function show(error: string | Error) {
    message.value = typeof error === 'string' ? error : error.message
  }

  function clear() {
    message.value = null
  }

  return {
    message,
    show,
    clear,
  }
}
