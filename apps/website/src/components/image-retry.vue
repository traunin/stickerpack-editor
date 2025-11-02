<template>
  <div class="image-retry">
    <LoadingAnimation v-if="isLoading" class="loading" />
    <picture v-else>
      <img :src="currentURL" :alt="alt">
    </picture>
  </div>
</template>

<script setup lang="ts">
import { useImage } from '@vueuse/core'
import { ref, watch } from 'vue'
import LoadingAnimation from '@/components/loading-animation.vue'

interface Props {
  url: string
  alt?: string
  retries?: number
  baseDelay?: number
}

const props = withDefaults(defineProps<Props>(), {
  retries: 3,
  baseDelay: 200,
})

let retryCount = 0
const maxRetries = props.retries
const baseDelay = props.baseDelay
const currentURL = ref(`${props.url}?r=${retryCount}`)

const { isLoading, error } = useImage(ref({ src: currentURL }))

watch(error, (failed) => {
  if (failed && retryCount < maxRetries) {
    const delay = baseDelay * 2 ** retryCount + baseDelay * Math.random()
    retryCount++
    setTimeout(() => {
      currentURL.value = `${props.url}?r=${retryCount}`
    }, delay)
  }
})
</script>

<style scoped>
picture {
  height: 100%;
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

img {
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
