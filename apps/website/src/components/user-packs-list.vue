<template>
  <div class="user-packs">
    <ModalLoading v-if="isDeleting" message="The stickerpack is deleting" />
    <ErrorMessage
      :error="deletionError"
      :cleanup-timeout="4000"
    />
    <div v-if="!authStore.isLoggedIn" class="unauthorized">
      Log in to see your packs
    </div>
    <div v-else class="packs-paginated">
      <div v-if="isLoading" class="results loading">
        <LoadingAnimation />
      </div>
      <div v-else-if="total === 0" class="results">
        You don't have any packs
      </div>
      <div v-else-if="isError" class="results">
        {{ error }}
      </div>
      <div v-else class="results packs">
        <div
          v-for="stickerpack in packs"
          :key="stickerpack.id"
          class="pack"
        >
          <StickerpackPreview
            :stickerpack="stickerpack"
          />
          <div class="delete" @click="confirmDelete(stickerpack.name)">
            X
          </div>
        </div>
        <div v-if="isFetchingNextPage" class="results loading">
          <LoadingAnimation />
        </div>
      </div>
      <div ref="scrollTrigger" style="height: 1px;" />
    </div>
    <ConfirmModal
      v-if="showConfirm"
      message="Are you sure you want to delete the pack?"
      @cancel="cancelDelete"
      @confirm="removePack"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import ConfirmModal from '@/components/confirm-modal.vue'
import ErrorMessage from '@/components/error-message.vue'
import LoadingAnimation from '@/components/loading-animation.vue'
import ModalLoading from '@/components/modal-loading.vue'
import StickerpackPreview from '@/components/stickerpack-preview.vue'
import { useDeletePackMutation } from '@/composables/use-delete-pack-mutation'
import { useScrollUserPacks } from '@/composables/use-scroll-user-packs'
import { useTgAuthStore } from '@/stores/use-tg-auth'

const showConfirm = ref(false)
const deletedPackName = ref('')

const authStore = useTgAuthStore()

const deletePackMutation = useDeletePackMutation()

const isDeleting = computed(() => deletePackMutation.isPending.value)
const deletionError = computed(() => deletePackMutation.deletionError)

function confirmDelete(name: string) {
  deletedPackName.value = name
  showConfirm.value = true
}

function cancelDelete() {
  showConfirm.value = false
  deletedPackName.value = ''
}

async function removePack() {
  try {
    await deletePackMutation.mutateAsync(deletedPackName.value)
  } catch (error) {
    console.error(error)
  } finally {
    showConfirm.value = false
    deletedPackName.value = ''
  }
}

const {
  packs,
  total,
  isLoading,
  isFetchingNextPage,
  hasMore,
  loadMore,
  isError,
  error,
} = useScrollUserPacks()
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

.pack {
  position: relative;
}

.delete {
  position: absolute;
  cursor: pointer;
  background: red;
  bottom: 2px;
  left: 2px;
  width: 32px;
  height: 32px;
  border-radius: 100%;
  font-size: 1.2em;
  display: flex;
  justify-content: center;
  align-items: center;
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

.unauthorized {
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 2em;
}

.error {
  position: fixed;
  top: 20px;
  left: 20px;
}

.v-enter-active,
.v-leave-active {
  transition: top 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  top: -15%;
}
</style>
