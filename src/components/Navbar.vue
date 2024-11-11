<template>    
    <nav>
        <h1>Stickerpack editor</h1>
        <ul>
            <li>
                <router-link 
                    :to="`/`"
                    class="link"
                    active-class="active"
                ><span>My stickerpacks</span></router-link>
            </li>
        </ul>

        <div class="tg-auth">
            <telegram-login :onAuth="handleAuth"></telegram-login>
        </div>
    </nav>
</template>

<script setup lang="ts">
import TelegramLogin from '@/components/TelegramLogin.vue';

function handleAuth(user: any) {
    alert(user)
}
</script>

<style scoped>
h1 {
    font-size: 1.5em;
    text-align: center;
    padding: 20px 40px;
}

nav {
    display: flex;
    background: var(--secondary);
    padding: 1px;
}

ul {
    display: flex;
    flex: 1;
    justify-content: center;
    margin-left: -280px;
    margin-right: -218px;
}

li {
    display: flex;
}

.link {
    padding: 20px;
    text-decoration: none;
    color: var(--text);
    font-size: 1.2em;
    border-left: 3px solid var(--secondary);
    border-right: 3px solid var(--secondary);
    display: flex;
    justify-content: center;
    align-items: center;
}

.link:hover {
    background: var(--accent);
    border-color: var(--accent)
}

.active {
    border-color: var(--text)
}

.active:hover {
    border-color: var(--text);
}

.tg-auth {
    margin: 0 40px;
    display: flex;
    align-items: center;
}
</style>