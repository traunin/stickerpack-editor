<template>
  <div class="sticker">
    <button class="close" @click="remove">
      ✖
    </button>
    <div class="preview">
      <img :src="model.full" :alt="model.name">
    </div>
    <input
      id="emojis"
      v-model="emojisInput"
      type="text"
    >
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Sticker } from '@/api/stickerpack-upload'

const emit = defineEmits<{
  (e: 'remove'): void
}>()

const model = defineModel<Sticker>({ required: true })

const emojisInput = computed({
  get: () => model.value.emoji_list?.join('') ?? '',
  set: (input: string) => {
    model.value.emoji_list = Array.from(input)
  },
})

function remove() {
  emit('remove')
}
</script>

<style scoped>
.sticker {
  display: flex;
  flex-direction: column;
  position: relative;
}

.preview {
  width: 192px;
  height: 192px;
  display: flex;
  align-items: center;
}

img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

input {
  width: 192px;
  padding: 5px;
  background: var(--input);
  border: 2px solid var(--primary);
  color: var(--text)
}

.close {
  background: #f00;
  color: var(--text);
  position: absolute;
  right: 0;
  aspect-ratio: 1;
  border: none;
  cursor: pointer;
}
</style>
