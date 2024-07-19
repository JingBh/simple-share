<script setup lang="ts">
import { useStore } from '../lib/store.ts'
import { computed } from 'vue'

import BiBoxArrowInLeft from 'bootstrap-icons/icons/box-arrow-in-left.svg?component'
import BiPersonCircle from 'bootstrap-icons/icons/person-circle.svg?component'

const store = useStore()

const username = computed<string | null>(() => {
  return store.userinfo?.username || store.userinfo?.subject.substring(0, 8) || null
})
</script>

<template>
  <header class="fixed top-0 left-0 w-full h-16 px-4 sm:px-8 flex items-center justify-between gap-4 bg-white dark:bg-neutral-800">
    <router-link
      class="font-bold text-2xl"
      to="/"
    >
      Simple Share
    </router-link>
    <div
      v-if="store.loggedIn"
      class="inline-flex items-center gap-2 select-none"
    >
      <bi-person-circle class="w-6 h-6 text-gray-400 dark:text-neutral-500" />
      <span
        class="text-gray-600 dark:text-neutral-300 font-semibold text-sm"
        v-text="username"
      />
    </div>
    <router-link
      v-else-if="!store.authDisabled"
      class="inline-flex items-center gap-1 text-gray-200 bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-400 text-sm font-medium rounded-lg px-3 py-2 dark:text-gray-800 dark:bg-gray-100 dark:hover:bg-gray-200"
      to="/login"
    >
      <bi-box-arrow-in-left class="w-4 h-4" />
      <span class="inline-block flex-1 text-center">
        Login
      </span>
    </router-link>
  </header>
</template>
