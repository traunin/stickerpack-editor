<template>
    <div class="stickerpacks">
        <stickerpack-preview
            v-for="stickerpack in stickerpacks"
            :name="stickerpack['name']"
        ></stickerpack-preview>
    </div>
</template>

<script setup lang = "ts">
import StickerpackPreview from '@/components/StickerpackPreview.vue';
import { ref } from 'vue'

type Stickerpack = {
    name: string
}

const loading = ref(false)
const stickerpacks = ref<null | Stickerpack[]>(null)

fetchStickerpacks();

async function fetchStickerpacks() {
  loading.value = true
  
  try {
    stickerpacks.value = 
        await fetch("http://7tvstickerpacks.website:8080/api/stickerpacks")  
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