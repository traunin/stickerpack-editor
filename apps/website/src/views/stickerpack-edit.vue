<template>
  <ErrorMessage
    :error="editError"
    :cleanup-timeout="4000"
  />
  <ModalProgress
    v-if="isEditing"
    :progress="progress"
  />

  <div v-if="!isLoading" class="edit-form">
    <div class="pack-parameters">
      <PackParametersForm
        v-model="packParams"
        class="params-form"
        :sticker-count="stickerCount"
        :max-stickers="maxStickers"
        :is-editing="true"
        @name-error="nameError = $event"
        @title-error="titleError = $event"
      />
      <PackDeleteButton :pack-name="props.name" class="delete-button" />
    </div>
    <div class="stickers">
      <div class="emote-source">
        <EmoteSource @sticker-selected="addSticker" />
      </div>
      <StickerListSelected v-model="stickers" />
    </div>

    <ButtonCreatePack
      class="edit-button"
      :error="buttonError"
      :disabled="!hasEdits"
      @click="edit"
    >
      Edit
    </ButtonCreatePack>

    <EmoteSourceModal
      :sticker-count="stickerCount"
      :max-stickers="maxStickers"
      @sticker-selected="addSticker"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, toRaw, watch } from 'vue'
import ButtonCreatePack from '@/components/button-create-pack.vue'
import EmoteSourceModal from '@/components/emote-source-modal.vue'
import EmoteSource from '@/components/emote-source.vue'
import ErrorMessage from '@/components/error-message.vue'
import ModalProgress from '@/components/modal-progress.vue'
import PackDeleteButton from '@/components/pack-delete-button.vue'
import PackParametersForm from '@/components/pack-parameters-form.vue'
import StickerListSelected from '@/components/sticker-list-selected.vue'
import { useLoadPack } from '@/composables/use-load-pack'
import { usePackDiff } from '@/composables/use-pack-diff'
import { usePackEdit } from '@/composables/use-pack-edit'
import { usePackValidation } from '@/composables/use-pack-validation'
import type { PackParameters } from '@/types/pack'
import type { Sticker } from '@/types/sticker'

const props = defineProps<{
  name: string
}>()

const packParams = ref<PackParameters>({
  title: '',
  isPublic: true,
})
const stickers = ref<Sticker[]>([])
const originalPack = {
  title: '',
  isPublic: true,
  stickers: [] as Sticker[],
}

const { data: pack, isLoading } = useLoadPack(props.name)

watch(pack, value => {
  if (!value) {
    return
  }

  packParams.value = {
    title: value.title,
    isPublic: value.isPublic,
  }

  stickers.value = value.stickers.map(sticker => ({ ...toRaw(sticker) }))
  originalPack.title = value.title
  originalPack.isPublic = value.isPublic
  originalPack.stickers = value.stickers.map(sticker => ({ ...toRaw(sticker) }))
}, { immediate: true })

const nameError = ref<string | null>(null)
const titleError = ref<string | null>(null)
const stickerCount = computed(() => stickers.value?.length ?? 0)
const maxStickers = 120

const { edits, hasEdits } = usePackDiff(
  stickers,
  packParams,
  originalPack,
)

const { progress, isEditing, editError, editPack } = usePackEdit()

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

async function edit() {
  if (buttonError.value)
    return
  await editPack(props.name, edits)
}
</script>

<style scoped>
.edit-form {
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

.edit-button {
  order: 0;
}

.pack-parameters {
  display: flex;
  gap: 10px;
  order: 1;
  align-items: center;
  flex-direction: column;
  align-items: stretch;
}

.params-form {
  flex: 1;
}

.delete-button {
  align-self: center;
}

@media screen and (min-width: 1000px) {
  .edit-form {
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

  .edit-button {
    order: 3;
  }

  .pack-parameters {
    order: 0;
    flex-direction: row;
    flex: 0;
  }

  .stickers {
    min-height: 0;
  }
}
</style>
