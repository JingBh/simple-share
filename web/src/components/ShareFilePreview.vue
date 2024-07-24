<script setup lang="ts">
import { computed } from 'vue'

import type { ShareFile } from '../types/share.ts'

import BiCloudDownload from 'bootstrap-icons/icons/cloud-download.svg?component'

type PreviewType = 'image' | 'audio' | 'video' | 'document' | 'text' | null

defineProps<{
  file: ShareFile
}>()

const emit = defineEmits([
  'download'
])

const previewType = computed<PreviewType>(() => {
  // const parts = props.file.path.split('.')
  // const ext = parts[parts.length - 1].toLowerCase()

  // if (['gif', 'jpg', 'jpeg', 'png', 'svg', 'webp'].includes(ext)) {
  //   return 'image'
  // }

  return null
})
</script>

<template>
  <div v-if="previewType === 'image'">
    <!-- TODO -->
  </div>
  <div v-else-if="previewType === 'audio'">
    <!-- TODO -->
  </div>
  <div v-else-if="previewType === 'video'">
    <!-- TODO -->
  </div>
  <div v-else-if="previewType === 'document'">
    <!-- TODO -->
  </div>
  <div v-else-if="previewType === 'text'">
    <!-- TODO -->
  </div>
  <div
    v-else
    class="flex flex-col items-center gap-3 py-8"
  >
    <p class="text-lg text-gray-500 dark:text-neutral-400">
      Preview not available
    </p>
    <button
      class="inline-flex items-center gap-2 text-gray-200 bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-400 font-medium rounded-lg px-4 py-2 dark:text-gray-800 dark:bg-gray-100 dark:hover:bg-gray-200"
      @click="emit('download')"
    >
      <bi-cloud-download class="flex-shrink-0 w-3 h-3" />
      <span>Download File</span>
    </button>
  </div>
</template>
