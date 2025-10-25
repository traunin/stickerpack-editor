<template>
  <div class="public-packs">
    <div v-if="isLoading" class="results loading">
      <LoadingAnimation />
    </div>
    <div v-else-if="total === 0" class="results">
      No packs were created
    </div>
    <div v-else-if="isError" class="results">
      {{ error }}
    </div>
    <div v-else class="results packs">
      <StickerpackPreview
        v-for="stickerpack in packs"
        :key="stickerpack.id"
        :stickerpack="stickerpack"
      />
      <div v-if="isFetchingNextPage" class="results loading">
        <LoadingAnimation />
      </div>
    </div>
    <div ref="scrollTrigger" style="height: 1px;" />
  </div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import LoadingAnimation from '@/components/loading-animation.vue'
import StickerpackPreview from '@/components/stickerpack-preview.vue'
import { useScrollPublicPacks } from '@/composables/use-scroll-public-packs'

const {
  packs,
  total,
  isLoading,
  isFetchingNextPage,
  hasMore,
  loadMore,
  isError,
  error,
} = useScrollPublicPacks()
const loadTriggerOffset = 500

const scrollTrigger = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

onMounted(async () => {
  await nextTick()
  observer = new IntersectionObserver(
    (entries) => {
      const target = entries[0]

      if (target.isIntersecting && hasMore.value && !isFetchingNextPage.value) {
        loadMore()
      }
    },
    {
      root: null,
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

watch([() => packs.value.length, isFetchingNextPage], ([length, fetching]) => {
  if (!fetching && length > 0) {
    setTimeout(checkAndLoadMore, 100)
  }
})

onBeforeUnmount(() => {
  if (observer) {
    observer.disconnect()
  }
})
</script>

<style scoped>
.results {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.packs {
  display: grid;
  grid-template-columns: repeat(auto-fill, 202px);
  gap: 15px;
}
</style>
