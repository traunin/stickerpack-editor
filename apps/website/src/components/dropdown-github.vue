<template>
  <a href="https://github.com/Traunin/stickerpack-editor" class="github">
    <GithubIcon class="github-icon" />
    <span>
      Report an issue
    </span>
  </a>
</template>

<script setup lang="ts">
import GithubIcon from '@/assets/icons/github.svg'
</script>

<style scoped>
.github {
  padding: 5px;
  cursor: pointer;
  flex: 1;
  display: flex;
  background: none;
  color: var(--text);
  border: none;
  font-size: 1.1em;
  align-items: center;
  text-decoration: none;
}

.github-icon {
  display: inline-block;
  margin: 5px 13px 5px 5px;
  border-radius: 50%;
  cursor: pointer;
}
</style>
