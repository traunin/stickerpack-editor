<template>
  <div class="dropdown">
    <CreateLink class="create" @click="open = false" />
    <div class="user-info" @click="open = !open">
      <img src="@/assets/hi.gif" alt="">
      {{ authStore.username }}
      <img :src="authStore.photoURL" alt="" class="profile-picture">
    </div>
    <Transition name="fade">
      <ul v-if="open">
        <li>
          <DropdownThemeButton />
        </li>
        <li class="github-link">
          <DropdownGithub />
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
import CreateLink from '@/components/create-link.vue'
import DropdownLogoutButton from '@/components/dropdown-logout-button.vue'
import DropdownThemeButton from '@/components/dropdown-theme-button.vue'
import { useTgAuthStore } from '@/stores/use-tg-auth'
import DropdownGithub from './dropdown-github.vue'

const open = ref(false)
const authStore = useTgAuthStore()
</script>

<style scoped>
.profile-picture {
  width: 36px;
  height: 36px;
  border: 2px solid var(--primary);
  border-radius: 100%;
}

.dropdown {
  position: relative;
  display: flex;
  gap: 20px;
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
  border-radius: 10px;
  overflow: hidden;
}

li {
  display: flex;
  cursor: pointer;
}

.github-link {
  border-bottom: 1px solid var(--border);
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
