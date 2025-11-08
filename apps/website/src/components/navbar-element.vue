<template>
  <RouterLink
    :to="route.path"
    class="link" :class="[orientation]"
    active-class="active"
    draggable="false"
  >
    <component :is="route.meta?.icon" v-if="showIcon" class="logo" />
    {{ route.name }}
  </RouterLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { RouteRecordRaw } from 'vue-router'

const props = defineProps<{
  route: RouteRecordRaw
  orientation?: 'horizontal' | 'vertical'
}>()

const showIcon = computed(
  () => props.route.meta?.icon && props.orientation === 'vertical',
)
</script>

<style scoped>
.link {
  padding: 18px;
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

.vertical {
  border: none;
  justify-content: flex-start;
  padding: 12px;
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

.logo {
  width: 20px;
  height: 20px;
  color: var(--text);
  margin-right: 16px;
}
</style>
