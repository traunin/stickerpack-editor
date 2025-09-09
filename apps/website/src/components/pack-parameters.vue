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
      <label for="watermark">Bot name watermark</label>
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
  font-size: 1.3em;
  display: flex;
}

input {
  align-self: stretch;
  flex: 0;
  aspect-ratio: 1;
  margin: 0;
  cursor: pointer;
}

.name, .title {
  flex: 1;
}

label {
  width: 145px;
  padding-left: 10px;
  cursor: pointer;
}

.sticker-count {
  font-size: 2em;
  align-self: center;
}
</style>
