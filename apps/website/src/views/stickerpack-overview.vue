<template>
  <div class="stickerpacks">
    <AddStickerpack />
    <StickerpackPreview
      v-for="stickerpack in stickerpacks"
      :key="stickerpack.name"
      :name="stickerpack.name"
    />
  </div>
</template>

<script setup lang = "ts">
import { ref } from 'vue'
import AddStickerpack from '@/components/create-link.vue'
import StickerpackPreview from '@/components/stickerpack-preview.vue'

interface Stickerpack {
  name: string
}

const loading = ref(false)
const stickerpacks = ref<null | Stickerpack[]>(null)

fetchStickerpacks()

async function fetchStickerpacks() {
  loading.value = true

  try {
    stickerpacks.value =
      await fetch('http://localhost:8080/api/stickerpacks')
        .then(res => res.json())
  } catch (err) {
    console.log(err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.stickerpacks {
  margin: 20px;
  display: flex;
  gap: 20px;
}
</style>
