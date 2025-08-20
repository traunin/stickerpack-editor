<template>
  <div v-if="!authStore.isLoggedIn" class="unauthorized">
    Log in to see your packs
  </div>
  <div v-else class="packs-paginated">
    <div class="back">
      <button :disabled="!hasPrevPage" @click="prev">
        &lt;
      </button>
    </div>
    <div v-if="loading" class="results loading">
      <LoadingAnimation />
    </div>

    <div v-else-if="error" class="results">
      {{ error }}
    </div>

    <div v-else ref="container" class="results packs">
      <StickerpackPreview
        v-for="stickerpack in publicPacks"
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
import { useTgAuthStore } from '@/stores/use-tg-auth'

const authStore = useTgAuthStore()

const container = ref<HTMLElement | null>(null)
const { pageSize, updatePageSize } = usePageSize(container)
const { publicPacks, error, page, maxPages, next, prev } = usePacksEndpoint(
  'user/packs',
  pageSize,
  computed(() => authStore.isLoggedIn),
)

const foundPacks = computed(() => publicPacks.value?.length !== 0)
const loading = computed(() => !foundPacks.value && error.value == null)
const hasPrevPage = computed(() => page.value > 1)
const hasNextPage = computed(() => page.value < maxPages.value)

watch(publicPacks, async (newPacks) => {
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
  height: 200px;
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
  cursor:default
}

.unauthorized {
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 2em;
}
</style>
