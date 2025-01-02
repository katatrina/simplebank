<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import FloatLabel from 'primevue/floatlabel';
import Button from 'primevue/button';
import { useToast } from 'primevue/usetoast';

import { ref, computed } from 'vue';
import axios from 'axios';
import type { User } from '@/types/user';
import store from '@/types/store';

interface LoginResponse {
  user: User
  access_token: string
  refresh_token: string
}

const toast = useToast();

const username = ref<string>('');
const password = ref<string>('');

const isLoginDisabled = computed(() => !username.value || !password.value);

const errorMessage = ref<string>('');

async function handleLogin() {
  try {
    const response = await axios.post<LoginResponse>('http://localhost:8080/v1/users/login', {
      username: username.value,
      password: password.value,
    });

    store.setUser(response.data.user, response.data.access_token, response.data.refresh_token)
    toast.add({
      severity: 'success',
      summary: `Hello, ${response.data.user.full_name}!`,
      detail: 'You have successfully logged in.',
      life: 3000
    });
  } catch (error: any) {
    if (error.response && (error.response.status === 404 || error.response.status === 401)) {
      // errorMessage.value = error.response.data.message;
      errorMessage.value = 'Invalid username or password. Please try again.';
    } else {
      errorMessage.value = 'An error occurred. Please try again later.';
    }

    toast.add({
      severity: 'error',
      summary: 'Login Failed',
      detail: errorMessage.value,
      life: 3000
    });
  } finally {
    username.value = '';
    password.value = '';
  }
}

</script>

<template>
  <form @submit.prevent="handleLogin" class="flex flex-col gap-4">
    <!--Username input -->
    <div>
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-user"></i>
        </InputGroupAddon>
        <FloatLabel variant="on">
          <InputText id="username" v-model="username" name="username" type="text" class="rounded-r-md"
            autocomplete="on" />
          <label for="username">Username</label>
        </FloatLabel>
      </InputGroup>
    </div>

    <!-- password input -->
    <div>
      <InputGroup>
        <InputGroupAddon><i class="pi pi-lock"></i></InputGroupAddon>
        <FloatLabel variant="on">
          <InputText id="password" type="password" v-model="password" autocomplete="off" />
          <label for="password">Password</label>
        </FloatLabel>
      </InputGroup>
    </div>

    <div>
      <Button type="submit" label="Login" class="w-full" :disabled="isLoginDisabled" />
    </div>
  </form>
</template>
