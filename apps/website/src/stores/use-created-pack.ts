import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { PackPreview } from '@/types/pack'

export const useCreatedPackStore = defineStore('created-pack', () => {
  const createdPack = ref<null | PackPreview>(null)

  return { createdPack }
})
