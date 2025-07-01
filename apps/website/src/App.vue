<template>
  <div class="container">
    <NavbarHeader />
    <div class="content">
      <RouterView />
    </div>
    <ThemeSwitch class="theme-switcher" />
  </div>
</template>

<script setup lang="ts">
import { useDark } from '@vueuse/core'
import { RouterView } from 'vue-router'
import NavbarHeader from '@/components/navbar-header.vue'
import ThemeSwitch from './components/theme-switch.vue'

useDark({
  selector: 'body',
  attribute: 'class',
  valueDark: 'dark',
  valueLight: 'light',
})
</script>

<style scoped>
.container {
  position: fixed;
  display: flex;
  flex-direction: column;
  background: var(--background);
  color: var(--text);
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.content {
  display: flex;
  flex-grow: 1;
  flex-direction: column;
  justify-content: flex-start;
  align-items: flex-start;
  min-height: 0;
}

.theme-switcher {
  position: fixed;
  right: 20px;
  bottom: 20px;
}
</style>
