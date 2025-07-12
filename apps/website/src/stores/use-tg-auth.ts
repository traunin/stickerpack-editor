import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface User {
  auth_date: number
  first_name: string
  hash: string
  id: number
  photo_url: string
  username: string
}

export const useTgAuthStore = defineStore('use-tg-auth', () => {
  const isLoggedIn = ref(false)
  // Only store necessary info
  const id = ref(-1)
  const username = ref('')
  const hash = ref('')

  function logIn(user: User) {
    isLoggedIn.value = true
    id.value = user.id
    username.value = user.username
    hash.value = user.hash
  }

  function logOut() {
    isLoggedIn.value = false
    id.value = -1
    username.value = ''
    hash.value = ''
  }

  return { isLoggedIn, id, username, hash, logIn, logOut }
}, {
  persist: {
    key: 'use-tg-auth',
    storage: localStorage,
  },
})
