<template>
  <Teleport to="body">
    <div class="backdrop">
      <div class="modal">
        <span>
          {{ message }}
        </span>
        <div class="actions">
          <button class="confirm" @click="confirm">
            Yes
          </button>
          <button class="cancel" @click="cancel">
            No
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
defineProps<{
  message: string
}>()
const emit = defineEmits<{
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()

function confirm() {
  emit('confirm')
}

function cancel() {
  emit('cancel')
}
</script>

<style scoped>
.backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 20px;
  border: 3px solid var(--primary);
  background: var(--background);
  font-size: 1.5em;
  color: var(--text);
  border-radius: 10px;
}

.actions {
  display: flex;
  gap: 20px;
}

button {
  color: var(--text);
  font-size: 1.2em;
  padding: 5px 20px;
  border: none;
  cursor: pointer;
  border-radius: 10px;
  border: 2px solid transparent
}

button:hover {
  border-color: var(--border);
}

.cancel {
  background: red
}

.confirm {
  background: green
}
</style>
