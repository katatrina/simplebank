<script setup lang="ts">
import Card from 'primevue/card';
import Divider from 'primevue/divider';
import Button from 'primevue/button';

import type { User } from '@/types/user';
import type { PropType } from 'vue';

defineProps({
  user: {
    type: Object as PropType<User>,
    required: true,
  },
})

</script>

<template>
  <Card>
    <template #title><span class="font-medium">User Profile</span></template>
    <template #content>
      <div class="flex flex-col gap-2">
        <div>
          <i class="pi pi-user mr-2"></i>
          <span class="">{{ user.full_name }}</span>
        </div>
        <div>
          <i class="pi pi-envelope mr-2"></i>
          <a :href="`mailto:${user.email}`">{{ user.email }}</a>
        </div>
      </div>

      <Divider type="solid" class="mt-4 mb-2" />
    </template>

    <template #footer>
      <Button label="Logout" icon="pi pi-sign-out" @click="$emit('logout', user)" />
    </template>
  </Card>
</template>
