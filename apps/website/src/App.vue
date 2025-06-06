<template>
    <div class="container dark">
        <navbar></navbar>
        <router-view></router-view>
    </div>
</template>

<script setup lang="ts">
import Navbar from '@/components/Navbar.vue';
</script>

<style scoped>
.container {
    position: fixed;
    background: var(--background);
    color: var(--text);
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
}
</style>