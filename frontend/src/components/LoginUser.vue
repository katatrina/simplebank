<script setup lang="ts">
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import FloatLabel from 'primevue/floatlabel';
import Button from 'primevue/button';

import { ref, computed } from 'vue';
import axios from 'axios';

const username = ref<string>('');
const password = ref<string>('');

const isLoginDisabled = computed(() => !username.value || !password.value);

async function handleLogin() {
  try {
    const response = await axios.post('http://localhost:8080/v1/users/login', {
      username: username.value,
      password: password.value,
    });

    console.log(response);
  } catch (error) {
    console.error(error);
  } finally {
    username.value = '';
    password.value = '';
  }
}

</script>

<template>
  <form @submit.prevent="handleLogin" class="flex flex-col gap-4">
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
      <Button label="Login" class="w-full" :disabled="isLoginDisabled" />
    </div>
  </form>
</template>
