import { defineStore } from 'pinia'
import { ref } from 'vue'
import { PackResponse } from '@/types/pack'

export const useCreatedPackStore = defineStore('created-pack', () => {
  const createdPack = ref<null | PackResponse>(null)

  return { createdPack }
})
