<template>
  <div class="emote-search">
    <input
      v-model="query"
      class="search"
      placeholder="Search 7TV emotes..."
    >
    <div ref="scroll-container" class="searched-emotes">
      <LoadingAnimation v-if="isLoading" class="centered" />
      <div v-else-if="total === 0" class="centered">
        No emotes found
      </div>
      <div v-else-if="isError" class="centered">
        {{ error }}
      </div>
      <template v-else>
        <SearchResult7TV
          v-for="emote in emotes"
          :key="emote.id"
          :emote="emote"
          @click="selectEmote(emote)"
        />
        <LoadingAnimation v-if="isFetchingNextPage" class="infinite-scroll-loading" />
        <div ref="scroll-trigger" style="height: 1px; flex-shrink: 0;" />
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDebounce } from '@vueuse/core'
import { onBeforeUnmount, onMounted, ref, useTemplateRef, watch } from 'vue'
import LoadingAnimation from '@/components/loading-animation.vue'
import SearchResult7TV from '@/components/search-result-7tv.vue'
import { useScroll7TVSearch } from '@/composables/use-scroll-7tv-search'
import { createSticker } from '@/types/sticker'
import type { Emote, Sticker } from '@/types/sticker'

const emit = defineEmits<{
  (e: 'sticker-selected', sticker: Sticker): void
}>()

const query = ref('')
const debounceQuery = useDebounce(query, 300)

const {
  emotes,
  total,
  isLoading,
  isFetchingNextPage,
  hasMore,
  loadMore,
  isError,
  error,
} = useScroll7TVSearch(debounceQuery, 10)

const loadTriggerOffset = 2000
const scrollTrigger = useTemplateRef('scroll-trigger')
let observer: IntersectionObserver | null = null
const scrollContainer = useTemplateRef('scroll-container')

function setupObserver() {
  if (!scrollTrigger.value || !scrollContainer.value || observer) {
    return
  }

  observer = new IntersectionObserver(
    (entries) => {
      const target = entries[0]

      if (target.isIntersecting && hasMore.value && !isFetchingNextPage.value) {
        loadMore()
      }
    },
    {
      root: scrollContainer.value,
      rootMargin: `${loadTriggerOffset}px`,
      threshold: 0,
    },
  )

  observer.observe(scrollTrigger.value)
}

const stopWatch = watch([scrollTrigger, scrollContainer], () => {
  if (scrollTrigger.value && scrollContainer.value) {
    setupObserver()
    stopWatch()
  }
})

onMounted(() => {
  setupObserver()
})

function checkAndLoadMore() {
  if (!scrollTrigger.value || !hasMore.value || isFetchingNextPage.value) {
    return
  }
  const rect = scrollTrigger.value.getBoundingClientRect()
  const isVisible = rect.top - loadTriggerOffset < window.innerHeight
  if (isVisible) {
    loadMore()
    setTimeout(checkAndLoadMore, 300)
  }
}

watch([() => emotes.value.length, isFetchingNextPage], ([length, fetching]) => {
  if (!fetching && length > 0) {
    setTimeout(checkAndLoadMore, 100)
  }
})

onBeforeUnmount(() => {
  if (observer) {
    observer.disconnect()
  }
})

function selectEmote(emote: Emote) {
  emit('sticker-selected', createSticker(emote, '7tv'))
}
</script>

<style scoped>
.emote-search {
  padding: 10px;
  display: flex;
  flex-direction: column;
  min-height: 0;
  background: var(--panel);
  border-radius: 10px;
}

.searched-emotes {
  flex: 1;
  display: flex;
  min-height: 0;
  flex-direction: column;
  overflow-y: auto;
  scrollbar-color: var(--accent) var(--input);
  scrollbar-width: thin;
}

.search {
  background: var(--input);
  color: var(--text);
  width: 100%;
  border: 1px solid transparent;
  padding: 5px 10px;
  font-size: 1.2em;
  margin-bottom: 10px;
}

.search:focus-visible {
  border: 1px solid var(--text);
}

.infinite-scroll-loading {
  align-self: center;
  margin: 20px;
}

.centered {
  margin: auto;
}
</style>
