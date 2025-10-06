<template>
  <header>
    <h1>Stickerpack<br>editor</h1>
    <nav>
      <NavbarElement
        v-for="route in navbarRoutes"
        :key="route.path"
        :route="route"
      />
    </nav>
    <UserInfoDesktop />
  </header>
</template>

<script setup lang="ts">
import { navbarRoutes } from '@/router'
import NavbarElement from '@/components/navbar-element.vue';
import UserInfoDesktop from '@/components/user-info-desktop.vue';
</script>

<style scoped>
h1 {
  font-size: 1.2em;
  border-left: 3px solid var(--primary);
  border-right: 3px solid var(--primary);
  align-self: stretch;
  display: flex;
  text-align: center;
  align-items: center;
  padding: 0 10px;
}

header {
  display: flex;
  border-bottom: 1px solid var(--border);
  padding: 0 20px;
  gap: 20px;
  align-items: center;
}

nav {
  flex: 1;
  display: flex;
}
</style>
