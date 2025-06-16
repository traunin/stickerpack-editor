<template>
  <div class="container dark">
    <NavbarHeader />
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import NavbarHeader from '@/components/navbar-header.vue'
</script>

<style scoped>
.container {
    position: fixed;
    background: var(--background);
    color: var(--text);
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
}
</style>
