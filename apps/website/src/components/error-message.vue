<template>
  <Transition>
    <div v-if="error.message.value" class="error">
      <div class="icon">
        !
      </div>
      <div class="message">
        {{ error.message }}
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import type { useErrorPopup } from '@/composables/use-error-popup'

const props = defineProps<{
  error: ReturnType<typeof useErrorPopup>
  cleanupTimeout?: number
}>()

const timeoutId = ref<number | null>(null)

function startCleanupTimer() {
  if (!props.cleanupTimeout)
    return
  if (timeoutId.value)
    clearTimeout(timeoutId.value)
  timeoutId.value = window.setTimeout(() => {
    props.error.clear()
  }, props.cleanupTimeout)
}

watch(
  () => props.error.message.value,
  (msg) => {
    if (msg)
      startCleanupTimer()
  },
)

onBeforeUnmount(() => {
  if (timeoutId.value)
    clearTimeout(timeoutId.value)
})
</script>

<style scoped>
.error {
  font-size: 1.4em;
  background: var(--background);
  border: 2px solid var(--primary);
  color: var(--text);
  display: flex;
  border-radius: 10px;
  overflow: hidden;
  position: fixed;
  top: 20px;
  left: 20px;
  right: 20px;
  z-index: 20;
}

.icon {
  background: red;
  aspect-ratio: 1;
  padding: 10px 20px;
  font-weight: 900;
  display: flex;
  align-items: center;
  justify-content: center;
}

.message {
  padding: 10px;
  display: flex;
  align-items: center;
}

.v-enter-active,
.v-leave-active {
  transition: top 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  top: -15%;
}
</style>
