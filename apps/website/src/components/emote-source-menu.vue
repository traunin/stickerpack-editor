<template>
  <nav class="sources">
    <ul>
      <li>
        <UserSearch7TV class="user-search" @source-select="selectUser" />
      </li>
      <li>
        <button @click="$emit('source-select', '7TV')">
          7TV emote search
        </button>
      </li>
      <li>
        <button @click="$emit('source-select', 'Tenor')">
          Tenor
        </button>
      </li>
    </ul>
  </nav>
</template>

<script setup lang="ts">
import UserSearch7TV from '@/components/search-user-7tv.vue'
import type { User7TV } from '@/types/user-7tv'

const emit = defineEmits<{
  (e: 'source-select', source: string, user7tv?: User7TV): void
}>()

function selectUser(user: User7TV) {
  emit('source-select', '7tvstreamer', user)
}
</script>

<style scoped>
.sources {
  display: flex;
  justify-content: center;
  align-items: center;
}

ul {
  border-radius: 10px;
  border: 2px solid var(--border);
  background: var(--panel);
}

button {
  padding: 10px;
  cursor: pointer;
  width: 100%;
  border-radius: 0;
  border: none;
  color: var(--text);
  background: none;
  font-size: 1.2em;
}

li {
  border-top: 1px solid var(--border);
  border-bottom: 1px solid var(--border);
}

li:hover {
  background: var(--panel);
}

li:first-child {
  border-top: none;
  border-radius: 10px 10px 0 0;
}

li:first-child:hover {
  background: none;
}

li:last-child {
  border-bottom: none;
  border-radius: 10px 10px 0 0;
  border-radius: 0 0 10px 10px;
}

.user-search {
  padding: 10px;
}
</style>
