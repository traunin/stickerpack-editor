<template>
    <div class="stickerpack-preview">
        {{ props.name }}
    </div>
</template>

<script setup lang = "ts">
const props = defineProps<{
    name: string
}>()
</script>

<style scoped>
.stickerpack-preview {
    width: 200px;
    height: 200px;
    border: 3px solid red;
    text-align: center;
    padding: 10px;
}
</style>