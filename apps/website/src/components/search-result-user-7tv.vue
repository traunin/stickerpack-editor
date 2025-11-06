<template>
  <div class="user-7tv">
    <img :src="user.avatarURL" alt="" width="34" height="34">
    <span>
      {{ user.name }}
    </span>
  </div>
</template>

<script setup lang="ts">
import type { User7TV } from '@/types/user-7tv'

defineProps<{
  user: User7TV
}>()
</script>

<style scoped>
.user-7tv {
  display: flex;
  align-items: center;
  padding: 7px 10px;
  cursor: pointer;
}

.user-7tv:hover {
  background: var(--panel)
}

img {
  border: 2px solid var(--primary);
  border-radius: 100%;
  user-select: none;
}

span {
  flex: 1;
  text-align: right;
  user-select: none;
  margin-left: 20px;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}
</style>
