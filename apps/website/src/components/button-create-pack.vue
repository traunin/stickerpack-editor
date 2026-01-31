<template>
  <button :disabled="!!error" class="create-button">
    <slot>Create</slot> {{ error ? ` | ${error}` : '' }}
  </button>
</template>

<script setup lang="ts">
defineProps<{
  error: null | string
}>()
</script>

<style scoped>
.create-button:disabled {
  background: var(--disabled);
  cursor: not-allowed;
}

.create-button {
  background: var(--input);
  color: var(--text);
  font-size: 1.3em;
  padding: 10px;
  border: none;
  margin-top: 0;
  cursor: pointer;
  border-radius: 10px;
}
</style>
