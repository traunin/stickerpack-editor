<template>
  <div class="creation-form">
    <PackParameters
      v-model:name="name"
      v-model:title="title"
      v-model:watermark="watermark"
      v-model:is-public="isPublic"
      :sticker-count="stickerCount"
      :max-stickers="maxStickers"
    />
    <div class="stickers">
      <div class="sticker-search">
        <EmoteSearch @emote-selected="addEmote" />
      </div>
      <div class="selected-stickers">
        <StickerCreate
          v-for="(emote, index) in stickers"
          :key="emote.id"
          v-model="stickers[index]"
          @remove="removeEmote(index)"
        />
      </div>
    </div>
    <button @click="createPack">
      Create
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, toRaw } from 'vue'
import { type Sticker, uploadPack } from '@/api/stickerpack-upload'
import EmoteSearch from '@/components/emote-search.vue'
import PackParameters from '@/components/pack-parameters.vue'
import StickerCreate from '@/components/sticker-create.vue'
import type { Emote } from '@/composables/use-emote-search'

const title = ref<string>('')
const name = ref<string>('')
const watermark = ref<boolean>(true)
const isPublic = ref<boolean>(true)
const stickers = ref<Sticker[]>([])
const stickerCount = computed(() => stickers.value.length)
const maxStickers = 50 // 200 is not supported yet

function addEmote(emote: Emote) {
  if (stickerCount.value < maxStickers) {
    stickers.value.push({ ...emote, emoji_list: ['😀'], source: '7tv' })
  }
}

function removeEmote(index: number) {
  stickers.value.splice(index, 1)
}

function createPack() {
  uploadPack({
    pack_name: name.value,
    title: title.value,
    emotes: stickers.value.map(e => toRaw(e)), // unwrapping for... reasons... emotes.value doesn't work
    has_watermark: watermark.value,
    is_public: isPublic.value,
  })
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
  padding: 10px;
  scrollbar-color: var(--accent) var(--input);
  scrollbar-width: thin;
}
</style>
