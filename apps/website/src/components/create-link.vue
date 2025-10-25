<template>
  <RouterLink
    to="create"
    class="link"
    draggable="false"
  >
    <AddIcon class="logo" />
    <span class="text">
      Create
    </span>
  </RouterLink>
</template>

<script setup lang="ts">
import AddIcon from '@/assets/icons/add.svg'
</script>

<style scoped>
.link {
  display: flex;
  justify-content: center;
  align-items: center;
  color: var(--text);
  text-decoration: none !important;
  font-size: 1.2em;
  background: var(--panel);
  border-radius: 10px;
  padding: 10px;
  gap: 10px;
  align-self: center;
  user-select: none;
}

.logo {
  width: 20px;
  height: 20px;
  color: var(--text);
}

.text {
  display: none;
}

@media screen and (min-width: 1000px) {
  .text {
    display: block;
  }
}
</style>
