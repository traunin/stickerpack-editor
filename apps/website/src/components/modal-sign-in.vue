<template>
  <Teleport to="body">
    <div class="backdrop" v-show="model">
      <div class="modal">
        Sign in
        <button
          aria-label="close"
          @click="model = false"
        >
          x
        </button>
        <TelegramLogin />
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import TelegramLogin from './telegram-login.vue';
const model = defineModel<boolean>({default: false})
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
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 30px;
  border-radius: 10px;
  background: var(--background);
  font-size: 1.5em;
  color: var(--text);
}

button {
  border: none;
  color: var(--color);
  position: absolute;
  background: none;
  padding: 5px;
  height: 48px;
  aspect-ratio: 1;
  font-size: 1.2em;
  cursor: pointer;
  top: 10px;
  right: 10px;
  border-radius: 5px;
  transition: background-color 0.05s;
}

button:hover {
  background-color: var(--panel);
}
</style>
