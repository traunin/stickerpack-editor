<template>
  <button
    type="button"
    aria-label="Toggle dark mode"
    class="theme"
    @click="toggleDark()"
  >
    <div class="theme-icon" />
    <span>Theme</span>
  </button>
</template>

<script setup lang="ts">
import { useDark, useToggle } from '@vueuse/core'

const isDark = useDark()
const toggleDark = useToggle(isDark)
</script>

<style scoped>
.theme {
  cursor: pointer;
  flex: 1;
  display: flex;
  background: none;
  color: var(--text);
  border: none;
  font-size: 1.1em;
  align-items: center;
}

.theme-icon {
  display: inline-block;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  cursor: pointer;
  margin: 5px 7px 5px 3px;
  background: center url("@/assets/icons/sun.svg") white no-repeat;
  background-size: 20px 20px;
}

.dark .theme-icon {
  background: center url("@/assets/icons/moon.svg") black no-repeat;
  background-size: 20px 20px;
}
</style>
