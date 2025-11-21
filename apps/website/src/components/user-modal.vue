<template>
  <div class="user-modal">
    <CreateLink class="create" @click="open = false" />
    <div class="user-info" @click="open = !open">
      <img src="@/assets/hi.gif" alt="">
      <img :src="authStore.photoURL" alt="" class="profile-picture">
    </div>
    <Transition name="fade">
      <Teleport to="body">
        <nav v-if="open">
          <ul>
            <li
              v-for="(route, i) in navbarRoutes"
              :key="route.path"
              :class="[{ 'route-last': i === navbarRoutes.length - 1 }, { 'route-first': i === 0 }]"
              @click="open = false"
            >
              <NavbarElement
                :route="route"
                orientation="vertical"
              />
            </li>
            <li class="theme-switch">
              <DropdownThemeButton />
            </li>
            <li class="github-link">
              <DropdownGithub />
            </li>
            <li class="sign-out" @click="open = false">
              <DropdownLogoutButton />
            </li>
          </ul>
        </nav>
      </Teleport>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import CreateLink from '@/components/create-link.vue'
import DropdownGithub from '@/components/dropdown-github.vue'
import DropdownLogoutButton from '@/components/dropdown-logout-button.vue'
import DropdownThemeButton from '@/components/dropdown-theme-button.vue'
import NavbarElement from '@/components/navbar-element.vue'
import { navbarRoutes } from '@/router'
import { useTgAuthStore } from '@/stores/use-tg-auth'

const open = ref(false)
const authStore = useTgAuthStore()
</script>

<style scoped>
.profile-picture {
  width: 36px;
  height: 36px;
  border: 2px solid var(--primary);
  border-radius: 100%;
}

.user-info {
  font-size: 1.2em;
  display: flex;
  align-items: center;
  height: 40px;
  border-radius: 10px;
  padding: 22px 15px;
  cursor: pointer;
  user-select: none;
}

.user-info:hover {
  background: var(--panel);
}

.user-modal {
  display: flex;
  gap: 10px;
}

nav {
  position: absolute;
  inset: 42px 0 0 0;
  background: var(--input);
  margin-top: 8px;
  border: var(--border) 2px solid;
  align-items: stretch;
  z-index: 60;
  isolation: unset;
}

ul {
  margin: 10px;
}

li {
  background: var(--panel);
  overflow: hidden;
}

.route-last {
  margin-bottom: 10px;
  border-radius: 0 0 10px 10px;
}

.route-first {
  border-radius: 10px 10px 0 0;
}

.theme-switch {
  margin-bottom: 10px;
  border-radius: 10px;
}

.sign-out {
  border-radius: 10px;
}

.github-link {
  margin-bottom: 10px;
  border-radius: 10px;
}

li > *:first-child {
  width: 100%
}

li > *:first-child:hover {
  background: var(--panel);
}
</style>
