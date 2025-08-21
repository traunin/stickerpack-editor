<template>
  <div class="overview">
    <div class="user-packs-panel">
      <h1>User packs</h1>
      <div class="user-packs-row">
        <AddStickerpack />
        <UserPacksList class="user-packs" />
      </div>
    </div>
    <div class="public-packs-panel">
      <h1>Shared packs</h1>
      <PublicPacksList class="public-packs" />
    </div>
  </div>
</template>

<script setup lang = "ts">
import AddStickerpack from '@/components/create-link.vue'
import PublicPacksList from '@/components/public-packs-list.vue'
import UserPacksList from '@/components/user-packs-list.vue'
</script>

<style scoped>
h1 {
  color: var(--text);
  font-size: 1.5em;
  border-bottom: 2px solid var(--text);
  padding: 5px;
  margin: 0 0 20px 0;
}

.overview {
  flex: 1;
  align-self: stretch;
  display: flex;
  flex-direction: column;
  min-height: 0;
  gap: 20px;
  justify-content: space-around;
}

.user-packs-row {
  display: flex;
  gap: 40px;
  align-self: stretch;
}

.public-packs {
  display: flex;
  flex-wrap: wrap;
  align-self: stretch;
  display: flex;
}

.user-packs {
  flex: 1;
}

.user-packs-panel, .public-packs-panel {
  padding: 20px;
  background: var(--panel)
}
</style>
