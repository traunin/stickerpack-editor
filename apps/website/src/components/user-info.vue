<template>
  <div v-if="authStore.isLoggedIn">
    <UserDropdown v-if="isDesktop" />
    <UserModal v-else />
  </div>
  <div v-else>
    <DropdownUnauthorized v-if="isDesktop" />
    <ModalUnauthorized v-else />
    <SignInButton />
  </div>
</template>

<script setup lang="ts">
import DropdownUnauthorized from '@/components/dropdown-unauthorized.vue'
import ModalUnauthorized from '@/components/modal-unauthorized.vue'
import SignInButton from '@/components/sign-in-button.vue'
import UserDropdown from '@/components/user-dropdown.vue'
import UserModal from '@/components/user-modal.vue'
import { useIsDesktop } from '@/composables/use-is-desktop'
import { useTgAuthStore } from '@/stores/use-tg-auth'

const isDesktop = useIsDesktop()
const authStore = useTgAuthStore()
</script>

<style scoped>
div {
  display: flex;
  align-items: center;
  gap: 10px;
}
</style>
