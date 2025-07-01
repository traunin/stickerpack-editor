<script setup lang="ts">
import { useDark, useToggle } from '@vueuse/core'

const isDark = useDark()
const toggleDark = useToggle(isDark)
</script>

<template>
  <button
    type="button"
    aria-label="Toggle dark mode"
    @click="toggleDark()"
  />
</template>

<style scoped>
button {
  display: inline-block;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  cursor: pointer;
  margin: 0;
  border:  2px solid var(--accent);
  background: center url("@/assets/icons/sun.svg") white no-repeat;
  background-size: 20px 20px;
}

.dark button {
  background: center url("@/assets/icons/moon.svg") black no-repeat;
  background-size: 20px 20px;
}
</style>
