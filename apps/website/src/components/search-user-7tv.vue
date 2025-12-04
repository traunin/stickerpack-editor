<template>
  <div class="user-search">
    <div v-if="debouncedQuery" class="search-response">
      <LoadingAnimation v-if="loading" class="offset" />
      <div v-else-if="error" class="results">
        {{ error }}
      </div>
      <div v-else-if="!usersFound" class="offset">
        No users found
      </div>
      <div v-else class="results users">
        <SearchResultUser7TV
          v-for="user in users"
          :key="user.id"
          :user="user"
          @click="selectUser(user)"
        />
      </div>
    </div>

    <input
      id="streamer"
      v-model="query"
      type="text"
      placeholder="Search 7TV users..."
    >
  </div>
</template>

<script setup lang="ts">
import { refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'
import LoadingAnimation from '@/components/loading-animation.vue'
import SearchResultUser7TV from '@/components/search-result-user-7tv.vue'
import { use7tvUserSearch } from '@/composables/use-7tv-user-search'
import type { User7TV } from '@/types/user-7tv'

const emit = defineEmits<{
  (e: 'source-select', user: User7TV): void
}>()

const query = ref('')
const debouncedQuery = refDebounced(query, 200)

const pageSize = 5
const { users, error, loading } = use7tvUserSearch(debouncedQuery, pageSize)
const usersFound = computed(() => users.value && users.value?.length > 0)

function selectUser(user: User7TV) {
  emit('source-select', user)
}
</script>

<style>
.user-search {
  position: relative;
}

.search-response {
  margin-bottom: 10px;
  background: var(--panel);
  display: flex;
  justify-content: center;
  border-radius: 10px;
  overflow: hidden;
}

.results {
  flex: 1;
}

input {
  width: 100%;
  padding: 10px;
  font-size: 1.2em;
  background: var(--input);
  outline: none;
  border: 1px solid transparent;
  color: var(--text);
  border-radius: 10px;
}

input:focus-visible {
  border: 1px solid var(--text);
}

.offset {
  margin: 10px;
}
</style>
