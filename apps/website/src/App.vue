<template>
  <div class="container dark">
    <NavbarHeader />
    <div class="content">
      <RouterView />
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import NavbarHeader from '@/components/navbar-header.vue'
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
}
</style>
