<template>
  <div class="stickerpack-preview">
    <div class="title">
      {{ trimmed }}
    </div>
    <div class="preview">
      <LoadingAnimation v-if="isLoading" />
      <div v-else-if="error" class="error">
        Failed to load
      </div>
      <img
        v-else-if="thumbnailData && !thumbnailData.isVideo"
        :src="thumbnailData.url"
        alt="Stickerpack preview"
      >
      <video
        v-else-if="thumbnailData?.isVideo"
        :src="thumbnailData.url"
        autoplay
        loop
        playsinline
      />
    </div>
    <a :href="tgLink" class="tg-link" target="_blank">
      <img src="@/assets/icons/tglogo.png" alt="Telegram icon" width="32" height="32">
    </a>
  </div>
</template>

<script setup lang="ts">
import LoadingAnimation from '@/components/loading-animation.vue'
import { useThumbnail } from '@/composables/use-thumbnail'
import { useTrimmedString } from '@/composables/use-trimmed-string'
import type { PackResponse } from '@/types/pack'

const props = defineProps<{ stickerpack: PackResponse }>()

const { trimmed } = useTrimmedString(props.stickerpack.title, 21)
const tgLink = `https://t.me/addstickers/${props.stickerpack.name}`

const { data: thumbnailData, isLoading, error } = useThumbnail(props.stickerpack.thumbnail_id)
</script>

<style scoped>
.stickerpack-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: var(--text);
  border: 2px solid var(--text);
  width: 202px;
  height: 202px;
  position: relative;
  background: var(--panel);
  border-radius: 10px;
  overflow: hidden;
}

.error {
  color: red;
}

.preview {
  width: 160px;
  height: 160px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview img, video {
  width: 100%;
  height: 100%;
  object-fit: contain;
  margin: 5px;
}

.title {
  background: var(--secondary);
  padding: 10px;
  font-size: 1em;
  text-align: center;
  width: 100%;
}

.tg-link {
  position: absolute;
  right: 0;
  bottom: 0;
  width: 32px;
  height: 32px;
}
</style>
