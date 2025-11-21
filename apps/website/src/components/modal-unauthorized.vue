<template>
  <div class="modal-unauthorized">
    <MenuIcon class="menu-icon" @click="open = !open" />
    <Transition name="fade">
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
        </ul>
      </nav>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import MenuIcon from '@/assets/icons/burger.svg'
import DropdownThemeButton from '@/components/dropdown-theme-button.vue'
import NavbarElement from '@/components/navbar-element.vue'
import { navbarRoutes } from '@/router'
import DropdownGithub from './dropdown-github.vue'

const open = ref(false)
</script>

<style scoped>
.menu-icon {
  width: 42px;
  height: 42px;
  color: var(--text);
  padding: 5px;
  cursor: pointer;
  border-radius: 10px;
}

.menu-icon:hover {
  background-color: var(--panel);
}

.modal-unauthorized {
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
  z-index: 10;
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

.github-link {
  border-radius: 10px;
}

li > *:first-child {
  width: 100%
}

li > *:first-child:hover {
  background: var(--panel);
}
</style>
