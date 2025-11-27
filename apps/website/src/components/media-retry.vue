<template>
  <div class="image-retry">
    <LoadingAnimation v-if="isLoading" class="loading" />
    <div v-else-if="error" class="error">
      Failed to load
    </div>
    <template v-else>
      <div v-if="data?.isVideo" class="video-container">
        <video
          :src="data?.url"
          :alt="alt"
          autoplay
          loop
          playsinline
          muted
        />
      </div>
      <picture v-else>
        <img :src="data?.url" :alt="alt">
      </picture>
    </template>
  </div>
</template>

<script setup lang="ts">
import LoadingAnimation from '@/components/loading-animation.vue'
import { useMedia } from '@/composables/use-media'

interface Props {
  url: string
  alt?: string
  retries?: number
}

const props = withDefaults(defineProps<Props>(), {
  retries: 3,
})

const { data, isLoading, error } = useMedia(props.url, props.retries)
</script>

<style scoped>
picture, .video-container {
  height: 100%;
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

img, video {
  object-fit: contain;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
}

.loading {
  margin: auto;
}

.image-retry {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
