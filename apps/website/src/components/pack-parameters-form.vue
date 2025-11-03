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
      :use-watermark="hasWatermark"
      @error="forwardTitleError"
    />
    <div class="watermark">
      <input id="watermark" v-model="hasWatermark" type="checkbox" checked>
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
import { toRefs } from 'vue'
import PackNameInput from '@/components/pack-name-input.vue'
import type { PackParameters } from '@/types/pack'
import PackTitleInput from './pack-title-input.vue'

defineProps<{
  stickerCount: number
  maxStickers: number
}>()

const emit = defineEmits<{
  (e: 'name-error', value: string | null): void
  (e: 'title-error', value: string | null): void
}>()

const params = defineModel<PackParameters>({ required: true })
const { name, title, hasWatermark, isPublic } = toRefs(params.value)

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
  gap: 15px;
  margin-bottom: 0;
  background: var(--panel);
  padding: 10px;
  border-radius: 10px;
}

.watermark, .public {
  font-size: 1.2em;
  display: flex;
  align-items: center;
}

input[type="checkbox"] {
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background-color: var(--input);
  border: none;
  cursor: pointer;
  display: inline-block;
  position: relative;
  transition: background-color 0.1s;
  margin: 0;
}

input[type="checkbox"]:checked {
  background-color: #fff;
}

input[type="checkbox"]:checked::after {
  content: "";
  position: absolute;
  top: 8px;
  left: 14px;
  width: 10px;
  height: 18px;
  border: solid #000;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

input[type="checkbox"]:active {
  transform: scale(0.95);
}

.name, .title {
  flex: 1;
}

label {
  width: 145px;
  padding-left: 10px;
  cursor: pointer;
  user-select: none;
}

.sticker-count {
  font-size: 2em;
  align-self: center;
}
</style>
