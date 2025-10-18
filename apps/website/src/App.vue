<template>
  <NavbarHeader />
  <main>
    <RouterView />
  </main>
</template>

<script setup lang="ts">
import { useDark } from '@vueuse/core'
import { RouterView } from 'vue-router'
import NavbarHeader from '@/components/navbar-header.vue'

useDark({
  selector: 'body',
  attribute: 'class',
  valueDark: 'dark',
  valueLight: 'light',
})
</script>

<style scoped>
main {
  display: flex;
  flex-grow: 1;
  flex-direction: column;
  justify-content: flex-start;
  align-items: flex-start;
  min-height: 0;
}
</style>
