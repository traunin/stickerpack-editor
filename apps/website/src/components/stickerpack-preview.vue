<template>
  <div class="stickerpack-preview">
    <div class="title">
      {{ trimmed }}
    </div>
    <div class="preview">
      <img
        v-if="!isVideo && mediaURL"
        :src="mediaURL"
        alt="Stickerpack preview"
      >
      <video
        v-else-if="isVideo"
        :src="mediaURL ?? ''"
        autoplay
        loop
        playsinline
      />
    </div>
    <a :href="tgLink" class="tg-link">
      <img src="@/assets/icons/tglogo.png" alt="Telegram icon">
    </a>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { API_URL } from '@/api/config'
import type { PublicPack } from '@/api/public-packs'
import { useTrimmedString } from '@/composables/use-trimmed-string'

const props = defineProps<{ stickerpack: PublicPack }>()

const mediaURL = ref<string | null>(null)
const isVideo = ref(false)
const { trimmed } = useTrimmedString(props.stickerpack.title, 18)
const tgLink = `https://t.me/addstickers/${props.stickerpack.name}`

onMounted(async () => {
  const url = `${API_URL}/thumbnail?thumbnail_id=${props.stickerpack.thumbnail_id}`
  const res = await fetch(url)
  const blob = await res.blob()

  isVideo.value = blob.type === 'video/webm'
  mediaURL.value = URL.createObjectURL(blob)
})
</script>

<style scoped>
.stickerpack-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: var(--text);
  border: 2px solid var(--text);
  width: 200px;
  height: 200px;
  position: relative;
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
  font-size: 1.1em;
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
