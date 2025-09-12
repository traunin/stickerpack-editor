<template>
  <div class="emote-search-selector">
    <div class="source-switch">
      <select id="source" v-model="source">
        <option value="7tv">
          7TV
        </option>
        <option value="tenor">
          Tenor
        </option>
      </select>
    </div>

    <component
      :is="activeComponent"
      @sticker-selected="onStickerSelected"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import EmoteSearch7tv from '@/components/emote-search-7tv.vue'
import EmoteSearchTenor from '@/components/emote-search-tenor.vue'
import type { Source, Sticker } from '@/types/sticker'

const emit = defineEmits<{
  (e: 'sticker-selected', sticker: Sticker): void
}>()

const source = ref<Source>('7tv')

const activeComponent = computed(() =>
  source.value === '7tv' ? EmoteSearch7tv : EmoteSearchTenor,
)

function onStickerSelected(sticker: Sticker) {
  emit('sticker-selected', sticker)
}
</script>

<style scoped>
.emote-search-selector {
  display: flex;
  flex-direction: column;
  height: 100%;
  flex: 1;
}

.source-switch {
  display: flex;
  align-items: center;
  cursor: pointer
}

select {
  background: var(--input);
  color: var(--text);
  border: 2px solid var(--primary);
  padding: 4px 8px;
  font-size: 1.2em;
  align-self: stretch;
  width: 100%;
  text-align: center;
  cursor: pointer;
}
</style>
