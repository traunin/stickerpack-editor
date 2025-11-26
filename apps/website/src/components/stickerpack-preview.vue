<template>
  <a class="stickerpack-preview" :href="tgLink" target="_blank">
    <div class="title">
      {{ stickerpack.title }}
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
      <RouterLink v-if="isEditable" :to="`/edit/${stickerpack.id}`" class="edit">
        <EditIcon class="edit-icon" />
      </RouterLink>
    </div>
  </a>
</template>

<script setup lang="ts">
import EditIcon from '@/assets/icons/edit.svg'
import LoadingAnimation from '@/components/loading-animation.vue'
import { useThumbnail } from '@/composables/use-thumbnail'
import type { PackResponse } from '@/types/pack'

const props = defineProps<{
  stickerpack: PackResponse
  isEditable?: boolean
}>()

const tgLink = `https://t.me/addstickers/${props.stickerpack.name}`

const { data: thumbnailData, isLoading, error } = useThumbnail(props.stickerpack.thumbnail_id)
</script>

<style scoped>
.stickerpack-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: var(--text);
  border: 2px solid transparent;
  width: 202px;
  height: 202px;
  position: relative;
  background: var(--panel);
  border-radius: 10px;
  overflow: hidden;
  text-decoration: none;
}

.stickerpack-preview:hover {
  border: 2px solid var(--border-hover);
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
  text-wrap: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-decoration: none;
}

.edit {
  position: absolute;
  right: 5px;
  bottom: 5px;
  width: 32px;
  height: 32px;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--primary);
  border-radius: 100%
}

.edit-icon {
  width: 20px;
  height: 20px;
  color: var(--text);
}
</style>
