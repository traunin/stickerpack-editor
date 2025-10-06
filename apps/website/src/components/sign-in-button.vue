<template>
  <button @click="showModal=true">
    Sign in
  </button>
  <ModalSignIn v-model="showModal" />
  <ModalLoading v-if="authStore.isLoading" message="Signing you in" />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ModalSignIn from '@/components/modal-sign-in.vue';
import ModalLoading from './modal-loading.vue';
import { useTgAuthStore } from '@/stores/use-tg-auth';

const showModal = ref(false)
const authStore = useTgAuthStore()
</script>

<style scoped>
button {
  color: var(--text);
  border: none;
  background: var(--primary);
  padding: 10px 20px;
  font-size: 1.1em;
  cursor: pointer;
  border-radius: 10px;
  transition: background-color 0.01s;
}

button:hover {
  background-color: var(--accent);
}
</style>