<template>
  <div class="error">
    <div class="icon">
      !
    </div>
    <div class="message">
      {{ message }}
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  message: string
}>()
</script>

<style scoped>
.error {
  font-size: 1.5em;
  background: var(--background);
  border: 3px solid var(--primary);
  color: var(--text);
  display: flex;
}

.icon {
  background: red;
  aspect-ratio: 1;
  padding: 10px;
  font-weight: 900;
  display: flex;
  align-items: center;
  justify-content: center;
}

.message {
  padding: 5px;
  display: flex;
  align-items: center;
}
</style>
