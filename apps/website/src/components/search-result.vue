<template>
  <div class="emote">
    <img :src="emote.preview" :alt="emote.name">
    <div class="name">
      {{ trimmed }}
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Emote } from '@/composables/use-emote-search'
import { useTrimmedString } from '@/composables/use-trimmed-string'

const props = defineProps<{
  emote: Emote
}>()

const { trimmed } = useTrimmedString(props.emote.name, 30)
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
  font-size: 1.1em;
}
</style>
