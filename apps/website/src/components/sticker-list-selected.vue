<template>
  <div class="selected-stickers">
    <draggable
      v-model="stickers"
      item-key="uuid"
      group="stickers"
      ghost-class="ghost-item"
      chosen-class="chosen-item"
      drag-class="drag-item"
      class="drag-area"
    >
      <template #item="{ index }">
        <div class="sticker-wrapper">
          <div class="drag-handle">
            ⋮⋮
          </div>
          <StickerCreate
            :model-value="stickers[index]"
            @update:model-value="updateSticker(index, $event)"
            @remove="removeSticker(index)"
          />
        </div>
      </template>
    </draggable>
  </div>
</template>

<script setup lang="ts">
import draggable from 'vuedraggable'
import StickerCreate from '@/components/sticker-create.vue'
import type { Sticker } from '@/types/sticker'

const stickers = defineModel<Sticker[]>({ required: true })

function updateSticker(index: number, value: Sticker) {
  stickers.value[index] = value
}

function removeSticker(index: number) {
  stickers.value.splice(index, 1)
}
</script>

<style scoped>
.selected-stickers {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: flex-start;
  align-self: flex-start;
  min-height: 0;
  overflow-y: auto;
  align-self: stretch;
  background: var(--panel);
  scrollbar-color: var(--accent) var(--input);
  scrollbar-width: thin;
  border-radius: 10px;
}

.drag-area > :first-child::before {
  content: "Preview";
  position: absolute;
  top: 5px;
  height: 26px;
  line-height: 26px;
  padding: 0 6px;
  background: var(--primary);
  color: var(--text);
  z-index: 5;
  left: 50%;
  transform: translate(-50%);
  border-radius: 5px;
}

.drag-area {
  display: grid;
  grid-template-columns: repeat(auto-fill, 202px);
  flex: 1;
  gap: 10px;
  margin: 15px;
  justify-content: center;
}

.sticker-wrapper {
  position: relative;
  padding: 5px;
  background: var(--panel);
  border-radius: 5px;
  cursor: grab;
}

.drag-handle {
  position: absolute;
  width: 26px;
  height: 26px;
  border-radius: 5px;
  text-align: center;
  color: var(--text);
  background: var(--primary);
  z-index: 10;
  font-size: 1.2em;
  line-height: 28px;
}
</style>
