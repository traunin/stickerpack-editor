<template>
  <div class="emote">
    <img :src="emote.preview" :alt="emote.name">
    <div class="name">
      {{ trimmedName }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Emote } from '@/composables/use-emote-search'

const props = defineProps<{
  emote: Emote
}>()

const MAX_NAME_LENGTH = 30
const TRIMMED_LENGTH = 27
const trimmedName = computed(() =>
  props.emote.name.length < MAX_NAME_LENGTH ?
    props.emote.name :
    `${props.emote.name.substring(0, TRIMMED_LENGTH)}...`,
)
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
