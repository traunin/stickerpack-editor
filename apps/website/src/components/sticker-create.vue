<template>
  <div class="sticker">
    <button class="remove" @click="remove">
      ✖
    </button>
    <div class="preview">
      <img :src="model.full" :alt="model.name">
    </div>
    <input
      id="emojis"
      v-model="emojisInput"
      type="text"
      placeholder="Enter emojis only..."
      @input="handleInput"
      @paste="handlePaste"
    >
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick } from 'vue'
import type { Sticker } from '@/types/sticker'

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

function extractEmojis(text: string): string {
  const emojiArray = [...text.matchAll(/([\p{Extended_Pictographic}\p{Emoji_Component}])+/gu)]
  return emojiArray.map(match => match[0]).join('')
}

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement
  const currentValue = target.value
  const filteredValue = extractEmojis(currentValue)

  if (currentValue !== filteredValue) {
    emojisInput.value = filteredValue

    nextTick(() => {
      target.value = filteredValue
      target.setSelectionRange(filteredValue.length, filteredValue.length)
    })
  }
}

function handlePaste(event: ClipboardEvent) {
  event.preventDefault()
  const pastedText = event.clipboardData?.getData('text') || ''
  const filteredEmojis = extractEmojis(pastedText)

  if (filteredEmojis) {
    const target = event.target as HTMLInputElement
    const start = target.selectionStart || 0
    const end = target.selectionEnd || 0
    const currentValue = emojisInput.value

    const newValue = currentValue.substring(0, start) + filteredEmojis + currentValue.substring(end)
    emojisInput.value = newValue

    nextTick(() => {
      target.setSelectionRange(start + filteredEmojis.length, start + filteredEmojis.length)
    })
  }
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
  color: var(--text);
  font-size: 1.2em;
}

.remove {
  width: 20px;
  line-height: 20px;
  background: #f00;
  color: var(--text);
  position: absolute;
  right: 0;
  aspect-ratio: 1;
  border: none;
  cursor: pointer;
}
</style>
