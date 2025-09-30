<template>
  <header>
    <h1>Stickerpack editor</h1>
    <nav>
      <NavbarElement
        v-for="route in navbarRoutes"
        :key="route.path"
        :route="route"
      />
    </nav>
    <TelegramLogin />
  </header>
</template>

<script setup lang="ts">
import TelegramLogin from '@/components/telegram-login.vue'
import { navbarRoutes } from '@/router'
import NavbarElement from '@/components/navbar-element.vue';
</script>

<style scoped>
h1 {
  font-size: 1.5em;
  text-align: center;
}

header {
  display: flex;
  background: var(--secondary);
  padding: 0 20px;
  gap: 20px;
  align-items: center;
}

nav {
  flex: 1;
  display: flex;
}
</style>
