<template>
  <ErrorMessage
    :error="createError"
    :cleanup-timeout="4000"
  />
  <ModalProgress
    v-if="isCreating"
    :progress="progress"
  />

  <div class="creation-form">
    <PackParametersForm
      v-model="packParams"
      class="pack-parameters"
      :sticker-count="stickerCount"
      :max-stickers="maxStickers"
      @name-error="nameError = $event"
      @title-error="titleError = $event"
    />
    <div class="stickers">
      <div class="emote-source">
        <EmoteSource @sticker-selected="addSticker" />
      </div>
      <StickerListSelected v-model="stickers" />
    </div>

    <ButtonCreatePack
      class="create-button"
      :error="buttonError"
      @click="create"
    />

    <EmoteSourceModal
      :sticker-count="stickerCount"
      :max-stickers="maxStickers"
      @sticker-selected="addSticker"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import ButtonCreatePack from '@/components/button-create-pack.vue'
import EmoteSourceModal from '@/components/emote-source-modal.vue'
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

const { progress, isCreating, createError, createPack } = usePackCreation()

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
  flex-direction: column-reverse;
  padding: 10px;
  gap: 10px;
  overflow: auto;
}

.stickers {
  flex: 1;
  align-self: stretch;
  display: flex;
  align-items: stretch;
  gap: 15px;
}

.emote-source {
  display: none;
}

.selected-stickers {
  flex: 2;
  order: 3;
}

.create-button {
  order: 0;
}

.pack-parameters {
  order: 1;
}

@media screen and (min-width: 1000px) {
  .creation-form {
    flex: 1;
    align-self: stretch;
    display: flex;
    flex-direction: column;
    min-height: 0;
    margin: 15px;
    gap: 15px;
    min-height: 0;
  }

  .emote-source {
    flex: 1;
    display: flex;
    align-items: stretch;
    justify-content: center;
    min-height: 0;
  }

  .selected-stickers {
    order: 2;
  }

  .create-button {
    order: 3;
  }

  .pack-parameters {
    order: 0;
  }

  .stickers {
    min-height: 0;
  }
}
</style>
