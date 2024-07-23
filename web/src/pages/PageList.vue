<script setup lang="ts">
import { ref } from 'vue'
import { useInfiniteScroll } from '@vueuse/core'

import { useAxiosInstance } from '../lib/axios.ts'
import FlowbiteSpinner from '../components/FlowbiteSpinner.vue'
import LayoutDashboard from '../layouts/LayoutDashboard.vue'
import ShareListItem from '../components/ShareListItem.vue'
import type { Share } from '../types/share.ts'

const data = ref<Share[]>([])

const nextCursor = ref<string | null>(null)

const { isLoading } = useInfiniteScroll(document.body, async () => {
  const res = (await useAxiosInstance().get<{
    data: Share[]
    cursor: string
  }>('/api/shares', {
    params: {
      cursor: nextCursor.value
    }
  })).data
  nextCursor.value = res.cursor
  data.value.push(...res.data)
}, {
  distance: 10,
  canLoadMore: () => {
    return nextCursor.value === null || nextCursor.value !== ''
  }
})
</script>

<template>
  <layout-dashboard>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-8">
      <share-list-item
        v-for="share of data"
        :key="share.name"
        :share="share"
      />
    </div>
    <p
      v-if="isLoading"
      class="flex items-center justify-center gap-2 py-8 text-sm text-gray-500 dark:text-neutral-400"
    >
      <flowbite-spinner class="w-5 h-5" />
      <span>Loading more data...</span>
    </p>
  </layout-dashboard>
</template>
