<script setup lang="ts">
import { ref, watch, watchEffect } from 'vue'

import { useAxiosInstance } from '../lib/axios.ts'
import { getContentUrl, getContentTypeUrl, getPreviewUrl } from '../utils/share-url.ts'
import FlowbiteSpinner from './FlowbiteSpinner.vue'
import type { Share, ShareFile } from '../types/share.ts'

import BiCloudDownload from 'bootstrap-icons/icons/cloud-download.svg?component'

type FileType = 'image' | 'audio' | 'video' | 'document' | 'text' | 'unknown'

const props = defineProps<{
  share: Share
  file: ShareFile
}>()

const emit = defineEmits([
  'download'
])

const isLoadingType = ref(false)

const isLoadingContent = ref(false)

const previewType = ref<FileType>('unknown')

const previewText = ref('')

watchEffect(() => {
  isLoadingType.value = true
  useAxiosInstance().get<{
    type: FileType
  }>(getContentTypeUrl(props.share, props.file.id))
    .then(({ data }) => {
      previewType.value = data.type
    })
    .finally(() => {
      isLoadingType.value = false
    })
})

watch(previewType, () => {
  if (previewType.value === 'text') {
    previewText.value = ''
    isLoadingContent.value = true
    useAxiosInstance().get<string>(getContentUrl(props.share, props.file.id), {
      responseType: 'text'
    }).then(({ data }) => {
      previewText.value = data
    }).finally(() => {
      isLoadingContent.value = false
    })
  }
}, {
  immediate: true
})
</script>

<template>
  <div
    v-if="previewType === 'image'"
    class="w-full aspect-[3/4] sm:aspect-[4/3] flex items-start justify-center"
  >
    <img
      class="max-w-full max-h-full"
      alt="image preview"
      :src="getPreviewUrl(share, file.id)"
      @error="previewType = 'unknown'"
    />
  </div>
  <div
    v-else-if="previewType === 'text' && !isLoadingContent"
    class="w-full py-4 text-gray-600 dark:text-neutral-300 text-xs font-mono whitespace-pre-wrap overflow-x-auto"
    v-text="previewText"
  />
  <div
    v-else-if="isLoadingType || isLoadingContent"
    class="py-4 sm:py-32 flex items-center justify-center gap-2"
  >
    <flowbite-spinner class="w-6 h-6" />
    <p class="text-lg text-gray-500 dark:text-neutral-400">
      Loading preview...
    </p>
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
