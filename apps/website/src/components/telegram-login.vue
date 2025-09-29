<template>
  <div
    v-show="!authStore.isLoggedIn"
    ref="telegramWidget"
    class="tg-auth-wrapper"
  />

  <div v-show="authStore.isLoggedIn" class="user-info">
    <div class="user-info">
      <img src="@/assets/hi.gif" alt="">
      {{ authStore.username }}
    </div>
    <button @click="authStore.logOut">
      <div class="logout" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { User } from '@/api/session'
import { useTgAuthStore } from '@/stores/use-tg-auth'

declare global {
  interface Window {
    onAuth: (user: User) => void
  }
}

const telegramWidget = ref<null | HTMLElement>(null)

const authStore = useTgAuthStore()

onMounted(() => {
  window.onAuth = authStore.logIn

  const script = document.createElement('script')

  script.setAttribute('data-telegram-login', import.meta.env.VITE_BOT_NAME)
  script.setAttribute('data-size', 'large')
  script.setAttribute('data-userpic', 'false')
  script.setAttribute('data-radius', '10')
  script.setAttribute('data-onauth', 'onAuth(user)')
  script.setAttribute('data-request-access', 'write')
  script.src = 'https://telegram.org/js/telegram-widget.js?22'
  script.async = true

  telegramWidget.value!.appendChild(script)
})
</script>

<style scoped>
.tg-auth-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
}

.user-info {
  font-size: 1.2em;
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--accent);
  height: 40px;
  border-radius: 10px;
}

img {
  margin-left: 15px;
}

button {
  background: var(--primary);
  border: none;
  cursor: pointer;
  border-radius: 0 10px 10px 0;
}

.logout {
  mask: no-repeat center url("@/assets/icons/logout.svg");
  background-color: var(--text);
  height: 40px;
  width: 60px;
}
</style>
