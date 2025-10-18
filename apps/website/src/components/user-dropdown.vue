<template>
  <div class="dropdown">
    <div class="user-info" @click="open = !open">
      <img src="@/assets/hi.gif" alt="">
      {{ authStore.username }}
    </div>
    <Transition name="fade">
      <ul v-if="open">
        <li>
          <DropdownThemeButton />
        </li>
        <li>
          <DropdownLogoutButton />
        </li>
      </ul>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useTgAuthStore } from '@/stores/use-tg-auth'
import DropdownLogoutButton from './dropdown-logout-button.vue'
import DropdownThemeButton from './dropdown-theme-button.vue'

const open = ref(false)
const authStore = useTgAuthStore()
</script>

<style scoped>
.dropdown {
  position: relative;
}

.user-info {
  font-size: 1.2em;
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  border-radius: 10px;
  padding: 22px 15px;
  cursor: pointer;
  user-select: none;
}

.user-info:hover {
  background: var(--panel);
}

ul {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--input);
  margin-top: 8px;
  border: var(--border) 2px solid;
  width: 200px;
  align-items: stretch;
  z-index: 10;
}

li {
  display: flex;
  cursor: pointer;
}

li:first-child {
  border-bottom: 1px solid var(--border);
}

li:hover {
  background: var(--panel);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.1s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
