<template>
  <ErrorMessage
    :error="uploadError"
    :cleanup-timeout="4000"
  />
  <ModalProgress
    v-if="isUploading"
    :progress="progress"
  />

  <div class="creation-form">
    <PackParametersForm
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

    <ButtonCreatePack :error="buttonError" @click="create" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import ButtonCreatePack from '@/components/button-create-pack.vue'
import EmoteSource from '@/components/emote-source.vue'
import ErrorMessage from '@/components/error-message.vue'
import ModalProgress from '@/components/modal-progress.vue'
import PackParametersForm from '@/components/pack-parameters-form.vue'
import StickerListSelected from '@/components/sticker-list-selected.vue'
import { usePackCreation } from '@/composables/use-pack-creation'
import { usePackValidation } from '@/composables/use-pack-validation'
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
const stickerCount = computed(() => stickers.value.length)
const maxStickers = 50

const { progress, isUploading, uploadError, createPack } = usePackCreation()

const buttonError = usePackValidation(
  nameError,
  titleError,
  stickerCount,
  maxStickers,
)

function addSticker(sticker: Sticker) {
  if (stickerCount.value < maxStickers) {
    stickers.value.push(sticker)
  }
}

async function create() {
  if (buttonError.value)
    return
  await createPack(packParams.value, stickers.value)
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
</style>
