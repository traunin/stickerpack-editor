<template>
  <div class="emote">
    <img :src="emote.preview" :alt="emote.name">
  </div>
</template>

<script setup lang="ts">
import type { Emote } from '@/types/sticker'

defineProps<{
  emote: Emote
}>()
</script>

<style scoped>
.emote {
  height: 96px;
  display: flex;
  cursor: pointer;
  align-items: center;
  justify-content: center;
  background: var(--panel);
}

.emote:hover {
  background: var(--accent);
}

img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
