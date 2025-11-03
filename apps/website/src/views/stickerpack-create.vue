<template>
  <Transition>
    <ErrorMessage
      v-if="uploadError"
      :message="uploadError"
      class="error"
    />
  </Transition>
  <ModalProgress
    v-if="isUploading"
    :message="loadingMessage"
    :total="progress.total"
    :progress="progress.done"
  />
  <div class="creation-form">
    <PackFormSettings
      v-model="packParams"
      :sticker-count="stickerCount"
      :max-stickers="maxStickers"
      @name-error="nameError = $event"
      @title-error="titleError = $event"
    />
    <div class="stickers">
      <div class="sticker-search">
        <EmoteSource @sticker-selected="addSticker" />
      </div>
      <StickerListSelected v-model="stickers" />
    </div>
    <ButtonCreatePack :error="buttonError" @click="createPack" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, toRaw } from 'vue'
import { useRouter } from 'vue-router'
import type { ProgressEvent } from '@/api/stickerpack-upload'
import ButtonCreatePack from '@/components/button-create-pack.vue'
import EmoteSource from '@/components/emote-source.vue'
import ErrorMessage from '@/components/error-message.vue'
import ModalProgress from '@/components/modal-progress.vue'
import PackFormSettings from '@/components/pack-form-settings.vue'
import StickerListSelected from '@/components/sticker-list-selected.vue'
import { usePackValidation } from '@/composables/use-pack-validation'
import { useUploadPackMutation } from '@/composables/use-upload-pack-mutation'
import { useCreatedPackStore } from '@/stores/use-created-pack'
import type { PackParameters } from '@/types/pack'
import type { Sticker } from '@/types/sticker'

const packParams = ref<PackParameters>({
  name: '',
  title: '',
  hasWatermark: true,
  isPublic: true,
})
const stickers = ref<Sticker[]>([])

const nameError = ref<string | null>(null)
const titleError = ref<string | null>(null)
const progress = ref<ProgressEvent>({ done: 0, total: 0 })

const stickerCount = computed(() => stickers.value.length)
const maxStickers = 50 // 200 is not supported yet

const router = useRouter()
const createdPack = useCreatedPackStore()

const uploadPackMutation = useUploadPackMutation()

const isUploading = computed(() => uploadPackMutation.isPending.value)
const uploadError = computed(() => uploadPackMutation.uploadError.value)

const loadingMessage = computed(() => {
  if (progress.value.total === 0) {
    return 'Starting pack creation...'
  }
  return `Processing stickers (${progress.value.done}/${progress.value.total})`
})

const buttonError = usePackValidation(nameError, titleError, stickerCount, maxStickers)

function addSticker(sticker: Sticker) {
  if (stickerCount.value < maxStickers) {
    stickers.value.push(sticker)
  }
}

async function createPack() {
  if (buttonError.value) {
    return
  }

  progress.value = { done: 0, total: stickerCount.value }

  try {
    const response = await uploadPackMutation.mutateAsync({
      request: {
        pack_name: packParams.value.name,
        title: packParams.value.title,
        emotes: stickers.value.map(e => toRaw(e)),
        has_watermark: packParams.value.hasWatermark,
        is_public: packParams.value.isPublic,
      },
      onProgress: (progressEvent: ProgressEvent) => {
        progress.value = progressEvent
      },
    })

    createdPack.createdPack = response.pack
    router.push({
      name: 'packCreated',
      params: { id: response.pack.id },
    })
  } catch (error) {
    console.error(error)
  } finally {
    progress.value = { done: 0, total: 0 }
  }
}
</script>

<style scoped>
.creation-form {
  flex: 1;
  align-self: stretch;
  display: flex;
  flex-direction: column;
  min-height: 0;
  margin: 15px;
  gap: 15px;
}

.stickers {
  flex: 1;
  align-self: stretch;
  display: flex;
  align-items: stretch;
  min-height: 0;
  gap: 15px;
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
