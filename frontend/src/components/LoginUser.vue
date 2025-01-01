<script setup lang="ts">
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';

</script>

<template>
  <Form class="flex flex-col gap-4">
    <div>
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-user"></i>
        </InputGroupAddon>
        <InputText name="username" type="text" placeholder="Username" class="rounded-r-md" />
      </InputGroup>
    </div>

    <div>
      <InputGroup>
        <InputGroupAddon><i class="pi pi-lock"></i></InputGroupAddon>
        <InputText name="password" type="password" placeholder="Password" />
      </InputGroup>
    </div>
  </Form>
</template>
