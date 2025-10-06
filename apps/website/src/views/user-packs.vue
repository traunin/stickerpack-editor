<template>
  <div class="user-packs-row">
    <UserPacksList class="user-packs" />
  </div>
</template>

<script setup lang = "ts">
import UserPacksList from '@/components/user-packs-list.vue'
</script>

<style scoped>
.user-packs-row {
  flex: 1;
  align-self: stretch;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-packs {
  flex: 1;
  padding: 20px;
  background: var(--panel)
}
</style>
