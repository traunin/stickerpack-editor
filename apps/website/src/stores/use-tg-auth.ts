import { defineStore } from 'pinia'
import { ref } from 'vue'
import { createSession, deleteSession, type User } from '@/api/session'

export const useTgAuthStore = defineStore('use-tg-auth', () => {
  const isLoggedIn = ref(false)
  const username = ref('')

  async function logIn(user: User) {
    if (await createSession(user)) {
      isLoggedIn.value = true
      username.value = user.username
    }
  }

  async function logOut() {
    if (await deleteSession()) {
      isLoggedIn.value = false
      username.value = ''
    }
  }

  // check if jwt is set?

  return { isLoggedIn, username, logIn, logOut }
}, {
  persist: {
    key: 'use-tg-auth',
    storage: localStorage,
  },
})
