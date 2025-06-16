<template>
  <nav>
    <h1>Stickerpack editor</h1>
    <ul>
      <li>
        <router-link
          to="`/`"
          class="link"
          active-class="active"
        >
          <span>My stickerpacks</span>
        </router-link>
      </li>
    </ul>

    <div class="tg-auth">
      <TelegramLogin :on-auth="handleAuth" />
    </div>
  </nav>
</template>

<script setup lang="ts">
import type { User } from '@/components/telegram-login.vue'
import TelegramLogin from '@/components/telegram-login.vue'

function handleAuth(user: User) {
  console.log(user)
  // alert('Logged in as ' + user.first_name + ' ' + user.last_name + ' (' + user.id + (user.username ? ', @' + user.username : '') + ')');
}
</script>

<style scoped>
h1 {
  font-size: 1.5em;
  text-align: center;
  padding: 20px;
  position: absolute;
  left: 20px;
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
  padding: 10px 20px;
  position: absolute;
  right: 20px;
  display: flex;
  align-items: center;
}
</style>
