<template>
  <div class="sticker">
    <button class="remove" @click="remove">
      ✖
    </button>
    <ImageRetry :url="model.full" :alt="model.name" class="preview" />
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
import ImageRetry from '@/components/image-retry.vue'
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
  border-radius: 5px;
  overflow: hidden;
  padding-top: 30px;
}

.preview {
  width: 192px;
  height: 192px;
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
</style>
