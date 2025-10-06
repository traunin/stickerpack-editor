<template>
  <div class="signed-in" v-if="authStore.isLoggedIn" >
    <CreateLink class="create" />
    <UserDropdownDesktop />
  </div>
  <SignInButton v-else />
</template>

<script setup lang="ts">
import { useTgAuthStore } from '@/stores/use-tg-auth';
import UserDropdownDesktop from '@/components/user-dropdown.vue';
import SignInButton from '@/components/sign-in-button.vue';
import CreateLink from '@/components/create-link.vue';


const authStore = useTgAuthStore()
</script>

<style scoped>
.signed-in {
  display: flex;
  gap: 20px;
}
</style>