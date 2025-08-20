<template>
  <div class="overview">
    <div class="user-packs-row">
      <AddStickerpack />
      <UserPacksList endpoint="user/packs" class="user-packs" />
    </div>
    <PublicPacksList endpoint="public/packs" class="public-packs" />
  </div>
</template>

<script setup lang = "ts">
import AddStickerpack from '@/components/create-link.vue'
import PublicPacksList from '@/components/public-packs-list.vue'
import UserPacksList from '@/components/user-packs-list.vue'
</script>

<style scoped>
.overview {
  flex: 1;
  align-self: stretch;
  display: flex;
  flex-direction: column;
  min-height: 0;
  margin: 20px;
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
</style>
