import { ref, watch } from 'vue'

type PackEventType = 'deleted'

interface PackEvent {
  type: PackEventType
  packName: string
}

const packEventBus = ref<PackEvent | null>(null)

export function usePackEvents() {
  function emitPackEvent(type: PackEventType, packName: string) {
    packEventBus.value = {
      type,
      packName,
    }
  }

  function onPackEvent(callback: (event: PackEvent) => void) {
    return watch(
      packEventBus,
      (event) => {
        if (event) {
          callback(event)
        }
      },
    )
  }

  return { emitPackEvent, onPackEvent }
}
