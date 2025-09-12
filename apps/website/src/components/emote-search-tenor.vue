<template>
  <div class="emote-search">
    <div class="searchbar">
      <input
        v-model="query"
        class="search"
        placeholder="Search tenor gifs..."
      >
    </div>
    <div v-if="!debounceQuery" class="results">
      Type in your query
    </div>

    <div v-else-if="loading" class="results loading">
      <LoadingAnimation />
    </div>

    <div v-else-if="error" class="results">
      {{ error }}
    </div>

    <div v-else class="results stickers">
      <SearchResultTenor
        v-for="emote in emotes"
        :key="emote.id"
        :emote="emote"
        @click="selectSticker(emote)"
      />
    </div>

    <div class="page-controls">
      <button :disabled="!hasPrev" @click="prev">
        &lt;
      </button>
      <button :disabled="!hasNext" @click="next">
        &gt;
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDebounce } from '@vueuse/core'
import { computed, ref } from 'vue'
import LoadingAnimation from '@/components/loading-animation.vue'
import { useTenorSearch } from '@/composables/use-tenor-emote-search'
import type { Emote, Sticker } from '@/types/sticker'
import SearchResultTenor from './search-result-tenor.vue'

const emit = defineEmits<{
  (e: 'sticker-selected', sticker: Sticker): void
}>()

const query = ref('')
const debounceQuery = useDebounce(query, 300)

const pageSize = 5
const { emotes, error, hasNext, hasPrev, next, prev } = useTenorSearch(debounceQuery, pageSize)

const foundStickers = computed(() => emotes.value?.length !== 0)
const loading = computed(() => !foundStickers.value && error.value == null)

function selectSticker(emote: Emote) {
  emit('sticker-selected', { ...emote, source: 'tenor', emoji_list: ['😀'] })
}
</script>

<style scoped>
.emote-search {
  padding: 10px;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  background: var(--panel);
}

.results {
  flex: 1;
  display: flex;
  min-height: 0;
  flex-direction: column;
  justify-content: space-around;
  overflow-y: auto;
  scrollbar-color: var(--accent) var(--input);
  scrollbar-width: thin;
}

.search {
  background: var(--input);
  color: var(--text);
  width: 100%;
  border: 2px solid var(--primary);
  padding: 5px;
  font-size: 1.3em;
}

.page-controls {
  display: flex;
  font-size: 1.3em;
  gap: 10px;
  justify-content: center;
}

.page-controls button {
  font-size: 1.3em;
  background: var(--primary);
  cursor: pointer;
  border: none;
  color: var(--text);
  height: 100%;
  line-height: 0.9em;
  aspect-ratio: 1/1;
  display: flex;
  align-items: center;
}

.text {
  text-align: center;
  align-self: center;
}

.loading {
  align-items: center
}

button:disabled {
  color: red;
}
</style>
