<template>
  <ModalLoading v-if="isDeleting" message="The stickerpack is deleting" />
  <Transition>
    <ErrorMessage
      v-if="deletionError != null"
      :message="deletionError"
      class="error"
    />
  </Transition>

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

    <div v-else-if="noPacks" class="results">
      You don't have any packs
    </div>

    <div v-else-if="error" class="results">
      {{ error }}
    </div>

    <div v-else ref="container" class="results packs">
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
    </div>

    <div class="forward">
      <button :disabled="!hasNextPage" @click="next">
        &gt;
      </button>
    </div>
  </div>

  <ConfirmModal
    v-if="showConfirm"
    message="Are you sure you want to delete the pack?"
    @cancel="cancelDelete"
    @confirm="removePack"
  />
</template>

<script setup lang = "ts">
import { computed, nextTick, ref, watch } from 'vue'
import { deletePack } from '@/api/pack-delete'
import LoadingAnimation from '@/components/loading-animation.vue'
import StickerpackPreview from '@/components/stickerpack-preview.vue'
import { usePacksEndpoint } from '@/composables/use-packs-endpoint'
import { usePackEvents } from '@/composables/use-packs-events'
import { usePageSize } from '@/composables/use-page-size'
import { useTgAuthStore } from '@/stores/use-tg-auth'
import ConfirmModal from './confirm-modal.vue'
import ModalLoading from './modal-loading.vue'

const showConfirm = ref(false)
const deletedPackName = ref('')
const isDeleting = ref(false)
const deletionError = ref('')

const authStore = useTgAuthStore()

const container = ref<HTMLElement | null>(null)
const { pageSize, updatePageSize } = usePageSize(container)
const { emitPackEvent } = usePackEvents()

const { packs, error, page, maxPages, next, prev } = usePacksEndpoint(
  'user/packs',
  pageSize,
  computed(() => authStore.isLoggedIn),
)

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

function confirmDelete(name: string) {
  deletedPackName.value = name
  showConfirm.value = true
}

function cancelDelete() {
  showConfirm.value = false
  deletedPackName.value = ''
}

async function removePack() {
  isDeleting.value = true
  deletionError.value = ''
  try {
    await deletePack(deletedPackName.value)
    emitPackEvent('deleted', deletedPackName.value)
  } catch (err: unknown) {
    if (err instanceof Error) {
      deletionError.value = err.message
    } else {
      deletionError.value = String(err)
    }
    setTimeout(() => deletionError.value = '', 4000)
  } finally {
    isDeleting.value = false
  }
  showConfirm.value = false
  deletedPackName.value = ''
}
</script>

<style scoped>
.packs-paginated {
  display: flex;
  gap: 20px;
  flex: 1;
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
