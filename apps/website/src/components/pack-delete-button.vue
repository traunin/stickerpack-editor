<template>
  <div class="delete-button">
    <button class="delete" @click="confirmDelete()">
      Delete pack
    </button>
    <ModalLoading v-if="isDeleting" message="The stickerpack is deleting" />
    <ErrorMessage
      :error="deletionError"
      :cleanup-timeout="4000"
    />
    <ConfirmModal
      v-if="showConfirm"
      message="Are you sure you want to delete the pack?"
      @cancel="cancelDelete"
      @confirm="removePack"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ErrorMessage from '@/components/error-message.vue'
import ConfirmModal from '@/components/modal-confirm.vue'
import ModalLoading from '@/components/modal-loading.vue'
import { usePackDeletion } from '@/composables/use-pack-deletion'

const props = defineProps<{
  packName: string
}>()

const showConfirm = ref(false)
const { isDeleting, deletionError, deletePack } = usePackDeletion()

function confirmDelete() {
  showConfirm.value = true
}

function cancelDelete() {
  showConfirm.value = false
}

async function removePack() {
  await deletePack(props.packName)
}
</script>

<style scoped>
.delete {
  cursor: pointer;
  background: red;
  border-radius: 10px;
  font-size: 1.2em;
  display: flex;
  justify-content: center;
  align-items: center;
  border: none;
  color: var(--text);
  border: 2px solid transparent;
  padding: 10px;
}

.delete:hover {
  border: 2px solid var(--border);
}

.error {
  position: fixed;
  top: 20px;
  left: 20px;
}

.v-enter-active,
.v-leave-active {
  transition: top 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  top: -15%;
}
</style>
