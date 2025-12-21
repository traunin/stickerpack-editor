<template>
  <div class="pack-title-input">
    <input
      id="pack-title-input"
      v-model="name"
      type="text"
      placeholder="Displayed title"
    >
    <StatusIcon :loading="false" :error="!!error" class="icon" />

    <div v-if="error" class="error">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { toRef, watch } from 'vue'
import StatusIcon from '@/components/status-icon.vue'
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
  flex-direction: column;
}

.icon {
  position: absolute;
  font-size: 1.2em;
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
  font-size: 1.2em;
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
