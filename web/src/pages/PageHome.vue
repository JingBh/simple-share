<script setup lang="ts">
import { useStore } from '../lib/store.ts'

import BiBoxArrowInLeft from 'bootstrap-icons/icons/box-arrow-in-left.svg?component'
import BiCloudPlus from 'bootstrap-icons/icons/cloud-plus.svg?component'
import BiShare from 'bootstrap-icons/icons/share.svg?component'

const store = useStore()
</script>

<template>
  <div class="w-full min-h-screen flex flex-col justify-center items-center gap-16 md:gap-24 p-8">
    <h1 class="font-title font-semibold text-4xl sm:text-5xl text-center">
      Simple Share
    </h1>
    <div
      class="flex flex-col items-stretch gap-6 transition-opacity"
      :class="{ 'opacity-0 pointer-events-none': store.userinfoLoading }"
    >
      <router-link
        v-if="!store.loggedIn && !store.authDisabled"
        class="inline-flex items-center gap-4 text-gray-200 bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-400 font-medium rounded-lg px-5 py-3 dark:text-gray-800 dark:bg-gray-100 dark:hover:bg-gray-200"
        to="/login"
      >
        <bi-box-arrow-in-left class="w-4 h-4" />
        <span class="inline-block flex-1 text-center">
          Login with OIDC
        </span>
      </router-link>
      <router-link
        v-if="store.loggedIn"
        class="inline-flex items-center gap-4 text-gray-200 bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-400 font-medium rounded-lg px-5 py-3 dark:text-gray-800 dark:bg-gray-100 dark:hover:bg-gray-200"
        to="/shares"
      >
        <bi-share class="w-4 h-4" />
        <span class="inline-block flex-1 text-center">
          All Shares
        </span>
      </router-link>
      <router-link
        v-if="store.loggedIn"
        class="inline-flex items-center gap-4 text-gray-200 bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-400 font-medium rounded-lg px-5 py-3 dark:text-gray-800 dark:bg-gray-100 dark:hover:bg-gray-200"
        to="/create"
      >
        <bi-cloud-plus class="w-4 h-4" />
        <span class="inline-block flex-1 text-center">
          Create Share
        </span>
      </router-link>
      <div
        v-if="!store.loggedIn"
        class="flex flex-col items-center gap-1 text-xs text-gray-400 dark:text-neutral-500"
      >
        <p>This is a personal file sharing service.</p>
        <p>To deploy your own, see the source <a class="underline underline-offset-2" href="https://github.com/JingBh/simple-share" target="_blank" rel="noreferrer">here</a>.</p>
      </div>
    </div>
  </div>
</template>
