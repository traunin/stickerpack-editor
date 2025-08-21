<template>
  <div class="parameters">
    <input
      id="name"
      v-model="name"
      type="text"
      placeholder="Pack name used in links, english letters and digits"
    >
    <input
      id="title"
      v-model="title"
      type="text"
      placeholder="Displayed pack name"
    >
    <div class="watermark">
      <input id="watermark" v-model="watermark" type="checkbox" checked>
      <label for="watermark">Use a bot name watermark</label>
    </div>
    <div class="public">
      <input id="public" v-model="isPublic" type="checkbox" checked>
      <label for="public">Show to other users</label>
    </div>
    <div class="sticker-count">
      {{ stickerCount }} / {{ maxStickers }}
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  stickerCount: number
  maxStickers: number
}>()

const name = defineModel('name')
const title = defineModel('title')
const watermark = defineModel('watermark')
const isPublic = defineModel('isPublic')
</script>

<style scoped>
.parameters {
  display: flex;
  gap: 20px;
  margin-bottom: 0;
  background: var(--panel);
  padding: 10px;
}

.watermark, .public {
  flex: 1;
  font-size: 1.5em;
  display: flex;
  gap: 20px;
}

.watermark input, .public input {
  align-self: stretch;
  flex: 0;
  aspect-ratio: 1;
}

label {
  flex: 1;
  align-self: center;
}

input {
  flex: 1;
  background: var(--input);
  color: var(--text);
  border: 3px solid var(--primary);
  font-size: 1.3em;
  padding: 5px;
}

.sticker-count {
  font-size: 2em;
  align-self: center;
}
</style>
