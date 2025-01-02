<script setup lang="ts">
import LoginForm from '@/components/LoginForm.vue'
import UserProfileCard from '@/components/UserProfileCard.vue'
import store from '@/types/store'
import type { User } from '@/types/user'
import { useToast } from 'primevue/usetoast'

const toast = useToast()

function onLogout(user: User) {
  store.clearUser()
  toast.add({
    severity: 'info',
    summary: `Goodbye, ${user.full_name}!`,
    detail: 'You have successfully logged out.',
    life: 3000
  })
}
</script>

<template>
  <main>
    <h1 class="green font-bold mb-5">Welcome to Simple Bank!</h1>
    <UserProfileCard v-if="store.authState.user" :user="store.authState.user" @logout="onLogout" />
    <LoginForm v-else />
  </main>
</template>
