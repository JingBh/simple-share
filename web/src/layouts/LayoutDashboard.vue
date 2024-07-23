<script setup lang="ts">
import AppHeader from '../components/AppHeader.vue'
</script>

<template>
  <div class="mx-auto max-w-screen-lg pt-20 sm:pt-24 pb-8 px-4 sm:px-8">
    <app-header />
    <slot />
  </div>
</template>
