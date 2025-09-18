<template>
  <Transition>
    <ErrorMessage
      v-if="error != null"
      :message="error"
      class="error"
    />
  </Transition>
  <ModalLoading v-if="isLoading" :message="loadingMessage" />
  <div class="creation-form">
    <PackParameters
      v-model:name="name"
      v-model:title="title"
      v-model:watermark="watermark"
      v-model:is-public="isPublic"
      :sticker-count="stickerCount"
      :max-stickers="maxStickers"
      @name-error="nameError = $event"
      @title-error="titleError = $event"
    />
    <div class="stickers">
      <div class="sticker-search">
        <EmoteSource @sticker-selected="addSticker" />
      </div>
      <div class="selected-stickers">
        <draggable
          v-model="stickers"
          item-key="uniqueId"
          group="stickers"
          handle=".drag-handle"
          ghost-class="ghost-item"
          chosen-class="chosen-item"
          drag-class="drag-item"
          class="drag-area"
        >
          <template #item="{ index }">
            <div class="sticker-wrapper">
              <div class="drag-handle">
                ⋮⋮
              </div>
              <StickerCreate
                v-model="stickers[index]"
                @remove="removeEmote(index)"
              />
            </div>
          </template>
        </draggable>
      </div>
    </div>
    <button :disabled="!!buttonError" @click="createPack">
      Create {{ buttonError ? ` | ${buttonError}` : '' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, toRaw } from 'vue'
import { useRouter } from 'vue-router'
import draggable from 'vuedraggable'
import { type ProgressEvent, uploadPack } from '@/api/stickerpack-upload'
import EmoteSource from '@/components/emote-source.vue'
import ErrorMessage from '@/components/error-message.vue'
import ModalLoading from '@/components/modal-loading.vue'
import PackParameters from '@/components/pack-parameters.vue'
import StickerCreate from '@/components/sticker-create.vue'
import { useCreatedPackStore } from '@/stores/use-created-pack'
import { useTgAuthStore } from '@/stores/use-tg-auth'
import type { Sticker } from '@/types/sticker'

const title = ref<string>('')
const name = ref<string>('')
const watermark = ref<boolean>(true)
const isPublic = ref<boolean>(true)
const stickers = ref<Sticker[]>([])

const nameError = ref<string | null>(null)
const titleError = ref<string | null>(null)

const stickerCount = computed(() => stickers.value.length)
const maxStickers = 50 // 200 is not supported yet

const router = useRouter()
const createdPack = useCreatedPackStore()
const authStore = useTgAuthStore()

const isLoading = ref(false)
const error = ref<string | null>(null)
const progress = ref<ProgressEvent>({ done: 0, total: 0 })

const loadingMessage = computed(() => {
  if (progress.value.total === 0) {
    return 'Starting pack creation...'
  }
  return `Processing stickers (${progress.value.done}/${progress.value.total})`
})

const buttonError = computed(() => {
  if (!authStore.isLoggedIn)
    return 'You are not logged in'
  if (nameError.value)
    return nameError.value
  if (titleError.value)
    return titleError.value
  if (stickerCount.value === 0)
    return 'No emotes selected'
  if (stickerCount.value > maxStickers)
    return `Too many emotes (max ${maxStickers})`
  return null
})

function addSticker(sticker: Sticker) {
  if (stickerCount.value < maxStickers) {
    stickers.value.push(sticker)
  }
}

function removeEmote(index: number) {
  stickers.value.splice(index, 1)
}

async function createPack() {
  if (buttonError.value) {
    return
  }

  isLoading.value = true
  error.value = null
  progress.value = { done: 0, total: stickerCount.value }

  try {
    const response = await uploadPack({
      pack_name: name.value,
      title: title.value,
      emotes: stickers.value.map(e => toRaw(e)),
      has_watermark: watermark.value,
      is_public: isPublic.value,
    }, (progressEvent: ProgressEvent) => {
      progress.value = progressEvent
    })

    createdPack.createdPack = response.pack
    router.push({
      name: 'packCreated',
      params: { id: response.pack.id },
    })
  } catch (err: unknown) {
    if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = String(err)
    }
    setTimeout(() => error.value = null, 4000)
  } finally {
    isLoading.value = false
    progress.value = { done: 0, total: 0 }
  }
}
</script>

<style scoped>
button {
  background: var(--primary);
  color: var(--text);
  font-size: 1.5em;
  padding: 10px;
  color: var(--text);
  border: none;
  margin-top: 0;
  cursor: pointer
}

.creation-form {
  flex: 1;
  align-self: stretch;
  display: flex;
  flex-direction: column;
  min-height: 0;
  margin: 20px;
  gap: 20px;
}

.stickers {
  flex: 1;
  align-self: stretch;
  display: flex;
  align-items: stretch;
  min-height: 0;
  gap: 20px;
}

.sticker-search {
  flex: 1;
  display: flex;
  align-items: stretch;
  justify-content: center;
  min-height: 0;
}

.selected-stickers {
  flex: 2;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: flex-start;
  align-self: flex-start;
  min-height: 0;
  overflow-y: auto;
  align-self: stretch;
  background: var(--panel);
  scrollbar-color: var(--accent) var(--input);
  scrollbar-width: thin;
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

button:disabled {
  background: grey;
  cursor: default;
}

.drag-area > :first-child::before {
  content: "Preview";
  position: absolute;
  top: 5px;
  height: 22.9px;
  line-height: 22.9px;
  padding: 0 4px;
  background: var(--primary);
  color: var(--text);
  z-index: 5;
  left: 50%;
  transform: translate(-50%);
}

.drag-area {
  display: grid;
  grid-template-columns: repeat(auto-fill, 202px);
  flex: 1;
  gap: 10px;
  margin: 20px;
  justify-content: center;
}

.sticker-wrapper {
  position: relative;
  padding: 5px;
  background: var(--panel);
}

.drag-handle {
  position: absolute;
  width: 20px;
  height: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  color: var(--text);
  background: var(--primary);
  cursor: grab;
  z-index: 10;
}
</style>
