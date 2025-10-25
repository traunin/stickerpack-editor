<template>
  <header>
    <h1>Stickerpack<br>editor</h1>
    <nav v-if="isDesktop">
      <NavbarElement
        v-for="route in navbarRoutes"
        :key="route.path"
        :route="route"
      />
    </nav>
    <UserInfo />
  </header>
</template>

<script setup lang="ts">
import NavbarElement from '@/components/navbar-element.vue'
import UserInfo from '@/components/user-info.vue'
import { useIsDesktop } from '@/composables/use-is-desktop'
import { navbarRoutes } from '@/router'

const isDesktop = useIsDesktop()
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
  padding: 6px 10px;
}

header {
  display: flex;
  border-bottom: 1px solid var(--border);
  padding: 0 10px 0 0;
  gap: 20px;
  align-items: center;
  justify-content: space-between;
}

nav {
  flex: 1;
  display: flex;
}

@media screen and (min-width: 1000px) {
  header {
    padding: 0 20px;
  }
}
</style>
