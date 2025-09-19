<template>
  <div class="icon">
    <img v-if="loading" :src="loadingIcon" alt="Loading...">
    <img v-else-if="!error" :src="okIcon" alt="Success">
    <img v-else :src="errorIcon" alt="Error">
  </div>
</template>

<script setup lang="ts">
defineProps<{
  loading: boolean
  error: boolean
}>()

// preload
const loadingIcon = new URL('@/assets/icons/loading.png', import.meta.url).href
const okIcon = new URL('@/assets/icons/ok.png', import.meta.url).href
const errorIcon = new URL('@/assets/icons/error.png', import.meta.url).href
</script>

<style scoped>
img {
  width: 32px;
  height: 32px;
}
</style>
