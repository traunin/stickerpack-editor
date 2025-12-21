<template>
  <div class="pack-name-input">
    <input
      id="pack-name-input"
      v-model="name"
      type="text"
      placeholder="Name in links (a-z, 0-9, _)"
    >
    <StatusIcon :loading="loading" :error="available === false" class="icon" />
    <div v-if="available === false" class="error">
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
const { available, error, loading } = usePackNameCheck(name)

watch([available, error], () => {
  emit('error', error.value)
}, { immediate: true })
</script>

<style scoped>
.pack-name-input {
  position: relative;
  display: flex;
  flex-direction: column;
}

.icon {
  position: absolute;
  top: 4px;
  display: flex;
  align-items: center;
  right: 10px;
}

.error {
  padding: 5px 10px;
  border: 1px solid red;
  background: var(--background);
  border-radius: 5px;
  margin-top: 5px
}

input {
  flex: 1;
  background: var(--input);
  color: var(--text);
  border: 1px solid transparent;
  border-radius: 10px;
  font-size: 20px;
  padding: 7px 10px;
  padding-right: 50px;
  width: 100%;
  outline: none;
}

input:focus-visible {
  border: 1px solid var(--text);
}

@media screen and (min-width: 1000px) {
  .error {
    position: absolute;
    bottom: -30px;
    right: 0;
  }
}
</style>
