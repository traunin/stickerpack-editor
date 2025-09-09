<template>
  <div class="pack-title-input">
    <input
      v-model="name"
      type="text"
      placeholder="Displayed title"
    >
    <div class="icon">
      <span v-if="error">❌</span>
      <span v-else>✅</span>
    </div>

    <div v-if="error" class="error">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { toRef, watch } from 'vue'
import { usePackTitleCheck } from '@/composables/use-pack-title-check'

const props = defineProps<{
  useWatermark: boolean
}>()

const emit = defineEmits<{
  (e: 'error', message: string | null): void
}>()

const name = defineModel<string>({ default: '' })
const { error } = usePackTitleCheck(name, toRef(props, 'useWatermark'))

watch(error, (val) => {
  emit('error', val)
}, { immediate: true })
</script>

<style scoped>
.pack-title-input {
  position: relative;
  display: flex;
}

.icon {
  position: absolute;
  font-size: 1.2em;
  top: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  right: 10px;
}

.error {
  position: absolute;
  bottom: -30px;
  right: 0;
  padding: 5px;
  border: 2px solid red;
  background: var(--background)
}

input {
  flex: 1;
  background: var(--input);
  color: var(--text);
  border: 3px solid var(--primary);
  font-size: 1.3em;
  padding: 5px;
  padding-right: 40px;
  width: 100%;
}
</style>
