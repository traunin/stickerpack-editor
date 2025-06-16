<template>
  <RouterLink
    to="create"
    class="link"
  >
    <span class="logo">
      +
    </span>
    <span class="text">
      Create a new pack
    </span>
  </RouterLink>
</template>

<script setup lang="ts">

</script>

<style scoped>
.link {
  display: flex;
  flex-direction: column;
  border: 3px solid var(--primary);
  justify-content: center;
  align-items: center;
  color: var(--text);
  text-decoration: none !important;
  padding: 20px;
  font-size: 1.3em;
}

.logo {
  font-size: 10em;
}
</style>
