<template>
  <Teleport to="body">
    <div class="backdrop">
      <div class="modal">
        {{ message }}
        <div v-if="total" class="progress-bar">
          <div
            class="progress"
            :style="{ right: `${100 - ((progress ?? 0) / total) * 100}%` }"
          />
        </div>
        <LoadingAnimation v-else />
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import LoadingAnimation from '@/components/loading-animation.vue'

defineProps<{
  message: string
  total?: number
  progress?: number
}>()
</script>

<style scoped>
.backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 20px;
  border: 3px solid var(--primary);
  background: var(--background);
  font-size: 1.5em;
  color: var(--text);
  border-radius: 10px;
}

.progress {
  position: absolute;
  background-color: var(--primary);
  color: var(--color);
  top: 0;
  bottom: 0;
  left: 0;
  border-radius: 10px;
}

.progress-bar {
  position: relative;
  background: grey;
  width: 100%;
  height: 30px;
}
</style>
