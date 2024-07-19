<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  location.pathname = '/auth/login'
})
</script>

<template>
  <div class="flex flex-col gap-2 items-center my-16">
    <h2 class="font-bold text-2xl">
      Redirecting...
    </h2>
    <p class="text-sm text-gray-500 dark:text-neutral-400">
      If nothing happens, <a class="underline underline-offset-2" href="/auth/login">click here</a>.
    </p>
  </div>
</template>
