<template>
  <div class="packs-paginated">
    <div class="back">
      <button :disabled="!hasPrevPage" @click="prev">
        &lt;
      </button>
    </div>
    <div v-if="loading" class="results loading">
      <LoadingAnimation />
    </div>

    <div v-else-if="noPacks" class="results">
      No packs were created
    </div>

    <div v-else-if="error" class="results">
      {{ error }}
    </div>

    <div v-else ref="container" class="results packs">
      <StickerpackPreview
        v-for="stickerpack in packs"
        :key="stickerpack.id"
        :stickerpack="stickerpack"
      />
    </div>

    <div class="forward">
      <button :disabled="!hasNextPage" @click="next">
        &gt;
      </button>
    </div>
  </div>
</template>

<script setup lang = "ts">
import { computed, nextTick, ref, watch } from 'vue'
import LoadingAnimation from '@/components/loading-animation.vue'
import StickerpackPreview from '@/components/stickerpack-preview.vue'
import { usePacksEndpoint } from '@/composables/use-packs-endpoint'
import { usePageSize } from '@/composables/use-page-size'

const container = ref<HTMLElement | null>(null)
const { pageSize, updatePageSize } = usePageSize(container)
const { packs, error, page, maxPages, next, prev } = usePacksEndpoint('public/packs', pageSize)

const foundPacks = computed(() => (packs.value?.length ?? 0) > 0)
const noPacks = computed(() => (packs.value?.length ?? 0) === 0)
const loading = computed(() => foundPacks.value == null && error.value == null)
const hasPrevPage = computed(() => page.value > 1)
const hasNextPage = computed(() => page.value < maxPages.value)

watch(packs, async (newPacks) => {
  if (newPacks && newPacks.length > 0) {
    await nextTick()
    updatePageSize()
  }
}, { immediate: true })
</script>

<style scoped>
.packs-paginated {
  display: flex;
  gap: 20px;
}

.results {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.packs {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-around;
}

button {
  height: 202px;
  padding: 10px;
  font-weight: 900;
  color: var(--text);
  background: var(--primary);
  font-size: 2em;
  border: none;
  font-family: "Roboto", serif;
  cursor: pointer;
}

button:disabled {
  color: red;
  cursor: default
}
</style>
