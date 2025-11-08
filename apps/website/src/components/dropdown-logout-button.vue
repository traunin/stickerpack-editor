<template>
  <button
    class="logout-button"
    @click="authStore.logOut"
  >
    <div class="logout" />
    <span>Sign out</span>
    <LoadingAnimation v-if="authStore.isLoading" class="loading" />
  </button>
</template>

<script setup lang="ts">
import { useTgAuthStore } from '@/stores/use-tg-auth'
import LoadingAnimation from './loading-animation.vue'

const authStore = useTgAuthStore()
</script>

<style scoped>
.logout-button {
  color: var(--text);
  background: none;
  border: none;
  cursor: pointer;
  border-radius: 0 10px 10px 0;
  display: flex;
  align-items: center;
  font-size: 1.1em;
  flex: 1;
  user-select: none;
}

.logout {
  mask: no-repeat center url("@/assets/icons/logout.svg");
  background-color: var(--text);
  height: 40px;
  width: 40px;
  margin-right: 4px;
}

.loading {
  width: 32px;
  height: 32px;
  margin-right: 5px;
}

span {
  flex: 1;
  text-align: left;
}
</style>
