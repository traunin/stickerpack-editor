<template>
  <div class="parameters">
    <PackNameInput
      v-model="name"
      class="name"
      @error="forwardNameError"
    />
    <PackTitleInput
      v-model="title"
      class="title"
      :use-watermark="watermark"
      @error="forwardTitleError"
    />
    <div class="watermark">
      <input id="watermark" v-model="watermark" type="checkbox" checked>
      <label for="watermark">Use a bot name watermark</label>
    </div>
    <div class="public">
      <input id="public" v-model="isPublic" type="checkbox" checked>
      <label for="public">Show to other users</label>
    </div>
    <div class="sticker-count">
      {{ stickerCount }} / {{ maxStickers }}
    </div>
  </div>
</template>

<script setup lang="ts">
import PackNameInput from '@/components/pack-name-input.vue'
import PackTitleInput from './pack-title-input.vue'

defineProps<{
  stickerCount: number
  maxStickers: number
}>()

const emit = defineEmits<{
  (e: 'name-error', value: string | null): void
  (e: 'title-error', value: string | null): void
}>()

const name = defineModel<string>('name', { default: '' })
const title = defineModel<string>('title', { default: '' })
const watermark = defineModel<boolean>('watermark', { default: true })
const isPublic = defineModel<boolean>('isPublic', { default: true })

function forwardNameError(e: string | null) {
  emit('name-error', e)
}

function forwardTitleError(e: string | null) {
  emit('title-error', e)
}
</script>

<style scoped>
.parameters {
  display: flex;
  gap: 20px;
  margin-bottom: 0;
  background: var(--panel);
  padding: 10px;
}

.watermark, .public {
  flex: 1;
  font-size: 1.5em;
  display: flex;
  gap: 20px;
}

input {
  align-self: stretch;
  flex: 0;
  aspect-ratio: 1;
}

.name, .title {
  flex: 1;
}

label {
  flex: 1;
  align-self: center;
}

.sticker-count {
  font-size: 2em;
  align-self: center;
}
</style>
