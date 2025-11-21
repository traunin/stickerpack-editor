<template>
  <div class="dropdown">
    <MenuIcon class="menu-icon" @click="open = !open" />
    <Transition name="fade">
      <ul v-if="open">
        <li>
          <DropdownThemeButton />
        </li>
        <li class="github-link">
          <DropdownGithub />
        </li>
      </ul>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import MenuIcon from '@/assets/icons/burger.svg'
import DropdownGithub from './dropdown-github.vue'
import DropdownThemeButton from './dropdown-theme-button.vue'

const open = ref(false)
</script>

<style scoped>
.menu-icon {
  width: 42px;
  height: 42px;
  color: var(--text);
  padding: 5px;
  cursor: pointer;
  border-radius: 10px;
}

.menu-icon:hover {
  background-color: var(--panel);
}

.dropdown {
  position: relative;
  display: flex;
  gap: 20px;
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
  border-radius: 10px;
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
