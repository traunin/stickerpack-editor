<template>
  <div class="emote">
    <img :src="emote.preview" :alt="emote.name">
    <div class="name">
      {{ emote.name }}
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Emote } from '@/composables/use-emote-search'

defineProps<{
  emote: Emote
}>()
</script>

<style scoped>
.emote {
  display: flex;
  cursor: pointer;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.emote:hover {
  background: var(--accent);
}

.name {
  margin-right: 20px;
}
</style>
