<template>
  <a class="stickerpack-preview" :href="tgLink" target="_blank">
    <div class="title">
      {{ stickerpack.title }}
    </div>
    <ImageRetry class="preview" :url="thumbnailUrl" />
    <RouterLink v-if="isEditable" :to="`/edit/${stickerpack.name}`" class="edit">
      <EditIcon class="edit-icon" />
    </RouterLink>
  </a>
</template>

<script setup lang="ts">
import { API_URL } from '@/api/config'
import EditIcon from '@/assets/icons/edit.svg'
import ImageRetry from '@/components/media-retry.vue'
import type { PackPreview } from '@/types/pack'

const props = defineProps<{
  stickerpack: PackPreview
  isEditable?: boolean
}>()

const tgLink = `https://t.me/addstickers/${props.stickerpack.name}`
const thumbnailUrl = `${API_URL}/media?file_id=${props.stickerpack.thumbnail_id}`
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

.stickerpack-preview:has(.edit:hover) {
  border-color: transparent;
}

.preview {
  width: 160px;
  height: 160px;
  display: flex;
  align-items: center;
  justify-content: center;
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
  border-radius: 100%;
  border: 2px solid transparent;
}

.edit:hover {
  border-color: var(--border-hover);
}

.edit-icon {
  width: 20px;
  height: 20px;
  color: var(--text);
}
</style>
