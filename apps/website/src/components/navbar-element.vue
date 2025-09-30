<template>
    <RouterLink
      :to="route.path"
      class="link"
      active-class="active"
    >
      <span>
        {{ route.name }}
      </span>
    </RouterLink>
</template>

<script setup lang="ts">
import { RouteRecordRaw } from 'vue-router';

defineProps<{
  route: RouteRecordRaw
}>()
</script>

<style scoped>
.link {
  padding: 20px;
  text-decoration: none;
  color: var(--text);
  font-size: 1.2em;
  border-left: 3px solid var(--secondary);
  border-right: 3px solid var(--secondary);
  display: flex;
  justify-content: center;
  align-items: center;
}

.link:hover {
  background: var(--accent);
  border-color: var(--accent)
}

.active {
  border-color: var(--text)
}

.active:hover {
  border-color: var(--text);
}
</style>