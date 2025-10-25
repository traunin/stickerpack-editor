<template>
  <div class="user-packs-scroll">
    <UserPacksList class="user-packs" />
  </div>
</template>

<script setup lang = "ts">
import UserPacksList from '@/components/user-packs-list.vue'
</script>

<style scoped>
.user-packs-scroll {
  flex: 1;
  align-self: stretch;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  overflow: auto;
  scrollbar-color: var(--accent) var(--input);
  scrollbar-width: thin;
}

.user-packs {
  flex: 1;
  padding: 20px;
  max-width: 1150px;
}
</style>
