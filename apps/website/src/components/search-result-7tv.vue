<template>
  <div class="emote">
    <ImageRetry class="image" :url="emote.preview" :alt="emote.name" />
    <div class="name">
      {{ trimmed }}
    </div>
  </div>
</template>

<script setup lang="ts">
import ImageRetry from '@/components/image-retry.vue'
import { useTrimmedString } from '@/composables/use-trimmed-string'
import type { Emote } from '@/types/sticker'

const props = defineProps<{
  emote: Emote
}>()

const trimmed = useTrimmedString(props.emote.name, 25)
</script>

<style scoped>
.emote {
  display: flex;
  cursor: pointer;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  height: 96px;
  flex-shrink: 0;
}

.emote:hover {
  background: var(--accent) !important;
}

.emote:nth-child(even) {
  background: var(--panel);
}

.name {
  margin-right: 20px;
  font-size: 1em;
  flex: 1;
  text-align: right;
}

.image {
  width: 96px;
  height: 96px;
}
</style>
