<template>
  <RouterLink
    :to="route.path"
    class="link"
    active-class="active"
    draggable="false"
  >
    {{ route.name }}
  </RouterLink>
</template>

<script setup lang="ts">
import type { RouteRecordRaw } from 'vue-router'

defineProps<{
  route: RouteRecordRaw
}>()
</script>

<style scoped>
.link {
  padding: 15px;
  text-decoration: none;
  color: var(--text);
  font-size: 1.1em;
  display: flex;
  justify-content: center;
  align-items: center;
  border-top: 2px solid transparent;
  border-bottom: 2px solid transparent;
  user-select: none;
}

.link:hover {
  border-bottom-color: var(--border-hover);
}

.active {
  border-bottom-color: var(--text);
}

.active:hover {
  border-bottom-color: var(--text);
}
</style>
