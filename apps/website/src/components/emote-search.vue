<template>
  <div class="emote-search">
    <div class="searchbar">
      <input
        v-model="query"
        class="search"
        placeholder="Search 7tv emotes..."
      >
    </div>
    <div v-if="loading" class="results">
      Searching...
    </div>

    <div v-else-if="error" class="results">
      {{ error }}
    </div>

    <div v-else class="results stickers">
      <SearchResult
        v-for="emote in emotes"
        :key="emote.id"
        :emote="emote"
        @click="selectEmote(emote)"
      />
    </div>

    <div class="page-controls">
      <button :disabled="page === 0" @click="prev">
        &lt;
      </button>
      <span class="text">
        <span>1..</span>
        <span>{{ page }}</span>
        <span>..{{ maxPages }}</span>
      </span>
      <button :disabled="emotes != null && (emotes.length < pageSize)" @click="next">
        &gt;
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDebounce } from '@vueuse/core'
import { computed, ref } from 'vue'
import SearchResult from '@/components/search-result.vue'
import type { Emote } from '@/composables/use-emote-search.ts'
import { useEmoteSearch } from '@/composables/use-emote-search.ts'

const emit = defineEmits<{
  (e: 'emote-selected', emote: Emote): void
}>()

const query = ref('')
const debounceQuery = useDebounce(query, 300)

const { emotes, error, page, next, prev, maxPages } = useEmoteSearch(debounceQuery, 10)
const pageSize = 10
const foundStickers = computed(() => emotes.value?.length)
const loading = computed(() => !foundStickers.value && error == null)

function selectEmote(emote: Emote) {
  emit('emote-selected', emote)
}
</script>

<style scoped>
.emote-search {
  display: flex;
  flex-direction: column;
  margin: 20px;
  flex: 1;
  min-height: 0;
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
}

.text {
  text-align: center;
  align-self: center;
}
</style>
