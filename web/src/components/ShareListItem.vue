<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useStore } from '../lib/store.ts'
import { formatSize } from '../utils/filesize.ts'
import ShareIcon from './ShareIcon.vue'
import type { Share } from '../types/share.ts'

import BiDot from 'bootstrap-icons/icons/dot.svg?component'
import BiLockFill from 'bootstrap-icons/icons/lock-fill.svg?component'

const store = useStore()

const router = useRouter()

const props = defineProps<{
  share: Share
}>()

const isOwner = computed(() => {
  return store.loggedIn && store.userinfo?.subject === props.share.creator?.subject
})

const onShow = () => {
  router.push(`/shares/${props.share.name}`)
}
</script>

<template>
  <div
    class="flex flex-col items-stretch gap-2 bg-white dark:bg-neutral-800 p-4 rounded-lg shadow select-none cursor-pointer hover:shadow-md"
    @click="onShow"
  >
    <h5
      class="inline-flex items-center gap-2 text-base font-name font-medium"
      :class="isOwner ? '' : 'text-gray-500 dark:text-gray-400'"
    >
      <share-icon :share-type="share.type" class="w-5 h-5" />
      <span v-text="share.displayName ?? share.name" />
    </h5>
    <p class="flex items-center flex-wrap text-xs text-gray-500 dark:text-neutral-400">
      <template v-if="share.files && share.files.length">
        <span v-if="share.files.length > 1">{{ share.files.length }} files</span>
        <span
          v-else
          v-text="share.files[0].path"
        />
        <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
      </template>
      <template v-if="share.size">
        <span v-text="formatSize(share.size)" />
        <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
      </template>
      <template v-if="share.password">
        <bi-lock-fill class="w-3 h-3" />
        <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
      </template>
      <template v-if="share.expiry">
        <span>{{ share.expiry }} days</span>
        <bi-dot class="w-3 h-3 mx-1 text-gray-400 dark:text-neutral-500 last:hidden" />
      </template>
    </p>
  </div>
</template>
