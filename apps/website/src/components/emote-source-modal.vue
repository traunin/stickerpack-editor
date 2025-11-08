<template>
  <Teleport to="body">
    <div class="source-modal" v-bind="$attrs">
      <button class="modal-control" @click="isOpen = !isOpen">
        {{ buttonMessage }} ({{ stickerCount }}/{{ maxStickers }})
      </button>
      <div v-if="isOpen" class="backdrop">
        <div class="modal">
          <EmoteSource @sticker-selected="onStickerSelected" />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Sticker } from '@/types/sticker'
import EmoteSource from './emote-source.vue'

defineProps<{
  stickerCount: number
  maxStickers: number
}>()
const emit = defineEmits<{
  (e: 'sticker-selected', sticker: Sticker): void
}>()

const isOpen = ref(false)
const buttonMessage = computed(() => isOpen.value ? 'Close' : 'Add stickers')

function onStickerSelected(sticker: Sticker) {
  emit('sticker-selected', sticker)
}
</script>

<style scoped>
.backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 10px;
  border: 2px solid var(--primary);
  background: var(--background);
  color: var(--text);
  border-radius: 10px;
  z-index: 15;
  inset: 10px;
}

.modal-control {
  border: 2px solid var(--primary);
  position: absolute;
  right: 5px;
  bottom: 5px;
  background-color: var(--input);
  color: var(--text);
  font-size: 1em;
  padding: 10px;
  cursor: pointer;
  z-index: 20;
  border-radius: 10px;
}

.source-modal {
  display: block;
}

@media screen and (min-width: 1000px) {
  .source-modal {
    display: none;
  }
}
</style>
