<template>
  <div class="user-emotes">
    <LoadingAnimation v-if="isLoading" />
    <div v-else-if="total === 0" class="results">
      No emotes found
    </div>
    <div v-else-if="isError" class="results">
      {{ error }}
    </div>
    <div v-else ref="scrollContainer" class="results packs">
      <SearchResult7TV
        v-for="emote in emotes"
        :key="emote.id"
        :emote="emote"
        @click="selectEmote(emote)"
      />
      <div v-if="isFetchingNextPage" class="results loading">
        <LoadingAnimation />
      </div>
    </div>
    <div ref="scrollTrigger" style="height: 1px; flex-shrink: 0;" />
  </div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import LoadingAnimation from '@/components/loading-animation.vue'
import SearchResult7TV from '@/components/search-result-7tv.vue'
import { useScroll7TVUserEmotes } from '@/composables/use-scroll-7tv-user-emotes'
import { createSticker } from '@/types/sticker'
import type { Emote, Sticker } from '@/types/sticker'
import type { User7TV } from '@/types/user-7tv'

const props = defineProps<{
  streamer: User7TV
}>()

const emit = defineEmits<{
  (e: 'sticker-selected', sticker: Sticker): void
}>()

const {
  emotes,
  total,
  isLoading,
  isFetchingNextPage,
  hasMore,
  loadMore,
  isError,
  error,
} = useScroll7TVUserEmotes(props.streamer.id, 20)

const loadTriggerOffset = 2000
const scrollTrigger = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null
const scrollContainer = ref<HTMLElement | null>(null)

onMounted(async () => {
  await nextTick()
  console.log(scrollContainer.value)
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
      // @ts-expect-error scroll-margin was not included in IntersectionObserverInit
      scrollMargin: `${loadTriggerOffset}px`,
      threshold: 0,
    },
  )
  if (scrollTrigger.value) {
    observer.observe(scrollTrigger.value)
  }
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
.user-emotes {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  background: var(--panel);
  border-radius: 10px;
  overflow: hidden;
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
</style>
