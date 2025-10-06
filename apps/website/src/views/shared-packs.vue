<template>
  <div class="public-packs-row">
    <PublicPacksList class="public-packs" />
  </div>
</template>

<script setup lang = "ts">
import PublicPacksList from '@/components/public-packs-list.vue'
</script>

<style scoped>
.public-packs-row {
  flex: 1;
  align-self: stretch;
  display: flex;
  align-items: center;
  justify-content: center;
}

.public-packs {
  flex: 1;
  padding: 20px;
  background: var(--panel)
}
</style>
