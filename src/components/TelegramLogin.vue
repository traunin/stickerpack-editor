<template>
    <div 
        ref="telegramWidget" 
        class="tg-auth-wrapper"
    ></div>
</template>

<script setup lang="ts">
declare global {
    interface Window { onAuth: any }
}

import { onMounted, ref } from 'vue';

const props = defineProps({
    onAuth: {
        type: Function,
        required: true
    }
})

const telegramWidget = ref<null | HTMLElement>(null);

onMounted(() => {
    window.onAuth = props.onAuth;

    const script = document.createElement('script');
    script.setAttribute('data-telegram-login', 'seventv_stickerpack_bot');
    script.setAttribute('data-size', 'large');
    script.setAttribute('data-userpic', 'false');
    script.setAttribute('data-radius', '10');
    script.setAttribute('data-onauth', 'window.onAuth(user)');
    script.setAttribute('data-request-access', 'write');
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.async = true;


    telegramWidget.value!.appendChild(script);
});
</script>

<style scoped>
.tg-auth-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>