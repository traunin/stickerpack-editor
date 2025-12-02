<template>
  <div class="created-pack">
    <StickerpackPreview
      v-if="createdPack"
      :stickerpack="createdPack"
      :is-editable="true"
    />
    <div v-else>
      <div class="create-link">
        You've not created a pack yet
        <CreateLink />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import CreateLink from '@/components/create-link.vue'
import StickerpackPreview from '@/components/stickerpack-preview.vue'
import { useCreatedPackStore } from '@/stores/use-created-pack'

const { createdPack } = useCreatedPackStore()
</script>

<style scoped>
.created-pack {
  display: flex;
  flex: 1;
  align-self: stretch;
  justify-content: center;
  align-items: center;
}

.create-link {
  display: flex;
  flex-direction: column;
  gap: 20px;
  font-size: 1.2em;
  padding: 30px;
  background: var(--panel)
}
</style>
