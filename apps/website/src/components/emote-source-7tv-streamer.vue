<template>
  <div class="user-emotes">
    <LoadingAnimation v-if="isLoading" class="centered" />
    <div v-else-if="total === 0" class="centered">
      No emotes found
    </div>
    <div v-else-if="isError" class="centered">
      {{ error }}
    </div>
    <div v-else ref="scroll-container" class="results packs">
      <SearchResult7TV
        v-for="emote in emotes"
        :key="emote.id"
        :emote="emote"
        @click="selectEmote(emote)"
      />
      <div v-if="isFetchingNextPage" class="results centered">
        <LoadingAnimation />
      </div>
      <div ref="scroll-trigger" style="height: 1px; flex-shrink: 0;" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, useTemplateRef, watch } from 'vue'
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

.centered {
  margin: auto;
}
</style>
