<template>
  <div class="sticker">
    <button
      class="remove"
      aria-label="Remove sticker"
      @click="remove"
    >
      ✖
    </button>
    <ImageRetry :url="model.full" :alt="model.name" class="preview" />
    <input
      id="emojis"
      :value="model.emoji_list?.join('') ?? ''"
      type="text"
      placeholder="Enter emojis only..."
      aria-label="Sticker emojis"
      @input="handleInput"
    >
  </div>
</template>

<script setup lang="ts">
import { splitEmojis } from '@/api/emoji'
import ImageRetry from '@/components/image-retry.vue'
import type { Sticker } from '@/types/sticker'

const emit = defineEmits<{
  (e: 'remove'): void
}>()

const model = defineModel<Sticker>({ required: true })

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement
  model.value.emoji_list = splitEmojis(target.value)
}

function remove() {
  emit('remove')
}
</script>

<style scoped>
.sticker {
  display: flex;
  flex-direction: column;
  position: relative;
  border-radius: 5px;
  overflow: hidden;
  padding-top: 30px;
}

.preview {
  width: 144px;
  height: 144px;
}

input {
  width: 192px;
  padding: 5px;
  background: var(--input);
  border: 1px solid transparent;
  color: var(--text);
  font-size: 1.2em;
  border-radius: 5px;
  margin-top: 5px;
  outline: none;
}

input:focus-visible {
  border: 1px solid var(--text);
}

.remove {
  width: 26px;
  border-radius: 5px;
  background: #f00;
  color: var(--text);
  position: absolute;
  right: 0;
  aspect-ratio: 1;
  border: none;
  cursor: pointer;
  top: 0;
  font-size: 1em;
}

@media screen and (min-width: 1000px) {
  .preview {
    width: 192px;
    height: 192px;
  }
}
</style>
