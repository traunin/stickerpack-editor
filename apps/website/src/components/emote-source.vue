<template>
  <div class="emote-search-selector">
    <div v-if="state === 'menu'" class="menu-header">
      Select emote source
    </div>
    <div v-else class="selected-source">
      <div class="source-name">
        <template v-if="typeof state === 'string'">
          {{ state }}
        </template>
        <template v-else-if="state.type === '7tvstreamer'">
          7TV user – {{ trimmedStreamerName }}
        </template>
      </div>
      <button @click="state = 'menu'">
        Back
      </button>
    </div>
    <component
      :is="currentView"
      class="source"
      v-bind="currentViewProps"
      @source-select="onSourceSelect"
      @sticker-selected="onStickerSelected"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import EmoteSearchTenor from '@/components/emote-search-tenor.vue'
import EmoteSearch7tv from '@/components/emote-source-7tv-search.vue'
import EmoteSource7TVStreamer from '@/components/emote-source-7tv-streamer.vue'
import EmoteSourceMenu from '@/components/emote-source-menu.vue'
import { useTrimmedString } from '@/composables/use-trimmed-string'
import type { Sticker } from '@/types/sticker'
import type { User7TV } from '@/types/user-7tv'

const emit = defineEmits<{
  (e: 'sticker-selected', sticker: Sticker): void
}>()

type State = 'menu' | '7TV' | 'Tenor' | { type: '7tvstreamer', streamer: User7TV }
const state = ref<State>('menu')
const views = {
  'menu': EmoteSourceMenu,
  '7TV': EmoteSearch7tv,
  'Tenor': EmoteSearchTenor,
  '7tvstreamer': EmoteSource7TVStreamer,
}

const currentView = computed(() => {
  if (typeof state.value === 'string') {
    return views[state.value]
  }
  return views[state.value.type]
})

const currentViewProps = computed(() => {
  if (typeof state.value === 'object' && state.value.type === '7tvstreamer')
    return { streamer: state.value.streamer }
  return {}
})

const trimmedStreamerName = computed(() => {
  if (typeof state.value === 'object' && state.value.type === '7tvstreamer') {
    return useTrimmedString(state.value.streamer.name, 15)
  }
  return ''
})

function onStickerSelected(sticker: Sticker) {
  emit('sticker-selected', sticker)
}

function onSourceSelect(source: string, user?: User7TV) {
  if (source === '7tvstreamer' && user) {
    state.value = { type: '7tvstreamer', streamer: user }
  } else {
    state.value = source as State
  }
}
</script>

<style scoped>
.menu-header {
  padding: 20px;
  font-size: 1.2em;
  text-align: center;
}

.emote-search-selector {
  background: var(--panel);
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  height: 100%;
  flex: 1;
}

button {
  padding: 5px 20px;
  cursor: pointer;
  border-radius: 10px;
  border: none;
  background: var(--input);
  color: var(--text);
  font-size: 1em;
}

.source {
  flex: 1;
  margin: 10px;
}

.selected-source {
  display: flex;
  font-size: 1.2em;
  padding: 10px;
  align-items: center;
  padding-bottom: 0;
}

.source-name {
  flex: 1;
}
</style>
