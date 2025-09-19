<template>
  <div class="pack-name-input">
    <input
      v-model="name"
      type="text"
      placeholder="Name in links (a-z, 0-9, _)"
    >
    <StatusIcon :loading="loading" :error="!!error" class="icon" />

    <div v-if="error" class="error">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import StatusIcon from '@/components/status-icon.vue'
import { usePackNameCheck } from '@/composables/use-pack-name-check'

const emit = defineEmits<{
  (e: 'error', message: string | null): void
}>()

const name = defineModel<string>({ default: '' })
const { error, loading } = usePackNameCheck(name)

watch(error, (val) => {
  emit('error', val)
}, { immediate: true })
</script>

<style scoped>
.pack-name-input {
  position: relative;
  display: flex;
}

.icon {
  position: absolute;
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
