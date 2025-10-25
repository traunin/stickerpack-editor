<template>
  <div class="shared-packs-scroll">
    <PublicPacksList class="shared-packs" />
  </div>
</template>

<script setup lang = "ts">
import PublicPacksList from '@/components/public-packs-list.vue'
</script>

<style scoped>
.shared-packs-scroll {
  flex: 1;
  align-self: stretch;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  overflow: auto;
  scrollbar-color: var(--accent) var(--input);
  scrollbar-width: thin;
}

.shared-packs {
  flex: 1;
  padding: 20px;
  max-width: 1150px;
}
</style>
