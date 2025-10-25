import { defineStore } from 'pinia'
import { ref } from 'vue'
import { createSession, deleteSession } from '@/api/session'
import type { User } from '@/api/session'

export const useTgAuthStore = defineStore('use-tg-auth', () => {
  const isLoggedIn = ref(false)
  const username = ref('')
  const photoURL = ref('')
  const isLoading = ref(false)

  async function logIn(user: User) {
    isLoading.value = true
    if (await createSession(user)) {
      isLoggedIn.value = true
      username.value = user.username
      photoURL.value = user.photo_url
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
        photoURL.value = ''
      }
    } catch (err) {
      console.log(err)
    } finally {
      isLoading.value = false
    }
  }

  // check if jwt is set?

  return { isLoggedIn, username, photoURL, logIn, logOut, isLoading }
}, {
  persist: {
    key: 'use-tg-auth',
    storage: localStorage,
    pick: ['username', 'photoURL', 'isLoggedIn'],
  },
})
