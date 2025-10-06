import { defineStore } from 'pinia'
import { ref } from 'vue'
import { createSession, deleteSession, type User } from '@/api/session'

export const useTgAuthStore = defineStore('use-tg-auth', () => {
  const isLoggedIn = ref(false)
  const username = ref('')
  const isLoading = ref(false)

  async function logIn(user: User) {
    isLoading.value = true
    if (await createSession(user)) {
      isLoggedIn.value = true
      username.value = user.username
    }
    isLoading.value = false
  }

  async function logOut() {
    isLoading.value = true
    try {
      const success = await deleteSession()
      if (success) {
        isLoggedIn.value = false
        username.value = ''
      }
    } catch (err) {
      console.log(err)
    } finally {
      isLoading.value = false
    }
  }

  // check if jwt is set?

  return { isLoggedIn, username, logIn, logOut, isLoading }
}, {
  persist: {
    key: 'use-tg-auth',
    storage: localStorage,
    pick: ['username', 'isLoggedIn']
  },
})
